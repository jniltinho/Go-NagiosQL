ARG NAGIOSQL_VERSION=3.5.0

# ════════════════════════════════════════════════════════════════
# Stage 1 — builder: empacota o NagiosQL como .deb
# ════════════════════════════════════════════════════════════════
FROM debian:trixie-slim AS builder

ENV DEBIAN_FRONTEND=noninteractive
ARG NAGIOSQL_VERSION

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates fakeroot dpkg-dev gzip \
    && rm -rf /var/lib/apt/lists/*

COPY nagios/build-debs.sh /build-debs.sh
COPY nagiosql.tar.gz /tmp/nagiosql.tar.gz
RUN mkdir -p /nagiosql-src \
    && tar -xzf /tmp/nagiosql.tar.gz --strip-components=1 -C /nagiosql-src \
    && rm /tmp/nagiosql.tar.gz
RUN chmod +x /build-debs.sh && /build-debs.sh

# ════════════════════════════════════════════════════════════════
# Stage 2 — runtime
# ════════════════════════════════════════════════════════════════
FROM debian:trixie-slim

ENV DEBIAN_FRONTEND=noninteractive
ARG NAGIOSQL_VERSION

# ── Nagios4 via pacotes Debian + Apache2 + PHP mod_php ────────
# nagios4-cgi depende de nagios4-common; monitoring-plugins inclui os checks
# libapache2-mod-php8.4: PHP integrado ao Apache (mod_php, sem fpm)
RUN apt-get update && apt-get install -y --no-install-recommends \
    nagios4-common nagios4-core monitoring-plugins nagios4-cgi \
    apache2 libapache2-mod-php8.4 \
    php8.4-mysql php8.4-gd php8.4-mbstring \
    php8.4-curl php8.4-xml php8.4-zip php-pear \
    supervisor default-mysql-client \
    ca-certificates openssl \
    iputils-ping dnsutils procps netcat-openbsd fping snmp \
    && rm -rf /var/lib/apt/lists/*

# ── Apache: módulos necessários ───────────────────────────────
# mpm_prefork: obrigatório para mod_php (não thread-safe)
# mod_cgi: executa os CGIs do Nagios4
# Desabilita o conf do pacote nagios4-cgi — usamos VirtualHost próprio
RUN a2dismod mpm_event 2>/dev/null; true \
    && a2enmod mpm_prefork php8.4 cgi \
    && a2disconf nagios4-cgi 2>/dev/null; true \
    && a2dissite 000-default

# ── Instalar NagiosQL .deb do stage builder ───────────────────
COPY --from=builder /tmp/debs/nagiosql_${NAGIOSQL_VERSION}_all.deb /tmp/
RUN dpkg -i --force-depends /tmp/nagiosql_${NAGIOSQL_VERSION}_all.deb \
    && rm /tmp/*.deb

# ── Snapshot do config NagiosQL para seed dos volumes ─────────
RUN cp -r /var/www/nagiosql/config /var/www/nagiosql/config.default

# ── /etc/nagiosql: estrutura de configs gerados pelo NagiosQL ─
# Separado de /etc/nagios4 conforme PDF §1.2 (base dir /etc/nagiosql)
RUN mkdir -p \
    /etc/nagiosql/hosts \
    /etc/nagiosql/services \
    /etc/nagiosql/backup/hosts \
    /etc/nagiosql/backup/services
COPY nagios/etc-extra/nagiosql/ /etc/nagiosql/
RUN cp -a /etc/nagiosql /etc/nagiosql.default

# ── Salvar defaults do Nagios4 e plugins para init dos volumes ─
RUN cp -r /etc/nagios4 /etc/nagios4.default \
    && cp -a /usr/lib/nagios/plugins /usr/lib/nagios/plugins.default

# ── Apache: sites (Nagios :80 e NagiosQL :8081) ───────────────
COPY nagios/apache-nagios4.conf  /etc/apache2/sites-available/nagios4.conf
COPY nagios/apache-nagiosql.conf /etc/apache2/sites-available/nagiosql.conf
RUN a2ensite nagios4 nagiosql

# ── Supervisord e scripts ─────────────────────────────────────
COPY nagios/supervisord.conf  /etc/supervisor/conf.d/nagios.conf
COPY nagios/reload-watcher.sh /usr/local/bin/reload-watcher.sh
COPY nagios/entrypoint.sh     /entrypoint.sh
RUN chmod +x /usr/local/bin/reload-watcher.sh /entrypoint.sh

WORKDIR /etc/nagiosql

EXPOSE 80 8081
ENTRYPOINT ["/entrypoint.sh"]
