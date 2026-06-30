package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sampleOPML = `<?xml version="1.0" encoding="UTF-8"?>
<opml version="1.0">
  <body>
    <outline text="Tech" title="Tech">
      <outline type="rss" title="Simon Willison" xmlUrl="https://simonwillison.net/atom/everything/"/>
      <outline type="rss" text="Hacker News" xmlUrl="https://news.ycombinator.com/rss"/>
    </outline>
    <outline type="rss" text="Standalone" xmlUrl="https://example.com/feed"/>
  </body>
</opml>`

func TestFeedsImportIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "cli.db")
	opmlPath := filepath.Join(dir, "subs.opml")
	if err := os.WriteFile(opmlPath, []byte(sampleOPML), 0o644); err != nil {
		t.Fatalf("write opml: %v", err)
	}

	// First import: all three are new.
	if out := run(t, dbPath, "feeds", "import", opmlPath); !strings.Contains(out, "3 new, 0 already present") {
		t.Errorf("first import: %q", out)
	}

	// Re-import the same file: nothing new, all skipped.
	if out := run(t, dbPath, "feeds", "import", opmlPath); !strings.Contains(out, "0 new, 3 already present") {
		t.Errorf("second import should be a no-op: %q", out)
	}

	// The feeds are actually queryable afterwards.
	if out := run(t, dbPath, "feeds", "list"); !strings.Contains(out, "simonwillison.net") {
		t.Errorf("list after import: %q", out)
	}
}
