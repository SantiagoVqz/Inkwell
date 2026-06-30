// Package opml parses OPML subscription exports (Feedly, Inoreader, etc.) into
// a flat list of feeds. OPML is XML with arbitrarily nested <outline> elements:
// a folder is an outline without an xmlUrl, a feed is an outline that has one.
package opml

import (
	"encoding/xml"
	"fmt"
	"io"
)

// Feed is one subscription extracted from an OPML file.
type Feed struct {
	Title string
	URL   string
}

// document mirrors the slice of the OPML structure we care about. encoding/xml
// fills these fields by matching struct tags to element/attribute names, and
// ignores anything we don't tag — so we declare only what we use. The
// `body>outline` path tag jumps straight to the <outline>s under <body>.
type document struct {
	XMLName  xml.Name  `xml:"opml"`
	Outlines []outline `xml:"body>outline"`
}

// outline is recursive: the `xml:"outline"` tag on Children matches nested
// <outline> elements, so one type models both folders and feeds at any depth.
type outline struct {
	Text     string    `xml:"text,attr"`
	Title    string    `xml:"title,attr"`
	XMLURL   string    `xml:"xmlUrl,attr"`
	Children []outline `xml:"outline"`
}

// Parse reads an OPML document and returns every feed it contains, flattening
// nested folders. Duplicate URLs within the file are returned as-is; de-duping
// against the database is the caller's concern (see CreateFeedIfNew).
func Parse(r io.Reader) ([]Feed, error) {
	var doc document
	if err := xml.NewDecoder(r).Decode(&doc); err != nil {
		return nil, fmt.Errorf("parse opml: %w", err)
	}

	var feeds []Feed
	var walk func([]outline)
	walk = func(outlines []outline) {
		for _, o := range outlines {
			if o.XMLURL != "" {
				feeds = append(feeds, Feed{Title: o.displayTitle(), URL: o.XMLURL})
			}
			walk(o.Children) // recurse into folders
		}
	}
	walk(doc.Outlines)
	return feeds, nil
}

// displayTitle prefers the `title` attribute, falling back to `text` — readers
// disagree on which they populate, and some set only one.
func (o outline) displayTitle() string {
	if o.Title != "" {
		return o.Title
	}
	return o.Text
}
