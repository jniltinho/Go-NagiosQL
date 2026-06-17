# Instalação do Nagios Core 4 com Nginx no Ubuntu 24.04

Guia completo para instalar o Nagios Core 4 com Nginx, PHP-FPM e fcgiwrap no Ubuntu 24.04 LTS.

**Referências:**
- https://github.com/NagiosEnterprises/nagioscore
- https://turndevopseasier.com/setting-up-nagios-core-with-nginx-on-ubuntu-24-04/

---

## Pré-requisitos

- Ubuntu 24.04 LTS (servidor limpo)
- Acesso root ou usuário com sudo
- Conexão com a internet

---

## 1. Atualizar o sistema e instalar dependências

```bash
sudo apt-get update && sudo apt-get upgrade -y

sudo apt-get install -y \
  autoconf gcc libc6 make wget curl unzip \
  libgd-dev openssl libssl-dev ufw

sudo apt-get install -y \
  php8.3 php8.3-fpm php8.3-common php8.3-cli \
  php8.3-mbstring php8.3-bcmath php8.3-mysql \
  php8.3-zip php8.3-gd php8.3-curl php8.3-xml

sudo apt-get install -y nginx fcgiwrap
```

---

## 2. Baixar o Nagios Core (versão mais recente)

```bash
mkdir -p ~/nagios-core && cd ~/nagios-core

wget -O nagioscore.tar.gz \
  $(curl -s "https://api.github.com/repos/NagiosEnterprises/nagioscore/releases/latest" \
    | grep '"browser_download_url":' \
    | grep -o 'https://[^"]*')

tar xzf nagioscore.tar.gz
```

> O diretório extraído será algo como `nagios-4.x.x/`. Ajuste o nome abaixo conforme a versão baixada.

```bash
cd nagios-4.*.*/
```

---

## 3. Compilar e instalar o Nagios Core

```bash
sudo ./configure --with-httpd-conf=/etc/nginx/sites-enabled

sudo make all

sudo make install-groups-users
sudo usermod -a -G nagios www-data

sudo make install
sudo make install-daemoninit
sudo make install-commandmode
sudo make install-config
```

---

## 4. Configurar o Nginx

Crie o arquivo de configuração do virtual host:

```bash
sudo nano /etc/nginx/sites-available/nagios.conf
```

Cole o conteúdo abaixo (ajuste `server_name` conforme seu hostname ou IP):

```nginx
server {
    server_name  nagios-server.srv.local;
    listen       80;
    root         /usr/local/nagios/share;

    access_log   /var/log/nginx/nagios.access.log;
    error_log    /var/log/nginx/nagios.error.log;

    auth_basic            "Nagios Auth";
    auth_basic_user_file  /usr/local/nagios/etc/htpasswd.users;

    index index.php index.html index.htm;

    rewrite ^/nagios/(.*) /$1;

    location / {
        try_files $uri $uri/ index.php;
    }

    location ~ ^/?(.*\.php)$ {
        try_files       $uri =404;
        fastcgi_pass    unix:/run/php/php8.3-fpm.sock;
        include         fastcgi.conf;
    }

    location ~ \.cgi$ {
        fastcgi_param   AUTH_USER    $remote_user;
        fastcgi_param   REMOTE_USER  $remote_user;
        include         fastcgi.conf;
        fastcgi_pass    unix:/run/fcgiwrap.socket;
    }
}
```

Ative o site e recarregue o Nginx:

```bash
sudo ln -s /etc/nginx/sites-available/nagios.conf /etc/nginx/sites-enabled/nagios.conf

# Remover o site padrão (opcional)
sudo rm -f /etc/nginx/sites-enabled/default

sudo nginx -t
sudo systemctl reload nginx
```

---

## 5. Criar usuário de autenticação

```bash
sudo sh -c "echo -n 'nagiosadmin:' >> /usr/local/nagios/etc/htpasswd.users"
sudo sh -c "openssl passwd -apr1 >> /usr/local/nagios/etc/htpasswd.users"
```

> Será solicitada a senha do usuário `nagiosadmin`. Anote-a.

Criar o symlink dos CGI:

```bash
sudo ln -s /usr/local/nagios/sbin /usr/local/nagios/share/cgi-bin
```

---

