# Spec Delta

## MODIFIED Requirements

### Requirement: Update saved filter
The system SHALL allow an authenticated user to partially update a saved filter's name, thumbnail, or filter_payload. The thumbnail MAY be updated by providing either `thumbnail_artpiece_id` (the handler resolves the artpiece's cover thumbnail R2 path) or `thumbnail_r2_path` (raw path, for internal use). When both are provided, `thumbnail_artpiece_id` takes precedence.

#### Scenario: Update name
- **WHEN** an authenticated user PATCHes `/saved_filters/:id` with a new name
- **THEN** the system SHALL return `200 OK` with the updated record reflecting the new name

#### Scenario: Clear thumbnail
- **WHEN** an authenticated user PATCHes `/saved_filters/:id` with `"thumbnail_r2_path": null`
- **THEN** the system SHALL return `200 OK` with `thumbnail_url` set to null

#### Scenario: Replace thumbnail via artpiece ID
- **WHEN** an authenticated user PATCHes `/saved_filters/:id` with a valid `thumbnail_artpiece_id`
- **THEN** the system SHALL resolve the artpiece's cover thumbnail and return `200 OK` with `thumbnail_url` reflecting the new thumbnail

#### Scenario: Artpiece has no cover thumbnail
- **WHEN** the `thumbnail_artpiece_id` refers to an artpiece with no cover thumbnail
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
