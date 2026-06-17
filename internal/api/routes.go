// Package api wires all HTTP routes onto the Echo instance.
package api

import (
	"net/http"

	"github.com/jniltinho/go-nagiosql/internal/api/handlers"
	apimw "github.com/jniltinho/go-nagiosql/internal/api/middleware"
	"github.com/jniltinho/go-nagiosql/internal/config"
	"github.com/jniltinho/go-nagiosql/internal/services/auth"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// RegisterRoutes mounts all API route groups onto e.
func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	authSvc := auth.New(db, cfg)

	// Health check — no auth.
	e.GET("/healthz", func(c *echo.Context) error {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "db_error"})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Public auth routes.
	authH := handlers.NewAuthHandler(authSvc, cfg.JWT.RefreshTTLDays)
	ag := e.Group("/api/v1/auth")
	ag.POST("/login", authH.Login)
	ag.POST("/refresh", authH.Refresh)
	ag.POST("/logout", authH.Logout)

	// All other routes require JWT.
	jwtMW := apimw.JWTAuth(authSvc)
	v1 := e.Group("/api/v1", jwtMW)

	// Identity.
	v1.GET("/me", func(c *echo.Context) error {
		cl := apimw.ClaimsFromContext(c)
		return c.JSON(http.StatusOK, map[string]any{
			"username":  cl.Username,
			"admin":     cl.Admin,
			"domain_id": cl.DomainID,
		})
	})

	// Monitoring dashboard.
	mon := handlers.NewMonitoringHandler(db)
	v1.GET("/monitoring/summary", mon.GetSummary)

	// Hosts.
	hh := handlers.NewHostHandler(db)
	v1.GET("/hosts", hh.List)
	v1.POST("/hosts", hh.Create)
	v1.GET("/hosts/:id", hh.Get)
	v1.PUT("/hosts/:id", hh.Update)
	v1.DELETE("/hosts/:id", hh.Delete)

	// Services.
	sh := handlers.NewServiceHandler(db)
	v1.GET("/services", sh.List)
	v1.POST("/services", sh.Create)
	v1.GET("/services/:id", sh.Get)
	v1.PUT("/services/:id", sh.Update)
	v1.DELETE("/services/:id", sh.Delete)

	// Commands.
	ch := handlers.NewCommandHandler(db)
	v1.GET("/commands", ch.List)
	v1.POST("/commands", ch.Create)
	v1.GET("/commands/:id", ch.Get)
	v1.PUT("/commands/:id", ch.Update)
	v1.DELETE("/commands/:id", ch.Delete)

	// Timeperiods.
	th := handlers.NewTimeperiodHandler(db)
	v1.GET("/timeperiods", th.List)
	v1.POST("/timeperiods", th.Create)
	v1.GET("/timeperiods/:id", th.Get)
	v1.PUT("/timeperiods/:id", th.Update)
	v1.DELETE("/timeperiods/:id", th.Delete)

	// Contacts.
	coh := handlers.NewContactHandler(db)
	v1.GET("/contacts", coh.List)
	v1.POST("/contacts", coh.Create)
	v1.GET("/contacts/:id", coh.Get)
	v1.PUT("/contacts/:id", coh.Update)
	v1.DELETE("/contacts/:id", coh.Delete)

	// Groups.
	gh := handlers.NewGroupHandler(db)
	v1.GET("/hostgroups", gh.ListHostgroups)
	v1.POST("/hostgroups", gh.CreateHostgroup)
	v1.GET("/hostgroups/:id", gh.GetHostgroup)
	v1.PUT("/hostgroups/:id", gh.UpdateHostgroup)
	v1.DELETE("/hostgroups/:id", gh.DeleteHostgroup)
	v1.PUT("/hostgroups/:id/members", gh.AddHostgroupMember)

	v1.GET("/servicegroups", gh.ListServicegroups)
	v1.POST("/servicegroups", gh.CreateServicegroup)
	v1.DELETE("/servicegroups/:id", gh.DeleteServicegroup)

	v1.GET("/contactgroups", gh.ListContactgroups)
	v1.POST("/contactgroups", gh.CreateContactgroup)
	v1.DELETE("/contactgroups/:id", gh.DeleteContactgroup)

	// Templates.
	tmpl := handlers.NewTemplateHandler(db)
	v1.GET("/hosttemplates", tmpl.ListHosttemplates)
	v1.POST("/hosttemplates", tmpl.CreateHosttemplate)
	v1.GET("/hosttemplates/:id", tmpl.GetHosttemplate)
	v1.DELETE("/hosttemplates/:id", tmpl.DeleteHosttemplate)

	v1.GET("/servicetemplates", tmpl.ListServicetemplates)
	v1.POST("/servicetemplates", tmpl.CreateServicetemplate)
	v1.DELETE("/servicetemplates/:id", tmpl.DeleteServicetemplate)

	v1.GET("/contacttemplates", tmpl.ListContacttemplates)
	v1.POST("/contacttemplates", tmpl.CreateContacttemplate)
	v1.DELETE("/contacttemplates/:id", tmpl.DeleteContacttemplate)

	// Users (admin-only).
	uh := handlers.NewUserHandler(db, authSvc)
	v1.GET("/users", uh.List, apimw.RequireAdmin())
	v1.POST("/users", uh.Create, apimw.RequireAdmin())
	v1.GET("/users/:id", uh.Get, apimw.RequireAdmin())
	v1.PUT("/users/:id/password", uh.ChangePassword)
	v1.DELETE("/users/:id", uh.Delete, apimw.RequireAdmin())

	// Logbook (read-only).
	lb := handlers.NewLogbookHandler(db)
	v1.GET("/logbook", lb.List)

	// Config write/verify/restart.
	cfgH := handlers.NewConfigHandler(db, cfg)
	v1.POST("/config/write", cfgH.WriteAll)
	v1.POST("/config/verify", cfgH.Verify)
	v1.POST("/config/restart", cfgH.Restart)

	// Import .cfg files.
	ih := handlers.NewImportHandler(db)
	v1.POST("/import", ih.Import)

	// Settings (admin-only write).
	seth := handlers.NewSettingsHandler(db)
	v1.GET("/settings", seth.Get)
	v1.PUT("/settings", seth.Update, apimw.RequireAdmin())

	// Variable definitions (custom _VAR macros).
	vh := handlers.NewVariableHandler(db)
	v1.GET("/variables", vh.List)
	v1.POST("/variables", vh.Create)
	v1.GET("/variables/:id", vh.Get)
	v1.PUT("/variables/:id", vh.Update)
	v1.DELETE("/variables/:id", vh.Delete)
}
