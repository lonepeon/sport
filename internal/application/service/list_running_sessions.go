package service

import (
	"context"

	"github.com/lonepeon/sport/internal/domain"
	"github.com/lonepeon/sport/internal/repository"
)

func ListRunningSessions(repo repository.Reader, ctx context.Context) ([]domain.RunningActivity, error) {
	return repo.ListRunningActivities(ctx)
}
