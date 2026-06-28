# LLM access: the official Anthropic Go SDK with bring-your-own-key

When the synthesis pipeline lands (v3, per ADR-0002), the LLM is reached through the **official `anthropic-sdk-go`** with a **user-supplied API key** (read from config/env, never bundled) — not by shelling out to the Claude Code CLI the way research-tool does. A public binary cannot assume the Claude CLI is installed, logged in, and on `PATH`. The SDK gives typed requests, streaming, retries, and in-process token accounting.

Default to the latest models with per-stage overrides mirroring research-tool's assignments: Haiku (4.5) for cheap classification, Sonnet (4.6) for clustering judgment, Opus (4.8) for synthesis. Before writing any API code, consult the `claude-api` reference for current model IDs, pricing, and SDK usage — never hardcode model facts from memory. Per ADR-0002, LLM calls stay confined to the synthesis pipeline; ingest never calls the LLM.

This **supersedes ROADMAP milestone v3-2** ("`internal/synth/` interface + `claude` CLI impl (shell out)"). The `internal/synth/` interface still stands; only the first concrete implementation changes from CLI-shell-out to SDK.

**Considered and rejected**: shelling out to the Claude CLI (research-tool's approach — assumes a developer-installed, authenticated dependency; unshippable to end users; brittle to CLI version/auth changes); bundling a vendor API key (insecure and unaffordable for an OSS binary — abuse and cost land on the author); an LLM-provider abstraction layer up front (premature — research-tool's GENERALIZE.md explicitly deferred multi-provider; add an interface only when a second provider is real).
