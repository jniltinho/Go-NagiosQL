//go:build integration

package integration

import (
	"fmt"
	"net/http"
	"testing"
)

func TestIntegration_ContactCRUD(t *testing.T) {
	srv, _, db := newTestServer(t)
	db.Exec("DELETE FROM tbl_contact WHERE contact_name = 'integ-ops'")
	tok := login(t, srv, "admin", "admin123")

	// Create.
	code, body := authPost(t, srv, tok, "/api/v1/contacts", `{
		"contact_name":"integ-ops",
		"alias":"Ops Team",
		"email":"ops@example.com",
		"host_notifications_enabled":1,
		"service_notifications_enabled":1,
		"host_notification_options":"d,u,r",
		"service_notification_options":"w,u,c,r",
		"active":"1","register":"1"
	}`)
	if code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d: %v", code, body)
	}
	id := int(body["id"].(float64))
	t.Cleanup(func() { db.Exec("DELETE FROM tbl_contact WHERE id = ?", id) })

	// Get.
	code, body = authGet(t, srv, tok, fmt.Sprintf("/api/v1/contacts/%d", id))
	if code != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", code)
	}
	if body["contact_name"] != "integ-ops" {
		t.Errorf("expected contact_name=integ-ops, got %v", body["contact_name"])
	}

	// List.
	code, _ = authGet(t, srv, tok, "/api/v1/contacts")
	if code != http.StatusOK {
		t.Errorf("list: expected 200, got %d", code)
	}

	// Update.
	code, body = authPut(t, srv, tok, fmt.Sprintf("/api/v1/contacts/%d", id), `{
		"contact_name":"integ-ops",
		"alias":"Operations",
		"email":"ops@example.com",
		"host_notifications_enabled":1,
		"service_notifications_enabled":1,
		"host_notification_options":"d,u,r",
		"service_notification_options":"w,u,c,r",
		"active":"1","register":"1"
	}`)
	if code != http.StatusOK {
		t.Fatalf("update: expected 200, got %d: %v", code, body)
	}
	if body["alias"] != "Operations" {
		t.Errorf("expected alias=Operations, got %v", body["alias"])
	}

	// Delete.
	code = authDelete(t, srv, tok, fmt.Sprintf("/api/v1/contacts/%d", id))
	if code != http.StatusNoContent {
		t.Errorf("delete: expected 204, got %d", code)
	}
}
