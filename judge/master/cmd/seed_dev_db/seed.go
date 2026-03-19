package main

import (
	"context"
	"fmt"
	db "master/db/postgres_db"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var queries *db.Queries
var ctx_bgd context.Context

func sedd_db() error {
	if err := seed_questions(); err != nil {
		return err
	}
	if err := seed_user(); err != nil {
		panic(err)
	}
	return nil
}

func seed_user() error {
	ctx := context.Background()

	seed := db.CreateUserAndReturnIdParams{
		Username: pgtype.Text{
			String: "dev",
			Valid:  true,
		},
		Email: pgtype.Text{
			String: "dev@example.com",
			Valid:  true,
		},
		EncryptedPassword: pgtype.Text{
			String: "",
			Valid:  true,
		},
		Role: pgtype.Text{
			String: "user",
			Valid:  true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
		LastOnline: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	}

	fmt.Println("run query")
	user_id, err := queries.CreateUserAndReturnId(ctx, seed)
	fmt.Println(err, user_id)
	fmt.Println("users seed done")

	return nil
}

func seed_questions() error {
	seed := db.CreateQuestionParams{
		QuestionID: 1681,
		QuestionCategory: pgtype.Text{
			String: "arrays",
			Valid:  true,
		},
		QuestionName:pgtype.Text{
			String: "question name",
			Valid:  true,
		},
		QuestionDescription: pgtype.Text{
			String: "Given an array of integers, return the maximum element.",
			Valid:  true,
		},
		InputDescription: pgtype.Text{
			String: "First line contains N. Second line contains N integers.",
			Valid:  true,
		},
		ConstraintsDescription: pgtype.Text{
			String: "1 <= N <= 10^5, -10^9 <= Ai <= 10^9",
			Valid:  true,
		},
		TimeConstraint:2000, 
		MemConstraint: 256000,
		ExampleInputs: []string{
			"5\n1 2 3 4 5",
		},
		ExampleOutputs: []string{
			"5",
		},
	}
	seed2 := db.CreateQuestionParams{
		QuestionID: 1069,
		QuestionCategory: pgtype.Text{
			String: "arrays",
			Valid:  true,
		},
		QuestionName:pgtype.Text{
			String: "question name2",
			Valid:  true,
		},
		QuestionDescription: pgtype.Text{
			String: "Given an array of integers, return the maximum element.",
			Valid:  true,
		},
		InputDescription: pgtype.Text{
			String: "First line contains N. Second line contains N integers.",
			Valid:  true,
		},
		ConstraintsDescription: pgtype.Text{
			String: "1 <= N <= 10^5, -10^9 <= Ai <= 10^9",
			Valid:  true,
		},
		TimeConstraint:2000, 
		MemConstraint: 256000,
		ExampleInputs: []string{
			"5\n1 2 3 4 5",
		},
		ExampleOutputs: []string{
			"5",
		},
	}

	fmt.Println("run query")
	ctx := context.Background()
	err := queries.CreateQuestion(ctx, seed)
	if err != nil {
		panic(err)
	}
	err = queries.CreateQuestion(ctx, seed2)
	if err != nil {
		panic(err)
	}

	fmt.Println("questions done")
	return nil
}

func init_db() {

	fmt.Println("connecting to db")
	conn, err := pgx.Connect(ctx_bgd, os.Getenv("PG_DATABASE_URL"))
	if err != nil {
		fmt.Println("err con", err)
		panic(err)
	}

	queries = db.New(conn)

}

func main() {
	ctx_bgd = context.Background()
	init_db()
	sedd_db()
}
