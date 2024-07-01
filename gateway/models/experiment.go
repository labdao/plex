package models

import "time"

type Experiment struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	Jobs           []Job     `gorm:"foreignKey:ExperimentID"`
	Name           string    `gorm:"type:varchar(255);"`
	Public         bool      `gorm:"type:boolean;not null;default:false"`
	RecordCID      string    `gorm:"column:record_cid;type:varchar(255);"`
	WalletAddress  string    `gorm:"not null"`
	User           User      `gorm:"foreignKey:WalletAddress"`
	LastModifiedAt time.Time `gorm:"autoUpdateTime"`
}
