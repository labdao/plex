package models

type Graph struct {
	CID           string `gorm:"primaryKey;column:cid;type:varchar(255);not null"` // column name specified to avoid GORM default snake case
	Jobs          []Job  `gorm:"foreignKey:GraphID"`
	WalletAddress string `gorm:"type:varchar(42);not null"`
}
