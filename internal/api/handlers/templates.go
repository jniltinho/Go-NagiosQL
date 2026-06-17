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

// TemplateHandler handles host, service, and contact template endpoints.
type TemplateHandler struct{ db *gorm.DB }

// NewTemplateHandler creates a TemplateHandler.
func NewTemplateHandler(db *gorm.DB) *TemplateHandler { return &TemplateHandler{db: db} }

// --- Host templates ---

// ListHosttemplates godoc
// @Summary      List host templates
// @Tags         templates
// @Produce      json
// @Param        page   query  int  false  "Page"
// @Param        limit  query  int  false  "Items per page"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Router       /hosttemplates [get]
func (h *TemplateHandler) ListHosttemplates(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "template_name")
	q := h.db.Model(&models.Hosttemplate{})
	if r := c.QueryParam("register"); r != "" {
		q = q.Where("active = ?", r)
	}
	var total int64
	q.Count(&total)
	var tpls []models.Hosttemplate
	if err := q.Order(fmt.Sprintf("%s %s", sp.Field, sp.Dir)).
		Offset((pp.Page - 1) * pp.Limit).Limit(pp.Limit).Find(&tpls).Error; err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, ListResponse{Data: tpls, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}

// GetHosttemplate godoc
// @Summary      Get a host template by ID
// @Tags         templates
// @Produce      json
// @Param        id  path  int  true  "Hosttemplate ID"
// @Security     BearerAuth
// @Success      200  {object}  models.Hosttemplate
// @Failure      404  {object}  map[string]string
// @Router       /hosttemplates/{id} [get]
func (h *TemplateHandler) GetHosttemplate(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var t models.Hosttemplate
	if err := h.db.First(&t, id).Error; err != nil {
		return NotFound(c)
	}
	return c.JSON(http.StatusOK, t)
}

// CreateHosttemplate godoc
// @Summary      Create a host template
// @Tags         templates
// @Accept       json
// @Produce      json
// @Param        body  body  models.Hosttemplate  true  "Hosttemplate"
// @Security     BearerAuth
// @Success      201  {object}  models.Hosttemplate
// @Router       /hosttemplates [post]
func (h *TemplateHandler) CreateHosttemplate(c *echo.Context) error {
	var t models.Hosttemplate
	if err := c.Bind(&t); err != nil {
		return BadRequest(c, "invalid body")
	}
	if t.TemplateName == "" {
		return BadRequest(c, "template_name is required")
	}
	t.LastModified = time.Now()
	if t.Active == "" {
		t.Active = "1"
	}
	if err := h.db.Create(&t).Error; err != nil {
		return Conflict(c, "template_name already exists")
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "create", "hosttemplate", t.TemplateName, fmt.Sprintf("id=%d", t.ID))
	return c.JSON(http.StatusCreated, t)
}

// DeleteHosttemplate godoc
// @Summary      Delete a host template
// @Tags         templates
// @Param        id  path  int  true  "Hosttemplate ID"
// @Security     BearerAuth
// @Success      204
// @Failure      404  {object}  map[string]string
// @Router       /hosttemplates/{id} [delete]
func (h *TemplateHandler) DeleteHosttemplate(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var t models.Hosttemplate
	if err := h.db.First(&t, id).Error; err != nil {
		return NotFound(c)
	}
	if err := h.db.Delete(&t).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "delete", "hosttemplate", t.TemplateName, fmt.Sprintf("id=%d", id))
	return c.NoContent(http.StatusNoContent)
}

// --- Service templates ---

// ListServicetemplates godoc
// @Summary      List service templates
// @Tags         templates
// @Produce      json
// @Param        page   query  int  false  "Page"
// @Param        limit  query  int  false  "Items per page"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Router       /servicetemplates [get]
func (h *TemplateHandler) ListServicetemplates(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "template_name")
	var total int64
	h.db.Model(&models.Servicetemplate{}).Count(&total)
	var tpls []models.Servicetemplate
	if err := h.db.Order(fmt.Sprintf("%s %s", sp.Field, sp.Dir)).
		Offset((pp.Page - 1) * pp.Limit).Limit(pp.Limit).Find(&tpls).Error; err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, ListResponse{Data: tpls, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}

