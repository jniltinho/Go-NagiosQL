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

// ExtendedHandler handles hostdependency, hostescalation, hostextinfo,
// servicedependency, serviceescalation and serviceextinfo endpoints.
type ExtendedHandler struct{ db *gorm.DB }

// NewExtendedHandler returns a new ExtendedHandler backed by the given database connection.
func NewExtendedHandler(db *gorm.DB) *ExtendedHandler { return &ExtendedHandler{db: db} }

// ── Host dependencies ────────────────────────────────────────────────────────

// ListHostdependencies godoc
// @Summary      List host dependencies
// @Tags         hostdependencies
// @Produce      json
// @Param        page   query  int  false  "Page"
// @Param        limit  query  int  false  "Items per page"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Failure      500  {object}  map[string]string
// @Router       /hostdependencies [get]
func (h *ExtendedHandler) ListHostdependencies(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "config_name")
	var total int64
	h.db.Model(&models.Hostdependency{}).Count(&total)
	var rows []models.Hostdependency
	if err := h.db.Order(fmt.Sprintf("%s %s", sp.Field, sp.Dir)).
		Offset((pp.Page - 1) * pp.Limit).Limit(pp.Limit).Find(&rows).Error; err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, ListResponse{Data: rows, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}

// GetHostdependency godoc
// @Summary      Get a host dependency by ID
// @Tags         hostdependencies
// @Produce      json
// @Param        id  path  int  true  "Host dependency ID"
// @Security     BearerAuth
// @Success      200  {object}  models.Hostdependency
// @Failure      404  {object}  map[string]string
// @Router       /hostdependencies/{id} [get]
func (h *ExtendedHandler) GetHostdependency(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var row models.Hostdependency
	if err := h.db.First(&row, id).Error; err != nil {
		return NotFound(c)
	}
	return c.JSON(http.StatusOK, row)
}

// CreateHostdependency godoc
// @Summary      Create a host dependency
// @Tags         hostdependencies
// @Accept       json
// @Produce      json
// @Param        body  body  models.Hostdependency  true  "Hostdependency"
// @Security     BearerAuth
// @Success      201  {object}  models.Hostdependency
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /hostdependencies [post]
func (h *ExtendedHandler) CreateHostdependency(c *echo.Context) error {
	var row models.Hostdependency
	if err := c.Bind(&row); err != nil {
		return BadRequest(c, "invalid body")
	}
	if row.ConfigName == "" {
		return BadRequest(c, "config_name is required")
	}
	row.LastModified = time.Now()
	if row.Active == "" {
		row.Active = "1"
	}
	if row.Register == "" {
		row.Register = "1"
	}
	if err := h.db.Create(&row).Error; err != nil {
		return Conflict(c, "config_name already exists")
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "create", "hostdependency", row.ConfigName, fmt.Sprintf("id=%d", row.ID))
	return c.JSON(http.StatusCreated, row)
}

// UpdateHostdependency godoc
// @Summary      Update a host dependency
// @Tags         hostdependencies
// @Accept       json
// @Produce      json
// @Param        id    path  int                    true  "Host dependency ID"
// @Param        body  body  models.Hostdependency  true  "Hostdependency"
// @Security     BearerAuth
// @Success      200  {object}  models.Hostdependency
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /hostdependencies/{id} [put]
func (h *ExtendedHandler) UpdateHostdependency(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var existing models.Hostdependency
	if err := h.db.First(&existing, id).Error; err != nil {
		return NotFound(c)
	}
	var upd models.Hostdependency
	if err := c.Bind(&upd); err != nil {
		return BadRequest(c, "invalid body")
	}
	upd.ID = id
	upd.LastModified = time.Now()
	if err := h.db.Save(&upd).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "update", "hostdependency", upd.ConfigName, fmt.Sprintf("id=%d", id))
	return c.JSON(http.StatusOK, upd)
}

// DeleteHostdependency godoc
// @Summary      Delete a host dependency
// @Tags         hostdependencies
// @Param        id  path  int  true  "Host dependency ID"
// @Security     BearerAuth
// @Success      204
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /hostdependencies/{id} [delete]
func (h *ExtendedHandler) DeleteHostdependency(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var row models.Hostdependency
	if err := h.db.First(&row, id).Error; err != nil {
		return NotFound(c)
	}
	if err := h.db.Delete(&row).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "delete", "hostdependency", row.ConfigName, fmt.Sprintf("id=%d", id))
	return c.NoContent(http.StatusNoContent)
}

