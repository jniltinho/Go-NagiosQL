// Package models — link tables connecting object rows to related objects.
// NagiosQL uses a "has_link" flag (uint8 count) on the parent table plus
// a separate link table that stores the many-to-many relationship.
// Each link table follows the naming convention:
//   tbl_lnk<Parent><Child>  (e.g. tbl_lnkHostToHostgroup)
package models

// LnkHostToHostgroup maps to tbl_lnkHostToHostgroup.
type LnkHostToHostgroup struct {
	HostID       uint `gorm:"column:idMaster;not null;index:idx_master" json:"host_id"`
	HostgroupID  uint `gorm:"column:idSlave;not null" json:"hostgroup_id"`
	Sequence     uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkHostToHostgroup) TableName() string { return "tbl_lnkHostToHostgroup" }

// LnkHostToContact maps to tbl_lnkHostToContact.
type LnkHostToContact struct {
	HostID    uint `gorm:"column:idMaster;not null;index:idx_master" json:"host_id"`
	ContactID uint `gorm:"column:idSlave;not null" json:"contact_id"`
	Sequence  uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkHostToContact) TableName() string { return "tbl_lnkHostToContact" }

// LnkHostToContactgroup maps to tbl_lnkHostToContactgroup.
type LnkHostToContactgroup struct {
	HostID         uint `gorm:"column:idMaster;not null;index:idx_master" json:"host_id"`
	ContactgroupID uint `gorm:"column:idSlave;not null" json:"contactgroup_id"`
	Sequence       uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkHostToContactgroup) TableName() string { return "tbl_lnkHostToContactgroup" }

// LnkHostToHost maps to tbl_lnkHostToHost (parent-child host relationships).
type LnkHostToHost struct {
	HostID       uint `gorm:"column:idMaster;not null;index:idx_master" json:"host_id"`
	ParentHostID uint `gorm:"column:idSlave;not null" json:"parent_host_id"`
	Sequence     uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkHostToHost) TableName() string { return "tbl_lnkHostToHost" }

// LnkHostToHosttemplate maps to tbl_lnkHostToHosttemplate.
type LnkHostToHosttemplate struct {
	HostID         uint `gorm:"column:idMaster;not null;index:idx_master" json:"host_id"`
	HosttemplateID uint `gorm:"column:idSlave;not null" json:"hosttemplate_id"`
	Sequence       uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkHostToHosttemplate) TableName() string { return "tbl_lnkHostToHosttemplate" }

// --- Service links ---

// LnkServiceToHost maps to tbl_lnkServiceToHost.
type LnkServiceToHost struct {
	ServiceID uint `gorm:"column:idMaster;not null;index:idx_master" json:"service_id"`
	HostID    uint `gorm:"column:idSlave;not null" json:"host_id"`
	Sequence  uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkServiceToHost) TableName() string { return "tbl_lnkServiceToHost" }

// LnkServiceToHostgroup maps to tbl_lnkServiceToHostgroup.
type LnkServiceToHostgroup struct {
	ServiceID   uint `gorm:"column:idMaster;not null;index:idx_master" json:"service_id"`
	HostgroupID uint `gorm:"column:idSlave;not null" json:"hostgroup_id"`
	Sequence    uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkServiceToHostgroup) TableName() string { return "tbl_lnkServiceToHostgroup" }

// LnkServiceToContact maps to tbl_lnkServiceToContact.
type LnkServiceToContact struct {
	ServiceID uint `gorm:"column:idMaster;not null;index:idx_master" json:"service_id"`
	ContactID uint `gorm:"column:idSlave;not null" json:"contact_id"`
	Sequence  uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkServiceToContact) TableName() string { return "tbl_lnkServiceToContact" }

// LnkServiceToContactgroup maps to tbl_lnkServiceToContactgroup.
type LnkServiceToContactgroup struct {
	ServiceID      uint `gorm:"column:idMaster;not null;index:idx_master" json:"service_id"`
	ContactgroupID uint `gorm:"column:idSlave;not null" json:"contactgroup_id"`
	Sequence       uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkServiceToContactgroup) TableName() string { return "tbl_lnkServiceToContactgroup" }

// LnkServiceToServicegroup maps to tbl_lnkServiceToServicegroup.
type LnkServiceToServicegroup struct {
	ServiceID      uint `gorm:"column:idMaster;not null;index:idx_master" json:"service_id"`
	ServicegroupID uint `gorm:"column:idSlave;not null" json:"servicegroup_id"`
	Sequence       uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkServiceToServicegroup) TableName() string { return "tbl_lnkServiceToServicegroup" }

