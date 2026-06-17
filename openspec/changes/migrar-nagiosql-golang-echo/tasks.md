## 1. Project Foundation

- [ ] 1.1 Initialize Go module `github.com/jniltinho/go-nagiosql` in `/home/nilton/Projetos/nilton/NagiosQL/go-nagiosql/` (do NOT modify `DOCUMENTS/`)
- [ ] 1.2 Add Go 1.26 dependencies: `labstack/echo/v5`, `gorm.io/gorm`, `gorm.io/driver/mysql`, `gorm.io/driver/sqlite` (tests only), `spf13/cobra`, `spf13/viper`, `golang-jwt/jwt/v5`, `golang.org/x/crypto`, `jlaffaye/ftp`, `golang.org/x/crypto/ssh`, `swaggo/swag`, `swaggo/http-swagger`, `jedib0t/go-pretty`, `leonelquinteros/gotext`
- [ ] 1.3 Create `cmd/root.go` — Cobra root with `--config` flag (default `config.toml`) and `--dev` flag for pretty logs
- [ ] 1.4 Create `cmd/serve.go` — start Echo server; wire middleware, routes, and Swagger UI
- [ ] 1.5 Create `cmd/migrate.go` — Cobra `migrate` command: runs `db.Migrate()` then `seeds.SeedRequired(db, cfg)`; flag `--sample` also runs `seeds.SeedSample(db)`; flag `--admin-user` (default `admin`) and `--admin-password` for the initial admin; exits 0 on success; idempotent (safe to run multiple times)
- [ ] 1.6 Create `cmd/version.go` — print `Version`, `BuildDate`, `GoVersion` injected via `-ldflags`
- [ ] 1.7 Create `cmd/config.go` — Cobra subcommands: `nagiosql config write [type|all]`, `nagiosql config verify`, `nagiosql config restart`
- [ ] 1.8 Create `cmd/import.go` — `nagiosql import <file> [--overwrite]`
- [ ] 1.9 Create `internal/config/config.go` — Viper load: `[server]`, `[jwt]`, `[database]`, `[nagios]` sections; `NAGIOSQL_` env overrides
- [ ] 1.10 Create `config.toml.example` with all sections: `[jwt]` (`secret`, `access_ttl_min=15`, `refresh_ttl_days=7`); `[nagios]` fields including `reload_trigger="/usr/local/nagios/var/reload.trigger"` (NOT nagios.cmd — see design D9/D13), `host_config_dir`, `service_config_dir`, `backup_dir`, `import_dir`, `resource_file`, `cgi_file`, `pid_file`
- [ ] 1.11 Create `internal/db/db.go` — GORM DSN builder + `Open()` with fail-fast on unreachable DB; set `sql_mode='NO_ENGINE_SUBSTITUTION'` and charset utf8
- [ ] 1.12 Configure Echo v5 built-in logger middleware: `e.Use(middleware.Logger())` in `cmd/serve.go`; use `middleware.RequestID()` for `X-Request-ID` header; set `e.Logger.SetLevel(log.DEBUG)` when `cfg.Server.Dev=true`; no custom logger package needed
- [ ] 1.13 Create `embed.go` with `//go:embed docs` (no frontend in this phase)
- [ ] 1.14 Create comprehensive `Makefile` with targets: `all`, `build`, `test`, `test-v`, `test-cover`, `test-integration`, `test-api`, `swagger`, `lint`, `vet`, `upx`, `clean`, `run`, `migrate`; embed `VERSION` and `BUILD_DATE` via `-ldflags`; include `.PHONY` and `## help` autodoc
- [ ] 1.15 Create `Dockerfile` (single-stage Go build + debian:trixie-slim runtime; no node stage)
- [ ] 1.16 Create `.golangci.yml` enabling: `errcheck`, `gosimple`, `govet`, `ineffassign`, `staticcheck`, `unused`, `gofmt`
- [ ] 1.17 Scaffold `test/api/` directory with placeholder `smoke.sh`, `auth.sh`, `hosts.sh`; create `test/fixtures/cfg/` with the real `.cfg` files copied from the reference stack (read-only reference material — copy, never import at runtime):
  - `test/fixtures/cfg/linux-host.cfg` ← copy of `DOCUMENTS/docker/nagios-core/nagios/etc-extra/nagiosql/hosts/linux-host.cfg`
  - `test/fixtures/cfg/gateway.cfg`, `google-dns.cfg`, `cloudflare-dns.cfg` ← same source dir
  - `test/fixtures/cfg/ping.cfg`, `dns.cfg`, `http.cfg`, `ssh.cfg` ← `etc-extra/nagiosql/services/`
  - `test/fixtures/cfg/commands.cfg` ← `DOCUMENTS/docker/nagios-core/nagios/etc-extra/nagiosql/commands.cfg`
  - **No SQL files needed** — schema is created by `db.Migrate()` and data by Go seed functions in `internal/db/seeds/`
- [ ] 1.18 Run `make vet` — zero errors; run `make build` — binary compiles and `nagiosql version` prints

## 2. GORM Models

