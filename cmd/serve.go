package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"

	"go-nagiosql/internal/api"
	"go-nagiosql/internal/config"
	"go-nagiosql/internal/db"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/spf13/cobra"
)

// DocsFS is set by main.go to expose the embedded Swagger docs directory.
var DocsFS fs.FS

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the NagiosQL HTTP API server",
	RunE:  runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	database, err := db.Open(cfg)
	if err != nil {
		return err
	}

	// Echo v5: Logger field is *slog.Logger.
	logLevel := slog.LevelWarn
	if cfg.Server.Dev || devMode {
		logLevel = slog.LevelDebug
	}
	slogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	e := echo.NewWithConfig(echo.Config{Logger: slogger})

	// Core middleware: panic recovery → request ID → request logging → CORS.
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLogger())
	// CORS: allow all origins (restrict in production by replacing "*" with a domain list).
	e.Use(middleware.CORS("*"))

	// Swagger UI: serve the embedded docs/ directory at /docs/
	if DocsFS != nil {
		sub, err := fs.Sub(DocsFS, "docs")
		if err == nil {
			e.GET("/docs/*", func(c *echo.Context) error {
				http.StripPrefix("/docs", http.FileServer(http.FS(sub))).ServeHTTP(c.Response(), c.Request())
				return nil
			})
			// Redirect /api/swagger to the UI HTML.
			e.GET("/api/swagger", func(c *echo.Context) error {
				return c.Redirect(http.StatusMovedPermanently, "/docs/swagger-ui.html")
			})
		}
	}

	// Register all API routes (auth, objects, config, import).
	api.RegisterRoutes(e, database, cfg)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("nagiosql %s starting on %s (dev=%v)", buildVersion, addr, cfg.Server.Dev || devMode)

	// e.Start handles OS signal listening and graceful shutdown internally.
	return e.Start(addr)
}
