package nagconfig_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jniltinho/go-nagiosql/internal/models"
	"github.com/jniltinho/go-nagiosql/internal/services/nagconfig"
	"github.com/jniltinho/go-nagiosql/internal/testhelpers"
)

func newGenerator(t *testing.T) (*nagconfig.Generator, string) {
	t.Helper()
	db := testhelpers.NewDB(t)
	tmp := t.TempDir()
	g := nagconfig.New(db, filepath.Join(tmp, "hosts"), filepath.Join(tmp, "services"), filepath.Join(tmp, "backup"))
	return g, tmp
}

func TestWriteHost_BasicFields(t *testing.T) {
	db := testhelpers.NewDB(t)
	tmp := t.TempDir()
	g := nagconfig.New(db, filepath.Join(tmp, "hosts"), filepath.Join(tmp, "services"), filepath.Join(tmp, "backup"))

	host := models.Host{
		HostName:     "web01",
		Alias:        "Web Server 01",
		Address:      "10.0.0.1",
		CheckCommand: "check-host-alive",
		Active:       "1",
		Register:     "1",
		LastModified: time.Now(),
	}
	db.Create(&host)

	if err := g.WriteHost(host.ID); err != nil {
		t.Fatalf("WriteHost: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmp, "hosts", "web01.cfg"))
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	content := string(data)

	// PHP NagiosQL format: 4-space indent, value at column 24 (key padded to 20).
	// register=1 is the Nagios default and is never written.
	for _, want := range []string{
		"define host {",
		"    host_name           web01",
		"    alias               Web Server 01",
		"    address             10.0.0.1",
		"    check_command       check-host-alive",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("missing %q in output:\n%s", want, content)
		}
	}
}

func TestWriteHost_OptionalFieldsOmitted(t *testing.T) {
	db := testhelpers.NewDB(t)
	tmp := t.TempDir()
	g := nagconfig.New(db, filepath.Join(tmp, "hosts"), filepath.Join(tmp, "services"), filepath.Join(tmp, "backup"))

	host := models.Host{
		HostName:     "minimal-host",
		Alias:        "Minimal",
		Address:      "10.0.0.2",
		CheckCommand: "check-host-alive",
		Active:       "1",
		Register:     "1",
		LastModified: time.Now(),
		// Optional pointer fields are nil — should not appear in output.
	}
	db.Create(&host)

	if err := g.WriteHost(host.ID); err != nil {
		t.Fatalf("WriteHost: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(tmp, "hosts", "minimal-host.cfg"))
	content := string(data)

	for _, absent := range []string{"display_name", "max_check_attempts", "check_interval", "retry_interval"} {
		if strings.Contains(content, absent) {
			t.Errorf("optional field %q should not appear when nil:\n%s", absent, content)
		}
	}
}

func TestWriteHost_NotFound(t *testing.T) {
	g, _ := newGenerator(t)
	err := g.WriteHost(9999)
	if err == nil {
		t.Error("expected error for non-existent host")
	}
}

func TestWriteServiceGroup(t *testing.T) {
	db := testhelpers.NewDB(t)
	tmp := t.TempDir()
	g := nagconfig.New(db, filepath.Join(tmp, "hosts"), filepath.Join(tmp, "services"), filepath.Join(tmp, "backup"))

	for _, desc := range []string{"PING", "HTTP", "DISK"} {
		db.Create(&models.Service{
			ServiceDescription: desc,
			ConfigName:         "web01",
			CheckCommand:       "check_" + strings.ToLower(desc),
			Active:             "1",
			Register:           "1",
			LastModified:       time.Now(),
		})
	}

	if err := g.WriteServiceGroup("web01"); err != nil {
		t.Fatalf("WriteServiceGroup: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmp, "services", "web01.cfg"))
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	content := string(data)

	// Services: value at column 25 (key padded to 21 chars).
	// "service_description" = 19 chars → 2 spaces before value.
	for _, want := range []string{
		"define service {",
		"    service_description  PING",
		"    service_description  HTTP",
		"    service_description  DISK",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("missing %q in service output:\n%s", want, content)
		}
	}
}

func TestWriteServiceGroup_Empty(t *testing.T) {
	g, _ := newGenerator(t)
	// No services in DB — should return nil without creating a file.
	if err := g.WriteServiceGroup("nonexistent"); err != nil {
		t.Errorf("expected nil for empty group, got %v", err)
	}
}

func TestWriteAllHosts(t *testing.T) {
	db := testhelpers.NewDB(t)
	tmp := t.TempDir()
	g := nagconfig.New(db, filepath.Join(tmp, "hosts"), filepath.Join(tmp, "services"), filepath.Join(tmp, "backup"))

	for i, name := range []string{"host-a", "host-b", "host-c"} {
		_ = i
		db.Create(&models.Host{
			HostName:     name,
			Alias:        name,
			Address:      "10.0.1.1",
			CheckCommand: "check-host-alive",
			Active:       "1",
			Register:     "1",
			LastModified: time.Now(),
		})
	}
	// Inactive host — must not be written.
	db.Create(&models.Host{
		HostName:     "inactive-host",
		Alias:        "Inactive",
		Address:      "10.0.1.99",
		CheckCommand: "check-host-alive",
		Active:       "0",
		Register:     "1",
		LastModified: time.Now(),
	})

	written, err := g.WriteAllHosts()
	if err != nil {
		t.Fatalf("WriteAllHosts: %v", err)
	}
	if written != 3 {
		t.Errorf("expected 3 written, got %d", written)
	}
	// Inactive host file must not exist.
	if _, err := os.Stat(filepath.Join(tmp, "hosts", "inactive-host.cfg")); !os.IsNotExist(err) {
		t.Error("inactive host file should not exist")
	}
}

func TestWriteAllServices(t *testing.T) {
	db := testhelpers.NewDB(t)
	tmp := t.TempDir()
	g := nagconfig.New(db, filepath.Join(tmp, "hosts"), filepath.Join(tmp, "services"), filepath.Join(tmp, "backup"))

	for _, host := range []string{"srv1", "srv2"} {
		db.Create(&models.Service{
			ServiceDescription: "PING",
			ConfigName:         host,
			CheckCommand:       "check_ping",
			Active:             "1",
			Register:           "1",
			LastModified:       time.Now(),
		})
	}

	written, err := g.WriteAllServices()
	if err != nil {
		t.Fatalf("WriteAllServices: %v", err)
	}
	if written != 2 {
		t.Errorf("expected 2 files written, got %d", written)
	}
}

func TestWriteWithBackup(t *testing.T) {
	db := testhelpers.NewDB(t)
	tmp := t.TempDir()
	g := nagconfig.New(db, filepath.Join(tmp, "hosts"), filepath.Join(tmp, "services"), filepath.Join(tmp, "backup"))

	host := models.Host{
		HostName:     "bkp-host",
		Alias:        "Backup Host",
		Address:      "10.0.2.1",
		CheckCommand: "check-host-alive",
		Active:       "1",
		Register:     "1",
		LastModified: time.Now(),
	}
	db.Create(&host)

	// Write once to create the original.
	g.WriteHost(host.ID)

	// Write again — should create a backup.
	g.WriteHost(host.ID)

	entries, _ := os.ReadDir(filepath.Join(tmp, "backup"))
	if len(entries) == 0 {
		t.Error("expected at least one backup file")
	}
}

func TestSanitizeName_SpecialChars(t *testing.T) {
	db := testhelpers.NewDB(t)
	tmp := t.TempDir()
	g := nagconfig.New(db, filepath.Join(tmp, "hosts"), filepath.Join(tmp, "services"), filepath.Join(tmp, "backup"))

	host := models.Host{
		HostName:     "host/with:spaces and\\slashes",
		Alias:        "Special Host",
		Address:      "10.0.3.1",
		CheckCommand: "check-host-alive",
		Active:       "1",
		Register:     "1",
		LastModified: time.Now(),
	}
	db.Create(&host)
	if err := g.WriteHost(host.ID); err != nil {
		t.Fatalf("WriteHost with special chars: %v", err)
	}
	// Filename must not contain unsafe characters.
	entries, _ := os.ReadDir(filepath.Join(tmp, "hosts"))
	if len(entries) == 0 {
		t.Fatal("expected a file to be created")
	}
	name := entries[0].Name()
	for _, bad := range []string{"/", "\\", ":"} {
		if strings.Contains(name, bad) {
			t.Errorf("filename %q contains unsafe char %q", name, bad)
		}
	}
}
