//go:build integration

package integration

import (
	"fmt"
	"net/http"
	"testing"
)

func TestIntegration_CommandList(t *testing.T) {
	srv, _, _ := newTestServer(t)
	token := login(t, srv, "admin", "admin123")

	code, body := authGet(t, srv, token, "/api/v1/commands")
	if code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %v", code, body)
	}
	if body["total"] == nil {
		t.Error("expected total field in response")
	}
}

func TestIntegration_CommandCRUD(t *testing.T) {
	srv, _, db := newTestServer(t)
	db.Exec("DELETE FROM tbl_command WHERE command_name = 'integ-check-ping'")
	token := login(t, srv, "admin", "admin123")

	// Create.
	code, body := authPost(t, srv, token, "/api/v1/commands",
		`{"command_name":"integ-check-ping","command_line":"$USER1$/check_ping -H $HOSTADDRESS$ -w 100,20% -c 200,60%","command_type":0,"active":"1","register":"1"}`)
	if code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d: %v", code, body)
	}
	id := int(body["id"].(float64))

	// Get.
	code, body = authGet(t, srv, token, fmt.Sprintf("/api/v1/commands/%d", id))
	if code != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", code)
	}
	if body["command_name"] != "integ-check-ping" {
		t.Errorf("expected command_name=integ-check-ping, got %v", body["command_name"])
	}

	// List with type filter.
	code, _ = authGet(t, srv, token, "/api/v1/commands?type=check")
	if code != http.StatusOK {
		t.Errorf("list with type filter: expected 200, got %d", code)
	}

	// Delete.
	code = authDelete(t, srv, token, fmt.Sprintf("/api/v1/commands/%d", id))
	if code != http.StatusNoContent {
		t.Errorf("delete: expected 204, got %d", code)
	}

	// Get deleted → 404.
	code, _ = authGet(t, srv, token, fmt.Sprintf("/api/v1/commands/%d", id))
	if code != http.StatusNotFound {
		t.Errorf("get after delete: expected 404, got %d", code)
	}
}

func TestIntegration_Command_Duplicate(t *testing.T) {
	srv, _, db := newTestServer(t)
	// Clean up any leftover from a previous run.
	db.Exec("DELETE FROM tbl_command WHERE command_name = 'integ-dup-cmd'")

	token := login(t, srv, "admin", "admin123")

	body := `{"command_name":"integ-dup-cmd","command_line":"x","command_type":0}`
	code, _ := authPost(t, srv, token, "/api/v1/commands", body)
	if code != http.StatusCreated {
		t.Fatalf("first create: expected 201, got %d", code)
	}
	t.Cleanup(func() { db.Exec("DELETE FROM tbl_command WHERE command_name = 'integ-dup-cmd'") })

	code, _ = authPost(t, srv, token, "/api/v1/commands", body)
	if code != http.StatusConflict {
		t.Errorf("duplicate: expected 409, got %d", code)
	}
}