## 6. Baixar e instalar os plugins do Nagios

```bash
mkdir -p ~/nagios-plugin && cd ~/nagios-plugin

wget -O nagios-plugins.tar.gz \
  $(curl -s "https://api.github.com/repos/nagios-plugins/nagios-plugins/releases/latest" \
    | grep '"browser_download_url":' \
    | grep -o 'https://[^"]*')

tar xzf nagios-plugins.tar.gz

cd nagios-plugins-*/

sudo ./configure \
  --with-nagios-user=nagios \
  --with-nagios-group=nagios

sudo make && sudo make install
```

---

## 7. Configurar o contato padrão (e-mail)

Edite o arquivo de contatos e altere o endereço de e-mail:

```bash
sudo nano /usr/local/nagios/etc/objects/contacts.cfg
```

Localize a linha `email` e substitua pelo seu endereço:

```
email    seuemail@exemplo.com
```

---

## 8. Validar configuração e iniciar o serviço

```bash
sudo /usr/local/nagios/bin/nagios -v /usr/local/nagios/etc/nagios.cfg
```

Se não houver erros:

```bash
sudo systemctl start nagios
sudo systemctl enable nagios
sudo systemctl status nagios
```

---

## 9. Configurar monitoramento de hosts

### Criar diretório para os hosts

```bash
sudo mkdir -p /usr/local/nagios/etc/objects/servers

sudo sed -i '0,/^#cfg_dir/ { /^#cfg_dir/ a\cfg_dir=/usr/local/nagios/etc/objects/servers}' \
  /usr/local/nagios/etc/nagios.cfg
```

### Exemplo de arquivo de host (`hosts.cfg`)

```bash
sudo nano /usr/local/nagios/etc/objects/servers/hosts.cfg
```

```nagios
define host {
    use                     linux-server
    host_name               meu-servidor
    alias                   Meu Servidor
    address                 192.168.1.100
    max_check_attempts      5
    check_period            24x7
    notification_interval   30
    notification_period     24x7
}

define hostgroup {
    hostgroup_name  linux-servers
    alias           Servidores Linux
    members         meu-servidor
}

define service {
    use                     generic-service
    host_name               meu-servidor
    service_description     PING
    check_command           check_ping!100.0,20%!500.0,60%
}

define service {
    use                     generic-service
    host_name               meu-servidor
    service_description     SSH
    check_command           check_ssh
    notifications_enabled   0
}
```

### Recarregar após alterações

```bash
sudo /usr/local/nagios/bin/nagios -v /usr/local/nagios/etc/nagios.cfg
sudo systemctl reload nagios
```

---

## 10. Acessar a interface web

Abra o navegador e acesse:

```
http://<IP-DO-SERVIDOR>/
```

Faça login com:
- **Usuário:** `nagiosadmin`
- **Senha:** a senha definida no passo 5

---

## Comandos úteis

| Comando | Descrição |
|---|---|
| `sudo systemctl start nagios` | Inicia o Nagios |
| `sudo systemctl stop nagios` | Para o Nagios |
| `sudo systemctl restart nagios` | Reinicia o Nagios |
| `sudo systemctl reload nagios` | Recarrega a configuração |
| `sudo systemctl status nagios` | Exibe o status do serviço |
| `sudo /usr/local/nagios/bin/nagios -v /usr/local/nagios/etc/nagios.cfg` | Valida a configuração |
| `sudo tail -f /usr/local/nagios/var/nagios.log` | Acompanha os logs em tempo real |

---

## Firewall (UFW)

Se o UFW estiver ativo, libere a porta HTTP:

```bash
sudo ufw allow 'Nginx HTTP'
sudo ufw reload
sudo ufw status
```

---

## Estrutura de diretórios importantes

| Caminho | Descrição |
|---|---|
| `/usr/local/nagios/etc/` | Arquivos de configuração |
| `/usr/local/nagios/etc/objects/` | Hosts, serviços, contatos |
| `/usr/local/nagios/etc/objects/servers/` | Hosts monitorados |
| `/usr/local/nagios/var/` | Logs e dados de estado |
| `/usr/local/nagios/share/` | Arquivos da interface web |
| `/usr/local/nagios/sbin/` | CGIs da interface web |
| `/usr/local/nagios/libexec/` | Plugins instalados |
