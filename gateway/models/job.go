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
	QueueTypeCPU QueueType = "cpu"
	QueueTypeGPU QueueType = "gpu"
)

type Job struct {
	ID            uint           `gorm:"primaryKey;autoIncrement"`
	BacalhauJobID string         `gorm:"type:varchar(255);index;not null"`
	State         JobState       `gorm:"type:varchar(255);default:'queued'"`
	Error         string         `gorm:"type:text;default:''"`
	WalletAddress string         `gorm:"type:varchar(255)"`
	ToolID        string         `gorm:"type:varchar(255);not null;index"`
	Tool          Tool           `gorm:"foreignKey:ToolID"`
	FlowID        string         `gorm:"type:varchar(255);not null;index"`
	Flow          Flow           `gorm:"foreignKey:FlowID"`
	Inputs        datatypes.JSON `gorm:"type:json"`
	InputFiles    []DataFile     `gorm:"many2many:job_inputs;foreignKey:BacalhauJobID;references:CID"`
	OutputFiles   []DataFile     `gorm:"many2many:job_outputs;foreignKey:BacalhauJobID;references:CID"`
	Queue         QueueType      `gorm:"type:varchar(255);not null"`
	CreatedAt     time.Time      `gorm:"type:timestamp"`
	Annotations   string         `gorm:"type:varchar(255)"`
}
