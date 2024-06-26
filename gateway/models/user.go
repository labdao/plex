package models

import "time"

type Tier int

const (
	TierFree Tier = iota
	TierPaid
)

type User struct {
	ID             uint         `gorm:"primaryKey;autoIncrement"`
	WalletAddress  string       `gorm:"column:wallet_address;type:varchar(255);unique;not null" json:"walletAddress"`
	DID            string       `gorm:"column:did;type:varchar(255);unique" json:"did"`
	CreatedAt      time.Time    `gorm:""`
	Admin          bool         `gorm:"column:admin;default:false" json:"admin"`
	Tier           Tier         `gorm:"type:int;not null;default:0" json:"tier"`
	ComputeTally   int          `gorm:"column:compute_tally;default:0" json:"computeTally"`
	StripeUserID   string       `gorm:"column:stripe_user_id;type:varchar(255)" json:"stripeUserId"`
	APIKeys        []APIKey     `gorm:"foreignKey:UserID"`
	UserFiles      []File       `gorm:"many2many:user_files;foreignKey:ID;joinForeignKey:user_id;inverseJoinForeignKey:file_id"`
	OrganizationID uint         `gorm:"column:organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
}
