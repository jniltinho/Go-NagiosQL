# Diagrama de Funcionamento: NagiosQL + Nagios Core

Este documento descreve o fluxo completo de como uma configuração criada no NagiosQL percorre o sistema até ser efetivamente monitorada pelo Nagios Core.

---

## Visão Geral da Arquitetura

NagiosQL e Nagios Core rodam no **mesmo container** (`nagios-core`), gerenciados pelo supervisord. Essa é a arquitetura natural do NagiosQL: acessa os arquivos `.cfg` e o binário diretamente, sem volumes compartilhados entre containers.

```
┌─────────────────────────────────────────────────────────────────────────┐
│  Host Docker                                                            │
│                                                                         │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │   Container: nagios-core                                         │  │
│  │                                                                  │  │
│  │   ┌───────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────┐  │  │
│  │   │   nginx   │  │ php-fpm  │  │ fcgiwrap │  │    nagios    │  │  │
│  │   │:80 → CGIs │  │:8081→PHP │  │  CGIs    │  │  Core daemon │  │  │
│  │   └───────────┘  └──────────┘  └──────────┘  └──────┬───────┘  │  │
│  │                                                      │          │  │
│  │   ┌──────────────────┐     ┌──────────────────────┐  │          │  │
│  │   │  NagiosQL (PHP)  │     │   reload-watcher.sh  │  │          │  │
│  │   │  Lê/Escreve BD   │     │   valida nagios.cfg  │  │          │  │
│  │   │  Gera .cfg files │     │   envia SIGHUP       ├──┘          │  │
│  │   └────────┬─────────┘     └──────────▲───────────┘             │  │
│  │            │                          │                         │  │
│  │   ┌────────▼──────────────────────────┴──────────────────────┐  │  │
│  │   │  ./volumes/nagios-etc/ (bind mount)                      │  │  │
│  │   │  nagiosql/hosts/*.cfg   ◄── NagiosQL escreve             │  │  │
│  │   │  nagiosql/services/*.cfg                                  │  │  │
│  │   │  objects/*.cfg          ◄── Nagios Core lê               │  │  │
│  │   │  nagios.cfg                                               │  │  │
│  │   │  var/reload.trigger     ◄── sinal de reload              │  │  │
│  │   └──────────────────────────────────────────────────────────┘  │  │
│  │                                                                  │  │
│  │   :8080 (Nagios Core)  ──────────────────────────────────────── │  │
│  │   :8081 (NagiosQL)     ──────────────────────────────────────── │  │
│  └──────────────────────────────────────────────────────────────────┘  │
│                                                                         │
│  ┌──────────────────────────────────────────┐                         │
│  │   Container: nagios-db (MariaDB 10.11)   │                         │
│  │   ./volumes/db-data/                     │                         │
│  │   tbl_host, tbl_service, tbl_contact     │                         │
│  └──────────────────────────────────────────┘                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Fluxo Completo: Do Formulário ao Monitoramento

```mermaid
flowchart TD
    A([Administrador\nacessa NagiosQL\n:8081]) --> B

    subgraph UI["NagiosQL — Interface Web (PHP)"]
        B[Preenche formulário\nHost / Serviço / Contato\nTemplate / Comando]
        B --> C{Salvar}
        C --> D[(MariaDB\ntbl_host\ntbl_service\ntbl_contact\ntbl_command\ntbl_timeperiod)]
    end

    D --> E

    subgraph GEN["Geração de Configuração"]
        E[Monitoring →\nWrite Config Files]
        E --> F[NagiosQL lê\nregistros do banco]
        F --> G[Gera arquivos .cfg\nnagiosql/hosts/meuhost.cfg\nnagiosql/services/ping.cfg\nnagiosql/commands.cfg]
    end

    G --> H

    subgraph VOL["Volume Compartilhado\n./volumes/nagios-etc/"]
        H[.cfg gravados no\nvolume compartilhado]
    end

    H --> I

    subgraph VER["Verificação"]
        I[Monitoring →\nVerify Configuration]
        I --> J["NagiosQL executa:\nnagios -v nagios.cfg\n(via SSH ou local)"]
        J --> K{Resultado}
        K -->|Total Errors: 0| L[✅ Configuração válida]
        K -->|Errors > 0| M[❌ Exibe erros\nCorrigir e regerar]
        M --> E
    end

    L --> N

    subgraph RELOAD["Reload do Nagios Core"]
        N[Monitoring →\nRestart Nagios]
        N --> O["NagiosQL escreve no\ncommandfile:\n/usr/local/nagios/var/reload.trigger"]
        O --> P[reload-watcher.sh\ndetecta o arquivo]
        P --> P2{"nagios -v nagios.cfg\n(validação prévia)"}
        P2 -->|válido| Q[Envia SIGHUP\nao processo nagios\n via pgrep]
        P2 -->|inválido| QX[❌ Reload cancelado\nNagios continua com\nconfig anterior]
    end

    Q --> R

    subgraph NAGIOS["Nagios Core — Processo de Reload"]
        R[Recebe SIGHUP]
        R --> S[Relê nagios.cfg\ne todos os cfg_file/cfg_dir]
        S --> T{Validação\ninterna}
        T -->|OK| U[Aplica nova\nconfiguração\nsem downtime]
        T -->|Falha| V[Mantém config\nanterior em memória\nregistra erro no log]
        U --> W[Agenda checks\npara novos hosts/serviços]
    end

    W --> X([Nagios monitora\nos novos hosts\nem :8080])
