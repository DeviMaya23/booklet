## Context

`worker.New` currently accepts domain-specific dependencies (`cleanupUsecase`, `threshold`) alongside infrastructure ones (`pool`, `logger`). The River client requires `Workers` and `PeriodicJobs` to be known at construction time, so all worker wiring is currently bundled into `New`. Adding a second worker would require adding its dependencies to `New`'s signature.

The sister project establishes the target pattern: `main.go` is the composition root for worker registration — the same role it already plays for HTTP route registration. River infrastructure (migrations, client creation) becomes a local `init*` function in `main.go`, consistent with how `initDB` is structured.

**Current shape:**
```
worker.New(ctx, pool, cleanupUsecase, threshold, logger)
  └─ runs migrations
  └─ registers PurgeExpiredUploadsWorker
  └─ wires PeriodicJob
  └─ creates river.Client
```

**Target shape:**
```
main.go
  └─ workers := river.NewWorkers()
  └─ river.AddWorker(workers, worker.NewPurgeExpiredUploadsWorker(uc, threshold))
  └─ periodicJobs := []*river.PeriodicJob{ ... }
  └─ initRiverClient(ctx, pool, workers, periodicJobs, logger)  ← local func, mirrors initDB
       └─ runs migrations
       └─ creates river.Client

internal/worker/
  └─ purge_upload.go  ← worker type + constructor only; no New(), no worker.go
```

## Goals / Non-Goals

**Goals:**
- `internal/worker` contains only worker types and constructors — no infrastructure wiring
- Worker registration in `main.go` mirrors HTTP handler registration: one call per worker type
- Each worker file exposes a constructor that encapsulates its own dependencies

**Non-Goals:**
- Changing stale-upload-purge behaviour in any way
- Introducing `usecase/job_args.go` or `usecase/job_enqueuer.go` — those belong when the first dispatch-style (explicitly enqueued) worker arrives
- Moving `PurgeExpiredUploadsArgs` or `cleanupUsecase` out of `worker/purge_upload.go` — the args struct has no external consumers, so it stays local

## Decisions

### 1. River infrastructure moves to `initRiverClient` in `main.go`; `worker.New` is removed

Rather than keeping `worker.New` with a changed signature, River setup (migrations + client creation) becomes `initRiverClient` — a local function in `main.go` alongside `initDB`. This mirrors the existing convention: GORM setup is not abstracted behind a `repository.New`; it lives in `main.go`. `internal/worker` is now purely worker definitions with no infrastructure concern.

**Alternative considered:** keeping `worker.New` with an updated signature accepting pre-built `*river.Workers` and `[]*river.PeriodicJob`. Rejected — `main.go` already imports `river` directly for `NewWorkers`/`AddWorker`/`PeriodicJob`; pulling `rivermigrate` and `riverpgxv5` into `main.go` is not meaningfully new coupling, and the `initDB` analogy makes the local function the more consistent choice.

### 2. Each worker file exposes a constructor

`NewPurgeExpiredUploadsWorker(uc cleanupUsecase, threshold time.Duration) *PurgeExpiredUploadsWorker` — `main.go` calls the constructor and passes the result to `river.AddWorker`. This keeps worker-specific wiring (dependencies, config) inside the worker package.

### 3. `cleanupUsecase` interface and `PurgeExpiredUploadsArgs` stay in `worker/purge_upload.go`

The purge job is periodic — River enqueues it automatically, no usecase code ever constructs `PurgeExpiredUploadsArgs{}`. There is no `usecase → worker` import needed. The worker package remains independent of the `usecase` package, consistent with how handlers are structured (each defines its own narrow interface; no import of the `usecase` package).

## Risks / Trade-offs

**`main.go` grows slightly** → Acceptable. It is already the composition root. Adding two lines per worker (constructor call + `river.AddWorker`) is the same pattern as adding a handler.

**Periodic job config (interval, `RunOnStart`) moves to `main.go`** → This is intentional — it's configuration, not worker logic. Worker files own behaviour; `main.go` owns schedule.
