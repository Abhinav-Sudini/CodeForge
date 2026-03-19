package handlers

import (
	"encoding/json"
	db "master/db/postgres_db"
	"master/types"
	"master/utils"
	"net/http"
	"strconv"
)

func (s *Server) GetQuestionSubmissionsHandler(w http.ResponseWriter, r *http.Request) {
	question_id,err := strconv.Atoi(r.PathValue("q_id"))
	if err != nil {
		http.Error(w,"no submissions_id pathvalue",http.StatusBadRequest)
		return
	}

	user_id := 1
	submissions, err := s.queries.GetAllSubmissionOfQuestion(r.Context(),
		db.GetAllSubmissionOfQuestionParams{
			QuestionID: int32(question_id),
			UserID:     int32(user_id),
		},
	)
	if err != nil {
		http.Error(w, "failed to get submissions", http.StatusInternalServerError)
		return
	}


	json_resp,err := json.Marshal(make_submissions_struct(submissions))
	if err != nil {
		http.Error(w, "failed to get submissions", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json_resp)
}

func make_submissions_struct(submissions []db.GetAllSubmissionOfQuestionRow) types.Get_question_submissions_resp_json {
	var temp_slice []types.Submission_verdict_resp_json
	for i := range submissions {
		q_stdin,q_stdout := utils.GetTestCaseStdinStdout(int(submissions[i].QuestionID),int(submissions[i].NotAcceptedTestCase.Int32))
		temp_slice = append(temp_slice,types.Submission_verdict_resp_json{
			Submission_id: int(submissions[i].SubmissionID),	
			QuestionId: int(submissions[i].QuestionID),
			Verdict: submissions[i].Verdict.String,
			Time_ms: int(submissions[i].TimeMs.Int32),
			Mem_usage: int(submissions[i].MemUsage.Int32),
			WA_Test_case: int(submissions[i].NotAcceptedTestCase.Int32),
			WA_Stdin: string(q_stdin),
			Required_Stdout: string(q_stdout),
			WA_Stdout: submissions[i].NotAcceptedTestCaseStdout.String,
			Stderr: submissions[i].Stderr.String,
			SubmissionTime: submissions[i].SubmissionTime.Time,
		})
	}
	return types.Get_question_submissions_resp_json{
		AllSubmissions: temp_slice,
	}
}
