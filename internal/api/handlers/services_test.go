package handlers_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jniltinho/go-nagiosql/internal/api/handlers"
	"github.com/jniltinho/go-nagiosql/internal/models"
	"github.com/labstack/echo/v5"
)

func TestServiceList(t *testing.T) {
	db := newTestDB(t)
	db.Create(&models.Service{
		ServiceDescription: "PING", ConfigName: "host1",
		CheckCommand: "check_ping", Active: "1", Register: "1", LastModified: time.Now(),
	})

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewServiceHandler(db)
	e.GET("/services", h.List)

	rec := doRequest(t, e, http.MethodGet, "/services", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestServiceList_FilterByConfigName(t *testing.T) {
	db := newTestDB(t)
	db.Create(&models.Service{ServiceDescription: "PING", ConfigName: "web01", CheckCommand: "check_ping", Active: "1", Register: "1", LastModified: time.Now()})
	db.Create(&models.Service{ServiceDescription: "HTTP", ConfigName: "db01", CheckCommand: "check_http", Active: "1", Register: "1", LastModified: time.Now()})

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewServiceHandler(db)
	e.GET("/services", h.List)

	rec := doRequest(t, e, http.MethodGet, "/services?config_name=web01", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestServiceCreate(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewServiceHandler(db)
	e.POST("/services", h.Create)

	rec := doRequest(t, e, http.MethodPost, "/services",
		`{"service_description":"DISK","config_name":"testhost","check_command":"check_disk","active":"1","register":"1"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestServiceGet_NotFound(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", false))
	h := handlers.NewServiceHandler(db)
	e.GET("/services/:id", h.Get)

	rec := doRequest(t, e, http.MethodGet, "/services/9999", "")
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestServiceDelete(t *testing.T) {
	db := newTestDB(t)
	svc := models.Service{ServiceDescription: "CPU", ConfigName: "srv1", CheckCommand: "check_cpu", Active: "1", Register: "1", LastModified: time.Now()}
	db.Create(&svc)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewServiceHandler(db)
	e.DELETE("/services/:id", h.Delete)

	rec := doRequest(t, e, http.MethodDelete, fmt.Sprintf("/services/%d", svc.ID), "")
	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestServiceCreate_MissingDescription(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewServiceHandler(db)
	e.POST("/services", h.Create)

	rec := doRequest(t, e, http.MethodPost, "/services", `{"config_name":"host1","check_command":"x"}`)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing service_description, got %d", rec.Code)
	}
}
