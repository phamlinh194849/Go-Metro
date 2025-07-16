package models

import (
	"go-metro/config"
	"time"
)

type Card struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CardID    string    `gorm:"uniqueIndex;not null" json:"card_id"`
	UserID    string    `json:"user_id"`
	Balance   float64   `json:"balance" gorm:"default:0"`
	Status    string    `json:"status" gorm:"default:'active'"` // active, inactive, blocked
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func MigrateCard() {
	config.DB.AutoMigrate(&Card{})
} 