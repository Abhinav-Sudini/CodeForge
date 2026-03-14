package handlers

import (
	"fmt"
	"master/config"
	db "master/db/postgres_db"
	jobschedular "master/job_schedular"

	"context"
	"master/types"
	"os"

	"github.com/jackc/pgx/v5"
)

type Server struct {
	Scedular jobschedular.JobScedular
	db_conn *pgx.Conn
	queries *db.Queries
}

func NewServer() (*Server,error) {
	ctx_bgd := context.Background()
	db_con,err := pgx.Connect(ctx_bgd, os.Getenv("PG_DATABASE_URL"))
	if err != nil {
		return nil,err
	}
	fmt.Println("db connected")

	queries := db.New(db_con)

	var job_chanels = make(map[string]chan(types.Worker_req_json))
	var worker_chanels = make(map[string]chan(types.Worker_info))
	for runtime := range config.AllRuntimes {
		job_chanels[runtime] = make(chan types.Worker_req_json,100)
		worker_chanels[runtime] = make(chan types.Worker_info,100)
	}

	return &Server{
		Scedular: jobschedular.JobScedular{
			Job_queue_channels:   job_chanels,
			Worker_pool_channels: worker_chanels,
			Ongoing_jobs:        make(map[int]types.Running_job_info,100),
		},
		db_conn: db_con,
		queries: queries,
	},nil
}
