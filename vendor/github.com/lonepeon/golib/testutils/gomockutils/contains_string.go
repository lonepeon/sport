package gomockutils

import (
	"fmt"
	"strings"
)

func ContainsString(s string) _GoMockStringContainsMatcher {
	return _GoMockStringContainsMatcher{value: s}
}

type _GoMockStringContainsMatcher struct {
	value string
}

func (m _GoMockStringContainsMatcher) Matches(v interface{}) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	return strings.Contains(s, m.value)
}

func (m _GoMockStringContainsMatcher) String() string {
	return fmt.Sprintf("contains=%s", m.value)
}
