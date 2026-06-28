# Idempotency: composite unique key on (feed_id, guid OR canonical_url); content_hash is diagnostic only

Entries are uniquely identified by `(feed_id, COALESCE(guid, canonical_url))`. URLs are canonicalized aggressively before storage — stripping `utm_*`, `fbclid`, and similar tracking params; lowercasing the host; stripping trailing slashes. A `content_hash` column (sha256 of title + body) is stored alongside but **does not participate in uniqueness** — it's purely diagnostic.

If the same canonical URL re-appears with a different content_hash, Inkwell logs a warning ("upstream edit detected") and **does not rewrite the note**. Entries are immutable per ADR-0003. This is the only correct policy when notes may have been annotated by the user — overwriting would silently destroy their additions.

**Considered and rejected**: GUID-only (≈25% of feeds in practice publish unusable GUIDs); URL-only without canonicalization (tracking params create flood of false duplicates); content_hash as part of the unique constraint (every typo-fix on the publisher side would create a duplicate Entry).
