# Spec Delta

## ADDED Requirements

### Requirement: filter_payload must be a JSON object
The system SHALL reject a `filter_payload` value that is not a JSON object (e.g., a string, array, or null literal) on both create and patch. The error SHALL be returned as HTTP 422 with a field-level error on `filter_payload`.

#### Scenario: Non-object payload on create
- **WHEN** an authenticated user POSTs to `/saved_filters` with a `filter_payload` that is not a JSON object (e.g., a string or an array)
- **THEN** the system SHALL return `422 Unprocessable Entity` with a field-level error on `filter_payload`

#### Scenario: Non-object payload on patch
- **WHEN** an authenticated user PATCHes `/saved_filters/:id` and supplies a `filter_payload` that is not a JSON object
- **THEN** the system SHALL return `422 Unprocessable Entity` with a field-level error on `filter_payload`

### Requirement: filter_payload size cap
The system SHALL reject a `filter_payload` exceeding 4 096 bytes on both create and patch. The error SHALL be returned as HTTP 422 with a field-level error on `filter_payload`.

#### Scenario: Oversized payload on create
- **WHEN** an authenticated user POSTs to `/saved_filters` with a `filter_payload` whose serialized size exceeds 4 096 bytes
- **THEN** the system SHALL return `422 Unprocessable Entity` with a field-level error on `filter_payload`

#### Scenario: Oversized payload on patch
- **WHEN** an authenticated user PATCHes `/saved_filters/:id` and supplies a `filter_payload` whose serialized size exceeds 4 096 bytes
- **THEN** the system SHALL return `422 Unprocessable Entity` with a field-level error on `filter_payload`
