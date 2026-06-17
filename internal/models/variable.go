package models

import "time"

// Variabledefinition maps to tbl_variabledefinition.
// Custom Nagios variables (e.g. _SNMP_COMMUNITY) are stored here and linked
// to host, service or contact objects via the tbl_lnkXxxToVariabledefinition tables.
type Variabledefinition struct {
	ID           uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"column:name;size:255;not null" json:"name"`
	Value        string    `gorm:"column:value;type:text;not null" json:"value"`
	Vartype      string    `gorm:"column:vartype;size:50;not null;default:'string'" json:"vartype"`
	Active       string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	DomainID     uint      `gorm:"column:domain_id;not null;default:0" json:"domain_id"`
	ConfigID     uint8     `gorm:"column:config_id;not null;default:0" json:"config_id"`
	LastModified time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
}

func (Variabledefinition) TableName() string { return "tbl_variabledefinition" }
