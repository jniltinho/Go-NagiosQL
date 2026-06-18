package models

// Link tables for the 6 extended Nagios object types:
// Hostdependency, Hostescalation, Servicedependency, Serviceescalation,
// Hostextinfo, Serviceextinfo.
//
// Tables that reference a service_description string (not an FK int) use
// strSlave instead of idSlave — these are resolved via resolveStrSlave.

// ── Hostdependency ───────────────────────────────────────────────────────────

type LnkHostdependencyToHostDH struct {
	MasterID uint `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveID  uint `gorm:"column:idSlave;not null"`
}

func (LnkHostdependencyToHostDH) TableName() string { return "tbl_lnkHostdependencyToHost_DH" }

type LnkHostdependencyToHostgroupDH struct {
	MasterID uint `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveID  uint `gorm:"column:idSlave;not null"`
}

func (LnkHostdependencyToHostgroupDH) TableName() string {
	return "tbl_lnkHostdependencyToHostgroup_DH"
}

type LnkHostdependencyToHostH struct {
	MasterID uint `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveID  uint `gorm:"column:idSlave;not null"`
}

func (LnkHostdependencyToHostH) TableName() string { return "tbl_lnkHostdependencyToHost_H" }

type LnkHostdependencyToHostgroupH struct {
	MasterID uint `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveID  uint `gorm:"column:idSlave;not null"`
}

func (LnkHostdependencyToHostgroupH) TableName() string {
	return "tbl_lnkHostdependencyToHostgroup_H"
}

// ── Hostescalation ───────────────────────────────────────────────────────────

type LnkHostescalationToHost struct {
	MasterID uint `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveID  uint `gorm:"column:idSlave;not null"`
}

func (LnkHostescalationToHost) TableName() string { return "tbl_lnkHostescalationToHost" }

type LnkHostescalationToHostgroup struct {
	MasterID uint `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveID  uint `gorm:"column:idSlave;not null"`
}

func (LnkHostescalationToHostgroup) TableName() string { return "tbl_lnkHostescalationToHostgroup" }

type LnkHostescalationToContact struct {
	MasterID uint `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveID  uint `gorm:"column:idSlave;not null"`
}

func (LnkHostescalationToContact) TableName() string { return "tbl_lnkHostescalationToContact" }

type LnkHostescalationToContactgroup struct {
	MasterID uint `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveID  uint `gorm:"column:idSlave;not null"`
}

func (LnkHostescalationToContactgroup) TableName() string {
	return "tbl_lnkHostescalationToContactgroup"
}

// ── Servicedependency ────────────────────────────────────────────────────────

type LnkServicedependencyToHostDH struct {
	MasterID uint `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveID  uint `gorm:"column:idSlave;not null"`
}

func (LnkServicedependencyToHostDH) TableName() string {
	return "tbl_lnkServicedependencyToHost_DH"
}

type LnkServicedependencyToHostgroupDH struct {
	MasterID uint `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveID  uint `gorm:"column:idSlave;not null"`
}

func (LnkServicedependencyToHostgroupDH) TableName() string {
	return "tbl_lnkServicedependencyToHostgroup_DH"
}

type LnkServicedependencyToServiceDS struct {
	MasterID uint   `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveStr string `gorm:"column:strSlave;size:255;not null"`
}

func (LnkServicedependencyToServiceDS) TableName() string {
	return "tbl_lnkServicedependencyToService_DS"
}

type LnkServicedependencyToServicegroupDS struct {
	MasterID uint `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveID  uint `gorm:"column:idSlave;not null"`
}

func (LnkServicedependencyToServicegroupDS) TableName() string {
	return "tbl_lnkServicedependencyToServicegroup_DS"
}

type LnkServicedependencyToHostH struct {
	MasterID uint `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveID  uint `gorm:"column:idSlave;not null"`
}

func (LnkServicedependencyToHostH) TableName() string { return "tbl_lnkServicedependencyToHost_H" }

type LnkServicedependencyToHostgroupH struct {
	MasterID uint `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveID  uint `gorm:"column:idSlave;not null"`
}

func (LnkServicedependencyToHostgroupH) TableName() string {
	return "tbl_lnkServicedependencyToHostgroup_H"
}

type LnkServicedependencyToServiceS struct {
	MasterID uint   `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveStr string `gorm:"column:strSlave;size:255;not null"`
}

func (LnkServicedependencyToServiceS) TableName() string {
	return "tbl_lnkServicedependencyToService_S"
}

type LnkServicedependencyToServicegroupS struct {
	MasterID uint `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveID  uint `gorm:"column:idSlave;not null"`
}

func (LnkServicedependencyToServicegroupS) TableName() string {
	return "tbl_lnkServicedependencyToServicegroup_S"
}

// ── Serviceescalation ────────────────────────────────────────────────────────

type LnkServiceescalationToHost struct {
	MasterID uint `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveID  uint `gorm:"column:idSlave;not null"`
}

func (LnkServiceescalationToHost) TableName() string { return "tbl_lnkServiceescalationToHost" }

type LnkServiceescalationToHostgroup struct {
	MasterID uint `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveID  uint `gorm:"column:idSlave;not null"`
}

func (LnkServiceescalationToHostgroup) TableName() string {
	return "tbl_lnkServiceescalationToHostgroup"
}

type LnkServiceescalationToService struct {
	MasterID uint   `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveStr string `gorm:"column:strSlave;size:255;not null"`
}

func (LnkServiceescalationToService) TableName() string { return "tbl_lnkServiceescalationToService" }

type LnkServiceescalationToServicegroup struct {
	MasterID uint `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveID  uint `gorm:"column:idSlave;not null"`
}

func (LnkServiceescalationToServicegroup) TableName() string {
	return "tbl_lnkServiceescalationToServicegroup"
}

type LnkServiceescalationToContact struct {
	MasterID uint `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveID  uint `gorm:"column:idSlave;not null"`
}

func (LnkServiceescalationToContact) TableName() string { return "tbl_lnkServiceescalationToContact" }

type LnkServiceescalationToContactgroup struct {
	MasterID uint `gorm:"column:idMaster;not null;index:idx_master"`
	SlaveID  uint `gorm:"column:idSlave;not null"`
}

func (LnkServiceescalationToContactgroup) TableName() string {
	return "tbl_lnkServiceescalationToContactgroup"
}
