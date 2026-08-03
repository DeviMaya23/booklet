## Why

After `DELETE /internal/users/:id` hard-deletes a user row, a still-valid Kinde JWT (24h TTL) can re-provision the account via `GetOrProvision` — effectively undoing deletion. A tombstone closes this window by keeping the user row alive (with all data stripped) until the JWT TTL has passed.

## What Changes

- Replace `is_pending_deletion` boolean column on `users` table with `account_state TEXT` enum (`active`, `pending_deletion`, `purged`) and add `purged_at TIMESTAMPTZ`
- `DELETE /internal/users/:id` no longer hard-deletes the user row; instead sets `account_state = 'purged'` + `purged_at = NOW()` after purging all user data
- Auth middleware checks `account_state != 'active'` (replaces `IsPendingDeletion` check) — any non-active state blocks the request with 401
- New daily River periodic job hard-deletes tombstone rows where `account_state = 'purged' AND purged_at < NOW() - 24h`
- `domain.User` replaces `IsPendingDeletion bool` with `AccountState AccountState` and `PurgedAt *time.Time`

## Capabilities

### New Capabilities

- `account-deletion-tombstone`: Tombstone lifecycle for deleted accounts — `account_state` transitions, purged row retention, and periodic cleanup

### Modified Capabilities

- `account-deletion`: The internal purge endpoint behavior changes (no longer hard-deletes the user row); `pending_deletion` flag mechanism replaced by `account_state`
- `user-auth`: Auth middleware block condition changes from `IsPendingDeletion` check to `account_state != 'active'`

## Impact

- **DB migration**: drop `is_pending_deletion`, add `account_state` (backfill from old column) + `purged_at`
- **`domain.User`**: field changes
- **`repository/user_repository.go`**: `SetPendingDeletion` and `DeleteAllUserData` both change behavior
- **`handler/middleware/auth.go`**: block condition updated
- **`usecase/user_usecase.go`**: `MarkPendingDeletion` unaffected in structure; `PurgeUserData` result changes (no longer returns keys from a deleted row, but tombstones it)
- **New River worker**: periodic tombstone cleanup job
- **`main.go`**: register new periodic worker
- **`openspec/specs/account-deletion/spec.md`**: update to reflect new endpoint behavior and state model
- **`openspec/specs/user-auth/spec.md`**: update middleware block condition
