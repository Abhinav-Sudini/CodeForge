package main

import (
	"context"
	"fmt"
	"master/config"
	"master/handlers"
	"net/http"
	"strconv"
)

var Middleware = []func(http.HandlerFunc)http.HandlerFunc {
	setCorsHeaderMiddleware,
}

func setCorsHeaderMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return 
		}

		next(w,r)
	}
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
	http.HandleFunc("/api/judge/", wrapMiddleware(server.Handle_new_job_req))
	http.HandleFunc("/api/submissions/{submission_id}/", wrapMiddleware(server.Get_submission_verdict_handler))

	http.HandleFunc("/api/question/", wrapMiddleware(server.GetAllQuestionsDetailsHandler))
	http.HandleFunc("/api/question/{q_id}/", wrapMiddleware(server.GetQuestionHandler))
	http.HandleFunc("/api/question/submissions/{q_id}/", wrapMiddleware(server.GetQuestionSubmissionsHandler))

	http.HandleFunc("/api/worker/register/", wrapMiddleware(server.Register_worker_handler))
	http.HandleFunc("/api/worker/verdict/", wrapMiddleware(server.Worker_verdict_handler))

	listend_add := "0.0.0.0" + ":" + strconv.Itoa(config.Server_Port)
	fmt.Println("running server at : ", listend_add)
	if err := http.ListenAndServe(listend_add, nil); err != nil {
		fmt.Println("master exiting with err :", err)
	}
}
