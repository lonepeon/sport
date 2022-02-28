package service_test

import (
	"context"
	"testing"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/application/service"
	"github.com/lonepeon/sport/internal/domain/domaintest"
	"github.com/lonepeon/sport/internal/repository/repositorytest"
)

func TestListRunningSessionsSuccess(t *testing.T) {
	repo := repositorytest.NewFake(t)
	activity1 := domaintest.NewRunningActivity(t).WithRawSlug("202101010000").Persist(repo)
	activity2 := domaintest.NewRunningActivity(t).WithRawSlug("202303030000").Persist(repo)
	activity3 := domaintest.NewRunningActivity(t).WithRawSlug("202202020000").Persist(repo)

	actualActivities, err := service.ListRunningSessions(repo, context.Background())

	testutils.AssertNoError(t, err, "can't get running sessions")
	testutils.AssertEqualInt(t, 3, len(actualActivities), "unexpected number of activities")

	domaintest.AssertEqualRunningActivity(t, activity2, actualActivities[0], "unexpected activity")
	domaintest.AssertEqualRunningActivity(t, activity3, actualActivities[1], "unexpected activity")
	domaintest.AssertEqualRunningActivity(t, activity1, actualActivities[2], "unexpected activity")
}

func TestListRunningSessionsNoEntries(t *testing.T) {
	repo := repositorytest.NewFake(t)

	actualActivities, err := service.ListRunningSessions(repo, context.Background())

	testutils.AssertNoError(t, err, "can't get running sessions")
	testutils.AssertEqualInt(t, 0, len(actualActivities), "unexpected number of activities")
}
