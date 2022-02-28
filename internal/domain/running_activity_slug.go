package domain

import (
	"fmt"
	"time"
)

var runningActivitySlugFormat = "200601021504"

type RunningActivitySlug struct {
	slug Slug
	time time.Time
}

func (r RunningActivitySlug) Equal(v interface{}) bool {
	other, ok := v.(RunningActivitySlug)
	if !ok {
		return false
	}

	return r.String() == other.String()
}

func NewRunnningActivitySlugFromString(s string) (RunningActivitySlug, error) {
	t, err := time.Parse(runningActivitySlugFormat, s)
	if err != nil {
		return RunningActivitySlug{}, fmt.Errorf("can't parse date format: %v", err)
	}

	return NewRunnningActivitySlugFromTime(t)
}

func NewRunnningActivitySlugFromTime(t time.Time) (RunningActivitySlug, error) {
	slug, err := NewSlug(t.Format(runningActivitySlugFormat))
	if err != nil {
		return RunningActivitySlug{}, fmt.Errorf("invalid date format: %v", err)
	}

	return RunningActivitySlug{slug: slug, time: t}, nil
}

func (s RunningActivitySlug) String() string {
	return s.slug.String()
}

func (s RunningActivitySlug) Time() time.Time {
	return s.time
}
