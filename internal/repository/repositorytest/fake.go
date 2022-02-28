package repositorytest

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"testing"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/domain"
	"github.com/lonepeon/sport/internal/domain/domaintest"
	"github.com/lonepeon/sport/internal/infrastructure/mapbox/mapboxtest"
)

type RunningActivityErrorResponse struct {
	Slug domain.RunningActivitySlug
	Err  error
}

type AssetErrorResponse struct {
	Filename string
	Err      error
}

type GenerateMapResponse struct {
	GPX domain.GPXFile
	Map domain.MapFile
	Err error
}

type CleanGPXFileResponse struct {
	Content []byte
	GPX     domain.GPXFile
	Err     error
}

type AnnotateMapWithStatsErrorResponse struct {
	Map domain.ShareableMapFile
	Err error
}

type RunningActivity struct {
	Activity domain.RunningActivity
	Deleted  bool
}

type Asset struct {
	Filename string
	Content  []byte
	Deleted  bool
}

type Fake struct {
	t *testing.T

	runs                   []RunningActivity
	cleanedGPXFiles        [][]byte
	assets                 []Asset
	generatedMaps          []domain.GPXFile
	annotatedMapsWithStats []domain.MapFile

	overrideRecordActivityResponse []RunningActivityErrorResponse
	overrideGetActivityResponse    []RunningActivityErrorResponse
	overrideListActivitiesResponse error
	overrideDeleteActivityResponse []RunningActivityErrorResponse
	overrideDeleteAssetResponse    []AssetErrorResponse
	overrideStoreAssetResponse     []AssetErrorResponse
	overrideGenerateMap            []GenerateMapResponse
	overrideCleanGPXFile           []CleanGPXFileResponse
	overrideAnnotateMapWithStats   []AnnotateMapWithStatsErrorResponse

	expectedCleanGPXFiles        [][]byte
	expectedGenerateMap          []domain.GPXFile
	expectedAnnotateMapWithStats []domain.MapFile
	expectedStoreAssets          []string
	expectedRecordActivities     []domain.RunningActivity
	expectedDeletedAssets        []string
	expectedDeletedActivities    []domain.RunningActivitySlug
}

func NewFake(t *testing.T) *Fake {
	return &Fake{t: t}
}

func (f *Fake) CleanGPXFile(ctx context.Context, r io.Reader) (domain.GPXFile, error) {
	content, err := ioutil.ReadAll(r)
	testutils.AssertNoError(f.t, err, "unexpected error while reading gpx file")

	for _, response := range f.overrideCleanGPXFile {
		if string(response.Content) == string(content) {
			if response.Err != nil {
				return domain.GPXFile{}, response.Err
			}

			f.cleanedGPXFiles = append(f.cleanedGPXFiles, content)
			return response.GPX, nil
		}
	}

	f.cleanedGPXFiles = append(f.cleanedGPXFiles, content)

	return domaintest.NewGPXFile(f.t).WithFileContent(content).Build(), nil
}

func (f *Fake) GetRunningActivity(ctx context.Context, slug domain.RunningActivitySlug) (domain.RunningActivity, error) {
	for _, response := range f.overrideGetActivityResponse {
		if response.Slug == slug {
			return domain.RunningActivity{}, response.Err
		}
	}

	for _, activity := range f.runs {
		if activity.Activity.Slug.String() == slug.String() {
			return activity.Activity, nil
		}
	}

	return domain.RunningActivity{}, domain.ErrCantGetRunningSession
}

func (f *Fake) ListRunningActivities(ctx context.Context) ([]domain.RunningActivity, error) {
	if f.overrideListActivitiesResponse != nil {
		return nil, f.overrideListActivitiesResponse
	}

	activities := make([]domain.RunningActivity, 0, len(f.runs))
	for _, activity := range f.runs {
		if activity.Deleted {
			continue
		}

		activities = append(activities, activity.Activity)
	}

	sort.Slice(activities, func(i int, j int) bool {
		return activities[i].RanAt.After(activities[j].RanAt)
	})

	return activities, nil
}

