package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"worker/handlers"
	"worker/utils"
)

// can be docker or somthing else
const worker_runtime = "dev"


func initRuntimeEnv(){
	allTasksDir := os.Getenv("QUESTIONS_DIR")

	if ok,_ := utils.DirExists(allTasksDir);ok == false {
		panic("can not intit worker")
	}
}

func runWorker() error {

	initRuntimeEnv()

	http.HandleFunc("/judge/", handlers.Compile_and_judge_handler)

	PORT := 8000

	fmt.Println("Starting server")
	fmt.Println("sering on http://localhost:", PORT)

	url_addr := "0.0.0.0:" + strconv.Itoa(PORT)
	err := http.ListenAndServe(url_addr, nil)
	return err

}

func main() {
	if err := runWorker(); err != nil {
		fmt.Println("[Worker] exited with error : ", err)
		os.Exit(1)
	}
}
