package models

import "time"

// Hostdependency represents a Nagios host-dependency definition stored in
// tbl_hostdependency. A host dependency tells Nagios to suppress checks or
// notifications for a dependent host when the master host is in a particular
// state.
//
// FK flag fields — DependentHostName, DependentHostgroupName, HostName, and
// HostgroupName — are uint8 presence flags (0 = no members set, 1 = at least
// one member linked). The actual host or hostgroup names are stored in the
// corresponding join tables and resolved at config-generation time:
//
//   - tbl_lnkHostdependencyToHost_DH    → dependent_host_name
//   - tbl_lnkHostdependencyToHostgroup_DH → dependent_hostgroup_name
//   - tbl_lnkHostdependencyToHost_H     → host_name
//   - tbl_lnkHostdependencyToHostgroup_H → hostgroup_name
type Hostdependency struct {
	ID                          uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ConfigName                  string    `gorm:"column:config_name;size:255;not null" json:"config_name"`
	DependentHostName           uint8     `gorm:"column:dependent_host_name;not null;default:0" json:"dependent_host_name"`
	DependentHostgroupName      uint8     `gorm:"column:dependent_hostgroup_name;not null;default:0" json:"dependent_hostgroup_name"`
	HostName                    uint8     `gorm:"column:host_name;not null;default:0" json:"host_name"`
	HostgroupName               uint8     `gorm:"column:hostgroup_name;not null;default:0" json:"hostgroup_name"`
	InheritsParent              uint8     `gorm:"column:inherits_parent;not null;default:0" json:"inherits_parent"`
	ExecutionFailureCriteria    string    `gorm:"column:execution_failure_criteria;size:255" json:"execution_failure_criteria"`
	NotificationFailureCriteria string    `gorm:"column:notification_failure_criteria;size:255" json:"notification_failure_criteria"`
	DependencyPeriod            uint      `gorm:"column:dependency_period;not null;default:0" json:"dependency_period"`
	Register                    string    `gorm:"column:register;type:enum('0','1');not null;default:'1'" json:"register"`
	Active                      string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	LastModified                time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
	AccessGroup                 uint      `gorm:"column:access_group;not null;default:0" json:"access_group"`
	ConfigID                    uint8     `gorm:"column:config_id;not null;default:0" json:"config_id"`
}

func (Hostdependency) TableName() string { return "tbl_hostdependency" }

// Hostescalation represents a Nagios host-escalation definition stored in
// tbl_hostescalation. A host escalation overrides normal notification
// behaviour — contacts, intervals, and options — after a host has been in a
// problem state for a configurable number of notifications.
//
// FK flag fields — HostName, HostgroupName, Contacts, and ContactGroups — are
// uint8 presence flags (0 = no members set, 1 = at least one member linked).
// The actual names are resolved from join tables at config-generation time:
//
//   - tbl_lnkHostescalationToHost         → host_name
//   - tbl_lnkHostescalationToHostgroup     → hostgroup_name
//   - tbl_lnkHostescalationToContact       → contacts
//   - tbl_lnkHostescalationToContactgroup  → contact_groups
type Hostescalation struct {
	ID                   uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ConfigName           string    `gorm:"column:config_name;size:255;not null" json:"config_name"`
	HostName             uint8     `gorm:"column:host_name;not null;default:0" json:"host_name"`
	HostgroupName        uint8     `gorm:"column:hostgroup_name;not null;default:0" json:"hostgroup_name"`
	Contacts             uint8     `gorm:"column:contacts;not null;default:0" json:"contacts"`
	ContactGroups        uint8     `gorm:"column:contact_groups;not null;default:0" json:"contact_groups"`
	FirstNotification    *int      `gorm:"column:first_notification" json:"first_notification,omitempty"`
	LastNotification     *int      `gorm:"column:last_notification" json:"last_notification,omitempty"`
	NotificationInterval *int      `gorm:"column:notification_interval" json:"notification_interval,omitempty"`
	EscalationPeriod     uint      `gorm:"column:escalation_period;not null;default:0" json:"escalation_period"`
	EscalationOptions    string    `gorm:"column:escalation_options;size:255" json:"escalation_options"`
	Register             string    `gorm:"column:register;type:enum('0','1');not null;default:'1'" json:"register"`
	Active               string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	LastModified         time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
	AccessGroup          uint      `gorm:"column:access_group;not null;default:0" json:"access_group"`
	ConfigID             uint8     `gorm:"column:config_id;not null;default:0" json:"config_id"`
}

func (Hostescalation) TableName() string { return "tbl_hostescalation" }

