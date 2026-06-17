// Package nagconfig generates Nagios .cfg files from the database,
// matching the output format of the PHP NagiosQL application.
package nagconfig

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jniltinho/go-nagiosql/internal/models"
	"gorm.io/gorm"
)

// Generator writes Nagios configuration files from the database.
type Generator struct {
	db         *gorm.DB
	hostDir    string
	serviceDir string
	backupDir  string
}

// New creates a Generator using explicit directory paths.
func New(db *gorm.DB, hostDir, serviceDir, backupDir string) *Generator {
	return &Generator{db: db, hostDir: hostDir, serviceDir: serviceDir, backupDir: backupDir}
}

// WriteHost renders and writes a single host .cfg file.
func (g *Generator) WriteHost(hostID uint) error {
	var host models.Host
	if err := g.db.First(&host, hostID).Error; err != nil {
		return fmt.Errorf("host %d not found: %w", hostID, err)
	}
	content, err := g.renderHost(&host)
	if err != nil {
		return err
	}
	filename := sanitizeName(host.HostName) + ".cfg"
	return g.writeWithBackup(filepath.Join(g.hostDir, filename), content)
}

// WriteServiceGroup writes one .cfg file containing all services with the given config_name.
func (g *Generator) WriteServiceGroup(configName string) error {
	var svcs []models.Service
	if err := g.db.Where("config_name = ? AND active = '1'", configName).
		Order("service_description").Find(&svcs).Error; err != nil {
		return fmt.Errorf("loading services %q: %w", configName, err)
	}
	if len(svcs) == 0 {
		return nil
	}
	var buf bytes.Buffer
	for _, svc := range svcs {
		chunk, err := g.renderService(&svc)
		if err != nil {
			return err
		}
		buf.Write(chunk)
		buf.WriteByte('\n')
	}
	filename := sanitizeName(configName) + ".cfg"
	return g.writeWithBackup(filepath.Join(g.serviceDir, filename), buf.Bytes())
}

// WriteAllHosts writes a .cfg file for every active host.
func (g *Generator) WriteAllHosts() (int, error) {
	var hosts []models.Host
	if err := g.db.Where("active = '1'").Find(&hosts).Error; err != nil {
		return 0, err
	}
	var written int
	for _, h := range hosts {
		if err := g.WriteHost(h.ID); err != nil {
			return written, err
		}
		written++
	}
	return written, nil
}

// WriteAllServices writes one .cfg file per distinct config_name.
func (g *Generator) WriteAllServices() (int, error) {
	var names []string
	g.db.Model(&models.Service{}).Distinct("config_name").Where("active = '1'").Pluck("config_name", &names)
	var written int
	for _, name := range names {
		if err := g.WriteServiceGroup(name); err != nil {
			return written, err
		}
		written++
	}
	return written, nil
}

// WriteAll writes hosts + services and returns total files written.
func (g *Generator) WriteAll() (int, error) {
	h, err := g.WriteAllHosts()
	if err != nil {
		return h, err
	}
	s, err := g.WriteAllServices()
	return h + s, err
}

// --- Rendering ---

