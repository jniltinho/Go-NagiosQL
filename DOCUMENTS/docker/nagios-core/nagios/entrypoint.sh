#!/bin/bash
set -e

# PHP timezone — segue TZ do container (definido no docker-compose.yml)
if [ -n "${TZ:-}" ]; then
    printf '[Date]\ndate.timezone = %s\n' "$TZ" \
        | tee /etc/php/8.4/fpm/conf.d/99-timezone.ini \
              /etc/php/8.4/cli/conf.d/99-timezone.ini > /dev/null
fi

NAGIOS_ETC=/usr/local/nagios/etc
NAGIOS_VAR=/usr/local/nagios/var
NAGIOS_LIB=/usr/local/nagios/libexec

DB_HOST="${DB_HOST:-db}"
DB_PORT="${DB_PORT:-3306}"
DB_NAME="${DB_NAME:-nagiosql}"
DB_USER="${DB_USER:-nagiosql}"
DB_PASSWORD="${DB_PASSWORD:-nagiosqlpass}"
NAGIOSQL_USER="${NAGIOSQL_USER:-admin}"
NAGIOSQL_PASSWORD="${NAGIOSQL_PASSWORD:-admin}"

MYSQL="mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASSWORD --skip-ssl $DB_NAME"
SETTINGS=/var/www/nagiosql/config/settings.php

# ════════════════════════════════════════════════════════════════
#  NAGIOS CORE — inicialização dos volumes
# ════════════════════════════════════════════════════════════════

echo "==> [nagios] Inicializando volumes..."

if [ ! -f "$NAGIOS_ETC/nagios.cfg" ]; then
    echo "    Copiando configurações padrão para $NAGIOS_ETC..."
    cp -r /usr/local/nagios/etc.default/. "$NAGIOS_ETC/"
    chown -R nagios:nagios "$NAGIOS_ETC"
    chmod -R 775 "$NAGIOS_ETC"
fi

# Nagios 4.x default: lock_file=/run/nagios.lock (diretório root-owned, nagios não pode gravar).
# Corrige para path no var/ que o nagios user controla.
sed -i 's|^lock_file=.*|lock_file=/usr/local/nagios/var/nagios.lock|' "$NAGIOS_ETC/nagios.cfg"

if [ -z "$(ls -A "$NAGIOS_LIB" 2>/dev/null)" ]; then
    echo "    Copiando plugins para $NAGIOS_LIB..."
    cp -a /usr/local/nagios/libexec.default/. "$NAGIOS_LIB/"
fi

mkdir -p \
    "$NAGIOS_VAR/rw" \
    "$NAGIOS_VAR/spool/checkresults" \
    "$NAGIOS_VAR/archives"
chown -R nagios:nagios "$NAGIOS_VAR"
# SGID em var/rw/: novos FIFOs (nagios.cmd) herdam grupo nagioscfg em vez do GID primário do processo
chgrp nagioscfg "$NAGIOS_VAR/rw"
chmod 775 "$NAGIOS_VAR/rw"
chmod g+s  "$NAGIOS_VAR/rw"
chmod 775 "$NAGIOS_VAR/spool/checkresults"

HTPASSWD="$NAGIOS_ETC/htpasswd.users"
if [ ! -f "$HTPASSWD" ]; then
    echo "    Criando usuário nagiosadmin..."
    PASS="${NAGIOS_ADMIN_PASSWORD:-nagiosadmin}"
    printf "nagiosadmin:%s\n" "$(openssl passwd -apr1 "$PASS")" > "$HTPASSWD"
    chown nagios:nagios "$HTPASSWD"
fi

mkdir -p \
    "$NAGIOS_ETC/nagiosql/hosts" \
    "$NAGIOS_ETC/nagiosql/services" \
    "$NAGIOS_ETC/nagiosql/backup/hosts" \
    "$NAGIOS_ETC/nagiosql/backup/services" \
    "$NAGIOS_ETC/import"

