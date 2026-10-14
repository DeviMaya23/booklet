# Tasks

## 1. Add `Ping` to `bookleaf.Client`

- [x] 1.1 Add `Ping(ctx context.Context) error` to `bookleaf.Client` in `backend/internal/bookleaf/client.go` — issues `GET <host>/health` (no auth header), returns nil on 2xx, error otherwise; verify it compiles
- [x] 1.2 Add unit test for `Ping` in `backend/internal/bookleaf/client_test.go` covering: successful ping (2xx → nil), non-2xx response (→ error), network error (→ error)

## 2. Extend `HealthHandler`

- [x] 2.1 Add `bookleafProbe func(ctx context.Context) error` field to `HealthHandler` and `Bookleaf string` field to `healthResponse` in `backend/internal/handler/health.go`; set the probe in `NewHealthHandler` using a `bookleafClient` parameter; verify it compiles
- [x] 2.2 Run the probe in `GetHealth` alongside the existing probes, populate `res.Bookleaf`, and set `res.Status = "degraded"` on failure; verify existing tests still pass (`go test ./internal/handler/...`)
- [x] 2.3 Add unit test scenarios in `backend/internal/handler/health_handler_test.go`: bookleaf probe fails → status degraded + bookleaf field contains error message; all probes pass → bookleaf field is "ok"

## 3. Wire in `main.go`

- [x] 3.1 Update `httphandler.NewHealthHandler` call in `backend/cmd/server/main.go` to pass `bookleafClient` (already constructed above); verify `go build ./cmd/server/...` succeeds

## 4. Lint

- [x] 4.1 Run `golangci-lint run ./...` from `backend/` and fix any reported issues
