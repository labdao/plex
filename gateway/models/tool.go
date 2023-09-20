package models

type Tool struct {
	ID            uint   `gorm:"primaryKey"`
	CID           string `gorm:"column:cid;type:varchar(255);not null"` // column name specified to avoid GORM default snake case
	ToolJSON      string `gorm:"type:varchar(255);not null"`
	WalletAddress string `gorm:"type:varchar(42);not null"`
}
