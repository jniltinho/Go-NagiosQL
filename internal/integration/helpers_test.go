//go:build integration

// Package integration contains end-to-end tests that require a live MariaDB instance.
// Run with:
//
//	make db-start
//	go test -tags integration ./internal/integration/... -v
//	make db-stop
package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jniltinho/go-nagiosql/internal/api"
	"github.com/jniltinho/go-nagiosql/internal/config"
	"github.com/jniltinho/go-nagiosql/internal/db"
	"github.com/jniltinho/go-nagiosql/internal/db/migrations"
	"github.com/jniltinho/go-nagiosql/internal/db/seeds"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// integrationCfg builds a config from environment variables with test defaults.
// Set NAGIOSQL_DATABASE_HOST, _PORT, _NAME, _USER, _PASSWORD to override.
func integrationCfg(t *testing.T) *config.Config {
	t.Helper()
	return &config.Config{
		Server: config.ServerConfig{Port: 8082, Dev: true},
		JWT: config.JWTConfig{
			Secret:         "integration-test-secret-key-minimum-32",
			AccessTTLMin:   15,
			RefreshTTLDays: 7,
		},
		Database: config.DatabaseConfig{
			Host:     envOr("NAGIOSQL_DATABASE_HOST", "127.0.0.1"),
			Port:     intEnvOr("NAGIOSQL_DATABASE_PORT", 3307),
			Name:     envOr("NAGIOSQL_DATABASE_NAME", "nagiosql_test"),
			User:     envOr("NAGIOSQL_DATABASE_USER", "nagiosql"),
			Password: envOr("NAGIOSQL_DATABASE_PASSWORD", "test"),
		},
		Nagios: config.NagiosConfig{
			ReloadTrigger: "/tmp/nagiosql-integ-reload.trigger",
			Binary:        "/bin/true",
		},
	}
}

// integrationDB opens and migrates a test MariaDB, seeding admin user.
func integrationDB(t *testing.T, cfg *config.Config) *gorm.DB {
	t.Helper()
	database, err := db.Open(cfg)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	if err := migrations.Migrate(database); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := seeds.SeedRequired(database, cfg, "admin", "admin123"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	return database
}

// newTestServer creates a full Echo server on an httptest.Server.
// Returns (server, config, db) so tests that need direct DB access can insert fixtures.
func newTestServer(t *testing.T) (*httptest.Server, *config.Config, *gorm.DB) {
	t.Helper()
	cfg := integrationCfg(t)
	database := integrationDB(t, cfg)

	e := echo.New()
	api.RegisterRoutes(e, database, cfg)

	srv := httptest.NewServer(e)
	t.Cleanup(srv.Close)
	return srv, cfg, database
}

// login POSTs to /api/v1/auth/login and returns the access token.
func login(t *testing.T, srv *httptest.Server, user, pass string) string {
	t.Helper()
	body := fmt.Sprintf(`{"username":%q,"password":%q}`, user, pass)
	resp, err := http.Post(srv.URL+"/api/v1/auth/login", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("login request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login failed: %d", resp.StatusCode)
	}
	var res map[string]any
	json.NewDecoder(resp.Body).Decode(&res)
	tok, _ := res["access_token"].(string)
	if tok == "" {
		t.Fatal("login: no access_token")
	}
	return tok
}

// authGet performs a GET with Bearer token and decodes JSON response.
func authGet(t *testing.T, srv *httptest.Server, token, path string) (int, map[string]any) {
	t.Helper()
	req, _ := http.NewRequest(http.MethodGet, srv.URL+path, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET %s: %v", path, err)
	}
	defer resp.Body.Close()
	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	return resp.StatusCode, result
}

// authPost performs a POST with Bearer token and decodes JSON response.
func authPost(t *testing.T, srv *httptest.Server, token, path, body string) (int, map[string]any) {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, srv.URL+path, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST %s: %v", path, err)
	}
	defer resp.Body.Close()
	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	return resp.StatusCode, result
}

// authPut performs a PUT with Bearer token and decodes JSON response.
func authPut(t *testing.T, srv *httptest.Server, token, path, body string) (int, map[string]any) {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPut, srv.URL+path, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT %s: %v", path, err)
	}
	defer resp.Body.Close()
	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	return resp.StatusCode, result
}

// authDelete performs a DELETE with Bearer token.
func authDelete(t *testing.T, srv *httptest.Server, token, path string) int {
	t.Helper()
	req, _ := http.NewRequest(http.MethodDelete, srv.URL+path, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE %s: %v", path, err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func intEnvOr(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		var n int
		fmt.Sscanf(v, "%d", &n)
		if n > 0 {
			return n
		}
	}
	return def
}
