# Design

## Context

See proposal.md — Why.

`HealthHandler` already follows a probe-function pattern: `dbProbe` and `r2Probe` are `func(ctx context.Context) error` fields set at construction, which makes them swappable in tests without any interface indirection. `bookleaf.Client` is already constructed in `main.go` and passed to several handlers; it is not yet passed to `HealthHandler`.

## Goals / Non-Goals

**Goals:**
- Add a `bookleaf` field to the health response body
- Wire a `bookleafProbe` into `HealthHandler` using the existing probe pattern
- Add a `Ping` method to `bookleaf.Client` that calls `GET <host>/health`

**Non-Goals:**
- Shallow/deep health check variants — not needed; the dependency is one-directional (Booklet depends on Bookleaf, not the reverse), so no loop risk exists
- Exposing Bookleaf's internal response body — Booklet only needs to know whether Bookleaf is reachable, not its component details

## Decisions

### Add `Ping` to `bookleaf.Client` rather than calling the URL directly in `HealthHandler`

`bookleaf.Client` already owns the Bookleaf host and encapsulates all outbound calls to Bookleaf. Putting the health call there keeps all Bookleaf communication in one place and lets tests stub it through the existing probe-func pattern without introducing a new interface.

Alternative considered: pass the host URL directly to `HealthHandler` and have it construct an ad-hoc HTTP request. Rejected — it bypasses the client abstraction and duplicates the host configuration.

### `Ping` calls `GET <host>/health` without the internal secret header

Bookleaf's `/health` is a public, unauthenticated endpoint. Sending the internal secret is unnecessary and would couple health checks to auth configuration.

### Probe timeout is shared (3 seconds) with existing probes

The same `healthProbeTimeout` constant applies to all probes via the single `context.WithTimeout` call in `GetHealth`. No per-probe timeout configuration is needed.

## Risks / Trade-offs

- **Bookleaf unreachable on startup** → health reports degraded immediately. This is correct behavior; no mitigation needed.
- **Flapping Bookleaf** → health oscillates between ok/degraded on each poll. Acceptable for a single-user app with no automated alerting; a circuit breaker would be over-engineering here.
