package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db_conn *pgxpool.Pool
var ctx_bgd context.Context

func init_db() {

	fmt.Println("connecting to db")
	var err error
	db_conn, err = pgxpool.New(ctx_bgd, os.Getenv("PG_DATABASE_URL"))
	if err != nil {
		fmt.Println("err con", err)
		panic(err)
	}
}

func run_migration() {
	sqlBytes, err := os.ReadFile("db/schema.sql")
	if err != nil {
		panic(err)
	}

	_,err = db_conn.Exec(ctx_bgd,string(sqlBytes))
	if err != nil {
		panic(err)
	}

}

func main() {
	ctx_bgd = context.Background()
	init_db()
	fmt.Println("db connected")
	run_migration()
	fmt.Println("migration done")
}
