package models

import (
	"time"
)

type DataFile struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CID           string    `gorm:"index;column:cid;type:varchar(255);not null"`
	WalletAddress string    `gorm:"type:varchar(42);not null"`
	Filename      string    `gorm:"type:varchar(255);not null"`
	Timestamp     time.Time `gorm:""`
	InputFiles    []Job     `gorm:"many2many:job_input_files;foreignKey:ID;joinForeignKey:data_file_id;inverseJoinForeignKey:job_id"`
	OutputFiles   []Job     `gorm:"many2many:job_output_files;foreignKey:ID;joinForeignKey:data_file_id;inverseJoinForeignKey:job_id"`
	Tags          []Tag     `gorm:"many2many:datafile_tags;foreignKey:ID;joinForeignKey:data_file_id;inverseJoinForeignKey:tag_name"`
}
