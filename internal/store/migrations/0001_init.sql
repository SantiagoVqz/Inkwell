-- 0001_init: the v1 ontology — feeds, entries, stories, attachments.
-- Timestamps are RFC3339 UTC TEXT. The application owns the format (no DB
-- default), so "what's stored" always matches what Go's time.RFC3339 emits.

-- feeds: a subscribed RSS/Atom source. One row per Feed.
CREATE TABLE feeds (
    id                   INTEGER PRIMARY KEY,
    url                  TEXT    NOT NULL UNIQUE,
    title                TEXT    NOT NULL DEFAULT '',
    active               BOOLEAN NOT NULL DEFAULT 1,
    created_at           TEXT    NOT NULL,

    -- failure policy (milestone 6): 24h backoff after 3 consecutive failures
    last_fetched_at      TEXT,                          -- NULL until first fetch
    consecutive_failures INTEGER NOT NULL DEFAULT 0,
    last_error           TEXT    NOT NULL DEFAULT ''
);

-- entries: one immutable item from one Feed; exactly one vault note per row.
CREATE TABLE entries (
    id             INTEGER PRIMARY KEY,
    feed_id        INTEGER NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    guid           TEXT,                                -- nullable: ~25% of feeds lack a usable GUID (ADR-0005)
    canonical_url  TEXT    NOT NULL,                    -- identity fallback; aggressively canonicalized
    title          TEXT    NOT NULL DEFAULT '',
    body           TEXT    NOT NULL DEFAULT '',         -- feed body as-is -> markdown (ADR-0007)
    published_at   TEXT,                                -- from the feed; nullable (feeds omit/lie)
    content_hash   TEXT    NOT NULL,                    -- sha256(title+body); diagnostic only (ADR-0005)
    note_path      TEXT    NOT NULL DEFAULT '',         -- vault-relative path once the note is written (M7)
    created_at     TEXT    NOT NULL
);

-- ADR-0005 identity key. Must be an expression index, NOT an inline UNIQUE,
-- because COALESCE() is illegal in a table-level UNIQUE clause. canonical_url
-- is NOT NULL, so the expression is never NULL and uniqueness always bites.
CREATE UNIQUE INDEX idx_entries_identity
    ON entries (feed_id, COALESCE(guid, canonical_url));

-- stories: a persistent, human-named narrative thread.
CREATE TABLE stories (
    id          INTEGER PRIMARY KEY,
    name        TEXT    NOT NULL UNIQUE,
    status      TEXT    NOT NULL DEFAULT 'open'
                   CHECK (status IN ('open', 'closed')),
    created_at  TEXT    NOT NULL,
    closed_at   TEXT                                    -- NULL while open
);

-- attachments: the Entry<->Story link, recording HOW it was made.
-- entry_id is the PRIMARY KEY -> enforces the "Entry 0..1 Story" rule for free.
CREATE TABLE attachments (
    entry_id    INTEGER PRIMARY KEY REFERENCES entries(id) ON DELETE CASCADE,
    story_id    INTEGER NOT NULL    REFERENCES stories(id) ON DELETE CASCADE,
    source      TEXT    NOT NULL
                   CHECK (source IN ('frontmatter', 'cli', 'embedding', 'llm')),
    created_at  TEXT    NOT NULL
);
