# Spec Delta

## MODIFIED Requirements

### Requirement: Init avatar upload
An authenticated user SHALL be able to initiate an avatar upload for a character they own. The system SHALL verify character ownership, create a `pending_character_avatar_uploads` record, generate an R2 object key (`{userID}/files/{uuid}`), and return a presigned PUT URL valid for 15 minutes.

Request body fields:
- `mime_type` — required; must be one of: `image/jpeg`, `image/png`
- `file_size` — required; integer bytes of the file to be uploaded

The system SHALL validate `file_size` against the image category size limit and return HTTP 422 if it is exceeded. The presigned PUT URL SHALL be signed with `ContentLength` equal to `file_size`.

Response body:
- `id` — the pending upload UUID (used as `:uploadID` in the complete endpoint)
- `upload_url` — the presigned R2 PUT URL
- `expires_at` — RFC3339 timestamp when the presigned URL expires

#### Scenario: Successful init
- **WHEN** an authenticated user sends `POST /characters/:id/avatar/init` with a valid `mime_type` and a `file_size` within the image limit for a character they own
- **THEN** the system returns 201 with `id`, `upload_url` (signed with `ContentLength`), and `expires_at`; a `pending_character_avatar_uploads` row exists in the database

#### Scenario: File size exceeds image limit
- **WHEN** an authenticated user sends `POST /characters/:id/avatar/init` with a `file_size` exceeding the image category size limit
- **THEN** the system returns HTTP 422

#### Scenario: Character not found or not owned
- **WHEN** an authenticated user sends `POST /characters/:id/avatar/init` for a character that does not exist or belongs to another user
- **THEN** the system returns 404

#### Scenario: Invalid character UUID path param
- **WHEN** an authenticated user sends `POST /characters/:id/avatar/init` with a non-UUID `:id`
- **THEN** the system returns 400

#### Scenario: Missing mime_type
- **WHEN** an authenticated user sends `POST /characters/:id/avatar/init` without `mime_type`
- **THEN** the system returns 422

#### Scenario: Unsupported mime_type
- **WHEN** an authenticated user sends `POST /characters/:id/avatar/init` with a mime_type not in the allowed set
- **THEN** the system returns 422
