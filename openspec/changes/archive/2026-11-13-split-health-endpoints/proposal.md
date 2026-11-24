# Proposal

## Why

`GET /health` is currently a public, unauthenticated deep check: every request runs a database query, an R2 call and an outbound call to Bookleaf. Anyone can use it to drive load on all three, and it returns raw error strings. It also always returns HTTP 200, so an HTTP-status monitor never sees a failure. A separate monitor app is planned to watch production liveness and readiness, and it needs a cheap public liveness check plus a protected readiness check that signals failure through the status code.

## What Changes

- `GET /health` becomes a shallow liveness check: no auth, no dependency calls, returns `200 {"status":"ok"}`.
- New `GET /health/deep` runs the dependency probes (db, R2, Bookleaf). It is protected by a shared-secret header and returns `503` when the database or R2 probe fails.
- **BREAKING**: `GET /health` no longer probes dependencies and no longer returns the `db`, `r2` and `bookleaf` fields. Nothing in the repo calls it (verified), so there is no known consumer to migrate.
- **BREAKING**: the "always HTTP 200" rule is removed for the deep check; failures of db or R2 return 503.
- Bookleaf is reported in the deep response body but does not affect the overall `status` or the HTTP status code (it is an optional dependency, used only for moodboards).
- Deep probes run in parallel, each with its own timeout, instead of sequentially on one shared context.
- `/health/deep` reuses the existing `NewInternalAuthMiddleware` (header `X-Booklet-Internal-Secret`), instantiated with a new, separate `HEALTH_CHECK_SECRET`. `INTERNAL_API_SECRET` is not reused because it authorizes user deletion.
- New env var `HEALTH_CHECK_SECRET`, loaded with `requireEnv`, so the server refuses to start without it. The deploy workflow maps it from a new Secret Manager secret.
- The internal-secret middleware compares secrets in constant time (`crypto/subtle`). Accept/reject behavior is unchanged for callers, including Bookleaf's account-deletion call. The middleware currently has no direct test, so one is added.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `health-check`: the single "always 200, public, deep" requirement is replaced by a shallow liveness requirement and a secret-protected deep readiness requirement with 503 semantics and Bookleaf excluded from overall status.

## Impact

- **Backend handler**: `backend/internal/handler/health.go` (shallow method, reworked deep method, parallel probes) and `health_handler_test.go`.
- **Routing**: `backend/cmd/server/main.go` registers `/health` and `/health/deep` (the latter behind the middleware).
- **Middleware**: `backend/internal/handler/middleware/internal_auth.go` (constant-time compare) plus a new `internal_auth_test.go`. This middleware also guards `DELETE /internal/users/:id`, so that route is covered by the same test and its caller (Bookleaf) is unaffected on the wire.
- **Config**: `backend/internal/platform/config/config.go` and `config_test.go` (new `HealthCheckSecret`), `backend/.env.example`.
- **Deploy**: `.github/workflows/build-push.yml` maps `HEALTH_CHECK_SECRET` from a new Secret Manager secret, which must exist before the deploy runs.
- **API collection**: `collection/health.bru` updated, new `collection/health-deep.bru`.
- **Frontend and extension**: none. No call sites of `/health` exist in the frontend or elsewhere in the repo.
- **Bookleaf**: no change. Booklet still calls Bookleaf's own `/health` from the deep check, and the secret-comparison change does not alter the wire contract.
