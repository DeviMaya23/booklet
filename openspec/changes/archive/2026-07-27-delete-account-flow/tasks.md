## 1. Domain & Config

- [x] 1.1 Add `IsPendingDeletion bool` field to `domain.User` with GORM column tag `is_pending_deletion`
- [x] 1.2 Write DB migration adding `is_pending_deletion boolean NOT NULL DEFAULT false` to the `users` table
- [x] 1.3 Add `BookletInternalSecret string` to `Config` struct and load it from `INTERNAL_API_SECRET` env var in `config.go`

## 2. Internal Auth Middleware

- [x] 2.1 Create `handler/middleware/internal_auth.go` with `NewInternalAuthMiddleware(secret string) echo.MiddlewareFunc` that validates `X-Booklet-Internal-Secret` header and returns 401 on mismatch

## 3. Auth Middleware Update

- [x] 3.1 Capture the returned `*domain.User` from `GetOrProvision` in the auth middleware handler
- [x] 3.2 After `GetOrProvision`, add an `if user.IsPendingDeletion` check that returns 401

## 4. Bookleaf Client

- [x] 4.1 Add `DeleteAccount(ctx context.Context, kindeUserID string) error` to `bookleaf/client.go`; call `DELETE /internal/accounts/:id` with `X-Bookleaf-Internal-Secret` header; return typed sentinel errors for 401 (`ErrUnauthorized`) and non-2xx (`ErrUnexpectedStatus`)

## 5. User Repository

- [x] 5.1 Add `SetPendingDeletion(ctx context.Context, id string) error` to `user_repository.go` — updates `is_pending_deletion = true` where `id = ?`
- [x] 5.2 Add `DeleteAllUserData(ctx context.Context, userID string) ([]string, error)` to `user_repository.go` — uses `dbFromContext` to operate on the caller's transaction; collects all R2 keys (image paths, thumbnail paths, character hero paths, pending upload keys), then hard-deletes pending_uploads, image_characters, images, character_folders, characters, and the user row (`Unscoped`); returns the collected keys

## 6. User Usecase

- [x] 6.1 Define `UserRepository` interface additions in `usecase/` for `SetPendingDeletion` and `DeleteAllUserData`
- [x] 6.2 Add `BookleafClient` interface to `usecase/` with `DeleteAccount(ctx, kindeUserID string) error`
- [x] 6.3 Implement `MarkPendingDeletion(ctx context.Context, userID string) error` in `user_usecase.go` — opens a GORM transaction, calls `SetPendingDeletion` within it, calls `bookleafClient.DeleteAccount`, commits on success or rolls back on any error; returns typed error distinguishing 401 (config error) from other failures
- [x] 6.4 Implement `PurgeUserData(ctx context.Context, userID string) ([]string, error)` in `user_usecase.go` — runs `DeleteAllUserData` in a transaction, returns the collected R2 keys on commit
- [x] 6.5 Write unit tests for `MarkPendingDeletion`:
  - Bookleaf returns success → verify commit path (no error returned)
  - Bookleaf returns 401 → verify rollback and specific config-error sentinel returned
  - Bookleaf returns non-2xx → verify rollback and appropriate error returned

## 7. Worker

- [x] 7.1 Create `worker/purge_user_storage.go` with `PurgeUserStorageArgs` (Kind: `purge_user_storage`, field `R2Keys []string`) and `PurgeUserStorageWorker`; iterate keys calling `storage.DeleteObject`; log each failure but do not return an error

## 8. User Handler

- [x] 8.1 Create `handler/user_handler.go` with `UserHandler` struct
- [x] 8.2 Implement `DeleteMe(c echo.Context) error` — extracts user ID from context, calls `userUsecase.MarkPendingDeletion`, maps errors to 502 (bookleaf non-2xx), 500 (401/config), 202 on success
- [x] 8.3 Implement `DeleteUserByID(c echo.Context) error` — extracts `:id` path param, calls `userUsecase.PurgeUserData`, enqueues `PurgeUserStorageArgs` River job with returned keys, returns 202; returns 404 if user not found
- [x] 8.4 Write unit tests for `DeleteMe`:
  - `MarkPendingDeletion` succeeds → assert 202 response
  - `MarkPendingDeletion` returns config error → assert 500 response
  - `MarkPendingDeletion` returns bookleaf error → assert 502 response
- [x] 8.5 Write unit tests for `DeleteUserByID`:
  - `PurgeUserData` succeeds → assert 202 response and River job enqueued with correct keys
  - `PurgeUserData` returns not-found error → assert 404 response

## 9. Wire Up in main.go

- [x] 9.1 Construct `userUsecase` with `bookleafClient` dependency (update `NewUserUsecase` signature and call site)
- [x] 9.2 Construct `UserHandler` and register `DELETE /me` on the existing `protected` group
- [x] 9.3 Create new internal echo group, apply `NewInternalAuthMiddleware(cfg.BookletInternalSecret)`, register `DELETE /internal/users/:id` on it
- [x] 9.4 Register `PurgeUserStorageWorker` with River workers

## 10. Bruno Collection

- [x] 10.1 Create `collection/me/delete_me.bru` for `DELETE /me` with bearer auth
- [x] 10.2 Create `collection/internal/delete_user.bru` for `DELETE /internal/users/:id` with `X-Booklet-Internal-Secret` header

## 11. Lint

- [x] 11.1 Run `golangci-lint run ./...` from `backend/` and fix any issues
