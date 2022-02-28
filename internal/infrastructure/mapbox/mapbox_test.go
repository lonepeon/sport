package mapbox_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/domain/domaintest"
	"github.com/lonepeon/sport/internal/infrastructure/mapbox"
)

type MapboxAPIMock struct {
	Response    []byte
	Status      int
	ExpectedURL string
}

func (m MapboxAPIMock) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.URL.String() != m.ExpectedURL {
		return nil, fmt.Errorf("wrong URL. want: %s; got: %s", m.ExpectedURL, r.URL.String())
	}
	recorder := httptest.NewRecorder()
	recorder.Body = bytes.NewBuffer(m.Response)
	recorder.Code = m.Status

	return recorder.Result(), nil
}

func TestURL(t *testing.T) {
	box := mapbox.New("<token>")

	url := box.URL(mapbox.Points{
		{Latitude: 38.5, Longitude: -120.2},
		{Latitude: 40.7, Longitude: -120.95},
		{Latitude: 43.252, Longitude: -126.453},
	})

	testutils.AssertEqualString(t, "https://api.mapbox.com/styles/v1/mapbox/outdoors-v11/static/path-3+f44-0.8(_p~iF~ps%7CU_ulLnnqC_mqNvxq%60%40)/auto/800x800@2x?logo=false&access_token=<token>&padding=100", url, "wrong mapbox url")
}

func TestGenerateMapFromPointsSuccess(t *testing.T) {
	box := mapbox.New("<token>")

	pts := mapbox.Points{
		{Latitude: 38.5, Longitude: -120.2},
		{Latitude: 40.7, Longitude: -120.95},
		{Latitude: 43.252, Longitude: -126.453},
	}

	gpxFile := domaintest.NewGPXFile(t).Build()

	mock := MapboxAPIMock{
		Status:      200,
		Response:    []byte("the image bytes"),
		ExpectedURL: box.URL(pts),
	}

	box.HTTPClient = &http.Client{Transport: mock}

	image, err := box.GenerateMap(context.Background(), gpxFile)
	testutils.AssertNoError(t, err, "can't generate image")
	content, err := ioutil.ReadAll(image.File())
	testutils.AssertNoError(t, err, "can't generate image content")
	testutils.AssertEqualString(t, "the image bytes", string(content), "wrong image content")
}

func TestGenerateMapFromPointsWrongToken(t *testing.T) {
	box := mapbox.New("<token>")

	pts := mapbox.Points{
		{Latitude: 38.5, Longitude: -120.2},
		{Latitude: 40.7, Longitude: -120.95},
		{Latitude: 43.252, Longitude: -126.453},
	}

	gpxFile := domaintest.NewGPXFile(t).Build()

	mock := MapboxAPIMock{
		Status:      401,
		Response:    []byte(`{"message":"Not Authorized - Invalid Token"}`),
		ExpectedURL: box.URL(pts),
	}

	box.HTTPClient = &http.Client{Transport: mock}

	image, err := box.GenerateMap(context.Background(), gpxFile)
	if err == nil {
		content, _ := ioutil.ReadAll(image.File())
		testutils.AssertHasError(t, err, "shouldn't generate image but got one. got: %v", string(content))
	}

	testutils.AssertErrorIs(t, mapbox.ErrInvalidToken, err, "wrong error")
}

func TestGenerateMapFromPointsWrongQuery(t *testing.T) {
	box := mapbox.New("<token>")

	pts := mapbox.Points{
		{Latitude: 38.5, Longitude: -120.2},
		{Latitude: 40.7, Longitude: -120.95},
		{Latitude: 43.252, Longitude: -126.453},
	}

	gpxFile := domaintest.NewGPXFile(t).Build()

	mock := MapboxAPIMock{
		Status:      422,
		Response:    []byte(`{"message":"Auto extent cannot be determined when GeoJSON has no features"}%`),
		ExpectedURL: box.URL(pts),
	}

	box.HTTPClient = &http.Client{Transport: mock}

	image, err := box.GenerateMap(context.Background(), gpxFile)
	if err == nil {
		content, _ := ioutil.ReadAll(image.File())
		testutils.AssertHasError(t, err, "shouldn't generate image but got one. got: %v", string(content))
	}

	testutils.AssertErrorIs(t, mapbox.ErrGeneric, err, "wrong error")
}
