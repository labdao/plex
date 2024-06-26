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
	ID            uint           `gorm:"primaryKey;autoIncrement"`
	RayJobID      string         `gorm:"type:varchar(255)"`
	State         JobState       `gorm:"type:varchar(255);default:'queued'"`
	Error         string         `gorm:"type:text;default:''"`
	WalletAddress string         `gorm:"type:varchar(255)"`
	ModelID       string         `gorm:"type:varchar(255);not null;index"`
	Model         Model          `gorm:"foreignKey:ModelID"`
	ExperimentID  uint           `gorm:"type:int;not null;index"`
	Experiment    Experiment     `gorm:"foreignKey:ExperimentID"`
	Inputs        datatypes.JSON `gorm:"type:json"`
	InputFiles    []File         `gorm:"many2many:job_input_files;foreignKey:ID;references:CID"`
	OutputFiles   []File         `gorm:"many2many:job_output_files;foreignKey:ID;references:CID"`
	Queue         QueueType      `gorm:"type:varchar(255)"`
	CreatedAt     time.Time      `gorm:""`
	StartedAt     time.Time      `gorm:""`
	CompletedAt   time.Time      `gorm:""`
	Annotations   string         `gorm:"type:varchar(255)"`
	Public        bool           `gorm:"type:boolean;not null;default:false"`
	RetryCount    int            `gorm:"type:int;not null;default:0"`
	JobType       JobType        `gorm:"type:varchar(255);not null;default:'bacalhau'"`
	ResultJSON    datatypes.JSON `gorm:"type:json"`
}
