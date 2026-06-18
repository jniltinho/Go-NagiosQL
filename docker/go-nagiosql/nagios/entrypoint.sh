#!/bin/bash
set -e

NAGIOS_ETC=/etc/nagios4
NAGIOSQL_ETC=/etc/nagiosql
NAGIOS_VAR=/var/lib/nagios4
NAGIOS_LIB=/usr/lib/nagios/plugins

DB_HOST="${NAGIOSQL_DATABASE_HOST:-${DB_HOST:-db}}"
DB_PORT="${NAGIOSQL_DATABASE_PORT:-${DB_PORT:-3306}}"
DB_USER="${NAGIOSQL_DATABASE_USER:-nagiosql}"
DB_PASS="${NAGIOSQL_DATABASE_PASSWORD:-nagiosqlpass}"
NAGIOSQL_USER="${NAGIOSQL_USER:-admin}"
NAGIOSQL_PASSWORD="${NAGIOSQL_PASSWORD:-admin}"

# ════════════════════════════════════════════════════════════════
#  NAGIOS CORE — inicialização dos volumes
# ════════════════════════════════════════════════════════════════

echo "==> [nagios] Inicializando volumes..."

if [ ! -f "$NAGIOS_ETC/nagios.cfg" ]; then
    echo "    Copiando configurações padrão para $NAGIOS_ETC..."
    cp -r /etc/nagios4.default/. "$NAGIOS_ETC/"
fi
mkdir -p "$NAGIOS_ETC/stylesheets"

chown -R root:nagios "$NAGIOS_ETC"
find "$NAGIOS_ETC" -type d -exec chmod 755 {} \;
find "$NAGIOS_ETC" -type f -exec chmod 644 {} \;

if [ -z "$(ls -A "$NAGIOSQL_ETC" 2>/dev/null)" ]; then
    echo "    Copiando configs de exemplo para $NAGIOSQL_ETC..."
    cp -rn /etc/nagiosql.default/. "$NAGIOSQL_ETC/"
fi

if [ -z "$(ls -A "$NAGIOS_LIB" 2>/dev/null)" ]; then
    echo "    Copiando plugins para $NAGIOS_LIB..."
    cp -a /usr/lib/nagios/plugins.default/. "$NAGIOS_LIB/"
fi

mkdir -p /run/nagios4
chown nagios:nagios /run/nagios4

mkdir -p "$NAGIOS_VAR/rw" "$NAGIOS_VAR/spool/checkresults"
chown -R nagios:nagios "$NAGIOS_VAR"
# rw/: SGID so files inherit www-data group; reload.trigger writable by www-data
chown nagios:www-data "$NAGIOS_VAR/rw"
chmod 2775 "$NAGIOS_VAR/rw"
touch "$NAGIOS_VAR/rw/reload.trigger"
chown nagios:www-data "$NAGIOS_VAR/rw/reload.trigger"
chmod 660 "$NAGIOS_VAR/rw/reload.trigger"
chown nagios:www-data "$NAGIOS_VAR/spool/checkresults"
chmod 775 "$NAGIOS_VAR/spool/checkresults"

mkdir -p /var/log/nagios4 /var/cache/nagios4
chown nagios:nagios /var/log/nagios4 /var/cache/nagios4

HTPASSWD="$NAGIOS_ETC/htpasswd.users"
if [ ! -f "$HTPASSWD" ]; then
    echo "    Criando usuário nagiosadmin..."
    PASS="${NAGIOS_ADMIN_PASSWORD:-nagiosadmin}"
    printf "nagiosadmin:%s\n" "$(openssl passwd -apr1 "$PASS")" > "$HTPASSWD"
fi
chown nagios:www-data "$HTPASSWD"
chmod 640 "$HTPASSWD"

# nagios.cfg + cgi.cfg: www-data precisa ler/escrever
chown www-data:nagios "$NAGIOS_ETC/nagios.cfg" "$NAGIOS_ETC/cgi.cfg" 2>/dev/null || true
chmod 640 "$NAGIOS_ETC/nagios.cfg" "$NAGIOS_ETC/cgi.cfg" 2>/dev/null || true
# Autenticação feita pelo Apache Basic Auth; evitar double-auth
sed -i 's/^use_authentication=1/use_authentication=0/' "$NAGIOS_ETC/cgi.cfg" 2>/dev/null || true

