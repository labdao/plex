package models

import "time"

type User struct {
	WalletAddress string    `gorm:"primaryKey;type:varchar(42);not null" json:"walletAddress"`
	CreatedAt     time.Time `gorm:""`
	APIKeys       []APIKey  `gorm:"foreignKey:UserID"`
}
