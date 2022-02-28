package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

type migration struct {
	Version string
	Script  string
}

type variables struct {
	CommandLine string
	Package     string
	FuncName    string
	Migrations  []migration
}

func run() error {
	pkg := os.Getenv("GOPACKAGE")
	cmdline := strings.Join(os.Args[1:], " ")
	fName := os.Getenv("GOFILE")
	fName = strings.TrimSuffix(fName, filepath.Ext(fName)) + "_migrations.go"
	funcName := "Migrations"

	flag.StringVar(&fName, "fname", fName, "name of the migration file")
	flag.StringVar(&funcName, "func", funcName, "name of the migration function")
	flag.Parse()

	folder := flag.Arg(0)
	if folder == "" {
		return fmt.Errorf("folder is required as first argument")
	}

	stat, err := os.Stat(folder)
	if err != nil {
		return fmt.Errorf("can't get folder information (folder=%s): %v", folder, err)
	}

	if !stat.IsDir() {
		return fmt.Errorf("folder is not a folder (folder=%s)", folder)
	}

	migrations, err := loadMigrations(folder)
	if err != nil {
		return fmt.Errorf("can't parse migration from folder (folder=%s): %v", folder, err)
	}

	return writeMigrationFile(fName, variables{
		CommandLine: cmdline,
		Package:     pkg,
		FuncName:    funcName,
		Migrations:  migrations,
	})
}

var tpl = `
// Code generated {{ .CommandLine }} DO NOT EDIT

package {{ .Package }}

import "github.com/lonepeon/golib/sqlutil"

// {{ .FuncName }} returns an ordered list of migrations to execute
func {{ .FuncName }}() []sqlutil.Migration {
	return []sqlutil.Migration {
		{{- range .Migrations }}
		{
			Version: "{{ .Version }}",
			Script: ` + "`" + `
{{- .Script }}
` + "`" + `,
		},
		{{- end }}
	}
}
`

func loadMigrations(folder string) ([]migration, error) {
	var migrations []migration
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != folder {
			if path == folder {
				return nil
			}
			return filepath.SkipDir
		}

		if filepath.Ext(info.Name()) != ".sql" {
			return nil
		}

		version := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("can't read sql migration (name=%s): %v", path, err)
		}

		migrations = append(migrations, migration{Version: version, Script: string(content)})

		return nil
	})

	return migrations, err
}

func writeMigrationFile(fName string, vars variables) error {
	var content bytes.Buffer

	err := template.Must(template.New("migrations.go").Parse(tpl)).Execute(&content, vars)
	if err != nil {
		return fmt.Errorf("can't generate migration file: %v", err)
	}

	fContent, err := format.Source(content.Bytes())
	if err != nil {
		return fmt.Errorf("can't format migration file: %v\n%s", err, content.String())
	}

	if err := ioutil.WriteFile(fName, fContent, 0644); err != nil {
		return fmt.Errorf("can't persist migration file (name=%s): %v", fName, err)
	}

	return nil
}
