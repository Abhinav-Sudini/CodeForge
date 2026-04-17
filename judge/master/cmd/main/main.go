package main

import (
	"context"
	"fmt"
	"master/config"
	"master/handlers"
	"master/middleware"
	"net/http"
	"strconv"
)

var Middleware = []func(http.HandlerFunc)http.HandlerFunc {
	middleware.SetCorsHeaderMiddleware,
}


func wrapMiddleware(h http.HandlerFunc) http.HandlerFunc{
	return_handler := h
	for _,m := range Middleware{
		return_handler = m(return_handler)
	}
	return return_handler
}

func main() {

	server, err := handlers.NewServer()
	if err != nil {
		panic(err)
	}
	defer server.Close()

	go server.Scedular.StartSchedular(context.Background())

	http.Handle("/", http.FileServer(http.Dir("/static/html")))

	// apis
	http.HandleFunc("/api/judge/", wrapMiddleware(middleware.WithJWTAuth(server.Handle_new_job_req)))
	http.HandleFunc("/api/submissions/{submission_id}/", wrapMiddleware(server.Get_submission_verdict_handler))

	http.HandleFunc("/api/question/", wrapMiddleware(server.GetAllQuestionsDetailsHandler))
	http.HandleFunc("/api/question/{q_id}/", wrapMiddleware(server.GetQuestionHandler))
	http.HandleFunc("/api/question/submissions/{q_id}/", wrapMiddleware(middleware.WithJWTAuth(server.GetQuestionSubmissionsHandler)))

	http.HandleFunc("/api/worker/register/", wrapMiddleware(server.Register_worker_handler))
	http.HandleFunc("/api/worker/verdict/", wrapMiddleware(server.Worker_verdict_handler))

	http.HandleFunc("/auth/register/", wrapMiddleware(server.RegisterUser))
	http.HandleFunc("/auth/login/", wrapMiddleware(server.UserLogin))

	listend_add := "0.0.0.0" + ":" + strconv.Itoa(config.Server_Port)
	fmt.Println("running server at : ", listend_add)
	if err := http.ListenAndServe(listend_add, nil); err != nil {
		fmt.Println("master exiting with err :", err)
	}
}
