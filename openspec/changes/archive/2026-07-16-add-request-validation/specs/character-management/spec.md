## MODIFIED Requirements

### Requirement: Create character
An authenticated user SHALL be able to create a character with a required name and optional biography, hero image path, and public visibility flag. The system SHALL assign a UUID and associate the character with the authenticated user.

#### Scenario: Successful creation
- **WHEN** an authenticated user sends `POST /characters` with a valid name
- **THEN** the system returns 201 with the created character including its generated UUID and all provided fields

#### Scenario: Missing required name
- **WHEN** an authenticated user sends `POST /characters` with no name field (or empty string)
- **THEN** the system returns 422 with a structured validation error

#### Scenario: Malformed request body
- **WHEN** an authenticated user sends `POST /characters` with invalid JSON
- **THEN** the system returns 400

---

## ADDED Requirements

### Requirement: Update character name must not be empty if provided
If the `name` field is present in a `PATCH /characters/:id` request body, it SHALL contain at least one character. An explicitly empty string SHALL be rejected.

#### Scenario: Empty name on update returns 422
- **WHEN** an authenticated user sends `PATCH /characters/:id` with `"name": ""`
- **THEN** the system returns 422 with a structured validation error for the `name` field

#### Scenario: Absent name on update is accepted
- **WHEN** an authenticated user sends `PATCH /characters/:id` without a `name` field
- **THEN** the system returns 200 and the character's existing name is unchanged
