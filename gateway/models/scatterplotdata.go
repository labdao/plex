package models

type ScatterPlotData struct {
	Plddt         float64 `json:"plddt"`
	IPae          float64 `json:"i_pae"`
	Checkpoint    string  `json:"checkpoint"`
	StructureFile string  `json:"structureFile"`
	PdbFilePath   string  `json:"pdbFilePath"`
	JobUUID       string  `json:"jobUUID"`
}
