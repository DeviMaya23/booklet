# Proposal

## Why

Booklet depends on Bookleaf for folder image data, but the health endpoint has no visibility into whether Bookleaf is reachable. A degraded Bookleaf silently breaks character folder features without surfacing in monitoring.

## What Changes

- `GET /health` adds a third probe: an HTTP call to Bookleaf's `GET /health` endpoint
- The response body gains a `bookleaf` field reporting the probe result
- `status` is set to `"degraded"` if the Bookleaf probe fails, alongside the existing db and R2 probes
- A new `Ping` method is added to `bookleaf.Client` to perform the outbound health check
- The `HealthHandler` is extended with a `bookleafProbe` func following the existing probe pattern

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `health-check`: adds the Bookleaf probe to the health endpoint and a `bookleaf` field to the response body

## Impact

- `backend/internal/bookleaf/client.go` — new `Ping` method
- `backend/internal/handler/health.go` — new `bookleafProbe` field, extended `healthResponse`
- `backend/internal/handler/health_handler_test.go` — new test scenarios for bookleaf probe pass/fail
- `backend/internal/bookleaf/client_test.go` — new test for `Ping`
- `backend/cmd/server/main.go` — passes `bookleaf.Client` to `NewHealthHandler`
