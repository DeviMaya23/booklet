# Design

## Context

See proposal.md for motivation. Current state:

- `HealthHandler.GetHealth` (`backend/internal/handler/health.go`) runs three probes sequentially on one shared 3s context and always returns 200. The probes are injectable function fields (`dbProbe`, `r2Probe`, `bookleafProbe`), which the existing handler test relies on.
- `main.go` registers `e.GET("/health", ...)` directly on the root echo instance, outside the `protected` and `internal` groups.
- `NewInternalAuthMiddleware(secret)` already exists and guards `DELETE /internal/users/:id` with `X-Booklet-Internal-Secret`. It uses a plain `!=` comparison and has no direct test. The secret comes from `INTERNAL_API_SECRET` via `requireEnv`.
- Env vars reach Cloud Run through the `secrets:` and `env_vars:` blocks in `.github/workflows/build-push.yml`.
- No code in the repo (backend, frontend, workflows, Dockerfile) calls `/health`. The only consumer will be the planned monitor app.

## Goals / Non-Goals

**Goals:**
- A cheap, public liveness endpoint that cannot be used to amplify load on the database, R2 or Bookleaf.
- A readiness endpoint, reachable only with a secret, whose HTTP status reflects whether db and R2 are usable.
- No new middleware, abstraction or dependency.

**Non-Goals:**
- Rate limiting or caching of the deep endpoint (the secret gate makes this unnecessary for a single monitor).
- Adding a Docker `HEALTHCHECK` or a Cloud Run probe.
- Any change to Bookleaf or to how Booklet calls Bookleaf's `/health`.
- Changing the secret-check middleware's header name, signature or status codes.

## Decisions

### Separate paths, not a query parameter
`/health` is shallow and `/health/deep` is deep. The access rule (public vs secret) follows the route, so the secret middleware attaches to one route and the cheap path cannot be turned expensive by a stray parameter. Alternative: `?deep=true` on one route. Rejected because it needs per-request branching on auth and makes accidental deep probing possible.

### `/health` becomes the shallow one
Nothing calls `/health` today, so repurposing the unqualified name costs nothing and means any future platform probe pointed at it gets the cheap behavior by default. Alternative: keep `/health` deep and add `/health/live`. Rejected because the unqualified name would remain the expensive one.

### Reuse `NewInternalAuthMiddleware` with a separate secret
The new route calls `NewInternalAuthMiddleware(cfg.HealthCheckSecret)`, attached per route (`e.GET("/health/deep", h.GetDeepHealth, mw)`) rather than via a third group. This adds no new middleware or pattern. The secret is a new `HEALTH_CHECK_SECRET`, not `INTERNAL_API_SECRET`, because the latter authorizes user deletion and the monitor is a separate app that will hold the value. The header name stays `X-Booklet-Internal-Secret`; a dedicated header would need a parameterized middleware for no real gain. Because each instance closes over its own secret, a secret valid for one route is rejected by the other.

### Constant-time comparison in the shared middleware
The comparison becomes `subtle.ConstantTimeCompare([]byte(header), []byte(secret)) == 1`. Accept/reject outcomes are identical to `!=` for every input, so Bookleaf's call to `DELETE /internal/users/:id` is unaffected and Bookleaf needs no change. Because `requireEnv` rejects empty secrets, the empty-equals-empty case cannot occur. Length is leaked on mismatch, which is the accepted trade-off for this API. This touches a file shared with the user-deletion route, so a new unit test (`internal_auth_test.go`) covers missing, wrong and correct headers directly. The `account-deletion` spec already has scenarios for a valid and for a missing or invalid secret, so it needs no delta.

### 503 on db or R2 failure; Bookleaf excluded
The deep response returns 503 when the db or R2 probe fails and 200 otherwise, matching the readiness-probe convention (2xx pass, anything else fail). Bookleaf is optional (moodboards only), so its failure is visible in the body but does not change `status` or the code. A monitor that only checks the status code will not alert on a Bookleaf outage; alerting on it requires reading the `bookleaf` field. The body shape is unchanged from today (`status`, `db`, `r2`, `bookleaf`) and error strings stay in it, since only a secret holder can see it.

### Parallel probes with per-probe timeouts
Each probe gets its own `context.WithTimeout(ctx, 3s)` derived from the request context and runs in its own goroutine, joined with `sync.WaitGroup`. Each goroutine writes only its own result variable, so no mutex is needed, and no new dependency (`errgroup`) is required. This removes the current flaw where a slow database leaves later probes with an already-expired context. Worst-case latency drops from the sum of the timeouts to one timeout. Because the goroutines run outside Echo's per-request recovery, each one has a `defer recover()` that records a panic as that probe's error (`probe panicked: <value>`), so a panicking probe fails its own field instead of crashing the process. The injectable probe fields are kept so handler tests stay table-driven without real dependencies.

### Handler shape
`GetHealth` becomes the shallow handler (static `{"status":"ok"}`); `GetDeepHealth` holds the probe logic and keeps `healthResponse`. `NewHealthHandler`'s signature is unchanged. The per-probe timeout moves from a package constant to a handler field defaulting to 3s, so the slow-probe test can shorten it instead of waiting three real seconds.

### Required config, deploy dependency
`HEALTH_CHECK_SECRET` is loaded with `requireEnv` into a new `Config.HealthCheckSecret`, following `INTERNAL_API_SECRET`. The server therefore refuses to start without it, which is fail-closed but means the value must exist in production before the deploy. `build-push.yml` gets a `HEALTH_CHECK_SECRET=<secret-name>:latest` line in its `secrets:` block, pointing at a new Secret Manager secret.

## Risks / Trade-offs

- **Deploy ordering** → The Cloud Run deploy fails or the new revision crashes at startup if the Secret Manager secret or the workflow mapping is missing. Mitigation: create the secret and grant the service account access before merging; the task list puts this first.
- **`/health` contract change** → Any external caller expecting the old body would silently get `{"status":"ok"}`. Mitigation: nothing in the repo calls it; confirm no external uptime tool is pointed at it before merge.
- **Bookleaf outage invisible to status-code monitors** → Intended by decision; the monitor must read the `bookleaf` field to alert on it.
- **Shared middleware edit** → The constant-time change sits on the user-deletion route. Mitigation: behavior is equivalent for all inputs, and the new unit test plus the existing `account-deletion` scenarios cover it.
- **Same header name for two secrets** → A misconfigured monitor sending the wrong secret gets a plain 401 with no hint which secret was expected. Acceptable for a single-user project.

## Migration Plan

1. Create the Secret Manager secret and grant the Cloud Run service account access.
2. Merge; the workflow deploys with `HEALTH_CHECK_SECRET` mapped.
3. Point the monitor at `GET /health` (liveness) and `GET /health/deep` with the header (readiness).

Rollback: redeploy the previous revision. It does not read `HEALTH_CHECK_SECRET`, so the extra secret mapping is harmless.
