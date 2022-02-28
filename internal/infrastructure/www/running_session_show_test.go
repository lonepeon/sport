package www_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lonepeon/golib/testutils/gomockutils"
	"github.com/lonepeon/golib/web/webtest"
	"github.com/lonepeon/sport/internal/application/applicationtest"
	"github.com/lonepeon/sport/internal/domain"
	"github.com/lonepeon/sport/internal/domain/domaintest"
	"github.com/lonepeon/sport/internal/infrastructure/www"
)

func TestRunningSessionShowInvalidSlug(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := webtest.NewMockContext(ctrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/running-session/{slug}", nil)

	ctx.EXPECT().Vars(r).Return(map[string]string{"slug": "invalid slug"})
	expected := webtest.MockedResponse("not found")
	ctx.EXPECT().
		NotFoundResponse(gomockutils.ContainsString("can't parse"), gomock.Any(), gomock.Any()).
		Return(expected)

	actual := www.RunningSessionsShow(nil)(ctx, w, r)

	webtest.AssertResponse(t, expected, actual, "unexpected response")
}

func TestRunningSessionShowActivityNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	application := applicationtest.NewMockApplication(ctrl)
	ctx := webtest.NewMockContext(ctrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/running-session/{slug}", nil)

	ctx.EXPECT().StdCtx().AnyTimes()
	ctx.EXPECT().Vars(r).Return(map[string]string{"slug": "202102122146"})
	application.EXPECT().
		GetRunningSession(gomock.Any(), domaintest.MatchRunningActivitySlug("202102122146")).
		Return(domain.RunningActivity{}, domain.ErrCantGetRunningSession)
	expected := webtest.MockedResponse("not found")
	ctx.EXPECT().
		NotFoundResponse(gomockutils.ContainsString("can't find"), gomock.Any(), gomock.Any()).
		Return(expected)

	actual := www.RunningSessionsShow(application)(ctx, w, r)

	webtest.AssertResponse(t, expected, actual, "unexpected response")
}

func TestRunningSessionShowActivityUnexpectedError(t *testing.T) {
	ctrl := gomock.NewController(t)
	application := applicationtest.NewMockApplication(ctrl)
	ctx := webtest.NewMockContext(ctrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/running-session/{slug}", nil)

	ctx.EXPECT().StdCtx().AnyTimes()
	ctx.EXPECT().Vars(r).Return(map[string]string{"slug": "202102122146"})
	application.EXPECT().
		GetRunningSession(gomock.Any(), domaintest.MatchRunningActivitySlug("202102122146")).
		Return(domain.RunningActivity{}, errors.New("boom"))
	expected := webtest.MockedResponse("server error")
	ctx.EXPECT().
		InternalServerErrorResponse(gomockutils.ContainsString("failed"), gomock.Any(), gomock.Any()).
		Return(expected)

	actual := www.RunningSessionsShow(application)(ctx, w, r)

	webtest.AssertResponse(t, expected, actual, "unexpected response")
}

func TestRunningSessionShowActivitySuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	application := applicationtest.NewMockApplication(ctrl)
	ctx := webtest.NewMockContext(ctrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/running-session/{slug}", nil)
	rawSlug := "202102122146"
	activity := domaintest.NewRunningActivity(t).WithRawSlug(rawSlug).Build()

	ctx.EXPECT().StdCtx().AnyTimes()
	ctx.EXPECT().Vars(r).Return(map[string]string{"slug": rawSlug})
	application.EXPECT().
		GetRunningSession(gomock.Any(), domaintest.MatchRunningActivitySlug(rawSlug)).
		Return(activity, nil)
	expected := webtest.MockedResponse("ok")
	ctx.EXPECT().
		Response(200, gomock.Any(), webtest.MatchDataContains("Activity", activity)).
		Return(expected)

	actual := www.RunningSessionsShow(application)(ctx, w, r)

	webtest.AssertResponse(t, expected, actual, "unexpected response")
}
