package service

import (
	"context"

	"github.com/lonepeon/sport/internal/domain"
	"github.com/lonepeon/sport/internal/repository"
)

func GetRunningSession(repo repository.Reader, ctx context.Context, slug domain.RunningActivitySlug) (domain.RunningActivity, error) {
	activity, err := repo.GetRunningActivity(ctx, slug)
	if err != nil {
		return domain.RunningActivity{}, err
	}

	return activity, nil
}
