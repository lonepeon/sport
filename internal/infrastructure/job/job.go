package job

import "github.com/lonepeon/golib/job"

// Enqueuer represents a client able to enqueue a job
type Enqueuer interface {
	Enqueue(job.Job) error
}
