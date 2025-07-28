package models

import (
	"time"

	"go-metro/config"
	"go-metro/consts"
)

type User struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	Username  string      `gorm:"uniqueIndex;not null" json:"username"`
	Password  string      `gorm:"not null" json:"-"` // "-" để không trả về password trong JSON
	Email     string      `gorm:"uniqueIndex" json:"email"`
	FullName  string      `json:"full_name"`
	Role      consts.Role `gorm:"default:2" json:"role"`          // "ADMIN", "USER", "STAFF"
	Status    string      `gorm:"default:'active'" json:"status"` // "active" hoặc "inactive"
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

func MigrateUser() {
	config.DB.AutoMigrate(&User{})
}
