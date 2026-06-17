package models

import "time"

// Contact maps to tbl_contact.
type Contact struct {
	ID                           uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ContactName                  string    `gorm:"column:contact_name;size:255;not null" json:"contact_name"`
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
	Name                         string    `gorm:"column:name;size:255;not null" json:"name"`
	UseVariables                 uint8     `gorm:"column:use_variables;not null;default:0" json:"use_variables"`
	UseTemplate                  uint8     `gorm:"column:use_template;not null;default:0" json:"use_template"`
	UseTemplateTploptions        uint8     `gorm:"column:use_template_tploptions;not null;default:2" json:"use_template_tploptions"`
	Register                     string    `gorm:"column:register;type:enum('0','1');not null;default:'1'" json:"register"`
	Active                       string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	LastModified                 time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
	AccessGroup                  uint      `gorm:"column:access_group;not null;default:0" json:"access_group"`
	ConfigID                     uint8     `gorm:"column:config_id;not null;default:0" json:"config_id"`
}

func (Contact) TableName() string { return "tbl_contact" }
