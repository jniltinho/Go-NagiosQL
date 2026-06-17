
# Building and Packaging Nagios Core 4.5.13 and Nagios Plugins 2.5 for Debian 13 (Trixie)

**Author:** Nilton Oliveira  
**Version:** 1.0  
**Target Platform:** Debian GNU/Linux 13 (Trixie)

---

# Table of Contents

1. Introduction
2. Architecture
3. Build Environment
4. Required Dependencies
5. Downloading Sources
6. Building Nagios Core
7. Creating Debian Packages
8. Building Nagios Plugins
9. Nagios Users and Permissions
10. Directory Layout
11. Systemd Integration
12. Nginx + fcgiwrap Configuration
13. SSL with Let's Encrypt
14. Initial Nagios Configuration
15. Host and Service Definitions
16. Service Management
17. Security Hardening
18. Reproducible Builds with sbuild
19. Internal APT Repository
20. GitHub Actions CI/CD
21. Troubleshooting
22. Upgrade Procedures
23. Backup Procedures
24. Production Checklist

---

# Introduction

This guide describes how to build, package, deploy, and maintain:

- Nagios Core 4.5.13
- Nagios Plugins 2.5

using Debian packaging standards.

Goals:

- Produce clean .deb packages
- Follow Debian Policy
- Avoid checkinstall
- Enable CI/CD automation
- Support internal repositories
- Use Nginx instead of Apache

---

# Architecture

```text
Internet
    |
Nginx (HTTPS)
    |
fcgiwrap
    |
Nagios CGI
    |
Nagios Core
    |
Nagios Plugins
    |
Monitored Hosts
```

---

# Build Environment

Recommended:

- Debian 13 (Trixie)
- 4 vCPU
- 8 GB RAM
- 20 GB Storage

Workspace:

```bash
mkdir -p ~/build
cd ~/build
```

---

# Required Dependencies

```bash
apt update

apt install -y  build-essential  devscripts  debhelper  dh-make  dpkg-dev  fakeroot  lintian  wget  curl  git  nginx  fcgiwrap  apache2-utils  libgd-dev  libssl-dev  apache2-dev
```

---

# Downloading Sources

## Nagios Core

```bash
wget https://github.com/NagiosEnterprises/nagioscore/releases/download/nagios-4.5.13/nagios-4.5.13.tar.gz
```

## Nagios Plugins

```bash
wget https://github.com/nagios-plugins/nagios-plugins/releases/download/release-2.5/nagios-plugins-2.5.tar.gz
```

---

# Building Nagios Core

```bash
tar xzf nagios-4.5.13.tar.gz
cd nagios-4.5.13

./configure   --prefix=/usr   --sysconfdir=/etc/nagios   --localstatedir=/var/lib/nagios   --with-command-group=nagios

make -j$(nproc)
```

Verify:

```bash
src/nagios --version
```

---

# Debian Packaging

Create skeleton:

```bash
dh_make --createorig -s -y
find debian -name "*.ex" -delete
```

## debian/control

```text
Source: nagios-core
Section: net
Priority: optional
Maintainer: Nilton Oliveira <jniltinho@gmail.com>

Build-Depends:
 debhelper-compat (= 13),
 libgd-dev,
 libssl-dev

Package: nagios-core
Architecture: any
Depends: ${shlibs:Depends}, ${misc:Depends}
Description: Nagios Core Monitoring System
```

## debian/rules

```make
#!/usr/bin/make -f

%:
	dh $@

override_dh_auto_configure:
	./configure 		--prefix=/usr 		--sysconfdir=/etc/nagios 		--localstatedir=/var/lib/nagios
```

Build:

```bash
dpkg-buildpackage -us -uc -b
```

---

# Building Nagios Plugins

```bash
tar xzf nagios-plugins-2.5.tar.gz
cd nagios-plugins-2.5

./configure  --prefix=/usr  --libexecdir=/usr/lib/nagios/plugins

make -j$(nproc)
```

Package:

```bash
dh_make --createorig -s -y
dpkg-buildpackage -us -uc -b
```

---

# Users and Permissions

```bash
groupadd --system nagios

useradd  --system  --gid nagios  --home /var/lib/nagios  --shell /usr/sbin/nologin  nagios
```

Directories:

```bash
mkdir -p  /etc/nagios  /var/lib/nagios  /var/log/nagios

chown -R nagios:nagios  /var/lib/nagios  /var/log/nagios
```

---

# Directory Layout

```text
/etc/nagios
/etc/nagios/conf.d
/usr/lib/nagios/plugins
/usr/lib/nagios/cgi
/var/log/nagios
/var/lib/nagios
/usr/share/nagios/html
```

---

# Systemd Service

File:

```text
/etc/systemd/system/nagios.service
```

