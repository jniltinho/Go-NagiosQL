package nagimport

import (
	"strconv"
	"time"

	"go-nagiosql/internal/models"
	"gorm.io/gorm"
)

// ImportObject persists one ParsedObject into the database.
// Returns (found/acted, wasNew, error).
// If overwrite=false and the object already exists, returns (false, false, nil).
func ImportObject(db *gorm.DB, obj ParsedObject, configID uint8, overwrite bool) (ok bool, wasNew bool, err error) {
	switch obj.Type {
	case "host":
		return importHost(db, obj.Fields, configID, overwrite)
	case "service":
		return importService(db, obj.Fields, configID, overwrite)
	case "command":
		return importCommand(db, obj.Fields, configID, overwrite)
	default:
		// Unknown type — skip silently.
		return false, false, nil
	}
}

func importHost(db *gorm.DB, fields map[string]string, configID uint8, overwrite bool) (bool, bool, error) {
	name := fields["host_name"]
	if name == "" {
		return false, false, nil
	}
	var existing models.Host
	found := db.Where("host_name = ?", name).First(&existing).Error == nil
	if found && !overwrite {
		return false, false, nil
	}
	h := models.Host{
		HostName:     name,
		Alias:        fields["alias"],
		Address:      fields["address"],
		DisplayName:  fields["display_name"],
		CheckCommand: fields["check_command"],
		Active:       "1",
		Register:     orDefault(fields["register"], "1"),
		ConfigID:     configID,
		LastModified: time.Now(),
	}
	if found {
		h.ID = existing.ID
		return true, false, db.Save(&h).Error
	}
	return true, true, db.Create(&h).Error
}

func importService(db *gorm.DB, fields map[string]string, configID uint8, overwrite bool) (bool, bool, error) {
	desc := fields["service_description"]
	if desc == "" {
		return false, false, nil
	}
	configName := orDefault(fields["host_name"], "imported")
	var existing models.Service
	found := db.Where("service_description = ? AND config_name = ?", desc, configName).First(&existing).Error == nil
	if found && !overwrite {
		return false, false, nil
	}

	s := models.Service{
		ServiceDescription: desc,
		ConfigName:         configName,
		CheckCommand:       fields["check_command"],
		Active:             "1",
		Register:           orDefault(fields["register"], "1"),
		ConfigID:           configID,
		LastModified:       time.Now(),
	}
	if v, _ := strconv.Atoi(fields["max_check_attempts"]); v > 0 {
		s.MaxCheckAttempts = &v
	}
	if v, _ := strconv.Atoi(fields["check_interval"]); v > 0 {
		s.CheckInterval = &v
	}
	if v, _ := strconv.Atoi(fields["retry_interval"]); v > 0 {
		s.RetryInterval = &v
	}
	if found {
		s.ID = existing.ID
		return true, false, db.Save(&s).Error
	}
	return true, true, db.Create(&s).Error
}

func importCommand(db *gorm.DB, fields map[string]string, configID uint8, overwrite bool) (bool, bool, error) {
	name := fields["command_name"]
	if name == "" {
		return false, false, nil
	}
	var existing models.Command
	found := db.Where("command_name = ?", name).First(&existing).Error == nil
	if found && !overwrite {
		return false, false, nil
	}
	cmd := models.Command{
		CommandName:  name,
		CommandLine:  fields["command_line"],
		Active:       "1",
		Register:     "1",
		ConfigID:     configID,
		LastModified: time.Now(),
	}
	if found {
		cmd.ID = existing.ID
		return true, false, db.Save(&cmd).Error
	}
	return true, true, db.Create(&cmd).Error
}

func orDefault(v, def string) string {
	if v == "" {
		return def
	}
	return v
}
