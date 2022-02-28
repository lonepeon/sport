package service_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/application/service"
	"github.com/lonepeon/sport/internal/domain"
	"github.com/lonepeon/sport/internal/domain/domaintest"
	"github.com/lonepeon/sport/internal/repository/repositorytest"
)

func TestTrackRunningSessionSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repositorytest.NewFake(t)

	gpxFileBytes := domaintest.GetGPXBytes()
	gpxFile := domaintest.NewGPXFile(t).WithFileContent(gpxFileBytes).Build()

	mapFileBytes := []byte("generated-map")
	mapFile := domain.NewMapFile(mapFileBytes)

	activity := domaintest.NewRunningActivity(t).
		WithDistanceMeters(gpxFile.Distance.Meters()).
		WithDuration(gpxFile.Duration).
		WithSpeedKmh(gpxFile.Speed.KilometersPerHour()).
		Build()

	ctx := context.Background()

	repo.OverrideCleanGPXFile(gpxFileBytes, gpxFile, nil)
	repo.OverrideGenerateMap(gpxFile, mapFile, nil)

	repo.ExpectCleanGPXFiles(gpxFileBytes)
	repo.ExpectGenerateMaps(gpxFile)
	repo.ExpectAnnotateMapsWithStats(mapFile)
	repo.ExpectStoreAssets(activity.GPXPath.String(), activity.MapPath.String(), activity.ShareableMapPath.String())
	repo.ExpectRecordActivities(activity)

	err := service.TrackRunningSession(repo, ctx, activity.RanAt, bytes.NewBuffer(gpxFileBytes))
	testutils.AssertNoError(t, err, "can't create running session")
}
