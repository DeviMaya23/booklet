# Spec Delta

## MODIFIED Requirements

### Requirement: Initiate file upload
The system SHALL allow an authenticated user to initiate a file upload by generating a presigned PUT URL for R2, optionally associating the upload with an artpiece, notes, and a name. A pending_file_upload record is created to track the in-flight upload. For image-type files, the system SHALL also generate a presigned PUT URL for the thumbnail and return it alongside the upload URL. The request SHALL include a `file_size` field (integer, bytes). The system SHALL validate `file_size` against the per-category size limit for the file's extension and return HTTP 422 if it is exceeded. When `file_size` is within the limit, the presigned PUT URL SHALL be signed with `ContentLength` equal to `file_size` so that R2 enforces the exact upload size.

#### Scenario: Successful initiation — image type
- **WHEN** a user posts a valid initiate request with an image `mime_type`, a `file_size` within the image category limit, and optional `artpiece_id`, `notes`, and `name`
- **THEN** the system creates a pending_file_upload record and returns a presigned PUT URL (signed with `ContentLength`) for the original file, a `thumbnail_upload_url` presigned PUT URL for the thumbnail, and an expiry

#### Scenario: Successful initiation — non-image type
- **WHEN** a user posts a valid initiate request with a non-image `mime_type`, a `file_size` within the limit for that extension's category, and optional `artpiece_id`, `notes`, and `name`
- **THEN** the system creates a pending_file_upload record and returns a presigned PUT URL (signed with `ContentLength`) with an expiry; `thumbnail_upload_url` is null

#### Scenario: File size exceeds category limit
- **WHEN** a user posts an initiate request with a `file_size` that exceeds the limit for the file's extension category
- **THEN** the system returns HTTP 422 with a body indicating the size limit was exceeded

#### Scenario: Artpiece not owned by user
- **WHEN** a user initiates an upload with an artpiece_id that does not belong to them
- **THEN** the system returns HTTP 422
