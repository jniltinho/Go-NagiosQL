//go:build integration

package integration

import (
	"fmt"
	"net/http"
	"testing"
)

func TestIntegration_VariableCRUD(t *testing.T) {
	srv, _, db := newTestServer(t)
	db.Exec("DELETE FROM tbl_variabledefinition WHERE name = 'INTEG_VAR'")
	tok := login(t, srv, "admin", "admin123")

	// Create.
	code, body := authPost(t, srv, tok, "/api/v1/variables",
		`{"name":"INTEG_VAR","value":"test-value","vartype":"string","active":"1"}`)
	if code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d: %v", code, body)
	}
	id := int(body["id"].(float64))
	t.Cleanup(func() { db.Exec("DELETE FROM tbl_variabledefinition WHERE id = ?", id) })

	// Get.
	code, body = authGet(t, srv, tok, fmt.Sprintf("/api/v1/variables/%d", id))
	if code != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", code)
	}
	if body["name"] != "INTEG_VAR" {
		t.Errorf("expected name=INTEG_VAR, got %v", body["name"])
	}

	// List.
	code, body = authGet(t, srv, tok, "/api/v1/variables")
	if code != http.StatusOK {
		t.Fatalf("list: expected 200, got %d", code)
	}
	if body["total"] == nil {
		t.Error("expected total field")
	}

	// Update.
	code, body = authPut(t, srv, tok, fmt.Sprintf("/api/v1/variables/%d", id),
		`{"name":"INTEG_VAR","value":"updated-value","vartype":"string","active":"1"}`)
	if code != http.StatusOK {
		t.Fatalf("update: expected 200, got %d: %v", code, body)
	}
	if body["value"] != "updated-value" {
		t.Errorf("expected value=updated-value, got %v", body["value"])
	}

	// Delete.
	code = authDelete(t, srv, tok, fmt.Sprintf("/api/v1/variables/%d", id))
	if code != http.StatusNoContent {
		t.Errorf("delete: expected 204, got %d", code)
	}

	// Get deleted → 404.
	code, _ = authGet(t, srv, tok, fmt.Sprintf("/api/v1/variables/%d", id))
	if code != http.StatusNotFound {
		t.Errorf("get after delete: expected 404, got %d", code)
	}
}

func TestIntegration_Variable_NameFilter(t *testing.T) {
	srv, _, db := newTestServer(t)
	db.Exec("DELETE FROM tbl_variabledefinition WHERE name LIKE 'INTEG_FILTER_%'")
	tok := login(t, srv, "admin", "admin123")

	for i := 0; i < 3; i++ {
		authPost(t, srv, tok, "/api/v1/variables",
			fmt.Sprintf(`{"name":"INTEG_FILTER_%d","value":"v%d","vartype":"string","active":"1"}`, i, i))
	}
	t.Cleanup(func() { db.Exec("DELETE FROM tbl_variabledefinition WHERE name LIKE 'INTEG_FILTER_%'") })

	code, body := authGet(t, srv, tok, "/api/v1/variables?name=INTEG_FILTER")
	if code != http.StatusOK {
		t.Fatalf("filter: expected 200, got %d", code)
	}
	total := int(body["total"].(float64))
	if total < 3 {
		t.Errorf("filter: expected at least 3, got %d", total)
	}
}
