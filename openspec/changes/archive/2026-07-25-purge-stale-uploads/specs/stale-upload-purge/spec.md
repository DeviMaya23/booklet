# Stale Upload Purge

## Purpose

Defines the background job that periodically finds and hard-deletes pending uploads that were never completed. A pending upload is considered stale once its age exceeds the presign TTL, meaning the client's upload window has closed and the record will never be completed.

---

## Requirements

### Requirement: Stale identification
A `PendingUpload` record SHALL be considered stale when its `created_at` is older than the configured threshold (equal to the presign TTL: 15 minutes). The cutoff is computed as `now - threshold` at the time the job runs.

#### Scenario: Record within threshold is not purged
- **WHEN** the purge job runs and a `pending_upload` row has `created_at` newer than `now - threshold`
- **THEN** that record is not deleted

#### Scenario: Record beyond threshold is purged
- **WHEN** the purge job runs and a `pending_upload` row has `created_at` older than `now - threshold`
- **THEN** that record is deleted from the database and its R2 object is deleted

---

### Requirement: Deletion order
For each stale pending upload, the system SHALL delete the R2 object before deleting the database record. This ensures that a failure at any point leaves the system in a state where the next retry can successfully clean up both resources.

#### Scenario: R2 deletion precedes DB deletion
- **WHEN** the purge job processes a stale record
- **THEN** the R2 object deletion is attempted before the database record is removed

#### Scenario: R2 object already absent
- **WHEN** the purge job attempts to delete an R2 object that does not exist (e.g. the client never uploaded the file)
- **THEN** the absent object is treated as a successful deletion and the database record is still removed

---

### Requirement: Cleanup logging
The usecase SHALL log the number of records deleted after each purge run, including runs where zero records were deleted.

#### Scenario: Records deleted
- **WHEN** the purge job completes and deleted one or more records
- **THEN** a log entry is emitted at Info level with the count of deleted records

#### Scenario: No stale records found
- **WHEN** the purge job runs and finds no stale records
- **THEN** a log entry is emitted at Info level with a count of zero

---

### Requirement: Retry on failure
If the purge job fails for any record (R2 error or DB error), the job SHALL be retried by the River worker. The retry is safe because R2 deletion is idempotent (missing object = success) and a DB delete on an already-deleted row affects zero rows without error.

#### Scenario: Failure triggers retry
- **WHEN** the purge job encounters an error during processing
- **THEN** River retries the job according to its configured retry policy
