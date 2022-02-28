package domaintest

import (
	"fmt"
	"testing"
	"time"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/domain"
	"github.com/lonepeon/sport/internal/repository"
	"golang.org/x/net/context"
)

type Speed struct {
	t     *testing.T
	speed domain.Speed
}

func NewSpeed(t *testing.T) Speed {
	kmh := 10.0
	speed, err := domain.NewSpeedFromKmh(kmh)
	testutils.AssertNoError(t, err, "can't generate speed of %f kn/h", kmh)

	return Speed{t: t, speed: speed}
}

func (s Speed) WithKmh(kmh float64) Speed {
	speed, err := domain.NewSpeedFromKmh(kmh)
	testutils.AssertNoError(s.t, err, "can't generate speed of %f km/h", speed)

	s.speed = speed
	return s
}

func (s Speed) Build() domain.Speed {
	return s.speed
}

type Distance struct {
	t        *testing.T
	distance domain.Distance
}

func NewDistance(t *testing.T) Distance {
	meters := intBetween(5000, 10000)
	distance, err := domain.NewDistanceFromMeters(meters)
	testutils.AssertNoError(t, err, "can't generate distance of %d meters", meters)

	return Distance{t: t, distance: distance}
}

func (d Distance) WithMeters(meters int) Distance {
	distance, err := domain.NewDistanceFromMeters(meters)
	testutils.AssertNoError(d.t, err, "can't generate distance of %d meters", meters)

	d.distance = distance
	return d
}

func (d Distance) Build() domain.Distance {
	return d.distance
}

type GPXFile struct {
	t        *testing.T
	content  []byte
	distance domain.Distance
	duration time.Duration
	speed    domain.Speed
	points   domain.GPXPoints
}

func NewGPXFile(t *testing.T) GPXFile {
	meters := intBetween(5000, 10000)
	distance, err := domain.NewDistanceFromMeters(meters)
	testutils.AssertNoError(t, err, "can't generate gpx file with a distance of %d meters", meters)

	kmh := 10.0
	speed, err := domain.NewSpeedFromKmh(kmh)
	testutils.AssertNoError(t, err, "can't generate gpx file with a speed of %f km/h", kmh)

	ranAt := time.Now().
		Add(-durationBetween(1, 24*30*12) * time.Hour)

	duration := durationBetween(30, 60) * time.Minute

	points := []domain.GPXPoint{
		{
			Latitude:  38.5,
			Longitude: -120.2,
			Time:      ranAt,
			Distance:  55,
			Elevation: 3,
			Speed:     10.25,
		},
		{
			Latitude:  40.7,
			Longitude: -120.95,
			Time:      ranAt.Add(20 * time.Second),
			Distance:  43,
			Elevation: 8,
			Speed:     9.80,
		},
		{
			Latitude:  43.252,
			Longitude: -126.453,
			Time:      ranAt.Add(40 * time.Second),
			Distance:  50,
			Elevation: 1,
			Speed:     10.01,
		},
	}

	return GPXFile{
		t:        t,
		content:  gpxContent1,
		distance: distance,
		duration: duration,
		speed:    speed,
		points:   points,
	}
}

func (g GPXFile) WithFileContent(content []byte) GPXFile {
	g.content = content
	return g
}

func (g GPXFile) WithPoints(points domain.GPXPoints) GPXFile {
	g.points = points
	return g
}

func (g GPXFile) Build() domain.GPXFile {
	return domain.NewGPXFile(g.content, g.distance, g.duration, g.speed, g.points)
}

type RunningActivity struct {
	t        *testing.T
	ranAt    time.Time
	duration time.Duration
	distance domain.Distance
	speed    domain.Speed
}

func NewRunningActivity(t *testing.T) RunningActivity {
	ranAt := time.Now().
		Add(-durationBetween(1, 24*30*12) * time.Hour)

	meters := intBetween(5000, 10000)
	distance, err := domain.NewDistanceFromMeters(meters)
	testutils.AssertNoError(t, err, "can't generate activity with a distance of %d meters", meters)

	kmh := 10.0
	speed, err := domain.NewSpeedFromKmh(kmh)
	testutils.AssertNoError(t, err, "can't generate activity with a speed of %f km/h", kmh)

	return RunningActivity{
		t:        t,
		ranAt:    ranAt,
		distance: distance,
		speed:    speed,
		duration: durationBetween(30, 60) * time.Minute,
	}
}

func (r RunningActivity) WithDuration(d time.Duration) RunningActivity {
	r.duration = d

	return r
}

func (r RunningActivity) WithRawSlug(s string) RunningActivity {
	slug, err := domain.NewRunnningActivitySlugFromString(s)
	testutils.AssertNoError(r.t, err, "invalid slug in running activity builder")

	r.ranAt = slug.Time()

	return r
}

func (r RunningActivity) WithSpeedKmh(kmh float64) RunningActivity {
	speed, err := domain.NewSpeedFromKmh(kmh)
	testutils.AssertNoError(r.t, err, "can't build speed of %f km/h", kmh)

	r.speed = speed

	return r
}

func (r RunningActivity) WithDistanceMeters(meters int) RunningActivity {
	distance, err := domain.NewDistanceFromMeters(meters)
	testutils.AssertNoError(r.t, err, "can't build distance of %d meters", meters)

	r.distance = distance

	return r
}

func (r RunningActivity) Build() domain.RunningActivity {
	activity, err := domain.NewRunningActivity(
		r.ranAt,
		r.duration,
		r.distance,
		r.speed,
		domain.GPXFilePath(fmt.Sprintf("runs/%s/run.gpx", r.ranAt.Format("2006-01-02.15h04"))),
		domain.MapFilePath(fmt.Sprintf("runs/%s/map.png", r.ranAt.Format("2006-01-02.15h04"))),
		domain.ShareableMapFilePath(fmt.Sprintf("runs/%s/share-map.png", r.ranAt.Format("2006-01-02.15h04"))),
	)

	testutils.AssertNoError(r.t, err, "can't generate activity")

	return activity
}

func (r RunningActivity) Persist(w repository.Writer) domain.RunningActivity {
	activity := r.Build()
	err := w.RecordRunningActivity(context.Background(), activity)
	testutils.AssertNoError(r.t, err, "can't persist activity")

	return activity
}
