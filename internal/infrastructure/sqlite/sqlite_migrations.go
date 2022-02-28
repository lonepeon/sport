// Code generated ./scripts DO NOT EDIT

package sqlite

import "github.com/lonepeon/golib/sqlutil"

// Migrations returns an ordered list of migrations to execute
func Migrations() []sqlutil.Migration {
	return []sqlutil.Migration{
		{
			Version: "202111212328",
			Script: `CREATE TABLE runs (
  id TEXT PRIMARY KEY,
  ran_at TEXT,
  duration TEXT,
  distance REAL,
  speed REAL,
  gpx_path TEXT,
  map_path TEXT,
  created_at TEXT
)

`,
		},
		{
			Version: "20211220002500",
			Script: `ALTER TABLE runs ADD COLUMN shareable_map_path TEXT;

`,
		},
	}
}
