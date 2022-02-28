package repository_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/domain"
	"github.com/lonepeon/sport/internal/domain/domaintest"
	"github.com/lonepeon/sport/internal/repository"
	"github.com/lonepeon/sport/internal/repository/repositorytest"
)

type FakeLogger struct {
	Infos []string
}

func (f *FakeLogger) Infof(msg string, vars ...interface{}) {
	f.Infos = append(f.Infos, fmt.Sprintf(msg, vars...))
}

func (f *FakeLogger) Info(msg string) {
	f.Infos = append(f.Infos, msg)
}

func TestStoreAssetSuccess(t *testing.T) {
	repo := repositorytest.NewFake(t)
	log := FakeLogger{}

	err := repository.NewLogger(&log, repo).StoreAsset(bytes.NewBuffer(nil), "myfile.txt")
	testutils.AssertNoError(t, err, "unexpected repository error")

	testutils.AssertEqualInt(t, 2, len(log.Infos), "unexpected number of info message")
	testutils.AssertContainsString(t, "uploads", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, "myfile.txt", log.Infos[0], "unexpected file name in info message")
	testutils.AssertContainsString(t, "uploaded", log.Infos[1], "unexpected info message")
}

func TestStoreAssetError(t *testing.T) {
	repo := repositorytest.NewFake(t)
	log := FakeLogger{}
	expectedErr := errors.New("boom")

	repo.OverrideStoreAsset("myfile.txt", expectedErr)

	err := repository.NewLogger(&log, repo).StoreAsset(nil, "myfile.txt")
	testutils.AssertErrorIs(t, expectedErr, err, "expected repository error")

	testutils.AssertEqualInt(t, 2, len(log.Infos), "unexpected number of info message")
	testutils.AssertContainsString(t, "uploads", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, "myfile.txt", log.Infos[0], "unexpected file name in info message")
	testutils.AssertContainsString(t, "failed", log.Infos[1], "unexpected error message")
	testutils.AssertContainsString(t, err.Error(), log.Infos[1], "unexpected error message")
}

func TestDeleteAssetSuccess(t *testing.T) {
	repo := repositorytest.NewFake(t)
	log := FakeLogger{}

	repo.ExpectDeleteAssets("myfile.txt")

	err := repository.NewLogger(&log, repo).DeleteAsset("myfile.txt")
	testutils.AssertNoError(t, err, "unexpected repository error")

	testutils.AssertEqualInt(t, 2, len(log.Infos), "unexpected number of info message")
	testutils.AssertContainsString(t, "deletes", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, "myfile.txt", log.Infos[0], "unexpected file name in info message")
	testutils.AssertContainsString(t, "deleted", log.Infos[1], "unexpected info message")
}

func TestDeleteAssetError(t *testing.T) {
	repo := repositorytest.NewFake(t)
	log := FakeLogger{}
	expectedErr := errors.New("boom")

	repo.OverrideDeleteAsset("myfile.txt", expectedErr)

	err := repository.NewLogger(&log, repo).DeleteAsset("myfile.txt")
	testutils.AssertErrorIs(t, expectedErr, err, "expected repository error")

	testutils.AssertEqualInt(t, 2, len(log.Infos), "unexpected number of info message")
	testutils.AssertContainsString(t, "deletes", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, "myfile.txt", log.Infos[0], "unexpected file name in info message")
	testutils.AssertContainsString(t, "failed", log.Infos[1], "unexpected error message")
	testutils.AssertContainsString(t, err.Error(), log.Infos[1], "unexpected error message")
}

func TestGetRunningActivitySuccess(t *testing.T) {
	repo := repositorytest.NewFake(t)
	log := FakeLogger{}
	expected := domaintest.NewRunningActivity(t).Persist(repo)

	actual, err := repository.NewLogger(&log, repo).GetRunningActivity(context.Background(), expected.Slug)
	testutils.AssertNoError(t, err, "unexpected repository error")

	domaintest.AssertEqualRunningActivity(t, expected, actual, "unexpected activity")
	testutils.AssertEqualInt(t, 2, len(log.Infos), "unexpected number of info message")
	testutils.AssertContainsString(t, "fetches", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, expected.Slug.String(), log.Infos[0], "unexpected file name in info message")
	testutils.AssertContainsString(t, "found", log.Infos[1], "unexpected info message")
}

