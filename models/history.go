package models

import (
	"go-metro/config"
	"go-metro/consts"
	"time"
)

type History struct {
	ID         uint              `gorm:"primaryKey" json:"id"`
	CardID     string            `json:"card_id"`
	Time       time.Time         `json:"time"`
	UserID     string            `json:"user_id"`
	Balance    float64           `json:"balance"`
	UserAction consts.UserAction `json:"user_action"`
	CardAction consts.CardAction `json:"card_action"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

func MigrateHistory() {
	config.DB.AutoMigrate(&History{})
}
