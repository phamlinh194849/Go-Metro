package models

import (
	"go-metro/config"
	"time"
)

type SellHistory struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CardID        string    `json:"card_id"`
	SellerID      uint      `json:"seller_id"`
	CardPriceSold float64   `json:"card_price_sold"`
	Time          time.Time `json:"time"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Foreign key relationships
	Card   Card `gorm:"foreignKey:CardID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"card"`
	Seller User `gorm:"foreignKey:SellerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"seller"`
}

func MigrateSellHistory() {
	config.DB.AutoMigrate(&SellHistory{})
} 