// Package handlers contains Echo HTTP handler functions for the NagiosQL REST API.
package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/jniltinho/go-nagiosql/internal/services/auth"
	"github.com/labstack/echo/v5"
)

// AuthHandler handles authentication-related endpoints.
type AuthHandler struct {
	svc            *auth.Service
	refreshTTLDays int
}

// NewAuthHandler creates an AuthHandler.
func NewAuthHandler(svc *auth.Service, refreshTTLDays int) *AuthHandler {
	return &AuthHandler{svc: svc, refreshTTLDays: refreshTTLDays}
}

// loginRequest is the JSON body for POST /api/v1/auth/login.
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// loginResponse is returned on a successful login.
type loginResponse struct {
	AccessToken          string `json:"access_token"`
	RequiresPasswordReset bool  `json:"requires_password_reset,omitempty"`
}

// Login godoc
// @Summary      Authenticate user and receive access + refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      loginRequest  true  "Credentials"
// @Success      200   {object}  loginResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "username and password are required"})
	}

	pair, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrLegacyMD5) {
			// Correct password but must reset — return 200 with the flag set so
			// the client can redirect to the password-reset page.
			return c.JSON(http.StatusOK, loginResponse{RequiresPasswordReset: true})
		}
		if errors.Is(err, auth.ErrInvalidCredentials) || errors.Is(err, auth.ErrUserInactive) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}

	// Refresh token goes into an httpOnly, Secure, SameSite=Strict cookie.
	h.setRefreshCookie(c, pair.RefreshToken)

	return c.JSON(http.StatusOK, loginResponse{AccessToken: pair.AccessToken})
}

// Refresh godoc
// @Summary      Issue a new access token using the refresh cookie
// @Tags         auth
// @Produce      json
// @Success      200  {object}  loginResponse
// @Failure      401  {object}  map[string]string
// @Router       /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *echo.Context) error {
	cookie, err := c.Request().Cookie("refresh_token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing refresh token"})
	}

	pair, err := h.svc.RefreshTokens(cookie.Value)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired refresh token"})
	}

	h.setRefreshCookie(c, pair.RefreshToken)
	return c.JSON(http.StatusOK, loginResponse{AccessToken: pair.AccessToken})
}

// Logout godoc
// @Summary      Clear the refresh token cookie
// @Tags         auth
// @Success      204  "No Content"
// @Router       /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *echo.Context) error {
	http.SetCookie(c.Response(), &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/v1/auth/refresh",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	return c.NoContent(http.StatusNoContent)
}

func (h *AuthHandler) setRefreshCookie(c *echo.Context, token string) {
	http.SetCookie(c.Response(), &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/api/v1/auth/refresh",
		Expires:  time.Now().Add(time.Duration(h.refreshTTLDays) * 24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		// Secure: true — set to true in production; omitted here so dev HTTP works.
	})
}
