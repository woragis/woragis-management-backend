package testutil

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/woragis/management/backend/server/internal/models"
	"gorm.io/gorm"
)

func OpenSQLite(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(
		&models.LeetcodeVideo{},
		&models.LeetcodeChannelSettings{},
		&models.WhatsappMessageTemplate{},
		&models.ContentThumbnail{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}
