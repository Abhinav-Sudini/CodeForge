package jobschedular

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"master/config"
	"master/types"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type JobScedular struct {
	Job_queue_channels    map[string]chan (types.Worker_req_json)
	Worker_pool_channels map[string]chan (types.Worker_info)
	Ongoing_jobs         map[int]types.Running_job_info // maps job_id to assigned worker
	mutex                sync.RWMutex
}

func (scedular *JobScedular) ProcessJob(job types.Worker_req_json) error {
	select {
	case scedular.Job_queue_channels[job.Runtime] <- job:
		return nil
	default:
		return errors.New("buffer full")
	}
}

func (scedular *JobScedular) AddToWorkerPool(worker types.Worker_info) error {

	select {
	case scedular.Worker_pool_channels[worker.Runtime] <- worker:
		return nil
	default:
		return errors.New("buffer full")
	}
}

func (scedular *JobScedular) RemoveFromOnGoing(job_id int) error {
	scedular.mutex.Lock()
	defer scedular.mutex.Unlock()

	_, ok := scedular.Ongoing_jobs[job_id]
	if ok == false {
		return errors.New("job not found")
	}

	delete(scedular.Ongoing_jobs, job_id)
	return nil
}

func (scedular *JobScedular) AddOnGoingJob(job *types.Worker_req_json, worker *types.Worker_info) {
	scedular.mutex.Lock()
	defer scedular.mutex.Unlock()

	scedular.Ongoing_jobs[job.JobId] = types.Running_job_info{
		Job:             *job,
		Assigned_worker: *worker,
		Assigned_time:   time.Now(),
	}
}

func (scedular *JobScedular) StartSchedular(ctx context.Context) {
	for runtime := range config.AllRuntimes {
		go func() {
			for {
				// if both channels have a job and a worker then run
				select {
				case job := <-scedular.Job_queue_channels[runtime]:
					select {
					case worker := <-scedular.Worker_pool_channels[runtime]:
						fmt.Printf("sending req to worker question_id: %v and job_id: %v  to worker at %v \n", job.QuestionId, job.JobId, worker.IP)

						scedular.AddOnGoingJob(&job, &worker)
						err := scedular.executeJob(&job, worker)
						if err != nil {
							if _,ok := err.(types.WorkerBusyError); ok==true {
								fmt.Println("[schedular alloc]","worker to busy err for worker ip:",worker.IP)
								scedular.Worker_pool_channels[runtime] <- worker
							}
							//add job back to the queue
							scedular.Job_queue_channels[runtime] <- job
							scedular.RemoveFromOnGoing(job.JobId)
							fmt.Println("job failed with err: ", err)
						}
					case <-ctx.Done():
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}
}

func (scheduler *JobScedular) executeJob(job *types.Worker_req_json, worker types.Worker_info) error {
	if err := sendWorkerJobRequest(job, worker); err != nil {
		return err
	}

	//worker will send the verdict to /worker/api when the request is done
	//so return
	return nil
}

func sendWorkerJobRequest(job *types.Worker_req_json, worker types.Worker_info) error {

	worker_url := "http://" + worker.IP + ":" + strconv.Itoa(worker.Port) + "/judge/"
	job_req_body, err := json.Marshal(&job)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*config.Max_request_wait_time)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, worker_url, bytes.NewReader(job_req_body))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return types.WorkerBusyError{}
	}

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return errors.New("job not accepted with MSG: " + string(body))
	}

	return nil
}
