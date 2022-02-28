package domain

import "math"

// Speed represents the average speed for an activity in km/h
type Speed struct {
	kilometersPerHour float64
}

// NewSpeedFromKmh builds a speed and validate the value is in a possible range
func NewSpeedFromKmh(kilomtersPerHour float64) (Speed, error) {
	if kilomtersPerHour < 0 {
		return Speed{}, ErrSpeedTooSmall
	}

	return Speed{kilometersPerHour: kilomtersPerHour}, nil

}

// KilometersPerHour converts speed to km/h
func (s Speed) KilometersPerHour() float64 {
	return s.kilometersPerHour
}

// MinutesPerKilometer converts speed from km/h to min/km
func (s Speed) MinutesPerKilometer() float64 {
	return math.Round(60/s.kilometersPerHour*100) / 100
}
