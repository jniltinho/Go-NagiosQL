#!/bin/bash
set -e

# PHP: timezone e include_path para NagiosQL (mod_php, Apache 8.4)
{
    [ -n "${TZ:-}" ] && printf '[Date]\ndate.timezone = %s\n' "$TZ"
    # NagiosQL usa PEAR bundled; garante que PEAR.php seja encontrado
    # mesmo se php-pear não estiver instalado no sistema
    printf 'include_path = ".:/usr/share/php:/var/www/nagiosql/libraries/pear"\n'
} > /etc/php/8.4/apache2/conf.d/99-timezone.ini

NAGIOS_ETC=/etc/nagios4        # config do daemon (pacote Debian)
NAGIOSQL_ETC=/etc/nagiosql     # configs gerados pelo NagiosQL (PDF §1.2)
NAGIOS_VAR=/var/lib/nagios4    # runtime: status, cmd pipe, spool
NAGIOS_LIB=/usr/lib/nagios/plugins

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

# /etc/nagios4: config principal do daemon (defaults do pacote Debian)
if [ ! -f "$NAGIOS_ETC/nagios.cfg" ]; then
    echo "    Copiando configurações padrão para $NAGIOS_ETC..."
    cp -r /etc/nagios4.default/. "$NAGIOS_ETC/"
fi
# Necessário para o Alias /nagios4/stylesheets no Apache (nagios4-cgi.conf)
mkdir -p "$NAGIOS_ETC/stylesheets"

# Permissões: root:nagios, dirs 755, arquivos 644 (padrão do pacote Debian).
# 644: www-data (NagiosQL check) e nagios (daemon) podem ler todos os arquivos.
# htpasswd.users tem ownership específico definido abaixo.
chown -R root:nagios "$NAGIOS_ETC"
find "$NAGIOS_ETC" -type d -exec chmod 755 {} \;
find "$NAGIOS_ETC" -type f -exec chmod 644 {} \;

# /etc/nagiosql: arquivos .cfg gerados pelo NagiosQL
if [ -z "$(ls -A "$NAGIOSQL_ETC" 2>/dev/null)" ]; then
    echo "    Copiando configs de exemplo para $NAGIOSQL_ETC..."
    cp -rn /etc/nagiosql.default/. "$NAGIOSQL_ETC/"
fi

# /usr/lib/nagios/plugins: plugins do monitoramento
if [ -z "$(ls -A "$NAGIOS_LIB" 2>/dev/null)" ]; then
    echo "    Copiando plugins para $NAGIOS_LIB..."
    cp -a /usr/lib/nagios/plugins.default/. "$NAGIOS_LIB/"
fi

# PID dir: não criado pelo systemd-tmpfiles em containers Docker
mkdir -p /run/nagios4
chown nagios:nagios /run/nagios4

# Diretórios de runtime no volume var
mkdir -p "$NAGIOS_VAR/rw" "$NAGIOS_VAR/spool/checkresults"
chown -R nagios:nagios "$NAGIOS_VAR"
# rw/: grupo www-data com SGID — Apache (www-data) escreve nagios.cmd e reload.trigger
chown nagios:www-data "$NAGIOS_VAR/rw"
chmod 2775 "$NAGIOS_VAR/rw"
# reload.trigger: NagiosQL verifica file_exists() antes de escrever — manter sempre presente
touch "$NAGIOS_VAR/rw/reload.trigger"
chown nagios:www-data "$NAGIOS_VAR/rw/reload.trigger"
chmod 660 "$NAGIOS_VAR/rw/reload.trigger"
# checkresults: www-data precisa escrever (NagiosQL roda nagios4 -v como www-data)
chown nagios:www-data "$NAGIOS_VAR/spool/checkresults"
chmod 775  "$NAGIOS_VAR/spool/checkresults"

mkdir -p /var/log/nagios4 /var/cache/nagios4
chown nagios:nagios /var/log/nagios4 /var/cache/nagios4

