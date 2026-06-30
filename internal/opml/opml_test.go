package opml

import (
	"strings"
	"testing"
)

// sample mixes a folder with two feeds and a top-level standalone feed, and
// exercises the title/text fallback: the first feed sets `title`, the second
// only `text`.
const sample = `<?xml version="1.0" encoding="UTF-8"?>
<opml version="1.0">
  <head><title>subscriptions</title></head>
  <body>
    <outline text="Tech" title="Tech">
      <outline type="rss" title="Simon Willison" text="ignored when title set"
               xmlUrl="https://simonwillison.net/atom/everything/" htmlUrl="https://simonwillison.net/"/>
      <outline type="rss" text="Hacker News" xmlUrl="https://news.ycombinator.com/rss"/>
    </outline>
    <outline type="rss" text="Standalone" xmlUrl="https://example.com/feed"/>
  </body>
</opml>`

func TestParseFlattensAndResolvesTitles(t *testing.T) {
	feeds, err := Parse(strings.NewReader(sample))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(feeds) != 3 {
		t.Fatalf("expected 3 feeds (folder flattened + standalone), got %d: %+v", len(feeds), feeds)
	}

	want := []Feed{
		{Title: "Simon Willison", URL: "https://simonwillison.net/atom/everything/"}, // title attr wins
		{Title: "Hacker News", URL: "https://news.ycombinator.com/rss"},              // text fallback
		{Title: "Standalone", URL: "https://example.com/feed"},
	}
	for i, w := range want {
		if feeds[i] != w {
			t.Errorf("feed[%d] = %+v, want %+v", i, feeds[i], w)
		}
	}
}

// A bare folder with no feed outlines should yield nothing, not an error.
func TestParseFolderWithoutFeeds(t *testing.T) {
	doc := `<opml><body><outline text="Empty Folder"></outline></body></opml>`
	feeds, err := Parse(strings.NewReader(doc))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(feeds) != 0 {
		t.Errorf("expected 0 feeds, got %d", len(feeds))
	}
}

func TestParseRejectsMalformedXML(t *testing.T) {
	if _, err := Parse(strings.NewReader("<opml><body>")); err == nil {
		t.Fatal("expected error on truncated XML, got nil")
	}
}
