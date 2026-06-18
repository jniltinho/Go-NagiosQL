# NagiosQL 3.5 — Installation Guide
### Debian Trixie (13) · Nagios 4 via apt · Apache2 + PHP 8.4

---

## 1. Prerequisites

| Component | Min version | Debian package |
|---|---|---|
| Apache2 | 2.4 | `apache2` |
| PHP | 8.2+ | `libapache2-mod-php8.4` |
| MySQL/MariaDB | 5.7+ | `mariadb-server` or external |
| Nagios | 4.x | `nagios4-core` |
| PHP modules | — | `php8.4-mysql php8.4-gd php8.4-mbstring php8.4-curl php8.4-xml php8.4-zip` |

Install all packages:

```bash
apt-get update
apt-get install -y \
    nagios4-common nagios4-core monitoring-plugins nagios4-cgi \
    apache2 libapache2-mod-php8.4 \
    php8.4-mysql php8.4-gd php8.4-mbstring php8.4-curl php8.4-xml php8.4-zip \
    mariadb-server
```

---

## 2. Debian nagios4 package — default paths

| Resource | Path |
|---|---|
| Main config | `/etc/nagios4/nagios.cfg` |
| CGI config | `/etc/nagios4/cgi.cfg` |
| Resource file | `/etc/nagios4/resource.cfg` |
| Plugins | `/usr/lib/nagios/plugins/` |
| CGI binaries | `/usr/lib/cgi-bin/nagios4/` |
| Static files | `/usr/share/nagios4/htdocs/` |
| Daemon binary | `/usr/sbin/nagios4` |
| Var / runtime | `/var/lib/nagios4/` |
| External command pipe | `/var/lib/nagios4/rw/nagios.cmd` |
| PID file | `/run/nagios4/nagios4.pid` |
| Logs | `/var/log/nagios4/nagios.log` |
| Object cache | `/var/cache/nagios4/objects.cache` |

---

## 3. NagiosQL directory structure

NagiosQL writes Nagios object files to **`/etc/nagiosql/`** — a directory separate from `/etc/nagios4/`.

```bash
mkdir -p \
    /etc/nagiosql/hosts \
    /etc/nagiosql/services \
    /etc/nagiosql/backup/hosts \
    /etc/nagiosql/backup/services
```

### 3.1 Permissions

Apache (`www-data`) needs write access; the Nagios daemon (`nagios`) needs read access:

```bash
chown -R www-data:nagios /etc/nagiosql
find /etc/nagiosql -type d -exec chmod 750 {} \;
find /etc/nagiosql -type f -exec chmod 640 {} \;
```

NagiosQL also needs write access to two core Nagios files:

```bash
chown www-data:nagios /etc/nagios4/nagios.cfg /etc/nagios4/cgi.cfg
chmod 640 /etc/nagios4/nagios.cfg /etc/nagios4/cgi.cfg
```

---

## 4. Apache2 configuration

### 4.1 Enable modules

`libapache2-mod-php8.4` requires `mpm_prefork` (mod_php is not thread-safe).
The `nagios4-cgi` package installs `/etc/apache2/conf-available/nagios4-cgi.conf` with
`Require ip` (private IPs only) and Digest Auth — disable it and use a custom VirtualHost:

```bash
a2dismod mpm_event
a2enmod mpm_prefork php8.4 cgi
a2disconf nagios4-cgi
```

### 4.2 VirtualHost — Nagios Core (port 80)

Create `/etc/apache2/sites-available/nagios4.conf`:

```apache
<VirtualHost *:80>
    RedirectMatch ^/$ /cgi-bin/nagios4/tac.cgi

    Alias /nagios4 /usr/share/nagios4/htdocs
    <Directory /usr/share/nagios4/htdocs>
        Options FollowSymLinks
        AllowOverride None
        AuthName "Nagios Access"
        AuthType Basic
        AuthUserFile /etc/nagios4/htpasswd.users
        Require valid-user
    </Directory>

    ScriptAlias /cgi-bin/nagios4 /usr/lib/cgi-bin/nagios4
    <Directory /usr/lib/cgi-bin/nagios4>
        Options ExecCGI
        AllowOverride None
        SetEnv NAGIOS_CGI_CONFIG /etc/nagios4/cgi.cfg
        AuthName "Nagios Access"
        AuthType Basic
        AuthUserFile /etc/nagios4/htpasswd.users
        Require valid-user
    </Directory>
</VirtualHost>
```

