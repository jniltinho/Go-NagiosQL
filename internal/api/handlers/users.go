package handlers

import (
	"fmt"
	"net/http"
	"time"

	apimw "go-nagiosql/internal/api/middleware"
	"go-nagiosql/internal/models"
	"go-nagiosql/internal/services/auth"
	"go-nagiosql/internal/services/logbook"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// UserHandler handles /api/v1/users endpoints (admin-only).
type UserHandler struct {
	db      *gorm.DB
	authSvc *auth.Service
}

// NewUserHandler creates a UserHandler.
func NewUserHandler(db *gorm.DB, authSvc *auth.Service) *UserHandler {
	return &UserHandler{db: db, authSvc: authSvc}
}

// List godoc
// @Summary      List users (admin only)
// @Tags         users
// @Produce      json
// @Param        page   query  int  false  "Page"
// @Param        limit  query  int  false  "Items per page"
// @Security     BearerAuth
// @Success      200  {object}  ListResponse
// @Router       /users [get]
func (h *UserHandler) List(c *echo.Context) error {
	pp := ParsePage(c)
	sp := ParseSort(c, "username")
	var total int64
	h.db.Model(&models.User{}).Count(&total)
	var users []models.User
	if err := h.db.Order(fmt.Sprintf("%s %s", sp.Field, sp.Dir)).
		Offset((pp.Page - 1) * pp.Limit).Limit(pp.Limit).Find(&users).Error; err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, ListResponse{Data: users, Total: int(total), Page: pp.Page, Limit: pp.Limit})
}

// Get godoc
// @Summary      Get a user by ID (admin only)
// @Tags         users
// @Produce      json
// @Param        id  path  int  true  "User ID"
// @Security     BearerAuth
// @Success      200  {object}  models.User
// @Failure      404  {object}  map[string]string
// @Router       /users/{id} [get]
func (h *UserHandler) Get(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		return NotFound(c)
	}
	return c.JSON(http.StatusOK, user)
}

// createUserRequest wraps models.User and adds a plain-text password field.
type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Admin    string `json:"admin"`
	DomainID uint   `json:"domain_id"`
	Active   string `json:"active"`
}

// Create godoc
// @Summary      Create a user (admin only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body  body  createUserRequest  true  "User"
// @Security     BearerAuth
// @Success      201  {object}  models.User
// @Router       /users [post]
func (h *UserHandler) Create(c *echo.Context) error {
	var req createUserRequest
	if err := c.Bind(&req); err != nil {
		return BadRequest(c, "invalid body")
	}
	if req.Username == "" || req.Password == "" {
		return BadRequest(c, "username and password are required")
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return InternalError(c, err)
	}
	user := models.User{
		Username:     req.Username,
		Password:     hash,
		Name:         req.Name,
		Email:        req.Email,
		Admin:        req.Admin,
		DomainID:     req.DomainID,
		Active:       "1",
		LastModified: time.Now(),
	}
	if req.Active != "" {
		user.Active = req.Active
	}
	if err := h.db.Create(&user).Error; err != nil {
		return Conflict(c, "username already exists")
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "create", "user", user.Username, fmt.Sprintf("id=%d", user.ID))
	return c.JSON(http.StatusCreated, user)
}

// changePasswordRequest is used for PUT /api/v1/users/:id/password.
type changePasswordRequest struct {
	NewPassword string `json:"new_password"`
}

// ChangePassword godoc
// @Summary      Change a user's password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path  int                    true  "User ID"
// @Param        body  body  changePasswordRequest  true  "New password"
// @Security     BearerAuth
// @Success      200  {object}  map[string]string
// @Router       /users/{id}/password [put]
func (h *UserHandler) ChangePassword(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}
	var req changePasswordRequest
	if err := c.Bind(&req); err != nil || req.NewPassword == "" {
		return BadRequest(c, "new_password is required")
	}
	if err := h.authSvc.UpgradePassword(id, req.NewPassword); err != nil {
		return InternalError(c, err)
	}
	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "change_password", "user", "", fmt.Sprintf("id=%d", id))
	return c.JSON(http.StatusOK, map[string]string{"status": "password updated"})
}

// Delete godoc
// @Summary      Delete a user (admin only)
// @Tags         users
// @Param        id  path  int  true  "User ID"
// @Security     BearerAuth
// @Success      204
// @Failure      409  {object}  map[string]string  "cannot delete own account"
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(c *echo.Context) error {
	id, ok := UintParam(c, "id")
	if !ok {
		return BadRequest(c, "invalid id")
	}

	// Prevent self-deletion.
	claims := apimw.ClaimsFromContext(c)
	var self models.User
	h.db.Where("username = ?", claims.Username).First(&self)
	if self.ID == id {
		return Conflict(c, "cannot delete your own account")
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		return NotFound(c)
	}
	if err := h.db.Delete(&user).Error; err != nil {
		return InternalError(c, err)
	}
	logbook.Write(h.db, claims.DomainID, claims.Username, "delete", "user", user.Username, fmt.Sprintf("id=%d", id))
	return c.NoContent(http.StatusNoContent)
}
