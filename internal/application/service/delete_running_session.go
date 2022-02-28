package service

import (
	"context"
	"fmt"

	"github.com/lonepeon/sport/internal/domain"
	"github.com/lonepeon/sport/internal/repository"
)

func DeleteRunningSession(repo repository.ReadWriter, ctx context.Context, slug domain.RunningActivitySlug) error {
	activity, err := repo.GetRunningActivity(ctx, slug)
	if err != nil {
		return fmt.Errorf("can't find run activity %s: %w", slug, err)
	}

	if err := repo.DeleteAsset(activity.GPXPath.String()); err != nil {
		return fmt.Errorf("can't delete gpx file %s for run %s: %w", activity.GPXPath, slug, err)
	}

	if err := repo.DeleteAsset(activity.MapPath.String()); err != nil {
		return fmt.Errorf("can't delete map file %s for run %s: %w", activity.MapPath, slug, err)
	}

	if err := repo.DeleteAsset(activity.ShareableMapPath.String()); err != nil {
		return fmt.Errorf("can't delete shareable map file %s for run %s: %w", activity.ShareableMapPath, slug, err)
	}

	if err := repo.DeleteRunningActivity(ctx, slug); err != nil {
		return fmt.Errorf("can't delete activity: %w", err)
	}

	return nil
}
