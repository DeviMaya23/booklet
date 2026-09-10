## MODIFIED Requirements

### Requirement: Partial update character
An authenticated user SHALL be able to fully replace the editable fields of a character they own using PUT semantics. The request body SHALL always include all updatable fields; the server SHALL write them exactly as received.

The `avatar_r2_path` field is NOT accepted in the update request body. Avatar images are set exclusively via the avatar upload flow.

Updatable fields: `name`, `biography`, `is_public`, `folder_ids`.

- `name` — required, non-empty string
- `biography` — nullable string; `null` clears the value
- `is_public` — required boolean
- `folder_ids` — required array of UUIDs; an empty array `[]` removes all folder assignments; each value SHALL be a valid UUID (format validation only)

#### Scenario: Successful full update
- **WHEN** an authenticated user sends `PUT /characters/:id` with all updatable fields
- **THEN** the system returns 200 with the updated character reflecting the submitted values

#### Scenario: Clearing biography
- **WHEN** an authenticated user sends `PUT /characters/:id` with `"biography": null`
- **THEN** the system returns 200 and the character's `biography` is null

#### Scenario: Flipping is_public to false
- **WHEN** an authenticated user sends `PUT /characters/:id` with `"is_public": false`
- **THEN** the system returns 200 and `is_public` is set to false on the character

#### Scenario: Updating folder_ids replaces assignments
- **WHEN** an authenticated user sends `PUT /characters/:id` with a new `folder_ids` array
- **THEN** the system returns 200 and the character's `folder_ids` reflect the new set

#### Scenario: Clearing all folder assignments
- **WHEN** an authenticated user sends `PUT /characters/:id` with `"folder_ids": []`
- **THEN** the system returns 200 and the character has no folder assignments

#### Scenario: Invalid UUID in folder_ids returns 422
- **WHEN** an authenticated user sends `PUT /characters/:id` with a `folder_ids` array containing a non-UUID value
- **THEN** the system returns 422 with a structured validation error

#### Scenario: Character not found or not owned
- **WHEN** an authenticated user sends `PUT /characters/:id` for a character that does not exist or belongs to another user
- **THEN** the system returns 404

#### Scenario: Malformed request body
- **WHEN** an authenticated user sends `PUT /characters/:id` with invalid JSON
- **THEN** the system returns 400

---

### Requirement: Update character name must not be empty if provided
If the `name` field is present in a `PUT /characters/:id` request body, it SHALL contain at least one character. An explicitly empty string SHALL be rejected.

`name` is always required in a PUT body; it cannot be omitted.

#### Scenario: Empty name on update returns 422
- **WHEN** an authenticated user sends `PUT /characters/:id` with `"name": ""`
- **THEN** the system returns 422 with a structured validation error for the `name` field
