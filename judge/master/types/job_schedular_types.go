package types

import "time"

type Running_job_info struct {
	Job             Worker_req_json
	Assigned_worker Worker_info
	Assigned_time   time.Time
}
