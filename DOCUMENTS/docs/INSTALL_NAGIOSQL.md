# Instalação do NagiosQL 3.5 com Nginx no Ubuntu 24.04

Interface web para gerenciamento de configurações do Nagios Core 4.

> **Pré-requisito:** Nagios Core 4 já instalado e funcionando conforme `INSTALL_NAGIOSCORE4.md`.

---

## 1. Instalar dependências

### MariaDB

```bash
sudo apt-get update
sudo apt-get install -y mariadb-server mariadb-client

sudo systemctl start mariadb
sudo systemctl enable mariadb

sudo mysql_secure_installation
```

### PHP 8.3 e módulos necessários

```bash
sudo apt-get install -y \
  php8.3-fpm \
  php8.3-cli \
  php8.3-common \
  php8.3-mbstring \
  php8.3-mysql \
  php8.3-gd \
  php8.3-curl \
  php8.3-xml \
  php8.3-zip \
  php8.3-gettext \
  php8.3-ftp
```

Verificar se os módulos foram carregados:

```bash
php8.3 -m | grep -E "mysqli|gettext|session|gd|ftp"
```

---

## 2. Criar banco de dados e usuário

```bash
sudo mysql -u root -p
```

```sql
CREATE DATABASE nagiosql CHARACTER SET utf8 COLLATE utf8_general_ci;
CREATE USER 'nagiosql'@'localhost' IDENTIFIED BY 'SuaSenhaForte123!';
GRANT ALL PRIVILEGES ON nagiosql.* TO 'nagiosql'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

---

## 3. Implantar os arquivos do NagiosQL

Copie os arquivos do NagiosQL para o diretório web:

```bash
sudo mkdir -p /var/www/nagiosql

# Se estiver usando o tar.gz do repositório local:
sudo tar -xzf nagiosql-3.5.0-git2023-06-18.tar.gz -C /var/www/nagiosql --strip-components=1

# Ou copie diretamente do diretório nagiosql/:
sudo cp -r nagiosql/. /var/www/nagiosql/
```

Ajustar owner e permissões:

```bash
sudo chown -R www-data:www-data /var/www/nagiosql
sudo chmod -R 755 /var/www/nagiosql

# O diretório config precisa ser gravável para salvar settings.php
sudo chmod -R 775 /var/www/nagiosql/config
```

---

## 4. Permissões para o Nagios Core

O usuário `www-data` precisa ter acesso de leitura/escrita nos diretórios de configuração do Nagios:

```bash
sudo usermod -a -G nagios www-data

sudo chown -R nagios:nagios /usr/local/nagios/etc
sudo chmod -R 775 /usr/local/nagios/etc

# Permitir que www-data execute o binário de verificação do Nagios
sudo chmod o+x /usr/local/nagios/bin/nagios
```

Reiniciar o PHP-FPM para aplicar o novo grupo:

```bash
sudo systemctl restart php8.3-fpm
```

---

## 5. Configurar o Nginx

```bash
sudo nano /etc/nginx/sites-available/nagiosql.conf
```

```nginx
server {
    listen       80;
    server_name  nagiosql.srv.local;

    root         /var/www/nagiosql;
    index        index.php index.html;

    access_log   /var/log/nginx/nagiosql.access.log;
    error_log    /var/log/nginx/nagiosql.error.log;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location ~ \.php$ {
        try_files       $uri =404;
        fastcgi_pass    unix:/run/php/php8.3-fpm.sock;
        fastcgi_index   index.php;
        include         fastcgi.conf;
    }

    # Bloquear acesso externo ao diretório install após a instalação
    location ^~ /install {
        allow 127.0.0.1;
        deny  all;
    }

    # Bloquear acesso a arquivos de configuração sensíveis
    location ~* /config/settings\.php {
        deny all;
        return 404;
    }
}
```

Ativar o site e testar:

```bash
sudo ln -s /etc/nginx/sites-available/nagiosql.conf /etc/nginx/sites-enabled/nagiosql.conf

