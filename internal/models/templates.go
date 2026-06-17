package models

import "time"

// Hosttemplate maps to tbl_hosttemplate.
type Hosttemplate struct {
	ID                     uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TemplateName           string    `gorm:"column:template_name;size:255;not null" json:"template_name"`
	Alias                  string    `gorm:"column:alias;size:255;not null" json:"alias"`
	DisplayName            string    `gorm:"column:display_name;size:255;not null" json:"display_name"`
	Address                string    `gorm:"column:address;size:255;not null" json:"address"`
	Parents                uint8     `gorm:"column:parents;not null;default:0" json:"parents"`
	ParentsTploptions      uint8     `gorm:"column:parents_tploptions;not null;default:2" json:"parents_tploptions"`
	Importance             *int      `gorm:"column:importance" json:"importance,omitempty"`
	Hostgroups             uint8     `gorm:"column:hostgroups;not null;default:0" json:"hostgroups"`
	HostgroupsTploptions   uint8     `gorm:"column:hostgroups_tploptions;not null;default:2" json:"hostgroups_tploptions"`
	CheckCommand           string    `gorm:"column:check_command;type:text;not null" json:"check_command"`
	UseTemplate            uint8     `gorm:"column:use_template;not null;default:0" json:"use_template"`
	UseTemplateTploptions  uint8     `gorm:"column:use_template_tploptions;not null;default:2" json:"use_template_tploptions"`
	InitialState           string    `gorm:"column:initial_state;size:20;not null" json:"initial_state"`
	MaxCheckAttempts       *int      `gorm:"column:max_check_attempts" json:"max_check_attempts,omitempty"`
	CheckInterval          *int      `gorm:"column:check_interval" json:"check_interval,omitempty"`
	RetryInterval          *int      `gorm:"column:retry_interval" json:"retry_interval,omitempty"`
	ActiveChecksEnabled    uint8     `gorm:"column:active_checks_enabled;not null;default:2" json:"active_checks_enabled"`
	PassiveChecksEnabled   uint8     `gorm:"column:passive_checks_enabled;not null;default:2" json:"passive_checks_enabled"`
	CheckPeriod            uint      `gorm:"column:check_period;not null;default:0" json:"check_period"`
	ObsessOverHost         uint8     `gorm:"column:obsess_over_host;not null;default:2" json:"obsess_over_host"`
	CheckFreshness         uint8     `gorm:"column:check_freshness;not null;default:2" json:"check_freshness"`
	FreshnessThreshold     *int      `gorm:"column:freshness_threshold" json:"freshness_threshold,omitempty"`
	EventHandler           uint      `gorm:"column:event_handler;not null;default:0" json:"event_handler"`
	EventHandlerEnabled    uint8     `gorm:"column:event_handler_enabled;not null;default:2" json:"event_handler_enabled"`
	LowFlapThreshold       *float64  `gorm:"column:low_flap_threshold" json:"low_flap_threshold,omitempty"`
	HighFlapThreshold      *float64  `gorm:"column:high_flap_threshold" json:"high_flap_threshold,omitempty"`
	FlapDetectionEnabled   uint8     `gorm:"column:flap_detection_enabled;not null;default:2" json:"flap_detection_enabled"`
	FlapDetectionOptions   string    `gorm:"column:flap_detection_options;size:255;not null" json:"flap_detection_options"`
	ProcessPerfData        uint8     `gorm:"column:process_perf_data;not null;default:2" json:"process_perf_data"`
	RetainStatusInfo       uint8     `gorm:"column:retain_status_information;not null;default:2" json:"retain_status_information"`
	RetainNonstatusInfo    uint8     `gorm:"column:retain_nonstatus_information;not null;default:2" json:"retain_nonstatus_information"`
	Contacts               uint8     `gorm:"column:contacts;not null;default:0" json:"contacts"`
	ContactsTploptions     uint8     `gorm:"column:contacts_tploptions;not null;default:2" json:"contacts_tploptions"`
	ContactGroups          uint8     `gorm:"column:contact_groups;not null;default:0" json:"contact_groups"`
	ContactGroupsTploptions uint8    `gorm:"column:contact_groups_tploptions;not null;default:2" json:"contact_groups_tploptions"`
	NotificationInterval   *int      `gorm:"column:notification_interval" json:"notification_interval,omitempty"`
	NotificationPeriod     uint      `gorm:"column:notification_period;not null;default:0" json:"notification_period"`
	FirstNotificationDelay *int      `gorm:"column:first_notification_delay" json:"first_notification_delay,omitempty"`
	NotificationOptions    string    `gorm:"column:notification_options;size:255;not null" json:"notification_options"`
	NotificationsEnabled   uint8     `gorm:"column:notifications_enabled;not null;default:2" json:"notifications_enabled"`
	StalkingOptions        string    `gorm:"column:stalking_options;size:255;not null" json:"stalking_options"`
	Notes                  string    `gorm:"column:notes;size:255;not null" json:"notes"`
	NotesURL               string    `gorm:"column:notes_url;size:255;not null" json:"notes_url"`
	ActionURL              string    `gorm:"column:action_url;size:255;not null" json:"action_url"`
	IconImage              string    `gorm:"column:icon_image;size:255;not null" json:"icon_image"`
	IconImageAlt           string    `gorm:"column:icon_image_alt;size:255;not null" json:"icon_image_alt"`
	VrmlImage              string    `gorm:"column:vrml_image;size:255;not null" json:"vrml_image"`
	StatusmapImage         string    `gorm:"column:statusmap_image;size:255;not null" json:"statusmap_image"`
	Coords2D               string    `gorm:"column:2d_coords;size:255;not null" json:"2d_coords"`
	Coords3D               string    `gorm:"column:3d_coords;size:255;not null" json:"3d_coords"`
	UseVariables           uint8     `gorm:"column:use_variables;not null;default:0" json:"use_variables"`
	Active                 string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	LastModified           time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
	AccessGroup            uint      `gorm:"column:access_group;not null;default:0" json:"access_group"`
	ConfigID               uint8     `gorm:"column:config_id;not null;default:0" json:"config_id"`
}

