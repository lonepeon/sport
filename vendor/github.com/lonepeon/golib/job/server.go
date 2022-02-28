package job

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lonepeon/golib/logger"
)

var (
	ErrGeneric = errors.New("something wrong happened")
)

type Server struct {
	registry *Registry
	db       *sql.DB
	log      *logger.Logger
	shutdown chan bool

	SleepDuration time.Duration
}

func NewServer(db *sql.DB, reg *Registry, log *logger.Logger) *Server {
	return &Server{
		db:            db,
		log:           log,
		registry:      reg,
		shutdown:      make(chan bool),
		SleepDuration: 5 * time.Second,
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.shutdown <- true
	return nil
}

func (s *Server) ListenAndServe() error {
	for {
		s.dequeue()

		select {
		case <-s.shutdown:
			return nil
		case <-time.After(s.SleepDuration):
		}
	}
}

func (c *Server) Client() *Client {
	return &Client{db: c.db}
}

func (s *Server) dequeue() {
	now := time.Now()

	job, err := s.fetchNextJob(now)
	if err != nil {
		return
	}

	handler, err := s.fetchJobHandler(now, job)
	if err != nil {
		return
	}

	if err = s.executeJobHandler(now, handler, job); err != nil {
		return
	}
}

func (s *Server) fetchNextJob(now time.Time) (Job, error) {
	row := s.db.QueryRow(`
			UPDATE jobs
			SET locked_until = $1
			WHERE id = (
				SELECT id
				FROM jobs
				WHERE (locked_until IS NULL OR locked_until <= $2)
					AND attempts < max_attempts
					AND at <= $2
					AND failed IS NULL
				ORDER BY at ASC
				LIMIT 1)
			RETURNING id, name, params, attempts, max_attempts`, now.Add(1*time.Minute), now)

	var job Job
	if err := row.Scan(&job.id, &job.Name, &job.params, &job.attempts, &job.MaxAttempts); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Job{}, err
		}

		s.log.Error(fmt.Sprintf("can't fetch/parse a job: %v", err))
		return Job{}, err
	}

	return job, nil
}

func (s *Server) fetchJobHandler(now time.Time, job Job) (HandlerFunc, error) {
	handler, ok := s.registry.Handler(job.Name)
	if !ok {
		s.log.Error(fmt.Sprintf("can't find registered handler for job (id=%s, name=%s, params=%#+v)", job.id, job.Name, string(job.params)))
		if _, err := s.db.Exec(`UPDATE jobs SET failed = $1 WHERE id = $2`, now, job.id); err != nil {
			s.log.Error(fmt.Sprintf("can't mark job as failed (id=%s, name=%s, params=%#+v): %v", job.id, job.Name, string(job.params), err))
		}
		return nil, fmt.Errorf("handler not found")
	}

	return handler, nil
}

func (s *Server) executeJobHandler(now time.Time, handler HandlerFunc, job Job) error {
	log := s.log.WithFields(logger.String("request-id", job.id))
	log.Info(fmt.Sprintf("executing job handler (id=%s, name=%s, params=%#+v)", job.id, job.Name, string(job.params)))
	if err := handler(context.Background(), job.params); err != nil {
		next, ok := job.ConfigureNextAttempt(time.Now())
		log.Error(fmt.Sprintf("failed to execute job handler (id=%s, name=%s, params=%#+v): %v", next.id, next.Name, string(next.params), err))
		if !ok {
			if _, err := s.db.Exec(`UPDATE jobs SET attempts = $1, failed = $2, locked_until = NULL WHERE id = $2`, next.attempts, now, next.id); err != nil {
				log.Error(fmt.Sprintf("can't mark job as failed (id=%s, name=%s, params=%#+v): %v", next.id, next.Name, string(next.params), err))
			}
			return fmt.Errorf("handler failed with no remaining attempts")
		}

		if _, err := s.db.Exec(`UPDATE jobs SET attempts = $1, at = $2, locked_until = NULL WHERE id = $3`, next.attempts, next.At, next.id); err != nil {
			log.Error(fmt.Sprintf("can't reschedule next attempt (id=%s, name=%s, params=%#+v): %v", next.id, next.Name, string(next.params), err))
		}

		return fmt.Errorf("handler failed and will be retried")
	}

	log.Info(fmt.Sprintf("job successfully processed (id=%s, name=%s, params=%#+v)", job.id, job.Name, string(job.params)))
	if _, err := s.db.Exec(`DELETE FROM jobs WHERE id = $1`, job.id); err != nil {
		log.Error(fmt.Sprintf("can't delete job after successful attempt (id=%s, name=%s, params=%#+v): %v", job.id, job.Name, string(job.params), err))
	}

	return nil
}
