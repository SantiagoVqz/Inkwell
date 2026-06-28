# research-tool → Inkwell feature map

A snapshot of everything the Python `research-tool` does today, mapped onto where
it lands in Inkwell. This is the **migration backlog**: each row is a unit of work
you can pull into a sprint. Read it alongside [`ROADMAP.md`](../ROADMAP.md) (the
Go-side milestones) and [`docs/adr/`](adr/) (the *why*).

> **Source of truth for behavior is the Python code.** When you port a row, open
> the listed `research-tool` script/module and treat it as the spec. Prompts in
> `research-tool/prompts/*.md` carry over verbatim. `research-tool` lives at
> `~/Documents/Personal/research-tool`.

## How to read this

- **Target** — the Inkwell version + milestone (per ROADMAP) where the feature lands, and any governing ADR.
  - **v1** = ingest-only (no LLM, no embeddings). **v2** = embeddings layer. **v3** = synthesis (clustering + LLM). **product** = distribution/UX layer that research-tool never had.
- **Port** — complexity of moving it: **Trivial / Low / Med / High**.
- **Fit** — how it ports:
  - **Direct** — same behavior, new language.
  - **Re-architected** — Inkwell does the same *job* a structurally different way (see divergences below).
  - **New** — no research-tool equivalent; net-new product work.
  - **Add-on (v4)** — content-pipeline-specific; not part of core. Ships as an optional v4 output plugin (ADR-0012), off by default for generic users.

---

## Three architectural divergences (this is NOT a 1:1 port)

Before the table — three places where Inkwell deliberately does the same job differently. These are the rows most likely to trip you up if you assume a straight port.

1. **Per-Entry LLM categorization → embeddings + clustering.** research-tool's `categorize.py` calls an LLM on *every* item to assign a taxonomy category + tier. Inkwell rejects this (ADR-0002): embeddings do the cheap deterministic work (dedup, similarity, retrieval) in v2, clustering discovers themes in v3, and the LLM runs **once per cluster**, not per Entry — ~100× fewer calls. So `categorize.py` does **not** port directly; its job is split across v2 (embed) and v3 (cluster + synthesize).

2. **`bundle.py` (LLM story clustering) → manual Stories (v1) + auto-attach (v2) + cluster discovery (v3).** research-tool clusters items into Story arcs with LLM judgment in one step. Inkwell splits this along its grain: Stories are a v1 schema with *manual* attachment (CLI + `story:` frontmatter), embedding-similarity auto-attach is v2, and ephemeral algorithmic Clusters that can promote into Stories are v3. Remember **Story ≠ Cluster** (ADR-0003).

3. **The content-drafting layer is an optional v4 add-on.** `draft.py`, `voice_lint.py`, `voice_tighten.py`, the `essays`/`cityfront` angles, and the LinkedIn/blog/Substack/threads formats are **Santiago-specific content generation**, not generic research tooling. research-tool's own `GENERALIZE.md` flags these as the heavy domain-specific lift. Inkwell core stops at synthesis + recommendation; drafting ships **last, as an optional v4 plugin behind an output-plugin interface** (ADR-0012) — off by default for generic users. It's a real track, but it must not pull core (v1–v3) sprints.

---

## A. Ingest — fetch sources → write notes  (Inkwell v1, the walking skeleton)

| Feature | research-tool source | Target | Port | Fit |
|---|---|---|---|---|
| RSS/Atom fetch from a source list | `scripts/ingest.py`, `feeds.yaml` | v1 · M5 fetcher (worker pool) | Med | Direct |
| Parallel fetch with concurrency cap | (feedparser, serial) | v1 · M5 — `errgroup` + semaphore | Med | Re-architected (Go concurrency is the upgrade) |
| Dedup so re-runs don't duplicate | URL-hash dedup in `ingest.py` | v1 · M6 — composite key, ADR-0005 | Med | Re-architected (canonical-URL + GUID key) |
| Per-feed failure backoff | (none — best-effort) | v1 · M6 — 24h backoff after 3 fails | Low | New |
| Write one markdown note per item | `ingest.py` → `_inbox/` | v1 · M7 writer, ADR-0006 (`YYYY/MM/`) | Med | Re-architected (vault layout differs) |
| Frontmatter on write | `src/research_pipeline/frontmatter.py` | v1 · M7 — `yaml.v3` encode | Low | Direct |
| Atomic / safe writes (never clobber) | `safe_write()` guardrail | v1 · M7 — `renameio`, immutable Entries | Low | Direct (ADR-0005: never rewrite a note) |
| Feed list management | `feeds.yaml` (hand-edited) | v1 · M3 — `inkwell feeds add/list/remove` | Low | New (CLI CRUD over SQLite) |
| OPML import/export | (none) | v1 · M4 — `feeds import x.opml` | Low | New |
| Body taken from RSS as-is (no refetch) | enrich does refetch; v1 Inkwell does not | v1 · ADR-0007 (full fetch → v2) | Trivial | Direct |

## B. Enrich — full article bodies + entities  (Inkwell v2)

