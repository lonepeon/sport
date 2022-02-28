package gomockutils

import (
	"fmt"

	"github.com/golang/mock/gomock"
)

type Equaler interface {
	Equal(interface{}) bool
}

func Equal(v Equaler) gomock.Matcher {
	return _GoMockEqMatcher{value: v}
}

type _GoMockEqMatcher struct {
	value Equaler
}

func (m _GoMockEqMatcher) Matches(v interface{}) bool {
	return m.value.Equal(v)
}

func (m _GoMockEqMatcher) String() string {
	return fmt.Sprintf("value=%#+v, type=%T", m.value, m.value)
}
