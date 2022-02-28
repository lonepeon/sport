package domain_test

import (
	"testing"
	"time"

	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/sport/internal/domain"
)

func TestNewRunnginActivityErrors(t *testing.T) {
	distance, err := domain.NewDistanceFromMeters(0)
	testutils.AssertNoError(t, err, "can't create distance of 0m")

	speed, err := domain.NewSpeedFromKmh(0)
	testutils.AssertNoError(t, err, "can't create speed of 0km/h")

	_, err = domain.NewRunningActivity(
		time.Now(),
		time.Duration(0),
		distance,
		speed,
		domain.GPXFilePath(""),
		domain.MapFilePath(""),
		domain.ShareableMapFilePath(""),
	)

	var inputErr *domain.InvalidInputErrors
	testutils.AssertErrorAs(t, &inputErr, err, "didn't get the expected error")
	testutils.AssertEqualBool(t, false, inputErr.IsEmpty(), "didn't expect to get an empty input error")
	errorMessages := inputErr.Detail()
	testutils.AssertEqualString(t, "speed must be greater than 0km/h", errorMessages[0], "wrong speed error")
	testutils.AssertEqualString(t, "distance must be greater than 0m", errorMessages[1], "wrong distance error")
	testutils.AssertEqualString(t, "duration must be greater than 0", errorMessages[2], "wrong duration error")
	testutils.AssertEqualString(t, "gpx path is required", errorMessages[3], "wrong gpx path error")
	testutils.AssertEqualString(t, "map path is required", errorMessages[4], "wrong map path error")
	testutils.AssertEqualString(t, "shareable map path is required", errorMessages[5], "wrong shareable map path error")
}
