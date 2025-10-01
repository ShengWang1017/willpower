package services

import (
	"log"
	"time"

	"gorm.io/gorm"
	"willpower-forge-api/internal/models"
)

type CleanupService struct {
	db *gorm.DB
}

func NewCleanupService(db *gorm.DB) *CleanupService {
	return &CleanupService{db: db}
}

// StartScheduledCleanup starts a background goroutine that periodically cleans up old deleted goals
func (s *CleanupService) StartScheduledCleanup() {
	go func() {
		ticker := time.NewTicker(24 * time.Hour) // Run once per day
		defer ticker.Stop()

		// Run immediately on start
		s.CleanupOldDeletedGoals()

		for range ticker.C {
			s.CleanupOldDeletedGoals()
		}
	}()
	log.Println("Scheduled cleanup service started")
}

// CleanupOldDeletedGoals permanently deletes goals that have been in recycle bin for more than 30 days
func (s *CleanupService) CleanupOldDeletedGoals() {
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	var goals []models.Goal
	result := s.db.Unscoped().
		Where("deleted_at IS NOT NULL AND deleted_at < ?", thirtyDaysAgo).
		Find(&goals)

	if result.Error != nil {
		log.Printf("Error finding old deleted goals: %v", result.Error)
		return
	}

	if len(goals) == 0 {
		log.Println("No old deleted goals to clean up")
		return
	}

	// Permanently delete these goals
	result = s.db.Unscoped().
		Where("deleted_at IS NOT NULL AND deleted_at < ?", thirtyDaysAgo).
		Delete(&models.Goal{})

	if result.Error != nil {
		log.Printf("Error permanently deleting old goals: %v", result.Error)
		return
	}

	log.Printf("Successfully cleaned up %d old deleted goals", result.RowsAffected)
}
