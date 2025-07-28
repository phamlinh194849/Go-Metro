package models

import (
	"go-metro/config"
	"go-metro/consts"
	"time"
)

type Card struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	RFID      string          `gorm:"uniqueIndex;not null" json:"rf_id"`
	Username  string          `gorm:"uniqueIndex;not null" json:"username"` // 1-1 nên unique
	Balance   float64         `json:"balance" gorm:"default:0"`
	Status    consts.Status   `json:"status"`
	Price     float64         `json:"price" gorm:"default:0"`
	Type      consts.CardType `json:"type"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`

	User User `gorm:"foreignKey:Username;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"` // quan hệ belongs-to
}

func MigrateCard() {
	config.DB.AutoMigrate(&Card{})
}
