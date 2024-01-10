package models

import (
	"time"
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
	BacalhauJobID string     `gorm:"primaryKey;type:varchar(255);not null"`
	State         JobState   `gorm:"type:varchar(255);default:'queued'"`
	Error         string     `gorm:"type:text;default:''"`
	WalletAddress string     `gorm:"type:varchar(255)"`
	ToolID        string     `gorm:"type:varchar(255);not null;index"`
	Tool          Tool       `gorm:"foreignKey:ToolID"`
	FlowID        string     `gorm:"type:varchar(255);not null;index"`
	Flow          Flow       `gorm:"foreignKey:FlowID"`
	Inputs        []DataFile `gorm:"many2many:job_inputs;foreignKey:BacalhauJobID;references:CID"`
	Outputs       []DataFile `gorm:"many2many:job_outputs;foreignKey:BacalhauJobID;references:CID"`
	Queue         QueueType  `gorm:"type:varchar(255);not null"`
	CreatedAt     time.Time  `gorm:"type:timestamp"`
}
