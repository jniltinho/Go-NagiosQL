## ADDED Requirements

### Requirement: Echo v5 router with /api/v1/ prefix
The REST API SHALL be served by Echo v5 with all routes grouped under `/api/v1/`. Routes SHALL be RESTful: `GET /api/v1/hosts`, `POST /api/v1/hosts`, `GET /api/v1/hosts/:id`, `PUT /api/v1/hosts/:id`, `DELETE /api/v1/hosts/:id`. The same pattern SHALL apply to all Nagios object types.

#### Scenario: List resource
- **WHEN** `GET /api/v1/hosts` is called with a valid session
- **THEN** HTTP 200 is returned with a JSON array of host objects and pagination metadata

#### Scenario: Create resource
- **WHEN** `POST /api/v1/hosts` is called with a valid JSON body
- **THEN** HTTP 201 is returned with the created host object

#### Scenario: Not found
- **WHEN** `GET /api/v1/hosts/999` is called and host id=999 does not exist
- **THEN** HTTP 404 is returned with `{"error": "not found"}`

### Requirement: Consistent JSON response format
All API responses SHALL use a consistent envelope format. Success responses with data SHALL return the resource directly or wrapped with pagination. Error responses SHALL always include an `error` string field.

#### Scenario: Paginated list response
- **WHEN** `GET /api/v1/hosts?page=2&limit=25` is called
- **THEN** the response includes `{"data": [...], "total": 100, "page": 2, "limit": 25}`

#### Scenario: Validation error response
- **WHEN** `POST /api/v1/hosts` is called with missing required fields
- **THEN** HTTP 400 is returned with `{"error": "validation failed", "fields": {"host_name": "required"}}`

### Requirement: Request validation
All POST and PUT endpoints SHALL validate the request body using Echo's binding and validation. Required fields SHALL return HTTP 400 with field-level error messages. Unknown fields SHALL be ignored.

#### Scenario: Invalid request body
- **WHEN** `POST /api/v1/commands` is called with a missing `command_name`
- **THEN** HTTP 400 is returned with a message identifying the missing field

### Requirement: JWT middleware on all API routes
All `/api/v1/` routes (except `/api/v1/auth/login` and `/api/v1/auth/refresh`) SHALL require a valid `Authorization: Bearer <token>` header. The middleware SHALL validate the HS256 signature, check expiry, and inject decoded claims into the Echo context. Missing, expired, or malformed tokens SHALL return HTTP 401.

#### Scenario: Missing Authorization header
- **WHEN** any protected `/api/v1/` request arrives without an Authorization header
- **THEN** HTTP 401 is returned with `{"error": "unauthorized"}`

#### Scenario: curl-testable auth
- **WHEN** `TOKEN=$(curl -s -X POST .../auth/login -d '{"username":"admin","password":"admin"}' | jq -r .access_token)` is run
- **THEN** subsequent `curl -H "Authorization: Bearer $TOKEN" .../api/v1/hosts` calls succeed with HTTP 200

### Requirement: CORS middleware
The API SHALL include CORS middleware configured to allow requests from the same origin. In development mode, it SHALL allow `http://localhost:5173` (Vite dev server).

#### Scenario: CORS preflight in development
- **WHEN** Vite dev server sends a preflight OPTIONS request to `/api/v1/hosts`
- **THEN** the server responds with appropriate CORS headers

### Requirement: Request logging middleware
All incoming requests SHALL be logged with method, path, status code, and response time using Echo's built-in logger middleware configured for structured output.

#### Scenario: Request logged
- **WHEN** any request is handled
- **THEN** a log line is emitted with method, path, status, and latency

### Requirement: Config generation and verify API endpoints
The API SHALL expose dedicated endpoints for config operations: `POST /api/v1/config/write` (write all configs), `POST /api/v1/config/write/:type/:id` (write single object), `POST /api/v1/config/verify` (run nagios -v), `POST /api/v1/config/restart` (restart Nagios).

#### Scenario: Write all configs via API
- **WHEN** `POST /api/v1/config/write` is called
- **THEN** config files for all active objects are generated and a summary response is returned

#### Scenario: Verify config via API
- **WHEN** `POST /api/v1/config/verify` is called
- **THEN** `nagios -v` output is returned as a JSON string with a `valid: true/false` field

### Requirement: Domain scoping on all resource endpoints
All resource list endpoints SHALL accept an optional `domain_id` query parameter. When provided, results SHALL be filtered to that domain plus domain 0 (common). When omitted, the user's default domain from the session SHALL be used.

#### Scenario: Domain-scoped host list
- **WHEN** `GET /api/v1/hosts?domain_id=2` is called
- **THEN** only hosts with `config_id=2` or `config_id=0` are returned

### Requirement: API routes for session management
`POST /api/v1/auth/login`, `POST /api/v1/auth/logout`, and `GET /api/v1/auth/me` SHALL be provided for session management. These routes SHALL NOT require prior authentication (except logout and me).

#### Scenario: Login via API
- **WHEN** `POST /api/v1/auth/login` is called with valid credentials
- **THEN** HTTP 200 is returned, the session cookie is set, and `{"user": {...}}` is returned

#### Scenario: Get current user
- **WHEN** `GET /api/v1/auth/me` is called with a valid session
- **THEN** the authenticated user's profile is returned

### Requirement: OpenAPI documentation served at /api/swagger/
All API handlers SHALL have swaggo doc comments. The generated OpenAPI spec SHALL be embedded in the binary and served at `/api/swagger/` via `swaggo/http-swagger`. The Swagger UI SHALL be accessible without authentication.

#### Scenario: Swagger UI accessible without login
- **WHEN** a browser visits `/api/swagger/` without a session
- **THEN** the Swagger UI renders with the full API listing

#### Scenario: All endpoints documented
- **WHEN** `swag init` is run
- **THEN** every `/api/v1/` route appears in `docs/swagger.json` with request/response schemas

### Requirement: Zerolog request logging middleware
Each request SHALL be logged by a zerolog-based Echo middleware with fields: method, path, status code, latency in milliseconds, and a unique request ID (UUID). The request ID SHALL also be set in the `X-Request-ID` response header.

#### Scenario: Request ID in response header
- **WHEN** any API request is handled
- **THEN** the response includes `X-Request-ID` header with a UUID value matching the log entry
