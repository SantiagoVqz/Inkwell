package store

import (
	"database/sql"
	"embed"
	"fmt"
	"sort"
	"strings"
)

// migrationsFS bundles the .sql files into the binary at compile time, so the
// shipped single binary carries its own schema — no files to install alongside
// it. The //go:embed directive is read by the compiler, not at runtime.
//
//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migrate applies every embedded migration that has not run yet, in filename
// order, each inside its own transaction. It is idempotent: a migration that is
// already recorded in schema_migrations is skipped, so calling Migrate on an
// up-to-date database is a no-op. This is what `inkwell migrate` will call.
func Migrate(db *sql.DB) error {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version    TEXT PRIMARY KEY,
		applied_at TEXT NOT NULL DEFAULT (datetime('now'))
	)`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names) // 0001_, 0002_, ... apply in lexical (= numeric) order

	for _, name := range names {
		applied, err := isApplied(db, name)
		if err != nil {
			return err
		}
		if applied {
			continue
		}
		if err := applyOne(db, name); err != nil {
			return fmt.Errorf("apply %s: %w", name, err)
		}
	}
	return nil
}

func isApplied(db *sql.DB, version string) (bool, error) {
	var n int
	err := db.QueryRow(
		`SELECT count(*) FROM schema_migrations WHERE version = ?`, version,
	).Scan(&n)
	return n > 0, err
}

func applyOne(db *sql.DB, name string) error {
	body, err := migrationsFS.ReadFile("migrations/" + name)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // no-op after a successful Commit; safety net on early return

	if _, err := tx.Exec(string(body)); err != nil {
		return err
	}
	if _, err := tx.Exec(
		`INSERT INTO schema_migrations (version) VALUES (?)`, name,
	); err != nil {
		return err
	}
	return tx.Commit()
}
