package models

import (
	"time"
)

type File struct {
	CID           string    `gorm:"primaryKey;column:cid;type:varchar(255);not null"`
	WalletAddress string    `gorm:"type:varchar(42);not null"`
	Filename      string    `gorm:"type:varchar(255);not null"`
	Timestamp     time.Time `gorm:""`
	InputFiles    []Job     `gorm:"many2many:job_input_files;foreignKey:CID;joinForeignKey:file_c_id;inverseJoinForeignKey:job_id"`
	OutputFiles   []Job     `gorm:"many2many:job_output_files;foreignKey:CID;joinForeignKey:file_c_id;inverseJoinForeignKey:job_id"`
	Tags          []Tag     `gorm:"many2many:file_tags;foreignKey:CID;joinForeignKey:file_c_id;inverseJoinForeignKey:tag_name"`
	Public        bool      `gorm:"type:boolean;not null;default:false"`
	UserFiles     []User    `gorm:"many2many:user_files;foreignKey:CID;joinForeignKey:c_id;inverseJoinForeignKey:wallet_address"`
	S3URI         string    `gorm:"type:varchar(255)"`
}
