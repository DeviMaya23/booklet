## Context

Currently, `DELETE /internal/users/:id` hard-deletes the user row after purging all associated data. Kinde JWTs have a 24h TTL. In the window between the hard-delete and JWT expiry, a request with a still-valid JWT will pass JWT validation and hit `GetOrProvision`, which re-creates the user row (via `GetOrCreate`) because it can no longer find the deleted row. The account is effectively un-deleted.

The fix is a tombstone: keep the user row alive after data purge, in a `purged` state, so the auth middleware can still find and block it. Once 24h have passed (JWT TTL is exhausted), a periodic job hard-deletes the row.

## Goals / Non-Goals

**Goals:**
- Close the JWT re-provisioning window after account deletion
- Replace the `is_pending_deletion` boolean with a typed `account_state` column covering the full lifecycle
- Keep the tombstone row until it is safe to remove (past JWT TTL)

**Non-Goals:**
- Changing how or when Bookleaf schedules Kinde deletion
- Adding retry logic or outbox patterns to the tombstone cleanup job
- Blocking non-authenticated (internal) endpoints based on account state

## Decisions

### Replace `is_pending_deletion` with `account_state`

`is_pending_deletion` is a partial model — it only tracks one transition. As we add a second terminal state (`purged`), a boolean becomes semantically incorrect (a purged account is not "pending" anything). A typed enum column is cleaner and future-proof.

**Migration strategy**: single migration — add `account_state TEXT NOT NULL DEFAULT 'active'`, backfill `WHERE is_pending_deletion = true → 'pending_deletion'`, drop `is_pending_deletion`. No data loss, no need for a two-phase migration.

**Alternative considered**: keeping `is_pending_deletion` and adding a separate `is_purged` column. Rejected — two boolean columns encoding a state machine is worse than one enum.

### `DELETE /internal/users/:id` tombstones instead of hard-deleting

After `DeleteAllUserData` strips the user's associated records, the user row itself is updated to `account_state = 'purged'`, `purged_at = NOW()` instead of being deleted. The row is a ghost: no data, but present to block re-provisioning.

`DeleteAllUserData` in the repository returns `gorm.ErrRecordNotFound` if the user row doesn't exist (currently, this is the 404 path). With tombstone, we need to keep that 404 path — if a user arrives at the internal endpoint that doesn't exist at all (neither active nor tombstoned), return 404. If the user is already `purged`, returning 404 is also acceptable since the data is already gone.

### Auth middleware checks `AccountState != Active`

Currently the middleware checks `user.IsPendingDeletion`. Replacing this with `user.AccountState != AccountStateActive` covers all non-active states in one condition. Both `pending_deletion` and `purged` rows block access.

**GetOrProvision behavior**: the middleware calls `GetOrProvision`, which tries `GetByID` first and only creates if not found. During the tombstone window (24h), the purged row is found, returned, and blocked by the state check. After the tombstone cleanup job removes the row (>24h), any JWT for that user is also expired by then — so `GetOrProvision` will never reach `GetOrCreate` for a legitimately deleted user.

### Periodic cleanup job runs once daily via River cron

Follows the same pattern as `PurgeExpiredUploadsWorker`. The job calls a usecase method that hard-deletes rows where `account_state = 'purged' AND purged_at < NOW() - INTERVAL '24 hours'`. Daily cadence is sufficient given the 24h TTL alignment.

**River schedule**: `river.PeriodicInterval(24 * time.Hour)` with `RunOnStart: false` (no need to run at boot).

### `type AccountState string` with typed constants

Mirrors the bookleaf pattern. Three constants: `AccountStateActive`, `AccountStatePendingDeletion`, `AccountStatePurged`. GORM stores as text with a `CHECK` constraint in the migration. Go code never uses raw strings.

## Risks / Trade-offs

- **Tombstone row retention**: If the periodic job fails or misses a run, tombstone rows linger past 24h. This is safe (access remains blocked) but means stale rows accumulate. River's built-in retry and at-least-once guarantee makes silent failure unlikely.
- **`DeleteAllUserData` behavior change**: the function currently hard-deletes the user row as its last step. After this change, it updates the row instead. The 404 contract (`gorm.ErrRecordNotFound` when user doesn't exist) is preserved by explicitly checking row existence before the purge, or by treating the "already purged" case as a 404 at the handler level.

## Migration Plan

1. Write migration `000008_replace_is_pending_deletion_with_account_state`:
   - Add `account_state TEXT NOT NULL DEFAULT 'active' CHECK (account_state IN ('active', 'pending_deletion', 'purged'))`
   - Add `purged_at TIMESTAMPTZ`
   - `UPDATE users SET account_state = 'pending_deletion' WHERE is_pending_deletion = true`
   - Drop `is_pending_deletion`
2. Update `domain.User` (field swap)
3. Update repository, usecase, middleware, worker in any order (all compile-break simultaneously after domain change)
4. Register new periodic worker in `main.go`

Rollback: a down migration that re-adds `is_pending_deletion`, backfills from `account_state`, and drops the new columns.
