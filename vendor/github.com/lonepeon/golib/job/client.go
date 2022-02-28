package job

import (
	"database/sql"
	"fmt"
)

type Client struct {
	db *sql.DB
}

func (c *Client) Enqueue(job Job) error {
	_, err := c.db.Exec(`
		INSERT INTO jobs (id, name, params, at, attempts, max_attempts)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, job.id, job.Name, job.params, job.At, job.attempts, job.MaxAttempts,
	)

	if err != nil {
		return fmt.Errorf("can't insert job for later use: %w: %v", ErrGeneric, err)
	}

	return nil
}
