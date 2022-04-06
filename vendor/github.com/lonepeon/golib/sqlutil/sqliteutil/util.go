package sqliteutil

import (
	"strings"

	"github.com/mattn/go-sqlite3"
)

func IsUniqueConstraintError(err error, index string) bool {
	e, ok := err.(sqlite3.Error)
	if !ok {
		return false
	}

	return strings.Contains(e.Error(), "UNIQUE ") && strings.Contains(e.Error(), index)
}
