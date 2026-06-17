//go:build integration

package integration

import (
	"fmt"
	"net/http"
	"testing"
)

func TestIntegration_HosttemplateCRUD(t *testing.T) {
	srv, _, db := newTestServer(t)
	db.Exec("DELETE FROM tbl_hosttemplate WHERE template_name = 'integ-generic-host'")
	tok := login(t, srv, "admin", "admin123")

	code, body := authPost(t, srv, tok, "/api/v1/hosttemplates",
		`{"template_name":"integ-generic-host","alias":"Generic Host Template","active":"1"}`)
	if code != http.StatusCreated {
		t.Fatalf("create hosttemplate: expected 201, got %d: %v", code, body)
	}
	id := int(body["id"].(float64))
	t.Cleanup(func() { db.Exec("DELETE FROM tbl_hosttemplate WHERE id = ?", id) })

	code, body = authGet(t, srv, tok, fmt.Sprintf("/api/v1/hosttemplates/%d", id))
	if code != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", code)
	}
	if body["template_name"] != "integ-generic-host" {
		t.Errorf("expected template_name=integ-generic-host, got %v", body["template_name"])
	}

	code, _ = authGet(t, srv, tok, "/api/v1/hosttemplates")
	if code != http.StatusOK {
		t.Errorf("list hosttemplates: expected 200, got %d", code)
	}

	code = authDelete(t, srv, tok, fmt.Sprintf("/api/v1/hosttemplates/%d", id))
	if code != http.StatusNoContent {
		t.Errorf("delete: expected 204, got %d", code)
	}
}

func TestIntegration_ServicetemplateCRUD(t *testing.T) {
	srv, _, db := newTestServer(t)
	db.Exec("DELETE FROM tbl_servicetemplate WHERE template_name = 'integ-generic-service'")
	tok := login(t, srv, "admin", "admin123")

	code, body := authPost(t, srv, tok, "/api/v1/servicetemplates",
		`{"template_name":"integ-generic-service","alias":"Generic Service Template","active":"1"}`)
	if code != http.StatusCreated {
		t.Fatalf("create servicetemplate: expected 201, got %d: %v", code, body)
	}
	id := int(body["id"].(float64))
	t.Cleanup(func() { db.Exec("DELETE FROM tbl_servicetemplate WHERE id = ?", id) })

	code, _ = authGet(t, srv, tok, "/api/v1/servicetemplates")
	if code != http.StatusOK {
		t.Errorf("list servicetemplates: expected 200, got %d", code)
	}

	code = authDelete(t, srv, tok, fmt.Sprintf("/api/v1/servicetemplates/%d", id))
	if code != http.StatusNoContent {
		t.Errorf("delete: expected 204, got %d", code)
	}
}

func TestIntegration_ContacttemplateCRUD(t *testing.T) {
	srv, _, db := newTestServer(t)
	db.Exec("DELETE FROM tbl_contacttemplate WHERE template_name = 'integ-generic-contact'")
	tok := login(t, srv, "admin", "admin123")

	code, body := authPost(t, srv, tok, "/api/v1/contacttemplates",
		`{"template_name":"integ-generic-contact","alias":"Generic Contact Template","active":"1"}`)
	if code != http.StatusCreated {
		t.Fatalf("create contacttemplate: expected 201, got %d: %v", code, body)
	}
	id := int(body["id"].(float64))
	t.Cleanup(func() { db.Exec("DELETE FROM tbl_contacttemplate WHERE id = ?", id) })

	code, _ = authGet(t, srv, tok, "/api/v1/contacttemplates")
	if code != http.StatusOK {
		t.Errorf("list contacttemplates: expected 200, got %d", code)
	}

	code = authDelete(t, srv, tok, fmt.Sprintf("/api/v1/contacttemplates/%d", id))
	if code != http.StatusNoContent {
		t.Errorf("delete: expected 204, got %d", code)
	}
}
