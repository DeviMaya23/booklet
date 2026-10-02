## Purpose

Allows a user to delete multiple files in a single request. DB rows are deleted synchronously; R2 object removal is handled asynchronously via a background job.

## Requirements

### Requirement: Bulk delete files
The system SHALL allow a user to delete multiple files in a single request. The DB rows are deleted synchronously; R2 object removal (file paths and thumbnail paths) is handled asynchronously via a background job.

#### Scenario: Successful bulk delete
- **WHEN** a user sends `DELETE /files` with a body containing a list of file IDs they own
- **THEN** the system deletes all matching DB rows, enqueues a single `purge_r2_objects` job carrying all collected R2 keys, and returns 204

#### Scenario: One or more IDs not owned by the user
- **WHEN** a user sends `DELETE /files` with a body that includes file IDs belonging to another user or that do not exist
- **THEN** the system returns 422 and deletes nothing

#### Scenario: Empty ID list
- **WHEN** a user sends `DELETE /files` with an empty `ids` array
- **THEN** the system returns 422

#### Scenario: R2 purge job enqueue fails
- **WHEN** the bulk delete DB operation succeeds but enqueueing the `purge_r2_objects` job fails
- **THEN** the system still returns 204, logs the enqueue error, and accepts the orphaned R2 objects as an acceptable tradeoff
