package job

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lonepeon/golib/job"
	"github.com/lonepeon/sport/internal/application"
	"github.com/lonepeon/sport/internal/domain"
)

const deleteRunningSessionJobName = "delete-running-session-job"

func EnqueueDeleteRunningSessionJob(client Enqueuer, input DeleteRunningSessionJobInput) error {
	j, err := job.NewJob(deleteRunningSessionJobName, input)
	if err != nil {
		return fmt.Errorf("can't build a new job (name=%s): %v", deleteRunningSessionJobName, err)
	}

	if err := client.Enqueue(j); err != nil {
		return fmt.Errorf("can't enqueue job (name=%s): %v", deleteRunningSessionJobName, err)
	}

	return nil
}

type DeleteRunningSessionJobInput struct {
	Slug string `json:"slug"`
}

type DeleteRunningSessionJob struct {
	application application.Application
}

func NewDeleteRunningSessionJob(app application.Application) *DeleteRunningSessionJob {
	return &DeleteRunningSessionJob{application: app}
}

func (j *DeleteRunningSessionJob) Name() string {
	return deleteRunningSessionJobName
}

func (j *DeleteRunningSessionJob) Handle(ctx context.Context, payload []byte) error {
	var input DeleteRunningSessionJobInput
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("can't parse input: %v", err)
	}

	slug, err := domain.NewRunnningActivitySlugFromString(input.Slug)
	if err != nil {
		return fmt.Errorf("can't parse slug: %v", err)
	}

	if err := j.application.DeleteRunningSession(ctx, slug); err != nil {
		return fmt.Errorf("can't delete running activity: %v", err)
	}

	return nil
}
