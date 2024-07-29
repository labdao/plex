package models

import (
	"time"

	"gorm.io/datatypes"
)

// event type can only be certain string values
const (
	EventTypeJobQueued     = "job_queued"
	EventTypeJobProcessing = "job_processing"
	EventTypeJobPending    = "job_pending"
	EventTypeJobRunning    = "job_running"
	EventTypeJobStopped    = "job_stopped"
	EventTypeJobSucceeded  = "job_succeeded"
	EventTypeJobFailed     = "job_failed"
)

// retry default 0?
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
	EventType    string         `gorm:"type:varchar(255);not null"`
}
