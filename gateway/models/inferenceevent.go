package models

import (
	"time"

	"gorm.io/datatypes"
)

type InferenceEvent struct {
	ID           uint           `gorm:"primaryKey;autoIncrement"`
	JobID        uint           `gorm:"not null"`
	Job          Job            `gorm:"foreignKey:JobID"`
	RayJobID     string         `gorm:"type:varchar(255);not null"`
	InputJson    datatypes.JSON `gorm:"type:json"`
	OutputJson   datatypes.JSON `gorm:"type:json"`
	RetryCount   int            `gorm:"not null"`
	JobStatus    JobState       `gorm:"type:varchar(255);default:'queued'"`
	ResponseCode int            `gorm:"type:int"`
	EventTime    time.Time      `gorm:""`
	EventMessage string         `gorm:"type:text"`
	EventType    string         `gorm:"type:text"`
}
