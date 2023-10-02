package models

// gorm supports the use of composite primary keys
// combination of wallet address and email must be unique
// https://gorm.io/docs/composite_primary_key.html

type User struct {
	WalletAddress string `gorm:"primaryKey;type:varchar(42);not null" json:"walletAddress"`
	EmailAddress  string `gorm:"type:text" json:"emailAddress"`
}
