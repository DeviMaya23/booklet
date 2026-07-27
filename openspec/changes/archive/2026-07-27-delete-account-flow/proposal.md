## Why

Booklet and Bookleaf share a Kinde identity ecosystem; deleting an account in one app must erase data in both. Currently there is no deletion flow — users have no way to remove their data, and no cross-app coordination exists.

## What Changes

- Add `is_pending_deletion` flag to the `User` domain entity
- Auth middleware returns 401 for users with `is_pending_deletion = true`
- New `DELETE /me` endpoint (JWT-protected): flags the account and coordinates with Bookleaf to schedule Kinde deletion
- New `DELETE /internal/users/:id` endpoint (X-Booklet-Internal-Secret auth): hard-deletes all user data from DB and enqueues R2 cleanup
- New `X-Booklet-Internal-Secret` internal secret middleware and `BOOKLET_INTERNAL_SECRET` env var
- New Bookleaf client method `DeleteAccount` calling `DELETE /internal/accounts/:id`
- New River job `PurgeUserStorage` that deletes R2 objects for a deleted user

## Capabilities

### New Capabilities

- `account-deletion`: User-initiated and Bookleaf-triggered account deletion, covering the limbo state, cross-app coordination, DB purge, and async R2 cleanup

### Modified Capabilities

- `user-auth`: Auth middleware gains an `is_pending_deletion` gate that returns 401 for accounts in the deletion limbo state

## Impact

- **domain**: `User` gains `IsPendingDeletion bool`
- **handler/middleware**: Auth middleware checks new flag after `GetOrProvision`
- **handler**: Two new handlers (`UserHandler` or similar) for `DELETE /me` and `DELETE /internal/users/:id`
- **usecase**: New user usecase methods for flagging and purging user data
- **repository**: New user repository methods (flag, hard-delete); new queries to collect R2 paths across images, characters, pending uploads
- **bookleaf/client**: New `DeleteAccount` method
- **worker**: New `PurgeUserStorageWorker` and `PurgeUserStorageArgs`
- **config**: `BOOKLET_INTERNAL_SECRET` env var added
- **main.go**: Register new routes, worker, and internal middleware
