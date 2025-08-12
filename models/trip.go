package models

import (
	"go-metro/config"
	"time"
)

type Trip struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Direction  string    `json:"direction"`
	TrainID    uint      `json:"train_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Foreign key relationships
	Train Train `gorm:"foreignKey:TrainID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"train"`
}

func MigrateTrip() {
	config.DB.AutoMigrate(&Trip{})
}