package main

import (
	"context"
	"fmt"
	"master/config"
	"master/handlers"
	"net/http"
	"strconv"
)

func main(){

	server,err := handlers.NewServer()
	if err != nil {
		panic(err)
	}
	go server.Scedular.StartSchedular(context.Background())

	http.HandleFunc("/judge/",server.Handle_new_job_req)
	http.HandleFunc("/api/worker/register/",server.Register_worker_handler)
	http.HandleFunc("/api/worker/verdict/",server.Worker_verdict_handler)

	listend_add := "0.0.0.0" + ":" + strconv.Itoa(config.Server_Port)
	fmt.Println("running server at : ",listend_add)
	http.ListenAndServe(listend_add,nil)
}
