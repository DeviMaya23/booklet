## MODIFIED Requirements

### Requirement: CompleteUpload
An authenticated user SHALL be able to complete an in-progress upload by calling `POST /images/:id/complete`, where `:id` is the `pending_upload` UUID returned by InitialUpload. The system SHALL, in a single database transaction: validate which of the stored character IDs belong to the user (silently dropping any that do not), validate the stored `artist_id` (silently setting it to null if the artist no longer exists or does not belong to the user), delete the `pending_upload` row, and insert a new row into `images` with `thumbnail_r2_path` set to null. After the transaction commits, the system SHALL enqueue a `generate_thumbnail` River job for the created image. The response SHALL be 201 with no body.

#### Scenario: Successful CompleteUpload
- **WHEN** an authenticated user sends `POST /images/:id/complete` for a pending upload they own
- **THEN** the system returns 201 with no body; the `pending_upload` row is deleted; a new `images` row exists with the same `r2_key` as `image_r2_path` and `mime_type` from the pending upload; `thumbnail_r2_path` is null; a `generate_thumbnail` job is in the River queue

#### Scenario: Artist ID silently nulled if not found
- **WHEN** the pending_upload has an `artist_id` that no longer exists or belongs to another user
- **THEN** the created image has `artist_id` set to null; no error is returned

#### Scenario: Artist ID carried through if valid
- **WHEN** the pending_upload has an `artist_id` that exists and belongs to the user
- **THEN** the created image has `artist_id` set to that value

#### Scenario: Character IDs silently filtered
- **WHEN** the pending_upload contains character IDs where some belong to the user and some do not
- **THEN** the created image is associated only with the valid character IDs; no error is returned

#### Scenario: All character IDs invalid
- **WHEN** the pending_upload contains character IDs none of which belong to the user
- **THEN** the created image is associated with no characters; no error is returned

#### Scenario: pending_upload not found or not owned
- **WHEN** an authenticated user sends `POST /images/:id/complete` for an ID that does not exist or belongs to another user
- **THEN** the system returns 404

#### Scenario: Invalid UUID path param
- **WHEN** an authenticated user sends `POST /images/:id/complete` with a non-UUID value for `:id`
- **THEN** the system returns 400
