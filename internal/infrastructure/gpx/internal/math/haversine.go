package math

import (
	"math"
)

// EarthRadiusMeters is the radius of Earth in meters
const EarthRadiusMeters = 6371e3

// HaversineOnEarth calculate the distance from two points using the haversine function
// based on the Earth sphere
func HaversineOnEarth(lat1, lon1, lat2, lon2 float64) float64 {
	return haversine(EarthRadiusMeters, lat1, lon1, lat2, lon2)
}

func haversine(r, lat1, lon1, lat2, lon2 float64) float64 {
	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	deltaPhi := (lat2 - lat1) * math.Pi / 180
	deltaLambda := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return float64(r) * c
}