// ── Host escalations ─────────────────────────────────────────────────────────

// ListHostescalations godoc
// @Summary      List host escalations
// @Tags         hostescalations
// @Produce      json
// @Param        page   query  int  false  "Page"
// @Param        limit  query  int  false  "Items per page"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Failure      500  {object}  map[string]string
// @Router       /hostescalations [get]
func (h *ExtendedHandler) ListHostescalations(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "config_name")
	var total int64
	h.db.Model(&models.Hostescalation{}).Count(&total)
	var rows []models.Hostescalation
	if err := h.db.Order(fmt.Sprintf("%s %s", sp.Field, sp.Dir)).
		Offset((pp.Page - 1) * pp.Limit).Limit(pp.Limit).Find(&rows).Error; err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, ListResponse{Data: rows, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}

// GetHostescalation godoc
// @Summary      Get a host escalation by ID
// @Tags         hostescalations
// @Produce      json
// @Param        id  path  int  true  "Host escalation ID"
// @Security     BearerAuth
// @Success      200  {object}  models.Hostescalation
// @Failure      404  {object}  map[string]string
// @Router       /hostescalations/{id} [get]
func (h *ExtendedHandler) GetHostescalation(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var row models.Hostescalation
	if err := h.db.First(&row, id).Error; err != nil {
		return NotFound(c)
	}
	return c.JSON(http.StatusOK, row)
}

// CreateHostescalation godoc
// @Summary      Create a host escalation
// @Tags         hostescalations
// @Accept       json
// @Produce      json
// @Param        body  body  models.Hostescalation  true  "Hostescalation"
// @Security     BearerAuth
// @Success      201  {object}  models.Hostescalation
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /hostescalations [post]
func (h *ExtendedHandler) CreateHostescalation(c *echo.Context) error {
	var row models.Hostescalation
	if err := c.Bind(&row); err != nil {
		return BadRequest(c, "invalid body")
	}
	if row.ConfigName == "" {
		return BadRequest(c, "config_name is required")
	}
	row.LastModified = time.Now()
	if row.Active == "" {
		row.Active = "1"
	}
	if row.Register == "" {
		row.Register = "1"
	}
	if err := h.db.Create(&row).Error; err != nil {
		return Conflict(c, "config_name already exists")
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "create", "hostescalation", row.ConfigName, fmt.Sprintf("id=%d", row.ID))
	return c.JSON(http.StatusCreated, row)
}

// UpdateHostescalation godoc
// @Summary      Update a host escalation
// @Tags         hostescalations
// @Accept       json
// @Produce      json
// @Param        id    path  int                    true  "Host escalation ID"
// @Param        body  body  models.Hostescalation  true  "Hostescalation"
// @Security     BearerAuth
// @Success      200  {object}  models.Hostescalation
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /hostescalations/{id} [put]
func (h *ExtendedHandler) UpdateHostescalation(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var existing models.Hostescalation
	if err := h.db.First(&existing, id).Error; err != nil {
		return NotFound(c)
	}
	var upd models.Hostescalation
	if err := c.Bind(&upd); err != nil {
		return BadRequest(c, "invalid body")
	}
	upd.ID = id
	upd.LastModified = time.Now()
	if err := h.db.Save(&upd).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "update", "hostescalation", upd.ConfigName, fmt.Sprintf("id=%d", id))
	return c.JSON(http.StatusOK, upd)
}

// DeleteHostescalation godoc
// @Summary      Delete a host escalation
// @Tags         hostescalations
// @Param        id  path  int  true  "Host escalation ID"
// @Security     BearerAuth
// @Success      204
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /hostescalations/{id} [delete]
func (h *ExtendedHandler) DeleteHostescalation(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var row models.Hostescalation
	if err := h.db.First(&row, id).Error; err != nil {
		return NotFound(c)
	}
	if err := h.db.Delete(&row).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "delete", "hostescalation", row.ConfigName, fmt.Sprintf("id=%d", id))
	return c.NoContent(http.StatusNoContent)
}

