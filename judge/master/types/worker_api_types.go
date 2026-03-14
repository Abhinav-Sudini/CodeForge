package types

type Worker_info struct {
	Port int    `json:"Port"`
	IP   string `json:"IP"`
}

type Worker_req_json struct {
	QuestionId     int    `json:"QuestionId"`
	Runtime        string `json:"runtime"`
	TimeConstrain  int    `json:"timeConstrain"`
	MemConstrain   int    `json:"memConstrain"`
	Code           string `json:"code"`
	JobId          int    `json:"jobId"`
	InternalApiKey string `json:"internalApiKey"`
}

type Worker_Response_json struct {
	JobId          int    `json:"jobId"`
	QuestionId     int    `json:"QuestionId"`
	Result         int    `json:"Result"`
	Mem_usage      int    `json:"Mem_usage"`
	Time_ms        int    `json:"Time_ms"`
	WA_Test_case   int    `json:"WA_Test_case"`
	WA_Stdout      string `json:"WA_Stdout"`
	Stderr         string `json:"Stderr"`
	MSG            string `json:"MSG"`
	InternalApiKey string `json:"internalApiKey"`
}
