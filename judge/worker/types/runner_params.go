package types

type RunnerParamsJson struct{
	CodeDir string `json:"code"`
	TestCasesDir string `json:"TestCasesDir"`
	Runtime string `json:"runtime"`
	TimeConstrain int `json:"timeConstrain"`
	MemConstrain int `json:"memConstrain"`
}
