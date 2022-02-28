package gpx_test

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/infrastructure/gpx"
)

func TestParseFile(t *testing.T) {
	fname := "testdata/valid.gpx"
	file, err := os.Open(fname)
	testutils.AssertNoError(t, err, "can't open test file: %v", err)

	segment, err := gpx.ParseTrackSegment(file)
	testutils.AssertNoError(t, err, "can't parse gpx file (file=%s): %v", fname, err)

	testutils.AssertEqualInt(t, 1467, len(segment.Points), "unexpected number of points")
	testutils.AssertEqualDuration(t, time.Duration(25*time.Minute+607*time.Millisecond), segment.Duration, "unexpected duration")
	testutils.AssertEqualInt(t, 4178, segment.Distance, "unexpected total distance")
	testutils.AssertEqualFloat64(t, 10.03, segment.Speed, "unexpected average speed")
}

func TestMarshalFile(t *testing.T) {
	fname := "testdata/valid.gpx"
	file, err := os.Open(fname)
	testutils.AssertNoError(t, err, "can't open test file: %v", err)

	segment, err := gpx.ParseTrackSegment(file)
	testutils.AssertNoError(t, err, "can't parse gpx file (file=%s): %v", fname, err)

	result, err := xml.Marshal(segment)
	testutils.AssertNoError(t, err, "can't marshal segment to gpx: %v", err)

	expected, err := ioutil.ReadFile("testdata/golden.gpx")
	testutils.AssertNoError(t, err, "can't load golden file: %v", err)

	testutils.AssertEqualString(t, strings.TrimSpace(string(expected)), string(result), "unexpected gpx result")
}
