package models

import (
	"time"

	"gorm.io/datatypes"
)

type JobState string

const (
	JobStateQueued    JobState = "queued"
	JobStateRunning   JobState = "running"
	JobStateFailed    JobState = "failed"
	JobStateCompleted JobState = "completed"
)

type QueueType string

const (
	QueueTypeRay QueueType = "ray"
)

type JobType string

const (
	JobTypeBacalhau JobType = "bacalhau"
	JobTypeRay      JobType = "ray"
)

type Job struct {
	ID             uint           `gorm:"primaryKey;autoIncrement"`
	RayJobID       string         `gorm:"type:varchar(255)"`
	JobStatus      JobState       `gorm:"type:varchar(255);default:'queued'"`
	CreatedAt      time.Time      `gorm:""`
	StartedAt      time.Time      `gorm:""`
	CompletedAt    time.Time      `gorm:""`
	LastModifiedAt time.Time      `gorm:"autoUpdateTime"`
	ExperimentID   uint           `gorm:"type:int;not null;index"`
	Experiment     Experiment     `gorm:"foreignKey:ExperimentID"`
	ModelID        int            `gorm:"type:int;not null;index"`
	Model          Model          `gorm:"foreignKey:ModelID"`
	UserID         uint           `gorm:"not null"`
	User           User           `gorm:"foreignKey:UserID"`
	Public         bool           `gorm:"type:boolean;not null;default:false"`
	RetryCount     int            `gorm:"type:int;not null;default:0"`
	Error          string         `gorm:"type:text;default:''"`
	Inputs         datatypes.JSON `gorm:"type:json"`
	InputFiles     []File         `gorm:"many2many:job_input_files;foreignKey:ID;joinForeignKey:job_id;References:ID;JoinReferences:file_id"`
	OutputFiles    []File         `gorm:"many2many:job_output_files;foreignKey:ID;references:ID"`
}
