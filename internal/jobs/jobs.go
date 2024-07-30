package jobs

import (
	"context"
	"fmt"
	"sort"

	"github.com/riverqueue/river"
)

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