- [ ] 2.1 `internal/models/host.go` — `tbl_host`; all Nagios host directives; `ENGINE=MyISAM`; `active ENUM('0','1')`, `last_modified DATETIME`, `config_id TINYINT`
- [ ] 2.2 `internal/models/service.go` — `tbl_service`; `config_name` field for per-group file generation
- [ ] 2.3 `internal/models/command.go` — `tbl_command`; `command_type` (1=check, 2=notify); `arg1_info`–`arg8_info`
- [ ] 2.4 `internal/models/timeperiod.go` + `timedefinition.go` — `tbl_timeperiod` + `tbl_timedefinition` (`tipId` FK)
- [ ] 2.5 `internal/models/contact.go` — `tbl_contact`; all address fields `address1`–`address6`; `statistik_url`
- [ ] 2.6 `internal/models/groups.go` — `tbl_contactgroup`, `tbl_hostgroup`, `tbl_servicegroup`
- [ ] 2.7 `internal/models/templates.go` — `tbl_hosttemplate`, `tbl_servicetemplate`, `tbl_contacttemplate`; default `register='0'`
- [ ] 2.8 `internal/models/links.go` — all `tbl_lnk*` join table structs (M:N relations)
- [ ] 2.9 `internal/models/advanced.go` — `tbl_hostdependency`, `tbl_servicedependency`, `tbl_hostescalation`, `tbl_serviceescalation`, `tbl_hostextinfo`, `tbl_serviceextinfo`; `import_hash VARCHAR(255)`
- [ ] 2.10 `internal/models/domain.go` + `configtarget.go` — `tbl_datadomain` (id=0 "common"); `tbl_configtarget` with ALL fields from schema: `target`, `alias`, `server`, `port`, `method` (1=local/2=FTP/3=SSH), `user`, `password`, `ssh_key_path`, `ftp_secure`, `basedir`, `hostconfig`, `serviceconfig`, `backupdir`, `hostbackup`, `servicebackup`, `nagiosbasedir`, `importdir`, `picturedir`, `commandfile` (= reload.trigger), `binaryfile`, `pidfile`, `conffile`, `cgifile`, `resourcefile`, `version` (3 or 4); default row `target='localhost'` populated by `nagiosQL_v35_db_mysql.sql` and patched by entrypoint.sh
- [ ] 2.11 `internal/models/user.go` + `group.go` + `user_group.go` — `tbl_user` (MD5 password field preserved), `tbl_group`, `tbl_lnkUserToGroup`
- [ ] 2.12 `internal/models/meta.go` — `tbl_settings`, `tbl_tablestatus`, `tbl_logbook`, `tbl_info`, `tbl_menu`
- [ ] 2.13 `internal/models/variabledefinition.go` — `tbl_variabledefinition` (`varName`, `varValue`, `varType`); 6 link tables from schema: `tbl_lnkHostToVariabledefinition`, `tbl_lnkServiceToVariabledefinition`, `tbl_lnkHosttemplateToVariabledefinition`, `tbl_lnkServicetemplateToVariabledefinition`, `tbl_lnkContactToVariabledefinition`, `tbl_lnkContacttemplateToVariabledefinition`
- [ ] 2.14 `internal/db/migrations.go` — `AllModels() []any` returns every model pointer in dependency order (no FK issue with MyISAM but keep order for clarity); `Migrate(db *gorm.DB) error` calls `db.Set("gorm:table_options","ENGINE=MyISAM CHARSET=utf8 COLLATE=utf8_unicode_ci").AutoMigrate(AllModels()...)`

## 3. Database Seeds (Go — no SQL files at runtime)

All data from `nagiosQL_v35_db_mysql.sql` (schema) and `import_nagios_sample.sql` (sample data) is expressed as Go structs and seeded via GORM. No `.sql` file is read at runtime. Seeds are idempotent: use `db.FirstOrCreate` or `db.Save` with known IDs.

### Required seeds — always run by `nagiosql migrate`

- [ ] 3.1 `internal/db/seeds/domain.go` — `SeedDomains(db) error`: seed `tbl_datadomain` with two rows:
  - `id=0`: `domain="common"`, `alias="Global common domain"`, `nodelete='1'`, `active='1'` — use raw `db.Exec("INSERT INTO tbl_datadomain ... ON DUPLICATE KEY UPDATE alias=alias")` because GORM skips zero-value PKs; follow with `db.Exec("UPDATE tbl_datadomain SET id=0 WHERE domain='common'")` and `db.Exec("ALTER TABLE tbl_datadomain AUTO_INCREMENT=1")`
  - `id=1`: `domain="localhost"`, `alias="Local installation"`, `targets=1`, `enable_common=1`, `nodelete='1'`  — use `db.FirstOrCreate`

- [ ] 3.2 `internal/db/seeds/configtarget.go` — `SeedConfigTarget(db *gorm.DB, cfg *config.Config) error`: seed one `tbl_configtarget` row with `target="localhost"`, `nodelete='1'`; all paths sourced from `cfg.Nagios.*`:
  - `basedir` ← `cfg.Nagios.BaseDir + "/nagiosql/"`  (e.g. `/usr/local/nagios/etc/nagiosql/`)
  - `hostconfig` ← `cfg.Nagios.HostConfigDir`
  - `serviceconfig` ← `cfg.Nagios.ServiceConfigDir`
  - `backupdir` ← `cfg.Nagios.BackupDir`
  - `hostbackup` ← `cfg.Nagios.BackupDir + "hosts/"`
  - `servicebackup` ← `cfg.Nagios.BackupDir + "services/"`
  - `nagiosbasedir` ← `cfg.Nagios.BaseDir + "/etc/"`
  - `importdir` ← `cfg.Nagios.ImportDir`
  - `commandfile` ← `cfg.Nagios.ReloadTrigger`
  - `binaryfile` ← `cfg.Nagios.Binary`
  - `pidfile` ← `cfg.Nagios.PidFile`
  - `conffile` ← `cfg.Nagios.ConfigFile`
  - `cgifile` ← `cfg.Nagios.CgiFile`
  - `resourcefile` ← `cfg.Nagios.ResourceFile`
  - `version` ← `4`; `method` ← `"1"` (local); use `db.Where("target=?","localhost").FirstOrCreate`

