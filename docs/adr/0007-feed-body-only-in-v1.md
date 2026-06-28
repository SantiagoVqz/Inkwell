# v1 ingests feed body as-is; full-article fetch is a v2 concern

Inkwell v1 takes whatever HTML is in the RSS body (`<content:encoded>` or `<description>`), converts it to markdown via `html-to-markdown/v2`, and writes that into the note. No second HTTP request to the article URL, no readability extraction. Excerpt feeds (most major publishers) produce excerpt notes; full-text feeds (most tech blogs) produce full notes — Inkwell does not normalize between them.

Excerpt notes are still useful as a "what was published" index — the title, source, date, and URL are there, and the reader clicks through when interested. Full-article fetching pairs naturally with embeddings in v2, where short bodies actually become a problem (weak embedding signal). Doing it in v1 doubles HTTP requests for marginal value, adds a fragile dependency (`go-readability` works on ~70% of sites; paywalls and Cloudflare break the rest), and forces a fallback policy v1 doesn't need.

**Considered and rejected**: always-fetch-full-HTML (doubles requests, brittle); short-body heuristic (fragile threshold, same complexity as always-fetch); per-feed flag (deferred — add only if needed in v2).
