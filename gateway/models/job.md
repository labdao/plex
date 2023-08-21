package models

import (
  "gorm.io/gorm"
  "gorm.io/datatypes"
)

type DAG struct {
	gorm.Model
  DAGID       string `json:"dag_id" gorm:"uniqueIndex"`
  CID string `json:"cid" gorm:"uniqueIndex"`
  Jobs string `fk`
}

type Job struct {
	gorm.Model
  Id
  BaclhauJobId
  State
  ErrMsg
  UserId
  Inputs FK to datafile
  Outputs Fk to datafile
  Tool Fk by CID
}

type DataFile struct {
	gorm.Model
  id
  filename
  cid
}
