## Context

NagiosQL 3.5 is a PHP web application that manages Nagios Core configuration through a MariaDB database and generates `.cfg` files on disk. The current Docker stack runs five supervisord processes in a single container: php-fpm, fcgiwrap, nagios daemon, nginx, and reload-watcher. The Go migration targets the PHP/Nginx/PHP-FPM layer only — Nagios Core, MariaDB, and the reload-watcher remain unchanged.

The reference implementation to follow is `github.com/jniltinho/go-postfixadmin`: a Go 1.26 binary with Echo v5, GORM, Cobra, Viper, swaggo, and an embedded Vue.js 3 + TypeScript + Tailwind CSS v4 + Pinia frontend that manages Postfix mail server configuration via a web UI backed by MariaDB. This project deviates from go-postfixadmin by using JWT (not gorilla/sessions) and Echo v5's built-in logger (not zerolog).

Current stack constraints:
- Nagios Core is hardcoded to `/usr/local/nagios` prefix throughout NagiosQL's DB (`tbl_configtarget`)
- The MariaDB schema uses MyISAM tables with `tbl_*` naming; GORM models must map to this schema
- The reload mechanism uses a file trigger (`/usr/local/nagios/var/reload.trigger`) polled by `reload-watcher.sh`
- NagiosQL's `config/settings.php` bootstrap file will be replaced by `config.toml`

## Goals / Non-Goals

**Goals:**
- Complete REST API backend (`nagiosql`) with JWT authentication, fully testable via `curl`/`httpie` from the console
- Feature parity with NagiosQL 3.5: all object types, all CRUD operations, config generation, import, verify/reload
- API endpoint naming compatible with NagiosQL page structure (`/api/v1/hosts`, `/api/v1/services`, `/api/v1/verify`, etc.)
- JWT auth: short-lived access token (15 min Bearer) + long-lived refresh token (7 days, httpOnly cookie)
- bcrypt password hashing replacing unsalted MD5; legacy MD5 migration flow on first login
- GORM auto-migration for schema management (create tables if missing, add columns if new)
- `config.toml` as the single configuration file (DB credentials, server port, Nagios paths, JWT secret)
- Cobra CLI with `serve`, `migrate`, `version`, `config write|verify|restart`, and `import` commands
- Comprehensive test suite: unit + handler (httptest) + integration (build tag) + curl smoke tests
- Comprehensive Makefile: `build`, `test`, `test-cover`, `test-integration`, `swagger`, `lint`, `upx`, `clean`
- Docker-ready: drop-in replacement for the NagiosQL PHP layer in the existing compose stack
- OpenAPI/Swagger docs via swaggo served at `/api/swagger/`

