package domain_test

import (
	"testing"
	"time"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/domain"
)

func TestRunningActivitySlugFromStringSuccess(t *testing.T) {
	slugs := []string{
		"202202231558",
	}

	for _, s := range slugs {
		t.Run(s, func(t *testing.T) {
			slug, err := domain.NewRunnningActivitySlugFromString(s)
			testutils.AssertNoError(t, err, "didn't expect error when parsing slug %s", s)
			testutils.AssertEqualString(t, s, slug.String(), "unexpected slug")
		})
	}
}

func TestRunningActivitySlugFromStringFailure(t *testing.T) {
	slugs := []string{
		"",
		"2021-02-11 21:55:00",
		"2021 Feb 11 21:55:00",
		"not-a-date",
	}

	for _, s := range slugs {
		t.Run(s, func(t *testing.T) {
			_, err := domain.NewRunnningActivitySlugFromString(s)
			testutils.AssertHasError(t, err, "expecting error when parsing slug '%s'", s)
		})
	}
}

func TestRunningActivitySlugFromTimeSuccess(t *testing.T) {
	d, err := time.Parse(time.RFC3339, "2022-03-19T23:42:16Z")
	testutils.AssertNoError(t, err, "can't parse date")

	slug, err := domain.NewRunnningActivitySlugFromTime(d)
	testutils.AssertNoError(t, err, "didn't expect error when parsing time %v", t)
	testutils.AssertEqualString(t, "202203192342", slug.String(), "unexpected slug")
}
