package models

import "time"

type User struct {
	WalletAddress string     `gorm:"primaryKey;type:varchar(42);not null" json:"walletAddress"`
	DID           string     `gorm:"column:did;type:varchar(255);unique" json:"did"`
	CreatedAt     time.Time  `gorm:""`
	APIKeys       []APIKey   `gorm:"foreignKey:UserID"`
	Admin         bool       `gorm:"column:admin;default:false" json:"admin"`
	UserDatafiles []DataFile `gorm:"many2many:user_datafiles;foreignKey:WalletAddress;joinForeignKey:wallet_address;inverseJoinForeignKey:c_id"`
}
