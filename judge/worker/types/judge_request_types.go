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

type JudgeCodeResponse struct {
	JobId          int    `json:"jobId"`
	QuestionId     int    `json:"QuestionId"`
	Result         int    `json:"Result"`
	Mem_usage      int    `json:"Mem_usage"`
	Time_ms        int    `json:"Time_ms"`
	WA_Test_case   int    `json:"WA_Test_case"`
	MSG            string `json:"MSG"`
	InternalApiKey string `json:"internalApiKey"`
}

type JudgeCodeReqVerdict struct {
	JobId   int    `json:"JobId"`
	Verdict string `json:"Verdict"`
	Msg     string `json:"Msg"`
}
