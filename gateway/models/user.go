package models

type User struct {
	ID            uint   `gorm:"primary_key" json:"id"`
	Username      string `json:"username"`
	WalletAddress string `gorm:"type:varchar(42);not null" json:"walletAddress"`
}
