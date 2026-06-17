## ADDED Requirements

### Requirement: Host CRUD
The system SHALL provide full create, read, update, and delete operations for Nagios hosts via `/admin/hosts` (UI) and `/api/v1/hosts` (API). All fields from `tbl_host` SHALL be supported including check_command with arguments, notification options, contact assignments, and template inheritance.

#### Scenario: Create host
- **WHEN** a user submits a valid host form with at minimum `host_name` and `address`
- **THEN** the host is saved to `tbl_host` with `active='1'` and the audit log is updated in `tbl_logbook`

#### Scenario: List hosts with pagination
- **WHEN** a user accesses `/admin/hosts`
- **THEN** hosts are listed with configurable records-per-page, sortable by name or alias

#### Scenario: Delete host with dependency check
- **WHEN** a user attempts to delete a host that is referenced by services or host groups
- **THEN** deletion is blocked and a message lists the dependent objects

#### Scenario: Copy host
- **WHEN** a user clicks "Copy" on an existing host
- **THEN** a new host is created with the same settings and a `_copy` suffix on the name

### Requirement: Service CRUD
The system SHALL provide full CRUD for Nagios services via `/admin/services` and `/api/v1/services`. Services SHALL be groupable by `config_name` for per-file generation. All fields from `tbl_service` SHALL be supported.

#### Scenario: Create service
- **WHEN** a user creates a service with `service_description` and `check_command`
- **THEN** the service is saved to `tbl_service` linked to the selected host

#### Scenario: Filter services by host
- **WHEN** a user filters the service list by a specific host
- **THEN** only services assigned to that host are shown

### Requirement: Command CRUD
The system SHALL provide CRUD for Nagios commands via `/admin/checkcommands` (check commands) and `/admin/misccommands` (notify/event commands), backed by `tbl_command`. Command type (1=check, 2=notify) SHALL be used to separate the two pages.

#### Scenario: Create check command
- **WHEN** a user creates a command with type=1, name, and command line
- **THEN** it is saved to `tbl_command` and available in host/service command selectors

### Requirement: Time period CRUD
The system SHALL provide CRUD for Nagios time periods via `/admin/timeperiods`, backed by `tbl_timeperiod` and `tbl_timedefinition`. Time ranges SHALL be editable inline.

#### Scenario: Create time period with ranges
- **WHEN** a user creates a time period and adds day definitions (e.g., `monday 00:00-24:00`)
- **THEN** the period and its ranges are saved to `tbl_timeperiod` and `tbl_timedefinition` respectively

### Requirement: Contact CRUD
The system SHALL provide CRUD for Nagios contacts via `/admin/contacts`, backed by `tbl_contact`. All address fields (`address1`–`address6`) and notification command assignments SHALL be supported.

#### Scenario: Create contact
- **WHEN** a user creates a contact with name, email, and notification commands
- **THEN** the contact is saved to `tbl_contact`

### Requirement: Group object CRUD (contact, host, service groups)
The system SHALL provide CRUD for contact groups (`/admin/contactgroups`), host groups (`/admin/hostgroups`), and service groups (`/admin/servicegroups`), backed by their respective `tbl_*` tables and `tbl_lnk*` link tables.

#### Scenario: Assign hosts to host group
- **WHEN** a user edits a host group and selects member hosts
- **THEN** entries are written to `tbl_lnkHostToHostgroup`

### Requirement: Template CRUD (host, service, contact)
The system SHALL provide CRUD for host templates (`/admin/hosttemplates`), service templates (`/admin/servicetemplates`), and contact templates (`/admin/contacttemplates`), each backed by their respective `tbl_*template` tables with `register='0'`.

#### Scenario: Template is excluded from config unless used
- **WHEN** a host template exists but no host uses it
- **THEN** the template name does NOT appear as a standalone host definition in generated `.cfg` files

