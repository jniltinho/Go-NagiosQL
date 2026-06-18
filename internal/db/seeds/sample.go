package seeds

import (
	"fmt"
	"log"
	"time"

	"go-nagiosql/internal/models"
	"gorm.io/gorm"
)

// SeedSample loads the standard Nagios sample objects equivalent to
// import_nagios_sample.sql: 24 commands, 5 timeperiods with time definitions,
// 5 host templates, 2 service templates, 4 hostgroups, 1 contactgroup,
// 1 contact (nagiosadmin), 4 hosts, and 21 services with all link rows.
func SeedSample(db *gorm.DB) error {
	now := time.Now()

	if err := seedSampleCommands(db, now); err != nil {
		return err
	}
	if err := seedSampleTimeperiods(db, now); err != nil {
		return err
	}
	if err := seedSampleHosttemplates(db, now); err != nil {
		return err
	}
	if err := seedSampleServicetemplates(db, now); err != nil {
		return err
	}
	if err := seedSampleContacttemplate(db, now); err != nil {
		return err
	}
	if err := seedSampleContactgroup(db, now); err != nil {
		return err
	}
	if err := seedSampleContact(db, now); err != nil {
		return err
	}
	if err := seedSampleHostgroups(db, now); err != nil {
		return err
	}
	if err := seedSampleHosts(db, now); err != nil {
		return err
	}
	if err := seedSampleServices(db, now); err != nil {
		return err
	}

	log.Println("sample seed complete")
	return nil
}

// p is a convenience helper to get a pointer to an int literal.
func p(n int) *int { return &n }

