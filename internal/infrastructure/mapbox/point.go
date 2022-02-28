package mapbox

import (
	"math"
	"strings"
)

var polylineFactor = math.Pow10(5)

// Points represents a list of GPS points
type Points []Point

// Point represents a GPS point
type Point struct {
	Latitude  float64
	Longitude float64
}

// PolylineEncode compress a list of GPS point using a polyline algorithm https://developers.google.com/maps/documentation/utilities/polylinealgorithm
func (pts Points) PolylineEncode() string {
	if len(pts) == 0 {
		return ""
	}

	var s strings.Builder

	var previousE5Latitude, previousE5Longitude int32

	for i := range pts {
		e5Latitude := toE5(pts[i].Latitude)
		e5Longitude := toE5(pts[i].Longitude)

		latitude := e5Latitude - previousE5Latitude
		longitude := e5Longitude - previousE5Longitude

		encodePolylineValue(&s, latitude)
		encodePolylineValue(&s, longitude)

		previousE5Latitude = e5Latitude
		previousE5Longitude = e5Longitude
	}

	return strings.ReplaceAll(s.String(), "\\", "\\\\")
}

func toE5(n float64) int32 {
	return int32(math.Round(n * polylineFactor))
}

func encodePolylineValue(line *strings.Builder, value int32) {
	var v = uint32(value << 1)
	if value < 0 {
		v = ^v
	}

	for v >= 0x20 {
		line.WriteRune(rune(0x20|(v&0x1f)) + 63)
		v = v >> 5
	}

	line.WriteRune(rune(v + 63))
}
