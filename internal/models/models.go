package models

// AllModels returns a slice of all GORM model instances.
// Pass this to gorm.AutoMigrate to create or update every table.
// Order matters: tables with no foreign keys first (MyISAM does not enforce FKs,
// but it's good practice to migrate parent tables before link tables).
func AllModels() []any {
	return []any{
		// System / meta tables
		&Datadomain{},
		&Configtarget{},
		&Settings{},
		&User{},
		&Logbook{},

		// Core Nagios object tables
		&Variabledefinition{},
		&Command{},
		&Timeperiod{},
		&Timedefinition{},
		&Contact{},
		&Contacttemplate{},
		&Contactgroup{},
		&Host{},
		&Hosttemplate{},
		&Hostgroup{},
		&Service{},
		&Servicetemplate{},
		&Servicegroup{},

		// Link tables — host
		&LnkHostToHost{},
		&LnkHostToHosttemplate{},
		&LnkHostToHostgroup{},
		&LnkHostToContact{},
		&LnkHostToContactgroup{},
		&LnkHostToVariabledefinition{},

		// Link tables — service
		&LnkServiceToHost{},
		&LnkServiceToHostgroup{},
		&LnkServiceToServicetemplate{},
		&LnkServiceToServicegroup{},
		&LnkServiceToContact{},
		&LnkServiceToContactgroup{},
		&LnkServiceToVariabledefinition{},

		// Link tables — hosttemplate
		&LnkHosttemplateToHosttemplate{},
		&LnkHosttemplateToHostgroup{},
		&LnkHosttemplateToContact{},
		&LnkHosttemplateToContactgroup{},
		&LnkHosttemplateToVariabledefinition{},

		// Link tables — servicetemplate
		&LnkServicetemplateToServicetemplate{},
		&LnkServicetemplateToContact{},
		&LnkServicetemplateToContactgroup{},
		&LnkServicetemplateToVariabledefinition{},

		// Link tables — contact
		&LnkContactToContacttemplate{},
		&LnkContactToContactgroup{},
		&LnkContactToCommandHost{},
		&LnkContactToCommandService{},
		&LnkContactToVariabledefinition{},

		// Link tables — contacttemplate
		&LnkContacttemplateToContactgroup{},
		&LnkContacttemplateToCommandHost{},
		&LnkContacttemplateToCommandService{},
		&LnkContacttemplateToVariabledefinition{},

		// Link tables — groups
		&LnkHostgroupToHost{},
		&LnkHostgroupToHostgroup{},
		&LnkServicegroupToService{},
		&LnkServicegroupToServicegroup{},
		&LnkContactgroupToContact{},
		&LnkContactgroupToContactgroup{},

		// Link tables — timeperiod
		&LnkTimeperiodToTimeperiod{},

		// Extended object tables
		&Hostdependency{},
		&Hostescalation{},
		&Hostextinfo{},
		&Servicedependency{},
		&Serviceescalation{},
		&Serviceextinfo{},

		// Link tables — hostdependency
		&LnkHostdependencyToHostDH{},
		&LnkHostdependencyToHostgroupDH{},
		&LnkHostdependencyToHostH{},
		&LnkHostdependencyToHostgroupH{},

		// Link tables — hostescalation
		&LnkHostescalationToHost{},
		&LnkHostescalationToHostgroup{},
		&LnkHostescalationToContact{},
		&LnkHostescalationToContactgroup{},

		// Link tables — servicedependency
		&LnkServicedependencyToHostDH{},
		&LnkServicedependencyToHostgroupDH{},
		&LnkServicedependencyToServiceDS{},
		&LnkServicedependencyToServicegroupDS{},
		&LnkServicedependencyToHostH{},
		&LnkServicedependencyToHostgroupH{},
		&LnkServicedependencyToServiceS{},
		&LnkServicedependencyToServicegroupS{},

		// Link tables — serviceescalation
		&LnkServiceescalationToHost{},
		&LnkServiceescalationToHostgroup{},
		&LnkServiceescalationToService{},
		&LnkServiceescalationToServicegroup{},
		&LnkServiceescalationToContact{},
		&LnkServiceescalationToContactgroup{},
	}
}
