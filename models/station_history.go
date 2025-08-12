package models

import (
	"go-metro/config"
	"time"
)

type StationHistory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Action     string    `json:"action"` // "checkin" hoặc "checkout"
	Time       time.Time `json:"time"`
	CardID     string    `json:"card_id"`
	StationID  uint      `json:"station_id"`
	UsedBalance float64  `json:"used_balance"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Foreign key relationships
	Card    Card    `gorm:"foreignKey:CardID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"card"`
	Station Station `gorm:"foreignKey:StationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"station"`
}

func MigrateStationHistory() {
	config.DB.AutoMigrate(&StationHistory{})
} 