package handlers

import (
	"encoding/json"
	"master/types"
	"master/utils"
	"net/http"
	"strconv"
)

func (s *Server) Get_submission_verdict_handler(w http.ResponseWriter, r *http.Request) {
	submission_id, err := strconv.Atoi(r.PathValue("submission_id"))
	if err != nil {
		http.Error(w, "no submissions_id pathvalue", http.StatusBadRequest)
		return
	}

	submission, err := s.queries.GetSubmissionVerdict(r.Context(), int32(submission_id))
	if err != nil {
		http.Error(w, "failed to get submission", http.StatusInternalServerError)
		return
	}

	q_stdin, q_stdout := utils.GetTestCaseStdinStdout(int(submission.QuestionID), int(submission.NotAcceptedTestCase.Int32))

	return_struct := types.Submission_verdict_resp_json{
		Submission_id:   int(submission.SubmissionID),
		QuestionId:      int(submission.QuestionID),
		Verdict:         submission.Verdict.String,
		Time_ms:         int(submission.TimeMs.Int32),
		Mem_usage:       int(submission.MemUsage.Int32),
		WA_Test_case:    int(submission.NotAcceptedTestCase.Int32),
		WA_Stdin:        string(q_stdin),
		Required_Stdout: string(q_stdout),
		WA_Stdout:       submission.NotAcceptedTestCaseStdout.String,
		Stderr:          submission.Stderr.String,
		SubmissionTime:  submission.SubmissionTime.Time,
	}


	json_resp, err := json.Marshal(return_struct)
	if err != nil {
		http.Error(w, "failed to get submission", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json_resp)
}
