## ADDED Requirements

### Requirement: Unit tests for JWT and auth service
The `internal/services/auth/` package SHALL have unit tests covering: JWT token generation, JWT token validation (valid/expired/malformed), bcrypt hash/compare, and legacy MD5 hash detection. Tests SHALL use only the standard `testing` package with no external dependencies.

#### Scenario: Token round-trip
- **WHEN** `GenerateAccessToken(claims)` is called and the result is passed to `ValidateToken(tokenString)`
- **THEN** the returned claims match the input without error

#### Scenario: Expired token rejected
- **WHEN** a token is generated with a past `exp` and passed to `ValidateToken`
- **THEN** an error is returned containing "token is expired"

#### Scenario: MD5 hash detected
- **WHEN** `IsLegacyMD5("5f4dcc3b5aa765d61d8327deb882cf99")` is called
- **THEN** it returns `true` (does not start with `$2`)

#### Scenario: bcrypt verify correct password
- **WHEN** `HashPassword("secret")` produces a hash and `CheckPassword("secret", hash)` is called
- **THEN** nil error is returned

### Requirement: Unit tests for config generation service
The `internal/services/nagconfig/` package SHALL have unit tests covering: template rendering for each object type, long-line continuation (>800 chars), per-entity file naming, Nagios 3 vs 4 directive selection, and custom variable injection.

#### Scenario: Host config output format
- **WHEN** `GenerateHost` is called with a `models.Host` fixture
- **THEN** the output string matches `define host {\n  host_name <name>\n  ...}\n` exactly

#### Scenario: Long-line continuation applied
- **WHEN** a check_command argument exceeds 800 characters
- **THEN** the output line ends with ` \` and continues on the next line

#### Scenario: Custom variable in output
- **WHEN** a host has a variable definition `_SNMP_COMMUNITY=public`
- **THEN** the generated `.cfg` contains `_SNMP_COMMUNITY public` inside the `define host { }` block

#### Scenario: Nagios 4 importance directive included
- **WHEN** `GenerateHost` is called with a v4 config target and `importance=10`
- **THEN** the output contains `importance 10`

#### Scenario: Nagios 3 target omits importance
- **WHEN** `GenerateHost` is called with a v3 config target and `importance=10`
- **THEN** the output does NOT contain `importance`

### Requirement: Unit tests for .cfg import parser
The `internal/services/nagimport/` package SHALL have unit tests covering: parsing of all 17 object types from fixture `.cfg` files, unknown directive skipping, multi-object single-file parsing, and relation name-to-ID resolution.

#### Scenario: Parse host definition
- **WHEN** the parser reads a `define host { host_name web01\n  address 10.0.0.1\n }` block
- **THEN** it returns a `ParsedHost{HostName:"web01", Address:"10.0.0.1"}` with no error

#### Scenario: Unknown directive skipped
- **WHEN** a `.cfg` block contains `unknown_directive value`
- **THEN** the object is parsed successfully with that field ignored

#### Scenario: Multi-block file
- **WHEN** a file contains 3 `define host { }` blocks
- **THEN** the parser returns a slice of 3 parsed hosts

### Requirement: Handler tests using httptest
All Echo route handlers in `internal/api/handlers/` SHALL have tests using `net/http/httptest`. Each handler test SHALL set up a mocked or in-memory GORM DB (using `gorm.io/driver/sqlite` for tests), create an Echo instance, and assert the HTTP status code and JSON response body.

#### Scenario: GET /api/v1/hosts returns 200
- **WHEN** the hosts handler is called via httptest with a valid JWT in the Authorization header
- **THEN** HTTP 200 is returned with a JSON array

#### Scenario: POST /api/v1/hosts with missing host_name returns 400
- **WHEN** the create host handler is called with `{"address":"10.0.0.1"}` (no host_name)
- **THEN** HTTP 400 is returned with `{"error":"validation failed","fields":{"host_name":"required"}}`

#### Scenario: DELETE /api/v1/hosts/:id with dependents returns 409
- **WHEN** a host referenced by a service is deleted
- **THEN** HTTP 409 is returned with a message listing the dependent service names

#### Scenario: Unauthenticated request returns 401
- **WHEN** any `/api/v1/` endpoint is called without Authorization header (JWT middleware active)
- **THEN** HTTP 401 is returned

### Requirement: Integration tests with MariaDB (build tag: integration)
Tests tagged with `//go:build integration` SHALL run against a live MariaDB instance (configured via `TEST_DSN` env var). They SHALL cover the full request lifecycle: login → CRUD → config write → verify.

#### Scenario: Full host lifecycle
- **WHEN** the integration test creates a host via API, writes its config, and runs `nagios -v`
- **THEN** `nagios -v` exits with code 0

