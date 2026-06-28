# Inkwell roadmap

Snippet-driven build. Each milestone ships something demonstrable, teaches specific Go concepts, and is sized for one or two evening sessions. See `CONTEXT.md` for vocabulary and `docs/adr/` for durable decisions.

> **Strategic framing (ADR-0008):** Inkwell is the shippable product that absorbs the Python `research-tool`'s full scope, slice by slice. v1's ingest-only scope below is unchanged — only the destination is now explicit. Distribution (ADR-0009), the web dashboard (ADR-0010), and direct API access (ADR-0011) are first-class, not "someday." The full feature inventory and sprint backlog live in [`docs/research-tool-feature-map.md`](docs/research-tool-feature-map.md).

## v1 — Ingest only (target: ~4-6 weeks elapsed, ~13-18 active evenings)

Goal: replace the Python `content-pipeline` cron job with `inkwell ingest`. No embeddings, no LLM, no clustering. Entries land in `Inbox/Inkwell/YYYY/MM/{slug}-{hash}.md`. Stories are manually attached via CLI or frontmatter.

| # | Milestone | Output you can show | New Go concepts | Effort |
|---|---|---|---|---|
| 0 | Bootstrap | `inkwell version` prints | Modules, build, `main`, `go.mod`, Makefile | 1 evening |
| 1 | Config + CLI scaffold | `inkwell config show` reads YAML with defaults | `internal/`, cobra, slog, struct tags, `errors.Is` | 1-2 evenings |
| 2 | Store layer | `inkwell migrate` creates schema; store tests pass | sqlc workflow, in-memory SQLite tests, generated code | 2-3 evenings |
| 3 | Feed CRUD | `inkwell feeds add/list/remove/activate/update` works | Cobra subcommand patterns, sqlc query bindings | 1 evening |
| 4 | OPML import/export | `inkwell feeds import feedly.opml` adds 30 feeds | `encoding/xml`, struct unmarshalling | 1 evening |
| 5 | Fetcher (worker pool) | `inkwell ingest --dry-run` fetches in parallel, prints planned inserts | **The inflection point:** `context.Context` propagation, `errgroup`, semaphore channel, `defer` for body close, error wrapping at boundaries | 2-3 evenings |
| 6 | Idempotency + failure policy | Re-running `ingest` doesn't duplicate; broken feeds get 24h backoff after 3 failures | Composite unique constraints, transactional inserts, table-driven tests | 1-2 evenings |
| 7 | Obsidian writer | Real notes appear in the vault under `YYYY/MM/` | `text/template` or yaml.v3 encode, `os.MkdirAll`, atomic file writes (`renameio`) | 1-2 evenings |
| 8 | Story attachment | `inkwell stories new/attach/close/list/show`; editing `story:` in Obsidian picked up on next sweep | More sqlc, frontmatter parsing on read side | 1-2 evenings |
| 9 | launchd integration + status | launchd runs ingest on schedule; `inkwell status` shows feed health and recent runs | OS integration via plist; nothing new in Go itself | 1 evening |
| 10 | Polish | 1-2 weeks of real daily use; fix surfaced edge cases | — | 1 week elapsed (low-active) |

### v1 acceptance criteria

- `inkwell ingest` invoked by launchd writes notes to the vault for 14 consecutive days with zero manual intervention.
- The Python `content-pipeline` can be disabled without losing any feeds.
- `inkwell status` accurately reports which feeds are healthy and which are erroring.
- All sqlc queries compile; `go test ./...` passes; `go vet` clean.

### v1 explicitly does NOT ship

- Embeddings of any kind (`internal/embed/`) → v2
- Clustering (`internal/cluster/`) → v3
- LLM synthesis (`internal/synth/`) → v3
- Storyline notes → v3
- Automatic Story attachment → v2
- Internal daemon mode (`run --daemon`) → if and only if v1 is migrated off macOS
- TUI (`bubbletea`) → optional *secondary* surface, later; the primary UI is the web dashboard (ADR-0010)
- Web dashboard → product milestone *after* the ingest walking skeleton is solid (ADR-0010), not v1
- Semantic dedup → v2
- `goreleaser` / single-binary release → first public release; distribution is first-class (ADR-0009), designed in early, just not coded in v1's first milestones

---

## v2 — Embeddings layer (outline, not a commitment)

Goal: every Entry has a vector; "find related" works; auto-attach to Stories.

| # | Milestone |
|---|---|
| v2-1 | `internal/embed/` interface + Ollama impl (first cut: `nomic-embed-text`) |
| v2-2 | `embedding BLOB` + `embedding_model TEXT` columns; backfill query plan |
| v2-3 | Embed during ingest (one call per new Entry); `inkwell entries reembed --model X` for model upgrades |
| v2-4 | Brute-force cosine similarity helpers; `inkwell entries similar <id>` shows neighbours |
| v2-5 | Auto-attach new Entries to open Stories via similarity threshold |
| v2-6 | Semantic dedup: warn / mark / merge syndicated stories |

## v3 — Synthesis (outline)

Goal: weekly Storyline note in the vault, clustered themes, LLM-generated narrative.

| # | Milestone |
|---|---|
| v3-1 | `internal/cluster/` interface + agglomerative-threshold impl |
| v3-2 | `internal/synth/` interface + **Anthropic Go SDK** impl (`anthropic-sdk-go`, BYO-key, prompt template, parse output) — ADR-0011 supersedes the earlier `claude` CLI shell-out plan |
| v3-3 | Synthesis pipeline orchestration; Storyline note layout in vault |
| v3-4 | Optional refinements: HDBSCAN swap for the clusterer; richer dashboard parity (timeline/heatmap). (Distribution, the web dashboard, and the Anthropic SDK are no longer parked here — they're first-class per ADR-0009/0010/0011 and land on their own tracks.) |

---

## v4 — Output add-on (optional, ADR-0012)

Goal: optional content-generation output, behind a plugin interface, for users who want drafts (Santiago's LinkedIn/blog use case). **Strictly optional and last — must not pull core sprints; v1–v3 must be fully usable without it.** Port the drafting logic from research-tool (`draft.py`, voice lint/tighten, angles, claim lint); prompts in `research-tool/prompts/formats/` carry over verbatim. See [`docs/research-tool-feature-map.md`](docs/research-tool-feature-map.md) §F.

| # | Milestone |
|---|---|
| v4-1 | `internal/output/` plugin interface (one consumer to start); behind build tag / config flag |
| v4-2 | Drafting plugin: multi-format publish bundle (LinkedIn / blog / Substack / threads / newsletter / notes) |
| v4-3 | Voice lint + tighten; content angles (`essays` / `cityfront`); claim lint |

---

## Working mode

- Hybrid snippet-driven build. One function / type / test at a time, with Go-idiom explainers the first time a concept appears.
- Skip language-agnostic explanations — closures, generics, modules, testing are familiar from TS/Python.
- Treat the conversation as pair-programming; the assistant proposes the next snippet, explains why, waits for the user to push back or absorb before continuing.
