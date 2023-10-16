package models

import (
	"time"
)

type DataFile struct {
	CID           string    `gorm:"primaryKey;column:cid;type:varchar(255);not null"` // column name specified to avoid GORM default snake case
	WalletAddress string    `gorm:"type:varchar(42);not null"`
	Filename      string    `gorm:"type:varchar(255);not null"`
	JobInputs     []Job     `gorm:"many2many:job_inputs;foreignKey:CID;references:BacalhauJobID"`
	JobOutputs    []Job     `gorm:"many2many:job_outputs;foreignKey:CID;references:BacalhauJobID"`
	Timestamp     time.Time `gorm:""`
}
