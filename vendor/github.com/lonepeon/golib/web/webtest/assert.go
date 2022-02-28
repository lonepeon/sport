package webtest

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/golib/web"
)

func AssertResponse(t *testing.T, want web.Response, got web.Response, format string, args ...interface{}) {
	t.Helper()

	explanation := func(msg string) string {
		return fmt.Sprintf("%s: %s", fmt.Sprintf(format, args...), msg)
	}

	testutils.AssertEqualString(t, want.Layout, got.Layout, explanation("unexpected response layout"))
	testutils.AssertEqualString(t, want.Template, got.Template, explanation("unexpected response template"))
	testutils.AssertEqualInt(t, want.HTTPCode, got.HTTPCode, explanation("unexpected response http code"))

	if !reflect.DeepEqual(want.Data, got.Data) {
		t.Errorf("%s\nwant:/n%#+v\ngot:\n%#+v\n", explanation("unexpected response data"), want.Data, got.Data)
	}
}
