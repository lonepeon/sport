package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/application/service"
	"github.com/lonepeon/sport/internal/domain"
	"github.com/lonepeon/sport/internal/domain/domaintest"
	"github.com/lonepeon/sport/internal/repository/repositorytest"
)

func TestDeleteRunningSessionActivityNotFound(t *testing.T) {
	repo := repositorytest.NewFake(t)
	activity := domaintest.NewRunningActivity(t).Build()

	err := service.DeleteRunningSession(repo, context.Background(), activity.Slug)

	testutils.AssertErrorIs(t, domain.ErrCantGetRunningSession, err, "unexpected running session result")
}

func TestDeleteRunningSessionActivityCantDeleteGPX(t *testing.T) {
	repo := repositorytest.NewFake(t)
	activity := domaintest.NewRunningActivity(t).Persist(repo)

	repo.OverrideDeleteAsset(activity.GPXPath.String(), errors.New("boom"))

	err := service.DeleteRunningSession(repo, context.Background(), activity.Slug)

	testutils.AssertHasError(t, err, "unexpected running session result")
	testutils.AssertContainsString(t, "boom", err.Error(), "unexpected error message")
	testutils.AssertContainsString(t, activity.GPXPath.String(), err.Error(), "unexpected error message")
}

func TestDeleteRunningSessionActivityCantDeleteMap(t *testing.T) {
	repo := repositorytest.NewFake(t)
	activity := domaintest.NewRunningActivity(t).Persist(repo)

	repo.ExpectDeleteAssets(activity.GPXPath.String())
	repo.OverrideDeleteAsset(activity.MapPath.String(), errors.New("boom"))

	err := service.DeleteRunningSession(repo, context.Background(), activity.Slug)

	testutils.AssertHasError(t, err, "unexpected running session result")
	testutils.AssertContainsString(t, "boom", err.Error(), "unexpected error message")
	testutils.AssertContainsString(t, activity.MapPath.String(), err.Error(), "unexpected error message")
}

func TestDeleteRunningSessionActivityCantDeleteShareableMap(t *testing.T) {
	repo := repositorytest.NewFake(t)
	activity := domaintest.NewRunningActivity(t).Persist(repo)

	repo.ExpectDeleteAssets(activity.GPXPath.String(), activity.MapPath.String())
	repo.OverrideDeleteAsset(activity.ShareableMapPath.String(), errors.New("boom"))

	err := service.DeleteRunningSession(repo, context.Background(), activity.Slug)

	testutils.AssertHasError(t, err, "unexpected running session result")
	testutils.AssertContainsString(t, "boom", err.Error(), "unexpected error message")
	testutils.AssertContainsString(t, activity.ShareableMapPath.String(), err.Error(), "unexpected error message")
}

func TestDeleteRunningSessionActivityCantDeleteActivityBecauseDoesNotExist(t *testing.T) {
	repo := repositorytest.NewFake(t)
	activity := domaintest.NewRunningActivity(t).Persist(repo)

	repo.ExpectDeleteAssets(
		activity.GPXPath.String(),
		activity.MapPath.String(),
		activity.ShareableMapPath.String(),
	)

	repo.OverrideDeleteActivity(activity.Slug, domain.ErrCantGetRunningSession)

	err := service.DeleteRunningSession(repo, context.Background(), activity.Slug)

	testutils.AssertErrorIs(t, domain.ErrCantGetRunningSession, err, "unexpected running session result")
}

func TestDeleteRunningSessionActivityCantDeleteActivityBecauseUnexpectedErrorHappened(t *testing.T) {
	repo := repositorytest.NewFake(t)
	activity := domaintest.NewRunningActivity(t).Persist(repo)

	repo.ExpectDeleteAssets(
		activity.GPXPath.String(),
		activity.MapPath.String(),
		activity.ShareableMapPath.String(),
	)

	repo.OverrideDeleteActivity(activity.Slug, errors.New("boom"))

	err := service.DeleteRunningSession(repo, context.Background(), activity.Slug)

	testutils.AssertHasError(t, err, "unexpected running session result")
	testutils.AssertContainsString(t, "boom", err.Error(), "unexpected error content")
}

func TestDeleteRunningSessionActivitySuccess(t *testing.T) {
	repo := repositorytest.NewFake(t)
	activity := domaintest.NewRunningActivity(t).Persist(repo)

	repo.ExpectDeleteActivities(activity.Slug)
	repo.ExpectDeleteAssets(
		activity.GPXPath.String(),
		activity.MapPath.String(),
		activity.ShareableMapPath.String(),
	)

	err := service.DeleteRunningSession(repo, context.Background(), activity.Slug)
	testutils.AssertNoError(t, err, "unexpected running session result")
}
