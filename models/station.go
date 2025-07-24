package models

import (
	"go-metro/config"
	"time"
)

type Station struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func MigrateStation() {
	config.DB.AutoMigrate(&Station{})
} 