- [ ] 3.3 `internal/db/seeds/settings.go` — `SeedSettings(db) error`: seed `tbl_settings`:
  - `{category:"install", name:"hash", value: sha256(uuid)}` — generate UUID in Go (`crypto/rand`), SHA256 hex; use `FirstOrCreate` on `name`
  - `{category:"db", name:"version", value:"3.5.0"}` — use `FirstOrCreate`

- [ ] 3.4 `internal/db/seeds/admin.go` — `SeedAdmin(db *gorm.DB, username, plainPassword string) error`:
  - Hash `plainPassword` using `bcrypt.GenerateFromPassword([]byte(plain), 12)` — result starts with `$2a$12$`
  - Insert `tbl_user` row: `admin_enable='1'`, `active='1'`, `nodelete='1'`, `language="1"`, `domain=1`
  - Use `db.Where("username=?", username).FirstOrCreate` — idempotent, never overwrites on subsequent `migrate` runs
  - **NEVER store MD5** — the legacy PHP entrypoint seeded with `MD5('$NAGIOSQL_PASSWORD')`; our Go seed always uses bcrypt from the start
  - Log: `log.Printf("admin user %q seeded with bcrypt hash", username)`

- [ ] 3.5 `internal/db/seeds/seeds.go` — `SeedRequired(db *gorm.DB, cfg *config.Config, adminUser, adminPass string) error`: calls 3.1 → 3.2 → 3.3 → 3.4 in order; logs each step via `log.Printf`

### Sample seeds — run only with `nagiosql migrate --sample`

All IDs match `import_nagios_sample.sql` exactly so that existing NagiosQL databases (migrated from PHP) are compatible.

- [ ] 3.6 `internal/db/seeds/sample_commands.go` — `SeedCommands(db) error`: seed all 24 `tbl_command` rows from `import_nagios_sample.sql`:
  - IDs 1–2: notify commands (type=2): `notify-host-by-email`, `notify-service-by-email`
  - IDs 3–22: check commands (type=1): `check-host-alive`, `check_local_disk`, `check_local_load`, `check_local_procs`, `check_local_users`, `check_local_swap`, `check_local_mrtgtraf`, `check_ftp`, `check_hpjd`, `check_snmp`, `check_http`, `check_ssh`, `check_dhcp`, `check_ping`, `check_pop`, `check_imap`, `check_smtp`, `check_tcp`, `check_udp`, `check_nt`
  - IDs 23–24: perf data commands (type=2): `process-host-perfdata`, `process-service-perfdata`
  - Use `db.Save` (upsert by PK); all `config_id=1`, `active='1'`, `register='1'`

- [ ] 3.7 `internal/db/seeds/sample_timeperiods.go` — `SeedTimeperiods(db) error`: seed 5 `tbl_timeperiod` rows + all `tbl_timedefinition` rows + 1 `tbl_lnkTimeperiodToTimeperiodUse` row:
  - Period id=1: `24x7` + 7 definitions (monday–sunday `00:00-24:00`)
  - Period id=2: `workhours` + 5 definitions (monday–friday `09:00-17:00`)
  - Period id=3: `none` (no definitions)
  - Period id=4: `us-holidays` + 3 definitions (jan 1, dec 25, jul 4 all `00:00-00:00`) + 2 dynamic (monday 1 september, thursday -1 november)
  - Period id=5: `24x7_sans_holidays` + 7 definitions (same as 24x7) + link to exclude period id=4

- [ ] 3.8 `internal/db/seeds/sample_contacts.go` — `SeedContacts(db) error`:
  - `tbl_contacttemplate` id=1: `generic-contact`, `register='0'`; link to `tbl_lnkContacttemplateToCommandHost` (idSlave=1) and `tbl_lnkContacttemplateToCommandService` (idSlave=2)
  - `tbl_contactgroup` id=1: `admins`, `alias="Nagios Administrators"`, `members=1`
  - `tbl_contact` id=1: `nagiosadmin`, `alias="Nagios Admin"`, `email="nagios@localhost"`, `register='1'`; links: `tbl_lnkContactToContacttemplate` (idSlave=1), `tbl_lnkContactgroupToContact` (idMaster=1 idSlave=1)

- [ ] 3.9 `internal/db/seeds/sample_hosttemplates.go` — `SeedHostTemplates(db) error`: seed 5 `tbl_hosttemplate` rows + all their link rows:
  - id=1: `generic-host`, `register='0'` — no contact groups
  - id=2: `linux-server`, `register='0'`, `use_template→id=1`, contact_group admins (id=1), `max_check_attempts=10`, `check_interval=5`
  - id=3: `windows-server`, `register='0'`, `use_template→id=1`, contact_group admins, hostgroup windows-servers (id=1)
  - id=4: `generic-printer`, `register='0'`, `use_template→id=1`, contact_group admins
  - id=5: `generic-switch`, `register='0'`, `use_template→id=1`, contact_group admins
  - Link rows: `tbl_lnkHosttemplateToContactgroup`, `tbl_lnkHosttemplateToHosttemplate`, `tbl_lnkHosttemplateToHostgroup`

- [ ] 3.10 `internal/db/seeds/sample_servicetemplates.go` — `SeedServiceTemplates(db) error`: seed 2 `tbl_servicetemplate` rows:
  - id=1: `generic-service`, `register='0'`, `max_check_attempts=3`, `check_interval=10`, `notification_interval=60`, `notification_options="w,u,c,r"`, `active_checks_enabled=1`, `process_perf_data=1`
  - id=2: `local-service`, `register='0'`, `use_template→id=1`, `max_check_attempts=4`, `check_interval=5`
  - Links: `tbl_lnkServicetemplateToContactgroup` (id=1→admins), `tbl_lnkServicetemplateToServicetemplate` (id=2→1)

