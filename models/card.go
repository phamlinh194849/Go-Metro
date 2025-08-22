package models

import (
  "time"

  "go-metro/config"
  "go-metro/consts"
)

type Card struct {
  ID        uint            `gorm:"primaryKey" json:"id"`
  UserID    uint            `gorm:"not null" json:"user_id"`
  RFID      string          `gorm:"uniqueIndex;not null" json:"rf_id"`
  Balance   float64         `json:"balance" gorm:"default:0"`
  Status    consts.Status   `json:"status"`
  Price     float64         `json:"price" gorm:"default:0"`
  Type      consts.CardType `json:"type"`
  CreatedAt time.Time       `json:"created_at"`
  UpdatedAt time.Time       `json:"updated_at"`

  User *User `gorm:"foreignKey:UserID" json:"user"`
}

func MigrateCard() {
  config.DB.AutoMigrate(&Card{})
}
