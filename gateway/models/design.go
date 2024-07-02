package models

import "gorm.io/datatypes"

type Design struct {
	ID                int            `json:"id"`
	JobID             int            `json:"job_id"`
	XAxisValue        string         `json:"x_axis_value"`
	YAxisValue        string         `json:"y_axis_value"`
	CheckpointPDBID   int            `json:"checkpoint_pdb_id"`
	AdditionalDetails datatypes.JSON `json:""`
}
