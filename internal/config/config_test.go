package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	// Create a minimal config file with only the required jwt.secret.
	f, err := os.CreateTemp("", "nagiosql-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.WriteString(`[jwt]
secret = "this-is-a-very-long-secret-key-32ch"
`)
	f.Close()

	cfg, err := Load(f.Name())
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.Server.Port != 8081 {
		t.Errorf("expected default port 8081, got %d", cfg.Server.Port)
	}
	if cfg.JWT.AccessTTLMin != 15 {
		t.Errorf("expected default access_ttl_min 15, got %d", cfg.JWT.AccessTTLMin)
	}
	if cfg.Database.Host != "db" {
		t.Errorf("expected default db host 'db', got %s", cfg.Database.Host)
	}
}

func TestLoadEnvOverride(t *testing.T) {
	f, err := os.CreateTemp("", "nagiosql-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString(`[jwt]
secret = "this-is-a-very-long-secret-key-32ch"
`)
	f.Close()

	os.Setenv("NAGIOSQL_DATABASE_HOST", "mariadb-test")
	defer os.Unsetenv("NAGIOSQL_DATABASE_HOST")

	cfg, err := Load(f.Name())
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.Database.Host != "mariadb-test" {
		t.Errorf("env override failed: got %s", cfg.Database.Host)
	}
}

func TestLoadMissingSecret(t *testing.T) {
	f, err := os.CreateTemp("", "nagiosql-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString("[server]\nport = 8081\n")
	f.Close()

	_, err = Load(f.Name())
	if err == nil {
		t.Error("expected error for missing jwt.secret")
	}
}

func TestLoadShortSecret(t *testing.T) {
	f, err := os.CreateTemp("", "nagiosql-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString(`[jwt]
secret = "tooshort"
`)
	f.Close()

	_, err = Load(f.Name())
	if err == nil {
		t.Error("expected error for short jwt.secret")
	}
}
