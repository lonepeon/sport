package jobtest

import (
	"encoding/json"
	"fmt"

	job "github.com/lonepeon/golib/job"
)

type JobMatcher struct {
	name     string
	receiver interface{}
	matcher  func(interface{}) bool
}

func NewJobMatcher(name string, receiver interface{}, matcher func(interface{}) bool) JobMatcher {
	return JobMatcher{name: name, receiver: receiver, matcher: matcher}
}

func (m JobMatcher) Matches(arg interface{}) bool {
	j, ok := arg.(job.Job)
	if !ok {
		return false
	}

	if m.name != j.Name {
		return false
	}

	if err := json.Unmarshal(j.EncodedParams(), &m.receiver); err != nil {
		return false
	}

	return m.matcher(m.receiver)
}

func (m JobMatcher) String() string {
	return fmt.Sprintf("job(name=%s, receiver=%#+v)", m.name, m.receiver)
}
