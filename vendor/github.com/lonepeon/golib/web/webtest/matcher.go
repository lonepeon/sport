package webtest

import (
	"fmt"
	"net/http"

	"github.com/golang/mock/gomock"
)

type HasRequestVerbMatcher struct {
	expectedVerb string
}

func (m HasRequestVerbMatcher) Matches(arg interface{}) bool {
	req := arg.(*http.Request)

	return m.expectedVerb == req.Method
}

func (m HasRequestVerbMatcher) String() string {
	return fmt.Sprintf("expected verb %s", m.expectedVerb)
}

func HasRequestVerb(expectedVerb string) gomock.Matcher {
	return HasRequestVerbMatcher{expectedVerb: expectedVerb}
}
