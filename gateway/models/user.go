package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID            uint   `gorm:"primaryKey"`
	WalletAddress string `gorm:"type:varchar(42);not null"`
}
