package types

type JudgeCodeRequest struct {
	QuestionId     int    `json:"QuestionId"`
	Runtime        string `json:"runtime"`
	TimeConstrain  int    `json:"timeConstrain"`
	MemConstrain   int    `json:"memConstrain"`
	Code           string `json:"code"`
	JobId          int    `json:"jobId"`
	InternalApiKey string `json:"internalApiKey"`
}

type JudgeCodeReqVerdict struct {
	JobId   int    `json:"JobId"`
	Verdict string `json:"Verdict"`
	Msg     string `json:"Msg"`
}
