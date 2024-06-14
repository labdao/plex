package models

import (
	"time"

	"gorm.io/datatypes"
)

type RequestTracker struct {
	ID           uint           `gorm:"primaryKey;autoIncrement"`
	JobID        uint           `gorm:"not null"`
	Job          Job            `gorm:"foreignKey:JobID"`
	JobResponse  datatypes.JSON `gorm:"type:json"`
	RayJobID     string         `gorm:"type:varchar(255);not null"`
	RetryCount   int            `gorm:"not null"`
	State        JobState       `gorm:"type:varchar(255);default:'queued'"`
	ResponseCode int            `gorm:"type:int"`
	CreatedAt    time.Time      `gorm:""`
	StartedAt    time.Time      `gorm:""`
	CompletedAt  time.Time      `gorm:""`
	ErrorMessage string         `gorm:"type:text;default:''"`
}
