package models

import "time"

// Timeperiod maps to tbl_timeperiod.
type Timeperiod struct {
	ID              uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TimeperiodName  string    `gorm:"column:timeperiod_name;size:255;not null" json:"timeperiod_name"`
	Alias           string    `gorm:"column:alias;size:255;not null" json:"alias"`
	Exclude         uint8     `gorm:"column:exclude;not null;default:0" json:"exclude"`
	UseTemplate     uint8     `gorm:"column:use_template;not null;default:0" json:"use_template"`
	Name            string    `gorm:"column:name;size:255;not null" json:"name"`
	Register        string    `gorm:"column:register;type:enum('0','1');not null;default:'1'" json:"register"`
	Active          string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	LastModified    time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
	AccessGroup     uint      `gorm:"column:access_group;not null;default:0" json:"access_group"`
	ConfigID        uint8     `gorm:"column:config_id;not null;default:0" json:"config_id"`

	// Definitions are the time ranges belonging to this period (tbl_timedefinition).
	Definitions []Timedefinition `gorm:"-" json:"definitions,omitempty"`
}

func (Timeperiod) TableName() string { return "tbl_timeperiod" }

// Timedefinition maps to tbl_timedefinition.
// Each row is one time range line within a timeperiod block,
// e.g. definition="monday" range="09:00-17:00".
type Timedefinition struct {
	ID           uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TipID        uint      `gorm:"column:tipId;not null" json:"tip_id"`
	Definition   string    `gorm:"column:definition;size:255;not null" json:"definition"`
	Range        string    `gorm:"column:range;type:text;not null" json:"range"`
	LastModified time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
}

func (Timedefinition) TableName() string { return "tbl_timedefinition" }
