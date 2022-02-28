package domaintest

import (
	"testing"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/domain"
)

func AssertEqualRunningActivity(t *testing.T, want domain.RunningActivity, got domain.RunningActivity, format string, args ...interface{}) {
	t.Helper()

	testutils.AssertEqualTime(t, want.RanAt, got.RanAt, format, args...)
	testutils.AssertEqualDuration(t, want.Duration, got.Duration, format, args...)
	testutils.AssertEqualInt(t, want.Distance.Meters(), got.Distance.Meters(), format, args...)
	testutils.AssertEqualFloat64(t, want.Speed.KilometersPerHour(), got.Speed.KilometersPerHour(), format, args...)
	testutils.AssertEqualString(t, want.GPXPath.String(), got.GPXPath.String(), format, args...)
	testutils.AssertEqualString(t, want.MapPath.String(), got.MapPath.String(), format, args...)
	testutils.AssertEqualString(t, want.ShareableMapPath.String(), got.ShareableMapPath.String(), format, args...)
}
