## MODIFIED Requirements

### Requirement: Update artist
An authenticated user SHALL be able to fully replace the editable fields of an artist they own using PUT semantics. The request body SHALL always include all updatable fields; the server SHALL write them exactly as received.

Updatable fields: `name`, `notes`, `artist_link`.

- `name` — required, non-empty string
- `notes` — nullable string; `null` clears the value
- `artist_link` — nullable string; `null` clears the value; if non-null, must be a valid URL

#### Scenario: Successful update
- **WHEN** an authenticated user sends `PUT /artists/:id` with all updatable fields
- **THEN** the system returns 200 with the updated artist reflecting the submitted values

#### Scenario: Clearing notes
- **WHEN** an authenticated user sends `PUT /artists/:id` with `"notes": null`
- **THEN** the system returns 200 and the artist's `notes` is null

#### Scenario: Clearing artist_link
- **WHEN** an authenticated user sends `PUT /artists/:id` with `"artist_link": null`
- **THEN** the system returns 200 and the artist's `artist_link` is null

#### Scenario: Duplicate name on update
- **WHEN** an authenticated user sends `PUT /artists/:id` with a name already used by another artist they own
- **THEN** the system returns 409

#### Scenario: Artist not found or not owned
- **WHEN** an authenticated user sends `PUT /artists/:id` for an artist that does not exist or belongs to another user
- **THEN** the system returns 404

#### Scenario: Invalid artist_link format
- **WHEN** an authenticated user sends `PUT /artists/:id` with `artist_link` set to a non-URL string
- **THEN** the system returns 422

#### Scenario: Malformed request body
- **WHEN** an authenticated user sends `PUT /artists/:id` with invalid JSON
- **THEN** the system returns 400
