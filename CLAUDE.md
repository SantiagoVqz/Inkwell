# Inkwell

A source-agnostic research engine: it ingests sources on a schedule, organizes
them into a knowledge vault, and (later) clusters + synthesizes them with an LLM.
Written in Go, shipped as a **single static binary**.

> **North star:** turn the working-but-unshippable Python `research-tool` into a
> tool a stranger can download and run. Inkwell is the product; `research-tool`
> is its reference spec and the brain we port from, one vertical slice at a time.

This file is the orientation layer. The durable detail lives in:
- [`CONTEXT.md`](CONTEXT.md) — the canonical glossary (ubiquitous language). **Read it before touching the domain.**
- [`ROADMAP.md`](ROADMAP.md) — milestones, acceptance criteria, what each version does *not* ship.
- [`docs/adr/`](docs/adr/) — Architecture Decision Records: the *why* behind each choice, with rejected alternatives.
- [`docs/research-tool-feature-map.md`](docs/research-tool-feature-map.md) — the migration backlog: every research-tool feature mapped to an Inkwell version/milestone, in sprint-pull order. **Start here when picking what to build next.**

---

## The two repos, and why both exist

| | `research-tool` (Python) | `Inkwell` (Go) — this repo |
|---|---|---|
| State | ~9k LOC, mature, **works today** | ~54 LOC, milestones 0–1, **barely started** |
| Role | The brain + reference spec. Daily driver. | The product we're building toward. |
| Strength | LLM stages, voice logic, the dashboard UI | single-binary distribution, speed, shippability |
| Weakness | can't ship (venv/pip/launchd/hand-edited YAML; assumes Claude CLI installed) | doesn't do the smart stuff yet |

Inkwell catches up **slice by slice** (ingest → enrich → classify → cluster →
synthesize → dashboard). Until it does, `research-tool` keeps running in parallel,
so no capability is ever lost (this is the premise of [ADR-0001](docs/adr/0001-v1-scope-ingest-only.md)).

`research-tool` lives at `~/Documents/Personal/research-tool`. When porting a
stage, read its Python as the spec — the prompts are text files that carry over
verbatim; what you re-derive is the glue, not the intelligence.

---

## The product decisions (settled, ADR-ratified)

These came out of the strategy discussion that motivated this file and are now
recorded as ADRs — the load-bearing record, not folklore.

1. **Language: Go.** The dominant constraint is *distribution* — a single static
   binary that cross-compiles to every OS collapses the two biggest shippability
   gaps (install + scheduling) for free. Go also fits the workload: concurrent
   fetch, stdlib HTTP for the dashboard, cobra for the CLI.
   ([ADR-0008](docs/adr/0008-inkwell-is-the-shippable-product.md))
2. **Inkwell is the product; `research-tool` is the reference spec**, run in
   parallel during migration so no capability is lost.
   ([ADR-0008](docs/adr/0008-inkwell-is-the-shippable-product.md))
3. **Distribution is first-class:** single static binary via `goreleaser` +
   releases + brew tap + `inkwell init`. Part of "done," not "someday."
   ([ADR-0009](docs/adr/0009-distribution-is-first-class.md))
4. **Form factor: a CLI that embeds a local web dashboard** on `127.0.0.1`
   (the Ollama / Syncthing / Jupyter pattern) — re-implementing research-tool's
   graph + recommendations + action allow-list. Pulled forward from "late v3+."
   ([ADR-0010](docs/adr/0010-cli-with-embedded-web-dashboard.md))
5. **LLM access: the official Anthropic Go SDK (`anthropic-sdk-go`),
   bring-your-own-key** — not a shelled-out Claude CLI. Latest models
   (Opus 4.8 / Sonnet 4.6 / Haiku 4.5); read the `claude-api` skill before
   writing API code.
   ([ADR-0011](docs/adr/0011-anthropic-go-sdk-byo-key.md))

**Still open (revisit, no ADR yet):** cross-platform scheduling. v1 stays on
`launchd` ([ADR-0004](docs/adr/0004-launchd-over-internal-daemon.md)) — correct
for a macOS-only start — but public reach means systemd/cron/Task Scheduler or
an internal scheduler later. ADR-0004 already keeps `Pipeline.Run(ctx)`
host-agnostic; write the ADR when we actually leave macOS.