// LnkServiceToServicetemplate maps to tbl_lnkServiceToServicetemplate.
type LnkServiceToServicetemplate struct {
	ServiceID         uint `gorm:"column:idMaster;not null;index:idx_master" json:"service_id"`
	ServicetemplateID uint `gorm:"column:idSlave;not null" json:"servicetemplate_id"`
	Sequence          uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkServiceToServicetemplate) TableName() string { return "tbl_lnkServiceToServicetemplate" }

// --- Hosttemplate links ---

// LnkHosttemplateToHosttemplate maps to tbl_lnkHosttemplateToHosttemplate.
type LnkHosttemplateToHosttemplate struct {
	HosttemplateID       uint `gorm:"column:idMaster;not null;index:idx_master" json:"hosttemplate_id"`
	ParentHosttemplateID uint `gorm:"column:idSlave;not null" json:"parent_hosttemplate_id"`
	Sequence             uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkHosttemplateToHosttemplate) TableName() string {
	return "tbl_lnkHosttemplateToHosttemplate"
}

// LnkHosttemplateToHostgroup maps to tbl_lnkHosttemplateToHostgroup.
type LnkHosttemplateToHostgroup struct {
	HosttemplateID uint `gorm:"column:idMaster;not null;index:idx_master" json:"hosttemplate_id"`
	HostgroupID    uint `gorm:"column:idSlave;not null" json:"hostgroup_id"`
	Sequence       uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkHosttemplateToHostgroup) TableName() string { return "tbl_lnkHosttemplateToHostgroup" }

// LnkHosttemplateToContact maps to tbl_lnkHosttemplateToContact.
type LnkHosttemplateToContact struct {
	HosttemplateID uint `gorm:"column:idMaster;not null;index:idx_master" json:"hosttemplate_id"`
	ContactID      uint `gorm:"column:idSlave;not null" json:"contact_id"`
	Sequence       uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkHosttemplateToContact) TableName() string { return "tbl_lnkHosttemplateToContact" }

// LnkHosttemplateToContactgroup maps to tbl_lnkHosttemplateToContactgroup.
type LnkHosttemplateToContactgroup struct {
	HosttemplateID uint `gorm:"column:idMaster;not null;index:idx_master" json:"hosttemplate_id"`
	ContactgroupID uint `gorm:"column:idSlave;not null" json:"contactgroup_id"`
	Sequence       uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkHosttemplateToContactgroup) TableName() string {
	return "tbl_lnkHosttemplateToContactgroup"
}

// --- Servicetemplate links ---

// LnkServicetemplateToServicetemplate maps to tbl_lnkServicetemplateToServicetemplate.
type LnkServicetemplateToServicetemplate struct {
	ServicetemplateID       uint `gorm:"column:idMaster;not null;index:idx_master" json:"servicetemplate_id"`
	ParentServicetemplateID uint `gorm:"column:idSlave;not null" json:"parent_servicetemplate_id"`
	Sequence                uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkServicetemplateToServicetemplate) TableName() string {
	return "tbl_lnkServicetemplateToServicetemplate"
}

// LnkServicetemplateToContact maps to tbl_lnkServicetemplateToContact.
type LnkServicetemplateToContact struct {
	ServicetemplateID uint `gorm:"column:idMaster;not null;index:idx_master" json:"servicetemplate_id"`
	ContactID         uint `gorm:"column:idSlave;not null" json:"contact_id"`
	Sequence          uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkServicetemplateToContact) TableName() string { return "tbl_lnkServicetemplateToContact" }

// LnkServicetemplateToContactgroup maps to tbl_lnkServicetemplateToContactgroup.
type LnkServicetemplateToContactgroup struct {
	ServicetemplateID uint `gorm:"column:idMaster;not null;index:idx_master" json:"servicetemplate_id"`
	ContactgroupID    uint `gorm:"column:idSlave;not null" json:"contactgroup_id"`
	Sequence          uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkServicetemplateToContactgroup) TableName() string {
	return "tbl_lnkServicetemplateToContactgroup"
}

// --- Contact links ---

// LnkContactToContactgroup maps to tbl_lnkContactToContactgroup.
type LnkContactToContactgroup struct {
	ContactID      uint `gorm:"column:idMaster;not null;index:idx_master" json:"contact_id"`
	ContactgroupID uint `gorm:"column:idSlave;not null" json:"contactgroup_id"`
	Sequence       uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkContactToContactgroup) TableName() string { return "tbl_lnkContactToContactgroup" }

// LnkContactToHostNotificationCommand maps to tbl_lnkContactToCommandHost.
type LnkContactToCommandHost struct {
	ContactID uint `gorm:"column:idMaster;not null;index:idx_master" json:"contact_id"`
	CommandID uint `gorm:"column:idSlave;not null" json:"command_id"`
	Sequence  uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkContactToCommandHost) TableName() string { return "tbl_lnkContactToCommandHost" }