| Feature | research-tool source | Target | Port | Fit |
|---|---|---|---|---|
| Full-article body extraction | `scripts/enrich.py` (trafilatura) | v2 — pairs with embeddings, ADR-0007 | High | Direct (readability is brittle; ~70% sites) |
| robots.txt honor + per-domain rate-limit | `enrich.py` | v2 | Med | Direct |
| Skip Twitter/YouTube/PDF | `enrich.py` | v2 | Low | Direct |
| Entity recognition / retag | `scripts/retag_entities.py`, `entities/` | v2+ | High | Re-architected (GENERALIZE.md "enrich" target) |

## C. Organize — classify + structure  (split: v2 embed / v3 LLM)

| Feature | research-tool source | Target | Port | Fit |
|---|---|---|---|---|
| Taxonomy config (categories, tiers, tags) | `config/taxonomy.yaml`, `taxonomy.py` | v2/v3 — config loader | Low | Direct (the config; not the LLM call) |
| Per-item category + tier assignment | `scripts/categorize.py` (LLM/item) | **v2 embed + v3 cluster** | High | **Re-architected — divergence #1** |
| MOC (map-of-content) note updates | `categorize.py` writes `topics/<cat>.md` | v3 — vault writer | Med | Re-architected |
| Tag namespaces / nested tags | `taxonomy.yaml`, migrations | v2 | Med | Direct |
| Tag community detection | `scripts/tag_communities.py` | v3 / dashboard analytics | Med | Re-architected |

## D. Stories — narrative arcs  (v1 manual → v2 auto → v3 discovery)

| Feature | research-tool source | Target | Port | Fit |
|---|---|---|---|---|
| Story schema + lifecycle (open/closed) | implicit in `stories/` notes | v1 · M8 — `stories` tables, ADR-0003 | Med | Re-architected (real schema vs notes) |
| Manual attach (CLI) | (frontmatter only today) | v1 · M8 — `inkwell stories attach` | Low | New |
| Manual attach (`story:` frontmatter sweep) | vault frontmatter | v1 · M8 — read-side parse | Med | Direct |
| Cluster items into arcs (LLM judgment) | `scripts/bundle.py` | **v2 auto-attach + v3 clusters** | High | **Re-architected — divergence #2** |
| Lookback window for clustering | `bundle.py --lookback-days` | v3 — synthesis window | Low | Re-architected |

## E. Synthesize — rollups + daily synthesis  (Inkwell v3)

| Feature | research-tool source | Target | Port | Fit |
|---|---|---|---|---|
| Weekly insight rollups | `scripts/insights.py`, `prompts/insights_weekly.md` | v3 — Storyline, ADR-0003 | High | Re-architected (Storyline concept) |
| Idempotent "nothing changed, skip" | `insights.py` early-exit | v3 | Low | Direct |
| Daily digest synthesis | `scripts/digest.py`, `prompts/digest.md` | v3 | High | Direct (prompt carries over) |
| Quiet-day gate (skip on low volume) | inline in `daily.py` | v3 — synthesis pipeline | Low | Direct |
| Long-running themes | `insights/<theme>.md` | v3 | Med | Re-architected |
| Synthesis orchestration | `daily.py` + `internal/synth` | v3 · v3-3 | High | Re-architected (two-pipeline, ADR-0002) |
| LLM access for synthesis | `claude` CLI shell-out | v3 · **ADR-0011** (Go SDK, BYO-key) | Med | Re-architected (SDK, not CLI) |

## F. Recommend + Output  (recommend = v3 core; drafting = optional v4 add-on)