# htpasswd para Basic Auth do Nagios Core no Apache
HTPASSWD="$NAGIOS_ETC/htpasswd.users"
if [ ! -f "$HTPASSWD" ]; then
    echo "    Criando usuário nagiosadmin..."
    PASS="${NAGIOS_ADMIN_PASSWORD:-nagiosadmin}"
    printf "nagiosadmin:%s\n" "$(openssl passwd -apr1 "$PASS")" > "$HTPASSWD"
fi
# Sempre corrigir: o chown -R root:nagios acima sobrescreve esta permissão.
# www-data precisa ler o htpasswd para autenticar requisições HTTP.
chown nagios:www-data "$HTPASSWD"
chmod 640 "$HTPASSWD"

# nagios.cfg e cgi.cfg: www-data precisa de write (PDF §1.3)
chown www-data:nagios "$NAGIOS_ETC/nagios.cfg" "$NAGIOS_ETC/cgi.cfg" 2>/dev/null || true
chmod 640 "$NAGIOS_ETC/nagios.cfg" "$NAGIOS_ETC/cgi.cfg" 2>/dev/null || true

# use_authentication=0: autenticação feita pelo Apache (Basic Auth).
# O padrão Debian é 0, mas o volume pode ter sido editado — garantir.
sed -i 's/^use_authentication=1/use_authentication=0/' "$NAGIOS_ETC/cgi.cfg" 2>/dev/null || true

# ── /etc/nagiosql: estrutura de dirs e permissões ───────────────────
# PDF §1.2: chown -R wwwrun.nagios /etc/nagiosql (www-data.nagios em Debian)
# PDF §1.2: dirs 750, arquivos 640
mkdir -p \
    "$NAGIOSQL_ETC/hosts" \
    "$NAGIOSQL_ETC/services" \
    "$NAGIOSQL_ETC/backup/hosts" \
    "$NAGIOSQL_ETC/backup/services"

for cfg in timeperiods contacts contactgroups contacttemplates \
           hosttemplates hostgroups hostextinfo hostescalations hostdependencies \
           servicetemplates servicegroups serviceextinfo serviceescalations servicedependencies; do
    [ -f "$NAGIOSQL_ETC/${cfg}.cfg" ] || touch "$NAGIOSQL_ETC/${cfg}.cfg"
done
# commands.cfg sempre vazio no início — NagiosQL gera o conteúdo via UI.
# Um arquivo com definições duplicadas (já existentes no nagios4/objects/)
# causa parse error que impede o daemon de iniciar.
[ -f "$NAGIOSQL_ETC/commands.cfg" ] || touch "$NAGIOSQL_ETC/commands.cfg"

chown -R www-data:nagios "$NAGIOSQL_ETC"
find "$NAGIOSQL_ETC" -type d -exec chmod 750 {} \;
find "$NAGIOSQL_ETC" -type f -exec chmod 640 {} \;

# ── Remover entradas padrão que conflitam com NagiosQL ──────────────
# NagiosQL gera versões completas de timeperiods, templates, contacts e
# localhost — os arquivos padrão do pacote Debian causam duplicate errors.
NAGIOS_CFG="$NAGIOS_ETC/nagios.cfg"
sed -i \
    -e 's|^cfg_dir=/etc/nagios-plugins/config|#cfg_dir=/etc/nagios-plugins/config|' \
    -e 's|^cfg_file=/etc/nagios4/objects/commands\.cfg|#cfg_file=/etc/nagios4/objects/commands.cfg|' \
    -e 's|^cfg_file=/etc/nagios4/objects/timeperiods\.cfg|#cfg_file=/etc/nagios4/objects/timeperiods.cfg|' \
    -e 's|^cfg_file=/etc/nagios4/objects/templates\.cfg|#cfg_file=/etc/nagios4/objects/templates.cfg|' \
    -e 's|^cfg_file=/etc/nagios4/objects/contacts\.cfg|#cfg_file=/etc/nagios4/objects/contacts.cfg|' \
    -e 's|^cfg_file=/etc/nagios4/objects/localhost\.cfg|#cfg_file=/etc/nagios4/objects/localhost.cfg|' \
    "$NAGIOS_CFG"

