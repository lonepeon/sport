package domaintest

import (
	"fmt"

	"github.com/lonepeon/sport/internal/domain"
)

func MatchRunningActivitySlug(s string) _GoMockRunningActivitySlug {
	return _GoMockRunningActivitySlug{slugStr: s}
}

type _GoMockRunningActivitySlug struct {
	slugStr string
}

func (m _GoMockRunningActivitySlug) Matches(v interface{}) bool {
	slug, ok := v.(domain.RunningActivitySlug)
	if !ok {
		return false
	}

	return m.slugStr == slug.String()
}

func (m _GoMockRunningActivitySlug) String() string {
	return fmt.Sprintf("slug=%s", m.slugStr)
}