```

---

## Detalhamento: Geração dos Arquivos .cfg

```mermaid
flowchart LR
    subgraph DB["MariaDB — Tabelas"]
        T1[(tbl_host)]
        T2[(tbl_service)]
        T3[(tbl_contact)]
        T4[(tbl_command)]
        T5[(tbl_timeperiod)]
        T6[(tbl_configtarget)]
    end

    subgraph PHP["NagiosQL — Write Config Files"]
        P1[Lê tbl_configtarget\npara saber os caminhos]
        P2[Gera define host\n{ ... }]
        P3[Gera define service\n{ ... }]
        P4[Gera define contact\n{ ... }]
    end

    subgraph FILES["Volume: nagios-etc/nagiosql/"]
        F1[hosts/meuhost.cfg]
        F2[hosts/servidor-web.cfg]
        F3[services/ping.cfg]
        F4[services/http.cfg]
        F5[contacts.cfg]
        F6[commands.cfg]
        F7[hosttemplates.cfg]
        F8[servicetemplates.cfg]
    end

    T6 --> P1
    T1 --> P2 --> F1 & F2
    T2 --> P3 --> F3 & F4
    T3 --> P4 --> F5
    T4 --> F6
    T5 --> F7 & F8
```

---

## Detalhamento: Mecanismo de Reload

```mermaid
sequenceDiagram
    actor Admin
    participant NagiosQL as NagiosQL (PHP)
    participant Trigger as reload.trigger
    participant Watcher as reload-watcher.sh
    participant Nagios as Nagios Core

    Admin->>NagiosQL: Clica "Restart Nagios"
    NagiosQL->>NagiosQL: Lê commandfile da tbl_configtarget
    Note over NagiosQL: commandfile =\n/usr/local/nagios/var/reload.trigger

    NagiosQL->>Trigger: cria arquivo reload.trigger

    loop Polling a cada 2s
        Watcher->>Trigger: [ -f reload.trigger ]
    end

    Watcher->>Trigger: Detecta arquivo → remove (reset)
    Watcher->>Watcher: pgrep -x nagios → obtém PID

    Watcher->>Watcher: nagios -v nagios.cfg (validação prévia)

    alt nagios.cfg válido
        Watcher->>Nagios: kill -HUP $PID (SIGHUP)
        Nagios->>Nagios: Relê todos os .cfg files
        Nagios->>Nagios: Valida configuração internamente
        alt Configuração válida
            Nagios->>Nagios: Aplica nova config sem downtime
            Nagios-->>Admin: Hosts/serviços aparecem em :8080
        else Erro interno (raro — já validado antes)
            Nagios->>Nagios: Mantém config anterior
            Nagios-->>Admin: Erro registrado em nagios.log
        end
    else nagios.cfg inválido
        Watcher-->>Watcher: Loga erro, NÃO envia SIGHUP
        Note over Watcher: Nagios continua rodando\ncom a config anterior
        Watcher-->>Admin: Erro visível em: docker logs nagios-core
    end
