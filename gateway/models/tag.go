package models

type Tag struct {
	Name string `gorm:"primaryKey;type:varchar(255);not null;unique"`
	Type string `gorm:"type:varchar(100);not null"`
}
