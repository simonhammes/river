package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{})
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		args := SortArgs{
			Strings: []string{
				"whale", "tiger", "bear",
			},
		}

		_, err = riverClient.Insert(ctx, args, nil)

		if err != nil {
			// TODO
		}

		fmt.Fprint(w, "Home")
	})

	http.ListenAndServe("localhost:8000", mux)
}

type SortArgs struct {
	Strings []string `json:"strings"`
}

func (SortArgs) Kind() string {
	return "sort"
}

type SortWorker struct {
	river.WorkerDefaults[SortArgs]
}

func (w *SortWorker) Work(ctx context.Context, job *river.Job[SortArgs]) error {
	sort.Strings(job.Args.Strings)
	fmt.Printf("Sorted strings: %+v\n", job.Args.Strings)
	return nil
}
