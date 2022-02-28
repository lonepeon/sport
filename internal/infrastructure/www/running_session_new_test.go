package www_test

import (
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lonepeon/golib/web/webtest"
	"github.com/lonepeon/sport/internal/infrastructure/www"
)

func TestRunningSessionNewSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := webtest.NewMockContext(ctrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/running-session/new", nil)

	expected := webtest.MockedResponse("ok response")
	ctx.EXPECT().Response(200, gomock.Any(), nil).Return(expected)

	actual := www.RunningSessionNew()(ctx, w, r)

	webtest.AssertResponse(t, expected, actual, "unexpected response")
}
