//go:build integration

package integration

import (
	"fmt"
	"net/http"
	"testing"
)

func TestIntegration_TimeperiodCRUD(t *testing.T) {
	srv, _, db := newTestServer(t)
	db.Exec("DELETE FROM tbl_timeperiod WHERE timeperiod_name = 'integ-24x7'")
	tok := login(t, srv, "admin", "admin123")

	// Create with inline ranges.
	code, body := authPost(t, srv, tok, "/api/v1/timeperiods", `{
		"timeperiod_name":"integ-24x7",
		"alias":"24 Hours / 7 Days",
		"active":"1","register":"1",
		"ranges":[
			{"day":"monday","time_def":"00:00-24:00"},
			{"day":"sunday","time_def":"00:00-24:00"}
		]
	}`)
	if code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d: %v", code, body)
	}
	id := int(body["id"].(float64))
	t.Cleanup(func() {
		db.Exec("DELETE FROM tbl_timedefinition WHERE tipId = ?", id)
		db.Exec("DELETE FROM tbl_timeperiod WHERE id = ?", id)
	})

	// Get — should include definitions.
	code, body = authGet(t, srv, tok, fmt.Sprintf("/api/v1/timeperiods/%d", id))
	if code != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", code)
	}
	if body["timeperiod_name"] != "integ-24x7" {
		t.Errorf("expected timeperiod_name=integ-24x7, got %v", body["timeperiod_name"])
	}

	// List.
	code, _ = authGet(t, srv, tok, "/api/v1/timeperiods")
	if code != http.StatusOK {
		t.Errorf("list: expected 200, got %d", code)
	}

	// Update.
	code, body = authPut(t, srv, tok, fmt.Sprintf("/api/v1/timeperiods/%d", id), `{
		"timeperiod_name":"integ-24x7",
		"alias":"All Day",
		"active":"1","register":"1",
		"ranges":[{"day":"monday","time_def":"00:00-24:00"}]
	}`)
	if code != http.StatusOK {
		t.Fatalf("update: expected 200, got %d: %v", code, body)
	}

	// Delete.
	code = authDelete(t, srv, tok, fmt.Sprintf("/api/v1/timeperiods/%d", id))
	if code != http.StatusNoContent {
		t.Errorf("delete: expected 204, got %d", code)
	}
}
