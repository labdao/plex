package models

import (
	"gorm.io/datatypes"
)

type Tool struct {
	CID           string         `gorm:"primaryKey;column:cid;type:varchar(255);not null"`
	Name          string         `gorm:"type:text;not null;unique"`
	WalletAddress string         `gorm:"type:varchar(42);not null"`
	ToolJson      datatypes.JSON `gorm:"type:json"`
	Container     string         `gorm:"type:text"`
	Memory        int            `gorm:"type:int"`
	Cpu           float64        `gorm:"type:float"`
	Gpu           int            `gorm:"type:int"`
	Network       bool           `gorm:"type:boolean"`
}
