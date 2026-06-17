## Why

NagiosQL 3.5 is a legacy PHP application (last upstream update 2022) with significant technical debt: unsalted MD5 passwords, SQL queries without prepared statements, MyISAM tables without transactions, and dependencies on obsolete PEAR/YUI libraries. Migrating to Go with Echo v5 and GORM produces a single self-contained binary (frontend embedded via `embed`), eliminates separate PHP/Nginx/PHP-FPM dependencies, modernizes security, and maintains functional and visual compatibility with the original NagiosQL.

## What Changes

- **New Go project** (`go-nagiosql`) replacing the entire PHP layer of NagiosQL
- Single binary (`nagiosql`) with Vue.js 3 + TypeScript + Tailwind CSS v4 + Lucide Icons + Pinia frontend embedded via `embed`
- HTTP backend: Echo v5 (stable v5.2.0) with REST routes equivalent to original NagiosQL URLs; OpenAPI/Swagger docs via swaggo
- ORM: GORM with MariaDB — schema compatible with existing NagiosQL 3.5 database
- Configuration via `config.toml` managed by Viper; CLI via Cobra (subcommands: `serve`, `migrate`, `version`, `config write|verify|restart`, `import`)
- JWT authentication (access token 15 min + refresh token 7 days in httpOnly cookie) with bcrypt passwords (replacing PHP unsalted MD5); testable via `curl` / `httpie` from the console
- Comprehensive test suite: unit tests (`go test`), handler tests (`httptest`), config generation tests, and curl-based API smoke tests
- Nagios `.cfg` file generation engine (replaces `NagConfigClass.php`)
- Existing `.cfg` file import engine (replaces `NagImportClass.php`)
- Nagios Core control: config validation (`nagios -v`) and reload via command pipe
- Custom object variables (`_VARNAME`) on hosts/services and resource variables (`$USER1$`–`$USERn$`) for `resource.cfg`
- Structured logging via `zerolog` (matching go-postfixadmin pattern)
- Project pattern following `github.com/jniltinho/go-postfixadmin`
- Docker container: replaces only the `nagios-core` PHP layer, keeping `nagios-db` (MariaDB)

## Capabilities

### New Capabilities

- `project-foundation`: Go 1.26 module `go-nagiosql`, Cobra CLI (`serve`, `version`, `migrate`, `config`, `import`), Viper with `config.toml`, GORM/MariaDB, zerolog structured logging, comprehensive Makefile (build/test/swagger/upx/lint), Dockerfile
- `authentication`: JWT access + refresh token auth (`golang-jwt/jwt/v5`), bcrypt passwords, user and group management, domain-scoped access control, legacy MD5 migration flow
- `nagios-objects`: Full CRUD for all Nagios objects — hosts, services, commands, time periods, contacts, contact groups, host groups, service groups, templates (host/service/contact), dependencies, escalations, extinfo, **custom object variables** (`_VARNAME`), and **resource variables** (`$USERn$` → `resource.cfg`); URLs and endpoint names compatible with NagiosQL
- `config-generation`: Nagios `.cfg` file generation engine in Go (replaces `NagConfigClass.php`); local/FTP/SSH target support; stale config detection; per-entity generation (individual host, service group)
- `nagios-control`: Config verification (`nagios -v`), restart command via pipe, reload via trigger file; equivalent to NagiosQL's `verify.php`; also exposed as Cobra `config` subcommands for scripting/cron use
- `import-engine`: Parser for existing Nagios `.cfg` files to populate the database; equivalent to `NagImportClass.php`; also available as `nagiosql import` CLI command
- `rest-api`: Echo v5 routes for all resources — RESTful `/api/v1/` prefix, JWT middleware, CORS, zerolog request logging, request validation, OpenAPI/Swagger docs via swaggo; fully testable via `curl`/`httpie`
- `testing`: Unit tests for all services and handlers (`httptest`); integration tests (build tag); curl-based API smoke test scripts; `make test-cover` with ≥70% coverage target
- `web-frontend`: *(DEFERRED — separate proposal)* Vue.js 3 + TypeScript + Tailwind CSS v4 + Pinia SPA embedded in binary

### Modified Capabilities

*(None — this is a new project, not a modification of the existing PHP NagiosQL)*

## Impact

- **Removed**: PHP 8.4-fpm, Nginx (for NagiosQL), PEAR, YUI, fcgiwrap dependencies in the container
- **Kept**: MariaDB 10.11 as the database; compatible `tbl_*` schema; Nagios Core binary and its directory structure at `/usr/local/nagios`
- **New binary**: `nagiosql` (Go) listens on a configurable port (default 8081); serves the API (frontend deferred)
- **Docker**: Simplified Dockerfile — no php-fpm, no dedicated Nginx for NagiosQL; just `nagiosql` + `nagios` + `reload-watcher` under supervisord
- **Database**: GORM auto-migrations for schema creation/update; backward compatible with existing NagiosQL 3.5 data
- **Passwords**: bcrypt replacing MD5; requires password reset on first post-migration login
