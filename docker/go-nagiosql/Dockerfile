# ── Stage 1: Build Go binary ────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
ARG BUILDDATE=unknown

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags "-s -w -X main.Version=${VERSION} -X main.BuildDate=${BUILDDATE}" \
    -o /nagiosql .

# ── Stage 2: Runtime — Debian trixie + nagios4 (apt) + nagiosql binary ──────
FROM debian:trixie-slim

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
    nagios4-common nagios4-core monitoring-plugins nagios4-cgi \
    apache2 \
    supervisor default-mysql-client \
    ca-certificates openssl \
    iputils-ping dnsutils procps netcat-openbsd fping snmp \
    && rm -rf /var/lib/apt/lists/*

# Apache: CGI only — no PHP, no mod_php
RUN a2dismod mpm_event 2>/dev/null; true \
    && a2enmod mpm_prefork cgi \
    && a2disconf nagios4-cgi 2>/dev/null; true \
    && a2dissite 000-default

# Snapshot Nagios4 defaults for volume seed on first boot
RUN cp -r /etc/nagios4 /etc/nagios4.default \
    && cp -a /usr/lib/nagios/plugins /usr/lib/nagios/plugins.default

# /etc/nagiosql: seed files + directory structure for first boot
RUN mkdir -p \
    /etc/nagiosql/hosts \
    /etc/nagiosql/services \
    /etc/nagiosql/backup/hosts \
    /etc/nagiosql/backup/services
COPY docker/go-nagiosql/nagios/etc-nagiosql/ /etc/nagiosql/
RUN cp -a /etc/nagiosql /etc/nagiosql.default

# Go binary + config — /opt/nagiosql is the API working directory
RUN mkdir -p /opt/nagiosql
COPY --from=builder /nagiosql /opt/nagiosql/nagiosql
COPY docker/go-nagiosql/config.toml.docker /opt/nagiosql/config.toml

# Apache vhost: Nagios Core CGI on port 80
COPY docker/go-nagiosql/nagios/apache-nagios4.conf /etc/apache2/sites-available/nagios4.conf
RUN a2ensite nagios4

# Supervisord and runtime scripts
COPY docker/go-nagiosql/nagios/supervisord.conf  /etc/supervisor/conf.d/nagios.conf
COPY docker/go-nagiosql/nagios/reload-watcher.sh /usr/local/bin/reload-watcher.sh
COPY docker/go-nagiosql/nagios/entrypoint.sh     /entrypoint.sh
RUN chmod +x /usr/local/bin/reload-watcher.sh /entrypoint.sh

WORKDIR /opt/nagiosql

# 80: Nagios Core web (Apache CGI)
# 8081: nagiosql REST API (Go server)
EXPOSE 80 8081
ENTRYPOINT ["/entrypoint.sh"]
