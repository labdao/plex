package models

type JobGroup struct {
	ID               uint   `gorm:"primaryKey"`
	InitialIoJsonCID string `gorm:"column:initial_io_json_cid;type:varchar(255);not null"`
	// does this have
	CompletedIoCID string `gorm:"column:completed_io_cid;type:varchar(255);not null;default:''"`
	Status         string `gorm:"type:varchar(255);default:'initialized'"`

	// doesn't need to be included because Job has JobGroupID as a fk

	Jobs   []Job  `gorm:"foreignKey:JobGroupID"` // slice of Job structs, one-to-many
	Name   string `gorm:"type:varchar(255);not null"`
	UserID uint   `gorm:"column:user_id;type:int;not null;foreignKey"`
}
