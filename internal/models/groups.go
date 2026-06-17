package models

import "time"

// Contactgroup maps to tbl_contactgroup.
type Contactgroup struct {
	ID               uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ContactgroupName string    `gorm:"column:contactgroup_name;size:255;not null" json:"contactgroup_name"`
	Alias            string    `gorm:"column:alias;size:255;not null" json:"alias"`
	Members          uint8     `gorm:"column:members;not null;default:0" json:"members"`
	ContactgroupMembers uint8  `gorm:"column:contactgroup_members;not null;default:0" json:"contactgroup_members"`
	Register         string    `gorm:"column:register;type:enum('0','1');not null;default:'1'" json:"register"`
	Active           string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	LastModified     time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
	AccessGroup      uint      `gorm:"column:access_group;not null;default:0" json:"access_group"`
	ConfigID         uint8     `gorm:"column:config_id;not null;default:0" json:"config_id"`
}

func (Contactgroup) TableName() string { return "tbl_contactgroup" }

// Hostgroup maps to tbl_hostgroup.
type Hostgroup struct {
	ID            uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	HostgroupName string    `gorm:"column:hostgroup_name;size:255;not null" json:"hostgroup_name"`
	Alias         string    `gorm:"column:alias;size:255;not null" json:"alias"`
	Members       uint8     `gorm:"column:members;not null;default:0" json:"members"`
	HostgroupMembers uint8  `gorm:"column:hostgroup_members;not null;default:0" json:"hostgroup_members"`
	Notes         string    `gorm:"column:notes;size:255;not null" json:"notes"`
	NotesURL      string    `gorm:"column:notes_url;size:255;not null" json:"notes_url"`
	ActionURL     string    `gorm:"column:action_url;size:255;not null" json:"action_url"`
	Register      string    `gorm:"column:register;type:enum('0','1');not null;default:'1'" json:"register"`
	Active        string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	LastModified  time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
	AccessGroup   uint      `gorm:"column:access_group;not null;default:0" json:"access_group"`
	ConfigID      uint8     `gorm:"column:config_id;not null;default:0" json:"config_id"`
}

func (Hostgroup) TableName() string { return "tbl_hostgroup" }

// Servicegroup maps to tbl_servicegroup.
type Servicegroup struct {
	ID               uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ServicegroupName string    `gorm:"column:servicegroup_name;size:255;not null" json:"servicegroup_name"`
	Alias            string    `gorm:"column:alias;size:255;not null" json:"alias"`
	Members          uint8     `gorm:"column:members;not null;default:0" json:"members"`
	ServicegroupMembers uint8  `gorm:"column:servicegroup_members;not null;default:0" json:"servicegroup_members"`
	Notes            string    `gorm:"column:notes;size:255;not null" json:"notes"`
	NotesURL         string    `gorm:"column:notes_url;size:255;not null" json:"notes_url"`
	ActionURL        string    `gorm:"column:action_url;size:255;not null" json:"action_url"`
	Register         string    `gorm:"column:register;type:enum('0','1');not null;default:'1'" json:"register"`
	Active           string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	LastModified     time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
	AccessGroup      uint      `gorm:"column:access_group;not null;default:0" json:"access_group"`
	ConfigID         uint8     `gorm:"column:config_id;not null;default:0" json:"config_id"`
}

func (Servicegroup) TableName() string { return "tbl_servicegroup" }
