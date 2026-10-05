## MODIFIED Requirements

### Requirement: Booklet exposes an internal endpoint for user data purge
The system SHALL expose `DELETE /internal/users/:id` protected by the `X-Booklet-Internal-Secret` header. When called, it SHALL synchronously delete all user data from the DB, tombstone the user row, and asynchronously delete all user R2 objects. The endpoint SHALL NOT check `account_state`.

The purge SHALL collect R2 keys and delete DB records for: character avatar R2 paths, character folders, and characters. It SHALL NOT collect or delete image R2 keys, pending-upload R2 keys, `image_characters` rows, `images` rows, or `pending_uploads` rows (those tables no longer exist).

#### Scenario: Successful data purge
- **WHEN** Bookleaf calls `DELETE /internal/users/:id` with a valid `X-Booklet-Internal-Secret` header and the user exists
- **THEN** the system collects all character avatar R2 keys belonging to the user, hard-deletes character folders and characters in a single transaction, sets `account_state = 'purged'` and `purged_at = NOW()` on the user row, enqueues a `PurgeUserStorageArgs` River job with the collected R2 keys, and returns 202

#### Scenario: User not found
- **WHEN** Bookleaf calls `DELETE /internal/users/:id` and no user with that ID exists in the DB
- **THEN** the system returns 404

#### Scenario: Missing or invalid internal secret
- **WHEN** `DELETE /internal/users/:id` is called without `X-Booklet-Internal-Secret` header or with an incorrect value
- **THEN** the system returns 401
