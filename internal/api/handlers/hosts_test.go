package handlers_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	apimw "go-nagiosql/internal/api/middleware"
	"go-nagiosql/internal/api/handlers"
	"go-nagiosql/internal/models"
	"go-nagiosql/internal/services/auth"
	"github.com/labstack/echo/v5"
)

// injectClaims is middleware that injects fake JWT claims for tests.
func injectClaims(username string, admin bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("jwt_claims", &auth.Claims{Username: username, Admin: admin})
			return next(c)
		}
	}
}

func TestHostList(t *testing.T) {
	db := newTestDB(t)
	db.Create(&models.Host{HostName: "list-host", Alias: "List", Address: "1.1.1.1", Active: "1", Register: "1", LastModified: time.Now()})

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewHostHandler(db)
	e.GET("/hosts", h.List)

	rec := doRequest(t, e, http.MethodGet, "/hosts", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHostCreate(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewHostHandler(db)
	e.POST("/hosts", h.Create)

	rec := doRequest(t, e, http.MethodPost, "/hosts",
		`{"host_name":"new-host","alias":"New","address":"2.2.2.2","check_command":"check-host-alive"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHostGet_NotFound(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", false))
	h := handlers.NewHostHandler(db)
	e.GET("/hosts/:id", h.Get)

	rec := doRequest(t, e, http.MethodGet, "/hosts/9999", "")
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestHostDelete_WithLinkedService(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewHostHandler(db)
	e.DELETE("/hosts/:id", h.Delete)

	// Create host.
	host := models.Host{HostName: "dep-host", Alias: "Dep", Address: "3.3.3.3", Active: "1", Register: "1", LastModified: time.Now()}
	db.Create(&host)

	// Link a service to it.
	svc := models.Service{ServiceDescription: "PING", ConfigName: "dep-host", CheckCommand: "check_ping", Active: "1", Register: "1", LastModified: time.Now()}
	db.Create(&svc)
	db.Create(&models.LnkServiceToHost{ServiceID: svc.ID, HostID: host.ID})

	rec := doRequest(t, e, http.MethodDelete, fmt.Sprintf("/hosts/%d", host.ID), "")
	if rec.Code != http.StatusConflict {
		t.Errorf("expected 409 (host has linked service), got %d", rec.Code)
	}
}

func TestHostUnauthenticated(t *testing.T) {
	db := newTestDB(t)
	cfg := testCfg()
	svc := auth.New(db, cfg)
	e := echo.New()
	e.Use(apimw.JWTAuth(svc))
	h := handlers.NewHostHandler(db)
	e.GET("/hosts", h.List)

	// No Bearer token in request → middleware must return 401.
	rec := doRequest(t, e, http.MethodGet, "/hosts", "")
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for unauthenticated request, got %d", rec.Code)
	}
}