#### Scenario: MD5 password migration
- **WHEN** a user with an MD5 hash in `tbl_user` attempts to login
- **THEN** the API returns `requires_password_reset: true`; after reset, a normal JWT login succeeds

### Requirement: Console curl smoke tests
`test/api/smoke.sh` SHALL be an executable bash script that tests the full API flow using `curl`. It SHALL require a running server (`BASE_URL` env var, default `http://localhost:8081`) and `jq`. The script SHALL exit non-zero if any step fails.

#### Scenario: Smoke test passes on a clean deployment
- **WHEN** `BASE_URL=http://localhost:8081 make test-api` is run against a fresh deployment
- **THEN** the script exits 0 after successfully: logging in, creating a host, listing hosts, updating the host, writing config, verifying config, deleting the host, and logging out

#### Scenario: Smoke test fails on wrong credentials
- **WHEN** the smoke test uses wrong credentials
- **THEN** the script exits non-zero with a clear error message on the failing step

### Requirement: Individual API shell scripts per resource
`test/api/` SHALL contain individual scripts: `auth.sh`, `hosts.sh`, `services.sh`, `commands.sh`, `timeperiods.sh`, `contacts.sh`, `verify.sh`. Each script SHALL demonstrate all CRUD operations for that resource using `curl` with Bearer token auth, suitable as copy-paste examples in documentation.

#### Scenario: hosts.sh runs standalone
- **WHEN** `TOKEN=<valid_jwt> bash test/api/hosts.sh` is run
- **THEN** it creates, reads, updates, and deletes a test host, printing each step's result

### Requirement: Test coverage ≥70% on services and handlers
`make test-cover` SHALL produce a `coverage.html` report. The coverage for `internal/services/` and `internal/api/handlers/` packages SHALL be ≥70%.

#### Scenario: Coverage threshold enforced
- **WHEN** `go test ./... -coverprofile=coverage.out` is run
- **THEN** `go tool cover -func=coverage.out` shows ≥70% total coverage for service and handler packages

### Requirement: Docker test environment with MariaDB 10.11
A `docker/test/docker-compose.test.yml` file SHALL provide a MariaDB 10.11 container for local and CI testing. The container SHALL be seeded with the NagiosQL schema and sample data. `make db-start` SHALL bring it up and wait until healthy. `make db-stop` SHALL tear it down with volume removal. `make test-integration` SHALL orchestrate `db-start → go test integration → db-stop` atomically.

#### Scenario: Test DB starts from scratch
- **WHEN** `make db-start` is run on a clean machine
- **THEN** MariaDB 10.11 starts on port 3307, the `nagiosql_test` database is created and seeded, and the healthcheck passes within 30 seconds

#### Scenario: Integration tests run against real DB
- **WHEN** `make test-integration` is run
- **THEN** `go test` with `-tags integration` connects to MariaDB on 3307, seeds fixtures, runs all integration tests, and cleans up

#### Scenario: CI pipeline
- **WHEN** `make ci` is run
- **THEN** both mariadb and nagiosql-test containers start; tests run inside the nagiosql-test container; the exit code of that container determines the overall result

### Requirement: Bash smoke test scripts cover all API endpoints
`test/api/` SHALL contain one bash script per resource group, plus a `smoke.sh` orchestrator and a `lib.sh` shared helper. Each script SHALL use `set -euo pipefail`, assert HTTP status codes, validate JSON fields with `jq`, and clean up created resources via a `trap` handler. `make test-api` SHALL run the full smoke suite against a live server.

#### Scenario: Smoke test exits 0 on a healthy backend
- **WHEN** `make test-api` is run with a running `nagiosql serve` and seeded DB
- **THEN** all 15 resource scripts execute successfully and `smoke.sh` exits 0

#### Scenario: Smoke test exits non-zero on API failure
- **WHEN** a bug causes `POST /api/v1/hosts` to return 500
- **THEN** `hosts.sh` prints `[FAIL] create host — expected 201 got 500` and exits non-zero, causing `smoke.sh` to exit non-zero

#### Scenario: Single script can run independently
- **WHEN** `BASE_URL=http://localhost:8081 TOKEN=$(bash test/api/lib.sh login) bash test/api/hosts.sh` is run
- **THEN** only the host CRUD tests execute and exit 0

### Requirement: go vet and golangci-lint pass with zero findings
`make vet` and `make lint` SHALL produce no output (zero findings). The golangci-lint configuration SHALL enable at minimum: `errcheck`, `gosimple`, `govet`, `ineffassign`, `staticcheck`, `unused`.

#### Scenario: Lint passes on fresh build
- **WHEN** `make lint` is run on a freshly cloned repository after `go mod download`
- **THEN** golangci-lint exits 0 with no findings printed