// ── Host extended info ───────────────────────────────────────────────────────

// ListHostextinfo godoc
// @Summary      List host extended info entries
// @Tags         hostextinfo
// @Produce      json
// @Param        page   query  int  false  "Page"
// @Param        limit  query  int  false  "Items per page"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Failure      500  {object}  map[string]string
// @Router       /hostextinfo [get]
func (h *ExtendedHandler) ListHostextinfo(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "host_name")
	var total int64
	h.db.Model(&models.Hostextinfo{}).Count(&total)
	var rows []models.Hostextinfo
	if err := h.db.Order(fmt.Sprintf("%s %s", sp.Field, sp.Dir)).
		Offset((pp.Page - 1) * pp.Limit).Limit(pp.Limit).Find(&rows).Error; err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, ListResponse{Data: rows, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}

// GetHostextinfo godoc
// @Summary      Get a host extended info entry by ID
// @Tags         hostextinfo
// @Produce      json
// @Param        id  path  int  true  "Hostextinfo ID"
// @Security     BearerAuth
// @Success      200  {object}  models.Hostextinfo
// @Failure      404  {object}  map[string]string
// @Router       /hostextinfo/{id} [get]
func (h *ExtendedHandler) GetHostextinfo(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var row models.Hostextinfo
	if err := h.db.First(&row, id).Error; err != nil {
		return NotFound(c)
	}
	return c.JSON(http.StatusOK, row)
}

// CreateHostextinfo godoc
// @Summary      Create a host extended info entry
// @Tags         hostextinfo
// @Accept       json
// @Produce      json
// @Param        body  body  models.Hostextinfo  true  "Hostextinfo"
// @Security     BearerAuth
// @Success      201  {object}  models.Hostextinfo
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /hostextinfo [post]
func (h *ExtendedHandler) CreateHostextinfo(c *echo.Context) error {
	var row models.Hostextinfo
	if err := c.Bind(&row); err != nil {
		return BadRequest(c, "invalid body")
	}
	if row.HostName == 0 {
		return BadRequest(c, "host_name (host id) is required")
	}
	row.LastModified = time.Now()
	if row.Active == "" {
		row.Active = "1"
	}
	if row.Register == "" {
		row.Register = "1"
	}
	if err := h.db.Create(&row).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "create", "hostextinfo", fmt.Sprintf("host_id=%d", row.HostName), fmt.Sprintf("id=%d", row.ID))
	return c.JSON(http.StatusCreated, row)
}

// UpdateHostextinfo godoc
// @Summary      Update a host extended info entry
// @Tags         hostextinfo
// @Accept       json
// @Produce      json
// @Param        id    path  int                 true  "Hostextinfo ID"
// @Param        body  body  models.Hostextinfo  true  "Hostextinfo"
// @Security     BearerAuth
// @Success      200  {object}  models.Hostextinfo
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /hostextinfo/{id} [put]
func (h *ExtendedHandler) UpdateHostextinfo(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var existing models.Hostextinfo
	if err := h.db.First(&existing, id).Error; err != nil {
		return NotFound(c)
	}
	var upd models.Hostextinfo
	if err := c.Bind(&upd); err != nil {
		return BadRequest(c, "invalid body")
	}
	upd.ID = id
	upd.LastModified = time.Now()
	if err := h.db.Save(&upd).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "update", "hostextinfo", fmt.Sprintf("host_id=%d", upd.HostName), fmt.Sprintf("id=%d", id))
	return c.JSON(http.StatusOK, upd)
}

// DeleteHostextinfo godoc
// @Summary      Delete a host extended info entry
// @Tags         hostextinfo
// @Param        id  path  int  true  "Hostextinfo ID"
// @Security     BearerAuth
// @Success      204
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /hostextinfo/{id} [delete]
func (h *ExtendedHandler) DeleteHostextinfo(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var row models.Hostextinfo
	if err := h.db.First(&row, id).Error; err != nil {
		return NotFound(c)
	}
	if err := h.db.Delete(&row).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "delete", "hostextinfo", fmt.Sprintf("host_id=%d", row.HostName), fmt.Sprintf("id=%d", id))
	return c.NoContent(http.StatusNoContent)
}

