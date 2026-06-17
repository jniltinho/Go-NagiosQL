#!/bin/bash
# build-debs.sh — Compila e empacota Nagios Core, Nagios Plugins e NagiosQL
# como pacotes .deb seguindo o Debian Policy Manual:
#   §5.2  — formato do arquivo de controle
#   §5.6  — campos obrigatórios e recomendados do DEBIAN/control
#   §6.1  — estrutura dos maintainer scripts (postinst/prerm)
#   §12.5 — arquivo copyright obrigatório
#   §12.7 — changelog.Debian.gz obrigatório
set -euo pipefail

NAGIOS_VERSION="${NAGIOS_VERSION:-4.5.13}"
PLUGINS_VERSION="${PLUGINS_VERSION:-2.5}"
NAGIOSQL_VERSION="${NAGIOSQL_VERSION:-3.5.0}"
NAGIOSQL_SRC="${NAGIOSQL_SRC:-/nagiosql-src}"
MAINTAINER="${MAINTAINER:-Nagios Docker Build <build@local>}"

BUILD_DIR="/tmp/build"
OUTPUT_DIR="/tmp/debs"
mkdir -p "$BUILD_DIR" "$OUTPUT_DIR"

# ─────────────────────────────────────────────────────────────────────────────
# Utilitários de packaging (Policy §5 e §12)
# ─────────────────────────────────────────────────────────────────────────────

# Calcula Installed-Size em KiB (1024 bytes) — Policy §5.6.20
installed_size_kb() {
    du -sk --apparent-size "$1" | cut -f1
}

# Gera DEBIAN/md5sums de todos os arquivos fora de DEBIAN/ — Policy §5.6.31
gen_md5sums() {
    local staging="$1"
    (cd "$staging" && \
        find . -type f ! -path './DEBIAN/*' | LC_ALL=C sort | \
        xargs md5sum | sed 's|^\./||') > "$staging/DEBIAN/md5sums"
    chmod 644 "$staging/DEBIAN/md5sums"
}

# Instala um maintainer script com permissão 755 — Policy §6.1
# Uso: install_script <caminho> << 'HEREDOC' ... HEREDOC
install_script() {
    local file="$1"
    cat > "$file"
    chmod 755 "$file"
}

# Cria usr/share/doc/<pkg>/copyright mínimo — Policy §12.5
gen_copyright() {
    local staging="$1" pkg="$2" upstream_url="$3" license="$4"
    mkdir -p "$staging/usr/share/doc/$pkg"
    cat > "$staging/usr/share/doc/$pkg/copyright" << EOF
Format: https://www.debian.org/doc/packaging-manuals/copyright-format/1.0/
Upstream-Contact: ${upstream_url}

Files: *
License: ${license}
EOF
    chmod 644 "$staging/usr/share/doc/$pkg/copyright"
}

# Cria usr/share/doc/<pkg>/changelog.Debian.gz mínimo — Policy §12.7
gen_changelog() {
    local staging="$1" pkg="$2" ver="$3"
    mkdir -p "$staging/usr/share/doc/$pkg"
    printf '%s (%s) unstable; urgency=low\n\n  * Package built for Docker deployment.\n\n -- %s  %s\n' \
        "$pkg" "$ver" "$MAINTAINER" "$(date -R)" \
        | gzip -9 > "$staging/usr/share/doc/$pkg/changelog.Debian.gz"
    chmod 644 "$staging/usr/share/doc/$pkg/changelog.Debian.gz"
}

