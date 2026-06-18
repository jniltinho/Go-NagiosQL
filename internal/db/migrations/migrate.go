// Package migrations runs GORM AutoMigrate for all NagiosQL tables.
// Tables are created with the server default storage engine (InnoDB on MariaDB 10.x),
// which supports transactions, foreign keys, and is portable across database backends.
package migrations

import (
	"fmt"
	"strings"

	"go-nagiosql/internal/models"
	"gorm.io/gorm"
)

// Migrate creates or updates every NagiosQL table via AutoMigrate.
// It is idempotent and safe to run on an already-migrated database.
func Migrate(db *gorm.DB) error {
	// Migrate each model individually so a "table already exists" from MariaDB
	// on a previously-migrated schema does not abort the whole migration.
	for _, model := range models.AllModels() {
		if err := db.AutoMigrate(model); err != nil {
			if strings.Contains(err.Error(), "already exists") {
				continue
			}
			return fmt.Errorf("AutoMigrate: %w", err)
		}
	}
	return nil
}
