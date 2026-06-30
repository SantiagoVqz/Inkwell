package store

import (
	"database/sql"
	"path/filepath"
	"testing"
)

// newTestDB opens a fresh, already-migrated database backed by a temp file that
// the test framework deletes on cleanup. We use a temp file rather than
// `:memory:` deliberately: database/sql pools connections, and every new
// connection to ":memory:" gets its own empty database — so a table created on
// one connection vanishes on the next query. A temp file sidesteps that entirely.
func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func mustExec(t *testing.T, db *sql.DB, query string) {
	t.Helper()
	if _, err := db.Exec(query); err != nil {
		t.Fatalf("exec %q: %v", query, err)
	}
}

func TestMigrateCreatesOntologyTables(t *testing.T) {
	db := newTestDB(t)

	for _, table := range []string{"feeds", "entries", "stories", "attachments"} {
		var name string
		err := db.QueryRow(
			`SELECT name FROM sqlite_master WHERE type='table' AND name=?`, table,
		).Scan(&name)
		if err != nil {
			t.Errorf("table %q missing after migrate: %v", table, err)
		}
	}
}

func TestMigrateIsIdempotent(t *testing.T) {
	db := newTestDB(t)
	// Open already migrated once; a second run must be a clean no-op.
	if err := Migrate(db); err != nil {
		t.Fatalf("second migrate failed: %v", err)
	}
}

// TestEntryIdentityIsUnique exercises the ADR-0005 idempotency key: re-inserting
// the same (feed_id, canonical_url) must be rejected by idx_entries_identity.
func TestEntryIdentityIsUnique(t *testing.T) {
	db := newTestDB(t)
	mustExec(t, db,
		`INSERT INTO feeds (url, created_at) VALUES ('https://example.com/feed', '2026-01-01T00:00:00Z')`)

	insert := `INSERT INTO entries (feed_id, canonical_url, content_hash, created_at)
	           VALUES (1, 'https://example.com/a', 'hash', '2026-01-01T00:00:00Z')`
	mustExec(t, db, insert)

	if _, err := db.Exec(insert); err == nil {
		t.Fatal("expected UNIQUE violation on duplicate (feed_id, canonical_url), got nil")
	}
}

// TestForeignKeysEnforced confirms the foreign_keys pragma is actually on:
// an entry pointing at a non-existent feed must be rejected.
func TestForeignKeysEnforced(t *testing.T) {
	db := newTestDB(t)
	_, err := db.Exec(
		`INSERT INTO entries (feed_id, canonical_url, content_hash, created_at)
		 VALUES (999, 'https://example.com/x', 'hash', '2026-01-01T00:00:00Z')`)
	if err == nil {
		t.Fatal("expected FK violation for entry referencing missing feed, got nil")
	}
}
