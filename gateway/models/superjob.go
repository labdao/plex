package models

import (
	"time"
)

type SuperJob struct {
	ID        int       `gorm:"primaryKey"`
	FlowID    int       `gorm:"not null"`
	JobUUID   string    `gorm:"type:uuid;not null"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
	Status    string    `gorm:"default:'pending'"`
	Jobs      []Job     `gorm:"foreignKey:ID"`
	FlowID    uint      `gorm:"type:int;not null;index"`
	Flow      Flow      `gorm:"foreignKey:FlowID"`
}
