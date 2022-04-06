package testutils

import (
	"testing"
	"time"
)

func AssertHasError(t *testing.T, got error, format string, args ...interface{}) {
	t.Helper()

	checksHasError(t.Errorf, t, got, format, args...)
}

func AssertErrorContains(t *testing.T, substring string, got error, format string, args ...interface{}) {
	t.Helper()

	checksErrorContains(t.Errorf, t, substring, got, format, args...)
}

func AssertErrorAs(t *testing.T, want interface{}, got error, format string, args ...interface{}) {
	t.Helper()

	checksErrorAs(t.Errorf, t, want, got, format, args...)
}

func AssertErrorIs(t *testing.T, want error, got error, format string, args ...interface{}) {
	t.Helper()

	checksErrorIs(t.Errorf, t, want, got, format, args...)
}

func AssertNoError(t *testing.T, err error, format string, args ...interface{}) {
	t.Helper()

	checksNoError(t.Errorf, t, err, format, args...)
}

func AssertEqualBool(t *testing.T, want bool, got bool, format string, args ...interface{}) {
	t.Helper()

	checksEqualBool(t.Errorf, t, want, got, format, args...)
}

func AssertEqualBytes(t *testing.T, want []byte, got []byte, format string, args ...interface{}) {
	t.Helper()

	checksEqualBytes(t.Errorf, t, want, got, format, args...)
}

func AssertEqualStrings(t *testing.T, want []string, got []string, format string, args ...interface{}) {
	t.Helper()

	checksEqualStrings(t.Errorf, t, want, got, format, args...)
}

func AssertContainsString(t *testing.T, want string, fullmsg string, format string, args ...interface{}) {
	t.Helper()

	checksContainsString(t.Errorf, t, want, fullmsg, format, args...)
}

func AssertEqualString(t *testing.T, want string, got string, format string, args ...interface{}) {
	t.Helper()

	checksEqualString(t.Errorf, t, want, got, format, args...)
}

func AssertEqualFloat64(t *testing.T, want float64, got float64, format string, args ...interface{}) {
	t.Helper()

	checksEqualFloat64(t.Errorf, t, want, got, format, args...)
}

func AssertEqualInt(t *testing.T, want int, got int, format string, args ...interface{}) {
	t.Helper()

	checksEqualInt(t.Errorf, t, want, got, format, args...)
}

func AssertNotEqualNil(t *testing.T, got interface{}, format string, args ...interface{}) {
	t.Helper()

	checksNotEqualNil(t.Errorf, t, got, format, args...)
}

func AssertEqualNil(t *testing.T, got interface{}, format string, args ...interface{}) {
	t.Helper()

	checksEqualNil(t.Errorf, t, got, format, args...)
}

func AssertEqualInt64(t *testing.T, want int64, got int64, format string, args ...interface{}) {
	t.Helper()

	checksEqualInt64(t.Errorf, t, want, got, format, args...)
}

func AssertEqualDuration(t *testing.T, want time.Duration, got time.Duration, format string, args ...interface{}) {
	t.Helper()

	checksEqualDuration(t.Errorf, t, want, got, format, args...)
}

func AssertEqualTime(t *testing.T, want time.Time, got time.Time, format string, args ...interface{}) {
	t.Helper()

	checksEqualTime(t.Errorf, t, want, got, format, args...)
}

func AssertNotEmptyString(t *testing.T, got string, format string, args ...interface{}) {
	t.Helper()

	checksNotEmptyString(t.Errorf, t, got, format, args...)
}
