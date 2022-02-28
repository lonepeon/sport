package domain_test

import (
	"fmt"
	"testing"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/domain"
)

func TestNewSpeedFromKmhSuccess(t *testing.T) {
	kmh := 10.1
	distance, err := domain.NewSpeedFromKmh(kmh)
	testutils.AssertNoError(t, err, "can't build speed of %f kmh", kmh)
	testutils.AssertEqualFloat64(t, kmh, distance.KilometersPerHour(), "unexpected number of kilometers per hour")
}

func TestNewSpeedFromKmhNegativeError(t *testing.T) {
	kmh := -5.0
	_, err := domain.NewSpeedFromKmh(kmh)
	testutils.AssertErrorIs(t, domain.ErrSpeedTooSmall, err, "unexpected speed error")
}

func TestSpeedMinutesPerKilometer(t *testing.T) {
	tcs := map[float64]float64{
		10:    6,
		12.45: 4.82,
	}

	for kmh, expected := range tcs {
		t.Run(fmt.Sprintf("%vkm/h -> %vm/km", kmh, expected), func(t *testing.T) {
			speed, err := domain.NewSpeedFromKmh(kmh)
			testutils.AssertNoError(t, err, "can't build speed of %f km/h", kmh)
			testutils.AssertEqualFloat64(t, expected, speed.MinutesPerKilometer(), "invalid conversion")
		})
	}
}
