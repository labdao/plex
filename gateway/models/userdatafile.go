package models

import "time"

type UserDatafile struct {
	WalletAddress string    `gorm:"primaryKey;type:varchar(42);not null" json:"walletAddress"`
	CID           string    `gorm:"primaryKey;type:varchar(255);not null" json:"cid"`
	CreatedAt     time.Time `gorm:"" json:"createdAt"`
	User          User      `gorm:"foreignKey:WalletAddress;references:WalletAddress"`
	DataFile      DataFile  `gorm:"foreignKey:CID;references:CID"`
}
