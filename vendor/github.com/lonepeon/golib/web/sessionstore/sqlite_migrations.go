// Code generated ./sqlite_scripts DO NOT EDIT

package sessionstore

import "github.com/lonepeon/golib/sqlutil"

// Migrations returns an ordered list of migrations to execute
func Migrations() []sqlutil.Migration {
	return []sqlutil.Migration{
		{
			Version: "20220330205112",
			Script: `CREATE TABLE sessions (
  id TEXT PRIMARY KEY,
  data TEXT NOT NULL,
  created_at TEXT NOT NULL,
  expired_at TEXT NOT NULL
)

`,
		},
	}
}
