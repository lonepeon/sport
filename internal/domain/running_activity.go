package domain

import (
	"fmt"
	"time"
)

// RunningActivity represents a running session
type RunningActivity struct {
	Slug             RunningActivitySlug
	RanAt            time.Time
	Duration         time.Duration
	Distance         Distance
	Speed            Speed
	GPXPath          GPXFilePath
	MapPath          MapFilePath
	ShareableMapPath ShareableMapFilePath
}

func NewRunningActivity(when time.Time, duration time.Duration, distance Distance, speed Speed, gpxPath GPXFilePath, mapPath MapFilePath, shareableMapPath ShareableMapFilePath) (RunningActivity, error) {
	var err InvalidInputErrors
	err.ValidatePositiveFloat64(speed.KilometersPerHour(), "speed must be greater than 0km/h")
	err.ValidatePositiveInt(distance.Meters(), "distance must be greater than 0m")
	err.ValidateRequiredDuration(duration, "duration must be greater than 0")
	err.ValidateRequiredString(gpxPath.String(), "gpx path is required")
	err.ValidateRequiredString(mapPath.String(), "map path is required")
	err.ValidateRequiredString(shareableMapPath.String(), "shareable map path is required")

	if !err.IsEmpty() {
		return RunningActivity{}, &err
	}

	slug, errSlug := NewRunnningActivitySlugFromTime(when)
	if errSlug != nil {
		return RunningActivity{}, fmt.Errorf("cant' create activity slug: %v", err)
	}

	return RunningActivity{
		Slug:             slug,
		RanAt:            when,
		Duration:         duration,
		Distance:         distance,
		Speed:            speed,
		GPXPath:          gpxPath,
		MapPath:          mapPath,
		ShareableMapPath: shareableMapPath,
	}, nil
}