func (Hosttemplate) TableName() string { return "tbl_hosttemplate" }

// Servicetemplate maps to tbl_servicetemplate.
type Servicetemplate struct {
	ID                     uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TemplateName           string    `gorm:"column:template_name;size:255;not null" json:"template_name"`
	ServiceDescription     string    `gorm:"column:service_description;size:255;not null" json:"service_description"`
	DisplayName            string    `gorm:"column:display_name;size:255;not null" json:"display_name"`
	Servicegroups          uint8     `gorm:"column:servicegroups;not null;default:0" json:"servicegroups"`
	ServicegroupsTploptions uint8    `gorm:"column:servicegroups_tploptions;not null;default:2" json:"servicegroups_tploptions"`
	UseTemplate            uint8     `gorm:"column:use_template;not null;default:0" json:"use_template"`
	UseTemplateTploptions  uint8     `gorm:"column:use_template_tploptions;not null;default:2" json:"use_template_tploptions"`
	CheckCommand           string    `gorm:"column:check_command;type:text;not null" json:"check_command"`
	IsVolatile             uint8     `gorm:"column:is_volatile;not null;default:0" json:"is_volatile"`
	InitialState           string    `gorm:"column:initial_state;size:255;not null" json:"initial_state"`
	MaxCheckAttempts       *int      `gorm:"column:max_check_attempts" json:"max_check_attempts,omitempty"`
	CheckInterval          *int      `gorm:"column:check_interval" json:"check_interval,omitempty"`
	RetryInterval          *int      `gorm:"column:retry_interval" json:"retry_interval,omitempty"`
	ActiveChecksEnabled    uint8     `gorm:"column:active_checks_enabled;not null;default:1" json:"active_checks_enabled"`
	PassiveChecksEnabled   uint8     `gorm:"column:passive_checks_enabled;not null;default:1" json:"passive_checks_enabled"`
	CheckPeriod            uint      `gorm:"column:check_period;not null;default:0" json:"check_period"`
	ObsessOverService      uint8     `gorm:"column:obsess_over_service;not null;default:1" json:"obsess_over_service"`
	CheckFreshness         uint8     `gorm:"column:check_freshness;not null;default:0" json:"check_freshness"`
	FreshnessThreshold     *int      `gorm:"column:freshness_threshold" json:"freshness_threshold,omitempty"`
	EventHandler           uint      `gorm:"column:event_handler;not null;default:0" json:"event_handler"`
	EventHandlerEnabled    uint8     `gorm:"column:event_handler_enabled;not null;default:2" json:"event_handler_enabled"`
	LowFlapThreshold       *float64  `gorm:"column:low_flap_threshold" json:"low_flap_threshold,omitempty"`
	HighFlapThreshold      *float64  `gorm:"column:high_flap_threshold" json:"high_flap_threshold,omitempty"`
	FlapDetectionEnabled   uint8     `gorm:"column:flap_detection_enabled;not null;default:1" json:"flap_detection_enabled"`
	FlapDetectionOptions   string    `gorm:"column:flap_detection_options;size:255;not null" json:"flap_detection_options"`
	ProcessPerfData        uint8     `gorm:"column:process_perf_data;not null;default:1" json:"process_perf_data"`
	RetainStatusInfo       uint8     `gorm:"column:retain_status_information;not null;default:1" json:"retain_status_information"`
	RetainNonstatusInfo    uint8     `gorm:"column:retain_nonstatus_information;not null;default:1" json:"retain_nonstatus_information"`
	NotificationInterval   *int      `gorm:"column:notification_interval" json:"notification_interval,omitempty"`
	FirstNotificationDelay *int      `gorm:"column:first_notification_delay" json:"first_notification_delay,omitempty"`
	NotificationPeriod     uint      `gorm:"column:notification_period;not null;default:0" json:"notification_period"`
	NotificationOptions    string    `gorm:"column:notification_options;size:255;not null" json:"notification_options"`
	NotificationsEnabled   uint8     `gorm:"column:notifications_enabled;not null;default:1" json:"notifications_enabled"`
	Contacts               uint8     `gorm:"column:contacts;not null;default:0" json:"contacts"`
	ContactsTploptions     uint8     `gorm:"column:contacts_tploptions;not null;default:2" json:"contacts_tploptions"`
	ContactGroups          uint8     `gorm:"column:contact_groups;not null;default:0" json:"contact_groups"`
	ContactGroupsTploptions uint8    `gorm:"column:contact_groups_tploptions;not null;default:2" json:"contact_groups_tploptions"`
	StalkingOptions        string    `gorm:"column:stalking_options;size:255;not null" json:"stalking_options"`
	Notes                  string    `gorm:"column:notes;size:255;not null" json:"notes"`
	NotesURL               string    `gorm:"column:notes_url;size:255;not null" json:"notes_url"`
	ActionURL              string    `gorm:"column:action_url;size:255;not null" json:"action_url"`
	IconImage              string    `gorm:"column:icon_image;size:255;not null" json:"icon_image"`
	IconImageAlt           string    `gorm:"column:icon_image_alt;size:255;not null" json:"icon_image_alt"`
	UseVariables           uint8     `gorm:"column:use_variables;not null;default:0" json:"use_variables"`
	Active                 string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	LastModified           time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
	AccessGroup            uint      `gorm:"column:access_group;not null;default:0" json:"access_group"`
	ConfigID               uint8     `gorm:"column:config_id;not null;default:0" json:"config_id"`
}

