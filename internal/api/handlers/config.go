package handlers

import (
	"net/http"

	"github.com/jniltinho/go-nagiosql/internal/config"
	"github.com/jniltinho/go-nagiosql/internal/models"
	"github.com/jniltinho/go-nagiosql/internal/services/nagconfig"
	"github.com/jniltinho/go-nagiosql/internal/services/nagios"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// ConfigHandler handles Nagios config write/verify/restart endpoints.
type ConfigHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewConfigHandler creates a ConfigHandler.
func NewConfigHandler(db *gorm.DB, cfg *config.Config) *ConfigHandler {
	return &ConfigHandler{db: db, cfg: cfg}
}

func (h *ConfigHandler) getTarget() models.Configtarget {
	var ct models.Configtarget
	h.db.Where("domain_id = 0").First(&ct)
	return ct
}

// WriteAll godoc
// @Summary      Write all Nagios config files
// @Tags         config
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]any
// @Router       /config/write [post]
func (h *ConfigHandler) WriteAll(c *echo.Context) error {
	ct := h.getTarget()
	gen := nagconfig.New(h.db, ct.HostPath, ct.ServicePath, ct.BackupPath)
	written, err := gen.WriteAll()
	if err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"written": written, "errors": []string{}})
}

// Verify godoc
// @Summary      Verify Nagios config with nagios -v
// @Tags         config
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  nagios.VerifyResult
// @Router       /config/verify [post]
func (h *ConfigHandler) Verify(c *echo.Context) error {
	ct := h.getTarget()
	result := nagios.Verify(ct.NagiosBin, ct.NagiosCfg)
	return c.JSON(http.StatusOK, result)
}

// Restart godoc
// @Summary      Verify config then trigger graceful Nagios reload
// @Tags         config
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]any
// @Router       /config/restart [post]
func (h *ConfigHandler) Restart(c *echo.Context) error {
	ct := h.getTarget()
	result, err := nagios.Restart(ct.NagiosBin, ct.NagiosCfg, ct.CommandFile)
	if err != nil {
		return InternalError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"valid":     result.Valid,
		"output":    result.Output,
		"restarted": result.Valid && err == nil,
	})
}
