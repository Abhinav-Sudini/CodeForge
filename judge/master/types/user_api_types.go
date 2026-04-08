package types

import "time"

type User_judge_req_json struct {
	QuestionId int    `json:"QuestionId"`
	Runtime    string `json:"runtime"`
	Code       string `json:"code"`
}

type User_judge_resp_json struct {
	Submission_id int    `json:"Submission_id"`
	Verdict       string `josn:"Verdict"`
}

// user submission verdict req and response json

type Submission_verdict_req_json struct {
	Submission_id int `json:"submission_id"`
}
type Submission_verdict_resp_json struct {
	Submission_id int    `json:"Submission_id"`
	SubmittedCode string `json:"Submitted_code"`
	QuestionId    int    `json:"QuestionId"`
	Verdict       string `josn:"Verdict"`
	Mem_usage     int    `json:"Mem_usage"`
	Time_ms       int    `json:"Time_ms"`
	WA_Test_case  int    `json:"WA_Test_case"`
	WA_Stdin     string `json:"WA_Stdin"`
	WA_Stdout     string `json:"WA_Stdout"`
	Required_Stdout string `json:"Required_Stdout"`
	Stderr        string `json:"Stderr"`
	SubmissionTime time.Time `json:"SubmissionTime"`
}

// get all questions with verdicts
type Get_all_questions_resp_json struct {
	Questions []Question_record_minimal `json:"Questions"`
}

type Question_record_minimal struct {
	QuestionId       int    `json:"QuestionId"`
	QuestionName     string    `json:"QuestionName"`
	QuestionCategory string `json:"QuestionCategory"`
}

// get question from id
type Get_question_resp_json struct {
	QuestionId             int32    `json:"QuestionID"`
	QuestionCategory       string   `json:"QuestionCategory"`
	QuestionName           string   `json:"QuestionName"`
	QuestionDescription    string   `json:"QuestionDescription"`
	InputDescription       string   `json:"InputDescription"`
	OutputDescription      string   `json:"OutputDescription"`
	ConstraintsDescription string   `json:"ConstraintsDescription"`
	TimeConstraint         int32    `json:"TimeConstraint"`
	MemConstraint          int32    `json:"MemConstraint"`
	ExampleInputs          []string `json:"ExampleInputs"`
	ExampleOutputs         []string `json:"ExampleOutputs"`
}

type Get_question_req_json struct {
	QuestionId int `json:"QuestionId"`
}

// get all submisions api
type Get_question_submissions_resp_json struct {
	AllSubmissions []Submission_verdict_resp_json `json:"AllSubmissions"`
}
type Get_question_submissions_req_json struct {
	QuestionId int `json:"QuestionId"`
}
