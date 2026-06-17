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

// VariableHandler handles /api/v1/variables.
type VariableHandler struct{ db *gorm.DB }

// NewVariableHandler creates a VariableHandler.
func NewVariableHandler(db *gorm.DB) *VariableHandler { return &VariableHandler{db: db} }

type variableRequest struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Vartype  string `json:"vartype"`
	Active   string `json:"active"`
	DomainID uint   `json:"domain_id"`
	ConfigID uint8  `json:"config_id"`
}

// List godoc
// @Summary      List variable definitions
// @Tags         variables
// @Produce      json
// @Param        page   query  int    false  "Page"
// @Param        limit  query  int    false  "Items per page"
// @Param        name   query  string false  "Filter by name (partial)"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Router       /variables [get]
func (h *VariableHandler) List(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "name")
	offset := (pp.Page - 1) * pp.Limit

	q := h.db.Model(&models.Variabledefinition{})
	if name := c.QueryParam("name"); name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}

	var total int64
	q.Count(&total)

	var vars []models.Variabledefinition
	q.Order(sp.Field + " " + sp.Dir).Offset(offset).Limit(pp.Limit).Find(&vars)

	return c.JSON(http.StatusOK, ListResponse{Data: vars, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}

// Create godoc
// @Summary      Create a variable definition
// @Tags         variables
// @Accept       json
// @Produce      json
// @Param        body  body  variableRequest  true  "Variable"
// @Security     BearerAuth
// @Success      201  {object}  models.Variabledefinition
// @Router       /variables [post]
func (h *VariableHandler) Create(c *echo.Context) error {
	var req variableRequest
	if err := c.Bind(&req); err != nil {
		return BadRequest(c, "invalid body")
	}
	if req.Name == "" {
		return BadRequest(c, "name is required")
	}
	vartype := req.Vartype
	if vartype == "" {
		vartype = "string"
	}
	v := models.Variabledefinition{
		Name:         req.Name,
		Value:        req.Value,
		Vartype:      vartype,
		Active:       orStr(req.Active, "1"),
		DomainID:     req.DomainID,
		ConfigID:     req.ConfigID,
		LastModified: time.Now(),
	}
	if err := h.db.Create(&v).Error; err != nil {
		return Conflict(c, "variable name already exists or constraint violated")
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "create", "variable", v.Name, "")
	return c.JSON(http.StatusCreated, v)
}

// Get godoc
// @Summary      Get a variable definition by ID
// @Tags         variables
// @Produce      json
// @Param        id   path  int  true  "Variable ID"
// @Security     BearerAuth
// @Success      200  {object}  models.Variabledefinition
// @Failure      404  {object}  map[string]string
// @Router       /variables/{id} [get]
func (h *VariableHandler) Get(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var v models.Variabledefinition
	if err := h.db.First(&v, id).Error; err != nil {
		return NotFound(c)
	}
	return c.JSON(http.StatusOK, v)
}

// Update godoc
// @Summary      Update a variable definition
// @Tags         variables
// @Accept       json
// @Produce      json
// @Param        id    path  int              true  "Variable ID"
// @Param        body  body  variableRequest  true  "Variable"
// @Security     BearerAuth
// @Success      200  {object}  models.Variabledefinition
// @Router       /variables/{id} [put]
func (h *VariableHandler) Update(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var v models.Variabledefinition
	if err := h.db.First(&v, id).Error; err != nil {
		return NotFound(c)
	}
	var req variableRequest
	if err := c.Bind(&req); err != nil {
		return BadRequest(c, "invalid body")
	}
	if req.Name != "" {
		v.Name = req.Name
	}
	v.Value = req.Value
	v.Vartype = orStr(req.Vartype, v.Vartype)
	v.Active = orStr(req.Active, v.Active)
	v.LastModified = time.Now()

	if err := h.db.Save(&v).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "update", "variable", v.Name, fmt.Sprintf("id=%d", id))
	return c.JSON(http.StatusOK, v)
}

// Delete godoc
// @Summary      Delete a variable definition
// @Tags         variables
// @Param        id  path  int  true  "Variable ID"
// @Security     BearerAuth
// @Success      204
// @Failure      404  {object}  map[string]string
// @Router       /variables/{id} [delete]
func (h *VariableHandler) Delete(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var v models.Variabledefinition
	if err := h.db.First(&v, id).Error; err != nil {
		return NotFound(c)
	}
	if err := h.db.Delete(&v).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "delete", "variable", v.Name, fmt.Sprintf("id=%d", id))
	return c.NoContent(http.StatusNoContent)
}

func orStr(v, def string) string {
	if v == "" {
		return def
	}
	return v
}