### Requirement: Advanced objects (dependencies, escalations, extinfo)
The system SHALL provide CRUD for host dependencies (`/admin/hostdependencies`), service dependencies (`/admin/servicedependencies`), host escalations (`/admin/hostescalations`), service escalations (`/admin/serviceescalations`), host extinfo (`/admin/hostextinfo`), and service extinfo (`/admin/serviceextinfo`). Duplicate detection via SHA1 hash SHALL be implemented as in the original `NagDataClass::updateHash()`.

#### Scenario: Duplicate dependency detection
- **WHEN** a user attempts to create a host dependency identical to an existing one
- **THEN** the system detects the duplicate via hash comparison and rejects the creation with an error message

### Requirement: Data domain management
The system SHALL provide CRUD for data domains (`/admin/datadomain`) and config targets (`/admin/configtargets`). Domain id=0 ("common") SHALL always exist and be non-deletable. Objects in domain 0 SHALL appear in all domains when `enable_common=1`.

#### Scenario: Common domain objects appear in all domains
- **WHEN** a user switches to any data domain with `enable_common=1`
- **THEN** objects from domain id=0 are included in queries alongside domain-specific objects

### Requirement: Custom object variables on hosts and services
The system SHALL support Nagios custom object variables (`_VARNAME` format) on hosts, services, host templates, and service templates. These are stored in `tbl_variabledefinition` and linked via `tbl_lnkHostToVariabledefinition`, `tbl_lnkServiceToVariabledefinition`, `tbl_lnkHosttemplateToVariabledefinition`, and `tbl_lnkServicetemplateToVariabledefinition`. The host and service forms SHALL provide a dynamic variable editor (add/remove rows of name=value pairs). Custom variables SHALL appear in generated `.cfg` files as `_VARNAME value` inside the object definition block.

#### Scenario: Add custom variable to host
- **WHEN** a user adds `_COMMUNITY snmp-community` to a host via the form
- **THEN** the variable is saved to `tbl_variabledefinition` and linked to the host in `tbl_lnkHostToVariabledefinition`

#### Scenario: Custom variable in generated config
- **WHEN** a host with custom variable `_SNMP_PORT 161` is written to a `.cfg` file
- **THEN** the file contains `_SNMP_PORT 161` inside the `define host { }` block

### Requirement: Resource macros ($USERn$ variables)
The system SHALL support editing Nagios resource macros (`$USER1$`–`$USERn$`) via `/admin/specials`. These map to `tbl_variabledefinition` records with `varType='nagiosresource'` and are written to `/usr/local/nagios/etc/resource.cfg` during config generation.

#### Scenario: Edit resource macro
- **WHEN** a user sets `$USER1$` to `/usr/local/nagios/libexec` via `/admin/specials`
- **THEN** the value is saved to `tbl_variabledefinition` with `varType='nagiosresource'`

#### Scenario: Resource.cfg written
- **WHEN** config generation is triggered and resource macros are defined
- **THEN** `/usr/local/nagios/etc/resource.cfg` is written with `$USERn$=value` entries

### Requirement: Nagios 4 importance directive
The system SHALL support the `importance` directive on hosts and services (Nagios 4-only) and `minimum_importance` on contacts, controlled by the Nagios version setting in `tbl_configtarget`. These fields SHALL only appear in generated `.cfg` files when version=4.

#### Scenario: Importance field in version 4 mode
- **WHEN** a config target has `version=4` and a host has `importance=10`
- **THEN** `importance 10` appears in the generated host `.cfg`

#### Scenario: Importance field absent in version 3 mode
- **WHEN** a config target has `version=3`
- **THEN** `importance` is omitted from generated `.cfg` files

### Requirement: Audit log
All create, update, delete, and config-write operations SHALL be logged to `tbl_logbook` with timestamp, username, action description, and interface (web/scripting). The log SHALL be viewable at `/admin/logbook` with date filtering.

#### Scenario: Host creation logged
- **WHEN** a user creates a host
- **THEN** a row is written to `tbl_logbook` with action=`add`, the username, and the host name
