package math_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/infrastructure/gpx/internal/math"
)

func TestKilometerPerHour(t *testing.T) {
	tcs := []struct {
		Meters   float64
		Duration time.Duration
		Speed    float64
	}{
		{
			Meters:   1,
			Duration: time.Second,
			Speed:    3.6,
		},
		{
			Meters:   4,
			Duration: time.Second,
			Speed:    14.4,
		},
		{
			Meters:   4178,
			Duration: time.Duration(25*time.Minute + 607*time.Millisecond),
			Speed:    10.02,
		},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("%f meters in %f seconds => %fkm/h", tc.Meters, tc.Duration.Seconds(), tc.Speed), func(t *testing.T) {
			speed := math.KilometerPerHour(tc.Meters, tc.Duration)

			testutils.AssertEqualFloat64(t, tc.Speed, speed, "wrong speed")
		})
	}
}
