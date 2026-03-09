package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"master/config"
	"master/db"
	"master/types"
	"net/http"
)

func (server *Server) Handle_new_job_req(w http.ResponseWriter,r *http.Request){
	var user_req types.User_req_json 
	json_body,err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("could not read body with err : ",err)
		http.Error(w,"could not read body ",http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(json_body,&user_req)
	if err != nil {
		fmt.Println("could not read json with err : ",err)
		http.Error(w,"could not read json",http.StatusBadRequest)
		return
	}

	worker_req := createNewWorkerJobRequest(user_req)

	fmt.Println("adding to queue")
	server.Scedular.ProcessJob(worker_req)
	w.WriteHeader(http.StatusAccepted)
	
}

func createNewWorkerJobRequest(user_req types.User_req_json) types.Worker_req_json {

	time_limit,mem_limit := db.GetTimeAndMemLimits(user_req.QuestionId)
	Job_id := get_new_job_id()
	return types.Worker_req_json{
		 QuestionId: user_req.QuestionId,
		 Runtime: user_req.Runtime,
		 TimeConstrain: time_limit,
		 MemConstrain: mem_limit,
		 Code: user_req.Code,
		 JobId: Job_id,
		 InternalApiKey: config.InternalApiKey,
	}

}

