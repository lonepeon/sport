package math_test

import (
	"fmt"
	"testing"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/infrastructure/gpx/internal/math"
)

func TestHaversineOnEarth(t *testing.T) {
	tcs := []struct {
		Lat1     float64
		Lon1     float64
		Lat2     float64
		Lon2     float64
		Distance float64
	}{
		{
			Lat1:     50.69410745181364,
			Lon1:     3.271642034714101,
			Lat2:     50.69410109551195,
			Lon2:     3.271613341362284,
			Distance: 2.141111255802794,
		},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("Pt(%f, %f) <-> Pt(%f, %f) => %f", tc.Lat1, tc.Lon1, tc.Lat2, tc.Lon2, tc.Distance), func(t *testing.T) {
			distance := math.HaversineOnEarth(tc.Lat1, tc.Lon1, tc.Lat2, tc.Lon2)

			testutils.AssertEqualFloat64(t, tc.Distance, distance, "wrong distance")
		})
	}
}
