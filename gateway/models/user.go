package models

type User struct {
	Username      string `json:"username"`
	WalletAddress string `gorm:"type:varchar(42);not null" json:"walletAddress"`
}
