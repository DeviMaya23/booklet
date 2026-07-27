## Context

Booklet and Bookleaf share a Kinde identity ecosystem. Bookleaf is the primary orchestrator for Kinde account deletion (it already handles deleting the Kinde user). Booklet holds its own DB records (users, characters, images, pending uploads) and R2 objects (image files, thumbnails, character hero images, stale upload keys).

Two deletion paths must be supported:
1. **User-triggered**: User calls `DELETE /me` from the FE → Booklet coordinates with Bookleaf → Bookleaf schedules Kinde deletion → Bookleaf calls back Booklet to purge data
2. **Bookleaf-triggered**: Bookleaf independently triggers deletion of Booklet data via `DELETE /internal/users/:id`

Both paths converge on the same purge logic. The difference is who initiates it.

## Goals / Non-Goals

**Goals:**
- Users can initiate account deletion from Booklet's FE
- All Booklet user data (DB rows + R2 objects) is deleted when Bookleaf triggers the purge
- Accounts in the deletion limbo state are blocked from further API usage (401)
- Internal endpoints are protected by a shared secret, not JWT

**Non-Goals:**
- Kinde user deletion (owned by Bookleaf)
- Bookleaf data deletion (owned by Bookleaf)
- Soft-delete or data export before deletion
- Retry orchestration for the Bookleaf callback (Bookleaf owns that)

## Decisions

### 1. Limbo state via `is_pending_deletion` boolean, not `DeletedAt`

GORM's `DeletedAt` (soft delete) auto-filters records from queries. Using it for the limbo state would cause `GetOrProvision` to miss the user and re-create them. A dedicated `IsPendingDeletion bool` keeps the semantics distinct: "flagged for future deletion" vs. "actually deleted."

The auth middleware checks this flag after `GetOrProvision` and returns 401. This mirrors Bookleaf's existing pattern.

### 2. Flag-then-commit only on Bookleaf 202

`DELETE /me` opens a GORM transaction, sets `is_pending_deletion = true` within it, calls Bookleaf, and only commits if Bookleaf returns 202. Any non-2xx rolls back the transaction, leaving the user unflagged and able to retry.

Alternative considered: flag first, unflag on failure. Rejected because a crash between flag and unflag would leave the user permanently locked out with no recourse.

Residual risk: Booklet crashes after committing the transaction but before responding to the FE. The user is flagged in DB but the FE sees an error. On retry the middleware blocks them (401). Acceptable — the account deletion is in progress at that point.

### 3. `DELETE /me` stays in the JWT-protected group; `DELETE /internal/users/:id` gets its own group

`DELETE /me` is user-initiated and requires identity. It stays in the existing `protected` echo group.

`DELETE /internal/users/:id` is called by Bookleaf server-to-server and has no JWT. It gets a dedicated echo group with a new `InternalAuthMiddleware` that validates the `X-Booklet-Internal-Secret` header against `BOOKLET_INTERNAL_SECRET` env var. This mirrors the `X-Bookleaf-Internal-Secret` pattern already used in the Bookleaf client.

### 4. Synchronous DB purge, async R2 cleanup

`DELETE /internal/users/:id`:
1. Collect all R2 keys from DB (images, thumbnails, character hero images, pending upload keys)
2. Hard-delete all DB records in a transaction (pending_uploads → image_characters → images → character_folders → characters → user, using `Unscoped()` for the user row to bypass GORM soft delete)
3. Commit
4. Enqueue a `PurgeUserStorageArgs` River job with the collected R2 keys
5. Return 202

The user's Kinde account will already be deleted by Bookleaf before this endpoint is called, so there is no risk of re-auth after the user row is hard-deleted.

R2 cleanup is best-effort. Failures are logged; orphaned R2 objects are acceptable (no PII, no functional impact). This follows the same pattern as the existing `PurgeExpiredUploadsWorker`.

### 5. Hard delete for the user row

The user row is deleted with `Unscoped().Delete()` to bypass GORM's soft-delete. A soft-delete tombstone would interfere with `GetOrProvision` (could re-create the user row on a future auth attempt with the same Kinde ID, even though the Kinde account is gone). Since Bookleaf deletes the Kinde user before calling this endpoint, re-registration with the same ID is not expected but hard delete is the safer default.

### 6. R2 key collection inside the transaction

All R2 keys are collected from DB inside the same transaction as the deletions. The collection queries run on the same snapshot as the deletes, so the returned key set is always consistent with what was removed. This keeps `DeleteAllUserData` self-contained and avoids passing pre-collected key lists across layer boundaries.

## Risks / Trade-offs

- **Booklet crash after commit, before enqueuing River job** → R2 objects are orphaned. Mitigation: acceptable per the best-effort policy on R2 cleanup; objects contain no PII beyond the image content itself.
- **HTTP call inside a DB transaction** (`DELETE /me` calls Bookleaf while holding a connection) → connection held for the duration of the network round-trip. Mitigation: account deletion is a low-frequency, user-initiated operation; the risk is negligible.
- **Bookleaf callback arrives before `DELETE /me` commits** → `DELETE /internal/users/:id` runs on a user who is not yet flagged. The endpoint does not check `is_pending_deletion`; it purges unconditionally. Mitigation: not a problem — the purge is idempotent and correct regardless of flag state.

## Migration Plan

1. Deploy DB migration adding `is_pending_deletion` column (default false, not null) — safe, additive
2. Deploy application code with new endpoints and middleware
3. Add `BOOKLET_INTERNAL_SECRET` to environment before deploy
4. No rollback concerns: the new column is ignored by existing code paths
