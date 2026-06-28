# Vault layout: Year/Month folders, one note per Entry

Inkwell writes each Entry as a separate markdown file at `{vault}/{subfolder}/YYYY/MM/{slug}-{shorthash}.md`. The slug is derived from the title; the shorthash (first 4 hex chars of content_hash) breaks collisions. The date lives in the folder path and in frontmatter — not in the filename — so filenames stay human-readable.

Feed name is *not* in the path; it's a frontmatter field (`source: TechCrunch`) and a tag. Filtering by feed happens via Obsidian Properties search or ripgrep, not by folder navigation. This avoids folder explosion as the feed list grows and keeps the same story syndicated to multiple feeds organized by time, which is the dominant access pattern.

**Considered and rejected**: flat folder (becomes 1000+ files in a year, hostile to scrolling); folder-per-feed (scatters the same story across folders, explodes with feed count); daily-aggregate note (forces Inkwell to mutate existing notes, breaks Obsidian linking model — Entry mutability is rejected on principle).