- [ ] 3.11 `internal/db/seeds/sample_hostgroups.go` — `SeedHostGroups(db) error`: seed 4 `tbl_hostgroup` rows:
  - id=1: `windows-servers`, `alias="Windows Servers"`, `members=1`
  - id=2: `switches`, `alias="Network Switches"`, `members=1`
  - id=3: `network-printers`, `alias="Network Printers"`, `members=1`
  - id=4: `linux-servers`, `alias="Linux Servers"`, `members=1`

- [ ] 3.12 `internal/db/seeds/sample_hosts.go` — `SeedHosts(db) error`: seed 4 `tbl_host` rows + all link rows:
  - id=1: `winserver`, `address=192.168.1.2`, `use_template→windows-server(id=3)`, `config_id=1`
  - id=2: `linksys-srw224p`, `address=192.168.1.253`, `use_template→generic-switch(id=5)`, `config_id=1`; link to hostgroup switches(id=2)
  - id=3: `hplj2605dn`, `address=192.168.1.30`, `use_template→generic-printer(id=4)`, `config_id=1`; link to hostgroup network-printers(id=3)
  - id=4: `localhost`, `address=127.0.0.1`, `use_template→linux-server(id=2)`, `config_id=1`; link to hostgroup linux-servers(id=4) via `tbl_lnkHostgroupToHost`
  - All link tables: `tbl_lnkHostToHosttemplate`, `tbl_lnkHostToHostgroup`, `tbl_lnkHostgroupToHost`

- [ ] 3.13 `internal/db/seeds/sample_services.go` — `SeedServices(db) error`: seed all 21 `tbl_service` rows with their `tbl_lnkServiceToHost` and `tbl_lnkServiceToServicetemplate` links:
  - IDs 1–7: `config_name="winserver"`, host_id=1 (winserver); services: NSClient++ Version, Uptime, CPU Load, Memory Usage, C:\\ Drive Space, W3SVC, Explorer — all use template generic-service(id=1)
  - IDs 8–11: `config_name="linksys-srw224p"`, host_id=2; services: PING, Uptime, Port 1 Link Status, Port 1 Bandwidth Usage — use generic-service(id=1)
  - IDs 12–13: `config_name="hplj2605dn"`, host_id=3; services: Printer Status, PING — use generic-service(id=1)
  - IDs 14–21: `config_name="localhost"`, host_id=4; services: PING, Root Partition, Current Users, Total Processes, Current Load, Swap Usage, SSH, HTTP — IDs 14–18 use generic-service(id=1), IDs 19–21 use local-service(id=2)

- [ ] 3.14 `internal/db/seeds/sample.go` — `SeedSample(db) error`: orchestrates 3.6 → 3.7 → 3.8 → 3.9 → 3.10 → 3.11 → 3.12 → 3.13; wraps in a DB transaction; if any step fails, rolls back and returns error; logs count of seeded rows via `log.Printf`

## 5. Auth Service and JWT Middleware

- [ ] 5.1 `internal/services/auth/jwt.go` — `GenerateAccessToken(userID, username, admin, domainID) (string, error)` HS256, 15-min TTL; `GenerateRefreshToken(userID) (string, error)` HS256, 7-day TTL; `ValidateToken(tokenString) (*Claims, error)`
- [ ] 5.2 `internal/services/auth/password.go`:
  - `HashPassword(plain string) (string, error)` — `bcrypt.GenerateFromPassword([]byte(plain), 12)`; result: `$2a$12$...`
  - `CheckPassword(plain, hash string) error` — `bcrypt.CompareHashAndPassword`
  - `IsLegacyMD5(hash string) bool` — `!strings.HasPrefix(hash, "$2")` (32-char hex MD5 has no `$2` prefix)
  - `VerifyMD5(plain, hash string) bool` — `fmt.Sprintf("%x", md5.Sum([]byte(plain))) == hash` using stdlib `crypto/md5`
  - No third-party password library — only `golang.org/x/crypto/bcrypt` + stdlib `crypto/md5`
- [ ] 5.3 `internal/api/middleware/auth.go` — Echo JWT middleware: extract Bearer from `Authorization` header; validate with `auth.ValidateToken`; inject `*Claims` into context key `"claims"`; return `{"error":"unauthorized"}` on failure
- [ ] 5.4 `internal/api/handlers/auth.go` — handlers:
  - `POST /api/v1/auth/login`: load user from DB; if `IsLegacyMD5(user.Password)`: verify with `VerifyMD5`, on match return `{"requires_password_reset":true}` with HTTP 200 and **no token**; if bcrypt: `CheckPassword`, on success issue access+refresh tokens
  - `POST /api/v1/auth/refresh`: validate httpOnly cookie, issue new access token
  - `POST /api/v1/auth/logout`: clear cookie (Max-Age=0)
  - `GET /api/v1/auth/me`: return claims from context (no DB query)
  - `POST /api/v1/auth/reset-password`: load user, verify old password (supports MD5 or bcrypt), store new `HashPassword`, issue tokens — this is the upgrade path from MD5 to bcrypt
- [ ] 5.5 `internal/api/routes.go` — register auth routes (public); apply JWT middleware to all other route groups; mount Swagger UI at `/api/swagger/`
- [ ] 5.6 Write unit tests `internal/services/auth/jwt_test.go` and `password_test.go` — round-trip, expiry, MD5 detection, VerifyMD5 match/mismatch, bcrypt verify/mismatch (table-driven)
- [ ] 5.7 Write handler test `internal/api/handlers/auth_test.go` — login with bcrypt (200+token), login with MD5 correct (200+requires_password_reset), wrong password (401), refresh valid cookie (200), refresh missing cookie (401), logout (200+clear cookie), me (200+claims)

