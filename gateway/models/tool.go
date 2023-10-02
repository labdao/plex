package models

type Tool struct {
	CID           string `gorm:"primaryKey;column:cid;type:varchar(255);not null"` // column name specified to avoid GORM default snake case
	ToolJSON      string `gorm:"type:text"`
	WalletAddress string `gorm:"type:varchar(42);not null"`
}
