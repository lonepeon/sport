package domaintest

import (
	"math/rand"
	"time"

	// embed is used to load hardcoded GPX files
	_ "embed"
)

//go:embed content.gpx
var gpxContent1 []byte

func GetGPXBytes() []byte {
	return gpxContent1
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func intBetween(min, max int) int {
	return random.Intn(max-min) + min
}

func durationBetween(min, max int) time.Duration {
	return time.Duration(intBetween(min, max))
}
