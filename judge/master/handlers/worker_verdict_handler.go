package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"master/types"
	"net/http"
)


func (server *Server) Worker_verdict_handler(w http.ResponseWriter,r *http.Request){
	body,_ := io.ReadAll(r.Body)
	var verdict types.Worker_Response_json
	err := json.Unmarshal(body,&verdict)
	if err != nil {
		fmt.Println("failed to parse json err:",err)
	}
	job_id := verdict.JobId
	server.Scedular.RemoveFromOnGoing(job_id)
	fmt.Println("verdict := ",verdict)
}