func (Servicetemplate) TableName() string { return "tbl_servicetemplate" }

// Contacttemplate maps to tbl_contacttemplate.
type Contacttemplate struct {
	ID                           uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TemplateName                 string    `gorm:"column:template_name;size:255;not null" json:"template_name"`
	Alias                        string    `gorm:"column:alias;size:255;not null" json:"alias"`
	Contactgroups                uint8     `gorm:"column:contactgroups;not null;default:0" json:"contactgroups"`
	ContactgroupsTploptions      uint8     `gorm:"column:contactgroups_tploptions;not null;default:2" json:"contactgroups_tploptions"`
	HostNotificationsEnabled     uint8     `gorm:"column:host_notifications_enabled;not null;default:2" json:"host_notifications_enabled"`
	ServiceNotificationsEnabled  uint8     `gorm:"column:service_notifications_enabled;not null;default:2" json:"service_notifications_enabled"`
	HostNotificationPeriod       uint      `gorm:"column:host_notification_period;not null;default:0" json:"host_notification_period"`
	ServiceNotificationPeriod    uint      `gorm:"column:service_notification_period;not null;default:0" json:"service_notification_period"`
	HostNotificationOptions      string    `gorm:"column:host_notification_options;size:255;not null" json:"host_notification_options"`
	ServiceNotificationOptions   string    `gorm:"column:service_notification_options;size:255;not null" json:"service_notification_options"`
	HostNotificationCommands     uint8     `gorm:"column:host_notification_commands;not null;default:0" json:"host_notification_commands"`
	HostNotifCmdTploptions       uint8     `gorm:"column:host_notification_commands_tploptions;not null;default:2" json:"host_notification_commands_tploptions"`
	ServiceNotificationCommands  uint8     `gorm:"column:service_notification_commands;not null;default:0" json:"service_notification_commands"`
	ServiceNotifCmdTploptions    uint8     `gorm:"column:service_notification_commands_tploptions;not null;default:2" json:"service_notification_commands_tploptions"`
	CanSubmitCommands            uint8     `gorm:"column:can_submit_commands;not null;default:2" json:"can_submit_commands"`
	RetainStatusInfo             uint8     `gorm:"column:retain_status_information;not null;default:2" json:"retain_status_information"`
	RetainNonstatusInfo          uint8     `gorm:"column:retain_nonstatus_information;not null;default:2" json:"retain_nonstatus_information"`
	Email                        string    `gorm:"column:email;size:255;not null" json:"email"`
	Pager                        string    `gorm:"column:pager;size:255;not null" json:"pager"`
	Address1                     string    `gorm:"column:address1;size:255;not null" json:"address1"`
	Address2                     string    `gorm:"column:address2;size:255;not null" json:"address2"`
	Address3                     string    `gorm:"column:address3;size:255;not null" json:"address3"`
	Address4                     string    `gorm:"column:address4;size:255;not null" json:"address4"`
	Address5                     string    `gorm:"column:address5;size:255;not null" json:"address5"`
	Address6                     string    `gorm:"column:address6;size:255;not null" json:"address6"`
	UseVariables                 uint8     `gorm:"column:use_variables;not null;default:0" json:"use_variables"`
	UseTemplate                  uint8     `gorm:"column:use_template;not null;default:0" json:"use_template"`
	UseTemplateTploptions        uint8     `gorm:"column:use_template_tploptions;not null;default:2" json:"use_template_tploptions"`
	Active                       string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	LastModified                 time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
	AccessGroup                  uint      `gorm:"column:access_group;not null;default:0" json:"access_group"`
	ConfigID                     uint8     `gorm:"column:config_id;not null;default:0" json:"config_id"`
}

func (Contacttemplate) TableName() string { return "tbl_contacttemplate" }