**Non-Goals:**
- Vue.js frontend (deferred to a separate proposal)
- Replacing Nagios Core itself
- Replacing MariaDB with a different database engine
- Implementing the Nagios Core CGI web UI (port 8080 stays unchanged)
- Full Nagios 4 feature additions beyond what NagiosQL 3.5 already supported
- Multi-node or distributed deployment
- Real-time monitoring dashboard (that is Nagios Core's responsibility)

## Decisions

### D1: Echo v5 over Gin or Chi

Echo v5 (stable at v5.2.0, released 2026-06-14) is chosen because it aligns with the `go-postfixadmin` reference pattern, has strong middleware ecosystem, built-in request validation, and clean route grouping for `/api/v1/` namespacing. Gin has more community usage but Echo's API is cleaner for this CRUD-heavy workload. Note: in Echo v5, `Context` is now a **struct** (not an interface) — handler signatures must use `echo.Context` directly, and custom context embedding is no longer supported via interface. Generic parameter binding (`PathParam[T]`, `QueryParam[T]`) is available for type-safe extraction.

### D2: GORM with existing schema instead of full redesign

The existing `tbl_*` schema is kept as-is to allow zero-downtime migration (existing NagiosQL data works immediately). GORM struct tags map to the exact column names. Auto-migration only adds missing tables/columns — it never drops or renames. The MyISAM engine is preserved for compatibility; InnoDB migration is a separate future concern.

```
tbl_host      → models.Host
tbl_service   → models.Service
tbl_command   → models.Command
tbl_timeperiod → models.Timeperiod
...
```

### D3: Vue.js 3 + TypeScript + Tailwind CSS v4 + Pinia + Lucide embedded via `embed`

The frontend build output (`dist/`) is embedded into the Go binary at compile time using `//go:embed frontend/dist`. The Echo server serves static files from the embedded FS. This matches the `go-postfixadmin` pattern and produces a single deployable artifact.

Vite is the build tool (fast, native ESM). TypeScript is used for all frontend code (matching go-postfixadmin). Pinia is the state management library (replacing Vuex — lighter, better TS support). Vue Router handles client-side navigation with routes that mirror NagiosQL's page IDs and URLs (e.g., `/admin/hosts`, `/admin/services`, `/admin/verify`). Tailwind CSS v4 (used by go-postfixadmin) instead of v3.

### D4: Cobra CLI structure

```
nagiosql
├── serve               # Start HTTP server (reads config.toml) [default]
├── migrate             # Run GORM auto-migrations
├── version             # Print version, build date, Go version
├── config
│   ├── write [type]    # Write config files (all, or specific type: host, service, etc.)
│   ├── verify          # Run nagios -v and print output
│   └── restart         # Validate then restart Nagios Core
└── import [file]       # Import .cfg file(s) into database
```

`serve` is the default command when no subcommand is given, matching the `go-postfixadmin` pattern for Docker entrypoint simplicity. The `config` and `import` subcommands replace NagiosQL's `do_config.php` scripting interface, enabling cron-based automation without the web UI.

### D5: bcrypt for all new passwords; detect-and-block legacy MD5

**New passwords** (admin seeded by `nagiosql migrate`, user creation via API, password reset): always stored as `bcrypt` cost=12 via `golang.org/x/crypto/bcrypt`. Hash starts with `$2a$12$`.

**Legacy NagiosQL MD5 hashes** (32-char hex, no `$2` prefix): detected at login time. The system verifies the raw MD5 match (`md5(input) == storedHash`) and — if correct — returns HTTP 200 with `{"requires_password_reset": true}` and **no tokens**. No MD5 password is ever issued a JWT. After the user calls `POST /api/v1/auth/reset-password` with their old password + new password, the bcrypt hash replaces the MD5 in the DB and a normal token is issued.

**Standard library only**: `golang.org/x/crypto/bcrypt` (bcrypt) + `crypto/md5` from stdlib (MD5 detection). No third-party password library needed.

### D6: JWT authentication (access token + refresh token)

The API uses `github.com/golang-jwt/jwt/v5` for stateless authentication:

- **Access token**: HS256-signed JWT, 15-minute TTL, returned in JSON response body (`{"access_token": "..."}`)
- **Refresh token**: HS256-signed JWT, 7-day TTL, set as an httpOnly Secure cookie (`refresh_token`)
- **API calls**: `Authorization: Bearer <access_token>` header — directly testable via `curl -H "Authorization: Bearer $TOKEN"`
- **Rotation**: `POST /api/v1/auth/refresh` validates the refresh cookie and issues a new access token
- **Logout**: `POST /api/v1/auth/logout` clears the refresh cookie; access tokens are short-lived so no server-side revocation is needed

Claims payload: `{ "sub": userID, "username": "...", "admin": bool, "domain_id": int, "exp": ... }`

Sessions were considered but JWT is stateless (no store needed), REST-idiomatic, and far easier to test from the console. `gorilla/sessions` and `echo-contrib` are NOT used.

### D7: Project structure (following go-postfixadmin)

```
go-nagiosql/
├── cmd/
│   ├── root.go
│   ├── serve.go
│   ├── migrate.go
│   ├── version.go
│   ├── config.go          # nagiosql config write|verify|restart
│   └── import.go          # nagiosql import [file]
├── internal/
│   ├── api/
│   │   ├── handlers/      # Echo route handlers (one file per object group)
│   │   │   ├── auth.go
│   │   │   ├── hosts.go
│   │   │   ├── services.go
│   │   │   └── ...
│   │   ├── middleware/    # JWT auth, Echo built-in logger, CORS, request-id
│   │   │   ├── auth.go
│   │   │   └── logger.go
│   │   └── routes.go      # Route registration + Swagger UI mount
│   ├── models/            # GORM models (one file per tbl_* table)
│   ├── services/
│   │   ├── nagconfig/     # .cfg file generation (NagConfigClass equivalent)
│   │   ├── nagimport/     # .cfg file import (NagImportClass equivalent)
│   │   ├── auth/          # JWT issue/validate, bcrypt, MD5 detection
│   │   └── logbook/       # Audit log writer
│   ├── db/                # GORM init, migrations
│   └── config/            # Viper config loading
├── docs/                  # swaggo-generated OpenAPI spec (git-committed)
├── test/
│   ├── api/               # curl/httpie smoke test scripts
│   │   ├── auth.sh        # login, refresh, logout, me
│   │   ├── hosts.sh       # CRUD hosts via API
│   │   └── smoke.sh       # full flow: login → create → write config → verify
│   └── fixtures/          # SQL fixtures and sample .cfg files for tests
├── embed.go               # //go:embed docs
├── config.toml.example
├── main.go
├── go.mod                 # requires go 1.26
├── Makefile
└── Dockerfile             # multi-stage: go builder → debian:trixie-slim runtime
```

### Makefile targets

```makefile
BINARY    := nagiosql
VERSION   := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE:= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS   := -ldflags "-X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE) -s -w"

all:        swagger build        ## Build everything
build:                           ## Compile binary
test:                            ## Run unit + handler tests
test-v:                          ## Run tests with -v
test-cover:                      ## Generate HTML coverage report
test-integration:                ## Run integration tests (needs MariaDB)
test-api:                        ## Run curl smoke tests against running server
swagger:                         ## Regenerate OpenAPI spec via swag init
lint:                            ## Run golangci-lint
vet:                             ## Run go vet
upx:        build                ## Compress binary with UPX --best
clean:                           ## Remove binary, coverage, docs
run:        build                ## Build and run serve
migrate:    build                ## Build and run migrate
```

### D8: Config generation in Go

`NagConfigClass.php` (2519 lines) is reimplemented as `internal/services/nagconfig/`. The PEAR template system is replaced with Go's `text/template`. Template files are embedded alongside the binary. Generation supports three deployment methods: local filesystem, FTP (`github.com/jlaffaye/ftp`), SSH (`golang.org/x/crypto/ssh`).

### D9: config.toml (with NAGIOSQL_ env override)

```toml
[server]
port = 8081
dev  = false              # pretty-print logs when true

[jwt]
secret          = "change-me-min-32-chars"
access_ttl_min  = 15     # access token TTL in minutes
refresh_ttl_days = 7     # refresh token TTL in days

[database]
host     = "db"
port     = 3306
name     = "nagiosql"
user     = "nagiosql"
password = ""

[nagios]
base_dir        = "/usr/local/nagios"
config_file     = "/usr/local/nagios/etc/nagios.cfg"
cgi_file        = "/usr/local/nagios/etc/cgi.cfg"
resource_file   = "/usr/local/nagios/etc/resource.cfg"
reload_trigger  = "/usr/local/nagios/var/reload.trigger"   # file polled by reload-watcher.sh (NOT the nagios.cmd FIFO)
binary          = "/usr/local/nagios/bin/nagios"
pid_file        = "/usr/local/nagios/var/nagios.lock"
host_config_dir     = "/usr/local/nagios/etc/nagiosql/hosts/"
service_config_dir  = "/usr/local/nagios/etc/nagiosql/services/"
backup_dir          = "/usr/local/nagios/etc/nagiosql/backup/"
import_dir          = "/usr/local/nagios/etc/import/"
```

> **Critical**: The reload mechanism uses `reload.trigger` (a plain file whose mtime is polled by `reload-watcher.sh`) — **not** `nagios.cmd` (the Nagios external command FIFO). Writing `RESTART_PROGRAM` to `nagios.cmd` is dangerous from the Go process because the FIFO must be opened in a specific order to avoid deadlock. The trigger file approach, already used by the existing entrypoint, is safe.

Viper supports both `config.toml` and environment variable overrides (e.g., `NAGIOSQL_DATABASE_PASSWORD`), making Docker integration clean.

### D10: Echo v5 built-in logger (no zerolog)

Echo v5's built-in `middleware.Logger()` and `middleware.RequestID()` are used for request logging and request ID injection. This avoids an extra dependency and a custom adapter. In development mode (`cfg.Server.Dev=true`), `e.Logger.SetLevel(log.DEBUG)` enables verbose output. Seed functions and CLI commands use the standard `log` package (`log.Printf`). `rs/zerolog` is NOT used.

### D11: OpenAPI/Swagger documentation via swaggo

`github.com/swaggo/swag` generates an OpenAPI 2.0 spec from Go doc comments on handler functions. The spec is committed to `docs/` and served by `swaggo/http-swagger` at `/api/swagger/`. This provides interactive API documentation without an external service. Swagger annotations are added to all handler functions at implementation time.

### D12: Custom object variables and resource.cfg

NagiosQL manages two types of user-defined variables:
1. **Custom object variables** (`_VARNAME` prefix in Nagios): stored in `tbl_variabledefinition` with M:N link tables (`tbl_lnkHostToVariabledefinition`, `tbl_lnkServiceToVariabledefinition`, `tbl_lnkHosttemplateToVariabledefinition`, `tbl_lnkServicetemplateToVariabledefinition`). These appear in generated `.cfg` files as `_VARNAME value` inside `define host { }` blocks.
2. **Resource macros** (`$USER1$`–`$USERn$`): stored separately, written to Nagios `resource.cfg`. NagiosQL manages these via the `specials.php` admin page (`admin/specials.php`), which also appears to handle additional edge-case configurations.

Both must be supported in GORM models, CRUD handlers, and the config generation engine.

### D13: Existing Docker stack as development and test environment

The reference Docker stack at `DOCUMENTS/docker/nagios-core/` already provides everything needed for development and integration testing. **NEVER modify files under `DOCUMENTS/`** — they are read-only reference material.

```
DOCUMENTS/docker/nagios-core/
├── docker-compose.yml          ← MariaDB 10.11 + Nagios Core; ports 8080 (Nagios CGI) + 8081 (NagiosQL)
├── .env.example                ← Copy to .env for credentials
├── nagiosql/install/sql/
│   ├── nagiosQL_v35_db_mysql.sql     ← Full schema (60+ tables, MyISAM, utf8_unicode_ci)
│   └── import_nagios_sample.sql      ← Sample data: 24 commands, 5 timeperiods, 5 templates,
│                                        4 hosts (winserver/linksys-srw224p/hplj2605dn/localhost),
│                                        21 services, contact nagiosadmin
└── nagios/etc-extra/nagiosql/
    ├── commands.cfg            ← define command { check_dns ... }
    ├── hosts/
    │   ├── linux-host.cfg      ← define host { use linux-server ... address 192.168.1.9 }
    │   ├── gateway.cfg
    │   ├── google-dns.cfg
    │   └── cloudflare-dns.cfg
    └── services/
        ├── ping.cfg            ← define service { check_command check_ping!100.0,20%!500.0,60% }
        ├── dns.cfg
        ├── http.cfg
        └── ssh.cfg
```

**Key entrypoint facts** (from `nagios/entrypoint.sh`):
- Schema is loaded from `nagiosQL_v35_db_mysql.sql` only on first run (`TABLE_COUNT=0`)
- Admin user is seeded via `MD5('$NAGIOSQL_PASSWORD')` — confirming the bcrypt migration path
- `tbl_configtarget` is updated with Docker-specific paths:
  - `hostconfig` → `/usr/local/nagios/etc/nagiosql/hosts/`
  - `serviceconfig` → `/usr/local/nagios/etc/nagiosql/services/`
  - `basedir` → `/usr/local/nagios/etc/nagiosql/`
  - `backupdir` / `hostbackup` / `servicebackup` → `/usr/local/nagios/etc/nagiosql/backup/[hosts|services]/`
  - `commandfile` → `/usr/local/nagios/var/reload.trigger` (file trigger, not FIFO)
  - `binaryfile` → `/usr/local/nagios/bin/nagios`
  - `conffile` → `/usr/local/nagios/etc/nagios.cfg`
  - `cgifile` → `/usr/local/nagios/etc/cgi.cfg`
  - `resourcefile` → `/usr/local/nagios/etc/resource.cfg`
  - `importdir` → `/usr/local/nagios/etc/import/`
  - `version` → `4` (Nagios 4 directive set)

**Schema totals**: 60+ tables. Main objects (8), templates (3), advanced (6), link tables (~45), meta (tbl_settings, tbl_tablestatus, tbl_logbook, tbl_info, tbl_menu, tbl_language, tbl_relationinformation), auth (tbl_user, tbl_group, tbl_lnkGroupToUser), variables (tbl_variabledefinition + 6 link tables for host/service/hosttemplate/servicetemplate/contact/contacttemplate), config (tbl_configtarget, tbl_datadomain).

**Test strategy using existing stack**:
- `make db-start` starts only the `db` service from the existing compose on host port 3307 via an override file
- SQL fixtures for `test/fixtures/` are **symlinked or copied from** `DOCUMENTS/docker/nagios-core/nagiosql/install/sql/`
- Import test fixtures are the `.cfg` files from `DOCUMENTS/docker/nagios-core/nagios/etc-extra/nagiosql/`
- The Makefile `test-integration` target uses `TEST_DSN=nagiosql:nagiosqlpass@tcp(127.0.0.1:3307)/nagiosql_test`

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| GORM auto-migration may conflict with MyISAM-specific syntax | Use `db.Set("gorm:table_options", "ENGINE=MyISAM")` on model registration; test migrations against a live NagiosQL 3.5 DB before deploying |
| bcrypt migration requires password resets | Detect legacy MD5 hashes on login attempt; redirect to password-reset flow; document clearly in migration guide |
| Config generation correctness (NagConfigClass has 2519 lines of edge cases) | Port incrementally: start with hosts and services (80% of usage), add remaining object types; run generated files through `nagios -v` as acceptance test |
| Frontend bundle size with full Vue.js + Tailwind + Lucide | Use Vite tree-shaking and Tailwind purge; Lucide's tree-shakeable ESM exports; target < 500KB gzipped |
| JWT secret rotation invalidates all active tokens | Document that changing `[jwt].secret` in config.toml forces all users to re-login; acceptable for single-binary deployment |
| Short access TTL (15 min) may frustrate CLI users piping curl in scripts | Provide `nagiosql config write` CLI command that handles auth internally without needing a token; document token refresh for scripting |
| go-postfixadmin pattern may differ from current project conventions | Inspect actual repo structure at implementation time; adapt where this project's golang skills diverge |
| Echo v5 Context is now a struct — custom context patterns from v4 tutorials won't compile | Use `echo.Context` directly throughout; no embedding; use `c.Get()`/`c.Set()` for request-scoped values |
| Custom logger middleware not needed | Echo v5 `middleware.Logger()` handles structured request logging out of the box |
| hostextinfo/serviceextinfo are **deprecated in Nagios 4** | Include in UI with a visible deprecation notice; still generate `.cfg` for backward compatibility with Nagios 3 targets |
| `admin/specials.php` purpose is unclear from source inspection | Research before implementing; likely handles `resource.cfg` `$USERn$` macros — model as `tbl_variabledefinition` with `nagiosresource` type |

## Migration Plan

1. **Phase 1 — Foundation**: Project scaffold, GORM models for all `tbl_*` tables, `config.toml`, Cobra CLI, `migrate` command
2. **Phase 2 — Auth + Basic UI**: Login/logout, session middleware, user/group CRUD, monitoring dashboard (overview page)
3. **Phase 3 — Core Objects**: Hosts, services, commands, time periods, contacts, groups — full CRUD with config generation for each
4. **Phase 4 — Advanced Objects**: Templates, dependencies, escalations, extinfo — CRUD + config generation
5. **Phase 5 — Verify & Control**: `verify.php` equivalent — write configs, validate, restart/reload Nagios
6. **Phase 6 — Import**: `.cfg` file import engine
7. **Phase 7 — Settings & Admin**: Config targets, data domains, NagiosQL settings, logbook
8. **Phase 8 — Docker Integration**: Replace PHP layer in Dockerfile; update supervisord.conf; end-to-end test

**Rollback**: Keep the original PHP NagiosQL container image tagged. Switch `docker-compose.yml` back to the PHP image. Database is shared and unchanged.

## Open Questions

- Should GORM auto-migration run at startup (inside `serve`) or only via explicit `migrate` command? Explicit is safer for production; `serve` should only auto-migrate in development mode.
- The original NagiosQL supports FTP and SSH config targets — are these needed in v1 or can local-only ship first?
- What exactly does `admin/specials.php` do beyond `$USERn$` resource macros? Inspect the PHP source before implementing the `specials` handler.
- Should i18n (NagiosQL supports 10 languages via gettext) be included in v1? go-postfixadmin uses `leonelquinteros/gotext` — include the dependency and wire PT-BR + EN at minimum; other locales are additive.
- Should UPX binary compression be applied in the Makefile's `build` target by default or only in a separate `compress` target? UPX can reduce binary size by ~60% but adds ~1s to build time.
