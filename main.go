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
		log.Fatal("Error: Could not load .env file")
	}

	ctx := context.Background()

	workers := river.NewWorkers()
	river.AddWorker(workers, &SortWorker{})

	dbPool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		// handle error
	}

	riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		},
		Workers: workers,
	})

	if err != nil {
		// handle error
	}

	// Run the client inline. All executed jobs will inherit from ctx:
	if err := riverClient.Start(ctx); err != nil {
		// handle error
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		_, err = riverClient.Insert(ctx, SortArgs{
			Strings: []string{
				"whale", "tiger", "bear",
			},
		}, nil)

		if err != nil {
			// handle error
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
