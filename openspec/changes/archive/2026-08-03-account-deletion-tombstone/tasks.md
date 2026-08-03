## 1. DB Migration

- [x] 1.1 Create `backend/migration/000008_add_account_state_to_users.up.sql`: add `account_state TEXT NOT NULL DEFAULT 'active' CHECK (account_state IN ('active', 'pending_deletion', 'purged'))`, add `purged_at TIMESTAMPTZ`, backfill `UPDATE users SET account_state = 'pending_deletion' WHERE is_pending_deletion = true`, drop `is_pending_deletion`
- [x] 1.2 Create `backend/migration/000008_add_account_state_to_users.down.sql`: restore `is_pending_deletion`, backfill from `account_state`, drop `account_state` and `purged_at`

## 2. Domain Model

- [x] 2.1 In `backend/internal/domain/user.go`, add `type AccountState string` with constants `AccountStateActive`, `AccountStatePendingDeletion`, `AccountStatePurged`
- [x] 2.2 Replace `IsPendingDeletion bool` with `AccountState AccountState` (gorm tag: `column:account_state;default:active`) and `PurgedAt *time.Time` (gorm tag: `column:purged_at`)

## 3. Repository

- [x] 3.1 In `user_repository.go`, update `SetPendingDeletion` to `UPDATE users SET account_state = 'pending_deletion'` (rename column reference only — method name and signature unchanged)
- [x] 3.2 In `user_repository.go`, update `DeleteAllUserData`: replace the final hard-delete of the user row with `UPDATE users SET account_state = 'purged', purged_at = NOW()` — keep `gorm.ErrRecordNotFound` when the row doesn't exist at all (check `RowsAffected == 0`)
- [x] 3.3 In `user_repository.go`, add `DeleteExpiredTombstones(ctx context.Context) error`: hard-deletes rows where `account_state = 'purged' AND purged_at < NOW() - INTERVAL '24 hours'`
- [x] 3.4 Add `DeleteExpiredTombstones` to the `UserRepository` interface in `usecase/user_usecase.go`

## 4. Auth Middleware

- [x] 4.1 In `handler/middleware/auth.go`, replace `user.IsPendingDeletion` check with `user.AccountState != domain.AccountStateActive`

## 5. Usecase

- [x] 5.1 In `usecase/user_usecase.go`, add `CleanupExpiredTombstones(ctx context.Context) error` method that calls `userRepo.DeleteExpiredTombstones`

## 6. Periodic Worker

- [x] 6.1 Create `backend/internal/worker/purge_tombstones.go`: define `PurgeTombstonesArgs` (Kind: `"purge_tombstones"`), `PurgeTombstonesWorker` with a `tombstoneUsecase` interface (`CleanupExpiredTombstones`), and `Work` method that delegates to the usecase
- [x] 6.2 In `cmd/server/main.go`, register `NewPurgeTombstonesWorker` via `river.AddWorker` and add a `river.NewPeriodicJob` with `24 * time.Hour` interval and `RunOnStart: false`

## 7. Unit Tests

- [x] 7.1 In `usecase/user_usecase_test.go`, update `spyUserRepository` to add `DeleteExpiredTombstones` stub
- [x] 7.2 Add `TestCleanupExpiredTombstones_CallsRepo` — verifies the usecase delegates to the repository
- [x] 7.3 In `handler/user_handler_test.go`, add `TestDeleteUserByID_Success_TombstonesUser` asserting the handler returns 202 when `PurgeUserData` returns keys (existing success test covers this but verify it still reflects tombstone behavior accurately; update if needed)
- [x] 7.4 Create `backend/internal/handler/middleware/auth_test.go` with unit tests for the updated block condition: `TestAuth_PendingDeletionAccount_Returns401` and `TestAuth_PurgedAccount_Returns401` and `TestAuth_ActiveAccount_Passes`

## 8. Repository Integration Tests

- [x] 8.1 In `user_repository_integration_test.go`, add a test that after `DeleteAllUserData` the user row remains with `account_state = 'purged'` and `purged_at` set, and all associated data is gone
- [x] 8.2 Add a test for `DeleteExpiredTombstones`: verifies rows older than 24h are deleted and rows younger than 24h are retained

## 9. Lint

- [x] 9.1 Run `golangci-lint run ./...` from `backend/` and fix any issues
