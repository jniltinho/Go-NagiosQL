# go-nagiosql

A Go rewrite of [NagiosQL](https://www.nagiosql.org/) — a web-based configuration manager for Nagios/Icinga.

**Stack:** Go 1.26 · Echo v5 · GORM · Cobra · Viper · MariaDB · JWT

---

## Features

- REST API with JWT authentication (Bearer access token + httpOnly refresh cookie)
- Full CRUD for all Nagios object types: hosts, services, commands, contacts, timeperiods, hostgroups, servicegroups, contactgroups, host/service/contact templates, variable definitions
- Nagios config file generator (`host.cfg`, `services.cfg`) with automatic backup
- Import from existing Nagios `.cfg` files
- Admin-only settings and user management
- Swagger UI at `/docs/swagger-ui.html`
- CLI: `serve`, `migrate`, `version`, `config`, `import`
- Docker-ready with a two-stage Alpine image

---

## Quick Start

### Prerequisites

- Go 1.26+
- MariaDB 10.6+ (or MySQL 8+)
- `swag` CLI (only for regenerating docs): `go install github.com/swaggo/swag/cmd/swag@latest`

### 1. Clone and configure

```bash
git clone https://github.com/jniltinho/go-nagiosql
cd go-nagiosql
cp config.toml.example config.toml
# Edit config.toml — at minimum change jwt.secret and database credentials.
```

### 2. Migrate the database

```bash
go run . migrate --admin-password yourpassword
```

Add `--sample` to insert example hosts/services/commands.

### 3. Run the server

```bash
go run . serve
# or build first:
make build && ./nagiosql serve
```

The API listens on `http://localhost:8081` by default.

---

## Configuration

All settings live in `config.toml` (or environment variables with the `NAGIOSQL_` prefix).

```toml
[server]
port = 8081
dev  = false   # enables Echo debug logger

[jwt]
secret           = "CHANGE-ME-to-a-random-string-of-at-least-32-chars"
access_ttl_min   = 15       # Bearer token lifetime (minutes)
refresh_ttl_days = 7        # httpOnly cookie lifetime (days)

[database]
host     = "127.0.0.1"
port     = 3306
name     = "nagiosql"
user     = "nagiosql"
password = "nagiosql"

[nagios]
base_dir           = "/usr/local/nagios"
config_file        = "/usr/local/nagios/etc/nagios.cfg"
host_config_dir    = "/usr/local/nagios/etc/hosts/"
service_config_dir = "/usr/local/nagios/etc/services/"
backup_dir         = "/usr/local/nagios/etc/backup/"
import_dir         = "/usr/local/nagios/etc/import/"
# reload_trigger is a plain file — touching it signals a reload watcher.
# NEVER point this at nagios.cmd (deadlock risk).
reload_trigger     = "/usr/local/nagios/var/reload.trigger"
binary             = "/usr/local/nagios/bin/nagios"
pid_file           = "/usr/local/nagios/var/nagios.lock"
```

### Environment variable override

Every config key maps to `NAGIOSQL_<SECTION>_<KEY>`, e.g.:

```bash
NAGIOSQL_DATABASE_HOST=10.0.0.1
NAGIOSQL_JWT_SECRET=mysecret
NAGIOSQL_SERVER_PORT=9090
```

---

## CLI Commands

```
nagiosql serve          Start the HTTP API server
nagiosql migrate        AutoMigrate schema; seed admin user
  --admin-password      Password for the admin user (default: admin123)
  --sample              Insert sample data
  --config              Path to config.toml
nagiosql import         Import Nagios .cfg files into the database
  --dir                 Directory containing .cfg files
  --overwrite           Overwrite existing objects (default: false)
nagiosql config write   Write current config to stdout
nagiosql config verify  Verify config file is readable
nagiosql config restart Touch reload_trigger to signal Nagios reload
nagiosql version        Print version and build date
```

---

## API Reference

### Authentication

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/auth/login` | Obtain access + refresh tokens |
| POST | `/api/v1/auth/logout` | Revoke refresh cookie |
| POST | `/api/v1/auth/refresh` | Exchange refresh cookie for new access token |

**Login request:**
```json
{ "username": "admin", "password": "admin123" }
```

**Login response:**
```json
{
  "access_token": "eyJ...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

If the account has a legacy MD5 password, login returns HTTP 200 with:
```json
{ "requires_password_reset": true }
```
No token is issued — the user must reset their password before continuing.

All protected endpoints require:
```
Authorization: Bearer <access_token>
```

### Hosts

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/hosts` | List hosts (paginated) |
| POST | `/api/v1/hosts` | Create host |
| GET | `/api/v1/hosts/:id` | Get host |
| PUT | `/api/v1/hosts/:id` | Update host |
| DELETE | `/api/v1/hosts/:id` | Delete host |

### Services

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/services` | List services |
| POST | `/api/v1/services` | Create service |
| GET | `/api/v1/services/:id` | Get service |
| PUT | `/api/v1/services/:id` | Update service |
| DELETE | `/api/v1/services/:id` | Delete service |

### Commands

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/commands` | List commands |
| POST | `/api/v1/commands` | Create command |
| GET | `/api/v1/commands/:id` | Get command |
| PUT | `/api/v1/commands/:id` | Update command |
| DELETE | `/api/v1/commands/:id` | Delete command |

### Contacts

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/contacts` | List contacts |
| POST | `/api/v1/contacts` | Create contact |
| GET | `/api/v1/contacts/:id` | Get contact |
| PUT | `/api/v1/contacts/:id` | Update contact |
| DELETE | `/api/v1/contacts/:id` | Delete contact |

### Timeperiods

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/timeperiods` | List timeperiods |
| POST | `/api/v1/timeperiods` | Create timeperiod (with inline ranges) |
| GET | `/api/v1/timeperiods/:id` | Get timeperiod |
| PUT | `/api/v1/timeperiods/:id` | Update timeperiod |
| DELETE | `/api/v1/timeperiods/:id` | Delete timeperiod |

### Groups

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/hostgroups` | List hostgroups |
| POST | `/api/v1/hostgroups` | Create hostgroup |
| GET | `/api/v1/hostgroups/:id` | Get hostgroup |
| PUT | `/api/v1/hostgroups/:id` | Update hostgroup |
| DELETE | `/api/v1/hostgroups/:id` | Delete hostgroup |
| POST | `/api/v1/hostgroups/:id/members` | Add host to hostgroup |
| GET | `/api/v1/servicegroups` | List servicegroups |
| POST | `/api/v1/servicegroups` | Create servicegroup |
| DELETE | `/api/v1/servicegroups/:id` | Delete servicegroup |
| GET | `/api/v1/contactgroups` | List contactgroups |
| POST | `/api/v1/contactgroups` | Create contactgroup |
| DELETE | `/api/v1/contactgroups/:id` | Delete contactgroup |

### Templates

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/hosttemplates` | List host templates |
| POST | `/api/v1/hosttemplates` | Create host template |
| GET | `/api/v1/hosttemplates/:id` | Get host template |
| DELETE | `/api/v1/hosttemplates/:id` | Delete host template |
| GET | `/api/v1/servicetemplates` | List service templates |
| POST | `/api/v1/servicetemplates` | Create service template |
| DELETE | `/api/v1/servicetemplates/:id` | Delete service template |
| GET | `/api/v1/contacttemplates` | List contact templates |
| POST | `/api/v1/contacttemplates` | Create contact template |
| DELETE | `/api/v1/contacttemplates/:id` | Delete contact template |

### Variable Definitions

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/variables` | List variable definitions |
| POST | `/api/v1/variables` | Create variable definition |
| GET | `/api/v1/variables/:id` | Get variable definition |
| PUT | `/api/v1/variables/:id` | Update variable definition |
| DELETE | `/api/v1/variables/:id` | Delete variable definition |

### Users (admin only)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/users` | List users |
| POST | `/api/v1/users` | Create user (bcrypt, cost 12) |
| GET | `/api/v1/users/:id` | Get user |
| PUT | `/api/v1/users/:id/password` | Change user password |
| DELETE | `/api/v1/users/:id` | Delete user |

### Settings (admin write, any auth read)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/settings` | Get global settings |
| PUT | `/api/v1/settings` | Update global settings (admin only) |

### Config generation

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/config/write` | Generate all `.cfg` files from DB |
| POST | `/api/v1/config/verify` | Run `nagios -v nagios.cfg` |
| POST | `/api/v1/config/restart` | Touch `reload_trigger` |

### Import

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/import` | Import Nagios `.cfg` files |

```json
{ "dir": "/etc/nagios/objects", "overwrite": false }
```

### Monitoring

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/monitoring/summary` | Object counts (hosts, services, …) |

### Logbook

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/logbook` | List audit log entries |

---

## Swagger UI

Start the server and open: `http://localhost:8081/docs/swagger-ui.html`

To regenerate after changing annotations:
```bash
make swagger
```

---

## Development

```bash
make test              # unit tests (SQLite in-memory, no external deps)
make db-start          # start MariaDB 10.11 on port 3307 via Docker
make test-integration  # integration tests against the test DB
make db-stop           # stop test DB
make check             # vet + build + unit tests (CI entry point)
make swagger           # regenerate OpenAPI docs
make run               # build and serve (reads config.toml)
```

### Running unit tests

No external database needed — tests use an SQLite in-memory DB:

```bash
go test -race -cover ./...
```

### Running integration tests

```bash
make db-start
go test -tags integration -v ./internal/integration/...
make db-stop
```

---

## Docker

### Production image

```bash
docker build -t go-nagiosql .
docker run -p 8081:8081 \
  -v $(pwd)/config.toml:/app/config.toml:ro \
  go-nagiosql
```

### With the Nagios stack (docker-compose)

The `docker-compose.override.yml` file adds `go-nagiosql` to an existing Nagios Core compose stack:

```bash
cp docker-compose.override.yml DOCUMENTS/docker/nagios-core/docker-compose.override.yml
cd DOCUMENTS/docker/nagios-core
docker compose up -d
```

Use `config.toml.docker` as the mounted config — it points to the shared `nagios-etc` and `nagios-var` volumes.

---

## Migrating from PHP NagiosQL

1. **Keep the existing database** — `go-nagiosql` is schema-compatible with the PHP version's `tbl_*` tables.
2. Run `nagiosql migrate` — this runs `AutoMigrate` which adds any missing columns without dropping data.
3. **Legacy passwords**: Accounts with MD5 passwords (`$2` prefix absent) will require a password reset on first login. The API returns `{requires_password_reset: true}` and no token is issued.
4. After verifying the Go API works, stop the PHP application and point your reverse proxy to port 8081.

---

## Security Notes

- Passwords are stored with **bcrypt cost 12** for all new/reset passwords.
- MD5 legacy passwords are never used to issue tokens — they force a reset.
- JWT secrets must be at least 32 characters. Rotate by changing `jwt.secret` and restarting.
- `reload_trigger` must be a regular file, not `nagios.cmd`. Writing to the FIFO from Go causes a deadlock.
- All writes to Nagios config files include automatic timestamped backups in `nagios.backup_dir`.

---

## License

MIT
