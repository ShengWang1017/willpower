package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"willpower-forge-api/internal/models"
)

var dbInstance *gorm.DB

// Connect initializes and caches the SQLite connection used across the app.
func Connect() (*gorm.DB, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}

	db, err := gorm.Open(sqlite.Open("willpower.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return nil, err
	}

	dbInstance = db
	return dbInstance, nil
}

// AutoMigrateModels ensures the schema matches the expected models.
func AutoMigrateModels(db *gorm.DB) {
	if err := db.AutoMigrate(&models.User{}, &models.Goal{}, &models.CheckIn{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	if db.Migrator().HasIndex(&models.CheckIn{}, "idx_goal_date") {
		if err := db.Migrator().DropIndex(&models.CheckIn{}, "idx_goal_date"); err != nil {
			log.Printf("failed to drop legacy check-in index: %v", err)
		}
	}
}
