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

// CommandHandler handles /api/v1/commands endpoints.
type CommandHandler struct{ db *gorm.DB }

// NewCommandHandler creates a CommandHandler.
func NewCommandHandler(db *gorm.DB) *CommandHandler { return &CommandHandler{db: db} }

// ListCommands godoc
// @Summary      List commands
// @Tags         commands
// @Produce      json
// @Param        type   query  string  false  "Filter: check|notify (maps to command_type 0|1)"
// @Param        page   query  int     false  "Page"
// @Param        limit  query  int     false  "Limit"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Router       /commands [get]
func (h *CommandHandler) List(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "command_name")

	q := h.db.Model(&models.Command{})
	switch c.QueryParam("type") {
	case "check":
		q = q.Where("command_type = 0")
	case "notify":
		q = q.Where("command_type = 1")
	}

	var total int64
	q.Count(&total)
	var cmds []models.Command
	if err := q.Order(fmt.Sprintf("%s %s", sp.Field, sp.Dir)).
		Offset((pp.Page - 1) * pp.Limit).Limit(pp.Limit).Find(&cmds).Error; err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, ListResponse{Data: cmds, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}

// GetCommand godoc
// @Summary      Get command by ID
// @Tags         commands
// @Param        id  path  int  true  "Command ID"
// @Security     BearerAuth
// @Success      200  {object}  models.Command
// @Failure      404  {object}  map[string]string
// @Router       /commands/{id} [get]
func (h *CommandHandler) Get(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var cmd models.Command
	if err := h.db.First(&cmd, id).Error; err != nil {
		return NotFound(c)
	}
	return c.JSON(http.StatusOK, cmd)
}

// CreateCommand godoc
// @Summary      Create a command
// @Tags         commands
// @Accept       json
// @Produce      json
// @Param        body  body  models.Command  true  "Command"
// @Security     BearerAuth
// @Success      201  {object}  models.Command
// @Router       /commands [post]
func (h *CommandHandler) Create(c *echo.Context) error {
	var cmd models.Command
	if err := c.Bind(&cmd); err != nil {
		return BadRequest(c, "invalid body")
	}
	if cmd.CommandName == "" || cmd.CommandLine == "" {
		return BadRequest(c, "command_name and command_line are required")
	}
	cmd.LastModified = time.Now()
	if cmd.Active == "" {
		cmd.Active = "1"
	}
	if cmd.Register == "" {
		cmd.Register = "1"
	}
	if err := h.db.Create(&cmd).Error; err != nil {
		return Conflict(c, "command_name already exists")
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "create", "command", cmd.CommandName, fmt.Sprintf("id=%d", cmd.ID))
	return c.JSON(http.StatusCreated, cmd)
}

// UpdateCommand godoc
// @Summary      Update a command
// @Tags         commands
// @Accept       json
// @Produce      json
// @Param        id    path  int             true  "Command ID"
// @Param        body  body  models.Command  true  "Command fields"
// @Security     BearerAuth
// @Success      200  {object}  models.Command
// @Router       /commands/{id} [put]
func (h *CommandHandler) Update(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var existing models.Command
	if err := h.db.First(&existing, id).Error; err != nil {
		return NotFound(c)
	}
	var upd models.Command
	if err := c.Bind(&upd); err != nil {
		return BadRequest(c, "invalid body")
	}
	upd.ID = id
	upd.LastModified = time.Now()
	if err := h.db.Save(&upd).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "update", "command", upd.CommandName, fmt.Sprintf("id=%d", id))
	return c.JSON(http.StatusOK, upd)
}

// DeleteCommand godoc
// @Summary      Delete a command
// @Tags         commands
// @Param        id  path  int  true  "Command ID"
// @Security     BearerAuth
// @Success      204  "No Content"
// @Router       /commands/{id} [delete]
func (h *CommandHandler) Delete(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var cmd models.Command
	if err := h.db.First(&cmd, id).Error; err != nil {
		return NotFound(c)
	}
	if err := h.db.Delete(&cmd).Error; err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "delete", "command", cmd.CommandName, fmt.Sprintf("id=%d", id))
	return c.NoContent(http.StatusNoContent)
}
