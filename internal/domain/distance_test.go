package domain_test

import (
	"fmt"
	"testing"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/domain"
)

func TestNewDistanceFromMetersSuccess(t *testing.T) {
	meters := 1200
	distance, err := domain.NewDistanceFromMeters(meters)
	testutils.AssertNoError(t, err, "can't build distance of %d meters", meters)
	testutils.AssertEqualInt(t, meters, distance.Meters(), "unexpected number of meters")
}

func TestNewDistanceFromMetersNegativeError(t *testing.T) {
	meters := -5
	_, err := domain.NewDistanceFromMeters(meters)
	testutils.AssertErrorIs(t, domain.ErrDistanceTooSmall, err, "unexected distance error")
}

func TestDistanceFromMetersToKilometers(t *testing.T) {
	tcs := map[int]float64{
		1200: 1.2,
		3564: 3.56,
		3568: 3.57,
		100:  0.1,
	}

	for meters, expected := range tcs {
		t.Run(fmt.Sprintf("%vm -> %vkm", meters, expected), func(t *testing.T) {
			distance, err := domain.NewDistanceFromMeters(meters)
			testutils.AssertNoError(t, err, "can't build distance with of %d meters", meters)
			testutils.AssertEqualFloat64(t, expected, distance.Kilometers(), "invalid conversion")
		})
	}
}
