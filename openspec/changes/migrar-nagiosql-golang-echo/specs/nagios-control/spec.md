## ADDED Requirements

### Requirement: Write all config files
The system SHALL provide an action at `/admin/verify` that triggers generation of all config files for all active objects in the current data domain. This SHALL be equivalent to NagiosQL's "Write monitoring data" + "Write additional data" buttons.

#### Scenario: Write all configs
- **WHEN** a user clicks "Write Config" on the verify page
- **THEN** config files are generated for all active hosts, services, commands, time periods, contacts, groups, templates, dependencies, escalations, and extinfo; per-entity files are written; a success/failure summary is shown

#### Scenario: Write single object config
- **WHEN** a user saves a host and clicks "Generate Config"
- **THEN** only that host's `.cfg` file is regenerated without affecting other files

### Requirement: Config validation (nagios -v)
The system SHALL run `nagios -v {conf_file}` and stream the output to the user interface. This is equivalent to NagiosQL's "Check configuration" button. The check SHALL be run locally or via SSH depending on the config target.

#### Scenario: Valid config
- **WHEN** `nagios -v` is run and returns exit code 0
- **THEN** the output is shown in the UI with a green success indicator

#### Scenario: Invalid config
- **WHEN** `nagios -v` returns a non-zero exit code
- **THEN** the full error output is displayed in the UI with a red failure indicator and no reload is performed

### Requirement: Nagios Core restart/reload
The system SHALL send a `RESTART_PROGRAM` command to the Nagios command pipe (`nagios.cmd`) to trigger a reload. Before writing the command, `nagios -v` SHALL run automatically and block the restart if validation fails. A reload trigger file SHALL also be written for the `reload-watcher.sh` compatibility layer.

#### Scenario: Restart after successful validation
- **WHEN** a user clicks "Restart Nagios" and `nagios -v` succeeds
- **THEN** `[timestamp] RESTART_PROGRAM` is written to the command pipe and the reload trigger file is written

#### Scenario: Restart blocked on invalid config
- **WHEN** a user clicks "Restart Nagios" and `nagios -v` returns errors
- **THEN** the restart is blocked, the validation errors are shown, and no command is written to the pipe

#### Scenario: Restart via SSH (remote target)
- **WHEN** the config target uses method=3 (SSH)
- **THEN** the restart command is executed on the remote host via SSH using configured credentials

### Requirement: NagiosQL settings editor
The system SHALL provide an editor for `nagios.cfg` and `cgi.cfg` at `/admin/nagioscfg`, equivalent to NagiosQL's `nagioscfg.php`. The files SHALL be read and written locally or via SSH depending on the config target.

#### Scenario: Edit nagios.cfg
- **WHEN** a user edits and saves `nagios.cfg` via the UI
- **THEN** the file is written to the path from `tbl_configtarget.conffile`

### Requirement: Verify page summary dashboard
The `/admin/verify` page SHALL display a summary table of all config objects with last-generated timestamps, stale status indicators, and per-object regeneration buttons. This matches the NagiosQL `verify.php` layout.

#### Scenario: Verify page shows stale objects
- **WHEN** the verify page is loaded and some host configs are outdated
- **THEN** those hosts are highlighted with a warning indicator and an individual "Regenerate" button
