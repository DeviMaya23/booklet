## MODIFIED Requirements

### Requirement: Partial update image
An authenticated user SHALL be able to fully replace the editable fields of an image they own using PUT semantics. The request body SHALL always include all updatable fields; the server SHALL write them exactly as received.

Updatable fields: `title`, `artist_id`, `notes`, `character_ids`.

- `title` — nullable string; `null` clears the value
- `notes` — nullable string; `null` clears the value
- `artist_id` — nullable UUID string; `null` clears the artist association; if non-null, the referenced artist SHALL exist and belong to the requesting user, otherwise the system SHALL return 422
- `character_ids` — required array of UUID strings; an empty array `[]` clears all character associations; each ID SHALL belong to the requesting user, otherwise the system SHALL return 422

The `thumbnail_r2_path` field is NOT accepted in the update request body; it is managed exclusively by the thumbnail worker.

#### Scenario: Successful full update
- **WHEN** an authenticated user sends `PUT /images/:id` with all updatable fields
- **THEN** the system returns 200 with the updated image reflecting the submitted values

#### Scenario: Clearing title
- **WHEN** an authenticated user sends `PUT /images/:id` with `"title": null`
- **THEN** the system returns 200 and the image's `title` is null

#### Scenario: Clearing notes
- **WHEN** an authenticated user sends `PUT /images/:id` with `"notes": null`
- **THEN** the system returns 200 and the image's `notes` is null

#### Scenario: Set artist by ID
- **WHEN** an authenticated user sends `PUT /images/:id` with a valid `artist_id` they own
- **THEN** the system returns 200 and the image's `artist` reflects the referenced artist

#### Scenario: Clear artist association
- **WHEN** an authenticated user sends `PUT /images/:id` with `"artist_id": null`
- **THEN** the system returns 200 and the image has no artist association

#### Scenario: Artist not owned by user
- **WHEN** an authenticated user sends `PUT /images/:id` with an `artist_id` that belongs to another user or does not exist
- **THEN** the system returns 422

#### Scenario: Replace character associations
- **WHEN** an authenticated user sends `PUT /images/:id` with `"character_ids": ["uuid-a", "uuid-b"]`
- **THEN** the system returns 200 and the image's characters are exactly `[uuid-a, uuid-b]`

#### Scenario: Clear all character associations
- **WHEN** an authenticated user sends `PUT /images/:id` with `"character_ids": []`
- **THEN** the system returns 200 and the image has no associated characters

#### Scenario: Character not owned by user
- **WHEN** an authenticated user sends `PUT /images/:id` with a `character_ids` list containing an ID that belongs to another user or does not exist
- **THEN** the system returns 422

#### Scenario: Image not found or not owned
- **WHEN** an authenticated user sends `PUT /images/:id` for an image that does not exist or belongs to another user
- **THEN** the system returns 404

#### Scenario: Malformed request body
- **WHEN** an authenticated user sends `PUT /images/:id` with invalid JSON
- **THEN** the system returns 400
