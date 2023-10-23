package models

type User struct {
	WalletAddress string `gorm:"primaryKey;type:varchar(42);not null" json:"walletAddress"`
	IsMember      bool   `gorm:"type:boolean;not null;default:false" json:"isMember"`
}