| Feature | research-tool source | Target | Port | Fit |
|---|---|---|---|---|
| "Top pick + reasoning" recommendation | `scripts/recommend.py`, `prompts/recommend.md` | v3+ | Med | Re-architected (generic "top lead") |
| Keep/skip feedback loop | `src/research_pipeline/feedback.py` | v2+ | Low | Direct |
| Output-plugin interface | (none — GENERALIZE.md Phase B) | v4 · v4-1, ADR-0012 | Med | New (the seam drafting plugs into) |
| Multi-format draft generation | `scripts/draft.py`, `prompts/formats/` | v4 · v4-2, ADR-0012 | High | **Add-on — divergence #3** |
| Voice lint + tighten | `voice_lint.py`, `voice_tighten.py`, `config/voice*.yaml` | v4 · v4-3 | High | Add-on (encodes one person's taste) |
| Content angles (essays/cityfront) | `personal/<angle>/` | v4 · v4-3 | Med | Add-on |
| Claim lint (factual QA) | `scripts/claim_lint.py` | v4 · v4-3 | Med | Add-on |
| Wrapper convenience scripts | `scripts/wrappers/*.sh` | replaced by cobra subcommands | Low | New |

## G. Dashboard + vault analytics  (Inkwell product, after ingest)

| Feature | research-tool source | Target | Port | Fit |
|---|---|---|---|---|
| Local web server (loopback, stdlib) | `scripts/dashboard.py` | product · ADR-0010 — `net/http` + `embed` | Med | Direct (same pattern) |
| Fixed action allow-list | `dashboard_api.py` `ACTIONS` | product · ADR-0010 | Med | Direct (closed set, no arbitrary exec) |
| Interactive entity graph | `scripts/dashboard/*.js`, `vault_graph.py` | product | High | Direct (reuse the JS, serve from Go) |
| Vault export → JSON (timeline/graph/heatmap) | `vault_export.py`, `export_graph.py` | product / v2 | High | Re-architected (read from SQLite, not vault crawl) |
| Obsidian map build | `scripts/build_obsidian_map.py` | product | Med | Re-architected |
| Recommendations view + one-click run | `dashboard.js` | product | Med | Direct |
| SSE log streaming for actions | `dashboard_api.py` | product | Med | Direct |

## H. Config + infra  (cross-cutting; v1 foundations)

| Feature | research-tool source | Target | Port | Fit |
|---|---|---|---|---|
| Path/config resolution (env > file > default) | `src/research_pipeline/paths.py` | v1 · M1 — ✅ partial (`internal/config`) | Low | Direct (XDG-aware, done) |
| YAML config loading | `paths.py`, `taxonomy.py` | v1 · M1 — `yaml.v3` (not wired yet) | Low | Direct |
| Structured logging | ad-hoc | v1 · M1 — `slog` | Low | New |
| SQLite store + schema | (vault is the store) | v1 · M2 — sqlc, `inkwell migrate` | Med | New (the seam, ADR-0002) |
| Schema/frontmatter migrations | `scripts/migrations/` | v1 · M2 + ongoing | Med | Re-architected (DB migrations) |
| Pipeline orchestration | `scripts/daily.py` | v1 · `Pipeline.Run(ctx)`, ADR-0004 | Med | Re-architected (two pipelines) |
| launchd scheduling | `launchd/*.plist` | v1 · M9, ADR-0004 | Low | Direct (macOS first; cross-platform later) |
| `inkwell status` (feed health) | (none) | v1 · M9 | Low | New |
| Vault maintenance (prune/orphans) | `prune_inbox_archive.py`, `cleanup_orphans.py` | v2+ | Low | Direct |
| Reset + fresh-run test harness | `reset_vault.sh`, `test_fresh_run.sh` | dev tooling | Low | New (Go test fixtures) |

## I. Distribution + product UX  (net-new — Inkwell's whole reason to exist)

| Feature | research-tool gap | Target | Port | Fit |
|---|---|---|---|---|
| Single static binary | venv + pip + editable install | product · ADR-0009 | — | New |
| `goreleaser` + GitHub releases + brew tap | "someday" | product · ADR-0009 | Med | New |
| `inkwell init` guided onboarding | hand-authored YAML | product | Med | New |
| BYO API key flow | assumes Claude CLI | product · ADR-0011 | Low | New |
| Cross-platform scheduling | launchd (macOS only) | post-v1 · revisit ADR-0004 | High | New |
| Versioned releases + changelog + update path | `git pull` | product · ADR-0009 | Low | New |
| Robustness: retries, graceful failures, actionable errors | happy-path scripts | spread across v1+ | Med | Re-architected |

---

## Sprint backlog view (pull from the top)

Ordered so each item is buildable when you reach it. Milestone numbers are ROADMAP v1 milestones; v2/v3/product items follow.

**v1 — ingest (the walking skeleton ships first)**
1. M2 — SQLite store + `inkwell migrate` (the seam) · §H
2. M1 finish — YAML load + cobra + slog + `config show` · §H
3. M3 — feeds CRUD · §A
4. M4 — OPML import · §A
5. M5 — fetcher worker pool (the inflection point) · §A
6. M6 — idempotency + backoff · §A
7. M7 — Obsidian writer (`YYYY/MM/`, atomic) · §A
8. M8 — Story schema + manual attach · §D
9. M9 — launchd + `inkwell status` · §H
10. M10 — 2 weeks real daily use; fix edge cases

**product (interleave once ingest is solid)**
- Single binary + goreleaser + brew (ADR-0009) · §I
- `inkwell init` + BYO-key plumbing (ADR-0011) · §I
- Web dashboard MVP: serve graph + Entry list from SQLite (ADR-0010) · §G

**v2 — embeddings**
- Embed Entries; similarity/dedup · §C
- Full-article enrich (now that short bodies hurt) · §B
- Auto-attach Entries to Stories · §D
- Entity recognition · §B
- Keep/skip feedback · §F

**v3 — synthesis**
- Clustering (agglomerative/HDBSCAN) · §D
- `internal/synth` + Anthropic Go SDK (ADR-0011) · §E
- Storyline weekly synthesis + daily digest · §E
- Recommendations (generic "top lead") · §F
- Full dashboard parity (timeline/heatmap/recommendations) · §G

**v4 — output add-on (optional, last, ADR-0012)**
- v4-1 — `internal/output/` plugin interface (the seam) · §F
- v4-2 — drafting plugin: multi-format publish bundle · §F
- v4-3 — voice lint/tighten, angles, claim lint · §F

---

## Maintenance note

This map reflects research-tool as of the snapshot date in the commit that added
it. When you port a row, mark it (✅/in-progress) and link the PR. When
research-tool changes materially, re-snapshot — it is the moving reference, and a
stale map is worse than none.
