package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"worker/handlers"

	"github.com/joho/godotenv"
)

//can be docker or somthing else
const worker_runtime = "dev"

func runWorker() error {
	if worker_runtime == "dev" {
		err := godotenv.Load("worker.env")
		if err != nil {
			return errors.New("Failed to load env file")
		}
	}

	http.HandleFunc("/judge/",handlers.Compile_and_judge_handler)

	PORT,err := strconv.Atoi(os.Getenv("WORKER_PORT"))
	if err != nil {
		return err
	}

	fmt.Println("Starting server")
	fmt.Println("sering on http://localhost:", PORT)

	url_addr := "localhost:" + strconv.Itoa(PORT)
	err = http.ListenAndServe(url_addr,nil)
	return err

}

func main(){
	if err := runWorker(); err != nil {
		fmt.Println("[Worker] exited with error : ",err)
		os.Exit(1)
	}
}
