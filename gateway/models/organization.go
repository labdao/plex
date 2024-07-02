package models

type Organization struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(255);not null;unique"`
	Description string `gorm:"type:text"`
}
