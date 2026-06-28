# launchd over internal daemon mode for v1

v1 ships `inkwell ingest` as a one-shot CLI invoked by `launchd` on macOS. No `inkwell run --daemon`, no internal scheduler, no `internal/scheduler/` package. A plist template lives in `contrib/launchd/`; `make install-launchd` symlinks it into `~/Library/LaunchAgents/`.

The macOS scheduler already handles sleep/wake survival, missed-run replay, log rotation, and process lifecycle — building an in-process cron loop reinvents work the OS does for free, and would suspend silently when the laptop sleeps. The orchestration code (`internal/ingest/`) exposes a `Pipeline.Run(ctx)` method so a future `--daemon` mode can be added as a thin entry point without refactoring, if we ever move to a non-launchd host (Linux server, Raspberry Pi).

**Considered and rejected**: shipping `--daemon` mode in v1 (operational overhead with no benefit on macOS); committing to launchd-only permanently (locks out portability for negligible code savings).
