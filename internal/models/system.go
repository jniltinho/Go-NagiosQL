package models

import "time"

// User maps to tbl_user.
// Password is always bcrypt ($2a$12$...) for new users.
// Legacy MD5 hashes from the PHP project start with a hex string (no "$2" prefix)
// and block login with a password-reset prompt.
type User struct {
	ID           uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"column:username;size:64;not null;uniqueIndex" json:"username"`
	Password     string    `gorm:"column:password;size:255;not null" json:"-"`
	Name         string    `gorm:"column:name;size:255;not null" json:"name"`
	Email        string    `gorm:"column:email;size:255;not null" json:"email"`
	Admin        string    `gorm:"column:admin;type:enum('0','1');not null;default:'0'" json:"admin"`
	DomainID     uint      `gorm:"column:domain_id;not null;default:0" json:"domain_id"`
	LogonTimeout uint      `gorm:"column:logon_timeout;not null;default:60" json:"logon_timeout"`
	Active       string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	LastModified time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
}

func (User) TableName() string { return "tbl_user" }

// Datadomain maps to tbl_datadomain.
// The row with ID=0 is the "common" domain shared by all; it must be
// inserted with db.Exec (not GORM Save) because GORM skips zero-value PKs.
type Datadomain struct {
	ID           uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"column:name;size:255;not null" json:"name"`
	Description  string    `gorm:"column:description;size:255;not null" json:"description"`
	Active       string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	LastModified time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
}

func (Datadomain) TableName() string { return "tbl_datadomain" }

// Configtarget maps to tbl_configtarget.
// One row per datadomain; row for domain_id=0 holds the global Nagios paths.
type Configtarget struct {
	ID                 uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DomainID           uint      `gorm:"column:domain_id;not null;default:0" json:"domain_id"`
	NagiosCfg          string    `gorm:"column:nagios_cfg;size:255;not null" json:"nagios_cfg"`
	CgiCfg             string    `gorm:"column:cgi_cfg;size:255;not null" json:"cgi_cfg"`
	ResourceCfg        string    `gorm:"column:resource_cfg;size:255;not null" json:"resource_cfg"`
	BrokerModule       string    `gorm:"column:broker_module;size:255;not null" json:"broker_module"`
	CommandFile        string    `gorm:"column:command_file;size:255;not null" json:"command_file"`
	HostPath           string    `gorm:"column:host_path;size:255;not null" json:"host_path"`
	ServicePath        string    `gorm:"column:service_path;size:255;not null" json:"service_path"`
	CheckPath          string    `gorm:"column:check_path;size:255;not null" json:"check_path"`
	UserPath           string    `gorm:"column:user_path;size:255;not null" json:"user_path"`
	ImportPath         string    `gorm:"column:import_path;size:255;not null" json:"import_path"`
	BackupPath         string    `gorm:"column:backup_path;size:255;not null" json:"backup_path"`
	LogPath            string    `gorm:"column:log_path;size:255;not null" json:"log_path"`
	NagiosBin          string    `gorm:"column:nagios_bin;size:255;not null" json:"nagios_bin"`
	NagiosPID          string    `gorm:"column:nagios_pid;size:255;not null" json:"nagios_pid"`
	TargetName         string    `gorm:"column:target_name;size:255;not null" json:"target_name"`
	TargetDescription  string    `gorm:"column:target_description;size:255;not null" json:"target_description"`
	WriteToFiles       string    `gorm:"column:write_to_files;type:enum('0','1');not null;default:'1'" json:"write_to_files"`
	WriteToDatabase    string    `gorm:"column:write_to_database;type:enum('0','1');not null;default:'1'" json:"write_to_database"`
	Active             string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	LastModified       time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
}

func (Configtarget) TableName() string { return "tbl_configtarget" }

// Settings maps to tbl_settings.
// One row total; install_hash identifies the installation.
type Settings struct {
	ID                uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	InstallHash       string    `gorm:"column:install_hash;size:64;not null" json:"install_hash"`
	BackupAge         uint      `gorm:"column:backup_age;not null;default:14" json:"backup_age"`
	DateFormat        string    `gorm:"column:date_format;size:50;not null;default:'Y-m-d H:i:s'" json:"date_format"`
	LogbookContent    uint      `gorm:"column:logbook_content;not null;default:250" json:"logbook_content"`
	CryptoKey         string    `gorm:"column:crypto_key;size:255;not null" json:"crypto_key"`
	AuthenticationKey string    `gorm:"column:authentication_key;size:128;not null" json:"authentication_key"`
	LastModified      time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
}

func (Settings) TableName() string { return "tbl_settings" }

// Logbook maps to tbl_logbook. Stores an audit trail of all config changes.
type Logbook struct {
	ID           uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID       uint      `gorm:"column:user_id;not null;default:0" json:"user_id"`
	Username     string    `gorm:"column:username;size:64;not null" json:"username"`
	ObjectType   string    `gorm:"column:object_type;size:50;not null" json:"object_type"`
	ObjectName   string    `gorm:"column:object_name;size:255;not null" json:"object_name"`
	Action       string    `gorm:"column:action;size:50;not null" json:"action"`
	Info         string    `gorm:"column:info;type:text;not null" json:"info"`
	CreatedAt    time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

func (Logbook) TableName() string { return "tbl_logbook" }
