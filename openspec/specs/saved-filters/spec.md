## Purpose

Defines the `saved_filters` resource — user-owned snapshots of a filter configuration (artist IDs, character IDs, match mode) with an optional thumbnail path and a human-readable name.

## ADDED Requirements

### Requirement: Create saved filter
The system SHALL allow an authenticated user to create a saved filter by providing a name, an optional thumbnail R2 path, and a filter payload.

#### Scenario: Successful creation
- **WHEN** an authenticated user POSTs to `/saved_filters` with a valid name and filter_payload
- **THEN** the system SHALL return `201 Created` with the new saved filter record including its generated ID, name, thumbnail_r2_path (null if omitted), filter_payload, created_at, and updated_at

#### Scenario: Missing required field
- **WHEN** the request body omits `name` or `filter_payload`
- **THEN** the system SHALL return `422 Unprocessable Entity` with field-level validation errors

#### Scenario: Unauthenticated request
- **WHEN** the request carries no valid session
- **THEN** the system SHALL return `401 Unauthorized`

---

### Requirement: List saved filters
The system SHALL allow an authenticated user to retrieve all of their saved filters.

#### Scenario: Successful list
- **WHEN** an authenticated user GETs `/saved_filters`
- **THEN** the system SHALL return `200 OK` with an array of the user's saved filter records ordered by `created_at` descending

#### Scenario: No saved filters
- **WHEN** the user has no saved filters
- **THEN** the system SHALL return `200 OK` with an empty array

#### Scenario: Isolation between users
- **WHEN** user A GETs `/saved_filters`
- **THEN** the response SHALL contain only records belonging to user A, never records belonging to user B

---

### Requirement: Get single saved filter
The system SHALL allow an authenticated user to retrieve a single saved filter by ID.

#### Scenario: Successful get
- **WHEN** an authenticated user GETs `/saved_filters/:id` for a filter they own
- **THEN** the system SHALL return `200 OK` with the saved filter record

#### Scenario: Not found
- **WHEN** the ID does not exist or belongs to a different user
- **THEN** the system SHALL return `404 Not Found`

---

### Requirement: Update saved filter
The system SHALL allow an authenticated user to partially update a saved filter's name, thumbnail_r2_path, or filter_payload.

#### Scenario: Update name
- **WHEN** an authenticated user PATCHes `/saved_filters/:id` with a new name
- **THEN** the system SHALL return `200 OK` with the updated record reflecting the new name

#### Scenario: Clear thumbnail
- **WHEN** an authenticated user PATCHes `/saved_filters/:id` with `"thumbnail_r2_path": null`
- **THEN** the system SHALL return `200 OK` with `thumbnail_r2_path` set to null

#### Scenario: Replace filter payload
- **WHEN** an authenticated user PATCHes `/saved_filters/:id` with a new `filter_payload` object
- **THEN** the system SHALL return `200 OK` with the updated record reflecting the new payload

#### Scenario: Empty name rejected
- **WHEN** an authenticated user PATCHes `/saved_filters/:id` with an empty string for name
- **THEN** the system SHALL return `422 Unprocessable Entity`

#### Scenario: Not found
- **WHEN** the ID does not exist or belongs to a different user
- **THEN** the system SHALL return `404 Not Found`

---

### Requirement: Delete saved filter
The system SHALL allow an authenticated user to permanently delete a saved filter by ID.

#### Scenario: Successful delete
- **WHEN** an authenticated user DELETEs `/saved_filters/:id` for a filter they own
- **THEN** the system SHALL return `204 No Content` and the record SHALL no longer be retrievable

#### Scenario: Not found
- **WHEN** the ID does not exist or belongs to a different user
- **THEN** the system SHALL return `404 Not Found`