func (f *Fake) DeleteRunningActivity(ctx context.Context, slug domain.RunningActivitySlug) error {
	for _, response := range f.overrideDeleteActivityResponse {
		if slug == response.Slug {
			return response.Err
		}
	}

	for i := range f.runs {
		if f.runs[i].Activity.Slug.String() == slug.String() {
			f.runs[i] = RunningActivity{Activity: f.runs[i].Activity, Deleted: true}

			return nil
		}
	}

	return domain.ErrCantGetRunningSession
}

func (f *Fake) RecordRunningActivity(ctx context.Context, activity domain.RunningActivity) error {
	for _, response := range f.overrideRecordActivityResponse {
		if response.Slug == activity.Slug {
			return response.Err
		}
	}

	for _, recordedActivity := range f.runs {
		if activity.Slug.String() == recordedActivity.Activity.Slug.String() {
			return fmt.Errorf("activity with the same slug already exist")
		}
	}

	f.runs = append(f.runs, RunningActivity{Activity: activity, Deleted: false})

	return nil
}

func (f *Fake) StoreAsset(content io.Reader, filename string) error {
	for _, response := range f.overrideStoreAssetResponse {
		if response.Filename == filename {
			return response.Err
		}
	}

	data, err := io.ReadAll(content)
	if err != nil {
		return err
	}

	for i := range f.assets {
		if filename == f.assets[i].Filename {
			f.assets[i] = Asset{Filename: filename, Content: data, Deleted: false}
			return nil
		}
	}

	f.assets = append(f.assets, Asset{Filename: filename, Content: data, Deleted: false})

	return nil
}

func (f *Fake) DeleteAsset(filename string) error {
	for _, response := range f.overrideDeleteAssetResponse {
		if filename == response.Filename {
			return response.Err
		}
	}

	for i := range f.assets {
		if f.assets[i].Filename == filename {
			f.assets[i] = Asset{
				Filename: f.assets[i].Filename,
				Content:  f.assets[i].Content,
				Deleted:  true,
			}

			return nil
		}
	}

	f.assets = append(f.assets, Asset{Filename: filename, Content: nil, Deleted: true})

	return nil
}

func (f *Fake) AnnotateMapWithStats(ctx context.Context, mapFile domain.MapFile, distance domain.Distance, speed domain.Speed) (domain.ShareableMapFile, error) {
	content, err := ioutil.ReadAll(mapFile.File())
	testutils.AssertNoError(f.t, err, "can't read map content")

	for _, response := range f.overrideAnnotateMapWithStats {
		responseContent, err := ioutil.ReadAll(response.Map.File())
		testutils.AssertNoError(f.t, err, "can't read override map content")
		if string(responseContent) == string(content) {
			return domain.ShareableMapFile{}, response.Err
		}
	}

	f.annotatedMapsWithStats = append(f.annotatedMapsWithStats, mapFile)

	return domain.NewSharableMapFile(content), nil
}

func (f *Fake) GenerateMap(ctx context.Context, gpx domain.GPXFile) (domain.MapFile, error) {
	content, err := ioutil.ReadAll(gpx.File())
	testutils.AssertNoError(f.t, err, "can't read gpx content")

	for _, response := range f.overrideGenerateMap {
		responseContent, err := ioutil.ReadAll(response.GPX.File())
		testutils.AssertNoError(f.t, err, "can't read override gpx content")
		if string(responseContent) == string(content) {
			if response.Err != nil {
				return response.Map, response.Err
			}

			f.generatedMaps = append(f.generatedMaps, gpx)
			return response.Map, nil
		}
	}

	f.generatedMaps = append(f.generatedMaps, gpx)

	return domain.NewMapFile(mapboxtest.GenerateMap()), nil
}

