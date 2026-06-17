## ADDED Requirements

### Requirement: Import Nagios .cfg files into database
The system SHALL parse existing Nagios `.cfg` files from the configured import directory and insert the parsed objects into the MariaDB database. This replaces `NagImportClass.php`. The import SHALL be accessible via `/admin/import`.

#### Scenario: List importable files
- **WHEN** a user accesses `/admin/import`
- **THEN** the system lists all `.cfg` files in the `tbl_configtarget.importdir` directory

#### Scenario: Import a host config file
- **WHEN** a user selects a `.cfg` file containing `define host { }` blocks and clicks "Import"
- **THEN** each host definition is parsed and inserted into `tbl_host`; relation fields (templates, contacts, groups) are resolved to database IDs

#### Scenario: Import with overwrite
- **WHEN** a user imports a file with `overwrite=1` and a host with the same name already exists
- **THEN** the existing record is updated with values from the imported file

#### Scenario: Import without overwrite
- **WHEN** a user imports a file with `overwrite=0` and an object already exists
- **THEN** the existing record is kept and a skip message is shown for that object

### Requirement: Parser for all Nagios object types
The import engine SHALL parse all Nagios object types that NagiosQL manages: `host`, `service`, `command`, `timeperiod`, `contact`, `contactgroup`, `hostgroup`, `servicegroup`, `hosttemplate`, `servicetemplate`, `contacttemplate`, `hostdependency`, `servicedependency`, `hostescalation`, `serviceescalation`, `hostextinfo`, `serviceextinfo`.

#### Scenario: Parse multi-object file
- **WHEN** a `.cfg` file contains multiple `define X { }` blocks of different types
- **WHEN** the file is imported
- **THEN** each block is parsed and inserted into the correct `tbl_*` table

#### Scenario: Unrecognized directive skipped
- **WHEN** a `.cfg` file contains a directive not in NagiosQL's schema
- **THEN** that directive is silently skipped and the rest of the object is imported

### Requirement: Remote import via FTP and SSH
The import engine SHALL retrieve `.cfg` files from FTP or SSH sources when the config target method is 2 or 3, copying files to the local import directory before parsing.

#### Scenario: Remote import via SSH
- **WHEN** the config target uses method=3 and a user triggers an import
- **THEN** the system connects via SFTP, downloads the selected file, and processes it locally

### Requirement: Import result reporting
After import, the system SHALL display a report listing the number of objects imported, skipped, and failed, along with error details for any parse failures.

#### Scenario: Import report displayed
- **WHEN** import of a file completes
- **THEN** the UI shows counts of inserted, updated, skipped, and failed records with object names for failures
