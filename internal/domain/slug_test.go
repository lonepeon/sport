package domain_test

import (
	"testing"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/domain"
)

func TestSlugNewSlugSuccess(t *testing.T) {
	slugs := []string{
		"slug",
		"a-slug",
		"a-nice-slug-with-number-1234",
		"1234.another-slug.with.number-1234",
	}

	for _, s := range slugs {
		t.Run(s, func(t *testing.T) {
			slug, err := domain.NewSlug(s)
			testutils.AssertNoError(t, err, "didn't expect error when parsing slug %s", s)
			testutils.AssertEqualString(t, s, slug.String(), "unexpected slug")
		})
	}
}

func TestSlugNewSlugFailure(t *testing.T) {
	slugs := []string{
		"",
		"a-slug-in-error-because-it-is-greater-than-the-allowed-size-limit",
		"a spaced slug",
		"a-UPPERCASE-slug",
		"a-double--dash",
		"a-double..dot",
		"-start-with-dash",
		".start-with-dot",
		"end-with-dash-",
		"end-with-dot.",
	}

	for _, s := range slugs {
		t.Run(s, func(t *testing.T) {
			_, err := domain.NewSlug(s)
			testutils.AssertHasError(t, err, "expecting error when parsing slug %s", s)
		})
	}
}
