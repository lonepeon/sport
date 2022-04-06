// Code generated ./sqlite_scripts DO NOT EDIT

package authenticationstore

import "github.com/lonepeon/golib/sqlutil"

// Migrations returns an ordered list of migrations to execute
func Migrations() []sqlutil.Migration {
	return []sqlutil.Migration{
		{
			Version: "20220402232201",
			Script: `CREATE TABLE authentication_user (
  ID TEXT PRIMARY KEY,
  username TEXT,
  password TEXT,
  salt TEXT
);

CREATE UNIQUE INDEX authentication_user_username ON authentication_user(username);

`,
		},
	}
}