## 6. REST API — Nagios Object Handlers

- [ ] 4.1 `internal/api/handlers/hosts.go` — `GET/POST /api/v1/hosts`, `GET/PUT/DELETE /api/v1/hosts/:id`; domain scoping via `config_id`; pagination (`?page=&limit=`); sort (`?sort=host_name&dir=asc`); dependency-check before delete; write `tbl_logbook` on all writes; swaggo annotations on every handler
- [ ] 4.2 `internal/api/handlers/services.go` — same pattern; filter by `?host_id=`; group listing by `config_name`
- [ ] 4.3 `internal/api/handlers/commands.go` — split listing by `?type=check|notify`; full CRUD
- [ ] 4.4 `internal/api/handlers/timeperiods.go` — inline `tbl_timedefinition` in create/update body (`ranges: [{definition:"monday", range:"00:00-24:00"}]`)
- [ ] 4.5 `internal/api/handlers/contacts.go` — all address fields, service/host notification command FKs
- [ ] 4.6 `internal/api/handlers/groups.go` — contact groups, host groups, service groups + `PUT /api/v1/hostgroups/:id/members` for member management
- [ ] 4.7 `internal/api/handlers/templates.go` — host, service, contact templates; filter `?register=0`
- [ ] 4.8 `internal/api/handlers/advanced.go` — dependencies, escalations, extinfo; SHA1 hash on create for duplicate detection; return 409 on duplicate
- [ ] 4.9 `internal/api/handlers/domains.go` — data domains (id=0 non-deletable) and config targets (method 1/2/3)
- [ ] 4.10 `internal/api/handlers/users.go` — user CRUD (admin-only); self-deletion 409; group membership via `PUT /api/v1/groups/:id/users`
- [ ] 4.11 `internal/api/handlers/settings.go` — `GET/PUT /api/v1/settings`; reads/writes `tbl_settings` key-value pairs
- [ ] 4.12 `internal/api/handlers/logbook.go` — `GET /api/v1/logbook?from=&to=&user=`; read-only
- [ ] 4.13 `internal/services/logbook/logbook.go` — `WriteLog(db, userID, action, objectType, objectName, source)` called by all write handlers
- [ ] 4.14 `internal/api/handlers/variables.go` — `GET/POST/DELETE /api/v1/hosts/:id/variables` and `GET/POST/DELETE /api/v1/services/:id/variables` for `_VARNAME` custom variables
- [ ] 4.15 `internal/api/handlers/specials.go` — `GET/PUT /api/v1/specials` for `$USERn$` resource macros; writes `resource.cfg` on PUT
- [ ] 4.16 `internal/api/handlers/monitoring.go` — `GET /api/v1/monitoring/summary` — active/inactive counts per object type for the dashboard

## 7. Config Generation Engine

