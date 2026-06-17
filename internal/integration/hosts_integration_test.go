//go:build integration

package integration

import (
	"fmt"
	"net/http"
	"testing"
)

func TestIntegration_HostCRUD(t *testing.T) {
	srv, _, _ := newTestServer(t)
	tok := login(t, srv, "admin", "admin123")

	// Create.
	code, body := authPost(t, srv, tok, "/api/v1/hosts", `{
		"host_name":    "integ-host-01",
		"alias":        "Integration Host 01",
		"address":      "10.99.1.1",
		"check_command":"check-host-alive",
		"active":"1","register":"1"
	}`)
	if code != http.StatusCreated {
		t.Fatalf("create host: got %d body=%v", code, body)
	}
	id := int(body["id"].(float64))

	// Get.
	code, body = authGet(t, srv, tok, fmt.Sprintf("/api/v1/hosts/%d", id))
	if code != http.StatusOK {
		t.Fatalf("get host: got %d", code)
	}
	if body["host_name"] != "integ-host-01" {
		t.Errorf("host_name: got %v", body["host_name"])
	}

	// List — must contain the created host.
	code, body = authGet(t, srv, tok, "/api/v1/hosts")
	if code != http.StatusOK {
		t.Fatalf("list hosts: got %d", code)
	}
	data := body["data"].([]any)
	found := false
	for _, item := range data {
		m := item.(map[string]any)
		if m["host_name"] == "integ-host-01" {
			found = true
		}
	}
	if !found {
		t.Error("integ-host-01 not found in list")
	}

	// Delete.
	code = authDelete(t, srv, tok, fmt.Sprintf("/api/v1/hosts/%d", id))
	if code != http.StatusNoContent {
		t.Fatalf("delete host: got %d", code)
	}

	// Get after delete → 404.
	code, _ = authGet(t, srv, tok, fmt.Sprintf("/api/v1/hosts/%d", id))
	if code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", code)
	}
}

func TestIntegration_HostDelete_WithLinkedService_Conflicts(t *testing.T) {
	srv, _, db := newTestServer(t)
	tok := login(t, srv, "admin", "admin123")

	// Create host.
	code, body := authPost(t, srv, tok, "/api/v1/hosts", `{
		"host_name":"integ-conflict-host","alias":"Conflict","address":"10.99.2.1",
		"check_command":"check-host-alive","active":"1","register":"1"
	}`)
	if code != http.StatusCreated {
		t.Fatalf("create host: %d %v", code, body)
	}
	hostID := int(body["id"].(float64))

	// Create service.
	code, body = authPost(t, srv, tok, "/api/v1/services", `{
		"service_description":"PING","config_name":"integ-conflict-host",
		"check_command":"check_ping","active":"1","register":"1"
	}`)
	if code != http.StatusCreated {
		t.Fatalf("create service: %d %v", code, body)
	}
	svcID := int(body["id"].(float64))

	// Insert the host-service link directly (tbl_lnkServiceToHost: idMaster=svc, idSlave=host).
	if err := db.Exec("INSERT INTO tbl_lnkServiceToHost (idMaster, idSlave, idSort) VALUES (?, ?, 0)", svcID, hostID).Error; err != nil {
		t.Fatalf("insert link: %v", err)
	}

	// Attempt to delete host while service is linked → 409.
	code = authDelete(t, srv, tok, fmt.Sprintf("/api/v1/hosts/%d", hostID))
	if code != http.StatusConflict {
		t.Errorf("expected 409 conflict, got %d", code)
	}

	// Cleanup: remove link, delete service, then host.
	db.Exec("DELETE FROM tbl_lnkServiceToHost WHERE idMaster = ? AND idSlave = ?", svcID, hostID)
	authDelete(t, srv, tok, fmt.Sprintf("/api/v1/services/%d", svcID))
	authDelete(t, srv, tok, fmt.Sprintf("/api/v1/hosts/%d", hostID))
}

func TestIntegration_MonitoringSummary(t *testing.T) {
	srv, _, _ := newTestServer(t)
	tok := login(t, srv, "admin", "admin123")

	code, body := authGet(t, srv, tok, "/api/v1/monitoring/summary")
	if code != http.StatusOK {
		t.Fatalf("summary: got %d", code)
	}
	if _, ok := body["hosts"]; !ok {
		t.Errorf("missing 'hosts' in summary: %v", body)
	}
	if _, ok := body["services"]; !ok {
		t.Errorf("missing 'services' in summary: %v", body)
	}
}
