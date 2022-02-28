package service_test

import (
	"context"
	"testing"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/application/service"
	"github.com/lonepeon/sport/internal/domain"
	"github.com/lonepeon/sport/internal/domain/domaintest"
	"github.com/lonepeon/sport/internal/repository/repositorytest"
)

func TestGetRunningSessionSuccess(t *testing.T) {
	repo := repositorytest.NewFake(t)
	expectedActivity := domaintest.NewRunningActivity(t).Persist(repo)

	actualActivity, err := service.GetRunningSession(repo, context.Background(), expectedActivity.Slug)

	testutils.AssertNoError(t, err, "can't get running session")
	domaintest.AssertEqualRunningActivity(t, expectedActivity, actualActivity, "unexpected activity")
}

func TestGetRunningSessionNotFound(t *testing.T) {
	repo := repositorytest.NewFake(t)
	slug, err := domain.NewRunnningActivitySlugFromString("202202162149")
	testutils.AssertNoError(t, err, "can't build slug")

	_, err = service.GetRunningSession(repo, context.Background(), slug)

	testutils.AssertErrorIs(t, domain.ErrCantGetRunningSession, err, "unexpected error")
}
