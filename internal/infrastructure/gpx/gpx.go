package gpx

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/lonepeon/sport/internal/domain"
)

type GPX struct {
}

func (GPX) CleanGPXFile(ctx context.Context, r io.Reader) (domain.GPXFile, error) {
	segment, err := ParseTrackSegment(r)
	if err != nil {
		return domain.GPXFile{}, fmt.Errorf("can't parse track segment: %v", err)
	}

	distance, err := domain.NewDistanceFromMeters(segment.Distance)
	if err != nil {
		return domain.GPXFile{}, fmt.Errorf("can't parse segment distance: %v", err)
	}

	speed, err := domain.NewSpeedFromKmh(segment.Speed)
	if err != nil {
		return domain.GPXFile{}, fmt.Errorf("can't parse segment speed: %v", err)
	}

	gpxSegment, err := xml.Marshal(segment)
	if err != nil {
		return domain.GPXFile{}, fmt.Errorf("can't build clean gpx file: %v", err)
	}

	return domain.NewGPXFile(
		gpxSegment,
		distance,
		segment.Duration,
		speed,
		gpxPointsToDomainPoints(segment.Points),
	), nil
}

func gpxPointsToDomainPoints(gpxPoints []TrackPoint) []domain.GPXPoint {
	domainPoints := make([]domain.GPXPoint, len(gpxPoints))
	for i := range gpxPoints {
		domainPoints[i] = domain.GPXPoint{
			Time:      gpxPoints[i].Time,
			Latitude:  gpxPoints[i].Coordinate.Latitude,
			Longitude: gpxPoints[i].Coordinate.Longitude,
			Duration:  gpxPoints[i].Duration,
			Distance:  gpxPoints[i].Distance,
			Elevation: gpxPoints[i].Elevation,
			Speed:     gpxPoints[i].Speed,
		}
	}

	return domainPoints
}
