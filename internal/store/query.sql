-- Named queries for sqlc. Each `-- name: X :kind` comment tells sqlc what Go
-- method to generate and what it returns: :one (single row), :many (slice),
-- :exec (no rows, just error). The `?` placeholders bind in argument order.
--
-- This first slice covers feeds (milestone 3, Feed CRUD) and the two entry
-- queries the fetcher will need (milestone 5). Story/attachment queries arrive
-- with their own milestone.

-- name: CreateFeed :one
INSERT INTO feeds (url, title, created_at)
VALUES (?, ?, ?)
RETURNING *;

-- name: CreateFeedIfNew :execrows
-- Idempotent insert for bulk import: ON CONFLICT does nothing when the url
-- already exists. :execrows returns the rows-affected count (1 = inserted,
-- 0 = already present), so the importer can tally new vs skipped.
INSERT INTO feeds (url, title, created_at)
VALUES (?, ?, ?)
ON CONFLICT (url) DO NOTHING;

-- name: GetFeed :one
SELECT * FROM feeds
WHERE id = ?;

-- name: GetFeedByURL :one
SELECT * FROM feeds
WHERE url = ?;

-- name: ListFeeds :many
SELECT * FROM feeds
ORDER BY id;

-- name: ListActiveFeeds :many
SELECT * FROM feeds
WHERE active = 1
ORDER BY id;

-- name: SetFeedActive :exec
UPDATE feeds
SET active = ?
WHERE id = ?;

-- name: DeleteFeed :exec
DELETE FROM feeds
WHERE id = ?;

-- name: CreateEntry :one
INSERT INTO entries (
    feed_id, guid, canonical_url, title, body, published_at, content_hash, created_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = ?;