func TestGetRunningActivityError(t *testing.T) {
	repo := repositorytest.NewFake(t)
	log := FakeLogger{}
	activity := domaintest.NewRunningActivity(t).Persist(repo)
	expectedErr := errors.New("boom")

	repo.OverrideGetActivity(activity.Slug, expectedErr)

	_, err := repository.NewLogger(&log, repo).GetRunningActivity(context.Background(), activity.Slug)
	testutils.AssertErrorIs(t, expectedErr, err, "expected repository error")

	testutils.AssertEqualInt(t, 2, len(log.Infos), "unexpected number of info message")
	testutils.AssertContainsString(t, "fetches", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, activity.Slug.String(), log.Infos[0], "unexpected file name in info message")
	testutils.AssertContainsString(t, "failed to find", log.Infos[1], "unexpected info message")
	testutils.AssertContainsString(t, err.Error(), log.Infos[1], "unexpected info message")
}

func TestListRunningActivitiesSuccess(t *testing.T) {
	repo := repositorytest.NewFake(t)
	log := FakeLogger{}
	domaintest.NewRunningActivity(t).WithRawSlug("202102182208").Persist(repo)
	domaintest.NewRunningActivity(t).WithRawSlug("202202182208").Persist(repo)
	domaintest.NewRunningActivity(t).WithRawSlug("202302182208").Persist(repo)

	activities, err := repository.NewLogger(&log, repo).ListRunningActivities(context.Background())
	testutils.AssertNoError(t, err, "unexpected repository error")

	testutils.AssertEqualInt(t, 3, len(activities), "unexpected number of activities")
	testutils.AssertEqualInt(t, 2, len(log.Infos), "unexpected number of info message")
	testutils.AssertContainsString(t, "fetches", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, "found", log.Infos[1], "unexpected info message")
	testutils.AssertContainsString(t, strconv.Itoa(3), log.Infos[1], "unexpected info message")
}

func TestListRunningActivitiesError(t *testing.T) {
	repo := repositorytest.NewFake(t)
	log := FakeLogger{}
	expectedErr := errors.New("boom")

	repo.OverrideListActivities(expectedErr)

	_, err := repository.NewLogger(&log, repo).ListRunningActivities(context.Background())
	testutils.AssertErrorIs(t, expectedErr, err, "expected repository error")

	testutils.AssertEqualInt(t, 2, len(log.Infos), "unexpected number of info message")
	testutils.AssertContainsString(t, "fetches", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, "failed to find", log.Infos[1], "unexpected info message")
	testutils.AssertContainsString(t, err.Error(), log.Infos[1], "unexpected info message")
}

func TestRecordRunningActivitySuccess(t *testing.T) {
	repo := repositorytest.NewFake(t)
	log := FakeLogger{}
	activity := domaintest.NewRunningActivity(t).WithRawSlug("202202190640").Build()

	err := repository.NewLogger(&log, repo).RecordRunningActivity(context.Background(), activity)
	testutils.AssertNoError(t, err, "unexpected repository error")

	testutils.AssertEqualInt(t, 2, len(log.Infos), "unexpected number of info message")
	testutils.AssertContainsString(t, "records", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, "202202190640", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, "recorded", log.Infos[1], "unexpected info message")
}

func TestRecordRunningActivityError(t *testing.T) {
	repo := repositorytest.NewFake(t)
	log := FakeLogger{}
	activity := domaintest.NewRunningActivity(t).WithRawSlug("202202190640").Build()
	expectedErr := errors.New("boom")

	repo.OverrideRecordActivity(activity.Slug, expectedErr)

	err := repository.NewLogger(&log, repo).RecordRunningActivity(context.Background(), activity)
	testutils.AssertErrorIs(t, expectedErr, err, "expected repository error")

	testutils.AssertEqualInt(t, 2, len(log.Infos), "unexpected number of info message")
	testutils.AssertContainsString(t, "records", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, "202202190640", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, "failed to record", log.Infos[1], "unexpected info message")
	testutils.AssertContainsString(t, err.Error(), log.Infos[1], "unexpected info message")
}

func TestDeleteRunningActivitySuccess(t *testing.T) {
	repo := repositorytest.NewFake(t)
	log := FakeLogger{}
	activity := domaintest.NewRunningActivity(t).WithRawSlug("202202190640").Persist(repo)

	err := repository.NewLogger(&log, repo).DeleteRunningActivity(context.Background(), activity.Slug)
	testutils.AssertNoError(t, err, "unexpected repository error")

	testutils.AssertEqualInt(t, 2, len(log.Infos), "unexpected number of info message")
	testutils.AssertContainsString(t, "deletes", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, "202202190640", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, "deleted", log.Infos[1], "unexpected info message")
}

