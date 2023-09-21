package models

type Job struct {
	ID            uint          `gorm:"primaryKey"`
	JobGroupID    uint          `gorm:"column:job_group_id;type:int; not null;foreignKey"`
	Inputs        []InputOutput `gorm:"foreignKey:JobID"`
	Outputs       []InputOutput `gorm:"foreignKey:JobID"`
	State         string        `gorm:"type:varchar(255);default:'initialized'"`
	ErrMsg        string        `gorm:"column:err_msg;type:text;default:''"`
	WalletAddress string        `gorm:"column:wallet_address;type:varchar(255)"`
	BacalhauJobID string        `gorm:"colun:bacalhau_job_id;type:varchar(255);not null"`
	ToolID        uint          `gorm:"column:tool_id;type:int;not null;foreignKey"`
}
