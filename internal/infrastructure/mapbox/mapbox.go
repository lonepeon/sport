package mapbox

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/lonepeon/sport/internal/domain"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrGeneric      = errors.New("something wrong happened")
)

type Mapbox struct {
	HTTPClient    *http.Client
	EndpointURL   string
	token         string
	theme         string
	size          string
	lineColor     string
	lineThinkness int
	lineOpacity   float64
	padding       int
}

func New(token string) *Mapbox {
	return &Mapbox{
		HTTPClient:    http.DefaultClient,
		EndpointURL:   "https://api.mapbox.com",
		token:         token,
		theme:         "mapbox/outdoors-v11",
		lineThinkness: 3,
		lineOpacity:   0.8,
		lineColor:     "f44",
		size:          "800x800@2x",
		padding:       100,
	}
}

func (m *Mapbox) URL(pts Points) string {
	line := url.QueryEscape(pts.PolylineEncode())
	return fmt.Sprintf(
		"%s/styles/v1/%s/static/path-%d+%s-%.1f(%s)/auto/%s?logo=false&access_token=%s&padding=%d",
		m.EndpointURL,
		m.theme,
		m.lineThinkness, m.lineColor, m.lineOpacity,
		line,
		m.size,
		m.token,
		m.padding,
	)
}

func (m *Mapbox) GenerateMap(ctx context.Context, gpx domain.GPXFile) (domain.MapFile, error) {
	url := m.URL(m.mapPoints(gpx.Points))

	resp, err := m.HTTPClient.Get(url)
	if err != nil {
		return domain.MapFile{}, fmt.Errorf("can't fetch image (url=%v): %v", url, err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusUnauthorized:
		body, _ := ioutil.ReadAll(resp.Body)
		return domain.MapFile{}, fmt.Errorf("can't fetch image (url=%v, status=%d, body=%v): %w", url, resp.StatusCode, body, ErrInvalidToken)

	default:
		body, _ := ioutil.ReadAll(resp.Body)
		return domain.MapFile{}, fmt.Errorf("can't fetch image (url=%v, status=%d, body=%v): %w", url, resp.StatusCode, body, ErrGeneric)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return domain.MapFile{}, fmt.Errorf("%w: can't write map image: %v", ErrGeneric, err)
	}

	return domain.NewMapFile(buf.Bytes()), nil
}

func (m *Mapbox) mapPoints(gpxPoints domain.GPXPoints) Points {
	points := make(Points, len(gpxPoints))
	for i := range gpxPoints {
		points[i] = Point{
			Latitude:  gpxPoints[i].Latitude,
			Longitude: gpxPoints[i].Longitude,
		}
	}

	return points
}
