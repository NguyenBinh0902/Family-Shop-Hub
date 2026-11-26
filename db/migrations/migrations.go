// Package migrations embeds the SQL migrations so the migrate binary and the
// integration tests run exactly the files in this directory, with no path
// lookups at runtime.
package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var files embed.FS

// FS returns the embedded migration files.
func FS() fs.FS { return files }

// OpenDB opens a database/sql handle (which goose requires) backed by the pgx
// driver. url must use the owner role, which may run DDL.
func OpenDB(url string) (*sql.DB, error) {
	cfg, err := pgx.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("parse migration database url: %w", err)
	}
	return stdlib.OpenDB(*cfg), nil
}

// NewProvider returns a goose provider that applies the embedded migrations
// to db. db must connect as the owner role, which may run DDL.
func NewProvider(db *sql.DB, opts ...goose.ProviderOption) (*goose.Provider, error) {
	p, err := goose.NewProvider(goose.DialectPostgres, db, files, opts...)
	if err != nil {
		return nil, fmt.Errorf("create migration provider: %w", err)
	}
	return p, nil
}

var fileName = regexp.MustCompile(`^(\d+)_.+\.sql$`)

// validName limits migration names to snake_case so file names stay portable.
var validName = regexp.MustCompile(`^[a-z0-9]+(_[a-z0-9]+)*$`)

const template = `-- +goose Up


-- +goose Down

`

// Create writes a new empty migration NNNN_name.sql in dir, numbered one past
// the highest existing version, and returns its path.
func Create(dir, name string) (string, error) {
	if !validName.MatchString(name) {
		return "", fmt.Errorf("migration name %q must be snake_case, e.g. add_products", name)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", dir, err)
	}
	var latest int64
	for _, e := range entries {
		m := fileName.FindStringSubmatch(e.Name())
		if m == nil {
			continue
		}
		v, err := strconv.ParseInt(m[1], 10, 64)
		if err != nil {
			return "", fmt.Errorf("parse version of %s: %w", e.Name(), err)
		}
		latest = max(latest, v)
	}

	path := filepath.Join(dir, fmt.Sprintf("%04d_%s.sql", latest+1, name))
	// O_EXCL: never overwrite an existing migration.
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return "", fmt.Errorf("create %s: %w", path, err)
	}
	if _, err := f.WriteString(template); err != nil {
		_ = f.Close()
		return "", fmt.Errorf("write %s: %w", path, err)
	}
	if err := f.Close(); err != nil {
		return "", fmt.Errorf("close %s: %w", path, err)
	}
	return path, nil
}