sudo nginx -t
sudo systemctl reload nginx
```

---

## 6. Executar o assistente de instalação

Abra o navegador e acesse:

```
http://nagiosql.srv.local/install/install.php
```

### Passo 1 — Verificação de requisitos

O assistente verificará automaticamente:
- Versão do PHP (requer 7.2+)
- Módulos: `mysqli`, `session`, `gettext`, `filter`, `gd`
- Permissão de escrita em `config/`

Corrija qualquer item marcado como falha antes de prosseguir.

### Passo 2 — Configuração do banco de dados e admin

Preencha os campos:

| Campo | Valor |
|---|---|
| Database Type | `mysql` |
| Database Server | `localhost` |
| Database Server Port | `3306` |
| Database name | `nagiosql` |
| NagiosQL DB User | `nagiosql` |
| NagiosQL DB Password | `SuaSenhaForte123!` |
| Administrative DB User | `root` |
| Administrative DB Password | *(senha root do MariaDB)* |
| NagiosQL Admin Username | `admin` |
| NagiosQL Admin Password | *(escolha uma senha forte)* |

### Passo 3 — Finalizar instalação

O assistente irá:
1. Conectar ao banco de dados com o usuário administrativo
2. Criar o banco `nagiosql` (ou recriar se já existir)
3. Importar o schema (`nagiosQL_v35_db_mysql.sql`)
4. Criar o usuário `nagiosql` no banco
5. Gravar o arquivo `config/settings.php`

Ao concluir, clique em **"Go to NagiosQL"**.

---

## 7. Remover o diretório de instalação

Por segurança, remova o diretório de instalação após concluir:

```bash
sudo rm -rf /var/www/nagiosql/install
```

Remova também o bloqueio do Nginx adicionado no passo 5, já que o diretório não existe mais:

```bash
# Edite /etc/nginx/sites-available/nagiosql.conf e remova o bloco:
# location ^~ /install { ... }

sudo nginx -t && sudo systemctl reload nginx
```

---

## 8. Configuração pós-instalação

### 8.1 Acessar a interface

```
http://nagiosql.srv.local/
```

Login com `admin` e a senha definida no passo 6.

### 8.2 Verificar configurações gerais

Acesse **Administration → Settings** e confirme:

- **Nagios binary:** `/usr/local/nagios/bin/nagios`
- **Nagios configuration file:** `/usr/local/nagios/etc/nagios.cfg`
- **NagiosQL configuration path:** `/usr/local/nagios/etc/nagiosql/`

### 8.3 Configurar o domínio (Domain)

Acesse **Administration → Domains** e configure:

| Campo | Valor |
|---|---|
| Domain name | `Local Nagios` |
| Config files base directory | `/usr/local/nagios/etc/nagiosql/` |
| Nagios binary | `/usr/local/nagios/bin/nagios` |
| Nagios config file | `/usr/local/nagios/etc/nagios.cfg` |
| Transfer method | `local` |

Crie o diretório base de configuração do NagiosQL:

```bash
sudo mkdir -p /usr/local/nagios/etc/nagiosql/{hosts,services}
sudo chown -R nagios:nagios /usr/local/nagios/etc/nagiosql
sudo chmod -R 775 /usr/local/nagios/etc/nagiosql
```

### 8.4 Atualizar o nagios.cfg

Substitua as entradas `cfg_dir` e `cfg_file` do Nagios Core pelo bloco gerenciado pelo NagiosQL.

Edite `/usr/local/nagios/etc/nagios.cfg`:

```bash
sudo nano /usr/local/nagios/etc/nagios.cfg
```

Substitua (ou adicione após comentar as linhas originais):

```ini
# Configurações gerenciadas pelo NagiosQL
cfg_file=/usr/local/nagios/etc/nagiosql/timeperiods.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/commands.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/contacts.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/contactgroups.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/contacttemplates.cfg

