package models

import "time"

type Tier int

const (
	TierFree Tier = iota
	TierPaid
)

type User struct {
	WalletAddress      string       `gorm:"primaryKey;type:varchar(42);not null" json:"walletAddress"`
	DID                string       `gorm:"column:did;type:varchar(255);unique" json:"did"`
	CreatedAt          time.Time    `gorm:""`
	APIKeys            []APIKey     `gorm:"foreignKey:UserID"`
	Admin              bool         `gorm:"column:admin;default:false" json:"admin"`
	UserFiles          []File       `gorm:"many2many:user_files;foreignKey:WalletAddress;joinForeignKey:wallet_address;inverseJoinForeignKey:file_id"`
	Tier               Tier         `gorm:"type:int;not null;default:0" json:"tier"`
	ComputeTally       int          `gorm:"column:compute_tally;default:0" json:"computeTally"`
	StripeUserID       *string      `gorm:"column:stripe_user_id;type:varchar(255)" json:"stripeUserId"`
	OrganizationID     uint         `gorm:"column:organization_id"`
	Organization       Organization `gorm:"foreignKey:OrganizationID"`
	SubscriptionStatus *string      `gorm:"column:subscription_status;type:varchar(255)" json:"subscriptionStatus"`
	SubscriptionID     *string      `gorm:"column:subscription_id;type:varchar(255)" json:"subscriptionId"`
}
