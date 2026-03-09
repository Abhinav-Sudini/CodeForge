package handlers

import (
	jobschedular "master/job_schedular"
	"master/types"
)

type Server struct {
	Scedular jobschedular.JobScedular
}

func NewServer() *Server {
	return &Server{
		Scedular: jobschedular.JobScedular{
			Job_queue_channel:   make(chan types.Worker_req_json,100),
			Worker_pool_channel: make(chan types.Worker_info,100),
			Ongoing_jobs:        make(map[int]types.Running_job_info,100),
		},
	}
}