// ── Service dependencies ─────────────────────────────────────────────────────

// ListServicedependencies godoc
// @Summary      List service dependencies
// @Tags         servicedependencies
// @Produce      json
// @Param        page   query  int  false  "Page"
// @Param        limit  query  int  false  "Items per page"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Failure      500  {object}  map[string]string
// @Router       /servicedependencies [get]
func (h *ExtendedHandler) ListServicedependencies(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "config_name")
	var total int64
	h.db.Model(&models.Servicedependency{}).Count(&total)
	var rows []models.Servicedependency
	if err := h.db.Order(fmt.Sprintf("%s %s", sp.Field, sp.Dir)).
		Offset((pp.Page - 1) * pp.Limit).Limit(pp.Limit).Find(&rows).Error; err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, ListResponse{Data: rows, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}

// GetServicedependency godoc
// @Summary      Get a service dependency by ID
// @Tags         servicedependencies
// @Produce      json
// @Param        id  path  int  true  "Service dependency ID"
// @Security     BearerAuth
// @Success      200  {object}  models.Servicedependency
// @Failure      404  {object}  map[string]string
// @Router       /servicedependencies/{id} [get]
func (h *ExtendedHandler) GetServicedependency(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var row models.Servicedependency
	if err := h.db.First(&row, id).Error; err != nil {
		return NotFound(c)
	}
	return c.JSON(http.StatusOK, row)
}

// CreateServicedependency godoc
// @Summary      Create a service dependency
// @Tags         servicedependencies
// @Accept       json
// @Produce      json
// @Param        body  body  models.Servicedependency  true  "Servicedependency"
// @Security     BearerAuth
// @Success      201  {object}  models.Servicedependency
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /servicedependencies [post]
func (h *ExtendedHandler) CreateServicedependency(c *echo.Context) error {
	var row models.Servicedependency
	if err := c.Bind(&row); err != nil {
		return BadRequest(c, "invalid body")
	}
	if row.ConfigName == "" {
		return BadRequest(c, "config_name is required")
	}
	row.LastModified = time.Now()
	if row.Active == "" {
		row.Active = "1"
	}
	if row.Register == "" {
		row.Register = "1"
	}
	if err := h.db.Create(&row).Error; err != nil {
		return Conflict(c, "config_name already exists")
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "create", "servicedependency", row.ConfigName, fmt.Sprintf("id=%d", row.ID))
	return c.JSON(http.StatusCreated, row)
}

// UpdateServicedependency godoc
// @Summary      Update a service dependency
// @Tags         servicedependencies
// @Accept       json
// @Produce      json
// @Param        id    path  int                       true  "Service dependency ID"
// @Param        body  body  models.Servicedependency  true  "Servicedependency"
// @Security     BearerAuth
// @Success      200  {object}  models.Servicedependency
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /servicedependencies/{id} [put]
func (h *ExtendedHandler) UpdateServicedependency(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var existing models.Servicedependency
	if err := h.db.First(&existing, id).Error; err != nil {
		return NotFound(c)
	}
	var upd models.Servicedependency
	if err := c.Bind(&upd); err != nil {
		return BadRequest(c, "invalid body")
	}
	upd.ID = id
	upd.LastModified = time.Now()
	if err := h.db.Save(&upd).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "update", "servicedependency", upd.ConfigName, fmt.Sprintf("id=%d", id))
	return c.JSON(http.StatusOK, upd)
}

// DeleteServicedependency godoc
// @Summary      Delete a service dependency
// @Tags         servicedependencies
// @Param        id  path  int  true  "Service dependency ID"
// @Security     BearerAuth
// @Success      204
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /servicedependencies/{id} [delete]
func (h *ExtendedHandler) DeleteServicedependency(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var row models.Servicedependency
	if err := h.db.First(&row, id).Error; err != nil {
		return NotFound(c)
	}
	if err := h.db.Delete(&row).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "delete", "servicedependency", row.ConfigName, fmt.Sprintf("id=%d", id))
	return c.NoContent(http.StatusNoContent)
}

// ── Service escalations ──────────────────────────────────────────────────────

