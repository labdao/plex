package ipwl

import (
	"encoding/json"
	"fmt"

	"github.com/labdao/plex/gateway/models"
	"gorm.io/gorm"
)

type ModelInput struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Array       bool        `json:"array"`
	Glob        []string    `json:"glob"`
	Default     interface{} `json:"default"`
	Min         string      `json:"min"`
	Max         string      `json:"max"`
	Example     string      `json:"example"`
	Grouping    string      `json:"grouping"`
	Position    string      `json:"position"`
	Required    bool        `json:"required"`
}

type ModelOutput struct {
	Type string   `json:"type"`
	Item string   `json:"item"`
	Glob []string `json:"glob"`
}

type ModelType string

const (
	ModelTypeBacalhau ModelType = "bacalhau"
	ModelTypeRay      ModelType = "ray"
)

type Model struct {
	Name                 string                 `json:"name"`
	Description          string                 `json:"description"`
	Guide                string                 `json:"guide"`
	Author               string                 `json:"author"`
	GitHub               string                 `json:"github"`
	Paper                string                 `json:"paper"`
	Task                 string                 `json:"task"`
	CheckpointCompatible bool                   `json:"checkpointCompatible"`
	Inputs               map[string]ModelInput  `json:"inputs"`
	Outputs              map[string]ModelOutput `json:"outputs"`
	TaskCategory         string                 `json:"taskCategory"`
	MaxRunningTime       int                    `json:"maxRunningTime"`
	ComputeCost          int                    `json:"computeCost"`
	ModelType            ModelType              `json:"modelType"`
	RayServiceEndpoint   string                 `json:"rayServiceEndpoint"`
	XAxis                string                 `json:"xAxis"`
	YAxis                string                 `json:"yAxis"`
}

func ReadModelConfig(modelID int, db *gorm.DB) (Model, string, error) {
	var ipwlmodel Model
	var dbModel models.Model
	var modelName string
	var err error

	err = db.Where("id = ?", modelID).First(&dbModel).Error
	if err != nil {
		return ipwlmodel, modelName, fmt.Errorf("failed to get model from database: %w", err)
	}
	modelName = dbModel.Name

	// bytes should come from model_json column of model
	err = json.Unmarshal(dbModel.ModelJson, &ipwlmodel)
	if err != nil {
		return ipwlmodel, modelName, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return ipwlmodel, modelName, nil
}
