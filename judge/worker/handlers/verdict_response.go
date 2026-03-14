package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
	"worker/config"
	"worker/types"
)

func PostResponseToMaster(verdict types.JudgeCodeResponse) error {
	//TODO
	fmt.Println("verdict from server", verdict)
	job_resp_body, err := json.Marshal(&verdict)
	if err != nil {
		return err
	}

	master_url := config.Master_url + config.VerdictApiLocation
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*config.Max_req_timeout)
	defer cancel()
	fmt.Println("tring to sent req to ",master_url)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, master_url, bytes.NewReader(job_resp_body))
	if err != nil {
		return err
	}

	fmt.Println("sending post :", string(job_resp_body))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return errors.New("job not accepted with MSG: " + string(body))
	}
	return nil
}

func PostWorkerIsFreeReqToMaster() error {
	fmt.Println("posting is free ")
	master_url := config.Master_url + config.WorkerRegisterApiLocation
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*config.Max_req_timeout)
	// defer cancel()
	ctx := context.Background()

	json_body := bytes.NewReader([]byte(`{
		"IP": "192.168.1.42",
		"Port": 8000
	}`))
	fmt.Println("sending is free post :to ",master_url)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, master_url, json_body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	fmt.Println(err)
	if err != nil {
		return err
	}
	fmt.Println(resp.StatusCode)
	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return errors.New("req not accepted with MSG: " + string(body))
	}
	return nil
}

func CronPostWorkerIsFree() {
	for {
		WorkerMutex.Lock()

		PostWorkerIsFreeReqToMaster()

		WorkerMutex.Unlock()
		time.Sleep(time.Second * 2)
	}
}

