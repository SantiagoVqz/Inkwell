package store

import (
	"context"
	"testing"
)

// TestFeedRoundTrip drives the generated CreateFeed/GetFeed through the typed
// Queries struct — proving the sqlc wiring, the BOOLEAN->bool default, and the
// RETURNING-populated id all work end to end.
func TestFeedRoundTrip(t *testing.T) {
	q := New(newTestDB(t))
	ctx := context.Background()

	created, err := q.CreateFeed(ctx, CreateFeedParams{
		Url:       "https://example.com/feed",
		Title:     "Example",
		CreatedAt: "2026-01-01T00:00:00Z",
	})
	if err != nil {
		t.Fatalf("CreateFeed: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("expected an auto-assigned id, got 0")
	}
	if !created.Active {
		t.Error("expected new feed to default to active=true")
	}

	got, err := q.GetFeed(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetFeed: %v", err)
	}
	if got.Url != created.Url {
		t.Errorf("url mismatch: got %q want %q", got.Url, created.Url)
	}
}

// TestCreateEntryNullableGuid proves the *string mapping: a nil Guid stores SQL
// NULL and round-trips as nil; a non-nil Guid round-trips its value.
func TestCreateEntryNullableGuid(t *testing.T) {
	q := New(newTestDB(t))
	ctx := context.Background()

	feed, err := q.CreateFeed(ctx, CreateFeedParams{
		Url: "https://example.com/feed", CreatedAt: "2026-01-01T00:00:00Z",
	})
	if err != nil {
		t.Fatalf("CreateFeed: %v", err)
	}

	noGuid, err := q.CreateEntry(ctx, CreateEntryParams{
		FeedID:       feed.ID,
		Guid:         nil, // -> SQL NULL; identity falls back to canonical_url
		CanonicalUrl: "https://example.com/a",
		ContentHash:  "hash-a",
		CreatedAt:    "2026-01-01T00:00:00Z",
	})
	if err != nil {
		t.Fatalf("CreateEntry (nil guid): %v", err)
	}
	if noGuid.Guid != nil {
		t.Errorf("expected nil guid round-trip, got %q", *noGuid.Guid)
	}

	g := "guid-123"
	withGuid, err := q.CreateEntry(ctx, CreateEntryParams{
		FeedID:       feed.ID,
		Guid:         &g,
		CanonicalUrl: "https://example.com/b",
		ContentHash:  "hash-b",
		CreatedAt:    "2026-01-01T00:00:00Z",
	})
	if err != nil {
		t.Fatalf("CreateEntry (with guid): %v", err)
	}
	if withGuid.Guid == nil || *withGuid.Guid != g {
		t.Errorf("expected guid %q to round-trip, got %v", g, withGuid.Guid)
	}
}
