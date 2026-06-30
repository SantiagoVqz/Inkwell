package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// version is the build's identifier. The default is "dev"; production builds
// override it at link time via `-ldflags "-X main.version=<git-tag>"` (see
// Makefile). It MUST be a var, not a const — the linker rewrites memory at a
// known address, and constants are inlined at compile time with no address.
var version = "dev"

func main() {
	// Cancel the root context on Ctrl-C / SIGTERM so long-running commands
	// (later: ingest) can stop cleanly. signal.NotifyContext is the modern
	// idiom — the returned ctx is Done() the moment a signal arrives.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Execute prints the error (and, for usage errors, the usage) itself; we
	// only need to translate a non-nil error into a non-zero exit code.
	if err := newRootCmd().ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