What has **not** changed: the ingest-first scope ([ADR-0001](docs/adr/0001-v1-scope-ingest-only.md)),
the two-pipeline + SQLite-seam architecture ([ADR-0002](docs/adr/0002-two-pipeline-architecture.md)),
the Entry/Story/Storyline ontology ([ADR-0003](docs/adr/0003-entry-story-storyline-ontology.md)),
and the vault/idempotency rules (ADR-0005, -0006, -0007). Those decisions stand.

---

## How we work (the professional loop)

This project is also a deliberate exercise in attacking a build the way a senior
engineer does. Inkwell's own structure is the template — keep it that way:

1. **Spec before code.** Know the one job and what's out of scope before typing.
2. **Ubiquitous language.** The domain nouns are fixed in `CONTEXT.md`. Use them
   exactly; never introduce a synonym (no "article" for Entry, no "topic" for Story).
3. **ADR every architectural decision** — with the alternatives you rejected and
   why. If a choice would surprise a future reader, it needs an ADR. CLAUDE.md
   records direction; ADRs record commitments.
4. **Walking skeleton first.** The thinnest end-to-end slice (one feed → one note
   on disk → clean exit) is wired and green before any layer is fleshed out.
5. **Vertical slices with written acceptance criteria** (ROADMAP already does
   this). Each milestone ships something demonstrable.
6. **Tests + CI from slice one.** Trunk-based, small PRs; every push runs
   `go test ./...`, `go vet ./...`, and a build.

Harness skills that map to this loop: `write-a-prd`, `grill-me-with-docs`
(sharpen glossary + ADRs), `prd-to-issues` (tracer-bullet slices), `tdd`,
`find-critical-gaps`, `commit-push-pr`.

### Pairing style (how to write code with the owner)

The owner is a senior TS/Node/Python engineer, new to Go. When implementing:

- **Write the working code first, then explain it** — don't pre-narrate code the
  owner would just copy-paste. Build it, run the tests, then walk through it.
- **Keep the walkthrough at "enough to own it" altitude:** what each piece does,
  how the moving parts connect, and where they'd make a small change. Not an
  expert Go tutorial.
- **Explain a new Go or SQL concept once, briefly, the first time it appears**
  (e.g. `//go:embed`, `sql.NullString`, expression indexes) — a sentence or two,
  not a deep dive. Skip language-agnostic basics.
- Verify before explaining: `go test ./...` green and `go vet` clean is the
  precondition for the walkthrough, so the explanation describes code that works.

---

## Current state

- **Milestone 0 (bootstrap):** done. `inkwell version` prints; Makefile builds
  with `-ldflags` version injection.
- **Milestone 1 (config + CLI scaffold):** partial. `internal/config` has typed
  `Config` + `Default()` + XDG-aware `DefaultPath()`. **Not yet:** YAML loading,
  cobra, slog wiring, a `config show` command.
- **Next: milestone 2 — the store layer** (sqlc + SQLite, in-memory tests,
  `inkwell migrate` creates the schema). This is the first real Go you write.

`go.mod` currently has **zero dependencies** (Go 1.26). cobra, sqlc, yaml.v3,
and the SQLite driver arrive as their milestones land — add them with intent,
not all at once.

---

## Domain vocabulary (canonical: `CONTEXT.md`)

Quick reference — `CONTEXT.md` is the source of truth, this is just a map:

- **Feed** → N **Entry** (one item, immutable, one note per Entry)
- **Entry** 0..1 → 1 **Story** via an **Attachment** (records *how*: frontmatter / cli / embedding / llm)
- **Story** — persistent, human-named narrative thread (open/closed). **≠ Cluster.**
- **Cluster** — ephemeral, algorithm-named grouping during a v3 synthesis run.
- **Storyline** — a v3 weekly synthesis note referencing many Stories + Entries.
- **Ingest pipeline** — fetch→normalize→write→(embed)→store. Frequent. **No LLM, ever.**
- **Synthesis pipeline** (v3) — cluster→summarize→write Storyline. Infrequent. **LLM only here.**

---

## Build & test

```bash
make build      # → bin/inkwell (static, version-stamped)
make run ARGS="version"
make test       # go test ./...
make vet        # go vet ./...
make fmt        # go fmt ./...
make tidy       # go mod tidy
make help       # list targets
```

Acceptance bar for any change: `go test ./...` passes, `go vet ./...` is clean,
the binary builds. Match the surrounding code's idiom and comment density —
the existing files comment the *why* (see `cmd/inkwell/main.go` on why `version`
is a `var`, not a `const`); keep that standard.
