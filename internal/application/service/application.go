package service

import (
	"context"
	"io"
	"time"

	"github.com/lonepeon/sport/internal/domain"
	"github.com/lonepeon/sport/internal/repository"
)

type Application struct {
	repo repository.ReadWriter
}

func NewApplication(repo repository.ReadWriter) Application {
	return Application{repo: repo}
}

func (a Application) DeleteRunningSession(ctx context.Context, slug domain.RunningActivitySlug) error {
	return DeleteRunningSession(a.repo, ctx, slug)
}

func (a Application) GetRunningSession(ctx context.Context, slug domain.RunningActivitySlug) (domain.RunningActivity, error) {
	return GetRunningSession(a.repo, ctx, slug)
}

func (a Application) ListRunningSessions(ctx context.Context) ([]domain.RunningActivity, error) {
	return ListRunningSessions(a.repo, ctx)
}

func (a Application) TrackRunningSession(ctx context.Context, ranAt time.Time, file io.Reader) error {
	return TrackRunningSession(a.repo, ctx, ranAt, file)
}
