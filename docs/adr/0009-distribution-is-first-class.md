# Distribution is first-class: a single static binary via goreleaser

Inkwell ships as a **single statically-linked Go binary**, released via `goreleaser` to tagged GitHub releases with a Homebrew tap, plus an `inkwell init` onboarding flow (guided first-run config + key setup). This collapses the two largest shippability gaps inherited from research-tool — multi-step install (`venv` + `pip install -e` + copy plist) and OS-specific scheduling — into one download. The binary is the unit of distribution; config and vault stay external (XDG paths, see `internal/config`).

To keep cross-compilation to darwin/linux/windows × amd64/arm64 trivial, **avoid CGO**: use a pure-Go SQLite driver (`modernc.org/sqlite`) rather than `mattn/go-sqlite3`. The whole reason for choosing Go (ADR-0008) is frictionless single-binary distribution; a CGO dependency would forfeit it.

This **amends the ROADMAP**, which previously deferred `goreleaser` to "when actually distributed." Distribution is designed in from the start because — the research-tool lesson — a tool that only runs on its author's laptop never crosses the shippable line, and packaging constraints (no CGO, external config, init flow) shape architecture early.

**Considered and rejected**: deferring packaging until "later" (packaging shapes architecture; retrofitting it is the expensive path research-tool is stuck in); CGO SQLite (`mattn/go-sqlite3`) for marginal performance (breaks the easy cross-compilation that is Go's entire advantage here); per-OS native installers (`.pkg`/`.msi`/`.deb`) before a binary release exists (defeats the single-binary advantage; revisit only if a release ever needs OS-level integration).