func (f *Fake) ExpectStoreAssets(filenames ...string) {
	f.t.Cleanup(f.VerifyStoreAssets)
	f.expectedStoreAssets = append(f.expectedStoreAssets, filenames...)
}

func (f *Fake) ExpectDeleteAssets(filenames ...string) {
	f.t.Cleanup(f.VerifyDeleteAssets)
	f.expectedDeletedAssets = append(f.expectedDeletedAssets, filenames...)
}

func (f *Fake) VerifyDeleteAssets() {
	for _, filename := range f.expectedDeletedAssets {
		var found bool

		for _, asset := range f.assets {
			if filename != asset.Filename {
				continue
			}

			found = true
			testutils.AssertEqualBool(f.t, true, asset.Deleted, "expecting asset %s to be deleted", filename)
		}

		if !found {
			testutils.AssertEqualBool(f.t, true, false, "expecting asset %s to be deleted but wasn't stored", filename)
		}
	}
}

func (f *Fake) ExpectCleanGPXFiles(contents ...[]byte) {
	f.t.Cleanup(f.VerifyCleanGPXFiles)
	f.expectedCleanGPXFiles = append(f.expectedCleanGPXFiles, contents...)
}

func (f *Fake) ExpectAnnotateMapsWithStats(files ...domain.MapFile) {
	f.t.Cleanup(f.VerifyAnnotateMapsWithStats)
	f.expectedAnnotateMapWithStats = append(f.expectedAnnotateMapWithStats, files...)
}

func (f *Fake) ExpectGenerateMaps(files ...domain.GPXFile) {
	f.t.Cleanup(f.VerifyGenerateMaps)
	f.expectedGenerateMap = append(f.expectedGenerateMap, files...)
}

func (f *Fake) ExpectRecordActivities(activities ...domain.RunningActivity) {
	f.t.Cleanup(f.VerifyRecordActivities)
	f.expectedRecordActivities = append(f.expectedRecordActivities, activities...)
}

func (f *Fake) ExpectDeleteActivities(slugs ...domain.RunningActivitySlug) {
	f.t.Cleanup(f.VerifyDeleteActivities)
	f.expectedDeletedActivities = append(f.expectedDeletedActivities, slugs...)
}

func (f *Fake) VerifyDeleteActivities() {
	for _, slug := range f.expectedDeletedActivities {
		var found bool

		for _, run := range f.runs {
			if slug != run.Activity.Slug {
				continue
			}

			found = true
			testutils.AssertEqualBool(f.t, true, run.Deleted, "expecting run %s to be deleted", slug)
		}

		if !found {
			testutils.AssertEqualBool(f.t, true, false, "expecting run %s to be deleted but wasn't recorded", slug)
		}
	}
}

func (f *Fake) VerifyStoreAssets() {
	for _, filename := range f.expectedStoreAssets {
		var found bool

		for _, asset := range f.assets {
			if filename != asset.Filename {
				continue
			}

			found = true
			testutils.AssertEqualBool(f.t, false, asset.Deleted, "expecting asset %s to be stored but was deleted", filename)
		}

		if !found {
			testutils.AssertEqualBool(f.t, true, false, "expecting asset %s to be stored", filename)
		}
	}
}

func (f *Fake) VerifyRecordActivities() {
	for _, expected := range f.expectedRecordActivities {
		var found bool

		for _, run := range f.runs {
			if expected.Slug != run.Activity.Slug {
				continue
			}

			found = true
			testutils.AssertEqualBool(f.t, false, run.Deleted, "expecting run %s to be recorded but was deleted", run.Activity.Slug)
			domaintest.AssertEqualRunningActivity(f.t, expected, run.Activity, "invalid recorded activity")
		}

		if !found {
			testutils.AssertEqualBool(f.t, true, false, "expecting run %s to be recorded", expected.Slug)
		}
	}
}

