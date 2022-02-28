package job

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/lonepeon/golib/job"
	"github.com/lonepeon/sport/internal/application"
)

// TrackRunningSessionJobName is the name of the TrackRunningSessionJob
const trackRunningSessionJobName = "track-running-session-job"

// EnqueueTrackRunningSessionJob enqueues a new job
func EnqueueTrackRunningSessionJob(client Enqueuer, input TrackRunningSessionJobInput) error {
	j, err := job.NewJob(trackRunningSessionJobName, input)
	if err != nil {
		return fmt.Errorf("can't build a new job (name=%s): %v", trackRunningSessionJobName, err)
	}

	if err := client.Enqueue(j); err != nil {
		return fmt.Errorf("can't enqueue job (name=%s): %v", trackRunningSessionJobName, err)
	}

	return nil
}

// TrackRunningSessionJobInput represents a job input
type TrackRunningSessionJobInput struct {
	When        time.Time `json:"when"`
	GPXFilepath string    `json:"filepath"`
}

// TrackRunningSessionJob represent a tracker worker in charge of parsing and storing a running session
type TrackRunningSessionJob struct {
	application application.Application
}

// NewTrackRunningSessionJob initializes a running session job handler
func NewTrackRunningSessionJob(app application.Application) *TrackRunningSessionJob {
	return &TrackRunningSessionJob{application: app}
}

func (j *TrackRunningSessionJob) Name() string {
	return trackRunningSessionJobName
}

// Handle implements job.Handler
func (j *TrackRunningSessionJob) Handle(ctx context.Context, payload []byte) error {
	var input TrackRunningSessionJobInput
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("can't parse input: %v", err)
	}

	f, err := os.Open(input.GPXFilepath)
	if err != nil {
		return fmt.Errorf("can't open gpxfile (path=%s): %v", input.GPXFilepath, err)
	}
	defer f.Close()

	if err := j.application.TrackRunningSession(ctx, input.When, f); err != nil {
		return fmt.Errorf("can'track running session: %v", err)
	}

	return nil
}
