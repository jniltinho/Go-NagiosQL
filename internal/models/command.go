// Package models contains GORM structs that map to the NagiosQL tbl_* schema.
// All tables use ENGINE=MyISAM as in the original NagiosQL 3.5 database.
package models

import "time"

// Command maps to tbl_command.
// command_type: 1 = check command, 2 = notification/event command.
type Command struct {
	ID          uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CommandName string    `gorm:"column:command_name;size:255;not null;uniqueIndex:config_name" json:"command_name"`
	CommandLine string    `gorm:"column:command_line;type:text;not null" json:"command_line"`
	CommandType uint8     `gorm:"column:command_type;not null;default:0" json:"command_type"`
	Arg1Info    *string   `gorm:"column:arg1_info;type:text" json:"arg1_info,omitempty"`
	Arg2Info    *string   `gorm:"column:arg2_info;type:text" json:"arg2_info,omitempty"`
	Arg3Info    *string   `gorm:"column:arg3_info;type:text" json:"arg3_info,omitempty"`
	Arg4Info    *string   `gorm:"column:arg4_info;type:text" json:"arg4_info,omitempty"`
	Arg5Info    *string   `gorm:"column:arg5_info;type:text" json:"arg5_info,omitempty"`
	Arg6Info    *string   `gorm:"column:arg6_info;type:text" json:"arg6_info,omitempty"`
	Arg7Info    *string   `gorm:"column:arg7_info;type:text" json:"arg7_info,omitempty"`
	Arg8Info    *string   `gorm:"column:arg8_info;type:text" json:"arg8_info,omitempty"`
	Register    string    `gorm:"column:register;type:enum('0','1');not null;default:'1'" json:"register"`
	Active      string    `gorm:"column:active;type:enum('0','1');not null;default:'1'" json:"active"`
	LastModified time.Time `gorm:"column:last_modified;not null" json:"last_modified"`
	AccessGroup uint      `gorm:"column:access_group;not null;default:0" json:"access_group"`
	ConfigID    uint8     `gorm:"column:config_id;not null;default:0;uniqueIndex:config_name" json:"config_id"`
}

func (Command) TableName() string { return "tbl_command" }
