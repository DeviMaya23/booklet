# Backend Constraints

Critical constraints that contradict common defaults. Full conventions in `/CONVENTIONS.md`.

## Go Version

- The `go` directive in `go.mod` is pinned deliberately — never bump it yourself.
- After any dependency change, run `git diff go.mod`. If `go` or `toolchain` changed, revert it.
- Look for an older compatible version of the dep instead. If none exists, stop and tell the user — do not bump the Go version.
