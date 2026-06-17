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

// GroupHandler handles hostgroup, servicegroup, and contactgroup endpoints.
type GroupHandler struct{ db *gorm.DB }

// NewGroupHandler creates a GroupHandler.
func NewGroupHandler(db *gorm.DB) *GroupHandler { return &GroupHandler{db: db} }

// --- Hostgroups ---

// ListHostgroups godoc
// @Summary      List hostgroups
// @Tags         hostgroups
// @Produce      json
// @Param        page   query  int  false  "Page"
// @Param        limit  query  int  false  "Items per page"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Router       /hostgroups [get]
func (h *GroupHandler) ListHostgroups(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "hostgroup_name")
	var total int64
	h.db.Model(&models.Hostgroup{}).Count(&total)
	var groups []models.Hostgroup
	if err := h.db.Order(fmt.Sprintf("%s %s", sp.Field, sp.Dir)).
		Offset((pp.Page - 1) * pp.Limit).Limit(pp.Limit).Find(&groups).Error; err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, ListResponse{Data: groups, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}

// GetHostgroup godoc
// @Summary      Get a hostgroup by ID
// @Tags         hostgroups
// @Produce      json
// @Param        id  path  int  true  "Hostgroup ID"
// @Security     BearerAuth
// @Success      200  {object}  models.Hostgroup
// @Failure      404  {object}  map[string]string
// @Router       /hostgroups/{id} [get]
func (h *GroupHandler) GetHostgroup(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var g models.Hostgroup
	if err := h.db.First(&g, id).Error; err != nil {
		return NotFound(c)
	}
	return c.JSON(http.StatusOK, g)
}

// CreateHostgroup godoc
// @Summary      Create a hostgroup
// @Tags         hostgroups
// @Accept       json
// @Produce      json
// @Param        body  body  models.Hostgroup  true  "Hostgroup"
// @Security     BearerAuth
// @Success      201  {object}  models.Hostgroup
// @Router       /hostgroups [post]
func (h *GroupHandler) CreateHostgroup(c *echo.Context) error {
	var g models.Hostgroup
	if err := c.Bind(&g); err != nil {
		return BadRequest(c, "invalid body")
	}
	if g.HostgroupName == "" {
		return BadRequest(c, "hostgroup_name is required")
	}
	g.LastModified = time.Now()
	if g.Active == "" {
		g.Active = "1"
	}
	if g.Register == "" {
		g.Register = "1"
	}
	if err := h.db.Create(&g).Error; err != nil {
		return Conflict(c, "hostgroup_name already exists")
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "create", "hostgroup", g.HostgroupName, fmt.Sprintf("id=%d", g.ID))
	return c.JSON(http.StatusCreated, g)
}

// UpdateHostgroup godoc
// @Summary      Update a hostgroup
// @Tags         hostgroups
// @Accept       json
// @Produce      json
// @Param        id    path  int               true  "Hostgroup ID"
// @Param        body  body  models.Hostgroup  true  "Hostgroup"
// @Security     BearerAuth
// @Success      200  {object}  models.Hostgroup
// @Router       /hostgroups/{id} [put]
func (h *GroupHandler) UpdateHostgroup(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var existing models.Hostgroup
	if err := h.db.First(&existing, id).Error; err != nil {
		return NotFound(c)
	}
	var upd models.Hostgroup
	if err := c.Bind(&upd); err != nil {
		return BadRequest(c, "invalid body")
	}
	upd.ID = id
	upd.LastModified = time.Now()
	if err := h.db.Save(&upd).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "update", "hostgroup", upd.HostgroupName, fmt.Sprintf("id=%d", id))
	return c.JSON(http.StatusOK, upd)
}

// DeleteHostgroup godoc
// @Summary      Delete a hostgroup
// @Tags         hostgroups
// @Param        id  path  int  true  "Hostgroup ID"
// @Security     BearerAuth
// @Success      204
// @Failure      404  {object}  map[string]string
// @Router       /hostgroups/{id} [delete]
func (h *GroupHandler) DeleteHostgroup(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var g models.Hostgroup
	if err := h.db.First(&g, id).Error; err != nil {
		return NotFound(c)
	}
	if err := h.db.Delete(&g).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "delete", "hostgroup", g.HostgroupName, fmt.Sprintf("id=%d", id))
	return c.NoContent(http.StatusNoContent)
}

// AddHostgroupMember godoc
// @Summary      Add a host member to a hostgroup
// @Tags         hostgroups
// @Accept       json
// @Param        id    path  int  true  "Hostgroup ID"
// @Security     BearerAuth
// @Success      204
// @Router       /hostgroups/{id}/members [put]
// AddHostgroupMember adds a host to a hostgroup (PUT /api/v1/hostgroups/:id/members).
func (h *GroupHandler) AddHostgroupMember(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var body struct {
		HostID uint `json:"host_id"`
	}
	if err := c.Bind(&body); err != nil || body.HostID == 0 {
		return BadRequest(c, "host_id is required")
	}
	link := models.LnkHostgroupToHost{HostgroupID: id, HostID: body.HostID}
	if err := h.db.Create(&link).Error; err != nil {
		return Conflict(c, "member already exists")
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "added"})
}

