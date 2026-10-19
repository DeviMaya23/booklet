# Spec Delta

## MODIFIED Requirements

### Requirement: Create saved filter
The system SHALL allow an authenticated user to create a saved filter by providing a name, an optional `thumbnail_artpiece_id`, and a filter payload. The handler resolves the artpiece's cover thumbnail R2 path at creation time and copies it into the record.

#### Scenario: Successful creation — with thumbnail
- **WHEN** an authenticated user POSTs to `/saved_filters` with a valid name, filter_payload, and a `thumbnail_artpiece_id` pointing to an artpiece that has a cover thumbnail
- **THEN** the system SHALL return `201 Created` with the new saved filter record including its generated ID, name, `thumbnail_url` (a presigned URL for the copied thumbnail), filter_payload, created_at, and updated_at

#### Scenario: Successful creation — without thumbnail
- **WHEN** an authenticated user POSTs to `/saved_filters` with a valid name and filter_payload but no `thumbnail_artpiece_id`
- **THEN** the system SHALL return `201 Created` with `thumbnail_url` set to null

#### Scenario: Artpiece has no cover thumbnail
- **WHEN** the `thumbnail_artpiece_id` refers to an artpiece that has no cover thumbnail
- **THEN** the system SHALL return `201 Created` with `thumbnail_url` set to null

#### Scenario: Missing required field
- **WHEN** the request body omits `name` or `filter_payload`
- **THEN** the system SHALL return `422 Unprocessable Entity` with field-level validation errors

#### Scenario: Duplicate name
- **WHEN** the authenticated user already has a saved filter with the same name
- **THEN** the system SHALL return `422 Unprocessable Entity` with a field-level error on `name`

#### Scenario: Unauthenticated request
- **WHEN** the request carries no valid session
- **THEN** the system SHALL return `401 Unauthorized`

---

### Requirement: Update saved filter
The system SHALL allow an authenticated user to partially update a saved filter's name, thumbnail, or filter_payload.

#### Scenario: Update name
- **WHEN** an authenticated user PATCHes `/saved_filters/:id` with a new name
- **THEN** the system SHALL return `200 OK` with the updated record reflecting the new name

#### Scenario: Clear thumbnail
- **WHEN** an authenticated user PATCHes `/saved_filters/:id` with `"thumbnail_r2_path": null`
- **THEN** the system SHALL return `200 OK` with `thumbnail_url` set to null

#### Scenario: Replace filter payload
- **WHEN** an authenticated user PATCHes `/saved_filters/:id` with a new `filter_payload` object
- **THEN** the system SHALL return `200 OK` with the updated record reflecting the new payload

#### Scenario: Empty name rejected
- **WHEN** an authenticated user PATCHes `/saved_filters/:id` with an empty string for name
- **THEN** the system SHALL return `422 Unprocessable Entity`

#### Scenario: Not found
- **WHEN** the ID does not exist or belongs to a different user
- **THEN** the system SHALL return `404 Not Found`

## ADDED Requirements

### Requirement: Saved filter response format
All `/saved_filters` endpoints that return a filter record SHALL include a `thumbnail_url` field (a presigned URL or null) in place of `thumbnail_r2_path`. The raw R2 path is never exposed in the response.

#### Scenario: thumbnail_url in list response
- **WHEN** an authenticated user GETs `/saved_filters`
- **THEN** each record SHALL include a `thumbnail_url` field containing a presigned URL if a thumbnail was stored, or null otherwise

#### Scenario: thumbnail_url in get-by-id response
- **WHEN** an authenticated user GETs `/saved_filters/:id` for a filter they own
- **THEN** the record SHALL include a `thumbnail_url` field (presigned URL or null)
