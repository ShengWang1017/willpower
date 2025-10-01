package models

import (
	"time"
	"gorm.io/gorm"
)

type Goal struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	Type      string         `gorm:"not null" json:"type"`
	Title     string         `gorm:"not null" json:"title"`
	Status    string         `gorm:"not null;default:'active'" json:"status"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
