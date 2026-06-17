// Package migrations runs GORM AutoMigrate for all NagiosQL tables.
// Tables are created with the server default storage engine (InnoDB on MariaDB 10.x),
// which supports transactions, foreign keys, and is portable across database backends.
package migrations

import (
	"fmt"

	"github.com/jniltinho/go-nagiosql/internal/models"
	"gorm.io/gorm"
)

// Migrate creates or updates every NagiosQL table via AutoMigrate.
// It is idempotent and safe to run on an already-migrated database.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(models.AllModels()...); err != nil {
		return fmt.Errorf("AutoMigrate: %w", err)
	}
	return nil
}