// LnkContactToCommandService maps to tbl_lnkContactToCommandService.
type LnkContactToCommandService struct {
	ContactID uint `gorm:"column:idMaster;not null;index:idx_master" json:"contact_id"`
	CommandID uint `gorm:"column:idSlave;not null" json:"command_id"`
	Sequence  uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkContactToCommandService) TableName() string { return "tbl_lnkContactToCommandService" }

// LnkContactToContacttemplate maps to tbl_lnkContactToContacttemplate.
type LnkContactToContacttemplate struct {
	ContactID         uint `gorm:"column:idMaster;not null;index:idx_master" json:"contact_id"`
	ContacttemplateID uint `gorm:"column:idSlave;not null" json:"contacttemplate_id"`
	Sequence          uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkContactToContacttemplate) TableName() string { return "tbl_lnkContactToContacttemplate" }

// --- Contacttemplate links ---

// LnkContacttemplateToContactgroup maps to tbl_lnkContacttemplateToContactgroup.
type LnkContacttemplateToContactgroup struct {
	ContacttemplateID uint `gorm:"column:idMaster;not null;index:idx_master" json:"contacttemplate_id"`
	ContactgroupID    uint `gorm:"column:idSlave;not null" json:"contactgroup_id"`
	Sequence          uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkContacttemplateToContactgroup) TableName() string {
	return "tbl_lnkContacttemplateToContactgroup"
}

// LnkContacttemplateToCommandHost maps to tbl_lnkContacttemplateToCommandHost.
type LnkContacttemplateToCommandHost struct {
	ContacttemplateID uint `gorm:"column:idMaster;not null;index:idx_master" json:"contacttemplate_id"`
	CommandID         uint `gorm:"column:idSlave;not null" json:"command_id"`
	Sequence          uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkContacttemplateToCommandHost) TableName() string {
	return "tbl_lnkContacttemplateToCommandHost"
}

// LnkContacttemplateToCommandService maps to tbl_lnkContacttemplateToCommandService.
type LnkContacttemplateToCommandService struct {
	ContacttemplateID uint `gorm:"column:idMaster;not null;index:idx_master" json:"contacttemplate_id"`
	CommandID         uint `gorm:"column:idSlave;not null" json:"command_id"`
	Sequence          uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkContacttemplateToCommandService) TableName() string {
	return "tbl_lnkContacttemplateToCommandService"
}

// --- Variabledefinition link tables (6 total) ---

// LnkHostToVariabledefinition maps to tbl_lnkHostToVariabledefinition.
type LnkHostToVariabledefinition struct {
	HostID              uint   `gorm:"column:idMaster;not null;index:idx_master" json:"host_id"`
	VarName             string `gorm:"column:varname;size:255;not null" json:"varname"`
	VarValue            string `gorm:"column:varvalue;size:255;not null" json:"varvalue"`
	Sequence            uint   `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkHostToVariabledefinition) TableName() string { return "tbl_lnkHostToVariabledefinition" }

// LnkServiceToVariabledefinition maps to tbl_lnkServiceToVariabledefinition.
type LnkServiceToVariabledefinition struct {
	ServiceID uint   `gorm:"column:idMaster;not null;index:idx_master" json:"service_id"`
	VarName   string `gorm:"column:varname;size:255;not null" json:"varname"`
	VarValue  string `gorm:"column:varvalue;size:255;not null" json:"varvalue"`
	Sequence  uint   `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkServiceToVariabledefinition) TableName() string {
	return "tbl_lnkServiceToVariabledefinition"
}

// LnkHosttemplateToVariabledefinition maps to tbl_lnkHosttemplateToVariabledefinition.
type LnkHosttemplateToVariabledefinition struct {
	HosttemplateID uint   `gorm:"column:idMaster;not null;index:idx_master" json:"hosttemplate_id"`
	VarName        string `gorm:"column:varname;size:255;not null" json:"varname"`
	VarValue       string `gorm:"column:varvalue;size:255;not null" json:"varvalue"`
	Sequence       uint   `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkHosttemplateToVariabledefinition) TableName() string {
	return "tbl_lnkHosttemplateToVariabledefinition"
}

// LnkServicetemplateToVariabledefinition maps to tbl_lnkServicetemplateToVariabledefinition.
type LnkServicetemplateToVariabledefinition struct {
	ServicetemplateID uint   `gorm:"column:idMaster;not null;index:idx_master" json:"servicetemplate_id"`
	VarName           string `gorm:"column:varname;size:255;not null" json:"varname"`
	VarValue          string `gorm:"column:varvalue;size:255;not null" json:"varvalue"`
	Sequence          uint   `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkServicetemplateToVariabledefinition) TableName() string {
	return "tbl_lnkServicetemplateToVariabledefinition"
}