// --- Servicegroups ---

// ListServicegroups godoc
// @Summary      List servicegroups
// @Tags         servicegroups
// @Produce      json
// @Param        page   query  int  false  "Page"
// @Param        limit  query  int  false  "Items per page"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Router       /servicegroups [get]
func (h *GroupHandler) ListServicegroups(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "servicegroup_name")
	var total int64
	h.db.Model(&models.Servicegroup{}).Count(&total)
	var groups []models.Servicegroup
	if err := h.db.Order(fmt.Sprintf("%s %s", sp.Field, sp.Dir)).
		Offset((pp.Page - 1) * pp.Limit).Limit(pp.Limit).Find(&groups).Error; err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, ListResponse{Data: groups, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}

// CreateServicegroup godoc
// @Summary      Create a servicegroup
// @Tags         servicegroups
// @Accept       json
// @Produce      json
// @Param        body  body  models.Servicegroup  true  "Servicegroup"
// @Security     BearerAuth
// @Success      201  {object}  models.Servicegroup
// @Router       /servicegroups [post]
func (h *GroupHandler) CreateServicegroup(c *echo.Context) error {
	var g models.Servicegroup
	if err := c.Bind(&g); err != nil {
		return BadRequest(c, "invalid body")
	}
	if g.ServicegroupName == "" {
		return BadRequest(c, "servicegroup_name is required")
	}
	g.LastModified = time.Now()
	if g.Active == "" {
		g.Active = "1"
	}
	if g.Register == "" {
		g.Register = "1"
	}
	if err := h.db.Create(&g).Error; err != nil {
		return Conflict(c, "servicegroup_name already exists")
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "create", "servicegroup", g.ServicegroupName, fmt.Sprintf("id=%d", g.ID))
	return c.JSON(http.StatusCreated, g)
}

// DeleteServicegroup godoc
// @Summary      Delete a servicegroup
// @Tags         servicegroups
// @Param        id  path  int  true  "Servicegroup ID"
// @Security     BearerAuth
// @Success      204
// @Router       /servicegroups/{id} [delete]
func (h *GroupHandler) DeleteServicegroup(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var g models.Servicegroup
	if err := h.db.First(&g, id).Error; err != nil {
		return NotFound(c)
	}
	if err := h.db.Delete(&g).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "delete", "servicegroup", g.ServicegroupName, fmt.Sprintf("id=%d", id))
	return c.NoContent(http.StatusNoContent)
}

// --- Contactgroups ---

// ListContactgroups godoc
// @Summary      List contactgroups
// @Tags         contactgroups
// @Produce      json
// @Param        page   query  int  false  "Page"
// @Param        limit  query  int  false  "Items per page"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Router       /contactgroups [get]
func (h *GroupHandler) ListContactgroups(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "contactgroup_name")
	var total int64
	h.db.Model(&models.Contactgroup{}).Count(&total)
	var groups []models.Contactgroup
	if err := h.db.Order(fmt.Sprintf("%s %s", sp.Field, sp.Dir)).
		Offset((pp.Page - 1) * pp.Limit).Limit(pp.Limit).Find(&groups).Error; err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, ListResponse{Data: groups, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}

// CreateContactgroup godoc
// @Summary      Create a contactgroup
// @Tags         contactgroups
// @Accept       json
// @Produce      json
// @Param        body  body  models.Contactgroup  true  "Contactgroup"
// @Security     BearerAuth
// @Success      201  {object}  models.Contactgroup
// @Router       /contactgroups [post]
func (h *GroupHandler) CreateContactgroup(c *echo.Context) error {
	var g models.Contactgroup
	if err := c.Bind(&g); err != nil {
		return BadRequest(c, "invalid body")
	}
	if g.ContactgroupName == "" {
		return BadRequest(c, "contactgroup_name is required")
	}
	g.LastModified = time.Now()
	if g.Active == "" {
		g.Active = "1"
	}
	if g.Register == "" {
		g.Register = "1"
	}
	if err := h.db.Create(&g).Error; err != nil {
		return Conflict(c, "contactgroup_name already exists")
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "create", "contactgroup", g.ContactgroupName, fmt.Sprintf("id=%d", g.ID))
	return c.JSON(http.StatusCreated, g)
}

// DeleteContactgroup godoc
// @Summary      Delete a contactgroup
// @Tags         contactgroups
// @Param        id  path  int  true  "Contactgroup ID"
// @Security     BearerAuth
// @Success      204
// @Router       /contactgroups/{id} [delete]
func (h *GroupHandler) DeleteContactgroup(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var g models.Contactgroup
	if err := h.db.First(&g, id).Error; err != nil {
		return NotFound(c)
	}
	if err := h.db.Delete(&g).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "delete", "contactgroup", g.ContactgroupName, fmt.Sprintf("id=%d", id))
	return c.NoContent(http.StatusNoContent)
}
