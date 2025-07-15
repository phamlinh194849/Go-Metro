package models

import (
	"go-metro/db"
	"time"
)

type History struct {
	ID      uint      `gorm:"primaryKey" json:"id"`
	CardID  string    `json:"card_id"`
	Time    time.Time `json:"time"`
	UserID  string    `json:"user_id"`
	Balance float64   `json:"balance"`
	Action  string    `json:"action"`
}

func MigrateHistory() {
	db.DB.AutoMigrate(&History{})
}
