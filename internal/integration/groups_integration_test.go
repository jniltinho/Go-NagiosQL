//go:build integration

package integration

import (
	"fmt"
	"net/http"
	"testing"
)

func TestIntegration_HostgroupCRUD(t *testing.T) {
	srv, _, db := newTestServer(t)
	db.Exec("DELETE FROM tbl_hostgroup WHERE hostgroup_name = 'integ-linux-servers'")
	tok := login(t, srv, "admin", "admin123")

	code, body := authPost(t, srv, tok, "/api/v1/hostgroups",
		`{"hostgroup_name":"integ-linux-servers","alias":"Linux Servers","active":"1","register":"1"}`)
	if code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d: %v", code, body)
	}
	id := int(body["id"].(float64))
	t.Cleanup(func() { db.Exec("DELETE FROM tbl_hostgroup WHERE id = ?", id) })

	code, body = authGet(t, srv, tok, fmt.Sprintf("/api/v1/hostgroups/%d", id))
	if code != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", code)
	}
	if body["hostgroup_name"] != "integ-linux-servers" {
		t.Errorf("expected hostgroup_name=integ-linux-servers, got %v", body["hostgroup_name"])
	}

	code, _ = authGet(t, srv, tok, "/api/v1/hostgroups")
	if code != http.StatusOK {
		t.Errorf("list: expected 200, got %d", code)
	}

	code, body = authPut(t, srv, tok, fmt.Sprintf("/api/v1/hostgroups/%d", id),
		`{"hostgroup_name":"integ-linux-servers","alias":"All Linux Servers","active":"1","register":"1"}`)
	if code != http.StatusOK {
		t.Fatalf("update: expected 200, got %d: %v", code, body)
	}

	code = authDelete(t, srv, tok, fmt.Sprintf("/api/v1/hostgroups/%d", id))
	if code != http.StatusNoContent {
		t.Errorf("delete: expected 204, got %d", code)
	}
}

func TestIntegration_ServicegroupCRUD(t *testing.T) {
	srv, _, db := newTestServer(t)
	db.Exec("DELETE FROM tbl_servicegroup WHERE servicegroup_name = 'integ-http-checks'")
	tok := login(t, srv, "admin", "admin123")

	code, body := authPost(t, srv, tok, "/api/v1/servicegroups",
		`{"servicegroup_name":"integ-http-checks","alias":"HTTP Checks","active":"1","register":"1"}`)
	if code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d: %v", code, body)
	}
	id := int(body["id"].(float64))
	t.Cleanup(func() { db.Exec("DELETE FROM tbl_servicegroup WHERE id = ?", id) })

	code, _ = authGet(t, srv, tok, "/api/v1/servicegroups")
	if code != http.StatusOK {
		t.Errorf("list: expected 200, got %d", code)
	}

	code = authDelete(t, srv, tok, fmt.Sprintf("/api/v1/servicegroups/%d", id))
	if code != http.StatusNoContent {
		t.Errorf("delete: expected 204, got %d", code)
	}
}

func TestIntegration_ContactgroupCRUD(t *testing.T) {
	srv, _, db := newTestServer(t)
	db.Exec("DELETE FROM tbl_contactgroup WHERE contactgroup_name = 'integ-admins'")
	tok := login(t, srv, "admin", "admin123")

	code, body := authPost(t, srv, tok, "/api/v1/contactgroups",
		`{"contactgroup_name":"integ-admins","alias":"Admin Team","active":"1","register":"1"}`)
	if code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d: %v", code, body)
	}
	id := int(body["id"].(float64))
	t.Cleanup(func() { db.Exec("DELETE FROM tbl_contactgroup WHERE id = ?", id) })

	code, _ = authGet(t, srv, tok, "/api/v1/contactgroups")
	if code != http.StatusOK {
		t.Errorf("list: expected 200, got %d", code)
	}

	code = authDelete(t, srv, tok, fmt.Sprintf("/api/v1/contactgroups/%d", id))
	if code != http.StatusNoContent {
		t.Errorf("delete: expected 204, got %d", code)
	}
}
