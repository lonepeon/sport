package www_test

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/golib/testutils/gomockutils"
	"github.com/lonepeon/golib/web"
	"github.com/lonepeon/golib/web/webtest"
	"github.com/lonepeon/sport/internal/application/applicationtest"
	"github.com/lonepeon/sport/internal/domain"
	"github.com/lonepeon/sport/internal/domain/domaintest"
	"github.com/lonepeon/sport/internal/infrastructure/job"
	"github.com/lonepeon/sport/internal/infrastructure/job/jobtest"
	"github.com/lonepeon/sport/internal/infrastructure/www"
)

func TestRunningSessionDeleteInvalidDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := webtest.NewMockContext(ctrl)
	response := httptest.NewRecorder()
	request := httptest.NewRequest("DELETE", "/running-session/{slug}", nil)

	expectedResponse := webtest.MockedResponse("redirection")
	ctx.EXPECT().Vars(request).Return(map[string]string{"slug": "wrong-date"})
	ctx.EXPECT().AddFlash(web.NewFlashMessageError("no activity recorded with slug 'wrong-date'"))
	ctx.EXPECT().Redirect(response, 303, "/").Return(expectedResponse)

	actualResponse := www.RunningSessionsDelete(nil, nil)(ctx, response, request)

	webtest.AssertResponse(t, expectedResponse, actualResponse, "invalid response")
	testutils.AssertContainsString(t, "can't parse activity time", actualResponse.LogMessage, "unexpected log message")
	testutils.AssertContainsString(t, "wrong-date", actualResponse.LogMessage, "unexpected log message")
}

func TestRunningSessionDeleteRunningSessionDoesNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	application := applicationtest.NewMockApplication(ctrl)
	ctx := webtest.NewMockContext(ctrl)
	response := httptest.NewRecorder()
	request := httptest.NewRequest("DELETE", "/running-session/{slug}", nil)
	slug, err := domain.NewRunnningActivitySlugFromString("202101101105")
	testutils.AssertNoError(t, err, "can't parse slug")

	expectedResponse := webtest.MockedResponse("redirection")
	ctx.EXPECT().Vars(request).Return(map[string]string{"slug": "202101101105"})
	ctx.EXPECT().StdCtx()
	application.EXPECT().
		GetRunningSession(gomock.Any(), slug).
		Return(domain.RunningActivity{}, domain.ErrCantGetRunningSession)
	ctx.EXPECT().AddFlash(web.NewFlashMessageError("no activity recorded with slug '202101101105'"))
	ctx.EXPECT().Redirect(response, 303, "/").Return(expectedResponse)

	actualResponse := www.RunningSessionsDelete(application, nil)(ctx, response, request)

	webtest.AssertResponse(t, expectedResponse, actualResponse, "invalid response")
	testutils.AssertContainsString(t, "can't find activity", actualResponse.LogMessage, "unexpected log message")
	testutils.AssertContainsString(t, "202101101105", actualResponse.LogMessage, "unexpected log message")
}

func TestRunningSessionDeleteCannotEnqueueJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	application := applicationtest.NewMockApplication(ctrl)
	enqueuer := jobtest.NewMockEnqueuer(ctrl)
	ctx := webtest.NewMockContext(ctrl)
	response := httptest.NewRecorder()
	request := httptest.NewRequest("DELETE", "/running-session/{slug}", nil)
	activity := domaintest.NewRunningActivity(t).Build()
	expectedJob := jobtest.NewJobMatcher(
		"delete-running-session-job",
		&job.DeleteRunningSessionJobInput{},
		func(arg interface{}) bool {
			input := arg.(*job.DeleteRunningSessionJobInput)

			return input.Slug == activity.Slug.String()
		})

	expectedResponse := webtest.MockedResponse("redirection")
	ctx.EXPECT().Vars(request).Return(map[string]string{"slug": activity.Slug.String()})
	ctx.EXPECT().StdCtx()
	application.EXPECT().
		GetRunningSession(gomock.Any(), gomockutils.Equal(activity.Slug)).
		Return(activity, nil)
	enqueuer.EXPECT().Enqueue(expectedJob).Return(fmt.Errorf("boom"))
	ctx.EXPECT().InternalServerErrorResponse(gomock.Any(), gomock.Any()).Return(expectedResponse)

	actualResponse := www.RunningSessionsDelete(application, enqueuer)(ctx, response, request)

	webtest.AssertResponse(t, expectedResponse, actualResponse, "invalid response")
}

func TestRunningSessionDeleteSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	application := applicationtest.NewMockApplication(ctrl)
	enqueuer := jobtest.NewMockEnqueuer(ctrl)
	ctx := webtest.NewMockContext(ctrl)
	response := httptest.NewRecorder()
	request := httptest.NewRequest("DELETE", "/running-session/{slug}", nil)
	activity := domaintest.NewRunningActivity(t).Build()
	expectedJob := jobtest.NewJobMatcher(
		"delete-running-session-job",
		&job.DeleteRunningSessionJobInput{},
		func(arg interface{}) bool {
			input := arg.(*job.DeleteRunningSessionJobInput)

			return input.Slug == activity.Slug.String()
		})

	expectedResponse := webtest.MockedResponse("redirection")
	ctx.EXPECT().Vars(request).Return(map[string]string{"slug": activity.Slug.String()})
	ctx.EXPECT().StdCtx()
	application.EXPECT().
		GetRunningSession(gomock.Any(), gomockutils.Equal(activity.Slug)).
		Return(activity, nil)
	enqueuer.EXPECT().Enqueue(expectedJob).Return(nil)
	ctx.EXPECT().AddFlash(web.NewFlashMessageSuccess("activity recorded with slug '%s' is being deleted", activity.Slug))
	ctx.EXPECT().Redirect(response, 303, "/").Return(expectedResponse)

	actualResponse := www.RunningSessionsDelete(application, enqueuer)(ctx, response, request)

	webtest.AssertResponse(t, expectedResponse, actualResponse, "invalid response")
}
