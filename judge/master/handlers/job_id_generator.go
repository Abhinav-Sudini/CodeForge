package handlers

var job_id = 1
func get_new_job_id() int {
	job_id++
	return job_id
}
