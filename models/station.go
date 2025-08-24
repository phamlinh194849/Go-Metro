package models

import (
  "time"

  "go-metro/config"
)

type Station struct {
  ID          uint      `gorm:"primaryKey" json:"id"`
  Name        string    `gorm:"not null" json:"name"`
  IPAddress   string    `json:"ip_address"`
  Capacity    int       `json:"capacity"`
  Address     string    `json:"address"`
  Status      string    `json:"status"`
  Description string    `json:"description"`
  ImageURL    string    `json:"image_url"`
  CreatedAt   time.Time `json:"created_at"`
  UpdatedAt   time.Time `json:"updated_at"`
}

func MigrateStation() {
  config.DB.AutoMigrate(&Station{})
}
