# Account Deletion Spec

## Purpose

Covers the full account deletion lifecycle: a user initiates deletion via `DELETE /me`, user data is synchronously purged from the DB via an internal endpoint called by Bookleaf, and R2 storage objects are asynchronously cleaned up via a River worker job.

## Requirements

### Requirement: User can initiate account deletion
The system SHALL allow an authenticated user to request deletion of their account via `DELETE /me`. The endpoint SHALL flag the account as pending deletion and coordinate with Bookleaf to schedule Kinde account deletion. No user data SHALL be deleted at this point.

#### Scenario: Successful deletion request
- **WHEN** an authenticated user calls `DELETE /me`
- **THEN** the system opens a DB transaction, sets `is_pending_deletion = true`, calls `DELETE /internal/accounts/:id` on Bookleaf with the user's Kinde ID and the Bookleaf internal secret header, commits the transaction only on Bookleaf 202, and returns 202 to the caller

#### Scenario: Bookleaf returns non-2xx (not 401)
- **WHEN** an authenticated user calls `DELETE /me` and Bookleaf returns a non-2xx, non-401 response
- **THEN** the system rolls back the transaction (user remains unflagged), logs the response, and returns 502

#### Scenario: Bookleaf returns 401
- **WHEN** an authenticated user calls `DELETE /me` and Bookleaf returns 401
- **THEN** the system rolls back the transaction, logs a configuration error (shared secret mismatch), and returns 500

### Requirement: Booklet exposes an internal endpoint for user data purge
The system SHALL expose `DELETE /internal/users/:id` protected by the `X-Booklet-Internal-Secret` header. When called, it SHALL synchronously delete all user data from the DB and asynchronously delete all user R2 objects. The endpoint SHALL NOT check `is_pending_deletion` status.

#### Scenario: Successful data purge
- **WHEN** Bookleaf calls `DELETE /internal/users/:id` with a valid `X-Booklet-Internal-Secret` header and the user exists
- **THEN** the system collects all R2 keys belonging to the user, hard-deletes all DB records for the user in a single transaction (pending_uploads, image_characters, images, character_folders, characters, user row), enqueues a `PurgeUserStorageArgs` River job with the collected R2 keys, and returns 202

#### Scenario: User not found
- **WHEN** Bookleaf calls `DELETE /internal/users/:id` and no user with that ID exists in the DB
- **THEN** the system returns 404

#### Scenario: Missing or invalid internal secret
- **WHEN** `DELETE /internal/users/:id` is called without `X-Booklet-Internal-Secret` header or with an incorrect value
- **THEN** the system returns 401

### Requirement: R2 objects are asynchronously purged after account deletion
The system SHALL process `PurgeUserStorageArgs` River jobs to delete R2 objects belonging to a deleted user. Deletion SHALL be attempted for each key individually. Failures SHALL be logged but SHALL NOT cause the job to fail.

#### Scenario: All R2 keys deleted successfully
- **WHEN** a `PurgeUserStorageArgs` job is processed with a list of R2 keys
- **THEN** the worker calls `DeleteObject` for each key and completes without error

#### Scenario: One or more R2 key deletions fail
- **WHEN** a `PurgeUserStorageArgs` job is processed and one or more `DeleteObject` calls return an error
- **THEN** the worker logs each failure and continues processing remaining keys, completing the job without returning an error

### Requirement: Internal endpoints are protected by shared secret
The system SHALL validate the `X-Booklet-Internal-Secret` header on all routes in the internal echo group. The expected value SHALL be read from the `BOOKLET_INTERNAL_SECRET` environment variable.

#### Scenario: Valid secret
- **WHEN** a request to an internal endpoint includes `X-Booklet-Internal-Secret` matching `BOOKLET_INTERNAL_SECRET`
- **THEN** the request proceeds to the handler

#### Scenario: Invalid or missing secret
- **WHEN** a request to an internal endpoint includes a wrong or absent `X-Booklet-Internal-Secret` header
- **THEN** the system returns 401 before reaching the handler
