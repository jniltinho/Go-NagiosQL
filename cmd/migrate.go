package cmd

import (
	"fmt"
	"log"

	"github.com/jniltinho/go-nagiosql/internal/config"
	"github.com/jniltinho/go-nagiosql/internal/db"
	"github.com/jniltinho/go-nagiosql/internal/db/migrations"
	"github.com/jniltinho/go-nagiosql/internal/db/seeds"
	"github.com/spf13/cobra"
)

var (
	migrateWithSample bool
	migrateAdminUser  string
	migrateAdminPass  string
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations and seed required data",
	Long: `migrate creates or updates all NagiosQL database tables via GORM AutoMigrate
and seeds required initial data: default data domains, config target paths,
application settings, and the initial admin user (with bcrypt password).

Use --sample to also load the standard Nagios sample objects (24 commands,
5 time periods, templates, 4 example hosts, and 21 services) matching the
content of the original import_nagios_sample.sql.

The command is idempotent: safe to run multiple times.`,
	RunE: runMigrate,
}

func init() {
	migrateCmd.Flags().BoolVar(&migrateWithSample, "sample", false, "also seed Nagios sample objects")
	migrateCmd.Flags().StringVar(&migrateAdminUser, "admin-user", "admin", "initial admin username")
	migrateCmd.Flags().StringVar(&migrateAdminPass, "admin-password", "", "initial admin password (required)")
	_ = migrateCmd.MarkFlagRequired("admin-password")
	rootCmd.AddCommand(migrateCmd)
}

func runMigrate(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	database, err := db.Open(cfg)
	if err != nil {
		return err
	}

	log.Println("running database migrations...")
	if err := migrations.Migrate(database); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	log.Println("migrations complete")

	log.Println("seeding required data...")
	if err := seeds.SeedRequired(database, cfg, migrateAdminUser, migrateAdminPass); err != nil {
		return fmt.Errorf("seed required: %w", err)
	}
	log.Println("required seed complete")

	if migrateWithSample {
		log.Println("seeding sample Nagios objects...")
		if err := seeds.SeedSample(database); err != nil {
			return fmt.Errorf("seed sample: %w", err)
		}
		log.Println("sample seed complete")
	}

	log.Println("migrate: all done")
	return nil
}