func seedSampleCommands(db *gorm.DB, now time.Time) error {
	commands := []models.Command{
		{CommandName: "check-host-alive", CommandLine: "$USER1$/check_ping -H $HOSTADDRESS$ -w 3000.0,80% -c 5000.0,100% -p 5", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "check_ping", CommandLine: "$USER1$/check_ping -H $HOSTADDRESS$ -w $ARG1$ -c $ARG2$ -p 5", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "check_local_disk", CommandLine: "$USER1$/check_disk -w $ARG1$ -c $ARG2$ -p $ARG3$", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "check_local_load", CommandLine: "$USER1$/check_load -w $ARG1$ -c $ARG2$", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "check_local_procs", CommandLine: "$USER1$/check_procs -w $ARG1$ -c $ARG2$ -s $ARG3$", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "check_local_users", CommandLine: "$USER1$/check_users -w $ARG1$ -c $ARG2$", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "check_local_swap", CommandLine: "$USER1$/check_swap -w $ARG1$ -c $ARG2$", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "check_local_mrtgtraf", CommandLine: "$USER1$/check_mrtgtraf -F $ARG1$ -a $ARG2$ -w $ARG3$ -c $ARG4$ -e $ARG5$", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "check_ftp", CommandLine: "$USER1$/check_ftp -H $HOSTADDRESS$ $ARG1$", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "check_hpjd", CommandLine: "$USER1$/check_hpjd -H $HOSTADDRESS$ -C $ARG1$", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "check_snmp", CommandLine: "$USER1$/check_snmp -H $HOSTADDRESS$ -o $ARG1$ $ARG2$", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "check_http", CommandLine: "$USER1$/check_http -I $HOSTADDRESS$ $ARG1$", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "check_ssh", CommandLine: "$USER1$/check_ssh $ARG1$ $HOSTADDRESS$", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "check_dhcp", CommandLine: "$USER1$/check_dhcp $ARG1$", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "check_pop", CommandLine: "$USER1$/check_pop -H $HOSTADDRESS$ $ARG1$", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "check_imap", CommandLine: "$USER1$/check_imap -H $HOSTADDRESS$ $ARG1$", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "check_smtp", CommandLine: "$USER1$/check_smtp -H $HOSTADDRESS$ $ARG1$", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "check_tcp", CommandLine: "$USER1$/check_tcp -H $HOSTADDRESS$ -p $ARG1$ $ARG2$", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "check_udp", CommandLine: "$USER1$/check_udp -H $HOSTADDRESS$ -p $ARG1$ $ARG2$", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "check_nt", CommandLine: "$USER1$/check_nt -H $HOSTADDRESS$ -p 12489 -v $ARG1$ $ARG2$", CommandType: 0, Register: "1", Active: "1", LastModified: now},
		{CommandName: "notify-host-by-email", CommandLine: "/usr/bin/printf \"%b\" \"***** Nagios *****\\n\\nNotification Type: $NOTIFICATIONTYPE$\\nHost: $HOSTNAME$\\nState: $HOSTSTATE$\\nAddress: $HOSTADDRESS$\\nInfo: $HOSTOUTPUT$\\n\\nDate/Time: $LONGDATETIME$\\n\" | /usr/bin/mail -s \"** $NOTIFICATIONTYPE$ Host Alert: $HOSTNAME$ is $HOSTSTATE$ **\" $CONTACTEMAIL$", CommandType: 1, Register: "1", Active: "1", LastModified: now},
		{CommandName: "notify-service-by-email", CommandLine: "/usr/bin/printf \"%b\" \"***** Nagios *****\\n\\nNotification Type: $NOTIFICATIONTYPE$\\n\\nService: $SERVICEDESC$\\nHost: $HOSTALIAS$\\nAddress: $HOSTADDRESS$\\nState: $SERVICESTATE$\\n\\nDate/Time: $LONGDATETIME$\\n\\nAdditional Info:\\n\\n$SERVICEOUTPUT$\\n\" | /usr/bin/mail -s \"** $NOTIFICATIONTYPE$ Service Alert: $HOSTALIAS$/$SERVICEDESC$ is $SERVICESTATE$ **\" $CONTACTEMAIL$", CommandType: 1, Register: "1", Active: "1", LastModified: now},
		{CommandName: "process-host-perfdata", CommandLine: "/usr/bin/printf \"%b\" \"$LASTHOSTCHECK$\\t$HOSTNAME$\\t$HOSTSTATE$\\t$HOSTATTEMPT$\\t$HOSTSTATETYPE$\\t$HOSTEXECUTIONTIME$\\t$HOSTOUTPUT$\\t$HOSTPERFDATA$\\n\" >> /var/lib/nagios4/host-perfdata.out", CommandType: 1, Register: "1", Active: "1", LastModified: now},
		{CommandName: "process-service-perfdata", CommandLine: "/usr/bin/printf \"%b\" \"$LASTSERVICECHECK$\\t$HOSTNAME$\\t$SERVICEDESC$\\t$SERVICESTATE$\\t$SERVICEATTEMPT$\\t$SERVICESTATETYPE$\\t$SERVICEEXECUTIONTIME$\\t$SERVICELATENCY$\\t$SERVICEOUTPUT$\\t$SERVICEPERFDATA$\\n\" >> /var/lib/nagios4/service-perfdata.out", CommandType: 1, Register: "1", Active: "1", LastModified: now},
	}

	for i := range commands {
		cmd := &commands[i]
		if err := db.Where(models.Command{CommandName: cmd.CommandName}).FirstOrCreate(cmd).Error; err != nil {
			return fmt.Errorf("seed command %q: %w", cmd.CommandName, err)
		}
	}
	log.Printf("seeded %d sample commands", len(commands))
	return nil
}

