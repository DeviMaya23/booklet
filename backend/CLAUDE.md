# Backend Constraints

Critical constraints that contradict common defaults. Full conventions in `/CONVENTIONS.md`.

## Go Version

- The `go` directive in `go.mod` is pinned deliberately — never bump it yourself.
- After any dependency change, run `git diff go.mod`. If `go` or `toolchain` changed, revert it.
- Look for an older compatible version of the dep instead. If none exists, stop and tell the user — do not bump the Go version.

## Database Migrations

- Never use GORM `AutoMigrate`. All schema changes go through a numbered
  `golang-migrate` SQL pair in `backend/migration/` — see `ARCHITECTURE.md`
  for numbering and the domain-struct sync rule.
- Migrations apply automatically on merge to `main` — there is no manual
  gate before they hit production. Write and test the `down.sql` as
  carefully as the `up.sql`.
