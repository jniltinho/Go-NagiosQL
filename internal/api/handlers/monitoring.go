package handlers

import (
	"net/http"

	"go-nagiosql/internal/models"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// MonitoringHandler handles dashboard summary endpoints.
type MonitoringHandler struct{ db *gorm.DB }

// NewMonitoringHandler creates a MonitoringHandler.
func NewMonitoringHandler(db *gorm.DB) *MonitoringHandler { return &MonitoringHandler{db: db} }

// summary holds per-object-type counts for the dashboard.
type summary struct {
	Hosts            int64 `json:"hosts"`
	ActiveHosts      int64 `json:"active_hosts"`
	Services         int64 `json:"services"`
	ActiveServices   int64 `json:"active_services"`
	Commands         int64 `json:"commands"`
	Timeperiods      int64 `json:"timeperiods"`
	Contacts         int64 `json:"contacts"`
	Contactgroups    int64 `json:"contactgroups"`
	Hostgroups       int64 `json:"hostgroups"`
	Servicegroups    int64 `json:"servicegroups"`
	Hosttemplates    int64 `json:"hosttemplates"`
	Servicetemplates int64 `json:"servicetemplates"`
}

// GetSummary godoc
// @Summary      Dashboard object counts
// @Tags         monitoring
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  summary
// @Router       /monitoring/summary [get]
func (h *MonitoringHandler) GetSummary(c *echo.Context) error {
	var s summary
	h.db.Model(&models.Host{}).Count(&s.Hosts)
	h.db.Model(&models.Host{}).Where("active = '1'").Count(&s.ActiveHosts)
	h.db.Model(&models.Service{}).Count(&s.Services)
	h.db.Model(&models.Service{}).Where("active = '1'").Count(&s.ActiveServices)
	h.db.Model(&models.Command{}).Count(&s.Commands)
	h.db.Model(&models.Timeperiod{}).Count(&s.Timeperiods)
	h.db.Model(&models.Contact{}).Count(&s.Contacts)
	h.db.Model(&models.Contactgroup{}).Count(&s.Contactgroups)
	h.db.Model(&models.Hostgroup{}).Count(&s.Hostgroups)
	h.db.Model(&models.Servicegroup{}).Count(&s.Servicegroups)
	h.db.Model(&models.Hosttemplate{}).Count(&s.Hosttemplates)
	h.db.Model(&models.Servicetemplate{}).Count(&s.Servicetemplates)
	return c.JSON(http.StatusOK, s)
}
