# Account Deletion Tombstone Spec

## Purpose

Covers tombstone behavior for deleted accounts: after all user data is purged, the user row is retained with `account_state = 'purged'` to prevent re-provisioning during the 24h JWT TTL window, and periodically cleaned up once the TTL has expired.

## Requirements

### Requirement: Deleted accounts are tombstoned until JWT TTL expires
After all user data is purged, the user row SHALL remain in the database with `account_state = 'purged'` and `purged_at` set to the purge timestamp. The row SHALL NOT be hard-deleted at purge time. This tombstone prevents a still-valid Kinde JWT from re-provisioning the account during the 24h JWT TTL window.

#### Scenario: Tombstone row exists after data purge
- **WHEN** `DELETE /internal/users/:id` is called and the purge succeeds
- **THEN** the user row remains in the database with `account_state = 'purged'` and `purged_at` set to the current timestamp, while all associated data (images, characters, character_folders, image_characters, pending_uploads) is deleted

### Requirement: Expired tombstone rows are hard-deleted by a daily periodic job
The system SHALL run a daily River periodic job that hard-deletes user rows where `account_state = 'purged'` AND `purged_at < NOW() - INTERVAL '24 hours'`. This permanently removes tombstones after the JWT TTL has expired.

#### Scenario: Tombstone older than 24 hours is cleaned up
- **WHEN** the periodic cleanup job runs and a user row has `account_state = 'purged'` and `purged_at` older than 24 hours
- **THEN** the job hard-deletes that user row

#### Scenario: Tombstone younger than 24 hours is retained
- **WHEN** the periodic cleanup job runs and a user row has `account_state = 'purged'` and `purged_at` within the last 24 hours
- **THEN** the job does not delete that row
