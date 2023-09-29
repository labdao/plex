package models

type Tool struct {
	CID           string `gorm:"primaryKey;column:cid;type:varchar(255);not null"` // column name specified to avoid GORM default snake case
	Name          string `gorm:"type:text;not null;unique"`
	WalletAddress string `gorm:"type:varchar(42);not null"`
}
