package job_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/application/applicationtest"
	"github.com/lonepeon/sport/internal/domain/domaintest"
	"github.com/lonepeon/sport/internal/infrastructure/job"
)

func TestHandleInvalidPayload(t *testing.T) {
	err := job.NewDeleteRunningSessionJob(nil).
		Handle(context.Background(), []byte(`{this is not a json}`))

	testutils.AssertErrorContains(t, "can't parse input", err, "unexpected error")
}

func TestHandleInvalidSlug(t *testing.T) {
	err := job.NewDeleteRunningSessionJob(nil).
		Handle(context.Background(), []byte(`{"slug": "invalid slug"}`))

	testutils.AssertErrorContains(t, "can't parse slug", err, "unexpected error")
}

func TestHandleCannotDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	application := applicationtest.NewMockApplication(ctrl)

	application.EXPECT().
		DeleteRunningSession(gomock.Any(), domaintest.MatchRunningActivitySlug("202202231558")).
		Return(errors.New("boom"))

	err := job.NewDeleteRunningSessionJob(application).
		Handle(context.Background(), []byte(`{"slug": "202202231558"}`))

	testutils.AssertErrorContains(t, "can't delete", err, "unexpected error")
}

func TestHandleSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	application := applicationtest.NewMockApplication(ctrl)

	application.EXPECT().
		DeleteRunningSession(gomock.Any(), domaintest.MatchRunningActivitySlug("202202231558")).
		Return(nil)

	err := job.NewDeleteRunningSessionJob(application).
		Handle(context.Background(), []byte(`{"slug": "202202231558"}`))

	testutils.AssertNoError(t, err, "unexpected error")
}
