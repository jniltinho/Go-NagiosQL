# NagiosQL 3.5 — Guia de Instalação
### Debian Trixie (13) · Nagios4 via pacote apt · Apache2 + PHP 8.4

---

## 1. Pré-requisitos

| Componente | Versão mínima | Pacote Debian |
|---|---|---|
| Apache2 | 2.4 | `apache2` |
| PHP | 8.2+ | `libapache2-mod-php8.4` |
| MySQL/MariaDB | 5.7+ | `mariadb-server` ou externo |
| Nagios | 4.x | `nagios4-core` |
| PHP módulos | — | `php8.4-mysql php8.4-gd php8.4-mbstring php8.4-curl php8.4-xml php8.4-zip` |

Instalar pacotes:

```bash
apt-get update
apt-get install -y \
    nagios4-common nagios4-core monitoring-plugins nagios4-cgi \
    apache2 libapache2-mod-php8.4 \
    php8.4-mysql php8.4-gd php8.4-mbstring php8.4-curl php8.4-xml php8.4-zip \
    mariadb-server
```

---

## 2. Caminhos padrão do pacote Debian nagios4

| Recurso | Caminho |
|---|---|
| Config principal | `/etc/nagios4/nagios.cfg` |
| Config CGI | `/etc/nagios4/cgi.cfg` |
| Resource file | `/etc/nagios4/resource.cfg` |
| Plugins | `/usr/lib/nagios/plugins/` |
| CGI binaries | `/usr/lib/cgi-bin/nagios4/` |
| Arquivos estáticos | `/usr/share/nagios4/htdocs/` |
| Binário daemon | `/usr/sbin/nagios4` |
| Var/runtime | `/var/lib/nagios4/` |
| Command pipe | `/var/lib/nagios4/rw/nagios.cmd` |
| PID file | `/run/nagios4/nagios4.pid` |
| Logs | `/var/log/nagios4/nagios.log` |
| Object cache | `/var/cache/nagios4/objects.cache` |

---

## 3. Estrutura de diretórios do NagiosQL

O NagiosQL armazena os arquivos de configuração do Nagios em **`/etc/nagiosql/`** — diretório separado do `/etc/nagios4/`.

```bash
mkdir -p \
    /etc/nagiosql/hosts \
    /etc/nagiosql/services \
    /etc/nagiosql/backup/hosts \
    /etc/nagiosql/backup/services
```

### 3.1 Permissões

O usuário do Apache (`www-data`) precisa de escrita; o daemon Nagios (`nagios`) precisa de leitura:

```bash
# www-data (owner) escreve, nagios (group) lê
chown -R www-data:nagios /etc/nagiosql
find /etc/nagiosql -type d -exec chmod 750 {} \;
find /etc/nagiosql -type f -exec chmod 640 {} \;
```

Arquivos de config do Nagios que o NagiosQL modifica:

```bash
chown www-data:nagios /etc/nagios4/nagios.cfg /etc/nagios4/cgi.cfg
chmod 640 /etc/nagios4/nagios.cfg /etc/nagios4/cgi.cfg
```

---

## 4. Configuração do Apache2

### 4.1 Habilitar módulos

O `libapache2-mod-php8.4` requer `mpm_prefork` (mod_php não é thread-safe).
O pacote `nagios4-cgi` instala `/etc/apache2/conf-available/nagios4-cgi.conf` com
`Require ip` (só IPs privados) e Digest Auth — desabilitar e usar VirtualHost próprio:

```bash
a2dismod mpm_event
a2enmod mpm_prefork php8.4 cgi
a2disconf nagios4-cgi   # desabilita o conf padrão do pacote
```

### 4.2 VirtualHost — Nagios Core (porta 80)

Criar `/etc/apache2/sites-available/nagios4.conf`:

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

### 4.3 VirtualHost — NagiosQL (porta 8081)

Criar `/etc/apache2/sites-available/nagiosql.conf`:

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

Ativar sites:

```bash
a2dissite 000-default
a2ensite nagios4 nagiosql
systemctl restart apache2
```

### 4.4 PHP — configurações obrigatórias

Criar `/etc/php/8.4/apache2/conf.d/99-nagiosql.ini`:

```ini
[Date]
date.timezone = America/Sao_Paulo

[Session]
session.auto_start = 0

[File]
file_uploads = On
```

---

## 5. Instalar o NagiosQL

Copiar os arquivos para `/var/www/nagiosql/`:

```bash
cd /opt
tar xzf nagiosql_350.tar.gz
mv nagiosql35 /var/www/nagiosql
chown -R www-data:nagios /var/www/nagiosql
chmod -R 755 /var/www/nagiosql
chmod 750 /var/www/nagiosql/config
```

Acessar o wizard de instalação:

```
http://servidor/nagiosql/install/index.php    (porta 80)
http://servidor:8081/install/index.php        (porta 8081)
```

---

## 6. Configuração do NagiosQL (Administration → Config targets)

Após a instalação, em **Administration → Config targets → localhost**, definir:

| Campo | Valor |
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
| **Nagios command file** | `/var/lib/nagios4/rw/nagios.cmd` |
| **Nagios binary** | `/usr/sbin/nagios4` |
| **Nagios process file** | `/run/nagios4/nagios4.pid` |
| **Nagios config file** | `/etc/nagios4/nagios.cfg` |
| **CGI config file** | `/etc/nagios4/cgi.cfg` |
| **Resource file** | `/etc/nagios4/resource.cfg` |
| **Nagios version** | `4` |

---

## 7. Configuração do nagios.cfg

Adicionar ao final de `/etc/nagios4/nagios.cfg` para que o Nagios leia os arquivos gerados pelo NagiosQL:

```
# Configurações gerenciadas pelo NagiosQL
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
```

---

## 8. Command pipe e reload do Nagios

O Debian nagios4 cria `/var/lib/nagios4/rw/` com grupo `www-data` — o Apache pode escrever nesse diretório por padrão.

Verificar que `check_external_commands=1` está em `nagios.cfg`:

```bash
grep check_external_commands /etc/nagios4/nagios.cfg
# check_external_commands=1
```

Validar e recarregar:

```bash
/usr/sbin/nagios4 -v /etc/nagios4/nagios.cfg
systemctl reload nagios4
```

---

## 9. cgi.cfg — autenticação via Apache

O pacote Debian define `use_authentication=0` em `/etc/nagios4/cgi.cfg` por padrão.
Com Basic Auth no Apache, esse é o modo correto: qualquer usuário autenticado pelo Apache
tem acesso completo ao Nagios. Verificar:

```bash
grep use_authentication /etc/nagios4/cgi.cfg
# use_authentication=0
```

Criar também o diretório de stylesheets referenciado no VirtualHost:

```bash
mkdir -p /etc/nagios4/stylesheets
```

---

## 10. htpasswd — usuário do Nagios Core

```bash
# Criar o arquivo com o usuário nagiosadmin
htpasswd -c /etc/nagios4/htpasswd.users nagiosadmin

# Permissões
chown nagios:www-data /etc/nagios4/htpasswd.users
chmod 640 /etc/nagios4/htpasswd.users
```

---

## 11. Verificação final

```
http://servidor:80/cgi-bin/nagios4/tac.cgi   → Nagios Core (Basic Auth)
http://servidor:8081/                          → NagiosQL (login próprio)
```

Verificar no NagiosQL: **Tools → Support page** mostra o status de todos os caminhos configurados.

---

## Apêndice — Resumo de diferenças para guias antigos (nagios3 / pré-2020)

| Guia antigo | Debian nagios4 / trixie |
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