// CreateServicetemplate godoc
// @Summary      Create a service template
// @Tags         templates
// @Accept       json
// @Produce      json
// @Param        body  body  models.Servicetemplate  true  "Servicetemplate"
// @Security     BearerAuth
// @Success      201  {object}  models.Servicetemplate
// @Router       /servicetemplates [post]
func (h *TemplateHandler) CreateServicetemplate(c *echo.Context) error {
	var t models.Servicetemplate
	if err := c.Bind(&t); err != nil {
		return BadRequest(c, "invalid body")
	}
	if t.TemplateName == "" {
		return BadRequest(c, "template_name is required")
	}
	t.LastModified = time.Now()
	if t.Active == "" {
		t.Active = "1"
	}
	if err := h.db.Create(&t).Error; err != nil {
		return Conflict(c, "template_name already exists")
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "create", "servicetemplate", t.TemplateName, fmt.Sprintf("id=%d", t.ID))
	return c.JSON(http.StatusCreated, t)
}

// DeleteServicetemplate godoc
// @Summary      Delete a service template
// @Tags         templates
// @Param        id  path  int  true  "Servicetemplate ID"
// @Security     BearerAuth
// @Success      204
// @Router       /servicetemplates/{id} [delete]
func (h *TemplateHandler) DeleteServicetemplate(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var t models.Servicetemplate
	if err := h.db.First(&t, id).Error; err != nil {
		return NotFound(c)
	}
	if err := h.db.Delete(&t).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "delete", "servicetemplate", t.TemplateName, fmt.Sprintf("id=%d", id))
	return c.NoContent(http.StatusNoContent)
}

// --- Contact templates ---

// ListContacttemplates godoc
// @Summary      List contact templates
// @Tags         templates
// @Produce      json
// @Param        page   query  int  false  "Page"
// @Param        limit  query  int  false  "Items per page"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Router       /contacttemplates [get]
func (h *TemplateHandler) ListContacttemplates(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "template_name")
	var total int64
	h.db.Model(&models.Contacttemplate{}).Count(&total)
	var tpls []models.Contacttemplate
	if err := h.db.Order(fmt.Sprintf("%s %s", sp.Field, sp.Dir)).
		Offset((pp.Page - 1) * pp.Limit).Limit(pp.Limit).Find(&tpls).Error; err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, ListResponse{Data: tpls, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}

// CreateContacttemplate godoc
// @Summary      Create a contact template
// @Tags         templates
// @Accept       json
// @Produce      json
// @Param        body  body  models.Contacttemplate  true  "Contacttemplate"
// @Security     BearerAuth
// @Success      201  {object}  models.Contacttemplate
// @Router       /contacttemplates [post]
func (h *TemplateHandler) CreateContacttemplate(c *echo.Context) error {
	var t models.Contacttemplate
	if err := c.Bind(&t); err != nil {
		return BadRequest(c, "invalid body")
	}
	if t.TemplateName == "" {
		return BadRequest(c, "template_name is required")
	}
	t.LastModified = time.Now()
	if t.Active == "" {
		t.Active = "1"
	}
	if err := h.db.Create(&t).Error; err != nil {
		return Conflict(c, "template_name already exists")
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "create", "contacttemplate", t.TemplateName, fmt.Sprintf("id=%d", t.ID))
	return c.JSON(http.StatusCreated, t)
}

// DeleteContacttemplate godoc
// @Summary      Delete a contact template
// @Tags         templates
// @Param        id  path  int  true  "Contacttemplate ID"
// @Security     BearerAuth
// @Success      204
// @Router       /contacttemplates/{id} [delete]
func (h *TemplateHandler) DeleteContacttemplate(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var t models.Contacttemplate
	if err := h.db.First(&t, id).Error; err != nil {
		return NotFound(c)
	}
	if err := h.db.Delete(&t).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "delete", "contacttemplate", t.TemplateName, fmt.Sprintf("id=%d", id))
	return c.NoContent(http.StatusNoContent)
}
