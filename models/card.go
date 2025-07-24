package models

import (
	"go-metro/config"
	"go-metro/consts"
	"time"
)

type Card struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	RFID      string        `gorm:"uniqueIndex;not null" json:"rf_id"`
	OwnerID   string        `json:"owner_id"`
	Balance   float64       `json:"balance" gorm:"default:0"`
	Status    consts.Status `json:"status"`
	Price     float64       `json:"price" gorm:"default:0"`
	Type      consts.CardType
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func MigrateCard() {
	config.DB.AutoMigrate(&Card{})
}
