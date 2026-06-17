# Nagios Core 4 + NagiosQL 3.5 — Stack Docker Compose

Stack completo de monitoramento com **Nagios Core 4.5.13**, **NagiosQL 3.5.0** e **MariaDB 10.11**, rodando em **2 containers** Docker sobre Debian trixie-slim.

> Para entender o fluxo interno de como o NagiosQL gera configs e o Nagios as carrega, consulte o [DIAGRAMA_NAGIOSQL.md](docs/DIAGRAMA_NAGIOSQL.md).

---

## Sumário

- [Visão geral](#visão-geral)
- [Pré-requisitos](#pré-requisitos)
- [Estrutura do projeto](#estrutura-do-projeto)
- [Configuração inicial](#configuração-inicial)
- [Subindo o stack](#subindo-o-stack)
- [Acessando as interfaces](#acessando-as-interfaces)
- [Credenciais padrão](#credenciais-padrão)
- [Volumes e dados persistentes](#volumes-e-dados-persistentes)
- [Gerenciamento do stack](#gerenciamento-do-stack)
- [Adicionando hosts e serviços](#adicionando-hosts-e-serviços)
- [Reload do Nagios Core](#reload-do-nagios-core)
- [Rebuild das imagens](#rebuild-das-imagens)
- [Solução de problemas](#solução-de-problemas)
- [Modelo de permissões e segurança](#modelo-de-permissões-e-segurança)
- [Arquitetura interna do container](#arquitetura-interna-do-container)

---

## Visão geral

| Componente | Versão | Porta host |
|---|---|---|
| Nagios Core | 4.5.13 | 8080 |
| NagiosQL | 3.5.0 | 8081 |
| MariaDB | 10.11 | (interno) |
| Nginx | 1.26 | — |
| PHP | 8.4 | — |

**Nagios Core** e **NagiosQL** rodam no **mesmo container** (`nagios-core`), gerenciados pelo supervisord. Essa é a arquitetura natural do NagiosQL, que foi projetado para rodar na mesma máquina que o Nagios Core: acessa os arquivos `.cfg` e o binário diretamente, sem volumes compartilhados ou workarounds de rede.

O MariaDB fica em um container separado (`nagios-db`), acessado pelo NagiosQL via hostname `db`.

---

## Pré-requisitos

- Docker Engine 24+
- `docker-compose` standalone v2+ **ou** plugin `docker compose`
- Usuário com permissão para executar `docker`
- Portas 8080 e 8081 livres no host

Verificar:

```bash
docker --version
docker-compose --version   # standalone
# ou
docker compose version     # plugin
```

---

## Estrutura do projeto

```
projeto_base/
├── build.sh                    # Constrói a imagem nagios-core:latest (sem subir containers)
├── docker/
│   └── nagios-core/
│       ├── docker-compose.yml  # Define os serviços nagios e db
│       ├── nagios/
│       │   ├── Dockerfile          # Multi-stage: builder (compila + .deb) → runtime (instala .deb)
│       │   ├── build-debs.sh       # Compila Nagios Core, Plugins e empacota NagiosQL como .deb
│       │   ├── entrypoint.sh       # Inicializa volumes, importa schema, inicia supervisord
│       │   ├── nginx.conf          # Dois server blocks: :80 (Nagios Core) e :8081 (NagiosQL)
│       │   ├── supervisord.conf    # Gerencia 5 processos: nginx, php-fpm, fcgiwrap, nagios, reload-watcher
│       │   ├── reload-watcher.sh   # Monitora reload.trigger, valida nagios.cfg e envia SIGHUP ao Nagios
│       │   └── etc-extra/          # Configs de exemplo embutidas na imagem
│       │       └── nagiosql/
│       │           ├── commands.cfg         # check_dns customizado
│       │           ├── hosts/               # 4 hosts de exemplo (gateway, dns, linux-host)
│       │           └── services/            # Serviços PING, HTTP, SSH, DNS
│       ├── nagiosql/           # Código-fonte NagiosQL 3.5.0 (contexto de build)
│       ├── .env.example        # Modelo de variáveis de ambiente
│       └── volumes/            # Dados persistentes (bind mounts — não versionar)
│           ├── db-data/                # Dados do MariaDB
│           ├── nagios-etc/             # Configs do Nagios (/usr/local/nagios/etc)
│           ├── nagios-var/             # Runtime do Nagios (status.dat, logs, spool)
│           ├── nagios-plugins/         # Plugins compilados (/usr/local/nagios/libexec)
│           └── nagiosql-config/        # Config do NagiosQL (settings.php, fieldvars.php, CSS, locale)
├── docs/
│   ├── DIAGRAMA_NAGIOSQL.md        # Fluxo completo NagiosQL → Nagios Core
│   ├── BUILD_NAGIOS_TRIXIE.md      # Roteiro de build manual no Debian trixie
│   ├── BUILD_NAGIOS_TRIXIE_COMPLETE.md
│   ├── INSTALL_NAGIOSCORE4.md      # Instalação do Nagios Core 4 no Ubuntu 24.04 com Nginx
│   ├── INSTALL_NAGIOSQL.md
│   ├── NAGIOS_BUILD_COMMANDS.md    # Comandos de compilação passo a passo
│   └── NAGIOSQL_DEBIAN_PACKAGING.md
├── .gitignore
└── README.md                   # Este arquivo
```

---

## Configuração inicial

### 1. Copiar e editar o arquivo de variáveis de ambiente

```bash
cp docker/nagios-core/.env.example docker/nagios-core/.env
nano docker/nagios-core/.env   # ou seu editor preferido
```

Variáveis disponíveis:

```ini
# MariaDB
MYSQL_ROOT_PASSWORD=rootpassword_mude_aqui
MYSQL_DATABASE=nagiosql
MYSQL_USER=nagiosql
MYSQL_PASSWORD=nagiosqlpass_mude_aqui
DB_PORT=3306

# Nagios Core — senha do usuário nagiosadmin na interface web
NAGIOS_ADMIN_PASSWORD=nagiosadmin_mude_aqui

# NagiosQL — usuário e senha do admin da interface web
NAGIOSQL_USER=admin
NAGIOSQL_PASSWORD=admin_mude_aqui

# Portas publicadas no host
NAGIOS_PORT=8080
NAGIOSQL_PORT=8081
```

> **Atenção:** Troque todas as senhas antes de usar em produção. O arquivo `.env` está no `.gitignore` e nunca deve ser versionado.

### 2. Verificar se as portas estão livres

```bash
ss -tlnp | grep -E '8080|8081'
```

---

---

## Subindo o stack

### Build e start (primeira vez)

```bash
cd docker/nagios-core
docker-compose up -d --build
```

O processo:
1. **Stage builder** — compila Nagios Core 4.5.13 e Plugins 2.5 do fonte, empacota os três componentes como `.deb` via `build-debs.sh` (pode levar 3–5 min)
2. **Stage runtime** — parte do `debian:trixie-slim` limpo, instala apenas dependências de execução e os `.deb` gerados; ferramentas de compilação não entram na imagem final
3. Sobe o MariaDB e aguarda o healthcheck
4. Sobe o container `nagios-core`, que no entrypoint:
   - Inicializa os volumes de configs, plugins e runtime
   - Aguarda a porta 3306 do MariaDB via `nc`
   - Importa o schema SQL, os dados de exemplo e configura os caminhos (primeira execução)
   - Cria o usuário admin do NagiosQL
   - Gera o `settings.php` com as credenciais do `.env`
   - Valida o `nagios.cfg` e inicia o supervisord

Acompanhar logs em tempo real:

```bash
# a partir de docker/nagios-core/
docker-compose logs -f
```

### Verificar status dos containers

```bash
# a partir de docker/nagios-core/
docker-compose ps
```

Saída esperada:

```
NAME          IMAGE                 STATUS
nagios-core   projeto_base-nagios   Up X minutes
nagios-db     mariadb:10.11         Up X minutes (healthy)
```

---

## Acessando as interfaces

| Interface | URL | Descrição |
|---|---|---|
| Nagios Core | http://localhost:8080 | Monitoramento (redireciona para tac.cgi) |
| NagiosQL | http://localhost:8081 | Gerenciamento de configurações |

> Se estiver acessando de outra máquina, substitua `localhost` pelo IP do host Docker.

---

## Credenciais padrão

### Nagios Core (autenticação HTTP Basic)

| Campo | Valor padrão |
|---|---|
| Usuário | `nagiosadmin` |
| Senha | valor de `NAGIOS_ADMIN_PASSWORD` no `.env` (padrão: `nagiosadmin`) |

### NagiosQL (formulário de login)

| Campo | Valor padrão |
|---|---|
| Usuário | valor de `NAGIOSQL_USER` no `.env` (padrão: `admin`) |
| Senha | valor de `NAGIOSQL_PASSWORD` no `.env` (padrão: `admin`) |

> As senhas padrão são definidas no `docker/nagios-core/.env.example`. Caso o `.env` não exista, o stack usa os fallbacks codificados no `docker-compose.yml`.

---

## Volumes e dados persistentes

Todos os dados ficam em subdiretórios de `docker/nagios-core/volumes/` (bind mounts locais):

| Diretório | Caminho no container | Conteúdo |
|---|---|---|
| `docker/nagios-core/volumes/db-data` | `/var/lib/mysql` | Banco de dados MariaDB completo |
| `docker/nagios-core/volumes/nagios-etc` | `/usr/local/nagios/etc` | Arquivos `.cfg` do Nagios, htpasswd |
| `docker/nagios-core/volumes/nagios-var` | `/usr/local/nagios/var` | status.dat, logs, spool, nagios.lock |
| `docker/nagios-core/volumes/nagios-plugins` | `/usr/local/nagios/libexec` | Plugins compilados (check_ping, check_http…) |
| `docker/nagios-core/volumes/nagiosql-config` | `/var/www/nagiosql/config` | settings.php, fieldvars.php, CSS, locale |

Os volumes são inicializados automaticamente no primeiro boot:

- **nagios-etc:** copiado de `/usr/local/nagios/etc.default`, que inclui as configs de exemplo de `etc-extra/` (hosts, serviços, comandos)
- **nagios-plugins:** copiado de `/usr/local/nagios/libexec.default`; `check_ping` e `check_icmp` têm setuid root (necessário para raw sockets ICMP); inclui `check_fping`, `check_snmp`, `check_mysql` e `check_mysql_query` além dos plugins padrão
- **nagios-var:** diretórios `rw/`, `spool/checkresults/` e `archives/` criados pelo entrypoint
- **nagiosql-config:** inicializado a partir de `config.default/` (preserva `fieldvars.php`, CSS, locale); `settings.php` gerado com as variáveis do `.env`
- **db-data:** schema importado, dados de exemplo pré-carregados (24 comandos, 5 timeperiods, templates `linux-server`/`generic-service`, 4 hosts, 21 serviços), caminhos configurados e usuário admin criado na primeira execução

### Backup

```bash
# Backup das configs do Nagios
tar -czf backup-nagios-etc-$(date +%Y%m%d).tar.gz docker/nagios-core/volumes/nagios-etc/

# Backup do banco de dados
docker exec nagios-db \
    mysqldump -u root -p"${MYSQL_ROOT_PASSWORD}" nagiosql \
    > backup-nagiosql-$(date +%Y%m%d).sql
```

---

## Gerenciamento do stack

Todos os comandos `docker-compose` devem ser executados a partir do diretório `docker/nagios-core/`.

### Parar o stack (sem remover dados)

```bash
cd docker/nagios-core
docker-compose stop
```

### Iniciar o stack parado

```bash
cd docker/nagios-core
docker-compose start
```

### Parar e remover containers (volumes preservados)

```bash
cd docker/nagios-core
docker-compose down
```

### Destruir tudo, incluindo os dados

```bash
cd docker/nagios-core
docker-compose down
rm -rf ./volumes/
```

> **Atenção:** Sempre execute `docker-compose down` **antes** de apagar `./volumes/`. Remover `db-data/` enquanto o MariaDB está rodando deixa o banco num estado corrompido em memória, causando falhas na próxima inicialização. Os dados são apagados irreversivelmente.

### Ver logs de um container específico

```bash
# a partir de docker/nagios-core/
docker-compose logs -f nagios
docker-compose logs -f db
```

### Abrir shell em um container

```bash
docker exec -it nagios-core bash
docker exec -it nagios-db bash
```

---

## Adicionando hosts e serviços

O stack já sobe com **5 hosts e 16 serviços** de exemplo configurados (gateway, google-dns, cloudflare-dns, linux-host, localhost). Para adicionar novos, há duas formas:

### Via NagiosQL (recomendado)

Acesse http://localhost:8081, faça login e use os menus:

- **Monitoring → Hosts** — cadastrar hosts
- **Monitoring → Services** — cadastrar serviços
- **Monitoring → Write Config Files** — gerar os arquivos `.cfg`
- **Monitoring → Verify Configuration** — validar antes de aplicar
- **Monitoring → Restart Nagios** — recarregar o Nagios Core

O banco já vem com templates prontos para uso: `linux-server`, `windows-server`, `generic-host`, `generic-service`, `local-service`, além de 24 comandos e 5 timeperiods. Ao cadastrar um novo host basta selecionar o template desejado.

> Para entender o fluxo completo de como os dados percorrem o NagiosQL até o monitoramento, consulte o [DIAGRAMA_NAGIOSQL.md](docs/DIAGRAMA_NAGIOSQL.md).

### Via arquivos `.cfg` diretamente

Os arquivos gerenciados pelo NagiosQL ficam em:

```
docker/nagios-core/volumes/nagios-etc/nagiosql/
├── hosts/          ← um arquivo .cfg por host
├── services/       ← um arquivo .cfg por serviço (ou grupo)
├── commands.cfg    ← comandos customizados
├── contacts.cfg    ← contatos
└── hosttemplates.cfg, servicetemplates.cfg, etc.
```

Exemplo de host (`docker/nagios-core/volumes/nagios-etc/nagiosql/hosts/meuhost.cfg`):

```nagios
define host {
    use                 linux-server
    host_name           meuhost
    alias               Meu Servidor Linux
    address             192.168.1.10
    max_check_attempts  3
    check_interval      5
    retry_interval      1
    contact_groups      admins
}
```

Exemplo de serviço (`docker/nagios-core/volumes/nagios-etc/nagiosql/services/ping.cfg`):

```nagios
define service {
    use                  local-service
    host_name            meuhost
    service_description  PING
    check_command        check_ping!100.0,20%!500.0,60%
    max_check_attempts   3
    contact_groups       admins
}
```

Após criar ou editar arquivos, **valide e recarregue** conforme a seção abaixo.

---

## Reload do Nagios Core

### Via NagiosQL (recomendado)

Use o menu **Monitoring → Restart Nagios**. O `reload-watcher.sh` intercepta a solicitação, **valida o `nagios.cfg` automaticamente** antes de enviar o SIGHUP. Se a configuração contiver erros, o reload é cancelado e o Nagios continua operando com a configuração anterior.

Para ver o resultado da validação:

```bash
docker logs nagios-core | grep reload-watcher
```

Saída esperada após um reload bem-sucedido:
```
[reload-watcher] Config válido. SIGHUP enviado ao PID 51
```

Se houver erro de configuração:
```
[reload-watcher] ERRO: nagios.cfg inválido — reload cancelado.
```

### Validar a configuração manualmente

```bash
docker exec nagios-core su -s /bin/bash nagios -c \
    "/usr/local/nagios/bin/nagios -v /usr/local/nagios/etc/nagios.cfg"
```

Saída esperada no final:
```
Total Warnings: 0
Total Errors:   0
Things look okay - No serious problems were detected during the pre-flight check
```

### Recarregar manualmente (sem downtime)

```bash
docker exec nagios-core supervisorctl signal HUP nagios
```

O reload é feito via SIGHUP: o Nagios Core relê todos os arquivos `.cfg` sem interromper os checks em execução.

---

## Rebuild das imagens

O Dockerfile usa **multi-stage build**: o stage `builder` compila e gera os `.deb`; o stage `runtime` apenas instala esses pacotes. Ferramentas de compilação (gcc, make, etc.) não entram na imagem final.

O `build-debs.sh` é executado automaticamente pelo builder e gera:

| Pacote | Tamanho |
|---|---|
| `nagios-core_4.5.13_amd64.deb` | ~1.7 MB |
| `nagios-plugins_2.5_amd64.deb` | ~900 KB |
| `nagiosql_3.5.0_all.deb` | ~5.3 MB |

### Rebuild após alterar o Dockerfile, build-debs.sh ou configs

```bash
# A partir de docker/nagios-core/

# Rebuild sem cache (garante atualização total)
docker-compose build --no-cache

# Rebuild com cache (reutiliza camadas inalteradas — mais rápido)
docker-compose build

# Alternativa: script simples na raiz do projeto — gera apenas a imagem (sem subir containers)
./build.sh             # com cache
./build.sh --no-cache  # sem cache
```

O Docker cache funciona por camada: se apenas o `build-debs.sh` mudar, o `apt-get install` do builder é reutilizado. Se apenas o código do NagiosQL mudar, as etapas de compilação do Nagios Core e Plugins são reutilizadas.

### Aplicar rebuild sem parar o banco de dados

```bash
# a partir de docker/nagios-core/
docker-compose up -d --no-deps --build nagios
```

> **Atenção:** Ao rebuildar o `nagios`, os volumes **não são apagados**. Para forçar reinicialização dos plugins (ex.: após adicionar `check_snmp`, `check_mysql` ou alterar a versão compilada), limpe o volume e reinicie:
>
> ```bash
> # a partir de docker/nagios-core/
> docker-compose stop nagios
> docker run --rm \
>   -v $(pwd)/volumes/nagios-plugins:/vol \
>   debian:trixie-slim sh -c "find /vol -mindepth 1 -delete"
> docker-compose up -d
> ```

---

## Solução de problemas

### Container nagios-core não sobe

Verificar os logs:

```bash
docker-compose logs nagios | grep -E "Error|error|FATAL"
```

Causas comuns e soluções:

| Erro | Causa | Solução |
|---|---|---|
| `Check result path ... is not a valid directory` | Volume `nagios-var` vazio, faltam subdiretórios | O entrypoint cria automaticamente; reiniciar o container |
| `Could not read object configuration data` | Arquivo `.cfg` com sintaxe inválida | Validar com `nagios -v nagios.cfg` |
| `Could not open pipe` no check_ping | `iputils-ping` não instalado no build | Rebuild com `--no-cache` |

### Nagios Core mostra hosts DOWN com "You need more args!!!"

O plugin `check_ping` precisa do setuid root para abrir raw sockets. Verificar:

```bash
# a partir de docker/nagios-core/
ls -la ./volumes/nagios-plugins/check_ping
# deve mostrar: -r-sr-xr-x 1 root nagios
```

Se o owner for `nagios` em vez de `root`, limpe o volume de plugins e faça rebuild:

```bash
# a partir de docker/nagios-core/
docker-compose stop nagios
docker run --rm \
  -v $(pwd)/volumes/nagios-plugins:/vol \
  debian:trixie-slim sh -c "find /vol -mindepth 1 -delete"
docker-compose up -d
```

### NagiosQL redireciona para /install/index.php

O registro de versão no banco não existe. Inserir manualmente:

```bash
docker exec nagios-db mysql -u nagiosql -pnagiosqlpass nagiosql \
    -e "INSERT INTO tbl_settings (category, name, value) \
        VALUES ('db','version','3.5.0') \
        ON DUPLICATE KEY UPDATE value='3.5.0';"
```

### NagiosQL retorna erro 500 após login (admin.php)

O volume `nagiosql-config` está incompleto — falta o `fieldvars.php`. Limpar e reiniciar:

```bash
# a partir de docker/nagios-core/
docker-compose stop nagios
docker run --rm \
  -v $(pwd)/volumes/nagiosql-config:/vol \
  debian:trixie-slim sh -c "find /vol -mindepth 1 -delete"
docker-compose up -d
```

### Verificar processos dentro do container

```bash
docker exec nagios-core supervisorctl status
```

Saída esperada:

```
fcgiwrap                         RUNNING   pid XXXX, uptime 0:XX:XX
nagios                           RUNNING   pid XXXX, uptime 0:XX:XX
nginx                            RUNNING   pid XXXX, uptime 0:XX:XX
php-fpm                          RUNNING   pid XXXX, uptime 0:XX:XX
reload-watcher                   RUNNING   pid XXXX, uptime 0:XX:XX
```

### Testar check manualmente

```bash
# check_ping
docker exec nagios-core su -s /bin/bash nagios -c \
    "/usr/local/nagios/libexec/check_ping -H 8.8.8.8 -w 200.0,20% -c 500.0,60% -p 3"

# check_http
docker exec nagios-core su -s /bin/bash nagios -c \
    "/usr/local/nagios/libexec/check_http -H google.com"

# check_ssh
docker exec nagios-core su -s /bin/bash nagios -c \
    "/usr/local/nagios/libexec/check_ssh 192.168.1.9"
```

### NagiosQL "Restart Nagios" não aplicou as mudanças

O `reload-watcher.sh` valida o `nagios.cfg` antes de recarregar. Se a configuração contiver erros, o reload é silenciosamente cancelado. Verificar:

```bash
docker logs nagios-core | grep reload-watcher | tail -5
```

Se a saída mostrar `ERRO: nagios.cfg inválido`, valide manualmente para ver a mensagem de erro completa:

```bash
docker exec nagios-core su -s /bin/bash nagios -c \
    "/usr/local/nagios/bin/nagios -v /usr/local/nagios/etc/nagios.cfg"
```

Corrija os erros no NagiosQL, regenere os arquivos (**Monitoring → Write Config Files**) e tente novamente.

### Ver estatísticas do Nagios em tempo real

```bash
docker exec nagios-core /usr/local/nagios/bin/nagiostats
```

---

## Modelo de permissões e segurança

O stack usa dois grupos Unix para separar responsabilidades:

| Grupo | GID | Membros | Para que serve |
|---|---|---|---|
| `nagios` | 3000 | `nagios` | Grupo primário do daemon Nagios; controla acesso aos binários e runtime |
| `nagioscfg` | 3001 | `nagios`, `www-data` | Grupo compartilhado para escrita nos diretórios de configuração gerados pelo NagiosQL |

**Por que dois grupos?** O `www-data` (PHP/NagiosQL) precisa escrever nos diretórios `etc/nagiosql/` para gerar os arquivos `.cfg`. Adicioná-lo ao grupo `nagios` daria acesso excessivo (binários, socket de comandos externos). O grupo `nagioscfg` limita o acesso exatamente ao necessário.

O pipe de comandos externos (`var/rw/nagios.cmd`) usa o grupo `nagioscfg` via SGID no diretório `var/rw/`, o que permite ao NagiosQL enviar comandos ao Nagios quando necessário.

---

## Arquitetura interna do container

### Container `nagios-core`

Contém **Nagios Core 4.5.13** e **NagiosQL 3.5.0** no mesmo container, gerenciados pelo **supervisord** com 5 processos:

| Processo | Função |
|---|---|
| `nginx` | Dois server blocks: porta 80 (Nagios Core) e porta 8081 (NagiosQL) |
| `php-fpm` | Executa scripts PHP — tanto a UI do Nagios quanto o NagiosQL |
| `fcgiwrap` | Executa os CGIs do Nagios (tac.cgi, status.cgi, etc.) |
| `nagios` | Processo principal do Nagios Core |
| `reload-watcher` | Monitora `reload.trigger`, valida `nagios.cfg` e envia SIGHUP ao Nagios |

O mecanismo de reload funciona assim: o NagiosQL escreve no arquivo `reload.trigger` → o `reload-watcher.sh` detecta a mudança → **valida o `nagios.cfg` com `nagios -v`** → se válido, envia `SIGHUP` ao Nagios → o Nagios relê todas as configs sem downtime. Se a validação falhar, o reload é cancelado e o Nagios permanece com a configuração anterior, protegendo o monitoramento em produção.

Por rodarem no mesmo container, o NagiosQL acessa `/usr/local/nagios/` diretamente — sem volumes compartilhados entre containers, sem latência de rede, exatamente como foi projetado.

### Container `nagios-db`

MariaDB 10.11 dedicado ao NagiosQL. Expõe apenas a porta 3306 na rede interna `nagios-net`.

### Rede interna

```
Host
├── :8080 → nagios-core:80    (Nginx → Nagios CGIs via fcgiwrap)
└── :8081 → nagios-core:8081  (Nginx → NagiosQL PHP via php-fpm)

nagios-core ──── nagios-net ────→ nagios-db:3306
```

### Hosts e serviços de exemplo incluídos

A imagem já contém configurações iniciais em `docker/nagios-core/nagios/etc-extra/`, copiadas para o volume no primeiro boot:

| Host | Endereço | Serviços monitorados |
|---|---|---|
| `localhost` | 127.0.0.1 | PING, Disk, Load, Users, Procs, Swap, SSH, HTTP |
| `linux-host` | 192.168.1.9 | PING, HTTP, SSH |
| `gateway` | 192.168.1.1 | PING |
| `google-dns` | 8.8.8.8 | PING, DNS |
| `cloudflare-dns` | 1.1.1.1 | PING, DNS |

Para customizar os hosts de exemplo antes do primeiro boot, edite os arquivos em `docker/nagios-core/nagios/etc-extra/nagiosql/` e faça rebuild da imagem.

---

## Documentação adicional

| Arquivo | Descrição |
|---|---|
| [README.md](README.md) | Este guia de operação |
| [docs/DIAGRAMA_NAGIOSQL.md](docs/DIAGRAMA_NAGIOSQL.md) | Fluxo completo: formulário → banco → .cfg → reload → monitoramento |
| [docs/INSTALL_NAGIOSCORE4.md](docs/INSTALL_NAGIOSCORE4.md) | Instalação manual do Nagios Core 4 no Ubuntu 24.04 com Nginx |
| [docs/INSTALL_NAGIOSQL.md](docs/INSTALL_NAGIOSQL.md) | Instalação manual do NagiosQL |
| [docs/BUILD_NAGIOS_TRIXIE.md](docs/BUILD_NAGIOS_TRIXIE.md) | Roteiro de compilação manual no Debian trixie |
| [docs/BUILD_NAGIOS_TRIXIE_COMPLETE.md](docs/BUILD_NAGIOS_TRIXIE_COMPLETE.md) | Roteiro completo de build (incluindo plugins e NagiosQL) |
| [docs/NAGIOS_BUILD_COMMANDS.md](docs/NAGIOS_BUILD_COMMANDS.md) | Comandos de compilação passo a passo |
| [docs/NAGIOSQL_DEBIAN_PACKAGING.md](docs/NAGIOSQL_DEBIAN_PACKAGING.md) | Empacotamento do NagiosQL como .deb |
