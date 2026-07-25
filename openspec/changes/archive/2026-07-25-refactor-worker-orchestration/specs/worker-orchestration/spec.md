# Worker Orchestration

## Purpose

Defines the structural contract for how River workers are registered and wired in this codebase. This is a structural spec — no observable system behaviour changes.

---

## Requirements

### Requirement: Worker registration at composition root
Worker registration (adding workers, wiring periodic jobs) and River infrastructure setup (migrations, client creation) SHALL be performed in `cmd/server/main.go` via a local `initRiverClient` function. The `internal/worker` package SHALL contain only worker types and their constructors — no River client creation, no migration logic, no `New` package-level function.

#### Scenario: Adding a new worker requires no change to initRiverClient
- **WHEN** a new background worker is introduced to the codebase
- **THEN** `initRiverClient`'s signature is unchanged; only `main.go`'s worker registration block and the new worker file are modified

### Requirement: Per-worker constructors
Each worker type SHALL expose a constructor function in its own file within `internal/worker/`. The constructor SHALL accept the worker's dependencies and return the concrete worker struct.

#### Scenario: Worker constructed via its own constructor
- **WHEN** a worker is registered in main.go
- **THEN** it is instantiated via its package-level constructor (e.g. `worker.NewPurgeExpiredUploadsWorker(...)`) before being passed to `river.AddWorker`
