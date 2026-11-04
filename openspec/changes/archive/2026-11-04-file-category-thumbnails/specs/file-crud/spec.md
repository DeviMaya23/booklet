# Spec Delta

## ADDED Requirements

### Requirement: File response includes extension
The system SHALL include an `extension` field (string or null) in every file response — GET, list, update, and complete-upload responses.

#### Scenario: File with a known extension
- **WHEN** a file has a non-null `extension` (e.g., `"clip"`)
- **THEN** every file response for that file includes `"extension": "clip"`

#### Scenario: File without a stored extension
- **WHEN** a file has a null `extension` (uploaded before this change or with no extension provided)
- **THEN** every file response for that file includes `"extension": null`
