package db

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDb(t *testing.T) {
	testDb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	db = testDb
	testDb.AutoMigrate(
		&User{},
		&BookWork{},
		&BookCopy{},
		&Loan{},
		&Hold{},
	)
}
