package models

type InputOutput struct {
	ID         uint   `gorm:"primaryKey"`
	JobID      uint   `gorm:"column:job_id;type:int;not null"`
	KeyName    string `gorm:"column:key_name;type:varchar(255);not null"`
	Class      string `gorm:"column:class;type:varchar(255);not null"`
	DatafileID uint   `gorm:"column:datafile_id;type:int;not null;foreignKey"`
}
