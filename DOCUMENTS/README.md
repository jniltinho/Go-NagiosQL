# go-nagiosql — Reference Documentation

This document is the detailed reference for the Go rewrite of NagiosQL.
For a project overview and quick start, see the [root README](../README.md).

---

## Table of Contents

- [Configuration](#configuration)
- [API Reference](#api-reference)
  - [Authentication](#authentication)
  - [Hosts](#hosts)
  - [Services](#services)
  - [Commands](#commands)
  - [Contacts](#contacts)
  - [Timeperiods](#timeperiods)
  - [Groups](#groups)
  - [Templates](#templates)
  - [Variable Definitions](#variable-definitions)
  - [Users](#users)
  - [Settings](#settings)
  - [Config Generation](#config-generation)
  - [Import](#import)
  - [Monitoring](#monitoring)
  - [Logbook](#logbook)
- [Migrating from PHP NagiosQL](#migrating-from-php-nagiosql)
- [Docker](#docker)
- [Security Notes](#security-notes)
- [PHP NagiosQL Docs](#php-nagiosql-docs)

---

## Configuration

All settings live in `config.toml`. Copy `config.toml.example` as a starting point.
Every key can be overridden with an environment variable using the `NAGIOSQL_` prefix.

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
# NEVER point this at nagios.cmd (FIFO deadlock risk).
reload_trigger     = "/usr/local/nagios/var/reload.trigger"
binary             = "/usr/local/nagios/bin/nagios"
pid_file           = "/usr/local/nagios/var/nagios.lock"
```

### Environment variable override

Every key maps to `NAGIOSQL_<SECTION>_<KEY>`:

```bash
NAGIOSQL_DATABASE_HOST=10.0.0.1
NAGIOSQL_JWT_SECRET=mysecret
NAGIOSQL_SERVER_PORT=9090
NAGIOSQL_NAGIOS_HOST_CONFIG_DIR=/etc/nagios/hosts/
```

---

## API Reference

All protected endpoints require:

```
Authorization: Bearer <access_token>
```

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

**Login response (success):**
```json
{
  "access_token": "eyJ...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

**Login response (legacy MD5 password):**
```json
{ "requires_password_reset": true }
```
No token is issued. The user must reset their password via `PUT /api/v1/users/:id/password` before continuing.

---

### Hosts

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/hosts` | List hosts (paginated) |
| POST | `/api/v1/hosts` | Create host |
| GET | `/api/v1/hosts/:id` | Get host |
| PUT | `/api/v1/hosts/:id` | Update host |
| DELETE | `/api/v1/hosts/:id` | Delete host |

**Create host (minimal):**
```json
{
  "host_name": "web01",
  "alias": "Web Server 01",
  "address": "10.0.0.1",
  "active_checks_enabled": 1,
  "max_check_attempts": 3
}
```

---

### Services

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/services` | List services |
| POST | `/api/v1/services` | Create service |
| GET | `/api/v1/services/:id` | Get service |
| PUT | `/api/v1/services/:id` | Update service |
| DELETE | `/api/v1/services/:id` | Delete service |

---

### Commands

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/commands` | List commands |
| POST | `/api/v1/commands` | Create command |
| GET | `/api/v1/commands/:id` | Get command |
| PUT | `/api/v1/commands/:id` | Update command |
| DELETE | `/api/v1/commands/:id` | Delete command |

---

### Contacts

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/contacts` | List contacts |
| POST | `/api/v1/contacts` | Create contact |
| GET | `/api/v1/contacts/:id` | Get contact |
| PUT | `/api/v1/contacts/:id` | Update contact |
| DELETE | `/api/v1/contacts/:id` | Delete contact |

---

### Timeperiods

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/timeperiods` | List timeperiods |
| POST | `/api/v1/timeperiods` | Create timeperiod (with inline ranges) |
| GET | `/api/v1/timeperiods/:id` | Get timeperiod |
| PUT | `/api/v1/timeperiods/:id` | Update timeperiod |
| DELETE | `/api/v1/timeperiods/:id` | Delete timeperiod |

---

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
| GET | `/api/v1/servicegroups/:id` | Get servicegroup |
| DELETE | `/api/v1/servicegroups/:id` | Delete servicegroup |
| GET | `/api/v1/contactgroups` | List contactgroups |
| POST | `/api/v1/contactgroups` | Create contactgroup |
| GET | `/api/v1/contactgroups/:id` | Get contactgroup |
| DELETE | `/api/v1/contactgroups/:id` | Delete contactgroup |

---

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

---

### Variable Definitions

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/variables` | List variable definitions |
| POST | `/api/v1/variables` | Create variable definition |
| GET | `/api/v1/variables/:id` | Get variable definition |
| PUT | `/api/v1/variables/:id` | Update variable definition |
| DELETE | `/api/v1/variables/:id` | Delete variable definition |

---

### Users

> Admin role required for all user management endpoints.

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/users` | List users |
| POST | `/api/v1/users` | Create user (bcrypt cost 12) |
| GET | `/api/v1/users/:id` | Get user |
| PUT | `/api/v1/users/:id/password` | Change user password |
| DELETE | `/api/v1/users/:id` | Delete user |

**Change password:**
```json
{ "new_password": "hunter2" }
```

A user cannot delete their own account.

---

### Settings

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/settings` | Get global settings (any authenticated user) |
| PUT | `/api/v1/settings` | Update global settings (admin only) |

---

### Config Generation

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/config/write` | Generate all `.cfg` files from DB |
| POST | `/api/v1/config/verify` | Run `nagios -v nagios.cfg` |
| POST | `/api/v1/config/restart` | Touch `reload_trigger` to signal reload |

The config generator produces output format-compatible with the original PHP NagiosQL:

```
define host {
    use                 linux-server
    host_name           web01
    alias               Web Server 01
    address             10.0.0.1
    max_check_attempts  3
    contact_groups      admins
}
```

- Template names are resolved from the `tbl_lnkHostToHosttemplate` join table.
- Contact groups and host links are resolved from their respective link tables.
- Fields with value `2` (inherit from template) are silently omitted.
- `register` is only written for template objects (`register 0`).
- Existing files are automatically backed up to `nagios.backup_dir` before overwriting.

---

### Import

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/import` | Import Nagios `.cfg` files into the database |

```json
{ "dir": "/etc/nagios/objects", "overwrite": false }
```

---

### Monitoring

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/monitoring/summary` | Object counts (hosts, services, contacts, …) |

---

### Logbook

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/logbook` | List audit log entries |

---

## Migrating from PHP NagiosQL

go-nagiosql is designed to be a drop-in backend replacement for the PHP application.

### Steps

1. **Keep the existing database.** go-nagiosql uses the same `tbl_*` tables as the PHP version. No data migration is required.

2. **Run `nagiosql migrate`.** This calls GORM's `AutoMigrate`, which adds any missing columns without dropping or altering existing data.

3. **Handle legacy passwords.** Accounts whose passwords are stored as MD5 hashes will trigger a `{requires_password_reset: true}` response on login — no token is issued. Use `PUT /api/v1/users/:id/password` to set a new bcrypt password.

4. **Point the reverse proxy.** Once the Go API is verified working, stop the PHP/nginx stack and proxy your domain to port `8081`.

5. **Config file format.** The Go generator produces byte-for-byte compatible output with the PHP NagiosQL config writer, so existing Nagios Core installations need no changes.

---

## Docker

### Production image

```bash
docker build -t go-nagiosql .
docker run -p 8081:8081 \
  -v $(pwd)/config.toml:/app/config.toml:ro \
  go-nagiosql
```

The image uses a two-stage Alpine build (`CGO_ENABLED=0`) and produces a ~16 MB binary with no external dependencies.

### Full Nagios stack (docker-compose)

`docker-compose.override.yml` adds `go-nagiosql` as a sidecar to the existing Nagios Core compose stack:

```bash
cp docker-compose.override.yml docker/nagios-core/docker-compose.override.yml
cd docker/nagios-core
docker compose up -d
```

Use `config.toml.docker` as the mounted config — it points paths to the shared `nagios-etc` and `nagios-var` volumes already defined in the Nagios Core compose file.

---

## Security Notes

| Topic | Behavior |
|-------|----------|
| **Passwords** | bcrypt cost 12 for all new/reset passwords. MD5 legacy hashes are never upgraded automatically. |
| **JWT secrets** | Must be at least 32 characters. Rotate by changing `jwt.secret` and restarting — all outstanding tokens are immediately invalidated. |
| **Refresh token** | Stored as an `httpOnly` cookie; never exposed to JavaScript. |
| **reload_trigger** | Must be a regular file, not `nagios.cmd`. Writing to the command FIFO from Go causes a deadlock. The reload watcher (supervisord or a shell loop) should read this file and issue the actual `nagios -s` signal. |
| **Config backups** | Every `.cfg` file overwrite is preceded by a timestamped backup in `nagios.backup_dir`. |
| **Admin endpoints** | User management and settings writes require the `admin` role embedded in the JWT claim. |

---

## PHP NagiosQL Docs

The `docs/` directory contains the original PHP NagiosQL infrastructure documentation:

| File | Description |
|------|-------------|
| [DIAGRAMA_NAGIOSQL.md](docs/DIAGRAMA_NAGIOSQL.md) | Full architecture diagram and 8-step config lifecycle |
| [INSTALL_NAGIOSCORE4.md](docs/INSTALL_NAGIOSCORE4.md) | Building and installing Nagios Core 4 from source |
| [INSTALL_NAGIOSQL.md](docs/INSTALL_NAGIOSQL.md) | Installing PHP NagiosQL |
| [BUILD_NAGIOS_TRIXIE.md](docs/BUILD_NAGIOS_TRIXIE.md) | Building on Debian Trixie |
| [BUILD_NAGIOS_TRIXIE_COMPLETE.md](docs/BUILD_NAGIOS_TRIXIE_COMPLETE.md) | Complete build walkthrough |
| [NAGIOS_BUILD_COMMANDS.md](docs/NAGIOS_BUILD_COMMANDS.md) | Reference build commands |
| [NAGIOSQL_DEBIAN_PACKAGING.md](docs/NAGIOSQL_DEBIAN_PACKAGING.md) | Debian packaging notes |