func (f *Fake) OverrideCleanGPXFile(content []byte, gpx domain.GPXFile, err error) {
	f.overrideCleanGPXFile = append(f.overrideCleanGPXFile, CleanGPXFileResponse{
		Content: content,
		GPX:     gpx,
		Err:     err,
	})
}

func (f *Fake) OverrideDeleteActivity(slug domain.RunningActivitySlug, err error) {
	f.overrideDeleteActivityResponse = append(f.overrideDeleteActivityResponse, RunningActivityErrorResponse{
		Slug: slug,
		Err:  err,
	})
}

func (f *Fake) OverrideRecordActivity(slug domain.RunningActivitySlug, err error) {
	f.overrideRecordActivityResponse = append(f.overrideRecordActivityResponse, RunningActivityErrorResponse{
		Slug: slug,
		Err:  err,
	})
}

func (f *Fake) OverrideGetActivity(slug domain.RunningActivitySlug, err error) {
	f.overrideGetActivityResponse = append(f.overrideGetActivityResponse, RunningActivityErrorResponse{
		Slug: slug,
		Err:  err,
	})
}

func (f *Fake) OverrideListActivities(err error) {
	f.overrideListActivitiesResponse = err
}

func (f *Fake) OverrideDeleteAsset(filename string, err error) {
	f.overrideDeleteAssetResponse = append(f.overrideDeleteAssetResponse, AssetErrorResponse{
		Filename: filename,
		Err:      err,
	})
}

func (f *Fake) OverrideStoreAsset(filename string, err error) {
	f.overrideStoreAssetResponse = append(f.overrideStoreAssetResponse, AssetErrorResponse{
		Filename: filename,
		Err:      err,
	})
}

func (f *Fake) OverrideGenerateMap(gpx domain.GPXFile, mapFile domain.MapFile, err error) {
	f.overrideGenerateMap = append(f.overrideGenerateMap, GenerateMapResponse{
		GPX: gpx,
		Map: mapFile,
		Err: err,
	})
}

func (f *Fake) VerifyGenerateMaps() {
	for _, expected := range f.expectedGenerateMap {
		var found bool

		expectedContent, err := ioutil.ReadAll(expected.File())
		testutils.AssertNoError(f.t, err, "unexpected error while reading expected gpx file")

		for _, actual := range f.generatedMaps {
			actualContent, err := ioutil.ReadAll(actual.File())
			testutils.AssertNoError(f.t, err, "unexpected error while reading actual gpx file")

			if string(expectedContent) != string(actualContent) {
				continue
			}

			found = true
			break
		}

		if !found {
			testutils.AssertEqualBool(f.t, true, false, "expecting map to have been generated with gpx %#+v", string(expectedContent))
		}
	}
}

func (f *Fake) VerifyCleanGPXFiles() {
	for _, expected := range f.expectedCleanGPXFiles {
		var found bool

		for _, actual := range f.cleanedGPXFiles {
			if string(actual) != string(expected) {
				continue
			}

			found = true
			break
		}

		if !found {
			testutils.AssertEqualBool(f.t, true, false, "expecting clean gpx file to have been generated with content:\n%#+v", string(expected))
		}
	}
}

func (f *Fake) VerifyAnnotateMapsWithStats() {
	for _, expected := range f.expectedAnnotateMapWithStats {
		var found bool

		expectedContent, err := ioutil.ReadAll(expected.File())
		testutils.AssertNoError(f.t, err, "unexpected error while reading expected map file")

		for _, actual := range f.annotatedMapsWithStats {
			actualContent, err := ioutil.ReadAll(actual.File())
			testutils.AssertNoError(f.t, err, "unexpected error while reading actual map file")

			if string(expectedContent) != string(actualContent) {
				continue
			}

			found = true
			break
		}

		if !found {
			testutils.AssertEqualBool(f.t, true, false, "expecting annotated map to have been generated from map %#+v", string(expectedContent))
		}
	}
}
