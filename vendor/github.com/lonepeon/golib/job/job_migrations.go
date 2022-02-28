// Code generated ./scripts DO NOT EDIT

package job

import "github.com/lonepeon/golib/sqlutil"

// Migrations returns an ordered list of migrations to execute
func Migrations() []sqlutil.Migration {
	return []sqlutil.Migration{
		{
			Version: "202111192159",
			Script: `CREATE TABLE jobs (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  params TEXT NOT NULL,
  at TEXT NOT NULL,
  attempts INTEGER NOT NULL,
  max_attempts INTEGER  NOT NULL,
  locked_until TEXT,
  failed TEXT
)

`,
		},
	}
}
