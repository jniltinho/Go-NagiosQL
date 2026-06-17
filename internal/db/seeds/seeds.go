// Package seeds provides idempotent database seed functions for NagiosQL.
//
//   - SeedRequired: seeds the minimum data needed to start the application
//     (data domains, config target paths, app settings, initial admin user).
//     Always called after migrate.
//
//   - SeedSample: seeds the standard Nagios sample objects from
//     import_nagios_sample.sql (commands, timeperiods, templates, hosts,
//     services). Only called when --sample is passed to nagiosql migrate.
package seeds

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"go-nagiosql/internal/config"
	"go-nagiosql/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedRequired seeds the minimum required rows. It is idempotent.
func SeedRequired(db *gorm.DB, cfg *config.Config, adminUser, adminPassword string) error {
	if err := seedDatadomains(db); err != nil {
		return err
	}
	if err := seedConfigtarget(db, cfg); err != nil {
		return err
	}
	if err := seedSettings(db); err != nil {
		return err
	}
	if err := seedAdmin(db, adminUser, adminPassword); err != nil {
		return err
	}
	return nil
}

// seedDatadomains creates the common domain (id=0) and one default domain (id=1).
// tbl_datadomain id=0 must be inserted via raw SQL because GORM skips zero-value PKs.
func seedDatadomains(db *gorm.DB) error {
	// Check if id=0 already exists.
	var count int64
	db.Raw("SELECT COUNT(*) FROM tbl_datadomain WHERE id = 0").Scan(&count)

	if count == 0 {
		now := time.Now()
		if err := db.Exec(
			"INSERT INTO tbl_datadomain (id, name, description, active, last_modified) VALUES (0, 'Common', 'Common data domain shared by all users', '1', ?)",
			now,
		).Error; err != nil {
			return fmt.Errorf("seed datadomain id=0: %w", err)
		}
		// Reset AUTO_INCREMENT so the next INSERT gets id=1 or higher.
		if err := db.Exec("ALTER TABLE tbl_datadomain AUTO_INCREMENT = 1").Error; err != nil {
			return fmt.Errorf("reset datadomain auto_increment: %w", err)
		}
		log.Println("seeded datadomain id=0 (common)")
	}

	// Default domain id=1.
	domain := models.Datadomain{
		Name:         "Default",
		Description:  "Default data domain",
		Active:       "1",
		LastModified: time.Now(),
	}
	result := db.Where(models.Datadomain{Name: "Default"}).FirstOrCreate(&domain)
	if result.Error != nil {
		return fmt.Errorf("seed datadomain Default: %w", result.Error)
	}
	return nil
}

// seedConfigtarget inserts the global config target row (domain_id=0) using
// paths from the application configuration. Updates the row if it already exists.
func seedConfigtarget(db *gorm.DB, cfg *config.Config) error {
	n := cfg.Nagios
	ct := models.Configtarget{
		DomainID:          0,
		NagiosCfg:         n.ConfigFile,
		CgiCfg:            n.CgiFile,
		ResourceCfg:       n.ResourceFile,
		BrokerModule:      "",
		CommandFile:       n.ReloadTrigger,
		HostPath:          n.HostConfigDir,
		ServicePath:       n.ServiceConfigDir,
		CheckPath:         n.BaseDir + "/libexec/",
		UserPath:          n.BaseDir + "/etc/nagiosql/users/",
		ImportPath:        n.ImportDir,
		BackupPath:        n.BackupDir,
		LogPath:           n.BaseDir + "/var/nagios.log",
		NagiosBin:         n.Binary,
		NagiosPID:         n.PidFile,
		TargetName:        "local",
		TargetDescription: "Local Nagios installation",
		WriteToFiles:      "1",
		WriteToDatabase:   "1",
		Active:            "1",
		LastModified:      time.Now(),
	}

	var existing models.Configtarget
	if err := db.Where("domain_id = 0").First(&existing).Error; err == gorm.ErrRecordNotFound {
		if err := db.Create(&ct).Error; err != nil {
			return fmt.Errorf("seed configtarget: %w", err)
		}
		log.Println("seeded configtarget for domain_id=0")
	} else if err == nil {
		// Update paths in case config.toml changed.
		db.Model(&existing).Updates(map[string]any{
			"nagios_cfg":    ct.NagiosCfg,
			"cgi_cfg":       ct.CgiCfg,
			"resource_cfg":  ct.ResourceCfg,
			"command_file":  ct.CommandFile,
			"host_path":     ct.HostPath,
			"service_path":  ct.ServicePath,
			"check_path":    ct.CheckPath,
			"import_path":   ct.ImportPath,
			"backup_path":   ct.BackupPath,
			"nagios_bin":    ct.NagiosBin,
			"nagios_pid":    ct.NagiosPID,
			"last_modified": ct.LastModified,
		})
	}
	return nil
}

// seedSettings inserts the single tbl_settings row if it doesn't exist.
// install_hash is a random hex string that identifies this installation.
func seedSettings(db *gorm.DB) error {
	var count int64
	db.Model(&models.Settings{}).Count(&count)
	if count > 0 {
		return nil
	}

	hashBytes := make([]byte, 32)
	if _, err := rand.Read(hashBytes); err != nil {
		return fmt.Errorf("generating install_hash: %w", err)
	}

	s := models.Settings{
		InstallHash:       hex.EncodeToString(hashBytes),
		BackupAge:         14,
		DateFormat:        "Y-m-d H:i:s",
		LogbookContent:    250,
		CryptoKey:         "",
		AuthenticationKey: "",
		LastModified:      time.Now(),
	}
	if err := db.Create(&s).Error; err != nil {
		return fmt.Errorf("seed settings: %w", err)
	}
	log.Println("seeded settings")
	return nil
}

// seedAdmin creates the initial admin user with a bcrypt-hashed password.
// Password is NEVER stored as MD5 — legacy MD5 detection is for login only.
func seedAdmin(db *gorm.DB, username, password string) error {
	var count int64
	db.Model(&models.User{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		log.Printf("admin user %q already exists, skipping", username)
		return nil
	}

	// bcrypt cost=12 per security policy.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return fmt.Errorf("hashing admin password: %w", err)
	}

	user := models.User{
		Username:     username,
		Password:     string(hash),
		Name:         "Administrator",
		Email:        "",
		Admin:        "1",
		DomainID:     0,
		LogonTimeout: 60,
		Active:       "1",
		LastModified: time.Now(),
	}
	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("seed admin user: %w", err)
	}
	log.Printf("seeded admin user %q (bcrypt)", username)
	return nil
}
