package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lonepeon/sport/internal/domain"
)

//go:generate go run ../../../vendor/github.com/lonepeon/golib/sqlutil/cmd/sql-migration ./scripts

// SQLite represents a repository interacting with a SQLite database
type SQLite struct {
	DB *sql.DB
}

// NewSQLite initializes a new SQLite repository
func New(db *sql.DB) SQLite {
	return SQLite{
		DB: db,
	}
}

type runningActivity struct {
	ID               string
	RanAt            string
	Duration         string
	Distance         int
	Speed            float64
	GPXPath          string
	MapPath          string
	ShareableMapPath string
}

func (r runningActivity) ToDomain() (domain.RunningActivity, error) {
	timeLayout := "2006-01-02 15:04:05.999999999-07:00"

	var activity domain.RunningActivity

	activity.GPXPath = domain.GPXFilePath(r.GPXPath)
	activity.MapPath = domain.MapFilePath(r.MapPath)
	activity.ShareableMapPath = domain.ShareableMapFilePath(r.ShareableMapPath)

	ranAt, err := time.Parse(timeLayout, r.RanAt)
	if err != nil {
		return domain.RunningActivity{}, fmt.Errorf("can't parse ran at for activity (id=%s): %v", r.ID, err)
	}
	activity.RanAt = ranAt

	slug, err := domain.NewRunnningActivitySlugFromTime(ranAt)
	if err != nil {
		return domain.RunningActivity{}, fmt.Errorf("can't build slug from ranAt for activity (id=%s): %v", r.ID, err)
	}
	activity.Slug = slug

	duration, err := time.ParseDuration(r.Duration)
	if err != nil {
		return domain.RunningActivity{}, fmt.Errorf("can't parse duration for activity (id=%s): %v", r.ID, err)
	}
	activity.Duration = duration

	speed, err := domain.NewSpeedFromKmh(r.Speed)
	if err != nil {
		return domain.RunningActivity{}, fmt.Errorf("can't parse speed for activity (id=%s): %v", r.ID, err)
	}
	activity.Speed = speed

	distance, err := domain.NewDistanceFromMeters(r.Distance)
	if err != nil {
		return domain.RunningActivity{}, fmt.Errorf("can't parse distance for activity (id=%s): %v", r.ID, err)
	}
	activity.Distance = distance

	return activity, nil
}

// GetRunningActivity returns a list of all running activity
func (r SQLite) GetRunningActivity(ctx context.Context, slug domain.RunningActivitySlug) (domain.RunningActivity, error) {
	statement := `
		SELECT id, ran_at, duration, distance, speed, gpx_path, map_path, shareable_map_path
		FROM runs
		WHERE ran_at = ?
		ORDER BY ran_at DESC`

	rows, err := r.DB.QueryContext(ctx, statement, slug.Time())
	if err != nil {
		return domain.RunningActivity{}, fmt.Errorf("can't get running activity: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return domain.RunningActivity{}, domain.ErrCantGetRunningSession
	}

	var activity runningActivity
	err = rows.Scan(&activity.ID, &activity.RanAt, &activity.Duration, &activity.Distance, &activity.Speed, &activity.GPXPath, &activity.MapPath, &activity.ShareableMapPath)
	if err != nil {
		return domain.RunningActivity{}, fmt.Errorf("can't scan activity: %v", err)
	}

	return activity.ToDomain()
}

func (r SQLite) DeleteRunningActivity(ctx context.Context, slug domain.RunningActivitySlug) error {
	statement := `DELETE FROM runs WHERE ran_at = ?`
	rst, err := r.DB.ExecContext(ctx, statement, slug.Time())
	if err != nil {
		return fmt.Errorf("can't delete activity: %v", err)
	}

	if count, _ := rst.RowsAffected(); count == 0 {
		return domain.ErrCantGetRunningSession
	}

	return nil
}

// ListRunningActivities returns a list of all running activities
func (r SQLite) ListRunningActivities(ctx context.Context) ([]domain.RunningActivity, error) {
	statement := `
		SELECT id, ran_at, duration, distance, speed, gpx_path, map_path, shareable_map_path
		FROM runs
		ORDER BY ran_at DESC`

	rows, err := r.DB.QueryContext(ctx, statement)
	if err != nil {
		return nil, fmt.Errorf("can't get running activities: %v", err)
	}
	defer rows.Close()

	var dbActivity runningActivity
	var activities []domain.RunningActivity
	for rows.Next() {
		err := rows.Scan(&dbActivity.ID, &dbActivity.RanAt, &dbActivity.Duration, &dbActivity.Distance, &dbActivity.Speed, &dbActivity.GPXPath, &dbActivity.MapPath, &dbActivity.ShareableMapPath)
		if err != nil {
			return nil, fmt.Errorf("can't scan activity: %v", err)
		}

		activity, err := dbActivity.ToDomain()
		if err != nil {
			return nil, err
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

// RecordRunningActivity persists the activity in database
func (r SQLite) RecordRunningActivity(ctx context.Context, activity domain.RunningActivity) error {
	statement := `INSERT INTO runs (id, ran_at, duration, distance, speed, gpx_path, map_path, shareable_map_path, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.DB.ExecContext(
		ctx,
		statement,
		uuid.NewString(),
		activity.RanAt,
		activity.Duration.String(),
		activity.Distance.Meters(),
		activity.Speed.KilometersPerHour(),
		activity.GPXPath.String(),
		activity.MapPath.String(),
		activity.ShareableMapPath.String(),
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("can't insert into table: %v", err)
	}

	return nil
}
