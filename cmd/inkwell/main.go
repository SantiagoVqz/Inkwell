package main

import "fmt"

// version is the build's identifier. The default is "dev"; production builds
// override it at link time via `-ldflags "-X main.version=<git-tag>"` (see
// Makefile). It MUST be a var, not a const — the linker rewrites memory at a
// known address, and constants are inlined at compile time with no address.
var version = "dev"

func main() {
	fmt.Println("inkwell", version)
}
