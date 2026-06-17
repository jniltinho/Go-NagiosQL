package handlers_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"go-nagiosql/internal/api/handlers"
	"go-nagiosql/internal/models"
	"github.com/labstack/echo/v5"
)

// ── Hostdependency ────────────────────────────────────────────────────────────

func TestHostdependencyList(t *testing.T) {
	db := newTestDB(t)
	db.Create(&models.Hostdependency{ConfigName: "dep-1", Register: "1", Active: "1", LastModified: time.Now()})

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.GET("/hostdependencies", h.ListHostdependencies)

	rec := doRequest(t, e, http.MethodGet, "/hostdependencies", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHostdependencyCreate(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.POST("/hostdependencies", h.CreateHostdependency)

	rec := doRequest(t, e, http.MethodPost, "/hostdependencies",
		`{"config_name":"dep-new","inherits_parent":1,"active":"1","register":"1"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHostdependencyCreate_MissingName(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.POST("/hostdependencies", h.CreateHostdependency)

	rec := doRequest(t, e, http.MethodPost, "/hostdependencies", `{"inherits_parent":1}`)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestHostdependencyGet(t *testing.T) {
	db := newTestDB(t)
	row := models.Hostdependency{ConfigName: "dep-get", Register: "1", Active: "1", LastModified: time.Now()}
	db.Create(&row)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.GET("/hostdependencies/:id", h.GetHostdependency)

	rec := doRequest(t, e, http.MethodGet, fmt.Sprintf("/hostdependencies/%d", row.ID), "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHostdependencyGet_NotFound(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.GET("/hostdependencies/:id", h.GetHostdependency)

	rec := doRequest(t, e, http.MethodGet, "/hostdependencies/9999", "")
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestHostdependencyUpdate(t *testing.T) {
	db := newTestDB(t)
	row := models.Hostdependency{ConfigName: "dep-upd", Register: "1", Active: "1", LastModified: time.Now()}
	db.Create(&row)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.PUT("/hostdependencies/:id", h.UpdateHostdependency)

	rec := doRequest(t, e, http.MethodPut, fmt.Sprintf("/hostdependencies/%d", row.ID),
		`{"config_name":"dep-upd-changed","inherits_parent":1,"active":"1","register":"1"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHostdependencyDelete(t *testing.T) {
	db := newTestDB(t)
	row := models.Hostdependency{ConfigName: "dep-del", Register: "1", Active: "1", LastModified: time.Now()}
	db.Create(&row)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.DELETE("/hostdependencies/:id", h.DeleteHostdependency)

	rec := doRequest(t, e, http.MethodDelete, fmt.Sprintf("/hostdependencies/%d", row.ID), "")
	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}
}

// ── Hostescalation ────────────────────────────────────────────────────────────

func TestHostescalationList(t *testing.T) {
	db := newTestDB(t)
	db.Create(&models.Hostescalation{ConfigName: "esc-1", Register: "1", Active: "1", LastModified: time.Now()})

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.GET("/hostescalations", h.ListHostescalations)

	rec := doRequest(t, e, http.MethodGet, "/hostescalations", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHostescalationCreate(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.POST("/hostescalations", h.CreateHostescalation)

	rec := doRequest(t, e, http.MethodPost, "/hostescalations",
		`{"config_name":"esc-new","escalation_options":"r,u","active":"1","register":"1"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHostescalationCreate_MissingName(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.POST("/hostescalations", h.CreateHostescalation)

	rec := doRequest(t, e, http.MethodPost, "/hostescalations", `{"escalation_options":"r"}`)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestHostescalationGet(t *testing.T) {
	db := newTestDB(t)
	row := models.Hostescalation{ConfigName: "esc-get", Register: "1", Active: "1", LastModified: time.Now()}
	db.Create(&row)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.GET("/hostescalations/:id", h.GetHostescalation)

	rec := doRequest(t, e, http.MethodGet, fmt.Sprintf("/hostescalations/%d", row.ID), "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHostescalationUpdate(t *testing.T) {
	db := newTestDB(t)
	row := models.Hostescalation{ConfigName: "esc-upd", Register: "1", Active: "1", LastModified: time.Now()}
	db.Create(&row)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.PUT("/hostescalations/:id", h.UpdateHostescalation)

	rec := doRequest(t, e, http.MethodPut, fmt.Sprintf("/hostescalations/%d", row.ID),
		`{"config_name":"esc-upd-changed","escalation_options":"w","active":"1","register":"1"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHostescalationDelete(t *testing.T) {
	db := newTestDB(t)
	row := models.Hostescalation{ConfigName: "esc-del", Register: "1", Active: "1", LastModified: time.Now()}
	db.Create(&row)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.DELETE("/hostescalations/:id", h.DeleteHostescalation)

	rec := doRequest(t, e, http.MethodDelete, fmt.Sprintf("/hostescalations/%d", row.ID), "")
	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}
}

// ── Hostextinfo ───────────────────────────────────────────────────────────────

func TestHostextinfoList(t *testing.T) {
	db := newTestDB(t)
	db.Create(&models.Hostextinfo{HostName: 1, Notes: "test notes", Register: "1", Active: "1", LastModified: time.Now()})

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.GET("/hostextinfo", h.ListHostextinfo)

	rec := doRequest(t, e, http.MethodGet, "/hostextinfo", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHostextinfoCreate(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.POST("/hostextinfo", h.CreateHostextinfo)

	rec := doRequest(t, e, http.MethodPost, "/hostextinfo",
		`{"host_name":1,"notes":"Some notes","active":"1","register":"1"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHostextinfoCreate_MissingHost(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.POST("/hostextinfo", h.CreateHostextinfo)

	rec := doRequest(t, e, http.MethodPost, "/hostextinfo", `{"notes":"orphan"}`)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestHostextinfoGet(t *testing.T) {
	db := newTestDB(t)
	row := models.Hostextinfo{HostName: 2, Register: "1", Active: "1", LastModified: time.Now()}
	db.Create(&row)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.GET("/hostextinfo/:id", h.GetHostextinfo)

	rec := doRequest(t, e, http.MethodGet, fmt.Sprintf("/hostextinfo/%d", row.ID), "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHostextinfoUpdate(t *testing.T) {
	db := newTestDB(t)
	row := models.Hostextinfo{HostName: 3, Register: "1", Active: "1", LastModified: time.Now()}
	db.Create(&row)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.PUT("/hostextinfo/:id", h.UpdateHostextinfo)

	rec := doRequest(t, e, http.MethodPut, fmt.Sprintf("/hostextinfo/%d", row.ID),
		`{"host_name":3,"notes":"updated notes","active":"1","register":"1"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHostextinfoDelete(t *testing.T) {
	db := newTestDB(t)
	row := models.Hostextinfo{HostName: 4, Register: "1", Active: "1", LastModified: time.Now()}
	db.Create(&row)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.DELETE("/hostextinfo/:id", h.DeleteHostextinfo)

	rec := doRequest(t, e, http.MethodDelete, fmt.Sprintf("/hostextinfo/%d", row.ID), "")
	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}
}

// ── Servicedependency ─────────────────────────────────────────────────────────

func TestServicedependencyList(t *testing.T) {
	db := newTestDB(t)
	db.Create(&models.Servicedependency{ConfigName: "svc-dep-1", Register: "1", Active: "1", LastModified: time.Now()})

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.GET("/servicedependencies", h.ListServicedependencies)

	rec := doRequest(t, e, http.MethodGet, "/servicedependencies", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestServicedependencyCreate(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.POST("/servicedependencies", h.CreateServicedependency)

	rec := doRequest(t, e, http.MethodPost, "/servicedependencies",
		`{"config_name":"svc-dep-new","inherits_parent":1,"active":"1","register":"1"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestServicedependencyCreate_MissingName(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.POST("/servicedependencies", h.CreateServicedependency)

	rec := doRequest(t, e, http.MethodPost, "/servicedependencies", `{"inherits_parent":1}`)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestServicedependencyGet(t *testing.T) {
	db := newTestDB(t)
	row := models.Servicedependency{ConfigName: "svc-dep-get", Register: "1", Active: "1", LastModified: time.Now()}
	db.Create(&row)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.GET("/servicedependencies/:id", h.GetServicedependency)

	rec := doRequest(t, e, http.MethodGet, fmt.Sprintf("/servicedependencies/%d", row.ID), "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestServicedependencyUpdate(t *testing.T) {
	db := newTestDB(t)
	row := models.Servicedependency{ConfigName: "svc-dep-upd", Register: "1", Active: "1", LastModified: time.Now()}
	db.Create(&row)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.PUT("/servicedependencies/:id", h.UpdateServicedependency)

	rec := doRequest(t, e, http.MethodPut, fmt.Sprintf("/servicedependencies/%d", row.ID),
		`{"config_name":"svc-dep-changed","inherits_parent":0,"active":"1","register":"1"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestServicedependencyDelete(t *testing.T) {
	db := newTestDB(t)
	row := models.Servicedependency{ConfigName: "svc-dep-del", Register: "1", Active: "1", LastModified: time.Now()}
	db.Create(&row)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.DELETE("/servicedependencies/:id", h.DeleteServicedependency)

	rec := doRequest(t, e, http.MethodDelete, fmt.Sprintf("/servicedependencies/%d", row.ID), "")
	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}
}

// ── Serviceescalation ─────────────────────────────────────────────────────────

func TestServiceescalationList(t *testing.T) {
	db := newTestDB(t)
	db.Create(&models.Serviceescalation{ConfigName: "svc-esc-1", Register: "1", Active: "1", LastModified: time.Now()})

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.GET("/serviceescalations", h.ListServiceescalations)

	rec := doRequest(t, e, http.MethodGet, "/serviceescalations", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestServiceescalationCreate(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.POST("/serviceescalations", h.CreateServiceescalation)

	rec := doRequest(t, e, http.MethodPost, "/serviceescalations",
		`{"config_name":"svc-esc-new","escalation_options":"w,c","active":"1","register":"1"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestServiceescalationCreate_MissingName(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.POST("/serviceescalations", h.CreateServiceescalation)

	rec := doRequest(t, e, http.MethodPost, "/serviceescalations", `{"escalation_options":"w"}`)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestServiceescalationGet(t *testing.T) {
	db := newTestDB(t)
	row := models.Serviceescalation{ConfigName: "svc-esc-get", Register: "1", Active: "1", LastModified: time.Now()}
	db.Create(&row)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.GET("/serviceescalations/:id", h.GetServiceescalation)

	rec := doRequest(t, e, http.MethodGet, fmt.Sprintf("/serviceescalations/%d", row.ID), "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestServiceescalationUpdate(t *testing.T) {
	db := newTestDB(t)
	row := models.Serviceescalation{ConfigName: "svc-esc-upd", Register: "1", Active: "1", LastModified: time.Now()}
	db.Create(&row)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.PUT("/serviceescalations/:id", h.UpdateServiceescalation)

	rec := doRequest(t, e, http.MethodPut, fmt.Sprintf("/serviceescalations/%d", row.ID),
		`{"config_name":"svc-esc-changed","escalation_options":"c","active":"1","register":"1"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestServiceescalationDelete(t *testing.T) {
	db := newTestDB(t)
	row := models.Serviceescalation{ConfigName: "svc-esc-del", Register: "1", Active: "1", LastModified: time.Now()}
	db.Create(&row)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.DELETE("/serviceescalations/:id", h.DeleteServiceescalation)

	rec := doRequest(t, e, http.MethodDelete, fmt.Sprintf("/serviceescalations/%d", row.ID), "")
	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}
}

// ── Serviceextinfo ────────────────────────────────────────────────────────────

func TestServiceextinfoList(t *testing.T) {
	db := newTestDB(t)
	db.Create(&models.Serviceextinfo{HostName: 1, ServiceDescription: 1, Notes: "svc notes", Register: "1", Active: "1", LastModified: time.Now()})

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.GET("/serviceextinfo", h.ListServiceextinfo)

	rec := doRequest(t, e, http.MethodGet, "/serviceextinfo", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestServiceextinfoCreate(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.POST("/serviceextinfo", h.CreateServiceextinfo)

	rec := doRequest(t, e, http.MethodPost, "/serviceextinfo",
		`{"host_name":1,"service_description":1,"notes":"cpu notes","active":"1","register":"1"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestServiceextinfoCreate_MissingHost(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.POST("/serviceextinfo", h.CreateServiceextinfo)

	rec := doRequest(t, e, http.MethodPost, "/serviceextinfo", `{"service_description":1,"notes":"x"}`)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestServiceextinfoCreate_MissingService(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.POST("/serviceextinfo", h.CreateServiceextinfo)

	rec := doRequest(t, e, http.MethodPost, "/serviceextinfo", `{"host_name":1,"notes":"x"}`)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestServiceextinfoGet(t *testing.T) {
	db := newTestDB(t)
	row := models.Serviceextinfo{HostName: 2, ServiceDescription: 2, Register: "1", Active: "1", LastModified: time.Now()}
	db.Create(&row)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.GET("/serviceextinfo/:id", h.GetServiceextinfo)

	rec := doRequest(t, e, http.MethodGet, fmt.Sprintf("/serviceextinfo/%d", row.ID), "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestServiceextinfoUpdate(t *testing.T) {
	db := newTestDB(t)
	row := models.Serviceextinfo{HostName: 3, ServiceDescription: 3, Register: "1", Active: "1", LastModified: time.Now()}
	db.Create(&row)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.PUT("/serviceextinfo/:id", h.UpdateServiceextinfo)

	rec := doRequest(t, e, http.MethodPut, fmt.Sprintf("/serviceextinfo/%d", row.ID),
		`{"host_name":3,"service_description":3,"notes":"updated","active":"1","register":"1"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestServiceextinfoDelete(t *testing.T) {
	db := newTestDB(t)
	row := models.Serviceextinfo{HostName: 4, ServiceDescription: 4, Register: "1", Active: "1", LastModified: time.Now()}
	db.Create(&row)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewExtendedHandler(db)
	e.DELETE("/serviceextinfo/:id", h.DeleteServiceextinfo)

	rec := doRequest(t, e, http.MethodDelete, fmt.Sprintf("/serviceextinfo/%d", row.ID), "")
	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}
}
