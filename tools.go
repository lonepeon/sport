//go:build tools

package tools

import (
	// generates a String method for the custom enums
	_ "golang.org/x/tools/cmd/stringer"
	// performs static analysis on the code base
	_ "honnef.co/go/tools/cmd/staticcheck"
	// performs a bunch of linting checks on the code base
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	// generates Mock used in tests
	_ "github.com/golang/mock/mockgen"
	// generates SQL migration files
	_ "github.com/lonepeon/golib/sqlutil/cmd/sql-migration"
)
