package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/simonhammes/river/internal/jobs"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error: Could not load .env file")
	}

	ctx := context.Background()

	workers := river.NewWorkers()
	river.AddWorker(workers, &jobs.SortWorker{})

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

	go forever()
	select {}
}

func forever() {
	for {
		time.Sleep(time.Second)
	}
}
