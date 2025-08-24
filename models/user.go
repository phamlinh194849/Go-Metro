package models

import (
  "time"

  "go-metro/config"
  "go-metro/consts"
)

type User struct {
  ID        uint        `gorm:"primaryKey" json:"id"`
  Password  string      `gorm:"not null" json:"-"` // "-" để không trả về password trong JSON
  Email     string      `gorm:"uniqueIndex" json:"email"`
  FullName  string      `json:"full_name"`
  Role      consts.Role `gorm:"default:3" json:"role"`          // "ADMIN", "STAFF", "USER"
  Status    string      `gorm:"default:'active'" json:"status"` // "active" hoặc "inactive"
  Avatar    string      `json:"avatar"`
  Phone     string      `json:"phone"`
  CreatedAt time.Time   `json:"created_at"`
  UpdatedAt time.Time   `json:"updated_at"`
}

func MigrateUser() {
  config.DB.AutoMigrate(&User{})
}
