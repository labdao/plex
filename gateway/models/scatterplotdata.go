package models

type ScatterPlotData struct {
	Plddt             float64 `json:"plddt"`
	IPae              float64 `json:"i_pae"`
	Checkpoint        string  `json:"checkpoint"`
	ProposedStructure string  `json:"proposedStructure"`
	PdbFilePath       string  `json:"pdbFilePath"`
}
