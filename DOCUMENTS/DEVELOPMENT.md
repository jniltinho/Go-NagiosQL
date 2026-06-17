# go-nagiosql — Development Guide

This document covers local setup, project structure, testing, and contribution workflow.
For API reference and configuration, see [README.md](./README.md).

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Local Setup](#local-setup)
- [Project Structure](#project-structure)
- [Makefile Targets](#makefile-targets)
- [Testing](#testing)
  - [Unit Tests](#unit-tests)
  - [Integration Tests](#integration-tests)
  - [API Smoke Tests](#api-smoke-tests)
  - [CI Entry Point](#ci-entry-point)
- [Swagger / OpenAPI](#swagger--openapi)
- [Code Style](#code-style)
- [Adding a New Object Type](#adding-a-new-object-type)

---

## Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.26+ | Compiler |
| Docker | any | Test MariaDB, Nagios stack |
| `swag` | latest | Regenerate OpenAPI docs |
| `make` | any | Task runner |

Install `swag`:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

---

## Local Setup

```bash
# 1. Clone and configure
git clone <repo>
cd NagiosQL
cp config.toml.example config.toml
# Edit config.toml — set jwt.secret and database credentials

# 2. Start the test database (MariaDB on :3307)
make db-start

# 3. Migrate schema and seed admin user
go run . migrate --admin-password admin123 --sample

# 4. Build and run
make build
./bin/nagiosql serve
```

API available at `http://localhost:8081`.
Swagger UI at `http://localhost:8081/docs/swagger-ui.html`.

---

## Project Structure

```
.
├── cmd/                          # Cobra CLI commands (serve, migrate, config, import, version)
├── docs/                         # Generated OpenAPI spec (docs.go, swagger.json, swagger.yaml)
├── internal/
│   ├── api/
│   │   ├── handlers/             # Echo HTTP handlers (one file per resource group)
│   │   │   ├── auth.go           # Login, logout, refresh
│   │   │   ├── commands.go       # Commands CRUD
│   │   │   ├── contacts.go       # Contacts CRUD
│   │   │   ├── extended.go       # hostdependencies, hostescalations, hostextinfo,
│   │   │   │                     #   servicedependencies, serviceescalations, serviceextinfo
│   │   │   ├── groups.go         # hostgroups, servicegroups, contactgroups
│   │   │   ├── hosts.go          # Hosts CRUD
│   │   │   ├── services.go       # Services CRUD
│   │   │   ├── templates.go      # Host/service/contact templates
│   │   │   ├── timeperiods.go    # Timeperiods CRUD
│   │   │   └── users.go          # User management (admin only)
│   │   ├── middleware/           # JWT auth, RequireAdmin
│   │   └── routes.go             # All route registrations
│   ├── buildinfo/                # Build-time version propagation
│   ├── config/                   # config.toml loading (Viper)
│   ├── db/
│   │   ├── migrations/           # GORM AutoMigrate
│   │   └── seeds/                # Admin user seed
│   ├── models/                   # GORM models (one file per domain area)
│   │   ├── extended.go           # Hostdependency, Hostescalation, Hostextinfo,
│   │   │                         #   Servicedependency, Serviceescalation, Serviceextinfo
│   │   ├── groups.go             # Hostgroup, Servicegroup, Contactgroup
│   │   ├── links.go              # tbl_lnk* join table models
│   │   └── ...
│   ├── services/
│   │   ├── auth/                 # JWT issue/verify, bcrypt+MD5 password check
│   │   ├── logbook/              # Audit log writes
│   │   ├── nagconfig/            # .cfg file generation (Generator)
│   │   └── nagimport/            # .cfg file parser (import)
│   └── testhelpers/              # Shared test utilities (SQLite in-memory DB)
├── test/api/                     # Bash smoke tests (one script per resource)
└── DOCUMENTS/                    # Documentation
    ├── README.md                 # API reference
    └── DEVELOPMENT.md            # This file
```

---

## Makefile Targets

```bash
make build            # Compile → bin/nagiosql (with version ldflags)
make test             # Unit tests with race detector and coverage
make test-integration # Integration tests against live MariaDB (needs make db-start)
make test-api         # API smoke tests (builds + starts server + runs bash scripts)
make check            # vet + build + test  (CI entry point)
make swagger          # Regenerate OpenAPI docs from code annotations
make db-start         # Start test MariaDB container on :3307
make db-stop          # Stop test MariaDB (keeps volume)
make db-reset         # Wipe and restart test DB
make fmt              # gofmt + goimports
make vet              # go vet
make lint             # golangci-lint (if installed)
make clean            # Remove bin/
```

---

## Testing

The project has three independent test layers.

### Unit Tests

No external dependencies — use SQLite in-memory via `internal/testhelpers.NewDB`.

```bash
# Run everything
go test ./...

# With race detector + coverage
make test

# Specific package
go test ./internal/api/handlers/...
go test ./internal/services/nagconfig/...
go test ./internal/services/auth/...

# Single test by name
go test ./internal/api/handlers/... -run TestLogin_Success
go test ./internal/api/handlers/... -run TestHostdependency
go test ./internal/services/nagconfig/... -run TestWriteHost

# Verbose output
go test ./internal/api/handlers/... -v -run TestLogin
```

**Coverage by package:**

| Package | Test file(s) | What is covered |
|---------|-------------|-----------------|
| `internal/api/handlers` | `auth_test.go`, `hosts_test.go`, `services_test.go`, `commands_test.go`, `users_test.go`, `extended_test.go` | Auth (6 scenarios), hosts, services, commands, users, all 6 extended object types (37 tests) |
| `internal/services/nagconfig` | `generator_test.go` | Host field output, optional field omission, service groups, write-all, backup rotation |
| `internal/services/nagimport` | `parser_test.go` | `.cfg` parser: single object, multiple objects, inline comments, empty file |
| `internal/config` | `config_test.go` | Defaults, env var overrides, missing secret, short secret validation |
| `internal/services/auth` | (within auth package) | JWT issue/verify, bcrypt and MD5 password verification |

**Auth scenarios in detail:**

```bash
go test ./internal/api/handlers/... -v -run "TestLogin|TestLogout"
```

| Test | Scenario | Expected |
|------|----------|---------|
| `TestLogin_Success` | Valid bcrypt password | 200 + `access_token` |
| `TestLogin_WrongPassword` | Wrong password | 401 |
| `TestLogin_LegacyMD5_CorrectPassword` | PHP MD5 hash, correct password | 200 + `requires_password_reset: true` (no token) |
| `TestLogin_LegacyMD5_WrongPassword` | PHP MD5 hash, wrong password | 401 |
| `TestLogin_MissingFields` | Empty body | 400 |
| `TestLogout` | POST `/logout` | 204 + refresh cookie cleared |

**Extended object type scenarios (37 tests):**

```bash
go test ./internal/api/handlers/... -v -run \
  "TestHostdependency|TestHostescalation|TestHostextinfo|\
   TestServicedependency|TestServiceescalation|TestServiceextinfo"
```

Each of the 6 types follows the same pattern:

| Test suffix | What it validates |
|-------------|------------------|
| `List` | GET returns 200 with paginated results |
| `Create` | POST returns 201 with created record |
| `Create_MissingName` / `Create_MissingHost` | POST with missing required field → 400 |
| `Get` | GET /:id returns 200 |
| `Get_NotFound` | GET /9999 → 404 |
| `Update` | PUT /:id returns 200 with updated record |
| `Delete` | DELETE /:id → 204 |

**Adding tables for new tests:**

New model tables must be declared in `internal/testhelpers/db.go` inside the `schema()` function — SQLite DDL, one `CREATE TABLE IF NOT EXISTS` per model.

---

### Integration Tests

Require a live MariaDB instance. Use the Docker Compose stack in `docker/test/`.

```bash
# 1. Start test database (MariaDB 10.11 on port 3307)
make db-start

# 2. Run integration tests (build tag: integration)
make test-integration

# 3. Clean up
make db-stop    # keeps volume
make db-reset   # wipes volume and restarts
```

Integration tests live in `internal/integration/` and are guarded by the `//go:build integration` tag — they never run during `go test ./...`.

---

### API Smoke Tests

Bash scripts in `test/api/` that test the running HTTP server end-to-end via `curl`.

```bash
# Automated: build + start server + run all scripts + stop server
make test-api

# Manual: point at any running server
./bin/nagiosql serve &
BASE_URL=http://localhost:8081 bash test/api/smoke.sh

# Run a single script
BASE_URL=http://localhost:8081 bash test/api/auth.sh
BASE_URL=http://localhost:8081 bash test/api/hosts.sh
```

**Scripts:**

| Script | Endpoints exercised |
|--------|-------------------|
| `auth.sh` | `/auth/login`, `/auth/refresh`, `/auth/logout` |
| `hosts.sh` | CRUD `/hosts` |
| `services.sh` | CRUD `/services` |
| `commands.sh` | CRUD `/commands` |
| `contacts.sh` | CRUD `/contacts` |
| `groups.sh` | `/hostgroups`, `/servicegroups`, `/contactgroups` |
| `timeperiods.sh` | CRUD `/timeperiods` |
| `users.sh` | CRUD `/users`, password change |
| `import.sh` | `POST /import` |
| `monitoring.sh` | `GET /monitoring/summary` |
| `logbook.sh` | `GET /logbook` |

---

### CI Entry Point

```bash
make check   # vet + build + unit tests (no external deps required)
```

---

## Swagger / OpenAPI

Swagger docs are generated from code annotations using `swaggo/swag`.

```bash
# Regenerate after adding or changing handler annotations
make swagger

# Docs land in docs/
#   docs/docs.go       — registered via blank import in main.go
#   docs/swagger.json  — OpenAPI 3 spec
#   docs/swagger.yaml  — same, YAML format
```

Swagger UI is served at `/docs/swagger-ui.html` when the server is running.

**Rules when adding a new handler:**

1. Add `// FuncName godoc` as the first line before `@Summary` (required by `swag fmt`).
2. Include `@Security BearerAuth` on every protected endpoint.
3. Include `@Failure 400` on Create/Update, `@Failure 404` on Get/Update/Delete, `@Failure 500` on all mutating endpoints.
4. Run `make swagger` — never commit stale docs.

---

## Code Style

- TAB indentation everywhere (enforced by `gofmt`).
- `max(1, N-len(key))` instead of `if pad < 1 { pad = 1 }` (Go 1.21+ builtin).
- Slices and maps always initialized: `names := []string{}`, never `var names []string` (nil slices serialize to `null` in JSON).
- Functions with 4+ arguments use one argument per line.
- Early return for errors — happy path at minimal indentation.

See the full style guide in `.claude/skills/golang-code-style/`.

---

## Adding a New Object Type

Checklist for adding a new Nagios object type (e.g. `hostextendedstuff`):

1. **Model** — add struct in `internal/models/extended.go` with GORM tags and `TableName()`.
2. **Test schema** — add `CREATE TABLE IF NOT EXISTS tbl_xxx` in `internal/testhelpers/db.go`.
3. **Handler** — add 5 methods (List, Get, Create, Update, Delete) in `internal/api/handlers/extended.go` following the existing pattern. Include all swagger annotations.
4. **Routes** — register 5 routes in `internal/api/routes.go`.
5. **Tests** — add 7 test functions in `internal/api/handlers/extended_test.go` (List, Create, Create_MissingXxx, Get, Get_NotFound, Update, Delete).
6. **Config writer** — add a `WriteXxx(outputFile string) error` method in `internal/services/nagconfig/generator.go` and register it in `WriteAll()`.
7. **Swagger** — run `make swagger`.
8. **Docs** — add endpoint table to `DOCUMENTS/README.md` and update `DOCUMENTS/DEVELOPMENT.md` test coverage table.
