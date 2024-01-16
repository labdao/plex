package models

import "time"

type Flow struct {
	CID           string    `gorm:"primaryKey;column:cid;type:varchar(255);not null"` // column name specified to avoid GORM default snake case
	Jobs          []Job     `gorm:"foreignKey:FlowID"`
	Name          string    `gorm:"type:varchar(255);"`
	WalletAddress string    `gorm:"type:varchar(42);not null"`
	StartTime     time.Time `gorm:""`
	EndTime       time.Time `gorm:""`
}
