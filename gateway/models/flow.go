package models

type Flow struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	CID           string `gorm:"column:cid;type:varchar(255);unique;not null"`
	Jobs          []Job  `gorm:"foreignKey:FlowID"`
	Name          string `gorm:"type:varchar(255);"`
	WalletAddress string `gorm:"type:varchar(42);not null"`
}
