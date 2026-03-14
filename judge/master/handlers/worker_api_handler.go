package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"master/config"
	"master/types"
	"net"
	"net/http"
)


func (server *Server) Register_worker_handler(w http.ResponseWriter,r *http.Request){
	worker,err := getWorkerInfo(r)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return 
	}
	fmt.Println("new worker added :",worker)
	err = server.Scedular.AddToWorkerPool(worker)
	if err != nil {
		fmt.Println("worker register failed with err: ",err)
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return 
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("worker registered\n"))
}

func getWorkerInfo(r *http.Request) (types.Worker_info,error) {
	body,_ := io.ReadAll(r.Body)
	var worker types.Worker_info
	err := json.Unmarshal(body,&worker)
	if err != nil {
		fmt.Println("can not read json err :",err)
		return worker,err
	}
	ip,_,_ := net.SplitHostPort(r.RemoteAddr)
	if _,ok := config.AllRuntimes[worker.Runtime];ok != true {
		return worker,fmt.Errorf("runtime does not exist")	
	}
	worker.IP = ip

	return worker,nil
}
