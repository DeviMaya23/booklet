## MODIFIED Requirements

### Requirement: Complete file upload
The system SHALL allow an authenticated user to complete a pending file upload. On completion, the pending record is deleted and a file row is created, carrying over `name`, `notes`, `artpiece_id`, and `mime_type` from the pending record. The file row is created with `thumbnail_gen_state` set based on mime type. For image-type files, an image_metadata row is inserted and a thumbnail generation job is enqueued. Cover auto-selection runs if the file is attached to an artpiece.

#### Scenario: Successful completion — image type
- **WHEN** a user completes a pending upload whose mime_type starts with `image/`
- **THEN** the file row is created (with `name` copied from the pending record and `thumbnail_gen_state = "pending"`), the pending record is deleted, an image_metadata row is inserted, and a generate_file_thumbnail job is enqueued

#### Scenario: Successful completion — non-image type
- **WHEN** a user completes a pending upload whose mime_type does not start with `image/`
- **THEN** the file row is created (with `name` copied from the pending record and `thumbnail_gen_state = "not_applicable"`) and the pending record is deleted; no thumbnail job is enqueued and no image_metadata row is created

#### Scenario: Pending upload not found or belongs to another user
- **WHEN** a user attempts to complete a pending upload that does not exist or belongs to another user
- **THEN** the system returns HTTP 404

---

### Requirement: Generate file thumbnail
The system SHALL generate a 600×600 JPEG thumbnail for image-type files after upload completion and store it in R2. The file's `thumbnail_r2_path` and `thumbnail_gen_state` are updated once the thumbnail is written. On final retry failure the worker SHALL update `thumbnail_gen_state` to `"failed"` and stop retrying.

#### Scenario: Successful thumbnail generation
- **WHEN** the generate_file_thumbnail worker processes a job for an image file
- **THEN** the system downloads the original from R2, resizes to fit within 600×600, encodes as JPEG at 85% quality, uploads the thumbnail to R2 at `users/{user_id}/thumbnails/{file_id}.jpg`, updates the file's `thumbnail_r2_path`, and sets `thumbnail_gen_state = "done"`

#### Scenario: Thumbnail generation fails on final attempt
- **WHEN** the generate_file_thumbnail worker encounters an error and `job.Attempt >= job.MaxAttempts`
- **THEN** the system sets the file's `thumbnail_gen_state = "failed"` and returns nil so River does not re-enqueue the job

#### Scenario: Thumbnail generation fails before final attempt
- **WHEN** the generate_file_thumbnail worker encounters an error and `job.Attempt < job.MaxAttempts`
- **THEN** the system returns the error so River retries the job; `thumbnail_gen_state` remains `"pending"`
