package models

import "time"

type CheckIn struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	GoalID      uint      `gorm:"not null;index" json:"goal_id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	Date        string    `gorm:"not null" json:"date"`
	Status      string    `gorm:"not null" json:"status"`
	ReviewNotes string    `json:"review_notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
