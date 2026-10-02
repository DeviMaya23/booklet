## MODIFIED Requirements

### Requirement: Initiate file upload
The system SHALL allow an authenticated user to initiate a file upload by generating a presigned PUT URL for R2, optionally associating the upload with an artpiece, notes, and a name. A pending_file_upload record is created to track the in-flight upload.

#### Scenario: Successful initiation
- **WHEN** a user posts a valid initiate request with a `mime_type` and optional `artpiece_id`, `notes`, and `name`
- **THEN** the system creates a pending_file_upload record (storing `name` if provided) and returns a presigned PUT URL with an expiry

#### Scenario: Artpiece not owned by user
- **WHEN** a user initiates an upload with an `artpiece_id` that does not belong to them
- **THEN** the system returns HTTP 422

### Requirement: Complete file upload
The system SHALL allow an authenticated user to complete a pending file upload. On completion, the pending record is deleted and a file row is created, carrying over `name`, `notes`, `artpiece_id`, and `mime_type` from the pending record. For image-type files, an image_metadata row is inserted and a thumbnail generation job is enqueued. Cover auto-selection runs if the file is attached to an artpiece.

#### Scenario: Successful completion — image type
- **WHEN** a user completes a pending upload whose mime_type starts with `image/`
- **THEN** the file row is created (with `name` copied from the pending record), the pending record is deleted, an image_metadata row is inserted, and a generate_file_thumbnail job is enqueued

#### Scenario: Successful completion — non-image type
- **WHEN** a user completes a pending upload whose mime_type does not start with `image/`
- **THEN** the file row is created (with `name` copied from the pending record) and the pending record is deleted; no thumbnail job is enqueued and no image_metadata row is created

#### Scenario: Pending upload not found or belongs to another user
- **WHEN** a user attempts to complete a pending upload that does not exist or belongs to another user
- **THEN** the system returns HTTP 404
