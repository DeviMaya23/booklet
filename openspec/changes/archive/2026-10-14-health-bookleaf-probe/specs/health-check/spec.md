# Spec Delta

## MODIFIED Requirements

### Requirement: Health endpoint
The system SHALL expose `GET /health` as an unauthenticated endpoint that probes the database, R2 storage, and Bookleaf. The response SHALL always return HTTP 200 regardless of probe outcomes; component health is communicated in the response body.

Response body fields:
- `status` — `"ok"` if all probes pass; `"degraded"` if any probe fails
- `db` — `"ok"` if the database probe passes; the error message string if it fails
- `r2` — `"ok"` if the R2 probe passes; the error message string if it fails
- `bookleaf` — `"ok"` if the Bookleaf probe passes; the error message string if it fails

Each probe runs with a 3-second timeout.

#### Scenario: All components healthy
- **WHEN** `GET /health` is called and the database, R2, and Bookleaf are all reachable
- **THEN** the system returns 200 with `{ "status": "ok", "db": "ok", "r2": "ok", "bookleaf": "ok" }`

#### Scenario: Database probe fails
- **WHEN** `GET /health` is called and the database is unreachable
- **THEN** the system returns 200 with `"status": "degraded"` and `"db"` containing the error message; `"r2"` and `"bookleaf"` reflect their actual probe results

#### Scenario: R2 probe fails
- **WHEN** `GET /health` is called and R2 is unreachable
- **THEN** the system returns 200 with `"status": "degraded"` and `"r2"` containing the error message; `"db"` and `"bookleaf"` reflect their actual probe results

#### Scenario: Bookleaf probe fails
- **WHEN** `GET /health` is called and Bookleaf is unreachable
- **THEN** the system returns 200 with `"status": "degraded"` and `"bookleaf"` containing the error message; `"db"` and `"r2"` reflect their actual probe results
