## ADDED Requirements

### Requirement: File response includes thumbnail_gen_state
The system SHALL include a `thumbnail_gen_state` field in every file response — GET, list, update, and complete-upload responses. For rows with a null `thumbnail_gen_state` (pre-migration), the response SHALL return `"failed"`.

#### Scenario: File with thumbnail_gen_state set
- **WHEN** a file has a non-null `thumbnail_gen_state`
- **THEN** the response includes `"thumbnail_gen_state": "<value>"`

#### Scenario: File with null thumbnail_gen_state
- **WHEN** a file has a null `thumbnail_gen_state`
- **THEN** the response includes `"thumbnail_gen_state": "failed"`
