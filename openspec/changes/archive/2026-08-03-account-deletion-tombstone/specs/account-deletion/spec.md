## MODIFIED Requirements

### Requirement: User can initiate account deletion
The system SHALL allow an authenticated user to request deletion of their account via `DELETE /me`. The endpoint SHALL flag the account as pending deletion and coordinate with Bookleaf to schedule Kinde account deletion. No user data SHALL be deleted at this point.

#### Scenario: Successful deletion request
- **WHEN** an authenticated user calls `DELETE /me`
- **THEN** the system opens a DB transaction, sets `account_state = 'pending_deletion'`, calls `DELETE /internal/accounts/:id` on Bookleaf with the user's Kinde ID and the Bookleaf internal secret header, commits the transaction only on Bookleaf 202, and returns 202 to the caller

#### Scenario: Bookleaf returns non-2xx (not 401)
- **WHEN** an authenticated user calls `DELETE /me` and Bookleaf returns a non-2xx, non-401 response
- **THEN** the system rolls back the transaction (user remains in `active` state), logs the response, and returns 502

#### Scenario: Bookleaf returns 401
- **WHEN** an authenticated user calls `DELETE /me` and Bookleaf returns 401
- **THEN** the system rolls back the transaction, logs a configuration error (shared secret mismatch), and returns 500

### Requirement: Booklet exposes an internal endpoint for user data purge
The system SHALL expose `DELETE /internal/users/:id` protected by the `X-Booklet-Internal-Secret` header. When called, it SHALL synchronously delete all user data from the DB, tombstone the user row, and asynchronously delete all user R2 objects. The endpoint SHALL NOT check `account_state`.

#### Scenario: Successful data purge
- **WHEN** Bookleaf calls `DELETE /internal/users/:id` with a valid `X-Booklet-Internal-Secret` header and the user exists
- **THEN** the system collects all R2 keys belonging to the user, hard-deletes all associated DB records (pending_uploads, image_characters, images, character_folders, characters) in a single transaction, sets `account_state = 'purged'` and `purged_at = NOW()` on the user row, enqueues a `PurgeUserStorageArgs` River job with the collected R2 keys, and returns 202

#### Scenario: User not found
- **WHEN** Bookleaf calls `DELETE /internal/users/:id` and no user with that ID exists in the DB (neither active nor tombstoned)
- **THEN** the system returns 404

#### Scenario: Missing or invalid internal secret
- **WHEN** `DELETE /internal/users/:id` is called without `X-Booklet-Internal-Secret` header or with an incorrect value
- **THEN** the system returns 401
