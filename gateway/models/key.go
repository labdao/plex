package models

import "time"

const (
	ScopeReadOnly  = "read-only"
	ScopeReadWrite = "read-write"
	ScopeAdmin     = "admin"
)

type APIKey struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Key       string    `gorm:"type:varchar(255);not null;unique"`
	Scope     string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:""`
	UserID    string    `gorm:"type:varchar(42);not null"`
	User      User      `gorm:"foreignKey:UserID"`
}