func seedSampleTimeperiods(db *gorm.DB, now time.Time) error {
	type tpDef struct {
		tp   models.Timeperiod
		defs []models.Timedefinition
	}

	timeperiods := []tpDef{
		{
			tp: models.Timeperiod{TimeperiodName: "24x7", Alias: "24 Hours A Day, 7 Days A Week", Register: "1", Active: "1", LastModified: now},
			defs: []models.Timedefinition{
				{Definition: "sunday", Range: "00:00-24:00", LastModified: now},
				{Definition: "monday", Range: "00:00-24:00", LastModified: now},
				{Definition: "tuesday", Range: "00:00-24:00", LastModified: now},
				{Definition: "wednesday", Range: "00:00-24:00", LastModified: now},
				{Definition: "thursday", Range: "00:00-24:00", LastModified: now},
				{Definition: "friday", Range: "00:00-24:00", LastModified: now},
				{Definition: "saturday", Range: "00:00-24:00", LastModified: now},
			},
		},
		{
			tp: models.Timeperiod{TimeperiodName: "workhours", Alias: "Standard Work Hours", Register: "1", Active: "1", LastModified: now},
			defs: []models.Timedefinition{
				{Definition: "monday", Range: "09:00-17:00", LastModified: now},
				{Definition: "tuesday", Range: "09:00-17:00", LastModified: now},
				{Definition: "wednesday", Range: "09:00-17:00", LastModified: now},
				{Definition: "thursday", Range: "09:00-17:00", LastModified: now},
				{Definition: "friday", Range: "09:00-17:00", LastModified: now},
			},
		},
		{
			tp: models.Timeperiod{TimeperiodName: "nonworkhours", Alias: "Non-Work Hours", Register: "1", Active: "1", LastModified: now},
			defs: []models.Timedefinition{
				{Definition: "sunday", Range: "00:00-24:00", LastModified: now},
				{Definition: "monday", Range: "00:00-09:00,17:00-24:00", LastModified: now},
				{Definition: "tuesday", Range: "00:00-09:00,17:00-24:00", LastModified: now},
				{Definition: "wednesday", Range: "00:00-09:00,17:00-24:00", LastModified: now},
				{Definition: "thursday", Range: "00:00-09:00,17:00-24:00", LastModified: now},
				{Definition: "friday", Range: "00:00-09:00,17:00-24:00", LastModified: now},
				{Definition: "saturday", Range: "00:00-24:00", LastModified: now},
			},
		},
		{
			tp:   models.Timeperiod{TimeperiodName: "us-holidays", Alias: "U.S. Holidays", Register: "1", Active: "1", LastModified: now},
			defs: []models.Timedefinition{},
		},
		{
			tp:   models.Timeperiod{TimeperiodName: "24x7_sans_holidays", Alias: "24x7 Sans Holidays", Exclude: 1, Register: "1", Active: "1", LastModified: now},
			defs: []models.Timedefinition{
				{Definition: "sunday", Range: "00:00-24:00", LastModified: now},
				{Definition: "monday", Range: "00:00-24:00", LastModified: now},
				{Definition: "tuesday", Range: "00:00-24:00", LastModified: now},
				{Definition: "wednesday", Range: "00:00-24:00", LastModified: now},
				{Definition: "thursday", Range: "00:00-24:00", LastModified: now},
				{Definition: "friday", Range: "00:00-24:00", LastModified: now},
				{Definition: "saturday", Range: "00:00-24:00", LastModified: now},
			},
		},
	}

	for i := range timeperiods {
		entry := &timeperiods[i]
		tp := &entry.tp
		if err := db.Where(models.Timeperiod{TimeperiodName: tp.TimeperiodName}).FirstOrCreate(tp).Error; err != nil {
			return fmt.Errorf("seed timeperiod %q: %w", tp.TimeperiodName, err)
		}
		if tp.TimeperiodName == "24x7_sans_holidays" {
			var usHolidays models.Timeperiod
			if err := db.Where("timeperiod_name = ?", "us-holidays").First(&usHolidays).Error; err == nil {
				link := LnkTimeperiodToTimeperiodSeed{MasterID: tp.ID, SlaveID: usHolidays.ID}
				db.Exec("INSERT IGNORE INTO tbl_lnkTimeperiodToTimeperiod (idMaster, idSlave, idSort) VALUES (?,?,0)", link.MasterID, link.SlaveID)
			}
		}
		for _, def := range entry.defs {
			def.TipID = tp.ID
			db.Where(models.Timedefinition{TipID: tp.ID, Definition: def.Definition}).FirstOrCreate(&def) //nolint:errcheck
		}
	}
	log.Printf("seeded %d sample timeperiods", len(timeperiods))
	return nil
}

// LnkTimeperiodToTimeperiodSeed is a transient helper for raw-SQL insertion.
type LnkTimeperiodToTimeperiodSeed struct{ MasterID, SlaveID uint }

