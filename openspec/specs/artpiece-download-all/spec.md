## Purpose

Defines the endpoint that streams all files attached to an artpiece as a zip archive, and the frontend flow that triggers a browser download of the resulting file.

## Requirements

### Requirement: Download all artpiece files as zip
The system SHALL provide an authenticated endpoint that streams all files attached to an artpiece as a zip archive without buffering the full archive in memory.

#### Scenario: Successful download — artpiece with files
- **WHEN** an authenticated user sends `GET /artpieces/:id/download` for an artpiece they own that has attached files
- **THEN** the system SHALL respond with HTTP 200, `Content-Type: application/zip`, `Content-Disposition: attachment; filename="<sanitized-title>.zip"`, and a streaming zip body containing one entry per file named `<sanitized-title>-<n>.<ext>` (1-indexed, extension derived from MimeType)

#### Scenario: Successful download — artpiece with no title
- **WHEN** an authenticated user downloads an artpiece whose title is null
- **THEN** the zip filename SHALL be `artpiece.zip` and each entry SHALL be named `artpiece-<n>.<ext>`

#### Scenario: Artpiece not found or not owned
- **WHEN** an authenticated user requests download for an artpiece ID that does not exist or belongs to another user
- **THEN** the system SHALL return HTTP 404

#### Scenario: Artpiece has no files
- **WHEN** an authenticated user downloads an artpiece with no attached files
- **THEN** the system SHALL return HTTP 200 with an empty zip archive

#### Scenario: R2 error mid-stream
- **WHEN** a storage error occurs after headers have been committed and the response body has started streaming
- **THEN** the system SHALL log the error and stop writing; the zip writer SHALL be closed to produce a valid end-of-central-directory record for completed entries; no HTTP error is returned
