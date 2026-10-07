## Purpose

This capability covers downloading a single file by generating a short-lived presigned URL via a dedicated endpoint.

## Requirements

### Requirement: Single file download endpoint
The system SHALL expose a `GET /files/:id/download` endpoint that returns a short-lived presigned download URL for the specified file.

#### Scenario: Successful download URL generation
- **WHEN** an authenticated user requests `GET /files/:id/download` for a file they own
- **THEN** the system SHALL respond HTTP 200 with `{ "download_url": "<presigned URL>" }` where the URL carries `Content-Disposition: attachment` and expires in 5 minutes

#### Scenario: File not found
- **WHEN** an authenticated user requests download for a file ID that does not exist or belongs to another user
- **THEN** the system SHALL respond HTTP 404

#### Scenario: Invalid file ID
- **WHEN** the `:id` path parameter is not a valid UUID
- **THEN** the system SHALL respond HTTP 400
