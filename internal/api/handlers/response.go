package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

// PageParams holds pagination query parameters.
type PageParams struct {
	Page  int
	Limit int
}

// SortParams holds sort query parameters.
type SortParams struct {
	Field string
	Dir   string
}

// ParsePage extracts ?page= and ?limit= from the request, with sane defaults.
func ParsePage(c *echo.Context) PageParams {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 500 {
		limit = 50
	}
	return PageParams{Page: page, Limit: limit}
}

// ParseSort extracts ?sort= and ?dir= with a default field.
func ParseSort(c *echo.Context, defaultField string) SortParams {
	field := c.QueryParam("sort")
	if field == "" {
		field = defaultField
	}
	dir := c.QueryParam("dir")
	if dir != "desc" {
		dir = "asc"
	}
	return SortParams{Field: field, Dir: dir}
}

// ListResponse wraps a paginated list with metadata.
type ListResponse struct {
	Data  any `json:"data"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// ErrResponse returns a JSON error response.
func ErrResponse(c *echo.Context, status int, msg string) error {
	return c.JSON(status, map[string]string{"error": msg})
}

// NotFound returns a 404 JSON error.
func NotFound(c *echo.Context) error {
	return ErrResponse(c, http.StatusNotFound, "not found")
}

// BadRequest returns a 400 JSON error.
func BadRequest(c *echo.Context, msg string) error {
	return ErrResponse(c, http.StatusBadRequest, msg)
}

// Conflict returns a 409 JSON error.
func Conflict(c *echo.Context, msg string) error {
	return ErrResponse(c, http.StatusConflict, msg)
}

// InternalError returns a 500 JSON error.
func InternalError(c *echo.Context, err error) error {
	return ErrResponse(c, http.StatusInternalServerError, err.Error())
}

// UintParam parses a named route parameter as uint. Returns 0 and false on failure.
func UintParam(c *echo.Context, name string) (uint, bool) {
	v, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil {
		return 0, false
	}
	return uint(v), true
}
