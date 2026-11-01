# Spec Delta

## ADDED Requirements

### Requirement: filter_payload includes verse_id
The `filter_payload` object stored in a saved filter SHALL support a `verse_id` field (a UUID string, or `null`). When the field is absent on read it SHALL be treated as `null`. The field is accepted on create and update without additional validation beyond type — it is not checked against the user's verses server-side, consistent with how `artist_id` and `character_ids` are already stored opaquely.

#### Scenario: Save filter with verse_id
- **WHEN** a user creates a saved filter with filter_payload containing a non-null `verse_id`
- **THEN** the system stores the payload and returns the record with filter_payload.verse_id matching the provided value

#### Scenario: Save filter with null verse_id
- **WHEN** a user creates a saved filter with filter_payload containing `verse_id: null`
- **THEN** the system stores the payload and returns the record with filter_payload.verse_id: null

#### Scenario: Saved filter without verse_id — treated as null
- **WHEN** a saved filter's filter_payload does not contain a verse_id field
- **THEN** the system treats verse_id as null when reading that filter
