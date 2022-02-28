package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type InvalidInputErrors []string

func (e *InvalidInputErrors) Append(msg string) {
	*e = append(*e, msg)
}

func (e *InvalidInputErrors) ValidateRequiredString(value string, errorMessage string) {
	if value == "" {
		e.Append(errorMessage)
	}
}

func (e *InvalidInputErrors) ValidatePositiveFloat64(value float64, errorMessage string) {
	if value <= 0 {
		e.Append(errorMessage)
	}
}
func (e *InvalidInputErrors) ValidatePositiveInt(value int, errorMessage string) {
	if value <= 0 {
		e.Append(errorMessage)
	}
}

func (e *InvalidInputErrors) ValidateRequiredDuration(value time.Duration, errorMessage string) {
	if value <= 0 {
		e.Append(errorMessage)
	}
}

// IsEmpty returns wether the error contains any error
func (e *InvalidInputErrors) IsEmpty() bool {
	return len(*e) == 0
}

// Detail return all the error messages
func (e *InvalidInputErrors) Detail() []string {
	return *e
}

// Error implements the Error interface
func (e *InvalidInputErrors) Error() string {
	details := strings.Join(*e, "; ")
	return fmt.Sprintf("invalid input: %s", details)
}

// ErrDistanceTooSmall is returned when a distance is built with a too small number
var ErrDistanceTooSmall = errors.New("distance can't be less than 0 meters")

// ErrSpeedTooSmall is returned when speed is built with a too small number
var ErrSpeedTooSmall = errors.New("speed can't be less than 0 km/h")

// ErrCantGetRunningSession is returned when a GetRunningSession usecase can't retrieve an activity
var ErrCantGetRunningSession = errors.New("running session not found")
