package math

import (
	"math"
	"time"
)

// MeterSecondsToKilometerHourRatio mutiplier to transform m/s to km/h
const MeterSecondsToKilometerHourRatio = 3.6

// KilometerPerHour calculate the speed in kilometer per hour based on the distance in meters and the elapsed duration
func KilometerPerHour(meters float64, duration time.Duration) float64 {
	return round2Decimals((meters / duration.Seconds()) * MeterSecondsToKilometerHourRatio)
}

func round2Decimals(n float64) float64 {
	return math.Round(n*100) / 100
}
