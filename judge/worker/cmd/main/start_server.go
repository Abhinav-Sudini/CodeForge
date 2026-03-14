package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"worker/handlers"
	"worker/utils"
)

// can be docker or somthing else
const PORT = 8000

func initRuntimeEnv() {
	allTasksDir := os.Getenv("QUESTIONS_DIR")

	if ok, _ := utils.DirExists(allTasksDir); ok == false {
		panic("can not intit worker questions dir not found")
	}
}

func StallTillWorkerIsRunning(n int) error {
	// wait some time for the server to be set up
	time.Sleep(time.Second)
	for _ = range n {
		url := fmt.Sprintf("http://localhost:%d/healthcheck/", PORT)
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return nil
		}

		time.Sleep(time.Second * 2)
	}
	return errors.New("failed to start worker")
}

func StallTillMasterIsUp(n int) error {
	for _ = range n {
		err := handlers.PostWorkerIsFreeReqToMaster()
		if err == nil {
			fmt.Println("connected to master")
			return nil
		}
		fmt.Println("connect to master failes retry..")
		time.Sleep(time.Second * 2)
	}
	return errors.New("failed to connect to master")
}

func runWorker() error {

	initRuntimeEnv()
	//function to post worker is free ever x seconds

	http.HandleFunc("/judge/", handlers.Compile_and_judge_handler)
	http.HandleFunc("/healthcheck/", handlers.Healthcheck)

	fmt.Println("Starting server")
	fmt.Println("sering on http://localhost:", PORT)

	go func() {
		var max_trys = 10
		if err := StallTillWorkerIsRunning(max_trys); err != nil {
			panic(err)
		}
		if err := StallTillMasterIsUp(max_trys); err != nil {
			panic(err)
		}
	}()

	url_addr := "0.0.0.0:" + strconv.Itoa(PORT)
	return http.ListenAndServe(url_addr, nil)
}

func main() {
	if err := runWorker(); err != nil {
		fmt.Println("[Worker] exited with error : ", err)
		os.Exit(1)
	}
}
