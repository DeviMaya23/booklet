# Health Check

## Purpose

Defines the health endpoint that probes system dependencies and reports their status. The endpoint is unauthenticated and always returns HTTP 200; component health is communicated in the response body.

---

## Requirements

### Requirement: Health endpoint
The system SHALL expose `GET /health` as an unauthenticated endpoint that probes the database and R2 storage. The response SHALL always return HTTP 200 regardless of probe outcomes; component health is communicated in the response body.

Response body fields:
- `status` — `"ok"` if all probes pass; `"degraded"` if any probe fails
- `db` — `"ok"` if the database probe passes; the error message string if it fails
- `r2` — `"ok"` if the R2 probe passes; the error message string if it fails

Each probe runs with a 3-second timeout.

#### Scenario: All components healthy
- **WHEN** `GET /health` is called and both the database and R2 are reachable
- **THEN** the system returns 200 with `{ "status": "ok", "db": "ok", "r2": "ok" }`

#### Scenario: Database probe fails
- **WHEN** `GET /health` is called and the database is unreachable
- **THEN** the system returns 200 with `"status": "degraded"` and `"db"` containing the error message; `"r2"` reflects its actual probe result

#### Scenario: R2 probe fails
- **WHEN** `GET /health` is called and R2 is unreachable
- **THEN** the system returns 200 with `"status": "degraded"` and `"r2"` containing the error message; `"db"` reflects its actual probe result
