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

func TestRunningSessionIndexError(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := webtest.NewMockContext(ctrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	usecase := applicationtest.NewMockApplication(ctrl)

	usecase.EXPECT().ListRunningSessions(gomock.Any()).Return(nil, errors.New("boom"))

	expected := webtest.MockedResponse("server error")
	ctx.EXPECT().StdCtx().AnyTimes()
	ctx.EXPECT().
		InternalServerErrorResponse(
			gomockutils.ContainsString("can't list"),
			gomock.Any(),
		).
		Return(expected)

	actual := www.RunningSessionsIndex(usecase)(ctx, w, r)

	webtest.AssertResponse(t, expected, actual, "unexpected response")
}

func TestRunningSessionIndexSuccess(t *testing.T) {
	tcs := map[string]struct {
		activities []domain.RunningActivity
	}{
		"noSessions": {},
		"withSessions": {
			activities: []domain.RunningActivity{
				domaintest.NewRunningActivity(t).Build(),
				domaintest.NewRunningActivity(t).Build(),
				domaintest.NewRunningActivity(t).Build(),
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := webtest.NewMockContext(ctrl)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			usecase := applicationtest.NewMockApplication(ctrl)

			usecase.EXPECT().ListRunningSessions(gomock.Any()).Return(tc.activities, nil)

			expected := webtest.MockedResponse("ok response")
			ctx.EXPECT().StdCtx().AnyTimes()
			ctx.EXPECT().
				Response(
					200,
					gomock.Any(),
					webtest.MatchDataContains("Activities", tc.activities),
				).
				Return(expected)

			actual := www.RunningSessionsIndex(usecase)(ctx, w, r)

			webtest.AssertResponse(t, expected, actual, "unexpected response")
		})
	}
}