# ═════════════════════════════════════════════════════════════════════════════
# Nagios Core
# ═════════════════════════════════════════════════════════════════════════════
build_nagios_core() {
    local name="nagios-${NAGIOS_VERSION}"
    local src="$BUILD_DIR/$name"
    local staging="$BUILD_DIR/pkg-nagios"
    local deb="$OUTPUT_DIR/nagios-core_${NAGIOS_VERSION}_amd64.deb"

    echo "==> [1/3] Nagios Core ${NAGIOS_VERSION}"

    cd "$BUILD_DIR"
    wget -q "https://github.com/NagiosEnterprises/nagioscore/releases/download/nagios-${NAGIOS_VERSION}/${name}.tar.gz"
    tar xzf "${name}.tar.gz" && cd "$src"

    echo "    configure..."
    # Ref: NAGIOS_BUILD_COMMANDS.md — usamos --prefix padrão (/usr/local/nagios) em vez
    # de --prefix=/usr porque NagiosQL é hardcoded para /usr/local/nagios/.
    # --with-command-group define o grupo autorizado a escrever no pipe de comandos externos.
    # Usamos nagioscfg (grupo compartilhado) em vez de nagios para limitar o acesso
    # de www-data apenas às operações necessárias. Ref: NAGIOSQL_DEBIAN_PACKAGING.md §File Permissions
    ./configure \
        --with-httpd-conf=/etc/nginx/sites-enabled \
        --with-command-group=nagioscfg \
        >/dev/null

    echo "    make -j$(nproc)..."
    make -j"$(nproc)" all >/dev/null

    echo "    make install..."
    # install-commandmode: define ownership/permissões de var/rw (pipe de comandos externos)
    make DESTDIR="$staging" install install-config install-commandmode >/dev/null

    echo "    packaging..."
    mkdir -p "$staging/DEBIAN"

    # control — campos obrigatórios + recomendados (Policy §5.6)
    # Installed-Size calculado após o make install para ser preciso
    cat > "$staging/DEBIAN/control" << EOF
Package: nagios-core
Version: ${NAGIOS_VERSION}
Architecture: amd64
Section: net
Priority: optional
Installed-Size: $(installed_size_kb "$staging")
Maintainer: ${MAINTAINER}
Homepage: https://www.nagios.org
Depends: libgd3, libssl3, libc6 (>= 2.17)
Description: open-source host, service and network monitoring daemon
 Nagios Core is a host/service/network monitoring application. It
 monitors specified hosts and services, alerting you when problems occur
 and when they are resolved.
 .
 This package includes the Nagios daemon, CGIs and web interface,
 compiled from source for Debian trixie (amd64).
EOF
    chmod 644 "$staging/DEBIAN/control"

    # postinst — Policy §6.1: deve usar "case $1 in configure)"
    install_script "$staging/DEBIAN/postinst" << 'POSTINST'
#!/bin/bash
set -e
case "$1" in
    configure)
        getent group nagios >/dev/null 2>&1 \
            || groupadd -g 3000 nagios
        # nagioscfg: grupo compartilhado para escrita nos dirs de config gerados pelo NagiosQL
        getent group nagioscfg >/dev/null 2>&1 \
            || groupadd -g 3001 nagioscfg
        id -u nagios >/dev/null 2>&1 \
            || useradd -u 3000 -g nagios -d /usr/local/nagios \
                       -s /bin/bash -c "Nagios monitoring daemon" nagios
        usermod -aG nagioscfg nagios 2>/dev/null || true
        chown nagios:nagios \
            /usr/local/nagios/bin/nagios \
            /usr/local/nagios/bin/nagiostats 2>/dev/null || true
        chmod 750 \
            /usr/local/nagios/bin/nagios \
            /usr/local/nagios/bin/nagiostats 2>/dev/null || true
        ;;
esac
exit 0
POSTINST

    # prerm — Policy §6.1: para o daemon antes de remover
    install_script "$staging/DEBIAN/prerm" << 'PRERM'
#!/bin/bash
set -e
case "$1" in
    remove|purge)
        pkill -x nagios 2>/dev/null || true
        ;;
esac
exit 0
PRERM

    gen_copyright "$staging" "nagios-core" "https://www.nagios.org" "GPL-2+"
    gen_changelog "$staging" "nagios-core" "$NAGIOS_VERSION"
    gen_md5sums   "$staging"

    fakeroot dpkg-deb --build "$staging" "$deb"

    rm -rf "$src" "${src}.tar.gz" "$staging"
    echo "    $(basename "$deb") — $(du -sh "$deb" | cut -f1)"
}

# ═════════════════════════════════════════════════════════════════════════════
# Nagios Plugins
# ═════════════════════════════════════════════════════════════════════════════
build_nagios_plugins() {
    local name="nagios-plugins-${PLUGINS_VERSION}"
    local src="$BUILD_DIR/$name"
    local staging="$BUILD_DIR/pkg-plugins"
    local deb="$OUTPUT_DIR/nagios-plugins_${PLUGINS_VERSION}_amd64.deb"

    echo "==> [2/3] Nagios Plugins ${PLUGINS_VERSION}"

    cd "$BUILD_DIR"
    wget -q "https://github.com/nagios-plugins/nagios-plugins/releases/download/release-${PLUGINS_VERSION}/${name}.tar.gz"
    tar xzf "${name}.tar.gz" && cd "$src"

    echo "    configure..."
    # Ref: NAGIOS_BUILD_COMMANDS.md — recomenda --libexecdir=/usr/lib/nagios/plugins (FHS);
    # usamos o padrão (/usr/local/nagios/libexec) por compatibilidade com NagiosQL.
    ./configure \
        --with-nagios-user=nagios \
        --with-nagios-group=nagios \
        --with-mysql \
        >/dev/null

    echo "    make -j$(nproc)..."
    make -j"$(nproc)" >/dev/null

    echo "    make install..."
    make DESTDIR="$staging" install >/dev/null

    echo "    packaging..."
    mkdir -p "$staging/DEBIAN"

    cat > "$staging/DEBIAN/control" << EOF
Package: nagios-plugins
Version: ${PLUGINS_VERSION}
Architecture: amd64
Section: net
Priority: optional
Installed-Size: $(installed_size_kb "$staging")
Maintainer: ${MAINTAINER}
Homepage: https://www.nagios-plugins.org
Depends: libc6 (>= 2.17), libssl3, iputils-ping, dnsutils, procps
Description: standard plugins for Nagios-compatible monitoring systems
 This package provides the standard Nagios plugins: check_ping,
 check_http, check_ssh, check_dns, check_disk, check_load,
 check_users, check_procs, check_swap, and many others.
 .
 check_ping and check_icmp require NET_RAW capability and are
 installed setuid root:nagios to open raw ICMP sockets.
EOF
    chmod 644 "$staging/DEBIAN/control"

    # postinst aplica setuid nos plugins que precisam de raw sockets —
    # elimina o RUN chown/chmod separado no Dockerfile
    install_script "$staging/DEBIAN/postinst" << 'POSTINST'
#!/bin/bash
set -e
case "$1" in
    configure)
        for plugin in check_ping check_icmp; do
            bin="/usr/local/nagios/libexec/${plugin}"
            [ -f "$bin" ] || continue
            chown root:nagios "$bin"
            chmod u+s,g+x   "$bin"
        done
        ;;
