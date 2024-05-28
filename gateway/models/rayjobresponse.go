package models

type RayJobResponse struct {
	UUID             string                 `json:"uuid"`
	PDB              FileDetail             `json:"pdb"`
	StructureMetrics FileDetail             `json:"structure_metrics"`
	Plots            []FileDetail           `json:"plots"`
	MSA              FileDetail             `json:"msa"`
	DynamicFields    map[string]interface{} `json:"-"`
}

type FileDetail struct {
	Key      string `json:"key"`
	Location string `json:"location"`
}
