# Entry, Story, and Storyline are three distinct concepts

An **Entry** is one item from one Feed, immutable once written, one note per Entry. A **Story** is a developing narrative thread that accumulates Entries over its lifetime (open or closed status, human-named slug). A **Storyline** is a v3 weekly-synthesis narrative that references multiple Stories and standalone Entries. These were initially conflated; sharpening them was the foundational step of v1 design.

The Story concept lands in v1 — schema (`stories`, `entry_stories`) plus manual attachment via both `inkwell stories ...` CLI and a `story:` frontmatter field read on each ingest sweep. Automatic attachment (embedding-similarity-based) is v2. Storyline synthesis is v3. Crucially, Story ≠ Cluster: Stories are persistent and human-named; Clusters are ephemeral and algorithm-named, only existing during a synthesis run.

See `CONTEXT.md` for the canonical glossary. This ADR exists so the *reasoning* survives if CONTEXT.md is ever edited.