### 4.3 VirtualHost — NagiosQL (port 8081)

Create `/etc/apache2/sites-available/nagiosql.conf`:

```apache
Listen 8081

<VirtualHost *:8081>
    DocumentRoot /var/www/nagiosql
    DirectoryIndex index.php

    <Directory /var/www/nagiosql>
        Options FollowSymLinks
        AllowOverride None
        Require all granted
    </Directory>

    <Directory /var/www/nagiosql/config>
        Require all denied
    </Directory>
</VirtualHost>
```

Enable sites:

```bash
a2dissite 000-default
a2ensite nagios4 nagiosql
systemctl restart apache2
```

### 4.4 PHP — required settings

Create `/etc/php/8.4/apache2/conf.d/99-nagiosql.ini`:

```ini
[Date]
date.timezone = America/Sao_Paulo

[Session]
session.auto_start = 0

[File]
file_uploads = On
```

---

## 5. Install NagiosQL

Copy files to `/var/www/nagiosql/`:

```bash
cd /opt
tar xzf nagiosql_350.tar.gz
mv nagiosql35 /var/www/nagiosql
chown -R www-data:nagios /var/www/nagiosql
chmod -R 755 /var/www/nagiosql
chmod 750 /var/www/nagiosql/config
```

Run the installation wizard:

```
http://server/nagiosql/install/index.php    (port 80)
http://server:8081/install/index.php        (port 8081)
```

---

## 6. NagiosQL configuration (Administration → Config targets)

After installation, go to **Administration → Config targets → localhost** and set:

| Field | Value |
|---|---|
| **Base directory** | `/etc/nagiosql/` |
| **Host directory** | `/etc/nagiosql/hosts/` |
| **Service directory** | `/etc/nagiosql/services/` |
| **Backup directory** | `/etc/nagiosql/backup/` |
| **Host backup dir** | `/etc/nagiosql/backup/hosts/` |
| **Service backup dir** | `/etc/nagiosql/backup/services/` |
| **Nagios base dir** | `/etc/nagios4/` |
| **Import directory** | `/etc/nagios4/conf.d/` |
| **Picture directory** | `/usr/share/nagios4/htdocs/images/logos/` |
| **Nagios command file** | `/var/lib/nagios4/rw/reload.trigger` |
| **Nagios binary** | `/usr/sbin/nagios4` |
| **Nagios process file** | `/run/nagios4/nagios4.pid` |
| **Nagios config file** | `/etc/nagios4/nagios.cfg` |
| **CGI config file** | `/etc/nagios4/cgi.cfg` |
| **Resource file** | `/etc/nagios4/resource.cfg` |
| **Nagios version** | `4` |

> **Note — command file:** NagiosQL uses a file-based reload trigger (`reload.trigger`) rather
> than the traditional external command pipe (`nagios.cmd`). The reload watcher script polls
> for this file and runs `nagios4 -v && systemctl reload nagios4` when it appears.

---

## 7. nagios.cfg — critical changes for NagiosQL

The Debian `nagios4` package ships default object files under `/etc/nagios4/objects/` that
**duplicate** the definitions NagiosQL generates. Loading both causes `duplicate definition`
errors that prevent the daemon from starting.

### 7.1 Comment out conflicting default entries

```bash
sed -i \
    -e 's|^cfg_dir=/etc/nagios-plugins/config|#cfg_dir=/etc/nagios-plugins/config|' \
    -e 's|^cfg_file=/etc/nagios4/objects/commands\.cfg|#cfg_file=/etc/nagios4/objects/commands.cfg|' \
    -e 's|^cfg_file=/etc/nagios4/objects/timeperiods\.cfg|#cfg_file=/etc/nagios4/objects/timeperiods.cfg|' \
    -e 's|^cfg_file=/etc/nagios4/objects/templates\.cfg|#cfg_file=/etc/nagios4/objects/templates.cfg|' \
    -e 's|^cfg_file=/etc/nagios4/objects/contacts\.cfg|#cfg_file=/etc/nagios4/objects/contacts.cfg|' \
    -e 's|^cfg_file=/etc/nagios4/objects/localhost\.cfg|#cfg_file=/etc/nagios4/objects/localhost.cfg|' \
    /etc/nagios4/nagios.cfg
```

Files that conflict and why:

