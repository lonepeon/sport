package mapbox_test

import (
	"testing"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/infrastructure/mapbox"
)

func TestPolyline(t *testing.T) {
	polyline := mapbox.Points{
		{Latitude: 38.5, Longitude: -120.2},
		{Latitude: 40.7, Longitude: -120.95},
		{Latitude: 43.252, Longitude: -126.453},
	}

	testutils.AssertEqualString(t, "_p~iF~ps|U_ulLnnqC_mqNvxq`@", polyline.PolylineEncode(), "wrong polyline")
}
