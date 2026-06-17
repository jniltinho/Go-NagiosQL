//go:build integration

package integration

import (
	"fmt"
	"net/http"
	"testing"
)

func TestIntegration_UserCRUD(t *testing.T) {
	srv, _, db := newTestServer(t)
	db.Exec("DELETE FROM tbl_user WHERE username = 'integ-testuser'")
	tok := login(t, srv, "admin", "admin123")

	// Create.
	code, body := authPost(t, srv, tok, "/api/v1/users", `{
		"username":"integ-testuser",
		"name":"Integration Test User",
		"email":"integ@example.com",
		"password":"IntegPass1!",
		"admin":"0",
		"active":"1"
	}`)
	if code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d: %v", code, body)
	}
	id := int(body["id"].(float64))
	t.Cleanup(func() { db.Exec("DELETE FROM tbl_user WHERE id = ?", id) })

	// Password must not appear in response.
	if body["password"] != nil {
		t.Error("password leaked in create response")
	}

	// Get.
	code, body = authGet(t, srv, tok, fmt.Sprintf("/api/v1/users/%d", id))
	if code != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", code)
	}
	if body["username"] != "integ-testuser" {
		t.Errorf("expected username=integ-testuser, got %v", body["username"])
	}

	// List.
	code, body = authGet(t, srv, tok, "/api/v1/users")
	if code != http.StatusOK {
		t.Fatalf("list: expected 200, got %d", code)
	}
	if body["total"] == nil {
		t.Error("expected total field in list")
	}

	// Change password.
	code, _ = authPut(t, srv, tok, fmt.Sprintf("/api/v1/users/%d/password", id),
		`{"new_password":"NewIntegPass2!"}`)
	if code != http.StatusOK {
		t.Errorf("change password: expected 200, got %d", code)
	}

	// Login with new password.
	code, body = authPost(t, srv, "", "/api/v1/auth/login",
		`{"username":"integ-testuser","password":"NewIntegPass2!"}`)
	if code != http.StatusOK {
		t.Errorf("login with new password: expected 200, got %d: %v", code, body)
	}

	// Delete.
	code = authDelete(t, srv, tok, fmt.Sprintf("/api/v1/users/%d", id))
	if code != http.StatusNoContent {
		t.Errorf("delete: expected 204, got %d", code)
	}
}

func TestIntegration_User_SelfDeleteConflict(t *testing.T) {
	srv, _, db := newTestServer(t)
	_ = db
	tok := login(t, srv, "admin", "admin123")

	// Get admin's own ID.
	code, body := authGet(t, srv, tok, "/api/v1/users?limit=100")
	if code != http.StatusOK {
		t.Fatalf("list users: %d", code)
	}
	data, _ := body["data"].([]any)
	var adminID int
	for _, item := range data {
		m := item.(map[string]any)
		if m["username"] == "admin" {
			adminID = int(m["id"].(float64))
		}
	}
	if adminID == 0 {
		t.Skip("admin user not found in list — skipping self-delete test")
	}

	code = authDelete(t, srv, tok, fmt.Sprintf("/api/v1/users/%d", adminID))
	if code != http.StatusConflict {
		t.Errorf("self-delete: expected 409, got %d", code)
	}
}

func TestIntegration_User_DuplicateUsername(t *testing.T) {
	srv, _, db := newTestServer(t)
	db.Exec("DELETE FROM tbl_user WHERE username = 'integ-dup-user'")
	tok := login(t, srv, "admin", "admin123")

	body := `{"username":"integ-dup-user","name":"Dup","email":"dup@example.com","password":"Pass1!","admin":"0","active":"1"}`
	code, _ := authPost(t, srv, tok, "/api/v1/users", body)
	if code != http.StatusCreated {
		t.Fatalf("first create: expected 201, got %d", code)
	}
	t.Cleanup(func() { db.Exec("DELETE FROM tbl_user WHERE username = 'integ-dup-user'") })

	code, _ = authPost(t, srv, tok, "/api/v1/users", body)
	if code != http.StatusConflict {
		t.Errorf("duplicate: expected 409, got %d", code)
	}
}