| Default file | Conflict |
|---|---|
| `objects/commands.cfg` | NagiosQL generates `commands.cfg` in `/etc/nagiosql/` |
| `objects/timeperiods.cfg` | NagiosQL generates `timeperiods.cfg` in `/etc/nagiosql/` |
| `objects/templates.cfg` | NagiosQL generates host/service templates |
| `objects/contacts.cfg` | NagiosQL generates `contacts.cfg` + `contactgroups.cfg` |
| `objects/localhost.cfg` | NagiosQL manages all host definitions |

### 7.2 Add NagiosQL cfg entries

Append to `/etc/nagios4/nagios.cfg` (idempotent — only if not already present):

```bash
grep -q "etc/nagiosql/hosts" /etc/nagios4/nagios.cfg || cat >> /etc/nagios4/nagios.cfg << 'EOF'

# Configurations managed by NagiosQL (/etc/nagiosql/)
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
```

Also ensure `check_external_commands=1` is set:

```bash
grep check_external_commands /etc/nagios4/nagios.cfg
# check_external_commands=1
```

---

## 8. Runtime directories and reload trigger

### 8.1 rw/ directory permissions

The `rw/` directory must be writable by both the Nagios daemon (`nagios`) and Apache
(`www-data`). Use SGID so new files inherit the group:

```bash
mkdir -p /var/lib/nagios4/rw /var/lib/nagios4/spool/checkresults
chown nagios:www-data /var/lib/nagios4/rw
chmod 2775 /var/lib/nagios4/rw        # SGID: new files inherit www-data group

chown nagios:www-data /var/lib/nagios4/spool/checkresults
chmod 775 /var/lib/nagios4/spool/checkresults
```

### 8.2 reload.trigger

NagiosQL checks `file_exists()` on the command file before writing. Pre-create the trigger
file so NagiosQL can write to it even before the first reload:

```bash
touch /var/lib/nagios4/rw/reload.trigger
chown nagios:www-data /var/lib/nagios4/rw/reload.trigger
chmod 660 /var/lib/nagios4/rw/reload.trigger
```

### 8.3 Validate and reload

```bash
/usr/sbin/nagios4 -v /etc/nagios4/nagios.cfg
systemctl reload nagios4
```

---

## 9. cgi.cfg — authentication mode

The Debian package sets `use_authentication=0` in `/etc/nagios4/cgi.cfg` by default.
With Apache Basic Auth in front, this is correct: Apache handles authentication and every
authenticated user gets full access to the Nagios CGI. Verify:

```bash
grep use_authentication /etc/nagios4/cgi.cfg
# use_authentication=0
```

Also create the stylesheets directory referenced in the VirtualHost:

```bash
mkdir -p /etc/nagios4/stylesheets
```

---

## 10. htpasswd — Nagios Core user

```bash
# Create the file with the nagiosadmin user
htpasswd -c /etc/nagios4/htpasswd.users nagiosadmin

# Fix permissions (Apache needs read access)
chown nagios:www-data /etc/nagios4/htpasswd.users
chmod 640 /etc/nagios4/htpasswd.users
```

---

## 11. Final verification

```
http://server:80/cgi-bin/nagios4/tac.cgi   → Nagios Core (Basic Auth)
http://server:8081/                          → NagiosQL (own login)
```

In NagiosQL: **Tools → Support page** shows the status of all configured paths.

Validate the Nagios config and confirm zero errors:

```bash
/usr/sbin/nagios4 -v /etc/nagios4/nagios.cfg 2>&1 | grep "Total Errors"
# Total Errors:   0
```

---

## Appendix — Differences from older guides (nagios3 / pre-2020)

| Old guide | Debian nagios4 / trixie |
|---|---|
| `/etc/nagios/` | `/etc/nagios4/` |
| `/opt/nagios/bin/nagios` | `/usr/sbin/nagios4` |
| `/var/nagios/` | `/var/lib/nagios4/` |
| `/var/nagios/rw/nagios.cmd` | `/var/lib/nagios4/rw/nagios.cmd` |
| `/var/nagios/nagios.lock` | `/run/nagios4/nagios4.pid` |
| `/usr/local/nagios/libexec/` | `/usr/lib/nagios/plugins/` |
| PHP 5.x / 7.x | PHP 8.4 |
| `php5-mysql` | `php8.4-mysql` |
| Apache `wwwrun` | Apache `www-data` |
| `nagios.cmd` as command file | `reload.trigger` (file-based reload) |
| Default object files loaded | Default object files **must be commented out** |
