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

// ContactHandler handles /api/v1/contacts endpoints.
type ContactHandler struct{ db *gorm.DB }

// NewContactHandler creates a ContactHandler.
func NewContactHandler(db *gorm.DB) *ContactHandler { return &ContactHandler{db: db} }

// List godoc
// @Summary      List contacts
// @Tags         contacts
// @Produce      json
// @Param        page   query  int  false  "Page"
// @Param        limit  query  int  false  "Items per page"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Router       /contacts [get]
func (h *ContactHandler) List(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "contact_name")
	var total int64
	h.db.Model(&models.Contact{}).Count(&total)
	var contacts []models.Contact
	if err := h.db.Order(fmt.Sprintf("%s %s", sp.Field, sp.Dir)).
		Offset((pp.Page - 1) * pp.Limit).Limit(pp.Limit).Find(&contacts).Error; err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, ListResponse{Data: contacts, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}

// Get godoc
// @Summary      Get a contact by ID
// @Tags         contacts
// @Produce      json
// @Param        id  path  int  true  "Contact ID"
// @Security     BearerAuth
// @Success      200  {object}  models.Contact
// @Failure      404  {object}  map[string]string
// @Router       /contacts/{id} [get]
func (h *ContactHandler) Get(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var contact models.Contact
	if err := h.db.First(&contact, id).Error; err != nil {
		return NotFound(c)
	}
	return c.JSON(http.StatusOK, contact)
}

// Create godoc
// @Summary      Create a contact
// @Tags         contacts
// @Accept       json
// @Produce      json
// @Param        body  body  models.Contact  true  "Contact"
// @Security     BearerAuth
// @Success      201  {object}  models.Contact
// @Router       /contacts [post]
func (h *ContactHandler) Create(c *echo.Context) error {
	var contact models.Contact
	if err := c.Bind(&contact); err != nil {
		return BadRequest(c, "invalid body")
	}
	if contact.ContactName == "" {
		return BadRequest(c, "contact_name is required")
	}
	contact.LastModified = time.Now()
	if contact.Active == "" {
		contact.Active = "1"
	}
	if contact.Register == "" {
		contact.Register = "1"
	}
	if err := h.db.Create(&contact).Error; err != nil {
		return Conflict(c, "contact_name already exists")
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "create", "contact", contact.ContactName, fmt.Sprintf("id=%d", contact.ID))
	return c.JSON(http.StatusCreated, contact)
}

// Update godoc
// @Summary      Update a contact
// @Tags         contacts
// @Accept       json
// @Produce      json
// @Param        id    path  int             true  "Contact ID"
// @Param        body  body  models.Contact  true  "Contact"
// @Security     BearerAuth
// @Success      200  {object}  models.Contact
// @Router       /contacts/{id} [put]
func (h *ContactHandler) Update(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var existing models.Contact
	if err := h.db.First(&existing, id).Error; err != nil {
		return NotFound(c)
	}
	var upd models.Contact
	if err := c.Bind(&upd); err != nil {
		return BadRequest(c, "invalid body")
	}
	upd.ID = id
	upd.LastModified = time.Now()
	if err := h.db.Save(&upd).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "update", "contact", upd.ContactName, fmt.Sprintf("id=%d", id))
	return c.JSON(http.StatusOK, upd)
}

// Delete godoc
// @Summary      Delete a contact
// @Tags         contacts
// @Param        id  path  int  true  "Contact ID"
// @Security     BearerAuth
// @Success      204
// @Failure      404  {object}  map[string]string
// @Router       /contacts/{id} [delete]
func (h *ContactHandler) Delete(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var contact models.Contact
	if err := h.db.First(&contact, id).Error; err != nil {
		return NotFound(c)
	}
	if err := h.db.Delete(&contact).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "delete", "contact", contact.ContactName, fmt.Sprintf("id=%d", id))
	return c.NoContent(http.StatusNoContent)
}
