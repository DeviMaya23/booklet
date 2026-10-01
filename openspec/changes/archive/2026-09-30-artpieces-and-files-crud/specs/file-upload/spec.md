## ADDED Requirements

### Requirement: Initiate file upload
The system SHALL allow an authenticated user to initiate a file upload by generating a presigned PUT URL for R2, optionally associating the upload with an artpiece and notes. A pending_file_upload record is created to track the in-flight upload.

#### Scenario: Successful initiation
- **WHEN** a user posts a valid initiate request with a mime_type and optional artpiece_id and notes
- **THEN** the system creates a pending_file_upload record and returns a presigned PUT URL with an expiry

#### Scenario: Artpiece not owned by user
- **WHEN** a user initiates an upload with an artpiece_id that does not belong to them
- **THEN** the system returns HTTP 422

---

### Requirement: Complete file upload
The system SHALL allow an authenticated user to complete a pending file upload. On completion, the pending record is deleted and a file row is created. For image-type files, an image_metadata row is inserted and a thumbnail generation job is enqueued. Cover auto-selection runs if the file is attached to an artpiece.

#### Scenario: Successful completion — image type
- **WHEN** a user completes a pending upload whose mime_type starts with `image/`
- **THEN** the file row is created, the pending record is deleted, an image_metadata row is inserted (initially with zero dimensions until the worker runs), and a generate_file_thumbnail job is enqueued

#### Scenario: Successful completion — non-image type
- **WHEN** a user completes a pending upload whose mime_type does not start with `image/`
- **THEN** the file row is created and the pending record is deleted; no thumbnail job is enqueued and no image_metadata row is created

#### Scenario: Pending upload not found or belongs to another user
- **WHEN** a user attempts to complete a pending upload that does not exist or belongs to another user
- **THEN** the system returns HTTP 404

---

### Requirement: Generate file thumbnail
The system SHALL generate a 600×600 JPEG thumbnail for image-type files after upload completion and store it in R2. The file's thumbnail_r2_path is updated once the thumbnail is written.

#### Scenario: Successful thumbnail generation
- **WHEN** the generate_file_thumbnail worker processes a job for an image file
- **THEN** the system downloads the original from R2, resizes to fit within 600×600, encodes as JPEG at 85% quality, uploads the thumbnail to R2 at `users/{user_id}/thumbnails/{file_id}.jpg`, and updates the file's thumbnail_r2_path

---

### Requirement: Stale pending file upload cleanup
The system SHALL periodically delete pending_file_upload records older than the presign TTL (15 minutes) and remove their corresponding R2 objects.

#### Scenario: Stale upload purged
- **WHEN** a pending_file_upload record's created_at is older than 15 minutes
- **THEN** the system deletes the R2 object and removes the pending record
