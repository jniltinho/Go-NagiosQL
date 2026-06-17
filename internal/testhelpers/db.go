// Package testhelpers provides shared test utilities.
// It must NOT be imported outside _test.go files or test packages.
package testhelpers

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB opens a private in-memory SQLite database for one test.
// MaxOpenConns(1) is critical: with `:memory:` each new connection creates
// a separate database, so we force GORM to reuse a single connection.
func NewDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	for _, sql := range schema() {
		if err := db.Exec(sql).Error; err != nil {
			t.Fatalf("create table: %v\nSQL: %.80s...", err, sql)
		}
	}
	return db
}

func schema() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS tbl_user (
			id             INTEGER PRIMARY KEY AUTOINCREMENT,
			username       TEXT NOT NULL UNIQUE,
			password       TEXT NOT NULL,
			name           TEXT NOT NULL DEFAULT '',
			email          TEXT NOT NULL DEFAULT '',
			admin          TEXT NOT NULL DEFAULT '0',
			active         TEXT NOT NULL DEFAULT '1',
			logon_timeout  INTEGER NOT NULL DEFAULT 60,
			domain_id      INTEGER NOT NULL DEFAULT 0,
			last_modified  DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_datadomain (
			id          INTEGER PRIMARY KEY,
			name        TEXT NOT NULL DEFAULT '',
			description TEXT NOT NULL DEFAULT '',
			active      TEXT NOT NULL DEFAULT '1'
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_configtarget (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			name          TEXT NOT NULL DEFAULT '',
			domain_id     INTEGER NOT NULL DEFAULT 0,
			host_path     TEXT NOT NULL DEFAULT '',
			service_path  TEXT NOT NULL DEFAULT '',
			backup_path   TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_settings (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			install_hash  TEXT NOT NULL DEFAULT '',
			nagios_bin    TEXT NOT NULL DEFAULT '',
			write_to_files TEXT NOT NULL DEFAULT '1'
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_logbook (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			domain_id   INTEGER NOT NULL DEFAULT 0,
			username    TEXT NOT NULL DEFAULT '',
			object_type TEXT NOT NULL DEFAULT '',
			action      TEXT NOT NULL DEFAULT '',
			object_name TEXT NOT NULL DEFAULT '',
			details     TEXT NOT NULL DEFAULT '',
			created_at  DATETIME
		)`,
		// Full tbl_host schema mirroring models.Host struct.
		`CREATE TABLE IF NOT EXISTS tbl_host (
			id                             INTEGER PRIMARY KEY AUTOINCREMENT,
			host_name                      TEXT NOT NULL DEFAULT '',
			alias                          TEXT NOT NULL DEFAULT '',
			display_name                   TEXT NOT NULL DEFAULT '',
			address                        TEXT NOT NULL DEFAULT '',
			parents                        INTEGER NOT NULL DEFAULT 0,
			parents_tploptions             INTEGER NOT NULL DEFAULT 2,
			importance                     INTEGER,
			hostgroups                     INTEGER NOT NULL DEFAULT 0,
			hostgroups_tploptions          INTEGER NOT NULL DEFAULT 2,
			check_command                  TEXT NOT NULL DEFAULT '',
			use_template                   INTEGER NOT NULL DEFAULT 0,
			use_template_tploptions        INTEGER NOT NULL DEFAULT 2,
			initial_state                  TEXT NOT NULL DEFAULT '',
			max_check_attempts             INTEGER,
			check_interval                 INTEGER,
			retry_interval                 INTEGER,
			active_checks_enabled          INTEGER NOT NULL DEFAULT 2,
			passive_checks_enabled         INTEGER NOT NULL DEFAULT 2,
			check_period                   INTEGER NOT NULL DEFAULT 0,
			obsess_over_host               INTEGER NOT NULL DEFAULT 2,
			check_freshness               INTEGER NOT NULL DEFAULT 2,
			freshness_threshold            INTEGER,
			event_handler                  INTEGER NOT NULL DEFAULT 0,
			event_handler_enabled          INTEGER NOT NULL DEFAULT 2,
			low_flap_threshold             REAL,
			high_flap_threshold            REAL,
			flap_detection_enabled         INTEGER NOT NULL DEFAULT 2,
			flap_detection_options         TEXT NOT NULL DEFAULT '',
			process_perf_data              INTEGER NOT NULL DEFAULT 2,
			retain_status_information      INTEGER NOT NULL DEFAULT 2,
			retain_nonstatus_information   INTEGER NOT NULL DEFAULT 2,
			contacts                       INTEGER NOT NULL DEFAULT 0,
			contacts_tploptions            INTEGER NOT NULL DEFAULT 2,
			contact_groups                 INTEGER NOT NULL DEFAULT 0,
			contact_groups_tploptions      INTEGER NOT NULL DEFAULT 2,
			notification_interval          INTEGER,
			notification_period            INTEGER NOT NULL DEFAULT 0,
			first_notification_delay       INTEGER,
			notification_options           TEXT NOT NULL DEFAULT '',
			notifications_enabled          INTEGER NOT NULL DEFAULT 2,
			stalking_options               TEXT NOT NULL DEFAULT '',
			notes                          TEXT NOT NULL DEFAULT '',
			notes_url                      TEXT NOT NULL DEFAULT '',
			action_url                     TEXT NOT NULL DEFAULT '',
			icon_image                     TEXT NOT NULL DEFAULT '',
			icon_image_alt                 TEXT NOT NULL DEFAULT '',
			vrml_image                     TEXT NOT NULL DEFAULT '',
			statusmap_image                TEXT NOT NULL DEFAULT '',
			"2d_coords"                    TEXT NOT NULL DEFAULT '',
			"3d_coords"                    TEXT NOT NULL DEFAULT '',
			use_variables                  INTEGER NOT NULL DEFAULT 0,
			name                           TEXT NOT NULL DEFAULT '',
			register                       TEXT NOT NULL DEFAULT '1',
			active                         TEXT NOT NULL DEFAULT '1',
			last_modified                  DATETIME,
			access_group                   INTEGER NOT NULL DEFAULT 0,
			config_id                      INTEGER NOT NULL DEFAULT 0
		)`,
		// Full tbl_service schema mirroring models.Service struct.
		`CREATE TABLE IF NOT EXISTS tbl_service (
			id                              INTEGER PRIMARY KEY AUTOINCREMENT,
			config_name                     TEXT NOT NULL DEFAULT '',
			host_name                       INTEGER NOT NULL DEFAULT 0,
			host_name_tploptions            INTEGER NOT NULL DEFAULT 2,
			hostgroup_name                  INTEGER NOT NULL DEFAULT 0,
			hostgroup_name_tploptions       INTEGER NOT NULL DEFAULT 2,
			service_description             TEXT NOT NULL DEFAULT '',
			display_name                    TEXT NOT NULL DEFAULT '',
			servicegroups                   INTEGER NOT NULL DEFAULT 0,
			servicegroups_tploptions        INTEGER NOT NULL DEFAULT 2,
			use_template                    INTEGER NOT NULL DEFAULT 0,
			use_template_tploptions         INTEGER NOT NULL DEFAULT 2,
			check_command                   TEXT NOT NULL DEFAULT '',
			is_volatile                     INTEGER NOT NULL DEFAULT 2,
			initial_state                   TEXT NOT NULL DEFAULT '',
			max_check_attempts              INTEGER,
			check_interval                  INTEGER,
			retry_interval                  INTEGER,
			active_checks_enabled           INTEGER NOT NULL DEFAULT 2,
			passive_checks_enabled          INTEGER NOT NULL DEFAULT 2,
			check_period                    INTEGER NOT NULL DEFAULT 0,
			parallelize_check               INTEGER NOT NULL DEFAULT 2,
			obsess_over_service             INTEGER NOT NULL DEFAULT 2,
			check_freshness                 INTEGER NOT NULL DEFAULT 2,
			freshness_threshold             INTEGER,
			event_handler                   INTEGER NOT NULL DEFAULT 0,
			event_handler_enabled           INTEGER NOT NULL DEFAULT 2,
			low_flap_threshold              REAL,
			high_flap_threshold             REAL,
			flap_detection_enabled          INTEGER NOT NULL DEFAULT 2,
			flap_detection_options          TEXT NOT NULL DEFAULT '',
			process_perf_data               INTEGER NOT NULL DEFAULT 2,
			retain_status_information       INTEGER NOT NULL DEFAULT 2,
			retain_nonstatus_information    INTEGER NOT NULL DEFAULT 2,
			notification_interval           INTEGER,
			first_notification_delay        INTEGER,
			notification_period             INTEGER NOT NULL DEFAULT 0,
			notification_options            TEXT NOT NULL DEFAULT '',
			notifications_enabled           INTEGER NOT NULL DEFAULT 2,
			contacts                        INTEGER NOT NULL DEFAULT 0,
			contacts_tploptions             INTEGER NOT NULL DEFAULT 2,
			contact_groups                  INTEGER NOT NULL DEFAULT 0,
			contact_groups_tploptions       INTEGER NOT NULL DEFAULT 2,
			stalking_options                TEXT NOT NULL DEFAULT '',
			notes                           TEXT NOT NULL DEFAULT '',
			notes_url                       TEXT NOT NULL DEFAULT '',
			action_url                      TEXT NOT NULL DEFAULT '',
			icon_image                      TEXT NOT NULL DEFAULT '',
			icon_image_alt                  TEXT NOT NULL DEFAULT '',
			use_variables                   INTEGER NOT NULL DEFAULT 0,
			name                            TEXT NOT NULL DEFAULT '',
			register                        TEXT NOT NULL DEFAULT '1',
			active                          TEXT NOT NULL DEFAULT '1',
			last_modified                   DATETIME,
			access_group                    INTEGER NOT NULL DEFAULT 0,
			config_id                       INTEGER NOT NULL DEFAULT 0,
			import_hash                     TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_command (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			command_name  TEXT NOT NULL DEFAULT '',
			command_line  TEXT NOT NULL DEFAULT '',
			command_type  INTEGER NOT NULL DEFAULT 0,
			arg1_info     TEXT,
			arg2_info     TEXT,
			arg3_info     TEXT,
			arg4_info     TEXT,
			arg5_info     TEXT,
			arg6_info     TEXT,
			arg7_info     TEXT,
			arg8_info     TEXT,
			register      TEXT NOT NULL DEFAULT '1',
			active        TEXT NOT NULL DEFAULT '1',
			last_modified DATETIME,
			access_group  INTEGER NOT NULL DEFAULT 0,
			config_id     INTEGER NOT NULL DEFAULT 0,
			UNIQUE(command_name, config_id)
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_timeperiod (
			id              INTEGER PRIMARY KEY AUTOINCREMENT,
			timeperiod_name TEXT NOT NULL DEFAULT '',
			alias           TEXT NOT NULL DEFAULT '',
			register        TEXT NOT NULL DEFAULT '1',
			active          TEXT NOT NULL DEFAULT '1',
			domain_id       INTEGER NOT NULL DEFAULT 0,
			config_id       INTEGER NOT NULL DEFAULT 0,
			last_modified   DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_timedefinition (
			id       INTEGER PRIMARY KEY AUTOINCREMENT,
			tip_id   INTEGER NOT NULL DEFAULT 0,
			day      TEXT NOT NULL DEFAULT '',
			time_def TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_contact (
			id                                INTEGER PRIMARY KEY AUTOINCREMENT,
			contact_name                      TEXT NOT NULL DEFAULT '',
			alias                             TEXT NOT NULL DEFAULT '',
			email                             TEXT NOT NULL DEFAULT '',
			pager                             TEXT NOT NULL DEFAULT '',
			address1                          TEXT NOT NULL DEFAULT '',
			address2                          TEXT NOT NULL DEFAULT '',
			address3                          TEXT NOT NULL DEFAULT '',
			address4                          TEXT NOT NULL DEFAULT '',
			address5                          TEXT NOT NULL DEFAULT '',
			address6                          TEXT NOT NULL DEFAULT '',
			host_notifications_enabled        TEXT NOT NULL DEFAULT '1',
			service_notifications_enabled     TEXT NOT NULL DEFAULT '1',
			host_notification_period          INTEGER NOT NULL DEFAULT 0,
			service_notification_period       INTEGER NOT NULL DEFAULT 0,
			host_notification_options         TEXT NOT NULL DEFAULT '',
			service_notification_options      TEXT NOT NULL DEFAULT '',
			host_notification_commands        INTEGER NOT NULL DEFAULT 0,
			service_notification_commands     INTEGER NOT NULL DEFAULT 0,
			use_variables                     INTEGER NOT NULL DEFAULT 0,
			register                          TEXT NOT NULL DEFAULT '1',
			active                            TEXT NOT NULL DEFAULT '1',
			domain_id                         INTEGER NOT NULL DEFAULT 0,
			config_id                         INTEGER NOT NULL DEFAULT 0,
			last_modified                     DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_contactgroup (
			id                    INTEGER PRIMARY KEY AUTOINCREMENT,
			contactgroup_name     TEXT NOT NULL DEFAULT '',
			alias                 TEXT NOT NULL DEFAULT '',
			members               INTEGER NOT NULL DEFAULT 0,
			contactgroup_members  INTEGER NOT NULL DEFAULT 0,
			register              TEXT NOT NULL DEFAULT '1',
			active                TEXT NOT NULL DEFAULT '1',
			access_group          INTEGER NOT NULL DEFAULT 0,
			domain_id             INTEGER NOT NULL DEFAULT 0,
			config_id             INTEGER NOT NULL DEFAULT 0,
			last_modified         DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_hostgroup (
			id                  INTEGER PRIMARY KEY AUTOINCREMENT,
			hostgroup_name      TEXT NOT NULL DEFAULT '',
			alias               TEXT NOT NULL DEFAULT '',
			members             INTEGER NOT NULL DEFAULT 0,
			hostgroup_members   INTEGER NOT NULL DEFAULT 0,
			notes               TEXT NOT NULL DEFAULT '',
			notes_url           TEXT NOT NULL DEFAULT '',
			action_url          TEXT NOT NULL DEFAULT '',
			register            TEXT NOT NULL DEFAULT '1',
			active              TEXT NOT NULL DEFAULT '1',
			access_group        INTEGER NOT NULL DEFAULT 0,
			domain_id           INTEGER NOT NULL DEFAULT 0,
			config_id           INTEGER NOT NULL DEFAULT 0,
			last_modified       DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_servicegroup (
			id                    INTEGER PRIMARY KEY AUTOINCREMENT,
			servicegroup_name     TEXT NOT NULL DEFAULT '',
			alias                 TEXT NOT NULL DEFAULT '',
			members               INTEGER NOT NULL DEFAULT 0,
			servicegroup_members  INTEGER NOT NULL DEFAULT 0,
			notes                 TEXT NOT NULL DEFAULT '',
			notes_url             TEXT NOT NULL DEFAULT '',
			action_url            TEXT NOT NULL DEFAULT '',
			register              TEXT NOT NULL DEFAULT '1',
			active                TEXT NOT NULL DEFAULT '1',
			access_group          INTEGER NOT NULL DEFAULT 0,
			domain_id             INTEGER NOT NULL DEFAULT 0,
			config_id             INTEGER NOT NULL DEFAULT 0,
			last_modified         DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_hosttemplate (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			template_name TEXT NOT NULL DEFAULT '',
			alias         TEXT NOT NULL DEFAULT '',
			register      TEXT NOT NULL DEFAULT '1',
			active        TEXT NOT NULL DEFAULT '1',
			domain_id     INTEGER NOT NULL DEFAULT 0,
			config_id     INTEGER NOT NULL DEFAULT 0,
			last_modified DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_servicetemplate (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			template_name TEXT NOT NULL DEFAULT '',
			alias         TEXT NOT NULL DEFAULT '',
			register      TEXT NOT NULL DEFAULT '1',
			active        TEXT NOT NULL DEFAULT '1',
			domain_id     INTEGER NOT NULL DEFAULT 0,
			config_id     INTEGER NOT NULL DEFAULT 0,
			last_modified DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_contacttemplate (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			template_name TEXT NOT NULL DEFAULT '',
			alias         TEXT NOT NULL DEFAULT '',
			register      TEXT NOT NULL DEFAULT '1',
			active        TEXT NOT NULL DEFAULT '1',
			domain_id     INTEGER NOT NULL DEFAULT 0,
			config_id     INTEGER NOT NULL DEFAULT 0,
			last_modified DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_lnkServiceToHost (
			idMaster  INTEGER NOT NULL DEFAULT 0,
			idSlave   INTEGER NOT NULL DEFAULT 0,
			idSort    INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_lnkHostToHostgroup (
			idMaster  INTEGER NOT NULL DEFAULT 0,
			idSlave   INTEGER NOT NULL DEFAULT 0,
			idSort    INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_lnkHostToContactgroup (
			idMaster  INTEGER NOT NULL DEFAULT 0,
			idSlave   INTEGER NOT NULL DEFAULT 0,
			idSort    INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_lnkServiceToServicegroup (
			idMaster  INTEGER NOT NULL DEFAULT 0,
			idSlave   INTEGER NOT NULL DEFAULT 0,
			idSort    INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_lnkContactToContactgroup (
			idMaster  INTEGER NOT NULL DEFAULT 0,
			idSlave   INTEGER NOT NULL DEFAULT 0,
			idSort    INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_lnkContactgroupToHost (
			idMaster  INTEGER NOT NULL DEFAULT 0,
			idSlave   INTEGER NOT NULL DEFAULT 0,
			idSort    INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_lnkContactgroupToService (
			idMaster  INTEGER NOT NULL DEFAULT 0,
			idSlave   INTEGER NOT NULL DEFAULT 0,
			idSort    INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_hostdependency (
			id                           INTEGER PRIMARY KEY AUTOINCREMENT,
			config_name                  TEXT NOT NULL DEFAULT '',
			dependent_host_name          INTEGER NOT NULL DEFAULT 0,
			dependent_hostgroup_name     INTEGER NOT NULL DEFAULT 0,
			host_name                    INTEGER NOT NULL DEFAULT 0,
			hostgroup_name               INTEGER NOT NULL DEFAULT 0,
			inherits_parent              INTEGER NOT NULL DEFAULT 0,
			execution_failure_criteria   TEXT NOT NULL DEFAULT '',
			notification_failure_criteria TEXT NOT NULL DEFAULT '',
			dependency_period            INTEGER NOT NULL DEFAULT 0,
			register                     TEXT NOT NULL DEFAULT '1',
			active                       TEXT NOT NULL DEFAULT '1',
			last_modified                DATETIME,
			access_group                 INTEGER NOT NULL DEFAULT 0,
			config_id                    INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_hostescalation (
			id                    INTEGER PRIMARY KEY AUTOINCREMENT,
			config_name           TEXT NOT NULL DEFAULT '',
			host_name             INTEGER NOT NULL DEFAULT 0,
			hostgroup_name        INTEGER NOT NULL DEFAULT 0,
			contacts              INTEGER NOT NULL DEFAULT 0,
			contact_groups        INTEGER NOT NULL DEFAULT 0,
			first_notification    INTEGER,
			last_notification     INTEGER,
			notification_interval INTEGER,
			escalation_period     INTEGER NOT NULL DEFAULT 0,
			escalation_options    TEXT NOT NULL DEFAULT '',
			register              TEXT NOT NULL DEFAULT '1',
			active                TEXT NOT NULL DEFAULT '1',
			last_modified         DATETIME,
			access_group          INTEGER NOT NULL DEFAULT 0,
			config_id             INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_hostextinfo (
			id              INTEGER PRIMARY KEY AUTOINCREMENT,
			host_name       INTEGER NOT NULL DEFAULT 0,
			notes           TEXT NOT NULL DEFAULT '',
			notes_url       TEXT NOT NULL DEFAULT '',
			action_url      TEXT NOT NULL DEFAULT '',
			icon_image      TEXT NOT NULL DEFAULT '',
			icon_image_alt  TEXT NOT NULL DEFAULT '',
			vrml_image      TEXT NOT NULL DEFAULT '',
			statusmap_image TEXT NOT NULL DEFAULT '',
			"2d_coords"     TEXT NOT NULL DEFAULT '',
			"3d_coords"     TEXT NOT NULL DEFAULT '',
			register        TEXT NOT NULL DEFAULT '1',
			active          TEXT NOT NULL DEFAULT '1',
			last_modified   DATETIME,
			access_group    INTEGER NOT NULL DEFAULT 0,
			config_id       INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_servicedependency (
			id                            INTEGER PRIMARY KEY AUTOINCREMENT,
			config_name                   TEXT NOT NULL DEFAULT '',
			dependent_host_name           INTEGER NOT NULL DEFAULT 0,
			dependent_hostgroup_name      INTEGER NOT NULL DEFAULT 0,
			dependent_service_description INTEGER NOT NULL DEFAULT 0,
			dependent_servicegroup_name   INTEGER NOT NULL DEFAULT 0,
			host_name                     INTEGER NOT NULL DEFAULT 0,
			hostgroup_name                INTEGER NOT NULL DEFAULT 0,
			service_description           INTEGER NOT NULL DEFAULT 0,
			servicegroup_name             INTEGER NOT NULL DEFAULT 0,
			inherits_parent               INTEGER NOT NULL DEFAULT 0,
			execution_failure_criteria    TEXT NOT NULL DEFAULT '',
			notification_failure_criteria TEXT NOT NULL DEFAULT '',
			dependency_period             INTEGER NOT NULL DEFAULT 0,
			register                      TEXT NOT NULL DEFAULT '1',
			active                        TEXT NOT NULL DEFAULT '1',
			last_modified                 DATETIME,
			access_group                  INTEGER NOT NULL DEFAULT 0,
			config_id                     INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_serviceescalation (
			id                    INTEGER PRIMARY KEY AUTOINCREMENT,
			config_name           TEXT NOT NULL DEFAULT '',
			host_name             INTEGER NOT NULL DEFAULT 0,
			hostgroup_name        INTEGER NOT NULL DEFAULT 0,
			service_description   INTEGER NOT NULL DEFAULT 0,
			servicegroup_name     INTEGER NOT NULL DEFAULT 0,
			contacts              INTEGER NOT NULL DEFAULT 0,
			contact_groups        INTEGER NOT NULL DEFAULT 0,
			first_notification    INTEGER,
			last_notification     INTEGER,
			notification_interval INTEGER,
			escalation_period     INTEGER NOT NULL DEFAULT 0,
			escalation_options    TEXT NOT NULL DEFAULT '',
			register              TEXT NOT NULL DEFAULT '1',
			active                TEXT NOT NULL DEFAULT '1',
			last_modified         DATETIME,
			access_group          INTEGER NOT NULL DEFAULT 0,
			config_id             INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS tbl_serviceextinfo (
			id                  INTEGER PRIMARY KEY AUTOINCREMENT,
			host_name           INTEGER NOT NULL DEFAULT 0,
			service_description INTEGER NOT NULL DEFAULT 0,
			notes               TEXT NOT NULL DEFAULT '',
			notes_url           TEXT NOT NULL DEFAULT '',
			action_url          TEXT NOT NULL DEFAULT '',
			icon_image          TEXT NOT NULL DEFAULT '',
			icon_image_alt      TEXT NOT NULL DEFAULT '',
			register            TEXT NOT NULL DEFAULT '1',
			active              TEXT NOT NULL DEFAULT '1',
			last_modified       DATETIME,
			access_group        INTEGER NOT NULL DEFAULT 0,
			config_id           INTEGER NOT NULL DEFAULT 0
		)`,
	}
}
