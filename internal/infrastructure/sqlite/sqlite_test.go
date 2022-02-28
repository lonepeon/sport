package sqlite_test

import (
	"context"
	"database/sql"
	"io/ioutil"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3" // sqlite3 adapter

	"github.com/lonepeon/golib/sqlutil"
	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/domain"
	"github.com/lonepeon/sport/internal/domain/domaintest"
	"github.com/lonepeon/sport/internal/infrastructure/sqlite"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Parallel()

	t.Run("ListRunningActivities", testListRunningActivities)
	t.Run("GetRunningActivitySuccess", testGetRunningActivitySuccess)
	t.Run("GetRunningActivityNotFound", testGetRunningActivityNotFound)
	t.Run("DeleteRunningActivitySuccess", testDeleteRunningActivitySuccess)
	t.Run("DeleteRunningActivityWhenActivityDoesNotMatch", testDeleteRunningActivityWhenActivityDoesNotMatch)
}

func testGetRunningActivitySuccess(t *testing.T) {
	repo, cleanup := setupDatabase(t)
	defer cleanup()
	expectedActivity := domaintest.NewRunningActivity(t).Build()

	recordActivity(t, repo, expectedActivity)

	actualActivity, err := repo.GetRunningActivity(context.Background(), expectedActivity.Slug)

	testutils.AssertNoError(t, err, "can't get activity")
	domaintest.AssertEqualRunningActivity(t, expectedActivity, actualActivity, "unexpected activity")
}

func testGetRunningActivityNotFound(t *testing.T) {
	repo, cleanup := setupDatabase(t)
	defer cleanup()
	activity := domaintest.NewRunningActivity(t).Build()
	recordActivity(t, repo, activity)

	slug, err := domain.NewRunnningActivitySlugFromString("202101010000")
	testutils.AssertNoError(t, err, "can't build slug")

	_, err = repo.GetRunningActivity(context.Background(), slug)

	testutils.AssertErrorIs(t, domain.ErrCantGetRunningSession, err, "can't get activity")
}

func testListRunningActivities(t *testing.T) {
	repo, cleanup := setupDatabase(t)
	defer cleanup()

	activity1 := domaintest.NewRunningActivity(t).WithRawSlug("202101010000").Build()
	activity2 := domaintest.NewRunningActivity(t).WithRawSlug("202303030000").Build()
	activity3 := domaintest.NewRunningActivity(t).WithRawSlug("202202020000").Build()

	recordActivity(t, repo, activity1)
	recordActivity(t, repo, activity2)
	recordActivity(t, repo, activity3)

	activities, err := repo.ListRunningActivities(context.Background())

	testutils.AssertNoError(t, err, "can't list activities")
	testutils.AssertEqualInt(t, 3, len(activities), "unexpected number of activities")

	domaintest.AssertEqualRunningActivity(t, activity2, activities[0], "unexpected activity")
	domaintest.AssertEqualRunningActivity(t, activity3, activities[1], "unexpected activity")
	domaintest.AssertEqualRunningActivity(t, activity1, activities[2], "unexpected activity")
}

func testDeleteRunningActivitySuccess(t *testing.T) {
	repo, cleanup := setupDatabase(t)
	defer cleanup()
	expectedActivity := domaintest.NewRunningActivity(t).Build()

	recordActivity(t, repo, expectedActivity)

	err := repo.DeleteRunningActivity(context.Background(), expectedActivity.Slug)
	testutils.AssertNoError(t, err, "can't delete activity")

	_, err = repo.GetRunningActivity(context.Background(), expectedActivity.Slug)
	testutils.AssertErrorIs(t, domain.ErrCantGetRunningSession, err, "activity should have been deleted")
}

func testDeleteRunningActivityWhenActivityDoesNotMatch(t *testing.T) {
	repo, cleanup := setupDatabase(t)
	defer cleanup()

	slug, err := domain.NewRunnningActivitySlugFromString("202101010000")
	testutils.AssertNoError(t, err, "can't build slug")

	expectedActivity := domaintest.NewRunningActivity(t).WithRawSlug("202303030000").Build()

	recordActivity(t, repo, expectedActivity)

	err = repo.DeleteRunningActivity(context.Background(), slug)
	testutils.AssertErrorIs(t, domain.ErrCantGetRunningSession, err, "activty should have been not found")
}

func setupDatabase(t *testing.T) (sqlite.SQLite, func()) {
	file, err := ioutil.TempFile("/tmp", "sqlite.XXXX")
	testutils.AssertNoError(t, err, "can't create sqlite temp file")

	db, err := sql.Open("sqlite3", file.Name())
	testutils.AssertNoError(t, err, "can't open sqlite connection")

	_, err = sqlutil.ExecuteMigrations(context.Background(), db, sqlite.Migrations())
	testutils.AssertNoError(t, err, "can't run migrations")

	return sqlite.New(db), func() {
		file.Close()
		os.Remove(file.Name())
		db.Close()
	}
}

func recordActivity(t *testing.T, repo sqlite.SQLite, activity domain.RunningActivity) {
	err := repo.RecordRunningActivity(context.Background(), activity)
	testutils.AssertNoError(t, err, "can't record activity")
}
