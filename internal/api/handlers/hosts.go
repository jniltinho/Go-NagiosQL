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

// HostHandler handles /api/v1/hosts endpoints.
type HostHandler struct{ db *gorm.DB }

// NewHostHandler creates a HostHandler.
func NewHostHandler(db *gorm.DB) *HostHandler { return &HostHandler{db: db} }

// ListHosts godoc
// @Summary      List hosts
// @Tags         hosts
// @Produce      json
// @Param        page      query  int     false  "Page number (default 1)"
// @Param        limit     query  int     false  "Items per page (default 50)"
// @Param        sort      query  string  false  "Sort field (default host_name)"
// @Param        dir       query  string  false  "Sort direction: asc|desc"
// @Param        config_id query  int     false  "Filter by data domain"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Router       /hosts [get]
func (h *HostHandler) List(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "host_name")

	q := h.db.Model(&models.Host{})
	if cid := c.QueryParam("config_id"); cid != "" {
		q = q.Where("config_id = ?", cid)
	}
	if active := c.QueryParam("active"); active != "" {
		q = q.Where("active = ?", active)
	}

	var total int64
	q.Count(&total)

	var hosts []models.Host
	if err := q.Order(fmt.Sprintf("%s %s", sp.Field, sp.Dir)).
		Offset((pp.Page - 1) * pp.Limit).Limit(pp.Limit).Find(&hosts).Error; err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, ListResponse{Data: hosts, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}

// GetHost godoc
// @Summary      Get a host by ID
// @Tags         hosts
// @Produce      json
// @Param        id  path  int  true  "Host ID"
// @Security     BearerAuth
// @Success      200  {object}  models.Host
// @Failure      404  {object}  map[string]string
// @Router       /hosts/{id} [get]
func (h *HostHandler) Get(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var host models.Host
	if err := h.db.First(&host, id).Error; err != nil {
		return NotFound(c)
	}
	return c.JSON(http.StatusOK, host)
}

// CreateHost godoc
// @Summary      Create a host
// @Tags         hosts
// @Accept       json
// @Produce      json
// @Param        body  body  models.Host  true  "Host object"
// @Security     BearerAuth
// @Success      201  {object}  models.Host
// @Failure      400  {object}  map[string]string
// @Router       /hosts [post]
func (h *HostHandler) Create(c *echo.Context) error {
	var host models.Host
	if err := c.Bind(&host); err != nil {
		return BadRequest(c, "invalid request body")
	}
	if host.HostName == "" {
		return BadRequest(c, "host_name is required")
	}
	host.LastModified = time.Now()
	if host.Active == "" {
		host.Active = "1"
	}
	if host.Register == "" {
		host.Register = "1"
	}

	if err := h.db.Create(&host).Error; err != nil {
		return Conflict(c, "host_name already exists or constraint violated")
	}

	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "create", "host", host.HostName, fmt.Sprintf("id=%d", host.ID))
	return c.JSON(http.StatusCreated, host)
}

// UpdateHost godoc
// @Summary      Update a host
// @Tags         hosts
// @Accept       json
// @Produce      json
// @Param        id    path  int         true  "Host ID"
// @Param        body  body  models.Host true  "Host fields to update"
// @Security     BearerAuth
// @Success      200  {object}  models.Host
// @Failure      404  {object}  map[string]string
// @Router       /hosts/{id} [put]
func (h *HostHandler) Update(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var existing models.Host
	if err := h.db.First(&existing, id).Error; err != nil {
		return NotFound(c)
	}

	var updates models.Host
	if err := c.Bind(&updates); err != nil {
		return BadRequest(c, "invalid request body")
	}
	updates.ID = existing.ID
	updates.LastModified = time.Now()

	if err := h.db.Save(&updates).Error; err != nil {
		return InternalError(c, err)
	}

	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "update", "host", updates.HostName, fmt.Sprintf("id=%d", id))
	return c.JSON(http.StatusOK, updates)
}

// DeleteHost godoc
// @Summary      Delete a host
// @Tags         hosts
// @Param        id  path  int  true  "Host ID"
// @Security     BearerAuth
// @Success      204  "No Content"
// @Failure      404  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Router       /hosts/{id} [delete]
func (h *HostHandler) Delete(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var host models.Host
	if err := h.db.First(&host, id).Error; err != nil {
		return NotFound(c)
	}

	// Check if any service references this host.
	var svcCount int64
	h.db.Model(&models.LnkServiceToHost{}).Where("idSlave = ?", id).Count(&svcCount)
	if svcCount > 0 {
		return Conflict(c, fmt.Sprintf("host has %d linked services; delete them first", svcCount))
	}

	if err := h.db.Delete(&host).Error; err != nil {
		return InternalError(c, err)
	}

	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "delete", "host", host.HostName, fmt.Sprintf("id=%d", id))
	return c.NoContent(http.StatusNoContent)
}
