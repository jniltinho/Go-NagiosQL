package handlers

import (
	"net/http"
	"time"

	apimw "go-nagiosql/internal/api/middleware"
	"go-nagiosql/internal/models"
	"go-nagiosql/internal/services/logbook"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// SettingsHandler handles /api/v1/settings.
type SettingsHandler struct{ db *gorm.DB }

// NewSettingsHandler creates a SettingsHandler.
func NewSettingsHandler(db *gorm.DB) *SettingsHandler { return &SettingsHandler{db: db} }

// settingsUpdateRequest holds the fields callers may change.
type settingsUpdateRequest struct {
	BackupAge      *uint   `json:"backup_age"`
	DateFormat     *string `json:"date_format"`
	LogbookContent *uint   `json:"logbook_content"`
}

// Get godoc
// @Summary      Get global settings
// @Tags         settings
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.Settings
// @Router       /settings [get]
func (h *SettingsHandler) Get(c *echo.Context) error {
	var s models.Settings
	if err := h.db.First(&s).Error; err != nil {
		return NotFound(c)
	}
	return c.JSON(http.StatusOK, s)
}

// Update godoc
// @Summary      Update global settings
// @Tags         settings
// @Accept       json
// @Produce      json
// @Param        body  body  settingsUpdateRequest  true  "Fields to update"
// @Security     BearerAuth
// @Success      200  {object}  models.Settings
// @Router       /settings [put]
func (h *SettingsHandler) Update(c *echo.Context) error {
	var s models.Settings
	if err := h.db.First(&s).Error; err != nil {
		return NotFound(c)
	}

	var req settingsUpdateRequest
	if err := c.Bind(&req); err != nil {
		return BadRequest(c, "invalid body")
	}

	if req.BackupAge != nil {
		s.BackupAge = *req.BackupAge
	}
	if req.DateFormat != nil {
		s.DateFormat = *req.DateFormat
	}
	if req.LogbookContent != nil {
		s.LogbookContent = *req.LogbookContent
	}
	s.LastModified = time.Now()

	if err := h.db.Save(&s).Error; err != nil {
		return InternalError(c, err)
	}

	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "update", "settings", "global", "")
	return c.JSON(http.StatusOK, s)
}