// Hostextinfo represents a Nagios extended-host-information definition stored
// in tbl_hostextinfo. Extended host info attaches presentation metadata —
// icon images, status-map coordinates, notes, and URLs — to a host without
// changing its check or notification behaviour.
//
// Unlike other structs in this package, HostName is NOT a uint8 FK flag.
// It is the integer primary-key ID of a row in tbl_host; the config generator
// calls resolveHostByID to look up the actual host_name string at generation
// time.
type Hostextinfo struct {
	ID             uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	HostName       uint      `gorm:"column:host_name;not null" json:"host_name"`
	Notes          string    `gorm:"column:notes;size:255" json:"notes"`
	NotesURL       string    `gorm:"column:notes_url;size:255" json:"notes_url"`
	ActionURL      string    `gorm:"column:action_url;size:255" json:"action_url"`
	IconImage      string    `gorm:"column:icon_image;size:255" json:"icon_image"`
	IconImageAlt   string    `gorm:"column:icon_image_alt;size:255" json:"icon_image_alt"`
	VrmlImage      string    `gorm:"column:vrml_image;size:255" json:"vrml_image"`
	StatusmapImage string    `gorm:"column:statusmap_image;size:255" json:"statusmap_image"`
	Coords2D       string    `gorm:"column:2d_coords;size:255" json:"2d_coords"`
	Coords3D       string    `gorm:"column:3d_coords;size:255" json:"3d_coords"`
	Register       string    `gorm:"column:register;type:enum('0','1');not null;default:'1'" json:"register"`
	Active         string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	LastModified   time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
	AccessGroup    uint      `gorm:"column:access_group;not null;default:0" json:"access_group"`
	ConfigID       uint8     `gorm:"column:config_id;not null;default:0" json:"config_id"`
}

func (Hostextinfo) TableName() string { return "tbl_hostextinfo" }

// Servicedependency represents a Nagios service-dependency definition stored
// in tbl_servicedependency. A service dependency tells Nagios to suppress
// checks or notifications for a dependent service when the master service is
// in a particular state.
//
// FK flag fields — DependentHostName, DependentHostgroupName,
// DependentServiceDescription, DependentServicegroupName, HostName,
// HostgroupName, ServiceDescription, and ServicegroupName — are uint8
// presence flags (0 = no members set, 1 = at least one member linked).
// The actual names are resolved from join tables at config-generation time:
//
//   - tbl_lnkServicedependencyToHost_DH        → dependent_host_name
//   - tbl_lnkServicedependencyToHostgroup_DH    → dependent_hostgroup_name
//   - tbl_lnkServicedependencyToService_DS      → dependent_service_description
//   - tbl_lnkServicedependencyToServicegroup_DS → dependent_servicegroup_name
//   - tbl_lnkServicedependencyToHost_H          → host_name
//   - tbl_lnkServicedependencyToHostgroup_H     → hostgroup_name
//   - tbl_lnkServicedependencyToService_S       → service_description
//   - tbl_lnkServicedependencyToServicegroup_S  → servicegroup_name
type Servicedependency struct {
	ID                          uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ConfigName                  string    `gorm:"column:config_name;size:255;not null" json:"config_name"`
	DependentHostName           uint8     `gorm:"column:dependent_host_name;not null;default:0" json:"dependent_host_name"`
	DependentHostgroupName      uint8     `gorm:"column:dependent_hostgroup_name;not null;default:0" json:"dependent_hostgroup_name"`
	DependentServiceDescription uint8     `gorm:"column:dependent_service_description;not null;default:0" json:"dependent_service_description"`
	DependentServicegroupName   uint8     `gorm:"column:dependent_servicegroup_name;not null;default:0" json:"dependent_servicegroup_name"`
	HostName                    uint8     `gorm:"column:host_name;not null;default:0" json:"host_name"`
	HostgroupName               uint8     `gorm:"column:hostgroup_name;not null;default:0" json:"hostgroup_name"`
	ServiceDescription          uint8     `gorm:"column:service_description;not null;default:0" json:"service_description"`
	ServicegroupName            uint8     `gorm:"column:servicegroup_name;not null;default:0" json:"servicegroup_name"`
	InheritsParent              uint8     `gorm:"column:inherits_parent;not null;default:0" json:"inherits_parent"`
	ExecutionFailureCriteria    string    `gorm:"column:execution_failure_criteria;size:255" json:"execution_failure_criteria"`
	NotificationFailureCriteria string    `gorm:"column:notification_failure_criteria;size:255" json:"notification_failure_criteria"`
	DependencyPeriod            uint      `gorm:"column:dependency_period;not null;default:0" json:"dependency_period"`
	Register                    string    `gorm:"column:register;type:enum('0','1');not null;default:'1'" json:"register"`
	Active                      string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	LastModified                time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
	AccessGroup                 uint      `gorm:"column:access_group;not null;default:0" json:"access_group"`
	ConfigID                    uint8     `gorm:"column:config_id;not null;default:0" json:"config_id"`
}

