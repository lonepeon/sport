package sqlutil

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrGeneric = errors.New("something wrong happened")
)

type Migration struct {
	Version string
	Script  string
}

func ExecuteMigrations(ctx context.Context, db *sql.DB, migrations []Migration) ([]string, error) {
	if err := initMigrationTable(ctx, db); err != nil {
		return nil, err
	}

	executedMigrations, err := loadAlreadyExecutedMigrations(ctx, db)
	if err != nil {
		return nil, err
	}

	var newMigrationVersions []string
	for _, migration := range migrations {
		if isAlreadyExecuted(executedMigrations, migration) {
			continue
		}

		newMigrationVersions = append(newMigrationVersions, migration.Version)
		if err := executeNewMigration(ctx, db, migration); err != nil {
			return nil, err
		}
	}

	return newMigrationVersions, nil
}

func initMigrationTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS migrations (version TEXT PRIMARY KEY, executed_at TEXT)`)
	if err != nil {
		return fmt.Errorf("can't create migrations table: %v", err)
	}

	return nil
}

func isAlreadyExecuted(executedMigrations map[string]interface{}, migration Migration) bool {
	_, ok := executedMigrations[migration.Version]
	return ok
}

func loadAlreadyExecutedMigrations(ctx context.Context, db *sql.DB) (map[string]interface{}, error) {
	rows, err := db.QueryContext(ctx, `SELECT version FROM migrations ORDER BY executed_at ASC`)
	if err != nil {
		return nil, fmt.Errorf("can't select already executed migrations: %v", err)
	}
	defer rows.Close()

	executedMigrations := make(map[string]interface{})
	var version string
	for rows.Next() {
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("can't scan past migrations result from sql: %v", err)
		}
		executedMigrations[version] = nil
	}

	return executedMigrations, nil
}

func executeNewMigration(ctx context.Context, db *sql.DB, migration Migration) error {
	if _, err := db.ExecContext(ctx, migration.Script); err != nil {
		return fmt.Errorf("can't execute migration (version=%s): %v", migration.Version, err)
	}

	_, err := db.ExecContext(
		ctx,
		`INSERT INTO migrations (version, executed_at) VALUES (?,?)`,
		migration.Version, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("can't mark migration as executed (version=%s): %v", migration.Version, err)
	}

	return nil
}
