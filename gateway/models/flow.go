package models

import "time"

type Flow struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	CID           string    `gorm:"column:cid;type:varchar(255);not null"`
	Jobs          []Job     `gorm:"foreignKey:FlowID"`
	Name          string    `gorm:"type:varchar(255);"`
	WalletAddress string    `gorm:"type:varchar(42);not null"`
	StartTime     time.Time `gorm:""`
	EndTime       time.Time `gorm:""`
	FlowUUID      string    `gorm:"type:uuid"`
}
