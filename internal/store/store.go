// Package store owns the SQLite database: opening it with the right pragmas,
// applying migrations, and (soon) the sqlc-generated typed queries.
package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // pure-Go driver, registers itself as "sqlite" (ADR-0013)
)

// Open opens the SQLite database at path, sets the pragmas the rest of the app
// assumes, and brings the schema up to date. The pragmas are per-connection, so
// they live in the DSN where the driver applies them to every pooled connection:
//   - foreign_keys(1): FK enforcement is OFF by default in SQLite; without this
//     the ON DELETE CASCADE clauses in the schema would be decorative.
//   - journal_mode(WAL): readers (the dashboard) don't block the writer (ingest).
//   - busy_timeout(5000): wait up to 5s for a lock instead of failing instantly.
func Open(path string) (*sql.DB, error) {
	// The default DB path lives under ~/.local/share/inkwell/, which won't exist
	// on a fresh machine. The driver creates the file but not its parent dirs.
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create db dir %s: %w", dir, err)
		}
	}

	dsn := fmt.Sprintf(
		"file:%s?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)",
		path,
	)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db %s: %w", path, err)
	}

	if err := Migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate db %s: %w", path, err)
	}

	return db, nil
}