# /etc/nagiosql: dirs e permissões
mkdir -p \
    "$NAGIOSQL_ETC/hosts" \
    "$NAGIOSQL_ETC/services" \
    "$NAGIOSQL_ETC/backup/hosts" \
    "$NAGIOSQL_ETC/backup/services"

for cfg in timeperiods commands contacts contactgroups contacttemplates \
           hosttemplates hostgroups hostextinfo hostescalations hostdependencies \
           servicetemplates servicegroups serviceextinfo serviceescalations servicedependencies; do
    [ -f "$NAGIOSQL_ETC/${cfg}.cfg" ] || touch "$NAGIOSQL_ETC/${cfg}.cfg"
done
chown -R www-data:nagios "$NAGIOSQL_ETC"
find "$NAGIOSQL_ETC" -type d -exec chmod 750 {} \;
find "$NAGIOSQL_ETC" -type f -exec chmod 640 {} \;

# ── Comentar arquivos padrão Debian que conflitam com nagiosql ──────────────
NAGIOS_CFG="$NAGIOS_ETC/nagios.cfg"
sed -i \
    -e 's|^cfg_dir=/etc/nagios-plugins/config|#cfg_dir=/etc/nagios-plugins/config|' \
    -e 's|^cfg_file=/etc/nagios4/objects/commands\.cfg|#cfg_file=/etc/nagios4/objects/commands.cfg|' \
    -e 's|^cfg_file=/etc/nagios4/objects/timeperiods\.cfg|#cfg_file=/etc/nagios4/objects/timeperiods.cfg|' \
    -e 's|^cfg_file=/etc/nagios4/objects/templates\.cfg|#cfg_file=/etc/nagios4/objects/templates.cfg|' \
    -e 's|^cfg_file=/etc/nagios4/objects/contacts\.cfg|#cfg_file=/etc/nagios4/objects/contacts.cfg|' \
    -e 's|^cfg_file=/etc/nagios4/objects/localhost\.cfg|#cfg_file=/etc/nagios4/objects/localhost.cfg|' \
    "$NAGIOS_CFG"

# ── Adicionar entradas cfg_dir/cfg_file ao nagios.cfg (idempotente) ─────────
if ! grep -q "etc/nagiosql/hosts" "$NAGIOS_CFG" 2>/dev/null; then
    echo "    Adicionando entradas nagiosql ao nagios.cfg..."
    cat >> "$NAGIOS_CFG" << 'EOF'

# Configurações gerenciadas pelo nagiosql (/etc/nagiosql/)
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
#  nagiosql — aguardar banco e inicializar (GORM migrate)
# ════════════════════════════════════════════════════════════════

echo "==> [nagiosql] Aguardando MariaDB em $DB_HOST:$DB_PORT..."
until mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" \
    --ssl=0 -e "SELECT 1" 2>/dev/null; do
    echo "    Waiting for MariaDB..."
    sleep 2
done
echo "    MariaDB pronto."

# migrate é idempotente — seguro rodar a cada boot
echo "==> [nagiosql] Rodando migrate..."
/opt/nagiosql/nagiosql migrate \
    --config /opt/nagiosql/config.toml \
    --admin-user  "$NAGIOSQL_USER" \
    --admin-password "$NAGIOSQL_PASSWORD" \
    --sample
echo "    Migrate concluído."

# Gerar todos os .cfg a partir do banco (idempotente)
echo "==> [nagiosql] Gerando arquivos .cfg do banco..."
/opt/nagiosql/nagiosql config write all \
    --config /opt/nagiosql/config.toml
echo "    Config write concluído."

# ════════════════════════════════════════════════════════════════
#  Validar e agendar reload
# ════════════════════════════════════════════════════════════════

# Valida nagios.cfg e escreve reload.trigger — reload-watcher envia SIGHUP
# assim que o nagios4 estiver rodando.
echo "==> [nagios] Validando e agendando reload..."
/opt/nagiosql/nagiosql config restart \
    --config /opt/nagiosql/config.toml
echo "    Reload agendado (reload.trigger escrito)."

echo "==> Iniciando supervisord..."
exec /usr/bin/supervisord -n -c /etc/supervisor/supervisord.conf
