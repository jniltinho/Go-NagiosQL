package models

import "time"

// Service maps to tbl_service.
// config_name groups services that are written to the same .cfg file
// under the services/ directory (one file per config_name value).
type Service struct {
	ID                     uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ConfigName             string    `gorm:"column:config_name;size:255;not null" json:"config_name"`
	HostName               uint8     `gorm:"column:host_name;not null;default:0" json:"host_name"`
	HostNameTploptions     uint8     `gorm:"column:host_name_tploptions;not null;default:2" json:"host_name_tploptions"`
	HostgroupName          uint8     `gorm:"column:hostgroup_name;not null;default:0" json:"hostgroup_name"`
	HostgroupNameTploptions uint8    `gorm:"column:hostgroup_name_tploptions;not null;default:2" json:"hostgroup_name_tploptions"`
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
	ParallelizeCheck       uint8     `gorm:"column:parallelize_check;not null;default:1" json:"parallelize_check"`
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
	Name                   string    `gorm:"column:name;size:255;not null" json:"name"`
	Register               string    `gorm:"column:register;type:enum('0','1');not null;default:'1'" json:"register"`
	Active                 string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	LastModified           time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
	AccessGroup            uint      `gorm:"column:access_group;not null;default:0" json:"access_group"`
	ConfigID               uint8     `gorm:"column:config_id;not null;default:0" json:"config_id"`
	ImportHash             string    `gorm:"column:import_hash;size:255;not null" json:"import_hash"`
}

func (Service) TableName() string { return "tbl_service" }