- [ ] 5.1 `internal/services/nagconfig/generator.go` — entry point: `NewGenerator(db, target)` returning a `Generator` with methods `WriteHost`, `WriteService`, `WriteAll`
- [ ] 5.2 Create Go `text/template` files embedded in the binary (in `internal/services/nagconfig/templates/`): `host.tpl`, `service.tpl`, `command.tpl`, `timeperiod.tpl`, `contact.tpl`, `contactgroup.tpl`, `hostgroup.tpl`, `servicegroup.tpl`, `hosttemplate_v3.tpl`, `hosttemplate_v4.tpl`, `servicetemplate_v3.tpl`, `servicetemplate_v4.tpl`, `contacttemplate.tpl`, `hostdependency.tpl`, `servicedependency.tpl`, `hostescalation.tpl`, `serviceescalation.tpl`, `hostextinfo.tpl`, `serviceextinfo.tpl`
- [ ] 5.3 `WriteHost(hostID int) error` — query host + its custom variables + templates; render template; backup existing file; write to `{hostconfig}/{host_name}.cfg`
- [ ] 5.4 `WriteServiceGroup(configName string) error` — query all services with that `config_name`; render; write to `{serviceconfig}/{config_name}.cfg`
- [ ] 5.5 `WriteAll(objectType string) error` — writes a single file per global type (commands, timeperiods, contacts, etc.)
- [ ] 5.6 Implement `wrapLongLine(line string, limit int) string` — splits at `limit` chars with ` \` continuation; unit-testable
- [ ] 5.7 Implement backup-before-overwrite: copy existing `.cfg` to `{hostbackup}/{filename}.{timestamp}.bak`
- [ ] 5.8 `IsStale(tableName string, filePath string) (bool, error)` — compare `tbl_tablestatus.last_modified` vs `os.Stat(filePath).ModTime()`
- [ ] 5.9 `internal/services/nagconfig/target_local.go` — local filesystem write (`os.WriteFile`)
- [ ] 5.10 `internal/services/nagconfig/target_ftp.go` — FTP upload via `jlaffaye/ftp`; FTPS optional
- [ ] 5.11 `internal/services/nagconfig/target_ssh.go` — SFTP via `golang.org/x/crypto/ssh`; key or password auth
- [ ] 5.12 `POST /api/v1/config/write` handler — write all; return `{written:N, errors:[]}` summary
- [ ] 5.13 `POST /api/v1/config/write/:type/:id` handler — write single object; return same summary shape
- [ ] 5.14 Write unit tests `nagconfig/*_test.go` — host template rendering, long-line wrap, stale detection (mock filesystem), custom variable injection, v3 vs v4 template selection

## 8. Nagios Control Endpoints

- [ ] 6.1 `internal/services/nagios/control.go` — `Verify(confFile string, target ConfigTarget) (output string, valid bool, err error)`: runs `{target.binaryfile} -v {target.conffile}` locally or via SSH; captures stdout+stderr; `valid = true` when output contains "Total Errors:   0"
- [ ] 6.2 `Restart(target ConfigTarget) error` — runs `Verify` first; on success calls `TriggerReload(target)` which writes the current Unix timestamp to `{target.commandfile}` (= `reload.trigger`); `reload-watcher.sh` polls this file's mtime and issues a graceful reload to the Nagios daemon. **DO NOT write to `nagios.cmd` FIFO** — that is managed by Nagios Core exclusively and opening it from Go risks deadlock.
- [ ] 6.3 `POST /api/v1/config/verify` — return `{"valid":true/false,"output":"..."}` (no restart)
- [ ] 6.4 `POST /api/v1/config/restart` — auto-verify first; if invalid return `{"valid":false,"output":"...","restarted":false}`; if valid call `Restart`, return `{"valid":true,"restarted":true}`
- [ ] 6.5 `GET /api/v1/config/nagios-cfg` + `PUT /api/v1/config/nagios-cfg` — read/write raw `nagios.cfg` content using `target.conffile` path (local or SSH); admin-only
- [ ] 6.6 `GET /api/v1/config/cgi-cfg` + `PUT /api/v1/config/cgi-cfg` — same for `target.cgifile`; admin-only

## 9. Import Engine

- [ ] 7.1 `internal/services/nagimport/parser.go` — `ParseFile(path string) ([]ParsedObject, error)`; detect `define X {` blocks; parse key-value pairs; skip unknown directives
- [ ] 7.2 Implement typed parsed structs for all 17 object types; `ParsedObject.Type` field identifies which struct to cast
- [ ] 7.3 `Importer.Import(parsed []ParsedObject, configID int, overwrite bool) ImportResult` — resolve relation names to DB IDs; upsert or skip; write `tbl_logbook`
- [ ] 7.4 `ImportResult` struct: `{Inserted int, Updated int, Skipped int, Failed []FailedObject}`
- [ ] 7.5 Remote file retrieval: `GetRemoteFile(target ConfigTarget, remotePath string) (localTempPath string, error)` — FTP or SFTP download to temp dir
- [ ] 7.6 `GET /api/v1/import/files` — list `.cfg` files in `tbl_configtarget.importdir`
- [ ] 7.7 `POST /api/v1/import` body: `{"file":"hosts.cfg","config_id":1,"overwrite":false}` — trigger import; return `ImportResult`
- [ ] 7.8 Write unit tests `nagimport/*_test.go` using fixture `.cfg` files in `test/fixtures/`

## 10. Swagger Documentation

- [ ] 8.1 Add `main.go` swaggo global annotations: `@title`, `@version`, `@description`, `@host`, `@BasePath /api/v1`, `@securityDefinitions.apikey BearerAuth`
- [ ] 8.2 Add swaggo doc comments to every handler in sections 3–6 above (summary, tags, param, success, failure, security)
- [ ] 8.3 Run `swag init -g main.go -o docs --parseDependency` — verify `docs/swagger.json` and `docs/swagger.yaml` are generated
- [ ] 8.4 Mount `swaggo/http-swagger` at `/api/swagger/` in `routes.go` — no JWT required for Swagger UI
- [ ] 8.5 Update `embed.go` to include `//go:embed docs`; verify Swagger UI loads from the binary without external files

## 11. Console API Tests (curl bash scripts)

Each script MUST: use `set -euo pipefail`; define `BASE_URL=${BASE_URL:-http://localhost:8081}`; use `jq` for JSON parsing; print each step with `echo "[STEP] ..."` before running; assert HTTP status codes with `curl -s -w "%{http_code}" -o /tmp/body`; print `[PASS]` or `[FAIL]` per step; clean up created objects in a `trap` handler.

- [ ] 9.1 `test/api/lib.sh` — shared helpers: `login()` (stores `$TOKEN` and `$REFRESH_COOKIE`), `api_get()`, `api_post()`, `api_put()`, `api_delete()`, `assert_status()`, `assert_field()`
- [ ] 9.2 `test/api/auth.sh` — steps: login (assert 200 + access_token), me (assert 200 + username), bad login (assert 401), refresh (assert 200 + new token), logout (assert 200 + clear cookie)
- [ ] 9.3 `test/api/hosts.sh` — steps: create host `web-test-01` (assert 201), list hosts (assert 200, assert `web-test-01` in response via jq), get by ID (assert 200), update alias (assert 200), delete (assert 200), get deleted (assert 404)
- [ ] 9.4 `test/api/services.sh` — steps: create service on `web-test-01`, list services (filter by host_id), update, delete; also test ping service with check_command `check_ping`
- [ ] 9.5 `test/api/commands.sh` — create check command (type=1), create notify command (type=2), list with `?type=check` (assert only check commands), list with `?type=notify`, delete both
- [ ] 9.6 `test/api/timeperiods.sh` — create 24x7 period with ranges, list, get, delete
- [ ] 9.7 `test/api/contacts.sh` — create contact with email, list, update, delete
- [ ] 9.8 `test/api/groups.sh` — create host group, assign host to group, verify member in GET response, remove member, delete group
- [ ] 9.9 `test/api/templates.sh` — create host template (`register=0`), list templates, delete
- [ ] 9.10 `test/api/variables.sh` — add `_SNMP_COMMUNITY=public` to host, list host variables, delete variable
- [ ] 9.11 `test/api/specials.sh` — get resource macros, set `$USER1$=/usr/local/nagios/libexec`, verify response
- [ ] 9.12 `test/api/config.sh` — steps: write all configs (assert 200 + written>0), verify (assert `valid:true`), get nagios.cfg content (assert 200 + non-empty), attempt restart (assert restarted:true OR valid:false with output)
- [ ] 9.13 `test/api/import.sh` — list import files, import fixture host file, assert host appears in `GET /api/v1/hosts`
- [ ] 9.14 `test/api/users.sh` — (admin token) create user, list users, update, change password, delete; self-delete (assert 409)
- [ ] 9.15 `test/api/logbook.sh` — perform some writes, query logbook, assert entries appear with correct action/user
- [ ] 9.16 `test/api/smoke.sh` — `source lib.sh`; run all scripts in dependency order; each script exits non-zero on failure which propagates; print summary table at end
- [ ] 9.17 Add Makefile targets:
  - `server-start`: build and start `nagiosql serve` in background, wait until `curl -sf http://localhost:8081/api/swagger/index.html` succeeds, write PID to `.server.pid`
  - `server-stop`: kill PID from `.server.pid`
  - `test-api`: `server-start` then `bash test/api/smoke.sh` then `server-stop`
- [ ] 9.18 `test/api/README.md` — document: prerequisites (jq, curl, running MariaDB via `make db-start`), how to start server (`make server-start`), how to run all scripts (`make test-api`), how to run single script (`TOKEN=$(bash test/api/lib.sh login) bash test/api/hosts.sh`)

## 12. Unit and Handler Tests

- [ ] 10.1 `internal/services/auth/jwt_test.go` — table-driven: GenerateAccessToken/ValidateToken round-trip; expired token; wrong secret; IsLegacyMD5 various inputs
- [ ] 10.2 `internal/services/auth/password_test.go` — HashPassword/CheckPassword round-trip; CheckPassword wrong password returns error; IsLegacyMD5 true for MD5, false for bcrypt
- [ ] 10.3 `internal/services/nagconfig/template_test.go` — host output format; service output; long-line wrap at 800 chars; custom variable injection; v3 vs v4 importance directive
- [ ] 10.4 `internal/services/nagimport/parser_test.go` — parse single host block; parse multi-block file; skip unknown directive; parse all 17 object types using fixture files
- [ ] 10.5 `internal/api/handlers/auth_test.go` — login success 200 + tokens; wrong password 401; MD5 user returns requires_password_reset; refresh 200; refresh no cookie 401; logout 200 clears cookie; me 200
- [ ] 10.6 `internal/api/handlers/hosts_test.go` — list 200; create 201; get by ID 200; get 404; update 200; delete 200; delete with dependents 409; unauthenticated 401
- [ ] 10.7 `internal/api/handlers/services_test.go`, `commands_test.go`, `contacts_test.go`, `groups_test.go`, `templates_test.go` — 200/201/404/409 for each CRUD op
- [ ] 10.8 `internal/api/handlers/config_test.go` — mock `nagios.Verifier` interface; verify returns valid:true/false; restart calls verify then write pipe; restart blocked when verify fails
- [ ] 10.9 `internal/api/middleware/auth_test.go` — valid token passes; expired token 401; malformed token 401; missing header 401
- [ ] 10.10 `internal/config/config_test.go` — load from TOML file; env var override (`NAGIOSQL_JWT_SECRET`); missing required field returns error

## 13. Test Infrastructure (Docker + MariaDB)

The existing `DOCUMENTS/docker/nagios-core/docker-compose.yml` already provides MariaDB 10.11. We create a minimal port-override file so the test DB is reachable from the host without touching `DOCUMENTS/`.

- [ ] 11.1 Create `docker/test/docker-compose.db.yml` — a **compose override** file with only the `db` service redeclared with host port mapping: `ports: ["3307:3306"]`, same env vars as `DOCUMENTS/docker/nagios-core/docker-compose.yml`; uses `MYSQL_DATABASE=nagiosql_test`; named volume `mariadb-test-data`. Image: `mariadb:10.11`. Healthcheck: `mysqladmin ping -h localhost`.
- [ ] 13.2 The Docker test DB starts **empty** (only the `nagiosql_test` database and user created via MariaDB env vars); no SQL init files are mounted — schema and seed data are applied by Go code in `TestMain` via `db.Migrate()` and `seeds.SeedRequired()` + `seeds.SeedSample()`
- [ ] 11.5 Create `docker/test/Dockerfile.ci` — `FROM golang:1.26-bookworm`; install `mariadb-client curl jq upx`; `WORKDIR /workspace`; `COPY . .`; `RUN go mod download`; `CMD ["make", "ci-inner"]` (runs vet + test + test-integration inside container)
- [ ] 11.6 Create `docker/test/docker-compose.ci.yml` — two services: `mariadb` (image: mariadb:10.11, same as 11.1 but on default port 3306 internal) + `nagiosql-ci` (Dockerfile.ci, `depends_on: mariadb: condition: service_healthy`, env `TEST_DSN=nagiosql:test@tcp(mariadb:3306)/nagiosql_test?parseTime=true`)
- [ ] 11.7 Add Makefile targets:
  - `db-start`: `docker compose -f docker/test/docker-compose.db.yml up -d --wait && echo "MariaDB ready on :3307"`
  - `db-stop`: `docker compose -f docker/test/docker-compose.db.yml down -v`
  - `db-reset`: `$(MAKE) db-stop && $(MAKE) db-start`
  - `test-integration`: `$(MAKE) db-start && TEST_DSN="nagiosql:test@tcp(127.0.0.1:3307)/nagiosql_test?parseTime=true" go test ./internal/integration/... -tags integration -v -count=1 -timeout 120s; CODE=$$?; $(MAKE) db-stop; exit $$CODE`
  - `ci`: `docker compose -f docker/test/docker-compose.ci.yml up --abort-on-container-exit --exit-code-from nagiosql-ci`
  - `ci-inner`: `$(MAKE) vet && $(MAKE) test && $(MAKE) test-integration`
- [ ] 13.8 Create `internal/integration/setup_test.go` — `//go:build integration`; `TestMain` reads `TEST_DSN` env (default `nagiosql:test@tcp(127.0.0.1:3307)/nagiosql_test?parseTime=true`); opens GORM; calls `db.Migrate()` (creates all tables); calls `seeds.SeedRequired(db, testCfg, "admin", "admin")` and `seeds.SeedSample(db)`; defers cleanup that calls `db.Exec("SET FOREIGN_KEY_CHECKS=0")` + truncates all `tbl_*` tables + `SET FOREIGN_KEY_CHECKS=1`; stores `testDB` in package var for all integration tests
- [ ] 13.9 `internal/integration/auth_test.go` — `//go:build integration`; login with `admin/admin` (seeded via `SeedAdmin` with bcrypt — NOT MD5); assert 200 + `access_token`; call `/me`; refresh; logout; also test: insert a raw MD5 user directly via `db.Exec`, then login and assert `requires_password_reset: true`
- [ ] 13.10 `internal/integration/hosts_test.go` — `//go:build integration`; assert seeded hosts from `SeedSample` are present (GET /api/v1/hosts returns at least 4); create new host (assert 201), update alias (assert 200), delete (assert 200), get deleted (assert 404)
- [ ] 13.11 `internal/integration/config_test.go` — `//go:build integration`; update `tbl_configtarget.hostconfig` to `t.TempDir()`; call `GenerateHost` for seeded `localhost` host (id=4); assert file written at `{tempdir}/localhost.cfg`; parse with `nagimport.ParseFile`; assert `host_name="localhost"` and `address="127.0.0.1"`

## 14. Docker Integration and Final Validation

The existing Docker stack lives in `DOCUMENTS/docker/nagios-core/`. Our Go binary slots in as a drop-in replacement for the PHP+nginx layer — we provide an updated Dockerfile and supervisord.conf alongside the existing compose, **without modifying any files under DOCUMENTS/**.

- [ ] 12.1 Create `go-nagiosql/Dockerfile` — multi-stage: stage 1 `FROM golang:1.26-bookworm AS builder` builds the binary; stage 2 `FROM debian:trixie-slim` installs `ca-certificates curl jq` and copies binary. Exposes port 8081. Entrypoint: `["/app/nagiosql"]`, CMD: `["serve"]`
- [ ] 12.2 Create `go-nagiosql/docker-compose.override.yml` — overrides the `nagios` service from `DOCUMENTS/docker/nagios-core/docker-compose.yml` to: build from `go-nagiosql/Dockerfile`, set `NAGIOSQL_JWT_SECRET`, `DB_HOST=db`, `DB_NAME=nagiosql`, `DB_USER=nagiosql`, `DB_PASSWORD`; remove the old PHP volumes (`/var/www/nagiosql`); keep the shared nagios volumes (`/usr/local/nagios/etc`, `/usr/local/nagios/var`); run `nagiosql migrate && nagiosql serve` as command
- [ ] 12.3 Create `go-nagiosql/supervisord.conf.example` — shows how to remove the `php-fpm`, `fcgiwrap`, and NagiosQL nginx block; add `[program:nagiosql]` with `command=/app/nagiosql serve`, `autorestart=true`, `priority=20`
- [ ] 12.4 Create `go-nagiosql/config.toml.docker` — ready-to-mount config for Docker; all `[nagios]` paths match the paths that `entrypoint.sh` sets in `tbl_configtarget`:
  - `reload_trigger = "/usr/local/nagios/var/reload.trigger"`
  - `host_config_dir = "/usr/local/nagios/etc/nagiosql/hosts/"`
  - `service_config_dir = "/usr/local/nagios/etc/nagiosql/services/"`
  - `binary = "/usr/local/nagios/bin/nagios"`
  - `config_file = "/usr/local/nagios/etc/nagios.cfg"`
  - `resource_file = "/usr/local/nagios/etc/resource.cfg"`
- [ ] 12.5 Start the full stack: `cd DOCUMENTS/docker/nagios-core && cp .env.example .env && docker compose up -d --wait`; verify MariaDB and Nagios Core are healthy; run `make test-api BASE_URL=http://localhost:8081` — all smoke tests pass
- [ ] 12.6 Import real `.cfg` fixtures: `nagiosql import test/fixtures/cfg/linux-host.cfg`; verify `linux-host` appears in `GET /api/v1/hosts`; then `nagiosql import test/fixtures/cfg/ping.cfg`; verify service appears in `GET /api/v1/services`
- [ ] 12.7 Write all configs: `POST /api/v1/config/write` → assert `written>0`; `POST /api/v1/config/verify` → assert `valid:true`; inspect generated `/usr/local/nagios/etc/nagiosql/hosts/linux-host.cfg` (via volume) for correctness
- [ ] 12.8 Run `make vet` — zero findings; run `make lint` — zero findings; run `make test-cover` — assert ≥70% on `internal/services/` and `internal/api/handlers/`
- [ ] 12.9 Run `make all` — verify binary size <50MB; run `make upx` — verify compressed binary <20MB
- [ ] 12.10 Verify Swagger UI at `http://localhost:8081/api/swagger/` documents all `/api/v1/` endpoints with `BearerAuth` security applied to protected routes
