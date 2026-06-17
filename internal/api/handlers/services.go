package handlers

import (
	"fmt"
	"net/http"
	"time"

	apimw "go-nagiosql/internal/api/middleware"
	"go-nagiosql/internal/models"
	"go-nagiosql/internal/services/logbook"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// ServiceHandler handles /api/v1/services endpoints.
type ServiceHandler struct{ db *gorm.DB }

// NewServiceHandler creates a ServiceHandler.
func NewServiceHandler(db *gorm.DB) *ServiceHandler { return &ServiceHandler{db: db} }

// ListServices godoc
// @Summary      List services
// @Tags         services
// @Produce      json
// @Param        host_id      query  int     false  "Filter by host ID"
// @Param        config_name  query  string  false  "Filter by config group name"
// @Param        page         query  int     false  "Page"
// @Param        limit        query  int     false  "Limit"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Router       /services [get]
func (h *ServiceHandler) List(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "service_description")

	q := h.db.Model(&models.Service{})
	if hostID := c.QueryParam("host_id"); hostID != "" {
		q = q.Joins("JOIN tbl_lnkServiceToHost lnk ON lnk.idMaster = tbl_service.id").
			Where("lnk.idSlave = ?", hostID)
	}
	if cn := c.QueryParam("config_name"); cn != "" {
		q = q.Where("config_name = ?", cn)
	}

	var total int64
	q.Count(&total)
	var svcs []models.Service
	if err := q.Order(fmt.Sprintf("%s %s", sp.Field, sp.Dir)).
		Offset((pp.Page - 1) * pp.Limit).Limit(pp.Limit).Find(&svcs).Error; err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, ListResponse{Data: svcs, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}

// GetService godoc
// @Summary      Get service by ID
// @Tags         services
// @Param        id  path  int  true  "Service ID"
// @Security     BearerAuth
// @Success      200  {object}  models.Service
// @Failure      404  {object}  map[string]string
// @Router       /services/{id} [get]
func (h *ServiceHandler) Get(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var svc models.Service
	if err := h.db.First(&svc, id).Error; err != nil {
		return NotFound(c)
	}
	return c.JSON(http.StatusOK, svc)
}

// CreateService godoc
// @Summary      Create a service
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        body  body  models.Service  true  "Service"
// @Security     BearerAuth
// @Success      201  {object}  models.Service
// @Router       /services [post]
func (h *ServiceHandler) Create(c *echo.Context) error {
	var svc models.Service
	if err := c.Bind(&svc); err != nil {
		return BadRequest(c, "invalid body")
	}
	if svc.ServiceDescription == "" || svc.ConfigName == "" {
		return BadRequest(c, "service_description and config_name are required")
	}
	svc.LastModified = time.Now()
	if svc.Active == "" {
		svc.Active = "1"
	}
	if svc.Register == "" {
		svc.Register = "1"
	}
	if err := h.db.Create(&svc).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "create", "service", svc.ServiceDescription, fmt.Sprintf("id=%d config=%s", svc.ID, svc.ConfigName))
	return c.JSON(http.StatusCreated, svc)
}

// UpdateService godoc
// @Summary      Update a service
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        id    path  int             true  "Service ID"
// @Param        body  body  models.Service  true  "Service fields"
// @Security     BearerAuth
// @Success      200  {object}  models.Service
// @Router       /services/{id} [put]
func (h *ServiceHandler) Update(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var existing models.Service
	if err := h.db.First(&existing, id).Error; err != nil {
		return NotFound(c)
	}
	var upd models.Service
	if err := c.Bind(&upd); err != nil {
		return BadRequest(c, "invalid body")
	}
	upd.ID = id
	upd.LastModified = time.Now()
	if err := h.db.Save(&upd).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "update", "service", upd.ServiceDescription, fmt.Sprintf("id=%d", id))
	return c.JSON(http.StatusOK, upd)
}

// DeleteService godoc
// @Summary      Delete a service
// @Tags         services
// @Param        id  path  int  true  "Service ID"
// @Security     BearerAuth
// @Success      204  "No Content"
// @Router       /services/{id} [delete]
func (h *ServiceHandler) Delete(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var svc models.Service
	if err := h.db.First(&svc, id).Error; err != nil {
		return NotFound(c)
	}
	// Remove link rows first.
	h.db.Where("idMaster = ?", id).Delete(&models.LnkServiceToHost{})
	if err := h.db.Delete(&svc).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "delete", "service", svc.ServiceDescription, fmt.Sprintf("id=%d", id))
	return c.NoContent(http.StatusNoContent)
}
