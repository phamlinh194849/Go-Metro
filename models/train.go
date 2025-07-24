package models

import (
	"go-metro/config"
	"time"
)

type Train struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Type      string    `json:"type"`
	Company   string    `json:"company"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func MigrateTrain() {
	config.DB.AutoMigrate(&Train{})
} 