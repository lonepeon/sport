package testutils

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

type Formatter func(string, ...interface{})

func checksHasError(log Formatter, t *testing.T, got error, format string, args ...interface{}) {
	t.Helper()

	if got == nil {
		log("%s\nerr: nil\n", fmt.Sprintf(format, args...))
	}
}

func checksErrorContains(log Formatter, t *testing.T, substring string, got error, format string, args ...interface{}) {
	t.Helper()

	if !strings.Contains(got.Error(), substring) {
		log("%s\nwant substring: %s\ngot: %s", fmt.Sprintf(format, args...), substring, got)
	}
}

func checksErrorAs(log Formatter, t *testing.T, want interface{}, got error, format string, args ...interface{}) {
	t.Helper()

	checksHasError(log, t, got, format, args...)

	if got == want || errors.As(got, want) {
		return
	}

	log("%s\nwant: %#+v\ngot:  %#+v\n", fmt.Sprintf(format, args...), want, got)
}

func checksErrorIs(log Formatter, t *testing.T, want error, got error, format string, args ...interface{}) {
	t.Helper()

	checksHasError(log, t, got, format, args...)

	if got == want || errors.Is(got, want) {
		return
	}

	log("%s\nwant: %#+v\ngot:  %#+v\n", fmt.Sprintf(format, args...), want, got)
}

func checksNoError(log Formatter, t *testing.T, err error, format string, args ...interface{}) {
	t.Helper()

	if err != nil {
		log("%s\nerr: %s\n", fmt.Sprintf(format, args...), err)
	}
}

func checksEqualBool(log Formatter, t *testing.T, want bool, got bool, format string, args ...interface{}) {
	t.Helper()

	if want != got {
		log("%s\nwant: %t\ngot:  %t\n", fmt.Sprintf(format, args...), want, got)
	}
}

func checksEqualStrings(log Formatter, t *testing.T, want []string, got []string, format string, args ...interface{}) {
	t.Helper()

	if len(want) != len(got) {
		log("%s\nwant: %d items\ngot:  %d items\n", fmt.Sprintf(format, args...), len(want), len(got))
		return
	}

	for i := range want {
		checksEqualString(log, t, want[i], got[i], fmt.Sprintf("item %d: %s", i, fmt.Sprintf(format, args...)))
	}
}

func checksContainsString(log Formatter, t *testing.T, want string, fullmsg string, format string, args ...interface{}) {
	t.Helper()

	if !strings.Contains(fullmsg, want) {
		log("%s\nwant:   %s\nstring: %s\n", fmt.Sprintf(format, args...), want, fullmsg)
	}
}

func checksEqualString(log Formatter, t *testing.T, want string, got string, format string, args ...interface{}) {
	t.Helper()

	if want != got {
		log("%s\nwant: %s\ngot:  %s\n", fmt.Sprintf(format, args...), want, got)
	}
}

func checksEqualFloat64(log Formatter, t *testing.T, want float64, got float64, format string, args ...interface{}) {
	t.Helper()

	if want != got {
		log("%s\nwant: %f\ngot:  %f\n", fmt.Sprintf(format, args...), want, got)
	}
}

func checksEqualInt(log Formatter, t *testing.T, want int, got int, format string, args ...interface{}) {
	t.Helper()

	if want != got {
		log("%s\nwant: %d\ngot:  %d\n", fmt.Sprintf(format, args...), want, got)
	}
}

func checksNotEqualNil(log Formatter, t *testing.T, got interface{}, format string, args ...interface{}) {
	t.Helper()

	if got == nil {
		log("%s\nwant: not nil\ngot:  nil\n", fmt.Sprintf(format, args...))
	}
}

func checksEqualNil(log Formatter, t *testing.T, got interface{}, format string, args ...interface{}) {
	t.Helper()

	if got != nil {
		log("%s\nwant: nil\ngot:  %v\n", fmt.Sprintf(format, args...), got)
	}
}

func checksEqualInt64(log Formatter, t *testing.T, want int64, got int64, format string, args ...interface{}) {
	t.Helper()

	if want != got {
		log("%s\nwant: %d\ngot:  %d\n", fmt.Sprintf(format, args...), want, got)
	}
}

func checksEqualDuration(log Formatter, t *testing.T, want time.Duration, got time.Duration, format string, args ...interface{}) {
	t.Helper()

	if want != got {
		log("%s\nwant: %s\ngot:  %s\n", fmt.Sprintf(format, args...), want.String(), got.String())
	}
}

func checksEqualTime(log Formatter, t *testing.T, want time.Time, got time.Time, format string, args ...interface{}) {
	t.Helper()

	if want.UnixMilli() != got.UnixMilli() {
		log("%s\nwant: %s\ngot:  %s\n", fmt.Sprintf(format, args...), want.String(), got.String())
	}
}
