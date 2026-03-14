package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"master/types"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func (server *Server) Get_verdict_handler(w http.ResponseWriter, r *http.Request) {
	var user_req types.Submission_verdict_req_json
	json_body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("could not read body with err : ", err)
		http.Error(w, "could not read body ", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(json_body, &user_req)
	if err != nil {
		fmt.Println("could not read json with err : ", err)
		http.Error(w, "could not read json", http.StatusBadRequest)
		return
	}

	verdict,err := server.queries.GetSubmissionVerdictAndQuestionid(context.Background(),int32(user_req.Submission_id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "can not find submission", http.StatusBadRequest)
			return
		}else{
			fmt.Println("querie failed with err: ",err)
			http.Error(w,"failed getting submission",http.StatusInternalServerError)
			return
		}
	}


	verdict_stats, err := server.queries.GetVerdictStats(context.Background(), int32(user_req.Submission_id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "can not find submission", http.StatusBadRequest)
			return
		}else{
			fmt.Println("querie failed with err: ",err)
			http.Error(w,"failed getting submission",http.StatusInternalServerError)
			return
		}
	}

	response := types.Submission_verdict_resp_json{
		Submission_id: int(verdict_stats.SubmissionID),
		QuestionId: int(verdict.QuestionID),
		Verdict: verdict.Verdict.String,
		Mem_usage: int(verdict_stats.MemUsage),
		Time_ms: int(verdict_stats.TimeMs),
		WA_Test_case: int(verdict_stats.NotAcceptedTestCase.Int32),
		WA_Stdout: verdict_stats.NotAcceptedTestCaseStdout.String,
		Stderr: verdict_stats.Stderr.String,
	}
	
	json_resp,err := json.Marshal(response)
	if err != nil {
		fmt.Println("json parse failed err: ",err)
			http.Error(w,"failed get response",http.StatusInternalServerError)
			return
	}
	
	w.Write(json_resp)
}