func seedSampleHosttemplates(db *gorm.DB, now time.Time) error {
	// generic-host must be created first (it's the parent of the others).
	templates := []models.Hosttemplate{
		{TemplateName: "generic-host", Alias: "Generic Host Template", MaxCheckAttempts: p(10), CheckInterval: p(5), RetryInterval: p(1), NotificationInterval: p(60), NotificationOptions: "d,u,r", Active: "1", LastModified: now},
		{TemplateName: "linux-server", Alias: "Linux Server Template", UseTemplate: 1, Active: "1", LastModified: now},
		{TemplateName: "windows-server", Alias: "Windows Server Template", UseTemplate: 1, Active: "1", LastModified: now},
		{TemplateName: "generic-switch", Alias: "Generic Switch Template", UseTemplate: 1, Active: "1", LastModified: now},
		{TemplateName: "generic-printer", Alias: "Generic Printer Template", UseTemplate: 1, Active: "1", LastModified: now},
	}
	for i := range templates {
		tmpl := &templates[i]
		if err := db.Where(models.Hosttemplate{TemplateName: tmpl.TemplateName}).FirstOrCreate(tmpl).Error; err != nil {
			return fmt.Errorf("seed hosttemplate %q: %w", tmpl.TemplateName, err)
		}
	}

	// Link child templates → generic-host parent.
	var genericHost models.Hosttemplate
	if err := db.Where("template_name = ?", "generic-host").First(&genericHost).Error; err != nil {
		return fmt.Errorf("resolve generic-host: %w", err)
	}
	// Link check_period and notification_period to 24x7 (ID queried, not hardcoded).
	var tp24x7 models.Timeperiod
	if db.Where("timeperiod_name = ?", "24x7").First(&tp24x7).Error == nil {
		db.Exec("UPDATE tbl_hosttemplate SET check_period=?, notification_period=? WHERE template_name='generic-host' AND check_period=0", tp24x7.ID, tp24x7.ID)
	}
	for _, childName := range []string{"linux-server", "windows-server", "generic-switch", "generic-printer"} {
		var child models.Hosttemplate
		if db.Where("template_name = ?", childName).First(&child).Error == nil {
			db.Exec("INSERT IGNORE INTO tbl_lnkHosttemplateToHosttemplate (idMaster, idSlave, idSort) VALUES (?,?,0)", child.ID, genericHost.ID)
		}
	}

	log.Printf("seeded %d host templates", len(templates))
	return nil
}

func seedSampleServicetemplates(db *gorm.DB, now time.Time) error {
	templates := []models.Servicetemplate{
		{TemplateName: "generic-service", ServiceDescription: "Generic Service Template", MaxCheckAttempts: p(3), CheckInterval: p(10), RetryInterval: p(2), NotificationInterval: p(60), NotificationOptions: "w,u,c,r", Active: "1", LastModified: now},
		{TemplateName: "local-service", ServiceDescription: "Local Service Template", UseTemplate: 1, Active: "1", LastModified: now},
	}
	for i := range templates {
		tmpl := &templates[i]
		if err := db.Where(models.Servicetemplate{TemplateName: tmpl.TemplateName}).FirstOrCreate(tmpl).Error; err != nil {
			return fmt.Errorf("seed servicetemplate %q: %w", tmpl.TemplateName, err)
		}
	}

	// Link check_period and notification_period on generic-service → 24x7.
	var tp24x7 models.Timeperiod
	if db.Where("timeperiod_name = ?", "24x7").First(&tp24x7).Error == nil {
		db.Exec("UPDATE tbl_servicetemplate SET check_period=?, notification_period=? WHERE template_name='generic-service' AND check_period=0", tp24x7.ID, tp24x7.ID)
	}

	// Link local-service → generic-service parent.
	var genericSvc, localSvc models.Servicetemplate
	if db.Where("template_name = ?", "generic-service").First(&genericSvc).Error == nil {
		if db.Where("template_name = ?", "local-service").First(&localSvc).Error == nil {
			db.Exec("INSERT IGNORE INTO tbl_lnkServicetemplateToServicetemplate (idMaster, idSlave, idSort) VALUES (?,?,0)", localSvc.ID, genericSvc.ID)
		}
	}

	log.Printf("seeded %d service templates", len(templates))
	return nil
}

