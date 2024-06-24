package models

import (
	"time"

	"gorm.io/datatypes"
)

type Model struct {
	CID                string         `gorm:"primaryKey;column:cid;type:varchar(255);not null"`
	Name               string         `gorm:"type:text;not null;unique"`
	WalletAddress      string         `gorm:"type:varchar(42);not null"`
	ModelJson          datatypes.JSON `gorm:"type:json"`
	Container          string         `gorm:"type:text"`
	Memory             int            `gorm:"type:int"`
	Cpu                float64        `gorm:"type:float"`
	Gpu                int            `gorm:"type:int"`
	Network            bool           `gorm:"type:boolean"`
	Timestamp          time.Time      `gorm:""`
	Display            bool           `gorm:"type:boolean;default:true"`
	TaskCategory       string         `gorm:"type:text;default:'community-models'"`
	DefaultModel       bool           `gorm:"type:boolean;default:false"`
	MaxRunningTime     int            `gorm:"type:int;default:2700"`
	ComputeCost        int            `gorm:"type:int;not null;default:0"`
	ModelType          string         `gorm:"type:varchar(255);default:'bacalhau'"`
	RayServiceEndpoint string         `gorm:"type:varchar(255)"`
	S3URI              string         `gorm:"type:varchar(255)"`
}
