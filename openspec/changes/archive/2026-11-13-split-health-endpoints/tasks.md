# Tasks

## 1. Deploy prerequisite (secret)

- [x] 1.1 Create the Secret Manager secret `BOOKLET_HEALTH_CHECK_SECRET` with a freshly generated random value, grant the Cloud Run service account access, and verify with `gcloud secrets describe BOOKLET_HEALTH_CHECK_SECRET` and a successful accessor binding check. This must exist before the change merges.

## 2. Config

- [x] 2.1 Add `HealthCheckSecret string` to `Config` in `backend/internal/platform/config/config.go`, loaded with `requireEnv("HEALTH_CHECK_SECRET")`, following how `INTERNAL_API_SECRET` is loaded. Verify `go build ./...` passes.
- [x] 2.2 In `config_test.go`, set `HEALTH_CHECK_SECRET` in `setRequiredEnvVars`, assert `cfg.HealthCheckSecret` is loaded in the all-vars-set test, and add a missing-var case asserting `Load` returns an error naming `HEALTH_CHECK_SECRET`. Verify `go test ./internal/platform/config/...` passes.
- [x] 2.3 Add `HEALTH_CHECK_SECRET` with a placeholder value to `backend/.env.example` near the other shared secrets. Verify the file lists it.

## 3. Constant-time internal secret check

- [x] 3.1 In `backend/internal/handler/middleware/internal_auth.go`, replace the `!=` header comparison with `crypto/subtle.ConstantTimeCompare`, keeping the header name, signature and 401 response unchanged. Verify `go build ./...` passes.
- [x] 3.2 Add `backend/internal/handler/middleware/internal_auth_test.go` (table-driven, `middleware_test` package like `auth_test.go`) asserting status code and that the next handler runs only when it should: missing header returns 401 and `next` is not called; wrong value returns 401 and `next` is not called; correct value calls `next` and returns its response. Also mount two routes with two middleware instances using different secrets and assert each secret is rejected on the other route (covers the `INTERNAL_API_SECRET` vs `HEALTH_CHECK_SECRET` scenarios). Verify `go test ./internal/handler/middleware/...` passes.

## 4. Health handler

- [x] 4.1 In `backend/internal/handler/health.go`, make `GetHealth` the shallow handler returning 200 `{"status":"ok"}` without touching any probe, and move the probe logic to a new `GetDeepHealth`. Verify `go build ./...` passes.
- [x] 4.2 In `GetDeepHealth`, run the db, R2 and Bookleaf probes in parallel goroutines joined with `sync.WaitGroup`, each with its own timeout context derived from the request context. Move the timeout from the `healthProbeTimeout` constant to a handler field defaulting to 3s (set in `NewHealthHandler`). Return 503 when the db or R2 probe fails and 200 otherwise; a Bookleaf failure sets only the `bookleaf` field. Verify with the tests in 4.3.
- [x] 4.3 Rewrite `health_handler_test.go` as table-driven handler tests asserting both status code and the decoded body for every scenario in the spec: shallow returns 200 `{"status":"ok"}` and invokes no probe (probes set to fail the test if called); deep all healthy returns 200 with all fields `"ok"`; db fails returns 503 with `status` `"degraded"` and `db` equal to the exact error text; r2 fails returns 503 likewise; only Bookleaf fails returns 200 with `status` `"ok"` and `bookleaf` equal to the exact error text; Bookleaf plus db fail returns 503; a probe that blocks until its context is done (with the timeout field shortened) reports `context deadline exceeded` in its own field while the other probes report their real results. Verify `go test ./internal/handler/...` passes.

## 5. Routing and deploy wiring

- [x] 5.1 In `backend/cmd/server/main.go`, register `e.GET("/health", healthHandler.GetHealth)` and `e.GET("/health/deep", healthHandler.GetDeepHealth, authmiddleware.NewInternalAuthMiddleware(cfg.HealthCheckSecret))` as per-route middleware, without adding a new group. Verify `go build ./...` passes and, with the server running locally, `curl` returns 200 for `/health`, 401 for `/health/deep` without the header, and the probe body with it.
- [x] 5.2 In `.github/workflows/build-push.yml`, add `HEALTH_CHECK_SECRET=BOOKLET_HEALTH_CHECK_SECRET:latest` to the `secrets:` block of the Cloud Run deploy step. Verify the YAML parses (for example `python3 -c "import yaml,sys; yaml.safe_load(open('.github/workflows/build-push.yml'))"`) and the line is present.

## 6. API collection

- [x] 6.1 Update `collection/health.bru` so it describes the shallow liveness check (URL `{{baseUrl}}/health`, no auth). Verify by running it in Bruno against a local server and getting 200 `{"status":"ok"}`.
- [x] 6.2 Add `collection/health-deep.bru` for `GET {{baseUrl}}/health/deep` with an `X-Booklet-Internal-Secret: {{healthCheckSecret}}` header, declaring `healthCheckSecret` in `vars:pre-request` with a placeholder matching `.env.example` (same pattern as `collection/internal/delete_user.bru`). Verify by running it in Bruno against a local server with `HEALTH_CHECK_SECRET` set and getting the probe body, and 401 with the header removed.

## 7. Final checks

- [x] 7.1 Run `golangci-lint run` from `backend/` and fix every issue it reports. Verify it exits clean.
- [x] 7.2 Run `go test ./...` from `backend/` and verify the whole suite passes, including the deletion-route and config tests.
- [x] 7.3 Run `git diff go.mod` and verify the `go` and `toolchain` lines are unchanged (no dependency was added; `sync` and `crypto/subtle` are stdlib).
