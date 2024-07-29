package models

import (
	"time"

	"gorm.io/datatypes"
)

type JobType string

const (
	JobTypeJob     JobType = "job"
	JobTypeService JobType = "service"
)

type Model struct {
	ID               int            `gorm:"primaryKey;autoIncrement"`
	Name             string         `gorm:"type:text;not null;unique"`
	WalletAddress    string         `gorm:"type:varchar(42);not null"`
	ModelJson        datatypes.JSON `gorm:"type:json"`
	CreatedAt        time.Time      `gorm:"autoCreateTime"`
	Display          bool           `gorm:"type:boolean;default:true"`
	TaskCategory     string         `gorm:"type:text;default:'community-models'"`
	DefaultModel     bool           `gorm:"type:boolean;default:false"`
	MaxRunningTime   int            `gorm:"type:int;default:2700"`
	ComputeCost      int            `gorm:"type:int;not null;default:0"`
	RayEndpoint      string         `gorm:"type:varchar(255)"`
	RayJobEntrypoint string         `gorm:"type:varchar(255)"`
	S3URI            string         `gorm:"type:varchar(255)"`
	JobType          JobType        `gorm:"type:text;default:'job'"`
}
