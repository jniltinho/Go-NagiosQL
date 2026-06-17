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

func TestCommandList(t *testing.T) {
	db := newTestDB(t)
	db.Create(&models.Command{
		CommandName: "check-host-alive", CommandLine: "$USER1$/check_ping -H $HOSTADDRESS$",
		CommandType: 0, Active: "1", Register: "1", LastModified: time.Now(),
	})

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewCommandHandler(db)
	e.GET("/commands", h.List)

	rec := doRequest(t, e, http.MethodGet, "/commands", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCommandList_TypeFilter(t *testing.T) {
	db := newTestDB(t)
	db.Create(&models.Command{CommandName: "check-ping", CommandLine: "...", CommandType: 0, Active: "1", Register: "1", LastModified: time.Now()})
	db.Create(&models.Command{CommandName: "notify-email", CommandLine: "...", CommandType: 1, Active: "1", Register: "1", LastModified: time.Now()})

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewCommandHandler(db)
	e.GET("/commands", h.List)

	// type=check → only check commands.
	rec := doRequest(t, e, http.MethodGet, "/commands?type=check", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	// type=notify → only notify commands.
	rec = doRequest(t, e, http.MethodGet, "/commands?type=notify", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestCommandCreate(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewCommandHandler(db)
	e.POST("/commands", h.Create)

	rec := doRequest(t, e, http.MethodPost, "/commands",
		`{"command_name":"new-check","command_line":"$USER1$/check_ping","command_type":0,"active":"1","register":"1"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCommandCreate_Duplicate(t *testing.T) {
	db := newTestDB(t)
	db.Create(&models.Command{CommandName: "dup-cmd", CommandLine: "x", CommandType: 0, Active: "1", Register: "1", LastModified: time.Now()})

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewCommandHandler(db)
	e.POST("/commands", h.Create)

	rec := doRequest(t, e, http.MethodPost, "/commands",
		`{"command_name":"dup-cmd","command_line":"x","command_type":0}`)
	if rec.Code != http.StatusConflict {
		t.Errorf("expected 409 for duplicate, got %d", rec.Code)
	}
}

func TestCommandDelete(t *testing.T) {
	db := newTestDB(t)
	cmd := models.Command{CommandName: "del-cmd", CommandLine: "x", CommandType: 0, Active: "1", Register: "1", LastModified: time.Now()}
	db.Create(&cmd)

	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewCommandHandler(db)
	e.DELETE("/commands/:id", h.Delete)

	rec := doRequest(t, e, http.MethodDelete, fmt.Sprintf("/commands/%d", cmd.ID), "")
	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}
}

func TestCommandGet_NotFound(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", false))
	h := handlers.NewCommandHandler(db)
	e.GET("/commands/:id", h.Get)

	rec := doRequest(t, e, http.MethodGet, "/commands/9999", "")
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestCommandCreate_MissingName(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin", true))
	h := handlers.NewCommandHandler(db)
	e.POST("/commands", h.Create)

	rec := doRequest(t, e, http.MethodPost, "/commands", `{"command_line":"x"}`)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing name, got %d", rec.Code)
	}
}

