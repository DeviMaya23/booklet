# Health Check

## Purpose

Defines the health endpoint that probes system dependencies and reports their status. The endpoint is unauthenticated and always returns HTTP 200; component health is communicated in the response body.

---

## Requirements

### Requirement: Liveness health endpoint
The system SHALL expose `GET /health` as an unauthenticated shallow liveness endpoint. It SHALL NOT probe the database, R2 storage or Bookleaf, and SHALL return HTTP 200 with the JSON body `{ "status": "ok" }` whenever the process is serving requests.

#### Scenario: Process is up
- **WHEN** `GET /health` is called without any headers
- **THEN** the system returns 200 with `{ "status": "ok" }`

#### Scenario: Dependencies are down
- **WHEN** `GET /health` is called while the database, R2 and Bookleaf are all unreachable
- **THEN** the system still returns 200 with `{ "status": "ok" }` and makes no call to any of them

### Requirement: Deep health endpoint probes
The system SHALL expose `GET /health/deep`, which probes the database, R2 storage and Bookleaf. Each probe SHALL run independently with its own 3-second timeout, so a slow or failed probe does not shorten or fail the others.

#### Scenario: All components healthy
- **WHEN** `GET /health/deep` is called with a valid secret and the database, R2 and Bookleaf are all reachable
- **THEN** the system returns 200 with `{ "status": "ok", "db": "ok", "r2": "ok", "bookleaf": "ok" }`

#### Scenario: One probe is slow
- **WHEN** `GET /health/deep` is called with a valid secret and one probe runs past its 3-second timeout while the others succeed
- **THEN** that probe's field contains its timeout error and the other probes report their actual results rather than a timeout error

### Requirement: Deep health response body
The deep endpoint SHALL return a JSON body with the fields `status`, `db`, `r2` and `bookleaf`. Each of `db`, `r2` and `bookleaf` SHALL be `"ok"` if its probe passes, or the probe's error message string if it fails. `status` SHALL be `"ok"` if the database and R2 probes pass and `"degraded"` if either fails.

#### Scenario: Probe error is surfaced
- **WHEN** `GET /health/deep` is called with a valid secret and the R2 probe fails with an error message
- **THEN** the body's `r2` field contains that exact error message and the `db` and `bookleaf` fields each contain `"ok"` or their own error message

### Requirement: Deep health status code
The deep endpoint SHALL return HTTP 503 when the database or R2 probe fails, otherwise HTTP 200. Bookleaf SHALL be reported in the body but SHALL NOT affect `status` or the HTTP status code.

#### Scenario: Database probe fails
- **WHEN** `GET /health/deep` is called with a valid secret and the database is unreachable
- **THEN** the system returns 503 with `"status": "degraded"` and `"db"` containing the error message; `"r2"` and `"bookleaf"` reflect their actual probe results

#### Scenario: R2 probe fails
- **WHEN** `GET /health/deep` is called with a valid secret and R2 is unreachable
- **THEN** the system returns 503 with `"status": "degraded"` and `"r2"` containing the error message; `"db"` and `"bookleaf"` reflect their actual probe results

#### Scenario: Only Bookleaf probe fails
- **WHEN** `GET /health/deep` is called with a valid secret, the database and R2 are reachable, and Bookleaf is unreachable
- **THEN** the system returns 200 with `"status": "ok"`, `"db": "ok"`, `"r2": "ok"` and `"bookleaf"` containing the error message

#### Scenario: Bookleaf fails alongside another component
- **WHEN** `GET /health/deep` is called with a valid secret and both Bookleaf and the database are unreachable
- **THEN** the system returns 503 with `"status": "degraded"`, `"db"` and `"bookleaf"` each containing their error message, and `"r2"` reflecting its actual result

### Requirement: Deep health endpoint secret
The system SHALL require the `X-Booklet-Internal-Secret` header on `GET /health/deep`, matching a dedicated `HEALTH_CHECK_SECRET` configuration value distinct from `INTERNAL_API_SECRET`. Requests with a missing or non-matching header SHALL receive 401 and SHALL NOT trigger any probe. The server SHALL refuse to start if `HEALTH_CHECK_SECRET` is not set.

#### Scenario: Missing secret header
- **WHEN** `GET /health/deep` is called without an `X-Booklet-Internal-Secret` header
- **THEN** the system returns 401 and no database, R2 or Bookleaf probe runs

#### Scenario: Wrong secret header
- **WHEN** `GET /health/deep` is called with an `X-Booklet-Internal-Secret` header that does not match `HEALTH_CHECK_SECRET`
- **THEN** the system returns 401 and no probe runs

#### Scenario: Internal API secret is not accepted
- **WHEN** `GET /health/deep` is called with an `X-Booklet-Internal-Secret` header equal to `INTERNAL_API_SECRET` while `HEALTH_CHECK_SECRET` is a different value
- **THEN** the system returns 401

#### Scenario: Health secret does not unlock internal routes
- **WHEN** `DELETE /internal/users/:id` is called with an `X-Booklet-Internal-Secret` header equal to `HEALTH_CHECK_SECRET` while `INTERNAL_API_SECRET` is a different value
- **THEN** the system returns 401 and no user data is deleted

#### Scenario: Valid secret header
- **WHEN** `GET /health/deep` is called with an `X-Booklet-Internal-Secret` header equal to `HEALTH_CHECK_SECRET`
- **THEN** the system runs the probes and responds as described by the deep health requirements

#### Scenario: Secret not configured
- **WHEN** the server is started without `HEALTH_CHECK_SECRET` set
- **THEN** configuration loading fails with an error naming `HEALTH_CHECK_SECRET` and the server does not start
