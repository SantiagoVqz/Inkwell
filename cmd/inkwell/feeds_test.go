package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

// run builds a fresh command tree (so flag state never leaks between calls),
// points it at dbPath, captures its output, and executes the given args.
func run(t *testing.T, dbPath string, args ...string) string {
	t.Helper()
	out := &bytes.Buffer{}
	root := newRootCmd()
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs(append(args, "--db", dbPath))
	if err := root.Execute(); err != nil {
		t.Fatalf("execute %v: %v\noutput: %s", args, err, out.String())
	}
	return out.String()
}

func TestFeedsLifecycle(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "cli.db")

	// add
	if out := run(t, dbPath, "feeds", "add", "https://example.com/feed", "--title", "Example"); !strings.Contains(out, "added feed 1") {
		t.Errorf("add: unexpected output %q", out)
	}

	// list shows it, active by default
	out := run(t, dbPath, "feeds", "list")
	if !strings.Contains(out, "example.com") || !strings.Contains(out, "true") {
		t.Errorf("list after add: %q", out)
	}

	// deactivate flips the active column
	run(t, dbPath, "feeds", "deactivate", "1")
	if out := run(t, dbPath, "feeds", "list"); !strings.Contains(out, "false") {
		t.Errorf("list after deactivate should show false: %q", out)
	}

	// remove empties the list
	run(t, dbPath, "feeds", "remove", "1")
	if out := run(t, dbPath, "feeds", "list"); !strings.Contains(out, "no feeds yet") {
		t.Errorf("list after remove: %q", out)
	}
}

func TestRemoveRejectsNonNumericID(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "cli.db")
	out := &bytes.Buffer{}
	root := newRootCmd()
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"feeds", "remove", "abc", "--db", dbPath})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error for non-numeric id, got nil")
	}
}
