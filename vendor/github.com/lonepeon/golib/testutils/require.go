package testutils

import (
	"testing"
	"time"
)

func RequireHasError(t *testing.T, got error, format string, args ...interface{}) {
	t.Helper()

	checksHasError(t.Fatalf, t, got, format, args...)
}

func RequireErrorContains(t *testing.T, substring string, got error, format string, args ...interface{}) {
	t.Helper()

	checksErrorContains(t.Fatalf, t, substring, got, format, args...)
}

func RequireErrorAs(t *testing.T, want interface{}, got error, format string, args ...interface{}) {
	t.Helper()

	checksErrorAs(t.Fatalf, t, want, got, format, args...)
}

func RequireErrorIs(t *testing.T, want error, got error, format string, args ...interface{}) {
	t.Helper()

	checksErrorIs(t.Fatalf, t, want, got, format, args...)
}

func RequireNoError(t *testing.T, err error, format string, args ...interface{}) {
	t.Helper()

	checksNoError(t.Fatalf, t, err, format, args...)
}

func RequireEqualBool(t *testing.T, want bool, got bool, format string, args ...interface{}) {
	checksEqualBool(t.Fatalf, t, want, got, format, args...)
}

func RequireEqualBytes(t *testing.T, want []byte, got []byte, format string, args ...interface{}) {
	checksEqualBytes(t.Fatalf, t, want, got, format, args...)
}

func RequireEqualStrings(t *testing.T, want []string, got []string, format string, args ...interface{}) {
	checksEqualStrings(t.Fatalf, t, want, got, format, args...)
}

func RequireContainsString(t *testing.T, want string, fullmsg string, format string, args ...interface{}) {
	checksContainsString(t.Fatalf, t, want, fullmsg, format, args...)
}

func RequireEqualString(t *testing.T, want string, got string, format string, args ...interface{}) {
	checksEqualString(t.Fatalf, t, want, got, format, args...)
}

func RequireEqualFloat64(t *testing.T, want float64, got float64, format string, args ...interface{}) {
	checksEqualFloat64(t.Fatalf, t, want, got, format, args...)
}

func RequireEqualInt(t *testing.T, want int, got int, format string, args ...interface{}) {
	checksEqualInt(t.Fatalf, t, want, got, format, args...)
}

func RequireNotEqualNil(t *testing.T, got interface{}, format string, args ...interface{}) {
	checksNotEqualNil(t.Fatalf, t, got, format, args...)
}

func RequireEqualNil(t *testing.T, got interface{}, format string, args ...interface{}) {
	checksEqualNil(t.Fatalf, t, got, format, args...)
}

func RequireEqualInt64(t *testing.T, want int64, got int64, format string, args ...interface{}) {
	checksEqualInt64(t.Fatalf, t, want, got, format, args...)
}

func RequireEqualDuration(t *testing.T, want time.Duration, got time.Duration, format string, args ...interface{}) {
	checksEqualDuration(t.Fatalf, t, want, got, format, args...)
}

func RequireEqualTime(t *testing.T, want time.Time, got time.Time, format string, args ...interface{}) {
	checksEqualTime(t.Fatalf, t, want, got, format, args...)
}

func RequireNotEmptyString(t *testing.T, got string, format string, args ...interface{}) {
	t.Helper()

	checksNotEmptyString(t.Fatalf, t, got, format, args...)
}
