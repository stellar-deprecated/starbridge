package store

import (
	"database/sql"
	"embed"
	"strings"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/stellar/go/support/db"
)

//go:embed migrations/*.sql
var migrations embed.FS

type DB struct {
	Session db.SessionInterface
}

func InitSchema(db *sql.DB) error {
	return Migrate(db, migrate.Up, 0)
}

func Migrate(db *sql.DB, dir migrate.MigrationDirection, max int) error {
	m := &migrate.AssetMigrationSource{
		Asset: migrations.ReadFile,
		AssetDir: func() func(string) ([]string, error) {
			return func(path string) ([]string, error) {
				dirEntry, err := migrations.ReadDir(path)
				if err != nil {
					return nil, err
				}
				entries := make([]string, 0)
				for _, e := range dirEntry {
					entries = append(entries, e.Name())
				}

				return entries, nil
			}
		}(),
		Dir: "migrations",
	}
	_, err := migrate.ExecMax(db, "postgres", m, dir, max)
	return err
}

func IsDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}
