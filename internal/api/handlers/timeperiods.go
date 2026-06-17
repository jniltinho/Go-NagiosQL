package handlers

import (
	"fmt"
	"net/http"
	"time"

	apimw "github.com/jniltinho/go-nagiosql/internal/api/middleware"
	"github.com/jniltinho/go-nagiosql/internal/models"
	"github.com/jniltinho/go-nagiosql/internal/services/logbook"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// TimeperiodHandler handles /api/v1/timeperiods endpoints.
type TimeperiodHandler struct{ db *gorm.DB }

// NewTimeperiodHandler creates a TimeperiodHandler.
func NewTimeperiodHandler(db *gorm.DB) *TimeperiodHandler { return &TimeperiodHandler{db: db} }

// timeperiodRequest is used for create/update — includes inline time definitions.
type timeperiodRequest struct {
	models.Timeperiod
	Ranges []models.Timedefinition `json:"ranges"`
}

// List godoc
// @Summary      List time periods
// @Tags         timeperiods
// @Produce      json
// @Param        page   query  int  false  "Page"
// @Param        limit  query  int  false  "Items per page"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Router       /timeperiods [get]
func (h *TimeperiodHandler) List(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "timeperiod_name")
	var total int64
	h.db.Model(&models.Timeperiod{}).Count(&total)
	var tps []models.Timeperiod
	if err := h.db.Order(fmt.Sprintf("%s %s", sp.Field, sp.Dir)).
		Offset((pp.Page - 1) * pp.Limit).Limit(pp.Limit).Find(&tps).Error; err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, ListResponse{Data: tps, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}

// Get godoc
// @Summary      Get a time period by ID
// @Tags         timeperiods
// @Produce      json
// @Param        id  path  int  true  "Timeperiod ID"
// @Security     BearerAuth
// @Success      200  {object}  models.Timeperiod
// @Failure      404  {object}  map[string]string
// @Router       /timeperiods/{id} [get]
func (h *TimeperiodHandler) Get(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var tp models.Timeperiod
	if err := h.db.First(&tp, id).Error; err != nil {
		return NotFound(c)
	}
	var defs []models.Timedefinition
	h.db.Where("tipId = ?", tp.ID).Find(&defs)
	tp.Definitions = defs
	return c.JSON(http.StatusOK, tp)
}

// Create godoc
// @Summary      Create a time period
// @Tags         timeperiods
// @Accept       json
// @Produce      json
// @Param        body  body  timeperiodRequest  true  "Timeperiod"
// @Security     BearerAuth
// @Success      201  {object}  models.Timeperiod
// @Router       /timeperiods [post]
func (h *TimeperiodHandler) Create(c *echo.Context) error {
	var req timeperiodRequest
	if err := c.Bind(&req); err != nil {
		return BadRequest(c, "invalid body")
	}
	if req.TimeperiodName == "" {
		return BadRequest(c, "timeperiod_name is required")
	}
	req.Timeperiod.LastModified = time.Now()
	if req.Active == "" {
		req.Active = "1"
	}
	if req.Register == "" {
		req.Register = "1"
	}
	if err := h.db.Create(&req.Timeperiod).Error; err != nil {
		return Conflict(c, "timeperiod_name already exists")
	}
	for i := range req.Ranges {
		req.Ranges[i].TipID = req.Timeperiod.ID
		req.Ranges[i].LastModified = time.Now()
		h.db.Create(&req.Ranges[i]) //nolint:errcheck
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "create", "timeperiod", req.TimeperiodName, fmt.Sprintf("id=%d", req.Timeperiod.ID))
	req.Timeperiod.Definitions = req.Ranges
	return c.JSON(http.StatusCreated, req.Timeperiod)
}

// Update godoc
// @Summary      Update a time period
// @Tags         timeperiods
// @Accept       json
// @Produce      json
// @Param        id    path  int                true  "Timeperiod ID"
// @Param        body  body  timeperiodRequest  true  "Timeperiod"
// @Security     BearerAuth
// @Success      200  {object}  models.Timeperiod
// @Router       /timeperiods/{id} [put]
func (h *TimeperiodHandler) Update(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var existing models.Timeperiod
	if err := h.db.First(&existing, id).Error; err != nil {
		return NotFound(c)
	}
	var req timeperiodRequest
	if err := c.Bind(&req); err != nil {
		return BadRequest(c, "invalid body")
	}
	req.Timeperiod.ID = id
	req.Timeperiod.LastModified = time.Now()
	if err := h.db.Save(&req.Timeperiod).Error; err != nil {
		return InternalError(c, err)
	}
	// Replace definitions.
	h.db.Where("tipId = ?", id).Delete(&models.Timedefinition{})
	for i := range req.Ranges {
		req.Ranges[i].TipID = id
		req.Ranges[i].LastModified = time.Now()
		h.db.Create(&req.Ranges[i]) //nolint:errcheck
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "update", "timeperiod", req.TimeperiodName, fmt.Sprintf("id=%d", id))
	req.Timeperiod.Definitions = req.Ranges
	return c.JSON(http.StatusOK, req.Timeperiod)
}

// Delete godoc
// @Summary      Delete a time period
// @Tags         timeperiods
// @Param        id  path  int  true  "Timeperiod ID"
// @Security     BearerAuth
// @Success      204
// @Failure      404  {object}  map[string]string
// @Router       /timeperiods/{id} [delete]
func (h *TimeperiodHandler) Delete(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var tp models.Timeperiod
	if err := h.db.First(&tp, id).Error; err != nil {
		return NotFound(c)
	}
	h.db.Where("tipId = ?", id).Delete(&models.Timedefinition{})
	if err := h.db.Delete(&tp).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "delete", "timeperiod", tp.TimeperiodName, fmt.Sprintf("id=%d", id))
	return c.NoContent(http.StatusNoContent)
}
