// Package logbook writes audit trail entries to tbl_logbook.
package logbook

import (
	"time"

	"go-nagiosql/internal/models"
	"gorm.io/gorm"
)

// Write appends one audit entry. Errors are intentionally swallowed so that a
// logbook failure never breaks the business operation that triggered it.
func Write(db *gorm.DB, userID uint, username, action, objectType, objectName, info string) {
	entry := models.Logbook{
		UserID:     userID,
		Username:   username,
		Action:     action,
		ObjectType: objectType,
		ObjectName: objectName,
		Info:       info,
		CreatedAt:  time.Now(),
	}
	db.Create(&entry) //nolint:errcheck
}