func TestDeleteRunningActivityError(t *testing.T) {
	repo := repositorytest.NewFake(t)
	log := FakeLogger{}
	activity := domaintest.NewRunningActivity(t).WithRawSlug("202202190640").Build()
	expectedErr := errors.New("boom")

	repo.OverrideDeleteActivity(activity.Slug, expectedErr)

	err := repository.NewLogger(&log, repo).DeleteRunningActivity(context.Background(), activity.Slug)
	testutils.AssertErrorIs(t, expectedErr, err, "expected repository error")

	testutils.AssertEqualInt(t, 2, len(log.Infos), "unexpected number of info message")
	testutils.AssertContainsString(t, "deletes", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, "202202190640", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, "failed to delete", log.Infos[1], "unexpected info message")
	testutils.AssertContainsString(t, err.Error(), log.Infos[1], "unexpected info message")
}

func TestGenerateMapSuccess(t *testing.T) {
	repo := repositorytest.NewFake(t)
	log := FakeLogger{}
	gpx := domaintest.NewGPXFile(t).Build()

	mapFile, err := repository.NewLogger(&log, repo).GenerateMap(context.Background(), gpx)
	testutils.AssertNoError(t, err, "unexpected repository error")
	testutils.AssertNotEqualNil(t, mapFile.File(), "unexpected result")

	testutils.AssertEqualInt(t, 2, len(log.Infos), "unexpected number of info message")
	testutils.AssertContainsString(t, "generates", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, "generated", log.Infos[1], "unexpected info message")
}

func TestGenerateMapError(t *testing.T) {
	repo := repositorytest.NewFake(t)
	log := FakeLogger{}
	expectedErr := errors.New("boom")
	gpx := domaintest.NewGPXFile(t).Build()

	repo.OverrideGenerateMap(gpx, domain.MapFile{}, expectedErr)

	_, err := repository.NewLogger(&log, repo).GenerateMap(context.Background(), gpx)
	testutils.AssertErrorIs(t, expectedErr, err, "expected repository error")

	testutils.AssertEqualInt(t, 2, len(log.Infos), "unexpected number of info message")
	testutils.AssertContainsString(t, "generates", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, "failed to generate", log.Infos[1], "unexpected info message")
	testutils.AssertContainsString(t, err.Error(), log.Infos[1], "unexpected info message")
}

func TestCleanGPXFileSuccess(t *testing.T) {
	repo := repositorytest.NewFake(t)
	log := FakeLogger{}
	expectedContent := []byte("gpx file")

	repo.ExpectCleanGPXFiles(expectedContent)

	actualGPX, err := repository.NewLogger(&log, repo).CleanGPXFile(context.Background(), bytes.NewBuffer(expectedContent))
	testutils.AssertNoError(t, err, "unexpected repository error")

	actualContent, err := ioutil.ReadAll(actualGPX.File())
	testutils.AssertNoError(t, err, "unexpected error while reading cleaned gpx file")

	testutils.AssertEqualString(t, string(expectedContent), string(actualContent), "unexpected result")

	testutils.AssertEqualInt(t, 2, len(log.Infos), "unexpected number of info message")
	testutils.AssertContainsString(t, "cleans", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, "cleaned", log.Infos[1], "unexpected info message")
}

func TestCleanGPXFileError(t *testing.T) {
	repo := repositorytest.NewFake(t)
	log := FakeLogger{}
	expectedErr := errors.New("boom")
	content := []byte("gpx file")

	repo.OverrideCleanGPXFile(content, domain.GPXFile{}, expectedErr)

	_, err := repository.NewLogger(&log, repo).CleanGPXFile(context.Background(), bytes.NewBuffer(content))
	testutils.AssertErrorIs(t, expectedErr, err, "unexpected repository error")

	testutils.AssertEqualInt(t, 2, len(log.Infos), "unexpected number of info message")
	testutils.AssertContainsString(t, "cleans", log.Infos[0], "unexpected info message")
	testutils.AssertContainsString(t, "failed to clean", log.Infos[1], "unexpected info message")
}