func (Servicedependency) TableName() string { return "tbl_servicedependency" }

// Serviceescalation represents a Nagios service-escalation definition stored
// in tbl_serviceescalation. A service escalation overrides normal notification
// behaviour — contacts, intervals, and options — after a service has been in a
// problem state for a configurable number of notifications.
//
// FK flag fields — HostName, HostgroupName, ServiceDescription,
// ServicegroupName, Contacts, and ContactGroups — are uint8 presence flags
// (0 = no members set, 1 = at least one member linked). The actual names are
// resolved from join tables at config-generation time:
//
//   - tbl_lnkServiceescalationToHost          → host_name
//   - tbl_lnkServiceescalationToHostgroup      → hostgroup_name
//   - tbl_lnkServiceescalationToService        → service_description
//   - tbl_lnkServiceescalationToServicegroup   → servicegroup_name
//   - tbl_lnkServiceescalationToContact        → contacts
//   - tbl_lnkServiceescalationToContactgroup   → contact_groups
type Serviceescalation struct {
	ID                   uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ConfigName           string    `gorm:"column:config_name;size:255;not null" json:"config_name"`
	HostName             uint8     `gorm:"column:host_name;not null;default:0" json:"host_name"`
	HostgroupName        uint8     `gorm:"column:hostgroup_name;not null;default:0" json:"hostgroup_name"`
	ServiceDescription   uint8     `gorm:"column:service_description;not null;default:0" json:"service_description"`
	ServicegroupName     uint8     `gorm:"column:servicegroup_name;not null;default:0" json:"servicegroup_name"`
	Contacts             uint8     `gorm:"column:contacts;not null;default:0" json:"contacts"`
	ContactGroups        uint8     `gorm:"column:contact_groups;not null;default:0" json:"contact_groups"`
	FirstNotification    *int      `gorm:"column:first_notification" json:"first_notification,omitempty"`
	LastNotification     *int      `gorm:"column:last_notification" json:"last_notification,omitempty"`
	NotificationInterval *int      `gorm:"column:notification_interval" json:"notification_interval,omitempty"`
	EscalationPeriod     uint      `gorm:"column:escalation_period;not null;default:0" json:"escalation_period"`
	EscalationOptions    string    `gorm:"column:escalation_options;size:255" json:"escalation_options"`
	Register             string    `gorm:"column:register;type:enum('0','1');not null;default:'1'" json:"register"`
	Active               string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	LastModified         time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
	AccessGroup          uint      `gorm:"column:access_group;not null;default:0" json:"access_group"`
	ConfigID             uint8     `gorm:"column:config_id;not null;default:0" json:"config_id"`
}

func (Serviceescalation) TableName() string { return "tbl_serviceescalation" }

// Serviceextinfo represents a Nagios extended-service-information definition
// stored in tbl_serviceextinfo. Extended service info attaches presentation
// metadata — icon images, notes, and URLs — to a service without changing its
// check or notification behaviour.
//
// Unlike other structs in this package, HostName and ServiceDescription are
// NOT uint8 FK flags. They are integer primary-key IDs of rows in tbl_host
// and tbl_service respectively; the config generator calls resolveHostByID and
// resolveServiceByID to look up the actual name strings at generation time.
type Serviceextinfo struct {
	ID                 uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	HostName           uint      `gorm:"column:host_name;not null" json:"host_name"`
	ServiceDescription uint      `gorm:"column:service_description;not null" json:"service_description"`
	Notes              string    `gorm:"column:notes;size:255" json:"notes"`
	NotesURL           string    `gorm:"column:notes_url;size:255" json:"notes_url"`
	ActionURL          string    `gorm:"column:action_url;size:255" json:"action_url"`
	IconImage          string    `gorm:"column:icon_image;size:500" json:"icon_image"`
	IconImageAlt       string    `gorm:"column:icon_image_alt;size:255" json:"icon_image_alt"`
	Register           string    `gorm:"column:register;type:enum('0','1');not null;default:'1'" json:"register"`
	Active             string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	LastModified       time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
	AccessGroup        uint      `gorm:"column:access_group;not null;default:0" json:"access_group"`
	ConfigID           uint8     `gorm:"column:config_id;not null;default:0" json:"config_id"`
}

func (Serviceextinfo) TableName() string { return "tbl_serviceextinfo" }
