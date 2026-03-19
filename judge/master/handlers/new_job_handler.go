package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"master/config"
	db "master/db/postgres_db"
	"master/types"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)


func validateRequest(user_req types.User_judge_req_json, queries *db.Queries) error {
	if _, ok := config.AllRuntimes[user_req.Runtime]; ok != true {
		return fmt.Errorf("Runtime does not exist :%s", user_req.Runtime)
	}
	if ok, err := queries.QuestionExists(context.Background(), int32(user_req.QuestionId)); err != nil || ok == false {
		return fmt.Errorf("question id does not exist %v with err :%v", user_req.QuestionId,err)
	}
	return nil
}

func (server *Server) Handle_new_job_req(w http.ResponseWriter, r *http.Request) {
	var user_req types.User_judge_req_json
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	decoder.DisallowUnknownFields()
	err := decoder.Decode(&user_req)
	io.Copy(io.Discard, r.Body)
	if err != nil {
		fmt.Println("could not read json with err : ", err)
		http.Error(w, "could not read json", http.StatusBadRequest)
		return
	}

	err = validateRequest(user_req, server.queries)
	if err != nil {
		fmt.Println("validation error with err : ", err)
		http.Error(w, fmt.Sprintf("invalid request with err: %v", err), http.StatusBadRequest)
		return
	}

	var user_id = 1
	worker_req, err := createNewWorkerJobRequest(user_req, user_id, server.queries)
	if err != nil {
		fmt.Println("could not create request with err : ", err)
		http.Error(w, "could not create job", http.StatusInternalServerError)
		return
	}

	// add job to the queue
	if err := server.Scedular.ProcessJob(worker_req); err != nil {
		fmt.Println("[info] job failed to added to pool questionid:", user_req.QuestionId)
	}

	resp := types.User_judge_resp_json{
		Verdict: "queued",
		Submission_id: worker_req.JobId,
	}
	resp_json,err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(resp_json)
}

func createNewWorkerJobRequest(user_req types.User_judge_req_json, user_id int, queries *db.Queries) (types.Worker_req_json, error) {

	ctx := context.Background()
	constraints, err := queries.GetTimeAndMemConstraints(ctx, int32(user_req.QuestionId))
	if err != nil {
		return types.Worker_req_json{}, err
	}

	params := db.CreateSubmissionAndReturnIdParams{
		UserID:     int32(user_id),
		QuestionID: int32(user_req.QuestionId),
		SubmissionTime: pgtype.Timestamp{
			Time:             time.Now(),
			InfinityModifier: 0,
			Valid:            true,
		},
		SubmitedCode: pgtype.Text{
			String: user_req.Code,
			Valid:  true,
		},
		CodeRuntime: pgtype.Text{String: user_req.Runtime,
			Valid: true,
		},
		Verdict: pgtype.Text{
			String: "queued",
			Valid:  true,
		},
	}

	Job_id, err := queries.CreateSubmissionAndReturnId(ctx, params)
	fmt.Println("new job id", Job_id)
	if err != nil {
		fmt.Println("failed to create submition err:", err)
		return types.Worker_req_json{}, err
	}

	return types.Worker_req_json{
		QuestionId:     user_req.QuestionId,
		Runtime:        user_req.Runtime,
		TimeConstrain:  int(constraints.TimeConstraint),
		MemConstrain:   int(constraints.MemConstraint),
		Code:           user_req.Code,
		JobId:          int(Job_id),
		InternalApiKey: config.InternalApiKey,
	}, nil
}