// LnkContactToVariabledefinition maps to tbl_lnkContactToVariabledefinition.
type LnkContactToVariabledefinition struct {
	ContactID uint   `gorm:"column:idMaster;not null;index:idx_master" json:"contact_id"`
	VarName   string `gorm:"column:varname;size:255;not null" json:"varname"`
	VarValue  string `gorm:"column:varvalue;size:255;not null" json:"varvalue"`
	Sequence  uint   `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkContactToVariabledefinition) TableName() string {
	return "tbl_lnkContactToVariabledefinition"
}

// LnkContacttemplateToVariabledefinition maps to tbl_lnkContacttemplateToVariabledefinition.
type LnkContacttemplateToVariabledefinition struct {
	ContacttemplateID uint   `gorm:"column:idMaster;not null;index:idx_master" json:"contacttemplate_id"`
	VarName           string `gorm:"column:varname;size:255;not null" json:"varname"`
	VarValue          string `gorm:"column:varvalue;size:255;not null" json:"varvalue"`
	Sequence          uint   `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkContacttemplateToVariabledefinition) TableName() string {
	return "tbl_lnkContacttemplateToVariabledefinition"
}

// LnkHostgroupToHostgroup maps to tbl_lnkHostgroupToHostgroup.
type LnkHostgroupToHostgroup struct {
	HostgroupID       uint `gorm:"column:idMaster;not null;index:idx_master" json:"hostgroup_id"`
	MemberHostgroupID uint `gorm:"column:idSlave;not null" json:"member_hostgroup_id"`
	Sequence          uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkHostgroupToHostgroup) TableName() string { return "tbl_lnkHostgroupToHostgroup" }

// LnkHostgroupToHost maps to tbl_lnkHostgroupToHost.
type LnkHostgroupToHost struct {
	HostgroupID uint `gorm:"column:idMaster;not null;index:idx_master" json:"hostgroup_id"`
	HostID      uint `gorm:"column:idSlave;not null" json:"host_id"`
	Sequence    uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkHostgroupToHost) TableName() string { return "tbl_lnkHostgroupToHost" }

// LnkServicegroupToServicegroup maps to tbl_lnkServicegroupToServicegroup.
type LnkServicegroupToServicegroup struct {
	ServicegroupID       uint `gorm:"column:idMaster;not null;index:idx_master" json:"servicegroup_id"`
	MemberServicegroupID uint `gorm:"column:idSlave;not null" json:"member_servicegroup_id"`
	Sequence             uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkServicegroupToServicegroup) TableName() string { return "tbl_lnkServicegroupToServicegroup" }

// LnkServicegroupToService maps to tbl_lnkServicegroupToService.
type LnkServicegroupToService struct {
	ServicegroupID uint `gorm:"column:idMaster;not null;index:idx_master" json:"servicegroup_id"`
	ServiceID      uint `gorm:"column:idSlave;not null" json:"service_id"`
	Sequence       uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkServicegroupToService) TableName() string { return "tbl_lnkServicegroupToService" }

// LnkContactgroupToContact maps to tbl_lnkContactgroupToContact.
type LnkContactgroupToContact struct {
	ContactgroupID uint `gorm:"column:idMaster;not null;index:idx_master" json:"contactgroup_id"`
	ContactID      uint `gorm:"column:idSlave;not null" json:"contact_id"`
	Sequence       uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkContactgroupToContact) TableName() string { return "tbl_lnkContactgroupToContact" }

// LnkContactgroupToContactgroup maps to tbl_lnkContactgroupToContactgroup.
type LnkContactgroupToContactgroup struct {
	ContactgroupID       uint `gorm:"column:idMaster;not null;index:idx_master" json:"contactgroup_id"`
	MemberContactgroupID uint `gorm:"column:idSlave;not null" json:"member_contactgroup_id"`
	Sequence             uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkContactgroupToContactgroup) TableName() string {
	return "tbl_lnkContactgroupToContactgroup"
}

// LnkTimeperiodToTimeperiod maps to tbl_lnkTimeperiodToTimeperiod (excluded periods).
type LnkTimeperiodToTimeperiod struct {
	TimeperiodID         uint `gorm:"column:idMaster;not null;index:idx_master" json:"timeperiod_id"`
	ExcludedTimeperiodID uint `gorm:"column:idSlave;not null" json:"excluded_timeperiod_id"`
	Sequence             uint `gorm:"column:idSort;not null;default:0" json:"sequence"`
}

func (LnkTimeperiodToTimeperiod) TableName() string { return "tbl_lnkTimeperiodToTimeperiod" }
