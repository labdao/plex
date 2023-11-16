package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type JobModel struct {
	gorm.Model
	NodeID string         `json:"NodeID"`
	Spec   datatypes.JSON `json:"Job" gorm:"column:spec"`    // Store Job object in JSON format
	JobID  string         `gorm:"column:job_id;uniqueIndex"` // Extracted Job.ID field
}

type Tabler interface {
	TableName() string
}

// TableName overrides the table name
func (JobModel) TableName() string {
	return "jobs"
}
