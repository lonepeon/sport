package domain

import "math"

// Distance represents how much meters has been traveled
type Distance struct {
	meters int
}

// NewDistanceFromMeters builds a distance and validate the value is in a possible range
func NewDistanceFromMeters(meters int) (Distance, error) {
	if meters < 0 {
		return Distance{}, ErrDistanceTooSmall
	}

	return Distance{meters: meters}, nil
}

// Meters converts a meter distance to its representing integer
func (d Distance) Meters() int {
	return d.meters
}

// Kilometers converts the distance from meters to kilometers
func (d Distance) Kilometers() float64 {
	return math.Round(float64(d.meters)/10) / 100.0
}
