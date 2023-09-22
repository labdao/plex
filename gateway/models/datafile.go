package models

import (
	"time"
)

type DataFile struct {
	CID           string    `gorm:"primaryKey;column:cid;type:varchar(255);not null"` // column name specified to avoid GORM default snake case
	WalletAddress string    `gorm:"type:varchar(42);not null"`
	Filename      string    `gorm:"type:varchar(255);not null"`
	Timestamp     time.Time `gorm:""`
	Public        bool      `gorm:"default:true"`
	Visible       bool      `gorm:"default:true"`
}
