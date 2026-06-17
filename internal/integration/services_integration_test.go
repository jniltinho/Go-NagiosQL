//go:build integration

package integration

import (
	"fmt"
	"net/http"
	"testing"
)

func TestIntegration_ServiceList(t *testing.T) {
	srv, _, _ := newTestServer(t)
	token := login(t, srv, "admin", "admin123")

	code, body := authGet(t, srv, token, "/api/v1/services")
	if code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %v", code, body)
	}
	if body["total"] == nil {
		t.Error("expected total field in response")
	}
}

func TestIntegration_ServiceCRUD(t *testing.T) {
	srv, _, db := newTestServer(t)
	db.Exec("DELETE FROM tbl_service WHERE service_description = 'INTEG-PING' AND config_name = 'integ-host'")
	token := login(t, srv, "admin", "admin123")

	// Create.
	code, body := authPost(t, srv, token, "/api/v1/services",
		`{"service_description":"INTEG-PING","config_name":"integ-host","check_command":"check_ping","active":"1","register":"1"}`)
	if code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d: %v", code, body)
	}
	id := int(body["id"].(float64))

	// Get.
	code, body = authGet(t, srv, token, fmt.Sprintf("/api/v1/services/%d", id))
	if code != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", code)
	}
	if body["service_description"] != "INTEG-PING" {
		t.Errorf("expected service_description=INTEG-PING, got %v", body["service_description"])
	}

	// Filter by config_name.
	code, _ = authGet(t, srv, token, "/api/v1/services?config_name=integ-host")
	if code != http.StatusOK {
		t.Errorf("filter: expected 200, got %d", code)
	}

	// Delete.
	code = authDelete(t, srv, token, fmt.Sprintf("/api/v1/services/%d", id))
	if code != http.StatusNoContent {
		t.Errorf("delete: expected 204, got %d", code)
	}

	// Get deleted → 404.
	code, _ = authGet(t, srv, token, fmt.Sprintf("/api/v1/services/%d", id))
	if code != http.StatusNotFound {
		t.Errorf("get after delete: expected 404, got %d", code)
	}
}

func TestIntegration_Service_Pagination(t *testing.T) {
	srv, _, _ := newTestServer(t)
	token := login(t, srv, "admin", "admin123")

	// Create a few services.
	for i := 0; i < 3; i++ {
		authPost(t, srv, token, "/api/v1/services",
			fmt.Sprintf(`{"service_description":"INTEG-PAGE-%d","config_name":"page-host","check_command":"x","active":"1","register":"1"}`, i))
	}

	// First page, limit=2.
	code, body := authGet(t, srv, token, "/api/v1/services?limit=2&page=1&config_name=page-host")
	if code != http.StatusOK {
		t.Fatalf("expected 200, got %d", code)
	}
	data, ok := body["data"].([]any)
	if !ok {
		t.Fatalf("data should be array, got %T", body["data"])
	}
	if len(data) > 2 {
		t.Errorf("expected at most 2 items per page, got %d", len(data))
	}
}