// ListServiceescalations godoc
// @Summary      List service escalations
// @Tags         serviceescalations
// @Produce      json
// @Param        page   query  int  false  "Page"
// @Param        limit  query  int  false  "Items per page"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Failure      500  {object}  map[string]string
// @Router       /serviceescalations [get]
func (h *ExtendedHandler) ListServiceescalations(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "config_name")
	var total int64
	h.db.Model(&models.Serviceescalation{}).Count(&total)
	var rows []models.Serviceescalation
	if err := h.db.Order(fmt.Sprintf("%s %s", sp.Field, sp.Dir)).
		Offset((pp.Page - 1) * pp.Limit).Limit(pp.Limit).Find(&rows).Error; err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, ListResponse{Data: rows, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}

// GetServiceescalation godoc
// @Summary      Get a service escalation by ID
// @Tags         serviceescalations
// @Produce      json
// @Param        id  path  int  true  "Service escalation ID"
// @Security     BearerAuth
// @Success      200  {object}  models.Serviceescalation
// @Failure      404  {object}  map[string]string
// @Router       /serviceescalations/{id} [get]
func (h *ExtendedHandler) GetServiceescalation(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var row models.Serviceescalation
	if err := h.db.First(&row, id).Error; err != nil {
		return NotFound(c)
	}
	return c.JSON(http.StatusOK, row)
}

// CreateServiceescalation godoc
// @Summary      Create a service escalation
// @Tags         serviceescalations
// @Accept       json
// @Produce      json
// @Param        body  body  models.Serviceescalation  true  "Serviceescalation"
// @Security     BearerAuth
// @Success      201  {object}  models.Serviceescalation
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /serviceescalations [post]
func (h *ExtendedHandler) CreateServiceescalation(c *echo.Context) error {
	var row models.Serviceescalation
	if err := c.Bind(&row); err != nil {
		return BadRequest(c, "invalid body")
	}
	if row.ConfigName == "" {
		return BadRequest(c, "config_name is required")
	}
	row.LastModified = time.Now()
	if row.Active == "" {
		row.Active = "1"
	}
	if row.Register == "" {
		row.Register = "1"
	}
	if err := h.db.Create(&row).Error; err != nil {
		return Conflict(c, "config_name already exists")
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "create", "serviceescalation", row.ConfigName, fmt.Sprintf("id=%d", row.ID))
	return c.JSON(http.StatusCreated, row)
}

// UpdateServiceescalation godoc
// @Summary      Update a service escalation
// @Tags         serviceescalations
// @Accept       json
// @Produce      json
// @Param        id    path  int                       true  "Service escalation ID"
// @Param        body  body  models.Serviceescalation  true  "Serviceescalation"
// @Security     BearerAuth
// @Success      200  {object}  models.Serviceescalation
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /serviceescalations/{id} [put]
func (h *ExtendedHandler) UpdateServiceescalation(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var existing models.Serviceescalation
	if err := h.db.First(&existing, id).Error; err != nil {
		return NotFound(c)
	}
	var upd models.Serviceescalation
	if err := c.Bind(&upd); err != nil {
		return BadRequest(c, "invalid body")
	}
	upd.ID = id
	upd.LastModified = time.Now()
	if err := h.db.Save(&upd).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "update", "serviceescalation", upd.ConfigName, fmt.Sprintf("id=%d", id))
	return c.JSON(http.StatusOK, upd)
}

// DeleteServiceescalation godoc
// @Summary      Delete a service escalation
// @Tags         serviceescalations
// @Param        id  path  int  true  "Service escalation ID"
// @Security     BearerAuth
// @Success      204
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /serviceescalations/{id} [delete]
func (h *ExtendedHandler) DeleteServiceescalation(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var row models.Serviceescalation
	if err := h.db.First(&row, id).Error; err != nil {
		return NotFound(c)
	}
	if err := h.db.Delete(&row).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "delete", "serviceescalation", row.ConfigName, fmt.Sprintf("id=%d", id))
	return c.NoContent(http.StatusNoContent)
}

// ── Service extended info ────────────────────────────────────────────────────

