package service

import (
	"context"
	"fmt"
	"io"
	"path"
	"time"

	"github.com/lonepeon/sport/internal/domain"
	"github.com/lonepeon/sport/internal/repository"
)

func TrackRunningSession(repo repository.Writer, ctx context.Context, when time.Time, gpxFile io.Reader) error {
	gpx, err := repo.CleanGPXFile(ctx, gpxFile)
	if err != nil {
		return fmt.Errorf("can't load gpx file: %v", err)
	}

	imageMap, err := repo.GenerateMap(ctx, gpx)
	if err != nil {
		return fmt.Errorf("can't generate image from gpx: %v", err)
	}

	shareableMap, err := repo.AnnotateMapWithStats(ctx, imageMap, gpx.Distance, gpx.Speed)
	if err != nil {
		return fmt.Errorf("can't generate shareable image from map: %v", err)
	}

	basePath := path.Join("runs", when.Format("2006-01-02.15h04"))
	mapPath := path.Join(basePath, "map.png")
	shareableMapPath := path.Join(basePath, "share-map.png")
	gpxPath := path.Join(basePath, "run.gpx")

	activity, err := domain.NewRunningActivity(
		when,
		gpx.Duration,
		gpx.Distance,
		gpx.Speed,
		domain.GPXFilePath(gpxPath),
		domain.MapFilePath(mapPath),
		domain.ShareableMapFilePath(shareableMapPath),
	)
	if err != nil {
		return fmt.Errorf("can't build activity: %v", err)
	}

	assets := map[string]io.Reader{
		activity.MapPath.String():          imageMap.File(),
		activity.ShareableMapPath.String(): shareableMap.File(),
		activity.GPXPath.String():          gpx.File(),
	}

	if err := uploadPNGs(repo, assets); err != nil {
		return err
	}

	err = repo.RecordRunningActivity(ctx, activity)
	if err != nil {
		return fmt.Errorf("can't persists run: %v", err)
	}

	return nil
}

func uploadPNGs(repo repository.Writer, assets map[string]io.Reader) error {
	for assetPath, assetContent := range assets {
		if err := repo.StoreAsset(assetContent, assetPath); err != nil {
			return fmt.Errorf("can't store png file (png=%s): %v", assetPath, err)
		}
	}

	return nil
}