esac
exit 0
POSTINST

    gen_copyright "$staging" "nagios-plugins" "https://www.nagios-plugins.org" "GPL-3+"
    gen_changelog "$staging" "nagios-plugins" "$PLUGINS_VERSION"
    gen_md5sums   "$staging"

    fakeroot dpkg-deb --build "$staging" "$deb"

    rm -rf "$src" "${src}.tar.gz" "$staging"
    echo "    $(basename "$deb") — $(du -sh "$deb" | cut -f1)"
}

# ═════════════════════════════════════════════════════════════════════════════
# NagiosQL
# ═════════════════════════════════════════════════════════════════════════════
package_nagiosql() {
    local staging="$BUILD_DIR/pkg-nagiosql"
    local deb="$OUTPUT_DIR/nagiosql_${NAGIOSQL_VERSION}_all.deb"

    echo "==> [3/3] NagiosQL ${NAGIOSQL_VERSION}"

    mkdir -p "$staging/var/www/nagiosql" "$staging/opt/nagiosql"
    cp -r "$NAGIOSQL_SRC/." "$staging/var/www/nagiosql/"
    # SQLs expostos em /opt/nagiosql para o entrypoint importar
    cp "$staging/var/www/nagiosql/install/sql/nagiosQL_v35_db_mysql.sql" \
       "$staging/opt/nagiosql/"
    cp "$staging/var/www/nagiosql/install/sql/import_nagios_sample.sql" \
       "$staging/opt/nagiosql/"
    rm -rf "$staging/var/www/nagiosql/install"

    mkdir -p "$staging/DEBIAN"

    cat > "$staging/DEBIAN/control" << EOF
Package: nagiosql
Version: ${NAGIOSQL_VERSION}
Architecture: all
Section: web
Priority: optional
Installed-Size: $(installed_size_kb "$staging")
Maintainer: ${MAINTAINER}
Homepage: https://www.nagiosql.org
Depends: nagios-core (>= 4), php8.4-fpm | php8.4-cli, php8.4-mysql, php8.4-mbstring, php8.4-gd, php8.4-curl, php8.4-xml, php8.4-zip
Description: web administration interface for Nagios Core
 NagiosQL is a browser-based administration tool for Nagios Core. It
 stores configuration objects (hosts, services, contacts, templates)
 in a MariaDB/MySQL database and writes them as .cfg files read
 directly by the Nagios daemon.
 .
 The database schema is provided in /opt/nagiosql/ for first-run
 import. No network connection is required between NagiosQL and
 Nagios Core — communication is entirely filesystem-based.
EOF
    chmod 644 "$staging/DEBIAN/control"

    # postinst ajusta ownership para www-data —
    # elimina o RUN chown/chmod separado no Dockerfile
    install_script "$staging/DEBIAN/postinst" << 'POSTINST'
#!/bin/bash
set -e
case "$1" in
    configure)
        # nagioscfg: grupo compartilhado entre www-data e nagios user para
        # escrita nos dirs de config. Ref: NAGIOSQL_DEBIAN_PACKAGING.md §File Permissions
        getent group nagioscfg >/dev/null 2>&1 \
            || groupadd -g 3001 nagioscfg
        usermod -aG nagioscfg www-data 2>/dev/null || true
        chown -R www-data:nagioscfg /var/www/nagiosql
        chmod -R 755 /var/www/nagiosql
        find /var/www/nagiosql/config -type d -exec chmod 775 {} \;
        ;;
esac
exit 0
POSTINST

    gen_copyright "$staging" "nagiosql" "https://www.nagiosql.org" "GPL-2+"
    gen_changelog "$staging" "nagiosql" "$NAGIOSQL_VERSION"
    gen_md5sums   "$staging"

    fakeroot dpkg-deb --build "$staging" "$deb"

    rm -rf "$staging"
    echo "    $(basename "$deb") — $(du -sh "$deb" | cut -f1)"
}

# ═════════════════════════════════════════════════════════════════════════════
# Main
# ═════════════════════════════════════════════════════════════════════════════
build_nagios_core
build_nagios_plugins
package_nagiosql

echo ""
echo "Pacotes gerados em $OUTPUT_DIR:"
ls -lh "$OUTPUT_DIR/"*.deb