// ListServiceextinfo godoc
// @Summary      List service extended info entries
// @Tags         serviceextinfo
// @Produce      json
// @Param        page   query  int  false  "Page"
// @Param        limit  query  int  false  "Items per page"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Failure      500  {object}  map[string]string
// @Router       /serviceextinfo [get]
func (h *ExtendedHandler) ListServiceextinfo(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "host_name")
	var total int64
	h.db.Model(&models.Serviceextinfo{}).Count(&total)
	var rows []models.Serviceextinfo
	if err := h.db.Order(fmt.Sprintf("%s %s", sp.Field, sp.Dir)).
		Offset((pp.Page - 1) * pp.Limit).Limit(pp.Limit).Find(&rows).Error; err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, ListResponse{Data: rows, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}

// GetServiceextinfo godoc
// @Summary      Get a service extended info entry by ID
// @Tags         serviceextinfo
// @Produce      json
// @Param        id  path  int  true  "Serviceextinfo ID"
// @Security     BearerAuth
// @Success      200  {object}  models.Serviceextinfo
// @Failure      404  {object}  map[string]string
// @Router       /serviceextinfo/{id} [get]
func (h *ExtendedHandler) GetServiceextinfo(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var row models.Serviceextinfo
	if err := h.db.First(&row, id).Error; err != nil {
		return NotFound(c)
	}
	return c.JSON(http.StatusOK, row)
}

// CreateServiceextinfo godoc
// @Summary      Create a service extended info entry
// @Tags         serviceextinfo
// @Accept       json
// @Produce      json
// @Param        body  body  models.Serviceextinfo  true  "Serviceextinfo"
// @Security     BearerAuth
// @Success      201  {object}  models.Serviceextinfo
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /serviceextinfo [post]
func (h *ExtendedHandler) CreateServiceextinfo(c *echo.Context) error {
	var row models.Serviceextinfo
	if err := c.Bind(&row); err != nil {
		return BadRequest(c, "invalid body")
	}
	if row.HostName == 0 {
		return BadRequest(c, "host_name (host id) is required")
	}
	if row.ServiceDescription == 0 {
		return BadRequest(c, "service_description (service id) is required")
	}
	row.LastModified = time.Now()
	if row.Active == "" {
		row.Active = "1"
	}
	if row.Register == "" {
		row.Register = "1"
	}
	if err := h.db.Create(&row).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "create", "serviceextinfo",
		fmt.Sprintf("host_id=%d svc_id=%d", row.HostName, row.ServiceDescription), fmt.Sprintf("id=%d", row.ID))
	return c.JSON(http.StatusCreated, row)
}

// UpdateServiceextinfo godoc
// @Summary      Update a service extended info entry
// @Tags         serviceextinfo
// @Accept       json
// @Produce      json
// @Param        id    path  int                    true  "Serviceextinfo ID"
// @Param        body  body  models.Serviceextinfo  true  "Serviceextinfo"
// @Security     BearerAuth
// @Success      200  {object}  models.Serviceextinfo
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /serviceextinfo/{id} [put]
func (h *ExtendedHandler) UpdateServiceextinfo(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var existing models.Serviceextinfo
	if err := h.db.First(&existing, id).Error; err != nil {
		return NotFound(c)
	}
	var upd models.Serviceextinfo
	if err := c.Bind(&upd); err != nil {
		return BadRequest(c, "invalid body")
	}
	upd.ID = id
	upd.LastModified = time.Now()
	if err := h.db.Save(&upd).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "update", "serviceextinfo",
		fmt.Sprintf("host_id=%d svc_id=%d", upd.HostName, upd.ServiceDescription), fmt.Sprintf("id=%d", id))
	return c.JSON(http.StatusOK, upd)
}

// DeleteServiceextinfo godoc
// @Summary      Delete a service extended info entry
// @Tags         serviceextinfo
// @Param        id  path  int  true  "Serviceextinfo ID"
// @Security     BearerAuth
// @Success      204
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /serviceextinfo/{id} [delete]
func (h *ExtendedHandler) DeleteServiceextinfo(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var row models.Serviceextinfo
	if err := h.db.First(&row, id).Error; err != nil {
		return NotFound(c)
	}
	if err := h.db.Delete(&row).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "delete", "serviceextinfo",
		fmt.Sprintf("host_id=%d svc_id=%d", row.HostName, row.ServiceDescription), fmt.Sprintf("id=%d", id))
	return c.NoContent(http.StatusNoContent)
}