func seedSampleContacttemplate(db *gorm.DB, now time.Time) error {
	ct := models.Contacttemplate{
		TemplateName:               "generic-contact",
		HostNotificationOptions:    "d,u,r",
		ServiceNotificationOptions: "w,u,c,r",
		HostNotificationCommands:   1,
		ServiceNotificationCommands: 1,
		Active:                     "1",
		LastModified:               now,
	}
	if err := db.Where(models.Contacttemplate{TemplateName: "generic-contact"}).FirstOrCreate(&ct).Error; err != nil {
		return fmt.Errorf("seed contacttemplate generic-contact: %w", err)
	}

	// Link notification periods → 24x7.
	var tp24x7 models.Timeperiod
	if db.Where("timeperiod_name = ?", "24x7").First(&tp24x7).Error == nil {
		db.Exec("UPDATE tbl_contacttemplate SET host_notification_period=?, service_notification_period=? WHERE template_name='generic-contact' AND host_notification_period=0", tp24x7.ID, tp24x7.ID)
	}

	// Link notification commands.
	var hostCmd, svcCmd models.Command
	if db.Where("command_name = ?", "notify-host-by-email").First(&hostCmd).Error == nil {
		db.Exec("INSERT IGNORE INTO tbl_lnkContacttemplateToCommandHost (idMaster, idSlave, idSort) VALUES (?,?,0)", ct.ID, hostCmd.ID)
	}
	if db.Where("command_name = ?", "notify-service-by-email").First(&svcCmd).Error == nil {
		db.Exec("INSERT IGNORE INTO tbl_lnkContacttemplateToCommandService (idMaster, idSlave, idSort) VALUES (?,?,0)", ct.ID, svcCmd.ID)
	}

	log.Printf("seeded contact template generic-contact")
	return nil
}

func seedSampleContactgroup(db *gorm.DB, now time.Time) error {
	cg := models.Contactgroup{
		ContactgroupName: "admins",
		Alias:            "Nagios Administrators",
		Members:          1,
		Register:         "1",
		Active:           "1",
		LastModified:     now,
	}
	if err := db.Where(models.Contactgroup{ContactgroupName: "admins"}).FirstOrCreate(&cg).Error; err != nil {
		return fmt.Errorf("seed contactgroup admins: %w", err)
	}
	return nil
}

func seedSampleContact(db *gorm.DB, now time.Time) error {
	contact := models.Contact{
		ContactName:      "nagiosadmin",
		Alias:            "Nagios Admin",
		Email:            "nagios@localhost",
		UseTemplate:      1,
		Contactgroups:    1,
		Register:         "1",
		Active:           "1",
		LastModified:     now,
	}
	if err := db.Where(models.Contact{ContactName: "nagiosadmin"}).FirstOrCreate(&contact).Error; err != nil {
		return fmt.Errorf("seed contact nagiosadmin: %w", err)
	}

	// Link to generic-contact template.
	var ct models.Contacttemplate
	if db.Where("template_name = ?", "generic-contact").First(&ct).Error == nil {
		db.Exec("INSERT IGNORE INTO tbl_lnkContactToContacttemplate (idMaster, idSlave, idSort) VALUES (?,?,0)", contact.ID, ct.ID)
	}

	// Link contact to admins contactgroup.
	var cg models.Contactgroup
	if err := db.Where("contactgroup_name = ?", "admins").First(&cg).Error; err == nil {
		db.Exec("INSERT IGNORE INTO tbl_lnkContactToContactgroup (idMaster, idSlave, idSort) VALUES (?,?,0)", contact.ID, cg.ID)
		db.Exec("INSERT IGNORE INTO tbl_lnkContactgroupToContact (idMaster, idSlave, idSort) VALUES (?,?,0)", cg.ID, contact.ID)
	}
	return nil
}

func seedSampleHostgroups(db *gorm.DB, now time.Time) error {
	groups := []models.Hostgroup{
		{HostgroupName: "linux-servers", Alias: "Linux Servers", Members: 1, Register: "1", Active: "1", LastModified: now},
		{HostgroupName: "windows-servers", Alias: "Windows Servers", Members: 1, Register: "1", Active: "1", LastModified: now},
		{HostgroupName: "network-switches", Alias: "Network Switches", Members: 1, Register: "1", Active: "1", LastModified: now},
		{HostgroupName: "network-printers", Alias: "Network Printers", Members: 1, Register: "1", Active: "1", LastModified: now},
	}
	for i := range groups {
		g := &groups[i]
		if err := db.Where(models.Hostgroup{HostgroupName: g.HostgroupName}).FirstOrCreate(g).Error; err != nil {
			return fmt.Errorf("seed hostgroup %q: %w", g.HostgroupName, err)
		}
	}
	log.Printf("seeded %d hostgroups", len(groups))
	return nil
}

