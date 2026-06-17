package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jniltinho/go-nagiosql/internal/api/handlers"
	"github.com/jniltinho/go-nagiosql/internal/config"
	"github.com/jniltinho/go-nagiosql/internal/models"
	"github.com/jniltinho/go-nagiosql/internal/services/auth"
	"github.com/jniltinho/go-nagiosql/internal/testhelpers"
	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// testCfg returns a minimal config for in-process testing.
func testCfg() *config.Config {
	return &config.Config{
		JWT: config.JWTConfig{
			Secret:         "test-secret-key-must-be-32chars!",
			AccessTTLMin:   15,
			RefreshTTLDays: 7,
		},
	}
}

// newTestDB opens an in-memory SQLite database with a minimal schema compatible
// with both SQLite and MariaDB (no MySQL-specific column types).
func newTestDB(t *testing.T) *gorm.DB {
	return testhelpers.NewDB(t)
}

// seedUser inserts a user with a bcrypt password into db.
func seedUser(t *testing.T, db *gorm.DB, username, password string, admin bool) models.User {
	t.Helper()
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 4) // cost=4 for test speed
	adminVal := "0"
	if admin {
		adminVal = "1"
	}
	u := models.User{
		Username:     username,
		Password:     string(hash),
		Name:         "Test User",
		Email:        username + "@example.com",
		Admin:        adminVal,
		Active:       "1",
		LastModified: time.Now(),
	}
	if err := db.Create(&u).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return u
}

// seedMD5User inserts a user with a raw MD5 hex password (simulating PHP legacy).
func seedMD5User(t *testing.T, db *gorm.DB, username, md5hash string) models.User {
	t.Helper()
	u := models.User{
		Username:     username,
		Password:     md5hash,
		Name:         "Legacy User",
		Email:        username + "@example.com",
		Admin:        "0",
		Active:       "1",
		LastModified: time.Now(),
	}
	if err := db.Create(&u).Error; err != nil {
		t.Fatalf("seed md5 user: %v", err)
	}
	return u
}

// doRequest sends a JSON request to the Echo handler and returns the recorder.
func doRequest(t *testing.T, e *echo.Echo, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != "" {
		buf.WriteString(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func TestLogin_Success(t *testing.T) {
	db := newTestDB(t)
	seedUser(t, db, "alice", "correct-horse", true)

	cfg := testCfg()
	svc := auth.New(db, cfg)
	h := handlers.NewAuthHandler(svc, cfg.JWT.RefreshTTLDays)

	e := echo.New()
	e.POST("/auth/login", h.Login)

	rec := doRequest(t, e, http.MethodPost, "/auth/login",
		`{"username":"alice","password":"correct-horse"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp map[string]any
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp["access_token"] == "" || resp["access_token"] == nil {
		t.Errorf("expected access_token in response: %v", resp)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	db := newTestDB(t)
	seedUser(t, db, "bob", "rightpass", false)

	cfg := testCfg()
	svc := auth.New(db, cfg)
	h := handlers.NewAuthHandler(svc, cfg.JWT.RefreshTTLDays)

	e := echo.New()
	e.POST("/auth/login", h.Login)

	rec := doRequest(t, e, http.MethodPost, "/auth/login",
		`{"username":"bob","password":"wrongpass"}`)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestLogin_LegacyMD5_CorrectPassword(t *testing.T) {
	db := newTestDB(t)
	// MD5("password") = 5f4dcc3b5aa765d61d8327deb882cf99
	seedMD5User(t, db, "legacy", "5f4dcc3b5aa765d61d8327deb882cf99")

	cfg := testCfg()
	svc := auth.New(db, cfg)
	h := handlers.NewAuthHandler(svc, cfg.JWT.RefreshTTLDays)

	e := echo.New()
	e.POST("/auth/login", h.Login)

	rec := doRequest(t, e, http.MethodPost, "/auth/login",
		`{"username":"legacy","password":"password"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 with reset flag, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp map[string]any
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp["requires_password_reset"] != true {
		t.Errorf("expected requires_password_reset=true: %v", resp)
	}
	if resp["access_token"] != nil && resp["access_token"] != "" {
		t.Errorf("expected no access_token for MD5 login: %v", resp)
	}
}

func TestLogin_LegacyMD5_WrongPassword(t *testing.T) {
	db := newTestDB(t)
	seedMD5User(t, db, "legacy2", "5f4dcc3b5aa765d61d8327deb882cf99")

	cfg := testCfg()
	svc := auth.New(db, cfg)
	h := handlers.NewAuthHandler(svc, cfg.JWT.RefreshTTLDays)

	e := echo.New()
	e.POST("/auth/login", h.Login)

	rec := doRequest(t, e, http.MethodPost, "/auth/login",
		`{"username":"legacy2","password":"wrongpass"}`)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for wrong MD5 password, got %d", rec.Code)
	}
}

func TestLogin_MissingFields(t *testing.T) {
	db := newTestDB(t)
	cfg := testCfg()
	svc := auth.New(db, cfg)
	h := handlers.NewAuthHandler(svc, cfg.JWT.RefreshTTLDays)

	e := echo.New()
	e.POST("/auth/login", h.Login)

	rec := doRequest(t, e, http.MethodPost, "/auth/login", `{"username":""}`)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestLogout(t *testing.T) {
	db := newTestDB(t)
	cfg := testCfg()
	svc := auth.New(db, cfg)
	h := handlers.NewAuthHandler(svc, cfg.JWT.RefreshTTLDays)

	e := echo.New()
	e.POST("/auth/logout", h.Logout)

	rec := doRequest(t, e, http.MethodPost, "/auth/logout", "")
	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}
}