cfg_dir=/usr/local/nagios/etc/nagiosql/hosts
cfg_file=/usr/local/nagios/etc/nagiosql/hosttemplates.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/hostgroups.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/hostextinfo.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/hostescalations.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/hostdependencies.cfg

cfg_dir=/usr/local/nagios/etc/nagiosql/services
cfg_file=/usr/local/nagios/etc/nagiosql/servicetemplates.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/servicegroups.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/serviceextinfo.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/serviceescalations.cfg
cfg_file=/usr/local/nagios/etc/nagiosql/servicedependencies.cfg
```

Valide e reinicie:

```bash
sudo /usr/local/nagios/bin/nagios -v /usr/local/nagios/etc/nagios.cfg
sudo systemctl restart nagios
```

---

## 9. Verificar a instalação

Acesse **Administration → Support** no NagiosQL. Esta página verifica:

- Conectividade com o banco de dados
- Permissões nos diretórios do Nagios
- Execução do binário `nagios -v`
- Módulos PHP necessários

Todos os itens devem aparecer como **OK**.

---

## 10. (Opcional) Instalar extensão SSH para deploy remoto

Caso precise gerenciar um servidor Nagios remoto via SSH:

```bash
sudo apt-get install -y libssh2-1-dev php8.3-dev php-pear make

sudo pecl install ssh2-1.3.1
```

Adicionar a extensão ao PHP:

```bash
echo "extension=ssh2.so" | sudo tee /etc/php/8.3/fpm/conf.d/30-ssh2.ini

sudo systemctl restart php8.3-fpm
```

### Configurar chave SSH para o servidor remoto

```bash
sudo mkdir -p /etc/nagiosql/ssh

sudo ssh-keygen -t rsa -m PEM -b 4096 -f /etc/nagiosql/ssh/id_rsa
# Deixe a passphrase vazia (Enter)

# Copiar a chave pública para o servidor remoto
sudo ssh-copy-id -i /etc/nagiosql/ssh/id_rsa.pub nagiosql_usr@servidor-remoto

# Ajustar permissões para o www-data
sudo chown www-data:www-data /etc/nagiosql/ssh/id_rsa
sudo chmod 600 /etc/nagiosql/ssh/id_rsa
```

No NagiosQL, edite o domínio (**Administration → Domains**):
- **Transfer method:** `SSH/SFTP`
- **User:** `nagiosql_usr`
- **SSH key directory:** `/etc/nagiosql/ssh/`

---

## 11. (Opcional) Importar configurações existentes do Nagios

Se o Nagios já possuía configurações manuais, importe-as via NagiosQL:

1. Acesse **Administration → Import**
2. Selecione o domínio configurado
3. Clique em **"Import"** — o NagiosQL lerá os arquivos `.cfg` existentes e populará o banco de dados

Após a importação, gere as configurações pelo NagiosQL para garantir consistência:

1. Acesse **Monitoring → Write Nagios Config**
2. Clique em **"Write Config"**
3. Valide e reinicie: `sudo systemctl restart nagios`

---

## Referências de diretórios

| Caminho | Descrição |
|---|---|
| `/var/www/nagiosql/` | Arquivos da aplicação NagiosQL |
| `/var/www/nagiosql/config/settings.php` | Configurações geradas pelo instalador |
| `/usr/local/nagios/etc/nagiosql/` | Configs do Nagios gerenciadas pelo NagiosQL |
| `/etc/nagiosql/ssh/` | Chaves SSH para deploy remoto (opcional) |
| `/var/log/nginx/nagiosql.*.log` | Logs do Nginx |

## Comandos úteis

| Comando | Descrição |
|---|---|
| `sudo systemctl restart php8.3-fpm` | Reinicia o PHP-FPM |
| `sudo systemctl reload nginx` | Recarrega o Nginx |
| `sudo systemctl restart nagios` | Reinicia o Nagios |
| `sudo mysql -u nagiosql -p nagiosql` | Acessa o banco do NagiosQL |
| `php8.3 -m \| grep -E "mysqli\|ssh2\|gettext"` | Lista módulos PHP ativos |
