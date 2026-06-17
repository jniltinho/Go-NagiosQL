package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jniltinho/go-nagiosql/internal/api/handlers"
	"github.com/jniltinho/go-nagiosql/internal/services/auth"
	"github.com/labstack/echo/v5"
)

func TestUserList_AdminOnly(t *testing.T) {
	db := newTestDB(t)
	seedUser(t, db, "admin1", "pass", true)

	e := echo.New()
	e.Use(injectClaims("admin1", true))
	cfg := testCfg()
	svc := auth.New(db, cfg)
	h := handlers.NewUserHandler(db, svc)
	e.GET("/users", h.List)

	rec := doRequest(t, e, http.MethodGet, "/users", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestUserCreate(t *testing.T) {
	db := newTestDB(t)
	seedUser(t, db, "admin2", "adminpass", true)

	e := echo.New()
	e.Use(injectClaims("admin2", true))
	cfg := testCfg()
	svc := auth.New(db, cfg)
	h := handlers.NewUserHandler(db, svc)
	e.POST("/users", h.Create)

	rec := doRequest(t, e, http.MethodPost, "/users",
		`{"username":"newuser","name":"New User","email":"new@example.com","password":"Password1!","admin":"0","active":"1"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	// Password must not leak in response.
	var resp map[string]any
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if _, hasPassword := resp["password"]; hasPassword {
		t.Error("password field leaked in response")
	}
}

func TestUserCreate_DuplicateUsername(t *testing.T) {
	db := newTestDB(t)
	seedUser(t, db, "dup-admin", "adminpass", true)
	seedUser(t, db, "existing", "pass", false)

	e := echo.New()
	e.Use(injectClaims("dup-admin", true))
	cfg := testCfg()
	svc := auth.New(db, cfg)
	h := handlers.NewUserHandler(db, svc)
	e.POST("/users", h.Create)

	rec := doRequest(t, e, http.MethodPost, "/users",
		`{"username":"existing","name":"Dup","email":"dup@example.com","password":"Pass1!","admin":"0","active":"1"}`)
	if rec.Code != http.StatusConflict {
		t.Errorf("expected 409 for duplicate username, got %d", rec.Code)
	}
}

func TestUserChangePassword(t *testing.T) {
	db := newTestDB(t)
	u := seedUser(t, db, "pwchange", "oldpass", false)

	e := echo.New()
	e.Use(injectClaims("pwchange", false))
	cfg := testCfg()
	svc := auth.New(db, cfg)
	h := handlers.NewUserHandler(db, svc)
	e.PUT("/users/:id/password", h.ChangePassword)

	rec := doRequest(t, e, http.MethodPut, fmt.Sprintf("/users/%d/password", u.ID),
		`{"new_password":"NewPass123!"}`)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestUserDelete_Self_Conflict(t *testing.T) {
	db := newTestDB(t)
	u := seedUser(t, db, "self-delete", "pass", true)

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// Inject claims with same username as the user we're trying to delete.
			c.Set("jwt_claims", &auth.Claims{Username: "self-delete", Admin: true})
			return next(c)
		}
	})
	cfg := testCfg()
	svc := auth.New(db, cfg)
	h := handlers.NewUserHandler(db, svc)
	e.DELETE("/users/:id", h.Delete)

	rec := doRequest(t, e, http.MethodDelete, fmt.Sprintf("/users/%d", u.ID), "")
	if rec.Code != http.StatusConflict {
		t.Errorf("expected 409 for self-delete, got %d", rec.Code)
	}
}

func TestUserGet_NotFound(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	e.Use(injectClaims("admin3", true))
	cfg := testCfg()
	svc := auth.New(db, cfg)
	h := handlers.NewUserHandler(db, svc)
	e.GET("/users/:id", h.Get)

	rec := doRequest(t, e, http.MethodGet, "/users/9999", "")
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}
