package models

import (
	"time"
)

type Transaction struct {
	ID            string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Amount        float64   `gorm:"type:float" json:"amount"`
	IsDebit       bool      `gorm:"type:boolean" json:"isDebit"`
	WalletAddress string    `gorm:"type:varchar(42);not null" json:"walletAddress"`
	Description   string    `gorm:"type:text" json:"description"`
	CreatedAt     time.Time `gorm:""`
	User          User      `gorm:"foreignKey:WalletAddress;references:WalletAddress"`
}
