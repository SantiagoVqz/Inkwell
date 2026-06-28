# Form factor: a CLI that embeds a local web dashboard, not a desktop GUI

Inkwell is a CLI that, on demand (`inkwell dashboard` / `inkwell serve`), starts a local HTTP server on `127.0.0.1` serving the read-and-act UI — the Ollama / Syncthing / Jupyter pattern. Static assets are compiled into the binary via Go's `embed`, so there is no separate frontend deploy and no Electron/Tauri runtime to bundle. This re-implements research-tool's dashboard (`scripts/dashboard.py` + `src/research_pipeline/dashboard_api.py`): the interactive entity graph, the recommendations view, and one-click pipeline actions.

Actions map to a **closed allow-list** mirroring research-tool's `ACTIONS` (`ingest`, `enrich`, `daily`, `export_graph`, `digest`, `insights`, …) — the server resolves an action name to a fixed argv/handler and **never executes arbitrary user-supplied commands**. Bound to loopback only; single-user, no auth, no multi-tenant.

This **amends the ROADMAP**, which listed the "HTTP read API" and TUI as "late v3+". The web dashboard is pulled forward to a defined product milestone **after the ingest walking skeleton is solid** (it needs Entries in the store to show), because it is the product's primary surface, not an afterthought. A `bubbletea` TUI remains a possible *secondary* surface for terminal-only ops, but the graph visualization wants a browser, so the web UI leads.

**Considered and rejected**: desktop GUI via Electron/Tauri (heavy runtime, bundling/signing pain, wrong shape for a batch/scheduled pipeline); a TUI as the *primary* surface (can't render the entity graph well; fine as a later secondary ops view); a hosted/multi-tenant web app (out of scope — single-user local tool, the same boundary research-tool's GENERALIZE.md drew).
