//go:build integration

package integration

import (
	"net/http"
	"testing"
)

func TestIntegration_Settings_GetAndUpdate(t *testing.T) {
	srv, _, _ := newTestServer(t)
	tok := login(t, srv, "admin", "admin123")

	// GET /settings.
	code, body := authGet(t, srv, tok, "/api/v1/settings")
	if code != http.StatusOK {
		t.Fatalf("get settings: expected 200, got %d: %v", code, body)
	}
	// Settings record should exist (seeded).
	if body["id"] == nil {
		t.Error("expected id field in settings response")
	}

	// PUT /settings — update backup_age.
	code, body = authPut(t, srv, tok, "/api/v1/settings",
		`{"backup_age":14}`)
	if code != http.StatusOK {
		t.Fatalf("update settings: expected 200, got %d: %v", code, body)
	}
	if ba, ok := body["backup_age"].(float64); !ok || int(ba) != 14 {
		t.Errorf("expected backup_age=14, got %v", body["backup_age"])
	}

	// Restore original value.
	authPut(t, srv, tok, "/api/v1/settings", `{"backup_age":7}`)
}

func TestIntegration_Settings_NonAdminForbidden(t *testing.T) {
	srv, _, db := newTestServer(t)
	db.Exec("DELETE FROM tbl_user WHERE username = 'integ-nonadmin-settings'")

	// Create a non-admin user.
	adminTok := login(t, srv, "admin", "admin123")
	code, body := authPost(t, srv, adminTok, "/api/v1/users", `{
		"username":"integ-nonadmin-settings",
		"name":"Non Admin",
		"email":"na@example.com",
		"password":"NonAdminPass1!",
		"admin":"0","active":"1"
	}`)
	if code != http.StatusCreated {
		t.Fatalf("create non-admin user: %d %v", code, body)
	}
	id := int(body["id"].(float64))
	t.Cleanup(func() { db.Exec("DELETE FROM tbl_user WHERE id = ?", id) })

	// Login as non-admin.
	code, body = authPost(t, srv, "", "/api/v1/auth/login",
		`{"username":"integ-nonadmin-settings","password":"NonAdminPass1!"}`)
	if code != http.StatusOK {
		t.Fatalf("login non-admin: %d %v", code, body)
	}
	naTok := body["access_token"].(string)

	// GET /settings — allowed for any authenticated user.
	code, _ = authGet(t, srv, naTok, "/api/v1/settings")
	if code != http.StatusOK {
		t.Errorf("non-admin GET settings: expected 200, got %d", code)
	}

	// PUT /settings — admin-only → 403.
	code, _ = authPut(t, srv, naTok, "/api/v1/settings", `{"backup_age":99}`)
	if code != http.StatusForbidden {
		t.Errorf("non-admin PUT settings: expected 403, got %d", code)
	}
}