func seedSampleHosts(db *gorm.DB, now time.Time) error {
	type hostDef struct {
		h         models.Host
		template  string
		hostgroup string
	}
	hosts := []hostDef{
		{models.Host{HostName: "localhost", Alias: "localhost", DisplayName: "localhost", Address: "127.0.0.1", UseTemplate: 1, Hostgroups: 1, ContactGroups: 1, Active: "1", Register: "1", LastModified: now}, "linux-server", "linux-servers"},
		{models.Host{HostName: "winserver", Alias: "Windows Server", DisplayName: "winserver", Address: "192.168.1.2", UseTemplate: 1, Hostgroups: 1, ContactGroups: 1, Active: "1", Register: "1", LastModified: now}, "windows-server", "windows-servers"},
		{models.Host{HostName: "linksys-srw224p", Alias: "Linksys SRW224P Switch", DisplayName: "linksys-srw224p", Address: "192.168.1.253", UseTemplate: 1, Hostgroups: 1, ContactGroups: 1, Active: "1", Register: "1", LastModified: now}, "generic-switch", "network-switches"},
		{models.Host{HostName: "hplj2605dn", Alias: "HP LaserJet 2605dn", DisplayName: "hplj2605dn", Address: "192.168.1.30", UseTemplate: 1, Hostgroups: 1, ContactGroups: 1, Active: "1", Register: "1", LastModified: now}, "generic-printer", "network-printers"},
	}

	for i := range hosts {
		entry := &hosts[i]
		h := &entry.h
		if err := db.Where(models.Host{HostName: h.HostName}).FirstOrCreate(h).Error; err != nil {
			return fmt.Errorf("seed host %q: %w", h.HostName, err)
		}

		// Link host → hosttemplate.
		var tmpl models.Hosttemplate
		if db.Where("template_name = ?", entry.template).First(&tmpl).Error == nil {
			db.Exec("INSERT IGNORE INTO tbl_lnkHostToHosttemplate (idMaster, idSlave, idSort) VALUES (?,?,0)", h.ID, tmpl.ID)
		}

		// Link host → hostgroup (both directions).
		var hg models.Hostgroup
		if db.Where("hostgroup_name = ?", entry.hostgroup).First(&hg).Error == nil {
			db.Exec("INSERT IGNORE INTO tbl_lnkHostToHostgroup (idMaster, idSlave, idSort) VALUES (?,?,0)", h.ID, hg.ID)
			db.Exec("INSERT IGNORE INTO tbl_lnkHostgroupToHost (idMaster, idSlave, idSort) VALUES (?,?,0)", hg.ID, h.ID)
		}

		// Link host → admins contactgroup.
		var cg models.Contactgroup
		if db.Where("contactgroup_name = ?", "admins").First(&cg).Error == nil {
			db.Exec("INSERT IGNORE INTO tbl_lnkHostToContactgroup (idMaster, idSlave, idSort) VALUES (?,?,0)", h.ID, cg.ID)
		}
	}
	log.Printf("seeded %d sample hosts", len(hosts))
	return nil
}

