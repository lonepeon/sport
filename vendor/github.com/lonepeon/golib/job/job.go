package job

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

//go:generate go run ../sqlutil/cmd/sql-migration ./scripts

const (
	DefaultMaxAttempts = 10
)

type Job struct {
	Name        string
	At          time.Time
	MaxAttempts int

	id       string
	params   []byte
	attempts int
}

func NewJob(name string, params interface{}) (Job, error) {
	id := uuid.NewString()

	p, err := json.Marshal(params)
	if err != nil {
		return Job{}, fmt.Errorf("can't marshal job params to json: %w: %v", ErrGeneric, err)
	}

	return Job{
		Name:        name,
		MaxAttempts: DefaultMaxAttempts,
		At:          time.Now(),

		id:       id,
		attempts: 1,
		params:   p,
	}, nil
}

func (j Job) EncodedParams() []byte {
	return j.params
}

func (j Job) ConfigureNextAttempt(now time.Time) (Job, bool) {
	j.attempts++

	if j.attempts > j.MaxAttempts {
		return j, false
	}

	j.At = now.Add(time.Duration(math.Pow(float64(j.attempts), 4)+5) * time.Second)

	return j, true
}
