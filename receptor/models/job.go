package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	JobID string         `json:"job_id" gorm:"uniqueIndex"`
	Spec  datatypes.JSON `json:"spec"`
}
