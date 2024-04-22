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
	BacalhauJobID string         `gorm:"type:varchar(255)"`
	State         JobState       `gorm:"type:varchar(255);default:'queued'"`
	Error         string         `gorm:"type:text;default:''"`
	WalletAddress string         `gorm:"type:varchar(255)"`
	ToolID        string         `gorm:"type:varchar(255);not null;index"`
	Tool          Tool           `gorm:"foreignKey:ToolID"`
	FlowID        uint           `gorm:"type:int;not null;index"`
	Flow          Flow           `gorm:"foreignKey:FlowID"`
	Inputs        datatypes.JSON `gorm:"type:json"`
	InputFiles    []DataFile     `gorm:"many2many:job_input_files;foreignKey:ID;references:CID"`
	OutputFiles   []DataFile     `gorm:"many2many:job_output_files;foreignKey:ID;references:CID"`
	Queue         QueueType      `gorm:"type:varchar(255)"`
	CreatedAt     time.Time      `gorm:""`
	StartedAt     time.Time      `gorm:""`
	CompletedAt   time.Time      `gorm:""`
	Annotations   string         `gorm:"type:varchar(255)"`
	JobUUID       string         `gorm:"type:uuid"`
	Public        bool           `gorm:"type:boolean;not null;default:false"`
	RetryCount    int            `gorm:"type:int;not null;default:0"`
}
