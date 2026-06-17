## ADDED Requirements

### Requirement: Go module and binary
The project SHALL be a Go 1.26 module named `go-nagiosql` with a binary named `nagiosql`, following the structure of `github.com/jniltinho/go-postfixadmin`.

#### Scenario: Binary is produced
- **WHEN** `make build` is executed
- **THEN** a single binary named `nagiosql` is produced in the project root

#### Scenario: Module name is correct
- **WHEN** `go.mod` is inspected
- **THEN** the module name is `github.com/jniltinho/go-nagiosql` and `go 1.26` is declared

### Requirement: Cobra CLI with serve, migrate, version, config, and import commands
The binary SHALL expose Cobra commands: `serve` (start HTTP server), `migrate` (run DB migrations), `version` (print version string), `config write|verify|restart` (scripting interface replacing `do_config.php`), and `import [file]` (CLI import). Running the binary with no arguments SHALL default to `serve`.

#### Scenario: Default command is serve
- **WHEN** `nagiosql` is run with no arguments
- **THEN** the HTTP server starts on the configured port

#### Scenario: Version command
- **WHEN** `nagiosql version` is run
- **THEN** the version string, build date, and Go version are printed to stdout

#### Scenario: Migrate command
- **WHEN** `nagiosql migrate` is run with a valid database connection
- **THEN** all GORM auto-migrations execute and exit with code 0

#### Scenario: Config write via CLI
- **WHEN** `nagiosql config write all` is run
- **THEN** all active config files are generated and a summary is printed to stdout using go-pretty table formatting

#### Scenario: Import via CLI
- **WHEN** `nagiosql import /path/to/hosts.cfg` is run
- **THEN** the file is parsed and imported into the database; a result summary is printed

### Requirement: Viper configuration via config.toml
The application SHALL read configuration from a `config.toml` file. Viper SHALL support environment variable overrides for every config key using the `NAGIOSQL_` prefix. The config file path SHALL be overridable via `--config` CLI flag.

#### Scenario: Config file is loaded
- **WHEN** `config.toml` exists in the working directory and `nagiosql serve` is run
- **THEN** all settings from the file are applied (port, DB credentials, Nagios paths)

#### Scenario: Environment variable override
- **WHEN** `NAGIOSQL_DATABASE_PASSWORD=secret nagiosql serve` is run without a config file
- **THEN** the database password is set to `secret`

#### Scenario: Missing config file
- **WHEN** no `config.toml` is present and no env vars are set
- **THEN** the application logs a clear error message and exits with a non-zero code

### Requirement: GORM with MariaDB connection
The application SHALL connect to MariaDB using GORM. The connection string SHALL be constructed from `config.toml` values. The application SHALL fail fast at startup if the database is unreachable.

#### Scenario: Successful DB connection
- **WHEN** valid MariaDB credentials are configured
- **THEN** GORM connects successfully and the server starts

#### Scenario: DB connection failure
- **WHEN** the database is unreachable at startup
- **THEN** the application logs the error and exits with a non-zero code

### Requirement: Vue.js 3 frontend embedded in binary
The compiled Vue.js 3 frontend (`frontend/dist/`) SHALL be embedded into the Go binary using `//go:embed`. The Echo server SHALL serve these static files from the embedded filesystem.

#### Scenario: Frontend served from binary
- **WHEN** a browser requests `/` from the running binary
- **THEN** the Vue.js index page is served without any external file access

#### Scenario: API routes are not shadowed
- **WHEN** a request is made to `/api/v1/hosts`
- **THEN** the route is handled by the Echo handler, not the static file server

### Requirement: zerolog structured logging
The application SHALL use `github.com/rs/zerolog` for all logging. Log output SHALL be JSON-formatted in production mode and pretty-printed in development mode (controlled by `[server].dev = true` in `config.toml`). A custom Echo middleware SHALL log each request with fields: method, path, status, latency, request_id.

#### Scenario: Request logged as JSON
- **WHEN** a request is handled in production mode
- **THEN** a single JSON log line is emitted with fields: `level`, `method`, `path`, `status`, `latency_ms`, `request_id`

#### Scenario: Startup logged
- **WHEN** the server starts
- **THEN** a zerolog INFO entry is emitted with the listening port and version

### Requirement: OpenAPI documentation via swaggo
The REST API SHALL be documented using `github.com/swaggo/swag` annotations on handler functions. The generated spec SHALL be committed to `docs/` and served at `/api/swagger/` by `swaggo/http-swagger`. The `make swagger` target SHALL regenerate the spec.

#### Scenario: Swagger UI accessible
- **WHEN** a browser visits `/api/swagger/`
- **THEN** the interactive Swagger UI is rendered showing all API endpoints

#### Scenario: Spec regenerated
- **WHEN** `make swagger` is run
- **THEN** `docs/swagger.json` and `docs/swagger.yaml` are updated from handler annotations

### Requirement: Makefile for build and development
A `Makefile` SHALL provide targets: `build` (compile binary), `dev` (run with live reload), `frontend` (build Vue.js), `all` (build frontend then binary), `swagger` (regenerate OpenAPI spec), `upx` (compress binary with UPX), `clean` (remove build artifacts).

#### Scenario: Full build
- **WHEN** `make all` is run
- **THEN** the Vue.js frontend is built first, then embedded into the Go binary

### Requirement: Dockerfile for container deployment
A `Dockerfile` SHALL produce a minimal image containing only the `nagiosql` binary. The image SHALL be based on `debian:trixie-slim` matching the existing stack.

#### Scenario: Container starts correctly
- **WHEN** the Docker image is run with a valid `config.toml` mounted or env vars set
- **THEN** the HTTP server starts and responds on the configured port
