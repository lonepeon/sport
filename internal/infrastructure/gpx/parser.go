package gpx

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/lonepeon/sport/internal/infrastructure/gpx/internal/math"
)

var (
	ErrFormat = errors.New("invalid file format")
)

// TrackSegment represents a GPX track segment: https://en.wikipedia.org/wiki/GPS_Exchange_Format
type TrackSegment struct {
	Points   []TrackPoint
	Speed    float64
	Duration time.Duration
	Distance int
}

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

// DistanceFrom calculates the distance in meters between two GPS coordinate, on earth
func (to Coordinate) DistanceFrom(from Coordinate) float64 {
	return math.HaversineOnEarth(from.Latitude, from.Longitude, to.Latitude, to.Longitude)
}

type TrackPoint struct {
	Time       time.Time
	Coordinate Coordinate
	Duration   time.Duration
	Distance   float64
	Elevation  float64
	Speed      float64
}

// ParseTrackSegment reads a GPX XML file and expect to find a track with one and only one segment.
// It returns a TrackSegment able to manipulate it.
// If it can't decode the content of the GPX file, it returns a ErrFormat
func ParseTrackSegment(r io.Reader) (TrackSegment, error) {
	trkpts, err := loadPointsFromXML(r)
	if err != nil {
		return TrackSegment{}, err
	}

	overallDuration := trkpts[len(trkpts)-1].Time.Sub(trkpts[0].Time)

	var overallDistance float64

	previousCoordinate := Coordinate{
		Latitude:  trkpts[0].Latitude,
		Longitude: trkpts[0].Longitude,
	}

	previousTime := trkpts[0].Time

	points := make([]TrackPoint, len(trkpts))
	for i := range trkpts {
		coordinate := Coordinate{
			Latitude:  trkpts[i].Latitude,
			Longitude: trkpts[i].Longitude,
		}

		distance := coordinate.DistanceFrom(previousCoordinate)
		overallDistance += distance

		duration := trkpts[i].Time.Sub(previousTime)

		points[i] = TrackPoint{
			Distance:   distance,
			Coordinate: coordinate,
			Duration:   duration,
			Time:       trkpts[i].Time,
			Elevation:  trkpts[i].Elevation,
			Speed:      math.KilometerPerHour(distance, duration),
		}

		previousCoordinate = points[i].Coordinate
		previousTime = points[i].Time
	}

	return TrackSegment{
		Speed:    math.KilometerPerHour(overallDistance, overallDuration),
		Duration: overallDuration,
		Distance: int(overallDistance),
		Points:   points,
	}, nil
}

func (s TrackSegment) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var segment XMLTrackSegment
	segment.Points = make([]XMLTrackPoint, len(s.Points))
	for i := range s.Points {
		segment.Points[i] = XMLTrackPoint{
			Latitude:  s.Points[i].Coordinate.Latitude,
			Longitude: s.Points[i].Coordinate.Longitude,
			Time:      s.Points[i].Time,
			Elevation: s.Points[i].Elevation,
		}
	}

	start.Name = xml.Name{Local: "gpx"}
	gpx := XMLGPX{
		Tracks: []XMLTrack{
			{
				Segments: []XMLTrackSegment{segment},
			},
		},
	}

	return e.EncodeElement(gpx, start)
}

func loadPointsFromXML(r io.Reader) ([]XMLTrackPoint, error) {
	var gpx XMLGPX
	if err := xml.NewDecoder(r).Decode(&gpx); err != nil {
		return nil, fmt.Errorf("can't decode file: %w: %v", ErrFormat, err)
	}

	if len(gpx.Tracks) != 1 {
		return nil, fmt.Errorf("expected one track (tracks=%d): %w", len(gpx.Tracks), ErrFormat)
	}

	trksegs := gpx.Tracks[0].Segments
	if len(trksegs) != 1 {
		return nil, fmt.Errorf("expected one segment (segments=%d): %w", len(trksegs), ErrFormat)
	}

	trkpts := trksegs[0].Points
	if len(trkpts) == 0 {
		return nil, fmt.Errorf("segment contains no points: %w", ErrFormat)
	}

	return trkpts, nil
}
