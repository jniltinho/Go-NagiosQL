// Package db provides the GORM database connection for NagiosQL.
package db

import (
	"fmt"
	"log"
	"time"

	"github.com/jniltinho/go-nagiosql/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open opens a GORM connection to MariaDB using the supplied configuration.
// InnoDB is used (MariaDB default) for transaction support and future portability.
// It fails fast if the database is unreachable — Docker healthchecks are expected
// to ensure MariaDB is ready before the binary starts.
func Open(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.Database.DSN()

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                      dsn,
		DefaultStringSize:        255,
		DisableDatetimePrecision: true,
		DontSupportRenameIndex:   true,
		DontSupportRenameColumn:  true,
	}), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("getting sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database at %s:%d: %w", cfg.Database.Host, cfg.Database.Port, err)
	}

	log.Printf("database connected: %s@%s:%d/%s",
		cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	return db, nil
}
