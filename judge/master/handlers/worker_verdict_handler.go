package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	db "master/db/postgres_db"
	"master/types"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Verdict int

const (
	VerdictAccepted Verdict = iota
	VerdictWrongAns
	VerdictCompilationError
	VerdictTLE
	VerdictMLE
	VerdictCodeToBig
	VerdictInternalError
	VerdictBadRequest
)

func (server *Server) Worker_verdict_handler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var verdict types.Worker_Response_json
	err := json.Unmarshal(body, &verdict)
	if err != nil {
		fmt.Println("failed to parse json err:", err)
	}
	job_id := verdict.JobId
	server.Scedular.RemoveFromOnGoing(job_id)

	err = UpdateSubmitionAndAddVerdictStats(verdict, server.db_conn, server.queries)
	if err != nil {
		http.Error(w,"code 500",http.StatusInternalServerError)
		return
	}
	fmt.Printf("verdict from worker for job_id:%v questionid:%v verdict:%v time:%v \n",
	verdict.JobId, verdict.QuestionId,  verdict.MSG, verdict.Time_ms)
}

func UpdateSubmitionAndAddVerdictStats(Verdict types.Worker_Response_json, db_con *pgx.Conn, q *db.Queries) error {

	ctx := context.Background()
	tx, err := db_con.Begin(ctx)
	if err != nil {
		return err
	}
	qtx := q.WithTx(tx)

	defer tx.Rollback(ctx)

	params := db.UpdateVerdictForSubmitionParams{
		SubmissionID: int32(Verdict.JobId),
		Verdict: pgtype.Text{
			String: Verdict.MSG,
			Valid:  true,
		},
	}
	err = qtx.UpdateVerdictForSubmition(ctx, params)
	if err != nil {
		return err
	}

	verdictParams := db.CreateVerdictStatsRecordParams{
		SubmissionID: int32(Verdict.JobId),
		TimeMs: int32(Verdict.Time_ms),
		MemUsage: int32(Verdict.Mem_usage),
		NotAcceptedTestCase: pgtype.Int4{
			Int32: int32(Verdict.WA_Test_case),
			Valid: true,
		},
		NotAcceptedTestCaseStdout: pgtype.Text{
			String: Verdict.WA_Stdout,
			Valid: true,
		},
		Stderr: pgtype.Text{
			String: Verdict.Stderr,
			Valid: true,
		},
	}

	err = qtx.CreateVerdictStatsRecord(ctx,verdictParams)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