```ini
[Unit]
Description=Nagios Core
After=network.target

[Service]
Type=forking
ExecStart=/usr/sbin/nagios /etc/nagios/nagios.cfg
PIDFile=/run/nagios.pid
Restart=always

[Install]
WantedBy=multi-user.target
```

Enable:

```bash
systemctl daemon-reload
systemctl enable nagios
```

---

# Nginx + fcgiwrap

Install:

```bash
apt install -y nginx fcgiwrap
```

Create password:

```bash
htpasswd -c /etc/nagios/htpasswd.users nagiosadmin
```

Virtual Host:

```nginx
server {
    listen 80;
    server_name nagios.example.com;

    root /usr/share/nagios/html;

    auth_basic "Nagios";
    auth_basic_user_file /etc/nagios/htpasswd.users;

    location / {
        index index.html;
    }

    location /cgi-bin/ {
        alias /usr/lib/nagios/cgi/;

        include fastcgi_params;

        fastcgi_param SCRIPT_FILENAME $request_filename;

        fastcgi_pass unix:/run/fcgiwrap.socket;
    }
}
```

Enable:

```bash
ln -s /etc/nginx/sites-available/nagios.conf       /etc/nginx/sites-enabled/

nginx -t
systemctl reload nginx
```

---

# SSL Configuration

```bash
apt install certbot python3-certbot-nginx

certbot --nginx -d nagios.example.com
```

Verify:

```bash
systemctl list-timers | grep certbot
```

---

# Initial Nagios Configuration

Validate:

```bash
nagios -v /etc/nagios/nagios.cfg
```

Expected:

```text
Things look okay - No serious problems were detected
```

---

# Example Host Definition

```cfg
define host {
    use             linux-server
    host_name       app01
    alias           Application Server
    address         10.0.0.10
}
```

---

# Example Service Definition

```cfg
define service {
    use                     generic-service
    host_name               app01
    service_description     PING
    check_command           check_ping!100.0,20%!500.0,60%
}
```

---

# Service Management

```bash
systemctl start nagios
systemctl stop nagios
systemctl restart nagios
systemctl status nagios
```

Logs:

```bash
journalctl -u nagios -f
```

---

# Security Hardening

Nginx:

```nginx
server_tokens off;

add_header X-Frame-Options SAMEORIGIN;
add_header X-Content-Type-Options nosniff;
add_header Referrer-Policy strict-origin;
```

Recommendations:

- HTTPS only
- Disable anonymous access
- Restrict CGI access
- Enable automatic security updates
- Use firewall rules

---

# Reproducible Builds

Install:

```bash
apt install sbuild
```

Create chroot:

```bash
sbuild-createchroot  trixie  /srv/chroot/trixie-amd64  http://deb.debian.org/debian
```

Build:

```bash
sbuild
```

---

# Internal APT Repository

Install:

```bash
apt install reprepro
```

Structure:

```text
repo/
├── conf
├── dists
└── pool
```

Add packages:

```bash
reprepro includedeb trixie nagios-core_4.5.13-1_amd64.deb
reprepro includedeb trixie nagios-plugins_2.5-1_amd64.deb
```

---

# GitHub Actions

```yaml
name: Build

on:
  push:

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Install Dependencies
        run: |
          sudo apt update
          sudo apt install -y devscripts debhelper

      - name: Build Package
        run: |
          dpkg-buildpackage -us -uc -b
```

---

# Backup Procedure

```bash
tar czf nagios-backup.tar.gz  /etc/nagios  /var/lib/nagios  /var/log/nagios
```

---

# Upgrade Procedure

```bash
apt update
apt install nagios-core nagios-plugins
```

Validate:

```bash
nagios -v /etc/nagios/nagios.cfg
```

Restart:

```bash
systemctl restart nagios
```

---

# Troubleshooting

## CGI 502 Error

```bash
systemctl status fcgiwrap
ls -la /run/fcgiwrap.socket
```

## Nagios Won't Start

```bash
nagios -v /etc/nagios/nagios.cfg
```

## Nginx 403

Check permissions:

```bash
ls -la /usr/share/nagios
```

## Plugin Not Found

```bash
ls -la /usr/lib/nagios/plugins
```

---

# Production Checklist

- [ ] Packages built successfully
- [ ] Lintian validation passed
- [ ] Nagios configuration validated
- [ ] HTTPS enabled
- [ ] Firewall configured
- [ ] Backups configured
- [ ] Monitoring tested
- [ ] Internal APT repository configured
- [ ] CI/CD pipeline operational
- [ ] Documentation stored in Git repository

---

# Deliverables

```text
nagios-core_4.5.13-1_amd64.deb
nagios-plugins_2.5-1_amd64.deb
```
