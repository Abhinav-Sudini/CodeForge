package config

const (
	SERVER_PORT = 8080
	Server_url = "http://127.0.0.1:7000"

	Question_details_api_url = Server_url + "/api/question/"
  Get_question_submissions_api = Server_url + "/api/question/submissions/"
  Code_submission_api = Server_url + "/api/judge/"
  Submission_verdict_api = Server_url + "/api/submissions/"
)