func (g *Generator) renderHost(h *models.Host) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("define host {\n")

	// use: resolved template name — must appear first per Nagios convention.
	if h.UseTemplate != 0 {
		if tpl := g.resolveHostTemplate(h.ID); tpl != "" {
			buf.WriteString(hostField("use", tpl))
		}
	}

	buf.WriteString(hostField("host_name", h.HostName))
	buf.WriteString(hostField("alias", h.Alias))

	if h.DisplayName != "" {
		buf.WriteString(hostField("display_name", h.DisplayName))
	}

	buf.WriteString(hostField("address", h.Address))

	// check_command: skip if '0' (inheriting from template via use) or empty.
	if cc := bangDecode(h.CheckCommand); cc != "" && cc != "0" {
		buf.WriteString(hostField("check_command", cc))
	}

	if h.MaxCheckAttempts != nil {
		buf.WriteString(hostField("max_check_attempts", fmt.Sprint(*h.MaxCheckAttempts)))
	}
	if h.CheckInterval != nil {
		buf.WriteString(hostField("check_interval", fmt.Sprint(*h.CheckInterval)))
	}
	if h.RetryInterval != nil {
		buf.WriteString(hostField("retry_interval", fmt.Sprint(*h.RetryInterval)))
	}

	// Value=2 means "inherit from template" — PHP skips these fields.
	if h.ActiveChecksEnabled != 2 {
		buf.WriteString(hostField("active_checks_enabled", fmt.Sprint(h.ActiveChecksEnabled)))
	}
	if h.PassiveChecksEnabled != 2 {
		buf.WriteString(hostField("passive_checks_enabled", fmt.Sprint(h.PassiveChecksEnabled)))
	}

	// contact_groups: resolved names from link table (FK is uint8 flag).
	if h.ContactGroups != 0 {
		if cg := g.resolveHostContactgroups(h.ID); cg != "" {
			buf.WriteString(hostField("contact_groups", cg))
		}
	}

	if h.NotificationOptions != "" {
		buf.WriteString(hostField("notification_options", h.NotificationOptions))
	}
	if h.Notes != "" {
		buf.WriteString(hostField("notes", h.Notes))
	}

	// register: '1' is the Nagios default so it's never written for real hosts.
	// Only templates have register=0.
	if h.Register == "0" {
		buf.WriteString(hostField("register", "0"))
	}

	buf.WriteString("}\n")
	return buf.Bytes(), nil
}

func (g *Generator) renderService(s *models.Service) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("define service {\n")

	// use: resolved template name.
	if s.UseTemplate != 0 {
		if tpl := g.resolveServiceTemplate(s.ID); tpl != "" {
			buf.WriteString(svcField("use", tpl))
		}
	}

	// host_name: resolved from tbl_lnkServiceToHost (comma-separated list).
	if hosts := g.resolveServiceHosts(s.ID); hosts != "" {
		buf.WriteString(svcField("host_name", hosts))
	}

	buf.WriteString(svcField("service_description", s.ServiceDescription))

	if s.DisplayName != "" {
		buf.WriteString(svcField("display_name", s.DisplayName))
	}

	if cc := bangDecode(s.CheckCommand); cc != "" && cc != "0" {
		buf.WriteString(svcField("check_command", cc))
	}

	if s.CheckInterval != nil {
		buf.WriteString(svcField("check_interval", fmt.Sprint(*s.CheckInterval)))
	}
	if s.RetryInterval != nil {
		buf.WriteString(svcField("retry_interval", fmt.Sprint(*s.RetryInterval)))
	}
	if s.MaxCheckAttempts != nil {
		buf.WriteString(svcField("max_check_attempts", fmt.Sprint(*s.MaxCheckAttempts)))
	}

	if s.ActiveChecksEnabled != 2 {
		buf.WriteString(svcField("active_checks_enabled", fmt.Sprint(s.ActiveChecksEnabled)))
	}
	if s.PassiveChecksEnabled != 2 {
		buf.WriteString(svcField("passive_checks_enabled", fmt.Sprint(s.PassiveChecksEnabled)))
	}

	if s.ContactGroups != 0 {
		if cg := g.resolveServiceContactgroups(s.ID); cg != "" {
			buf.WriteString(svcField("contact_groups", cg))
		}
	}

	if s.NotificationOptions != "" {
		buf.WriteString(svcField("notification_options", s.NotificationOptions))
	}
	if s.Notes != "" {
		buf.WriteString(svcField("notes", s.Notes))
	}

	if s.Register == "0" {
		buf.WriteString(svcField("register", "0"))
	}

	buf.WriteString("}\n")
	return buf.Bytes(), nil
}

