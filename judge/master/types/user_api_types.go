package types

type User_judge_req_json struct {
	QuestionId int    `json:"QuestionId"`
	Runtime    string `json:"runtime"`
	Code       string `json:"code"`
}

type Submission_verdict_req_json struct {
	Submission_id int `json:"submission_id"`
}
type Submission_verdict_resp_json struct {
	Submission_id int    `json:"Submission_id"`
	QuestionId    int    `json:"QuestionId"`
	Verdict       string `josn:"Verdict"`
	Mem_usage     int    `json:"Mem_usage"`
	Time_ms       int    `json:"Time_ms"`
	WA_Test_case  int    `json:"WA_Test_case"`
	WA_Stdout     string `json:"WA_Stdout"`
	Stderr        string `json:"Stderr"`
}
