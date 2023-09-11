package models

type Job struct {
	InitialIoCID    string `gorm:"column:initial_io_cid;type:varchar(255);not null"`
	InitialIoJson   string `gorm:"column:initial_io_json;type:varchar(255);not null"`
	CompletedIoCID  string `gorm:"column:completed_io_cid;type:varchar(255);not null;default:''"`
	CompletedIoJson string `gorm:"column:completed_io_json;type:varchar(255);not null;default:''"`
	Status          string `gorm:"type:varchar(255);default:'initialized'"`
	// TimeStamp of job submission
	// TimeStamp of job completion
}