```

---

## Estrutura de Arquivos no Volume Compartilhado

```
./volumes/nagios-etc/
│
├── nagios.cfg                      ← Arquivo principal (lido pelo Nagios Core)
│   │  cfg_file=.../timeperiods.cfg     │
│   │  cfg_dir=.../nagiosql/hosts/      │── NagiosQL inclui esses paths
│   │  cfg_dir=.../nagiosql/services/   │
│   └─ ...                              │
│
├── objects/                        ← Defaults do Nagios Core (não editados pelo NagiosQL)
│   ├── commands.cfg                    Comandos padrão (check_ping, check_http, etc.)
│   ├── contacts.cfg                    Contato nagiosadmin
│   ├── templates.cfg                   Templates linux-server, generic-service
│   ├── timeperiods.cfg                 24x7, workhours
│   └── localhost.cfg                   Host padrão localhost
│
├── nagiosql/                       ← Gerenciado exclusivamente pelo NagiosQL
│   ├── commands.cfg                    Comandos customizados (ex: check_dns)
│   ├── contacts.cfg                    Contatos criados via NagiosQL
│   ├── hosttemplates.cfg               Templates de host
│   ├── servicetemplates.cfg            Templates de serviço
│   ├── timeperiods.cfg                 Períodos customizados
│   ├── hostgroups.cfg                  Grupos de hosts
│   ├── servicegroups.cfg               Grupos de serviços
│   ├── hosts/                      ← Um .cfg por host
│   │   ├── gateway.cfg
│   │   ├── google-dns.cfg
│   │   └── linux-host.cfg
│   ├── services/                   ← .cfg agrupados por tipo de serviço
│   │   ├── ping.cfg
│   │   ├── http.cfg
│   │   └── ssh.cfg
│   └── backup/                     ← NagiosQL guarda versões anteriores aqui
│       ├── hosts/
│       └── services/
│
└── var/
    ├── reload.trigger              ← NagiosQL escreve aqui para disparar reload
    ├── nagios.log                  ← Log principal do Nagios
    ├── status.dat                  ← Estado atual de todos os hosts/serviços
    ├── rw/
    │   ├── nagios.cmd              ← Pipe de comandos externos (grupo nagioscfg)
    │   └── nagios.qh               ← Query handler socket
    └── spool/checkresults/         ← Resultados dos checks ativos
```

---

## Tabelas do Banco de Dados (MariaDB)

```mermaid
erDiagram
    tbl_configtarget {
        int id PK
        varchar target
        varchar basedir
        varchar hostconfig
        varchar serviceconfig
        varchar commandfile
        varchar binaryfile
        varchar conffile
        varchar pidfile
    }

    tbl_host {
        int id PK
        varchar host_name
        varchar alias
        varchar address
        int use_template
        int active
    }

    tbl_service {
        int id PK
        varchar service_description
        varchar check_command
        int use_template
        int active
    }

    tbl_contact {
        int id PK
        varchar contact_name
        varchar alias
        varchar email
    }

    tbl_command {
        int id PK
        varchar command_name
        varchar command_line
    }

    tbl_timeperiod {
        int id PK
        varchar timeperiod_name
        varchar alias
    }

    tbl_user {
        int id PK
        varchar username
        varchar password
        tinyint admin_enable
        tinyint active
    }

    tbl_settings {
        varchar category
        varchar name
        varchar value
    }

    tbl_host ||--o{ tbl_service : "host_name referência"
    tbl_host }o--|| tbl_configtarget : "escrito em basedir/hosts/"
    tbl_service }o--|| tbl_configtarget : "escrito em basedir/services/"
    tbl_command ||--o{ tbl_service : "check_command"
    tbl_timeperiod ||--o{ tbl_host : "check_period"
    tbl_contact }o--o{ tbl_host : "contact_groups"
```

---

## Resumo do Ciclo Completo

| Etapa | Responsável | O que acontece |
|---|---|---|
| **1. Cadastro** | NagiosQL (usuário) | Dados inseridos via formulário → salvos no MariaDB |
| **2. Geração** | NagiosQL (PHP) | Lê o banco → gera arquivos `.cfg` no volume |
| **3. Verificação** | NagiosQL → Nagios | Executa `nagios -v nagios.cfg` → exibe erros/warnings |
| **4. Sinalização** | NagiosQL | Cria `reload.trigger` para disparar o reload |
| **5. Validação prévia** | reload-watcher.sh | Detecta o trigger → executa `nagios -v` antes de recarregar |
| **6. SIGHUP** | reload-watcher.sh | Só enviado se a validação passar — protege config anterior |
| **7. Reload** | Nagios Core | Relê todos os `.cfg` → valida → aplica sem downtime |
| **8. Monitoramento** | Nagios Core | Agenda e executa checks dos novos hosts/serviços |
