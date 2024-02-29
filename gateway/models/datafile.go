package models

import (
	"time"
)

type DataFile struct {
	CID           string    `gorm:"primaryKey;column:cid;type:varchar(255);not null"`
	WalletAddress string    `gorm:"type:varchar(42);not null"`
	Filename      string    `gorm:"type:varchar(255);not null"`
	Timestamp     time.Time `gorm:""`
	InputFiles    []Job     `gorm:"many2many:job_input_files;foreignKey:CID;joinForeignKey:data_file_c_id;inverseJoinForeignKey:job_id"`
	OutputFiles   []Job     `gorm:"many2many:job_output_files;foreignKey:CID;joinForeignKey:data_file_c_id;inverseJoinForeignKey:job_id"`
	Tags          []Tag     `gorm:"many2many:datafile_tags;foreignKey:CID;joinForeignKey:data_file_c_id;inverseJoinForeignKey:tag_name"`
	Public        bool      `gorm:"type:boolean;not null;default:false"`
}
