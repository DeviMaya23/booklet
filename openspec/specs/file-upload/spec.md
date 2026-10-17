## Purpose

This capability covers the full lifecycle of file uploads: initiating an upload (generating a presigned R2 PUT URL), completing an upload (creating the file record and enqueuing post-processing), generating thumbnails for image files, and cleaning up stale pending upload records.

## Requirements

### Requirement: Initiate file upload
The system SHALL allow an authenticated user to initiate a file upload by generating a presigned PUT URL for R2, optionally associating the upload with an artpiece, notes, and a name. A pending_file_upload record is created to track the in-flight upload. For image-type files, the system SHALL also generate a presigned PUT URL for the thumbnail and return it alongside the upload URL.

#### Scenario: Successful initiation — image type
- **WHEN** a user posts a valid initiate request with an image `mime_type` and optional `artpiece_id`, `notes`, and `name`
- **THEN** the system creates a pending_file_upload record and returns a presigned PUT URL for the original file, a `thumbnail_upload_url` presigned PUT URL for the thumbnail, and an expiry

#### Scenario: Successful initiation — non-image type
- **WHEN** a user posts a valid initiate request with a non-image `mime_type` and optional `artpiece_id`, `notes`, and `name`
- **THEN** the system creates a pending_file_upload record and returns a presigned PUT URL with an expiry; `thumbnail_upload_url` is null

#### Scenario: Artpiece not owned by user
- **WHEN** a user initiates an upload with an artpiece_id that does not belong to them
- **THEN** the system returns HTTP 422

---

### Requirement: Complete file upload
The system SHALL allow an authenticated user to complete a pending file upload. On completion, the pending record is deleted and a file row is created, carrying over `name`, `notes`, `artpiece_id`, and `mime_type` from the pending record. For image-type files, `thumbnail_r2_path` is set deterministically from the known key pattern only when the client signals `thumbnail_uploaded: true`; when `false` (e.g. generation failed on the client), the path is left null. An image_metadata row is inserted when width/height are provided. Cover auto-selection runs if the file is attached to an artpiece.

#### Scenario: Successful completion — image type, thumbnail uploaded
- **WHEN** a user completes a pending upload whose mime_type starts with `image/` and `thumbnail_uploaded` is `true`
- **THEN** the file row is created (with `thumbnail_r2_path` set to `users/{user_id}/thumbnails/{file_id}.<ext>`), the pending record is deleted, and an image_metadata row is inserted; no thumbnail job is enqueued

#### Scenario: Successful completion — image type, thumbnail not uploaded
- **WHEN** a user completes a pending upload whose mime_type starts with `image/` and `thumbnail_uploaded` is `false`
- **THEN** the file row is created with `thumbnail_r2_path` null, the pending record is deleted, and an image_metadata row is inserted if dimensions were provided; no thumbnail job is enqueued

#### Scenario: Successful completion — non-image type
- **WHEN** a user completes a pending upload whose mime_type does not start with `image/`
- **THEN** the file row is created (with `name` copied from the pending record and `thumbnail_r2_path` null) and the pending record is deleted; no image_metadata row is created

#### Scenario: Pending upload not found or belongs to another user
- **WHEN** a user attempts to complete a pending upload that does not exist or belongs to another user
- **THEN** the system returns HTTP 404

---

### Requirement: Stale pending file upload cleanup
The system SHALL periodically delete pending_file_upload records older than the presign TTL (15 minutes) and remove their corresponding R2 objects.

#### Scenario: Stale upload purged
- **WHEN** a pending_file_upload record's created_at is older than 15 minutes
- **THEN** the system deletes the R2 object and removes the pending record
