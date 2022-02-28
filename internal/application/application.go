package application

import (
	"context"
	"io"
	"time"

	"github.com/lonepeon/sport/internal/domain"
)

type Application interface {
	DeleteRunningSession(context.Context, domain.RunningActivitySlug) error
	GetRunningSession(context.Context, domain.RunningActivitySlug) (domain.RunningActivity, error)
	ListRunningSessions(context.Context) ([]domain.RunningActivity, error)
	TrackRunningSession(context.Context, time.Time, io.Reader) error
}