for cfg in timeperiods commands contacts contactgroups contacttemplates \
           hosttemplates hostgroups hostextinfo hostescalations hostdependencies \
           servicetemplates servicegroups serviceextinfo serviceescalations servicedependencies; do
    PLACEHOLDER="$NAGIOS_ETC/nagiosql/${cfg}.cfg"
    [ -f "$PLACEHOLDER" ] || touch "$PLACEHOLDER"
done

# nagioscfg: grupo compartilhado nagios+www-data para escrita nos dirs de config
chown -R nagios:nagioscfg "$NAGIOS_ETC/nagiosql" "$NAGIOS_ETC/import" 2>/dev/null || true
chmod -R 775 "$NAGIOS_ETC/nagiosql" "$NAGIOS_ETC/import"

NAGIOS_CFG="$NAGIOS_ETC/nagios.cfg"
if ! grep -q "nagiosql/hosts" "$NAGIOS_CFG" 2>/dev/null; then
    echo "    Adicionando entradas NagiosQL ao nagios.cfg..."
    cat >> "$NAGIOS_CFG" << 'EOF'

# Configurações gerenciadas pelo NagiosQL
cfg_file=/usr/local/nagios/etc/nagiosql/timeperiods.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/commands.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/contacts.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/contactgroups.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/contacttemplates.cfg
cfg_dir=/usr/local/nagios/etc/nagiosql/hosts
cfg_file=/usr/local/nagios/etc/nagiosql/hosttemplates.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/hostgroups.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/hostextinfo.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/hostescalations.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/hostdependencies.cfg
cfg_dir=/usr/local/nagios/etc/nagiosql/services
cfg_file=/usr/local/nagios/etc/nagiosql/servicetemplates.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/servicegroups.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/serviceextinfo.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/serviceescalations.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/servicedependencies.cfg
EOF
fi

# ════════════════════════════════════════════════════════════════
#  NAGIOSQL — aguardar banco e inicializar
# ════════════════════════════════════════════════════════════════

echo "==> [nagiosql] Aguardando MariaDB em $DB_HOST:$DB_PORT..."
until nc -z "$DB_HOST" "$DB_PORT" 2>/dev/null; do
    echo "    Waiting for MariaDB..."
    sleep 2
done
echo "    MariaDB pronto."

TABLE_COUNT=$($MYSQL -sN -e \
    "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='$DB_NAME';" \
    2>/dev/null || echo "0")

if [ "$TABLE_COUNT" -eq 0 ]; then
    echo "==> [nagiosql] Importando schema..."
    $MYSQL --init-command="SET SESSION sql_mode='NO_ENGINE_SUBSTITUTION'" \
        < /opt/nagiosql/nagiosQL_v35_db_mysql.sql
    echo "    Schema importado."

    echo "    Importando dados de exemplo (templates, comandos, timeperiods)..."
    # Contém: 24 comandos, 5 timeperiods, templates generic-host/linux-server/
    # generic-service/local-service, contactgroup admins, 4 hostgroups,
    # contact nagiosadmin e 4 hosts de exemplo com 21 serviços.
    # Nota: o host 'localhost' do sample pode conflitar com objects/localhost.cfg
    # se o usuário gerar os .cfg sem remover um dos dois primeiro.
    $MYSQL --init-command="SET SESSION sql_mode='NO_ENGINE_SUBSTITUTION'" \
        < /opt/nagiosql/import_nagios_sample.sql
    echo "    Dados de exemplo importados."

    echo "    Ajustando caminhos do Nagios Core..."
    $MYSQL << SQL
