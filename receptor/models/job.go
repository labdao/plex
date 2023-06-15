package models

import (
  "gorm.io/gorm"
  "gorm.io/datatypes"
)

type Job struct {
	gorm.Model
  JobID       string `json:"job_id"`
  Spec        datatypes.JSON `json:"spec"`
}