func seedSampleServices(db *gorm.DB, now time.Time) error {
	type svcDef struct {
		s        models.Service
		hostName string
		template string
	}
	services := []svcDef{
		// localhost — local-service template
		{models.Service{ConfigName: "localhost", ServiceDescription: "PING", CheckCommand: "check_ping!100.0,20%!500.0,60%", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "localhost", "local-service"},
		{models.Service{ConfigName: "localhost", ServiceDescription: "Root Partition", CheckCommand: "check_local_disk!20%!10%!/", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "localhost", "local-service"},
		{models.Service{ConfigName: "localhost", ServiceDescription: "Current Users", CheckCommand: "check_local_users!20!50", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "localhost", "local-service"},
		{models.Service{ConfigName: "localhost", ServiceDescription: "Total Processes", CheckCommand: "check_local_procs!250!400!RSZDT", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "localhost", "local-service"},
		{models.Service{ConfigName: "localhost", ServiceDescription: "Current Load", CheckCommand: "check_local_load!5.0,4.0,3.0!10.0,6.0,4.0", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "localhost", "local-service"},
		{models.Service{ConfigName: "localhost", ServiceDescription: "Swap Usage", CheckCommand: "check_local_swap!20%!10%", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "localhost", "local-service"},
		{models.Service{ConfigName: "localhost", ServiceDescription: "SSH", CheckCommand: "check_ssh", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "localhost", "local-service"},
		{models.Service{ConfigName: "localhost", ServiceDescription: "HTTP", CheckCommand: "check_http", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "localhost", "local-service"},
		// winserver — generic-service template
		{models.Service{ConfigName: "winserver", ServiceDescription: "NSClient++ Version", CheckCommand: "check_nt!CLIENTVERSION", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "winserver", "generic-service"},
		{models.Service{ConfigName: "winserver", ServiceDescription: "Uptime", CheckCommand: "check_nt!UPTIME", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "winserver", "generic-service"},
		{models.Service{ConfigName: "winserver", ServiceDescription: "CPU Load", CheckCommand: "check_nt!CPULOAD!-l 5,80,90", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "winserver", "generic-service"},
		{models.Service{ConfigName: "winserver", ServiceDescription: "Memory Usage", CheckCommand: "check_nt!MEMUSE!-w 80 -c 90", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "winserver", "generic-service"},
		{models.Service{ConfigName: "winserver", ServiceDescription: "C:\\ Drive Space", CheckCommand: "check_nt!USEDDISKSPACE!-l c -w 80 -c 90", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "winserver", "generic-service"},
		{models.Service{ConfigName: "winserver", ServiceDescription: "W3SVC", CheckCommand: "check_nt!SERVICESTATE!-d SHOWALL -l W3SVC", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "winserver", "generic-service"},
		{models.Service{ConfigName: "winserver", ServiceDescription: "Explorer", CheckCommand: "check_nt!PROCSTATE!-d SHOWALL -l Explorer.exe", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "winserver", "generic-service"},
		// linksys-srw224p — generic-service template
		{models.Service{ConfigName: "linksys-srw224p", ServiceDescription: "PING", CheckCommand: "check_ping!200.0,20%!600.0,60%", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "linksys-srw224p", "generic-service"},
		{models.Service{ConfigName: "linksys-srw224p", ServiceDescription: "Port 1 Bandwidth Usage", CheckCommand: "check_local_mrtgtraf!/var/lib/mrtg/192.168.1.253_1.log!AVG!1000000,1000000!5000000,5000000!10", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "linksys-srw224p", "generic-service"},
		{models.Service{ConfigName: "linksys-srw224p", ServiceDescription: "Port 1 Link Status", CheckCommand: "check_snmp!ifOperStatus.1!public", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "linksys-srw224p", "generic-service"},
		{models.Service{ConfigName: "linksys-srw224p", ServiceDescription: "Uptime", CheckCommand: "check_snmp!sysUpTime.0!public", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "linksys-srw224p", "generic-service"},
		// hplj2605dn — generic-service template
		{models.Service{ConfigName: "hplj2605dn", ServiceDescription: "PING", CheckCommand: "check_ping!3000.0,80%!5000.0,100%", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "hplj2605dn", "generic-service"},
		{models.Service{ConfigName: "hplj2605dn", ServiceDescription: "Printer Status", CheckCommand: "check_hpjd!public", Active: "1", Register: "1", UseTemplate: 1, ContactGroups: 1, LastModified: now}, "hplj2605dn", "generic-service"},
	}

	for i := range services {
		entry := &services[i]
		svc := &entry.s
		if err := db.Where(models.Service{ConfigName: svc.ConfigName, ServiceDescription: svc.ServiceDescription}).FirstOrCreate(svc).Error; err != nil {
			return fmt.Errorf("seed service %q/%q: %w", svc.ConfigName, svc.ServiceDescription, err)
		}

		// Link service → host.
		var host models.Host
		if db.Where("host_name = ?", entry.hostName).First(&host).Error == nil {
			db.Exec("INSERT IGNORE INTO tbl_lnkServiceToHost (idMaster, idSlave, idSort) VALUES (?,?,0)", svc.ID, host.ID)
		}

		// Link service → servicetemplate.
		var tmpl models.Servicetemplate
		if db.Where("template_name = ?", entry.template).First(&tmpl).Error == nil {
			db.Exec("INSERT IGNORE INTO tbl_lnkServiceToServicetemplate (idMaster, idSlave, idSort) VALUES (?,?,0)", svc.ID, tmpl.ID)
		}

		// Link service → admins contactgroup.
		var cg models.Contactgroup
		if db.Where("contactgroup_name = ?", "admins").First(&cg).Error == nil {
			db.Exec("INSERT IGNORE INTO tbl_lnkServiceToContactgroup (idMaster, idSlave, idSort) VALUES (?,?,0)", svc.ID, cg.ID)
		}
	}
	log.Printf("seeded %d sample services", len(services))
	return nil
}
