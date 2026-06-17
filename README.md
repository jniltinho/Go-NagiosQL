# Go-NagiosQL

A Go rewrite of [NagiosQL](https://www.nagiosql.org/) — a web-based configuration manager for Nagios/Icinga.

**Stack:** Go 1.26 · Echo v5 · GORM · Cobra · Viper · MariaDB · JWT

**96 REST endpoints** covering all Nagios object types · **23 `.cfg` files** generated from DB · byte-for-byte compatible with PHP NagiosQL output

---

## Quick Start

```bash
# 1. Configure
cp config.toml.example config.toml
# Edit config.toml — change jwt.secret and database credentials

# 2. Migrate database
go run . migrate --admin-password yourpassword --sample

# 3. Run
make build
./bin/nagiosql serve
```

API is available at `http://localhost:8081`.
Swagger UI at `http://localhost:8081/docs/swagger-ui.html`.

---

## CLI

```
nagiosql serve              Start the HTTP API server
nagiosql migrate            AutoMigrate schema; seed admin user
nagiosql import             Import Nagios .cfg files into the database
nagiosql config write       Write current config to stdout
nagiosql config verify      Run nagios -v nagios.cfg
nagiosql config restart     Touch reload_trigger to signal reload
nagiosql version            Print version and build date
```

---

## Development

```bash
make build             # compile → bin/nagiosql
make test              # unit tests (SQLite in-memory, no external deps)
make db-start          # start MariaDB 10.11 on :3307 via Docker
make test-integration  # integration tests against the test DB
make check             # vet + build + test  (CI entry point)
make swagger           # regenerate OpenAPI docs
```

---

## Documentation

Detailed reference lives in [`DOCUMENTS/`](DOCUMENTS/README.md):

| Document | Description |
|----------|-------------|
| [Configuration reference](DOCUMENTS/README.md#configuration) | All `config.toml` keys and environment variables |
| [API reference](DOCUMENTS/README.md#api-reference) | All 96 endpoints with request/response examples |
| [Extended object types](DOCUMENTS/README.md#extended-object-types) | hostdependencies, hostescalations, hostextinfo, servicedependencies, serviceescalations, serviceextinfo |
| [Config generation](DOCUMENTS/README.md#config-generation) | 23 generated `.cfg` files, format details, FK resolution |
| [Migrating from PHP NagiosQL](DOCUMENTS/README.md#migrating-from-php-nagiosql) | Schema compatibility, legacy passwords, cutover steps |
| [Docker deployment](DOCUMENTS/README.md#docker) | Production image and full Nagios stack compose |
| [Security notes](DOCUMENTS/README.md#security-notes) | Passwords, JWT, reload trigger |
| [Development guide](DOCUMENTS/DEVELOPMENT.md) | Local setup, testing (unit/integration/smoke), Swagger, code style, adding new object types |
| [Architecture — PHP NagiosQL](DOCUMENTS/docs/DIAGRAMA_NAGIOSQL.md) | Original system diagram and config lifecycle |

---

## Architecture Diagrams

<p align="center">
  <img src="DOCUMENTS/screenshots/go-nagiosql-arch.svg" width="640" alt="Go-NagiosQL Architecture"/>
</p>
<p align="center"><em>Go-NagiosQL — API flow from JWT auth to Nagios config generation</em></p>

<p align="center">
  <img src="DOCUMENTS/screenshots/php-nagiosql-arch.svg" width="640" alt="PHP NagiosQL Architecture"/>
</p>
<p align="center"><em>PHP NagiosQL — original container stack (nginx · php-fpm · supervisord)</em></p>

<p align="center">
  <img src="DOCUMENTS/screenshots/nagios-core-arch.svg" width="640" alt="Nagios Core 4 Architecture"/>
</p>
<p align="center"><em>Nagios Core 4 — config parser · scheduler · checks · notifications</em></p>

---

## License

MIT
