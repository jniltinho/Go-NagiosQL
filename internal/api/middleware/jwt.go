// Package middleware contains Echo middleware for NagiosQL.
package middleware

import (
	"net/http"
	"strings"

	"go-nagiosql/internal/services/auth"
	"github.com/labstack/echo/v5"
)

const claimsKey = "jwt_claims"

// JWTAuth returns an Echo middleware that validates Bearer tokens on every
// request. The parsed *auth.Claims are stored in the context under claimsKey
// so that handlers can retrieve them with ClaimsFromContext.
func JWTAuth(svc *auth.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing Authorization header"})
			}

			const prefix = "Bearer "
			if !strings.HasPrefix(header, prefix) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization header must use Bearer scheme"})
			}

			token := strings.TrimPrefix(header, prefix)
			claims, err := svc.ValidateAccessToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
			}

			c.Set(claimsKey, claims)
			return next(c)
		}
	}
}

// RequireAdmin returns a middleware that allows only admin users through.
// Must be chained after JWTAuth.
func RequireAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			claims := ClaimsFromContext(c)
			if claims == nil || !claims.Admin {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "admin privileges required"})
			}
			return next(c)
		}
	}
}

// ClaimsFromContext retrieves the parsed JWT claims stored by JWTAuth.
// Returns nil if no claims are present (unauthenticated route).
func ClaimsFromContext(c *echo.Context) *auth.Claims {
	v := c.Get(claimsKey)
	if v == nil {
		return nil
	}
	claims, _ := v.(*auth.Claims)
	return claims
}
