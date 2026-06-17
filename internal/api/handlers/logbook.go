package handlers

import (
	"fmt"
	"net/http"

	"go-nagiosql/internal/models"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// LogbookHandler handles GET /api/v1/logbook (read-only).
type LogbookHandler struct{ db *gorm.DB }

// NewLogbookHandler creates a LogbookHandler.
func NewLogbookHandler(db *gorm.DB) *LogbookHandler { return &LogbookHandler{db: db} }

// ListLogbook godoc
// @Summary      List audit log entries
// @Tags         logbook
// @Produce      json
// @Param        from         query  string  false  "Start date (RFC3339)"
// @Param        to           query  string  false  "End date (RFC3339)"
// @Param        user         query  string  false  "Filter by username"
// @Param        object_type  query  string  false  "Filter by object type"
// @Param        page         query  int     false  "Page"
// @Param        limit        query  int     false  "Limit"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Router       /logbook [get]
func (h *LogbookHandler) List(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "created_at")

	q := h.db.Model(&models.Logbook{})
	if from := c.QueryParam("from"); from != "" {
		q = q.Where("created_at >= ?", from)
	}
	if to := c.QueryParam("to"); to != "" {
		q = q.Where("created_at <= ?", to)
	}
	if user := c.QueryParam("user"); user != "" {
		q = q.Where("username = ?", user)
	}
	if ot := c.QueryParam("object_type"); ot != "" {
		q = q.Where("object_type = ?", ot)
	}

	var total int64
	q.Count(&total)

	var entries []models.Logbook
	if err := q.Order(fmt.Sprintf("%s %s", sp.Field, sp.Dir)).
		Offset((pp.Page - 1) * pp.Limit).Limit(pp.Limit).Find(&entries).Error; err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, ListResponse{Data: entries, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}
