## ADDED Requirements

### Requirement: JWT login endpoint
`POST /api/v1/auth/login` SHALL accept `{"username":"...","password":"..."}` and on success return HTTP 200 with `{"access_token":"<jwt>","expires_in":900,"token_type":"Bearer","user":{...}}` plus a `Set-Cookie: refresh_token=<jwt>; HttpOnly; Secure; SameSite=Strict; Max-Age=604800` header. On failure it SHALL return HTTP 401 with `{"error":"invalid credentials"}`.

#### Scenario: Successful login
- **WHEN** `POST /api/v1/auth/login` is called with correct username and password
- **THEN** HTTP 200 is returned with `access_token` (15-min JWT) in body and `refresh_token` (7-day JWT) as httpOnly cookie

#### Scenario: Wrong password
- **WHEN** `POST /api/v1/auth/login` is called with wrong password
- **THEN** HTTP 401 is returned; no tokens are issued; constant-time comparison is used to prevent timing attacks

#### Scenario: Unknown username
- **WHEN** `POST /api/v1/auth/login` is called with a username that does not exist
- **THEN** HTTP 401 is returned with the same error message as wrong password (no username enumeration)

#### Scenario: Console login with curl
- **WHEN** `curl -s -X POST http://localhost:8081/api/v1/auth/login -H 'Content-Type: application/json' -d '{"username":"admin","password":"admin"}' | jq -r .access_token` is run
- **THEN** the access token is printed to stdout and can be stored in `$TOKEN`

### Requirement: JWT access token as Bearer header
All protected endpoints SHALL require `Authorization: Bearer <access_token>`. The JWT middleware SHALL validate the HS256 signature, check expiry, and extract claims. Invalid or expired tokens SHALL return HTTP 401.

#### Scenario: Valid Bearer token
- **WHEN** `curl -H "Authorization: Bearer $TOKEN" http://localhost:8081/api/v1/hosts` is called
- **THEN** HTTP 200 is returned with the host list

#### Scenario: Expired access token
- **WHEN** an access token past its `exp` claim is used
- **THEN** HTTP 401 is returned with `{"error":"token expired"}`

#### Scenario: Malformed token
- **WHEN** a token with an invalid signature is used
- **THEN** HTTP 401 is returned with `{"error":"invalid token"}`

### Requirement: JWT refresh endpoint
`POST /api/v1/auth/refresh` SHALL read the `refresh_token` cookie, validate it, and return a new access token. The refresh token's expiry is NOT extended (no sliding window). On invalid/missing cookie it SHALL return HTTP 401.

#### Scenario: Refresh valid cookie
- **WHEN** `POST /api/v1/auth/refresh` is called with a valid `refresh_token` cookie
- **THEN** HTTP 200 is returned with a new `access_token`

#### Scenario: Missing refresh cookie
- **WHEN** `POST /api/v1/auth/refresh` is called without a cookie
- **THEN** HTTP 401 is returned

### Requirement: Logout endpoint
`POST /api/v1/auth/logout` SHALL clear the `refresh_token` cookie by setting `Max-Age=0`. Access tokens are short-lived and require no server-side revocation.

#### Scenario: Logout clears cookie
- **WHEN** `POST /api/v1/auth/logout` is called
- **THEN** HTTP 200 is returned and the `Set-Cookie` header sets `refresh_token` with `Max-Age=0`

### Requirement: Current user endpoint
`GET /api/v1/auth/me` SHALL return the authenticated user's profile decoded from the JWT claims without a database query.

#### Scenario: Me returns user from token
- **WHEN** `GET /api/v1/auth/me` is called with a valid Bearer token
- **THEN** HTTP 200 is returned with `{"id":1,"username":"admin","admin":true,"domain_id":1}`

### Requirement: bcrypt password hashing
User passwords SHALL be stored as bcrypt hashes (cost=12). The system SHALL detect legacy NagiosQL MD5 hashes (hash does not start with `$2`) and return `{"requires_password_reset":true}` on login instead of issuing tokens.

#### Scenario: New password stored as bcrypt
- **WHEN** a user's password is created or updated
- **THEN** the stored value starts with `$2a$12$` (bcrypt prefix)

#### Scenario: Legacy MD5 hash detected on login
- **WHEN** a user with an MD5-hashed password logs in with the correct password
- **THEN** HTTP 200 is returned with `{"requires_password_reset":true}` and no tokens are issued

#### Scenario: Password reset completes migration
- **WHEN** `POST /api/v1/auth/reset-password` is called with the MD5-verified old password and a new password
- **THEN** the new bcrypt hash is stored and a normal JWT login response is returned

### Requirement: JWT claims include domain and admin flag
The access token payload SHALL include: `sub` (user ID as string), `username`, `admin` (bool), `domain_id` (int), `iat`, `exp`. The middleware SHALL inject these claims into the Echo context under key `"claims"` for use by handlers without additional DB lookups.

#### Scenario: Claims accessible in handler
- **WHEN** any protected handler runs
- **THEN** `c.Get("claims")` returns a typed `*jwt.RegisteredClaims`-extended struct with username and admin flag

### Requirement: User management (CRUD)
`GET/POST /api/v1/users` and `GET/PUT/DELETE /api/v1/users/:id` SHALL provide full user CRUD. Only admin users (`admin=true` in JWT claims) SHALL access these endpoints.

#### Scenario: Create user
- **WHEN** an admin calls `POST /api/v1/users` with `{"username":"ops","password":"secret","admin":false}`
- **THEN** HTTP 201 is returned; the password is stored as bcrypt

#### Scenario: Non-admin blocked
- **WHEN** a non-admin JWT calls `POST /api/v1/users`
- **THEN** HTTP 403 is returned with `{"error":"forbidden"}`

#### Scenario: Self-deletion prevented
- **WHEN** a user calls `DELETE /api/v1/users/:id` with their own ID
- **THEN** HTTP 409 is returned with `{"error":"cannot delete own account"}`

### Requirement: Group management (CRUD)
`GET/POST /api/v1/groups` and `GET/PUT/DELETE /api/v1/groups/:id` SHALL manage user groups. `GET/PUT /api/v1/groups/:id/users` SHALL manage group membership. Admin-only endpoints.

#### Scenario: Assign user to group
- **WHEN** an admin calls `PUT /api/v1/groups/1/users` with `{"user_ids":[2,3]}`
- **THEN** HTTP 200 is returned and `tbl_lnkUserToGroup` is updated

### Requirement: JWT middleware protects all /api/v1/ routes except auth
The Echo router SHALL apply JWT validation middleware to all groups under `/api/v1/` EXCEPT `/api/v1/auth/login` and `/api/v1/auth/refresh`. The middleware SHALL reject missing, expired, or malformed tokens with HTTP 401 before the handler runs.

#### Scenario: Auth routes are public
- **WHEN** `POST /api/v1/auth/login` is called without any Authorization header
- **THEN** the request reaches the login handler (no 401 from middleware)

#### Scenario: All other routes require token
- **WHEN** `GET /api/v1/commands` is called without Authorization header
- **THEN** HTTP 401 is returned before the handler executes
