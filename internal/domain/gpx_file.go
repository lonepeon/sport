package domain

import (
	"bytes"
	"io"
	"time"
)

type GPXPoint struct {
	Time      time.Time
	Latitude  float64
	Longitude float64
	Duration  time.Duration
	Distance  float64
	Elevation float64
	Speed     float64
}

type GPXPoints []GPXPoint

type GPXFile struct {
	Distance Distance
	Duration time.Duration
	Speed    Speed
	Points   GPXPoints

	content []byte
}

func NewGPXFile(content []byte, distance Distance, duration time.Duration, speed Speed, points GPXPoints) GPXFile {
	return GPXFile{
		Distance: distance,
		Duration: duration,
		Speed:    speed,
		Points:   points,

		content: content,
	}
}

func (g GPXFile) File() io.Reader {
	return bytes.NewBuffer(g.content)
}