UPDATE \`tbl_configtarget\` SET
    \`basedir\`       = '/usr/local/nagios/etc/nagiosql/',
    \`hostconfig\`    = '/usr/local/nagios/etc/nagiosql/hosts/',
    \`serviceconfig\` = '/usr/local/nagios/etc/nagiosql/services/',
    \`backupdir\`     = '/usr/local/nagios/etc/nagiosql/backup/',
    \`hostbackup\`    = '/usr/local/nagios/etc/nagiosql/backup/hosts/',
    \`servicebackup\` = '/usr/local/nagios/etc/nagiosql/backup/services/',
    \`nagiosbasedir\` = '/usr/local/nagios/etc/',
    \`importdir\`     = '/usr/local/nagios/etc/import/',
    \`picturedir\`    = '/usr/local/nagios/share/images/logos/',
    \`commandfile\`   = '/usr/local/nagios/var/reload.trigger',
    \`binaryfile\`    = '/usr/local/nagios/bin/nagios',
    \`pidfile\`       = '/usr/local/nagios/var/nagios.lock',
    \`conffile\`      = '/usr/local/nagios/etc/nagios.cfg',
    \`cgifile\`       = '/usr/local/nagios/etc/cgi.cfg',
    \`resourcefile\`  = '/usr/local/nagios/etc/resource.cfg',
    \`version\`       = 4
WHERE \`target\` = 'localhost';
SQL

    echo "    Criando usuário admin no NagiosQL..."
    $MYSQL << SQL
INSERT INTO \`tbl_user\`
    (\`id\`, \`username\`, \`alias\`, \`password\`,
     \`admin_enable\`, \`wsauth\`, \`active\`, \`nodelete\`,
     \`language\`, \`domain\`, \`last_login\`, \`last_modified\`)
VALUES
    (1, '$NAGIOSQL_USER', 'Administrator', MD5('$NAGIOSQL_PASSWORD'),
     '1', '0', '1', '1', '1', '1', '2000-01-01 00:00:00', NOW())
ON DUPLICATE KEY UPDATE
    \`password\` = MD5('$NAGIOSQL_PASSWORD');
SQL

    $MYSQL << SQL
INSERT INTO \`tbl_settings\` (\`category\`, \`name\`, \`value\`)
VALUES ('db', 'version', '3.5.0')
ON DUPLICATE KEY UPDATE \`value\` = '3.5.0';
SQL
    echo "    Configuração inicial concluída."
fi

CONFIG_DIR=$(dirname "$SETTINGS")
if [ ! -f "$CONFIG_DIR/fieldvars.php" ]; then
    echo "==> [nagiosql] Inicializando config/ a partir do config.default..."
    cp -rn /var/www/nagiosql/config.default/. "$CONFIG_DIR/"
    chown -R www-data:www-data "$CONFIG_DIR"
    chmod -R 775 "$CONFIG_DIR"
fi

if [ ! -f "$SETTINGS" ]; then
    echo "==> [nagiosql] Gerando config/settings.php..."
    cat > "$SETTINGS" << EOF
<?php
exit;
?>
;///////////////////////////////////////////////////////////////////////////////
;
; NagiosQL
;
;///////////////////////////////////////////////////////////////////////////////
[db]
type            = mysqli
server          = ${DB_HOST}
port            = ${DB_PORT}
database        = ${DB_NAME}
username        = ${DB_USER}
password        = ${DB_PASSWORD}
[path]
protocol        = http
tempdir         = /tmp
base_url        = /
base_path       = /var/www/nagiosql/
[data]
locale          = en_GB
encoding        = utf-8
[security]
logofftime      = 3600
wsauth          = 0
[common]
pagelines       = 15
seldisable      = 1
tplcheck        = 0
updcheck        = 0
[network]
proxy           = 0
proxyserver     =
proxyuser       =
proxypasswd     =
onlineupdate    = 0
[performance]
parents         = 1
EOF
    chown www-data:www-data "$SETTINGS"
    echo "    settings.php gerado."
fi

# ════════════════════════════════════════════════════════════════
#  Validar e iniciar
# ════════════════════════════════════════════════════════════════

echo "==> [nagios] Validando nagios.cfg..."
if su -s /bin/bash nagios -c \
    "/usr/local/nagios/bin/nagios -v $NAGIOS_CFG" 2>&1 | grep -q "^Total Errors:   0"; then
    echo "    Configuração OK."
else
    echo "    Aviso: erros encontrados — verifique os logs."
fi

echo "==> Iniciando supervisord..."
mkdir -p /run/php
exec /usr/bin/supervisord -n -c /etc/supervisor/supervisord.conf
