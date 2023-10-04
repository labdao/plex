package models

type User struct {
<<<<<<< HEAD
	WalletAddress string `gorm:"primaryKey;type:varchar(42);not null" json:"walletAddress"`
	EmailAddress  string `gorm:"type:text" json:"emailAddress"`
=======
	Username      string `json:"username"`
	WalletAddress string `gorm:"type:varchar(42);not null" json:"walletAddress"`
>>>>>>> main
}
