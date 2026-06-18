#!/bin/bash
# build-debs.sh — Empacota o NagiosQL como .deb
set -euo pipefail

NAGIOSQL_VERSION="${NAGIOSQL_VERSION:-3.5.0}"
NAGIOSQL_SRC="${NAGIOSQL_SRC:-/nagiosql-src}"
MAINTAINER="${MAINTAINER:-Nagios Docker Build <build@local>}"

BUILD_DIR="/tmp/build"
OUTPUT_DIR="/tmp/debs"
mkdir -p "$BUILD_DIR" "$OUTPUT_DIR"

installed_size_kb() { du -sk --apparent-size "$1" | cut -f1; }

gen_md5sums() {
    local staging="$1"
    (cd "$staging" && \
        find . -type f ! -path './DEBIAN/*' | LC_ALL=C sort | \
        xargs md5sum | sed 's|^\./||') > "$staging/DEBIAN/md5sums"
    chmod 644 "$staging/DEBIAN/md5sums"
}

install_script() { local file="$1"; cat > "$file"; chmod 755 "$file"; }

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

gen_changelog() {
    local staging="$1" pkg="$2" ver="$3"
    mkdir -p "$staging/usr/share/doc/$pkg"
    printf '%s (%s) unstable; urgency=low\n\n  * Package built for Docker deployment.\n\n -- %s  %s\n' \
        "$pkg" "$ver" "$MAINTAINER" "$(date -R)" \
        | gzip -9 > "$staging/usr/share/doc/$pkg/changelog.Debian.gz"
    chmod 644 "$staging/usr/share/doc/$pkg/changelog.Debian.gz"
}

package_nagiosql() {
    local staging="$BUILD_DIR/pkg-nagiosql"
    local deb="$OUTPUT_DIR/nagiosql_${NAGIOSQL_VERSION}_all.deb"

    echo "==> NagiosQL ${NAGIOSQL_VERSION}"

    mkdir -p "$staging/var/www/nagiosql" "$staging/opt/nagiosql"
    cp -r "$NAGIOSQL_SRC/." "$staging/var/www/nagiosql/"
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
Depends: nagios4-core
Description: web administration interface for Nagios Core
 NagiosQL is a browser-based administration tool for Nagios Core. It
 stores configuration objects (hosts, services, contacts, templates)
 in a MariaDB/MySQL database and writes them as .cfg files read
 directly by the Nagios daemon.
EOF
    chmod 644 "$staging/DEBIAN/control"

    install_script "$staging/DEBIAN/postinst" << 'POSTINST'
#!/bin/bash
set -e
case "$1" in
    configure)
        chown -R www-data:nagios /var/www/nagiosql
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

package_nagiosql

echo ""
echo "Pacotes gerados em $OUTPUT_DIR:"
ls -lh "$OUTPUT_DIR/"*.deb