# ── Adicionar entradas cfg_dir/cfg_file ao nagios.cfg ─────────────
if ! grep -q "etc/nagiosql/hosts" "$NAGIOS_CFG" 2>/dev/null; then
    echo "    Adicionando entradas NagiosQL ao nagios.cfg..."
    cat >> "$NAGIOS_CFG" << 'EOF'

# Configurações gerenciadas pelo NagiosQL (/etc/nagiosql/)
cfg_file=/etc/nagiosql/timeperiods.cfg
cfg_file=/etc/nagiosql/commands.cfg
cfg_file=/etc/nagiosql/contacts.cfg
cfg_file=/etc/nagiosql/contactgroups.cfg
cfg_file=/etc/nagiosql/contacttemplates.cfg
cfg_dir=/etc/nagiosql/hosts
cfg_file=/etc/nagiosql/hosttemplates.cfg
cfg_file=/etc/nagiosql/hostgroups.cfg
cfg_file=/etc/nagiosql/hostextinfo.cfg
cfg_file=/etc/nagiosql/hostescalations.cfg
cfg_file=/etc/nagiosql/hostdependencies.cfg
cfg_dir=/etc/nagiosql/services
cfg_file=/etc/nagiosql/servicetemplates.cfg
cfg_file=/etc/nagiosql/servicegroups.cfg
cfg_file=/etc/nagiosql/serviceextinfo.cfg
cfg_file=/etc/nagiosql/serviceescalations.cfg
cfg_file=/etc/nagiosql/servicedependencies.cfg
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

    echo "    Importando dados de exemplo..."
    $MYSQL --init-command="SET SESSION sql_mode='NO_ENGINE_SUBSTITUTION'" \
        < /opt/nagiosql/import_nagios_sample.sql
    echo "    Dados de exemplo importados."

    echo "    Ajustando caminhos para pacote Debian nagios4 + NagiosQL 3.5..."
    $MYSQL << SQL
UPDATE \`tbl_configtarget\` SET
    \`basedir\`       = '/etc/nagiosql/',
    \`hostconfig\`    = '/etc/nagiosql/hosts/',
    \`serviceconfig\` = '/etc/nagiosql/services/',
    \`backupdir\`     = '/etc/nagiosql/backup/',
    \`hostbackup\`    = '/etc/nagiosql/backup/hosts/',
    \`servicebackup\` = '/etc/nagiosql/backup/services/',
    \`nagiosbasedir\` = '/etc/nagios4/',
    \`importdir\`     = '/etc/nagios4/conf.d/',
    \`picturedir\`    = '/usr/share/nagios4/htdocs/images/logos/',
    \`commandfile\`   = '/var/lib/nagios4/rw/reload.trigger',
    \`binaryfile\`    = '/usr/sbin/nagios4',
    \`pidfile\`       = '/run/nagios4/nagios4.pid',
    \`conffile\`      = '/etc/nagios4/nagios.cfg',
    \`cgifile\`       = '/etc/nagios4/cgi.cfg',
    \`resourcefile\`  = '/etc/nagios4/resource.cfg',
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

    # check_dns não vem no import padrão mas é usado pelos hosts de amostra
    $MYSQL << 'SQL'
INSERT IGNORE INTO `tbl_command`
    (`command_name`, `command_line`, `command_type`, `register`, `active`, `last_modified`, `access_group`, `config_id`)
VALUES
    ('check_dns', '$USER1$/check_dns -H www.google.com -s $HOSTADDRESS$ $ARG1$', 0, '1', '1', NOW(), 0, 0);
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
    "/usr/sbin/nagios4 -v $NAGIOS_CFG" 2>&1 | grep -q "^Total Errors:   0"; then
    echo "    Configuração OK."
else
    echo "    Aviso: erros encontrados — verifique os logs."
fi

echo "==> Iniciando supervisord..."
exec /usr/bin/supervisord -n -c /etc/supervisor/supervisord.conf
