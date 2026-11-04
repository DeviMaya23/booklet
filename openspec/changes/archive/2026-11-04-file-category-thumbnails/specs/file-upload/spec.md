# Spec Delta

## MODIFIED Requirements

### Requirement: Initiate file upload
The system SHALL allow an authenticated user to initiate a file upload by generating a presigned PUT URL for R2, optionally associating the upload with an artpiece, notes, a name, and an extension. A pending_file_upload record is created to track the in-flight upload.

#### Scenario: Successful initiation — image type
- **WHEN** a user posts a valid initiate request with an image `mime_type` and optional `artpiece_id`, `notes`, and `name`
- **THEN** the system creates a pending_file_upload record and returns a presigned PUT URL for the original file, a `thumbnail_upload_url` presigned PUT URL for the thumbnail, and an expiry

#### Scenario: Successful initiation — non-image type
- **WHEN** a user posts a valid initiate request with a non-image `mime_type` and optional `artpiece_id`, `notes`, and `name`
- **THEN** the system creates a pending_file_upload record and returns a presigned PUT URL with an expiry; `thumbnail_upload_url` is null

#### Scenario: Successful initiation
- **WHEN** a user posts a valid initiate request with a `mime_type` and optional `artpiece_id`, `notes`, `name`, and `extension`
- **THEN** the system creates a pending_file_upload record (storing `name` and `extension` if provided) and returns a presigned PUT URL with an expiry

#### Scenario: Artpiece not owned by user
- **WHEN** a user initiates an upload with an artpiece_id that does not belong to them
- **THEN** the system returns HTTP 422

### Requirement: Complete file upload
The system SHALL allow an authenticated user to complete a pending file upload. On completion, the pending record is deleted and a file row is created, carrying over `name`, `notes`, `artpiece_id`, `mime_type`, and `extension` from the pending record. The file row is created with `thumbnail_gen_state` set based on mime type. For image-type files, an image_metadata row is inserted and a thumbnail generation job is enqueued. Cover auto-selection runs if the file is attached to an artpiece.

#### Scenario: Successful completion — image type, thumbnail uploaded
- **WHEN** a user completes a pending upload whose mime_type starts with `image/` and `thumbnail_uploaded` is `true`
- **THEN** the file row is created (with `thumbnail_r2_path` set to `users/{user_id}/thumbnails/{file_id}.<ext>`), the pending record is deleted, and an image_metadata row is inserted; no thumbnail job is enqueued

#### Scenario: Successful completion — image type, thumbnail not uploaded
- **WHEN** a user completes a pending upload whose mime_type starts with `image/` and `thumbnail_uploaded` is `false`
- **THEN** the file row is created with `thumbnail_r2_path` null, the pending record is deleted, and an image_metadata row is inserted if dimensions were provided; no thumbnail job is enqueued

#### Scenario: Successful completion — image type
- **WHEN** a user completes a pending upload whose mime_type starts with `image/`
- **THEN** the file row is created (with `name` and `extension` copied from the pending record and `thumbnail_gen_state = "pending"`), the pending record is deleted, an image_metadata row is inserted, and a generate_file_thumbnail job is enqueued

#### Scenario: Successful completion — non-image type
- **WHEN** a user completes a pending upload whose mime_type does not start with `image/`
- **THEN** the file row is created (with `name` and `extension` copied from the pending record and `thumbnail_gen_state = "not_applicable"`) and the pending record is deleted; no thumbnail job is enqueued and no image_metadata row is created

#### Scenario: Pending upload not found or belongs to another user
- **WHEN** a user attempts to complete a pending upload that does not exist or belongs to another user
- **THEN** the system returns HTTP 404
