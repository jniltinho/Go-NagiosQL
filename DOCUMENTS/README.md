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
  - [Extended Object Types](#extended-object-types)
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
| PUT | `/api/v1/hostgroups/:id/members` | Add host to hostgroup |
| GET | `/api/v1/servicegroups` | List servicegroups |
| POST | `/api/v1/servicegroups` | Create servicegroup |
| DELETE | `/api/v1/servicegroups/:id` | Delete servicegroup |
| GET | `/api/v1/contactgroups` | List contactgroups |
| POST | `/api/v1/contactgroups` | Create contactgroup |
| DELETE | `/api/v1/contactgroups/:id` | Delete contactgroup |

---

### Extended Object Types

Full CRUD (GET / POST / GET:id / PUT:id / DELETE:id) for six additional Nagios object types:

#### Host Dependencies

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/hostdependencies` | List host dependencies |
| POST | `/api/v1/hostdependencies` | Create host dependency |
| GET | `/api/v1/hostdependencies/:id` | Get host dependency |
| PUT | `/api/v1/hostdependencies/:id` | Update host dependency |
| DELETE | `/api/v1/hostdependencies/:id` | Delete host dependency |

**Create (minimal):**
```json
{
  "config_name": "web-depends-on-db",
  "inherits_parent": 1,
  "execution_failure_criteria": "o",
  "notification_failure_criteria": "o"
}
```
FK fields (`dependent_host_name`, `host_name`, `dependent_hostgroup_name`, `hostgroup_name`, `dependency_period`) store the integer count flag (`0`=none, `1`=linked). The actual names are resolved from the corresponding `tbl_lnkHostdependencyTo*` link tables when generating `.cfg` files.

#### Host Escalations

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/hostescalations` | List host escalations |
| POST | `/api/v1/hostescalations` | Create host escalation |
| GET | `/api/v1/hostescalations/:id` | Get host escalation |
| PUT | `/api/v1/hostescalations/:id` | Update host escalation |
| DELETE | `/api/v1/hostescalations/:id` | Delete host escalation |

**Create (minimal):**
```json
{
  "config_name": "host-esc-admins",
  "first_notification": 3,
  "last_notification": 0,
  "notification_interval": 60,
  "escalation_options": "r,u"
}
```

#### Host Extended Info

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/hostextinfo` | List host extended info entries |
| POST | `/api/v1/hostextinfo` | Create host extended info entry |
| GET | `/api/v1/hostextinfo/:id` | Get host extended info entry |
| PUT | `/api/v1/hostextinfo/:id` | Update host extended info entry |
| DELETE | `/api/v1/hostextinfo/:id` | Delete host extended info entry |

**Create:**
```json
{
  "host_name": 4,
  "notes": "Primary web server",
  "notes_url": "http://wiki/web01",
  "icon_image": "base/linux40.png",
  "statusmap_image": "base/linux40.gd2"
}
```
`host_name` is the integer ID from `tbl_host`.

#### Service Dependencies

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/servicedependencies` | List service dependencies |
| POST | `/api/v1/servicedependencies` | Create service dependency |
| GET | `/api/v1/servicedependencies/:id` | Get service dependency |
| PUT | `/api/v1/servicedependencies/:id` | Update service dependency |
| DELETE | `/api/v1/servicedependencies/:id` | Delete service dependency |

**Create (minimal):**
```json
{
  "config_name": "http-depends-on-dns",
  "inherits_parent": 1,
  "execution_failure_criteria": "c,u",
  "notification_failure_criteria": "c,u"
}
```

#### Service Escalations

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/serviceescalations` | List service escalations |
| POST | `/api/v1/serviceescalations` | Create service escalation |
| GET | `/api/v1/serviceescalations/:id` | Get service escalation |
| PUT | `/api/v1/serviceescalations/:id` | Update service escalation |
| DELETE | `/api/v1/serviceescalations/:id` | Delete service escalation |

**Create (minimal):**
```json
{
  "config_name": "svc-esc-critical",
  "first_notification": 2,
  "last_notification": 0,
  "notification_interval": 60,
  "escalation_options": "w,u,c"
}
```

#### Service Extended Info

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/serviceextinfo` | List service extended info entries |
| POST | `/api/v1/serviceextinfo` | Create service extended info entry |
| GET | `/api/v1/serviceextinfo/:id` | Get service extended info entry |
| PUT | `/api/v1/serviceextinfo/:id` | Update service extended info entry |
| DELETE | `/api/v1/serviceextinfo/:id` | Delete service extended info entry |

**Create:**
```json
{
  "host_name": 4,
  "service_description": 7,
  "notes": "CPU Load service",
  "notes_url": "http://wiki/cpu-load",
  "icon_image": "cpu.png"
}
```
`host_name` and `service_description` are integer IDs from `tbl_host` and `tbl_service`.

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

`POST /config/write` generates **23 files** in total:

| File | Object type | Dir |
|------|-------------|-----|
| `commands.cfg` | Commands | base |
| `contactgroups.cfg` | Contact groups | base |
| `contacts.cfg` | Contacts | base |
| `contacttemplates.cfg` | Contact templates | base |
| `hostgroups.cfg` | Host groups | base |
| `hosttemplates.cfg` | Host templates | base |
| `hostdependencies.cfg` | Host dependencies | base |
| `hostescalations.cfg` | Host escalations | base |
| `hostextinfo.cfg` | Host extended info | base |
| `servicetemplates.cfg` | Service templates | base |
| `servicegroups.cfg` | Service groups | base |
| `servicedependencies.cfg` | Service dependencies | base |
| `serviceescalations.cfg` | Service escalations | base |
| `serviceextinfo.cfg` | Service extended info | base |
| `timeperiods.cfg` | Time periods | base |
| `hosts/<name>.cfg` | One per host | host_config_dir |
| `services/<name>.cfg` | One per host with services | service_config_dir |

All files include the standard NagiosQL header/footer with generation timestamp and version:

```
###############################################################################
#
# Host configuration file
#
# Created by: Go-NagiosQL Version v1.0.0
# Date:	      2026-06-17 18:00:00
# Version:    Nagios 4.x config file
#
# --- DO NOT EDIT THIS FILE BY HAND ---
# Nagios QL will overwite all manual settings during the next update
#
###############################################################################

define host {
	host_name                      	web01
	alias                          	Web Server 01
	address                        	10.0.0.1
	max_check_attempts             	3
	contact_groups                 	admins
	register                       	1
}

###############################################################################
#
# Host configuration file
#
# END OF FILE
#
###############################################################################
```

- Output is byte-for-byte compatible with the original PHP NagiosQL writer.
- FK fields are resolved via `tbl_lnkXxx` join tables (n:n) or direct ID lookup (1:1).
- `service_description` in service dependencies/escalations uses `strSlave` string storage.
- `servicegroup members` expands hostgroup references to individual host+service pairs.
- Fields with value `2` (inherit from template) are silently omitted.
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

## Development

For local setup, project structure, testing (unit / integration / smoke), Swagger regeneration, code style rules, and the checklist for adding new object types, see **[DEVELOPMENT.md](./DEVELOPMENT.md)**.

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
