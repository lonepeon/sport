package repository

import (
	"context"
	"io"

	"github.com/lonepeon/sport/internal/domain"
)

type Logging interface {
	Info(string)
	Infof(string, ...interface{})
}

type Logger struct {
	logger Logging
	repo   ReadWriter
}

func NewLogger(log Logging, repo ReadWriter) Logger {
	return Logger{logger: log, repo: repo}
}

func (l Logger) StoreAsset(file io.Reader, fileName string) error {
	l.logger.Infof("repository uploads file %s", fileName)
	err := l.repo.StoreAsset(file, fileName)
	if err != nil {
		l.logger.Infof("repository failed to upload file: %v", err)
		return err
	}

	l.logger.Info("repository uploaded file")
	return nil
}

func (l Logger) DeleteAsset(fileName string) error {
	l.logger.Infof("repository deletes file %s", fileName)
	err := l.repo.DeleteAsset(fileName)
	if err != nil {
		l.logger.Infof("repository failed to delete file: %v", err)
		return err
	}

	l.logger.Info("repository deleted file")
	return nil
}

func (l Logger) GetRunningActivity(ctx context.Context, slug domain.RunningActivitySlug) (domain.RunningActivity, error) {
	l.logger.Infof("repository fetches running activity from slug %s", slug)
	activity, err := l.repo.GetRunningActivity(ctx, slug)

	if err != nil {
		l.logger.Infof("repository failed to find running activity: %v", err)
		return activity, err
	}

	l.logger.Infof("repository found running activity")
	return activity, nil
}

func (l Logger) ListRunningActivities(ctx context.Context) ([]domain.RunningActivity, error) {
	l.logger.Info("repository fetches all running activities")
	activities, err := l.repo.ListRunningActivities(ctx)
	if err != nil {
		l.logger.Infof("repository failed to find running activities: %v", err)
		return activities, err
	}

	l.logger.Infof("repository found %d running activities", len(activities))
	return activities, nil
}

func (l Logger) RecordRunningActivity(ctx context.Context, activity domain.RunningActivity) error {
	l.logger.Infof("repository records a new running activity at %s", activity.Slug)
	err := l.repo.RecordRunningActivity(ctx, activity)
	if err != nil {
		l.logger.Infof("repository failed to record the running activity: %v", err)
		return err
	}

	l.logger.Infof("repository recorded running activity")
	return nil
}

func (l Logger) DeleteRunningActivity(ctx context.Context, slug domain.RunningActivitySlug) error {
	l.logger.Infof("repository deletes running activity with slug %s", slug)
	if err := l.repo.DeleteRunningActivity(ctx, slug); err != nil {
		l.logger.Infof("repository failed to delete the running activity: %v", err)
		return err
	}

	l.logger.Infof("repository deleted running activity")
	return nil
}

func (l Logger) GenerateMap(ctx context.Context, gpx domain.GPXFile) (domain.MapFile, error) {
	l.logger.Info("repository generates map from points")
	mapFile, err := l.repo.GenerateMap(ctx, gpx)
	if err != nil {
		l.logger.Infof("repository failed to generate map from points: %v", err)
		return mapFile, err
	}

	l.logger.Info("repository generated map file")
	return mapFile, nil
}

func (l Logger) CleanGPXFile(ctx context.Context, r io.Reader) (domain.GPXFile, error) {
	l.logger.Info("repository cleans gpx file")
	gpx, err := l.repo.CleanGPXFile(ctx, r)
	if err != nil {
		l.logger.Infof("repository failed to clean gpx file: %v", err)
		return gpx, err
	}

	l.logger.Info("repository cleaned gpx file")
	return gpx, nil
}

func (l Logger) AnnotateMapWithStats(ctx context.Context, file domain.MapFile, distance domain.Distance, speed domain.Speed) (domain.ShareableMapFile, error) {
	return l.repo.AnnotateMapWithStats(ctx, file, distance, speed)
}
