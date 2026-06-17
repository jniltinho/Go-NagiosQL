package models

import "time"

// Hostdependency maps to tbl_hostdependency.
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

// Hostescalation maps to tbl_hostescalation.
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

// Hostextinfo maps to tbl_hostextinfo.
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

// Servicedependency maps to tbl_servicedependency.
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

// Serviceescalation maps to tbl_serviceescalation.
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

// Serviceextinfo maps to tbl_serviceextinfo.
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
