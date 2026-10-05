## MODIFIED Requirements

### Requirement: Delete artpiece
The system SHALL allow an authenticated user to delete an artpiece. An optional `delete_files` query parameter controls whether attached files are also permanently deleted. When `delete_files=true`, all files attached to the artpiece are deleted from the database and their R2 objects are enqueued for asynchronous purge; the artpiece is then deleted in the same transaction. When the parameter is absent or false, files remain with `artpiece_id` set to NULL (existing behavior).

#### Scenario: Successful deletion without delete_files
- **WHEN** a user deletes an artpiece they own without the `delete_files` query parameter
- **THEN** the system deletes the artpiece and returns HTTP 204, and all previously attached files remain with artpiece_id NULL

#### Scenario: Successful deletion with delete_files=true, artpiece has files
- **WHEN** a user deletes an artpiece they own with `?delete_files=true` and the artpiece has one or more attached files
- **THEN** the system deletes all attached file DB rows and the artpiece in a single transaction, enqueues a `purge_r2_objects` job for all collected R2 keys, and returns HTTP 204

#### Scenario: Successful deletion with delete_files=true, artpiece has no files
- **WHEN** a user deletes an artpiece they own with `?delete_files=true` and the artpiece has no attached files
- **THEN** the system deletes the artpiece and returns HTTP 204 (no file deletion or R2 job needed)

#### Scenario: R2 purge job enqueue fails
- **WHEN** `delete_files=true`, the transaction commits successfully, but enqueueing the `purge_r2_objects` job fails
- **THEN** the system returns HTTP 204, logs the enqueue error, and accepts the orphaned R2 objects as an acceptable tradeoff

#### Scenario: Artpiece not found or belongs to another user
- **WHEN** a user deletes an artpiece that does not exist or belongs to another user
- **THEN** the system returns HTTP 404