// --- FK resolvers ---

func (g *Generator) resolveHostTemplate(hostID uint) string {
	var name string
	g.db.Raw(`SELECT t.template_name FROM tbl_hosttemplate t
		INNER JOIN tbl_lnkHostToHosttemplate l ON l.idSlave = t.id
		WHERE l.idMaster = ? ORDER BY l.idSort LIMIT 1`, hostID).Scan(&name)
	return name
}

func (g *Generator) resolveServiceTemplate(serviceID uint) string {
	var name string
	g.db.Raw(`SELECT t.template_name FROM tbl_servicetemplate t
		INNER JOIN tbl_lnkServiceToServicetemplate l ON l.idSlave = t.id
		WHERE l.idMaster = ? ORDER BY l.idSort LIMIT 1`, serviceID).Scan(&name)
	return name
}

func (g *Generator) resolveHostContactgroups(hostID uint) string {
	var names []string
	g.db.Raw(`SELECT c.contactgroup_name FROM tbl_contactgroup c
		INNER JOIN tbl_lnkHostToContactgroup l ON l.idSlave = c.id
		WHERE l.idMaster = ? ORDER BY l.idSort`, hostID).Scan(&names)
	return strings.Join(names, ",")
}

func (g *Generator) resolveServiceContactgroups(serviceID uint) string {
	var names []string
	g.db.Raw(`SELECT c.contactgroup_name FROM tbl_contactgroup c
		INNER JOIN tbl_lnkServiceToContactgroup l ON l.idSlave = c.id
		WHERE l.idMaster = ? ORDER BY l.idSort`, serviceID).Scan(&names)
	return strings.Join(names, ",")
}

func (g *Generator) resolveServiceHosts(serviceID uint) string {
	var names []string
	g.db.Raw(`SELECT h.host_name FROM tbl_host h
		INNER JOIN tbl_lnkServiceToHost l ON l.idSlave = h.id
		WHERE l.idMaster = ? ORDER BY l.idSort`, serviceID).Scan(&names)
	return strings.Join(names, ",")
}

// --- Formatting helpers ---

// hostField formats a key=value pair matching PHP NagiosQL host output:
// 4-space indent, value aligned at column 24 (key padded to 20 chars).
func hostField(key, value string) string {
	pad := 20 - len(key)
	if pad < 1 {
		pad = 1
	}
	return "    " + key + strings.Repeat(" ", pad) + value + "\n"
}

// svcField formats a key=value pair matching PHP NagiosQL service output:
// 4-space indent, value aligned at column 25 (key padded to 21 chars).
func svcField(key, value string) string {
	pad := 21 - len(key)
	if pad < 1 {
		pad = 1
	}
	return "    " + key + strings.Repeat(" ", pad) + value + "\n"
}

// bangDecode converts NagiosQL's internal bang-escape back to the '!'
// separator that Nagios uses for command arguments.
func bangDecode(s string) string {
	s = strings.ReplaceAll(s, `\::bang::`, "!")
	s = strings.ReplaceAll(s, "::bang::", "!")
	return s
}

// sanitizeName removes characters unsafe for filenames.
func sanitizeName(name string) string {
	r := strings.NewReplacer("/", "_", "\\", "_", ":", "_", " ", "_")
	return r.Replace(name)
}

// writeWithBackup copies the existing file to backup/ then overwrites it.
func (g *Generator) writeWithBackup(path string, content []byte) error {
	if _, err := os.Stat(path); err == nil {
		if err := os.MkdirAll(g.backupDir, 0755); err == nil {
			ts := time.Now().Format("20060102_150405")
			bak := filepath.Join(g.backupDir, filepath.Base(path)+"."+ts+".bak")
			data, _ := os.ReadFile(path)
			_ = os.WriteFile(bak, data, 0644)
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", filepath.Dir(path), err)
	}
	return os.WriteFile(path, content, 0644)
}
