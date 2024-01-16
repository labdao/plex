package models

import "time"

type Job struct {
	BacalhauJobID string     `gorm:"primaryKey;type:varchar(255);not null"`
	State         string     `gorm:"type:varchar(255);default:'processing'"`
	Error         string     `gorm:"type:text;default:''"`
	WalletAddress string     `gorm:"type:varchar(255)"`
	ToolID        string     `gorm:"type:varchar(255);not null;index"`
	Tool          Tool       `gorm:"foreignKey:ToolID"`
	FlowID        string     `gorm:"type:varchar(255);not null;index"`
	Flow          Flow       `gorm:"foreignKey:FlowID"`
	Inputs        []DataFile `gorm:"many2many:job_inputs;foreignKey:BacalhauJobID;references:CID"`
	Outputs       []DataFile `gorm:"many2many:job_outputs;foreignKey:BacalhauJobID;references:CID"`
	CreatedAt     time.Time  `gorm:""`
	StartedAt     time.Time  `gorm:""`
	CompletedAt   time.Time  `gorm:""`
}
