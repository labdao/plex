package models

import "time"

type RequestTracker struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	JobID        uint      `gorm:"not null"`
	Job          Job       `gorm:"foreignKey:JobID"`
	JobResponse  string    `gorm:"type:text"`
	RayJobUUID   string    `gorm:"type:varchar(255);not null"`
	RetryCount   int       `gorm:"not null"`
	RayJobStatus string    `gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time `gorm:""`
}
