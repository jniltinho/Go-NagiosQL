# ARGs globais — acessíveis em todos os stages
ARG NAGIOS_VERSION=4.5.13
ARG PLUGINS_VERSION=2.5
ARG NAGIOSQL_VERSION=3.5.0

# ════════════════════════════════════════════════════════════════
# Stage 1 — builder: compila e empacota os .deb
# ════════════════════════════════════════════════════════════════
FROM debian:trixie-slim AS builder

ENV DEBIAN_FRONTEND=noninteractive
ARG NAGIOS_VERSION
ARG PLUGINS_VERSION
ARG NAGIOSQL_VERSION

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    autoconf gcc make wget unzip fakeroot \
    libgd-dev openssl libssl-dev \
    iputils-ping dnsutils procps fping \
    snmp libmariadb-dev \
    && rm -rf /var/lib/apt/lists/* \
    && ln -sf /usr/bin/mariadb_config /usr/local/bin/mysql_config \
    && ln -sf /usr/lib/x86_64-linux-gnu/libmariadb.so \
              /usr/local/lib/libmysqlclient.so \
    && ldconfig

# nagios + nagioscfg necessários no builder: make install resolve grupos no DESTDIR
RUN groupadd -g 3000 nagios && \
    groupadd -g 3001 nagioscfg && \
    useradd -u 3000 -g nagios -G nagioscfg -d /usr/local/nagios -s /bin/bash nagios

# Script de build isolado do código-fonte para melhor cache de camadas
COPY nagios/build-debs.sh /build-debs.sh
COPY nagiosql/ /nagiosql-src/
RUN chmod +x /build-debs.sh && /build-debs.sh

# ════════════════════════════════════════════════════════════════
# Stage 2 — runtime: instala os .deb e configura o container
# ════════════════════════════════════════════════════════════════
FROM debian:trixie-slim

ENV DEBIAN_FRONTEND=noninteractive
ARG NAGIOS_VERSION
ARG PLUGINS_VERSION
ARG NAGIOSQL_VERSION

# ── Dependências de runtime ──────────────────────────────────────
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates libgd3 openssl \
    nginx fcgiwrap supervisor \
    php8.4-fpm php8.4-gd \
    php8.4-cli php8.4-common \
    php8.4-mbstring php8.4-mysql \
    php8.4-curl php8.4-xml php8.4-zip \
    php-pear \
    default-mysql-client \
    iputils-ping dnsutils procps netcat-openbsd fping \
    snmp libmariadb3 \
    && rm -rf /var/lib/apt/lists/*

# ── Grupo/usuário nagios + nagioscfg ─────────────────────────────
# nagioscfg: grupo compartilhado para escrita nos dirs de config
# (nagios user + www-data) — evita dar www-data acesso ao grupo nagios inteiro.
# Ref: NAGIOSQL_DEBIAN_PACKAGING.md §File Permissions
RUN groupadd -g 3000 nagios && \
    groupadd -g 3001 nagioscfg && \
    useradd -u 3000 -g nagios -G nagioscfg -d /usr/local/nagios -s /bin/bash nagios && \
    usermod -a -G nagioscfg www-data

# ── Instalar os .deb gerados no stage builder ────────────────────
COPY --from=builder /tmp/debs/nagios-core_${NAGIOS_VERSION}_amd64.deb     /tmp/
COPY --from=builder /tmp/debs/nagios-plugins_${PLUGINS_VERSION}_amd64.deb /tmp/
COPY --from=builder /tmp/debs/nagiosql_${NAGIOSQL_VERSION}_all.deb        /tmp/

RUN dpkg -i \
        /tmp/nagios-core_${NAGIOS_VERSION}_amd64.deb \
        /tmp/nagios-plugins_${PLUGINS_VERSION}_amd64.deb \
        /tmp/nagiosql_${NAGIOSQL_VERSION}_all.deb && \
    rm /tmp/*.deb

# ── Snapshot do config NagiosQL para seed dos volumes ────────────
# Executado após dpkg (postinst já aplicou chown/chmod em /var/www/nagiosql)
RUN cp -r /var/www/nagiosql/config /var/www/nagiosql/config.default

# ── Configs de exemplo (hosts e serviços iniciais) ───────────────
COPY nagios/etc-extra/ /usr/local/nagios/etc/

# ── Salvar defaults para inicialização dos volumes ────────────────
RUN cp -r /usr/local/nagios/etc     /usr/local/nagios/etc.default && \
    cp -a /usr/local/nagios/libexec /usr/local/nagios/libexec.default

# ── Nginx (Nagios Core :80 + NagiosQL :8081) ─────────────────────
COPY nagios/nginx.conf /etc/nginx/sites-available/nagios.conf
RUN rm -f /etc/nginx/sites-enabled/default && \
    ln -sf /etc/nginx/sites-available/nagios.conf /etc/nginx/sites-enabled/nagios.conf && \
    ln -sf /usr/local/nagios/sbin /usr/local/nagios/share/cgi-bin

# ── Supervisord e scripts ─────────────────────────────────────────
COPY nagios/supervisord.conf  /etc/supervisor/conf.d/nagios.conf
COPY nagios/reload-watcher.sh /usr/local/nagios/bin/reload-watcher.sh
COPY nagios/entrypoint.sh     /entrypoint.sh
RUN chmod +x /usr/local/nagios/bin/reload-watcher.sh /entrypoint.sh

EXPOSE 80 8081
ENTRYPOINT ["/entrypoint.sh"]
