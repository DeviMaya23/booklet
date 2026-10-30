# Spec Delta

## MODIFIED Requirements

### Requirement: Create artpiece
The system SHALL allow an authenticated user to create an artpiece with optional title, artist, notes, a set of characters, an optional set of initial file IDs, and an optional `is_reference` flag (boolean, default false). When no file IDs are provided the artpiece is created with no files and no cover.

#### Scenario: Successful creation with all fields
- **WHEN** a user posts a valid create request with title, artist_id, notes, character_ids, and is_reference true
- **THEN** the system creates the artpiece owned by the user, sets is_reference to true, and returns it with HTTP 201

#### Scenario: Successful creation with no fields
- **WHEN** a user posts a create request with an empty body
- **THEN** the system creates an artpiece with all nullable fields null, is_reference false, and returns it with HTTP 201

#### Scenario: Artist not owned by user
- **WHEN** a user posts a create request with an artist_id that does not belong to them
- **THEN** the system returns HTTP 422 with an error indicating the artist does not belong to the user

#### Scenario: Character not owned by user
- **WHEN** a user posts a create request with a character_id that does not belong to them
- **THEN** the system returns HTTP 422 with an error indicating the character does not belong to the user

---

### Requirement: Update artpiece
The system SHALL allow an authenticated user to update an artpiece's title, artist, notes, character set, and `is_reference` flag. When `is_reference` is omitted from the request the field is left unchanged.

#### Scenario: Successful update
- **WHEN** a user sends a valid update request for an artpiece they own
- **THEN** the system updates the fields and returns the updated artpiece with HTTP 200

#### Scenario: Set is_reference to true
- **WHEN** a user sends an update request with is_reference: true
- **THEN** the system sets is_reference to true and returns the updated artpiece

#### Scenario: Set is_reference to false
- **WHEN** a user sends an update request with is_reference: false on an artpiece that had is_reference true
- **THEN** the system sets is_reference to false and returns the updated artpiece

#### Scenario: Omit is_reference — field unchanged
- **WHEN** a user sends an update request that does not include is_reference
- **THEN** the system does not modify the is_reference field

#### Scenario: Clear artist
- **WHEN** a user sends an update request with artist_id explicitly null
- **THEN** the system sets artist_id to NULL on the artpiece

#### Scenario: Replace character set
- **WHEN** a user sends an update request with a new character_ids array
- **THEN** the system replaces the full character set (not merges) with the provided IDs

#### Scenario: Artpiece not found or belongs to another user
- **WHEN** a user updates an artpiece that does not exist or belongs to another user
- **THEN** the system returns HTTP 404

#### Scenario: Artist not owned by user
- **WHEN** a user updates an artpiece with an artist_id that does not belong to them
- **THEN** the system returns HTTP 422

#### Scenario: Character not owned by user
- **WHEN** a user updates an artpiece with a character_id that does not belong to them
- **THEN** the system returns HTTP 422

---

### Requirement: List artpieces
The system SHALL return all artpieces belonging to the authenticated user, with optional filters for artist, character, and `is_reference`. Each artpiece in the list includes its cover file's thumbnail URL, its commission_id (null if not linked), and — when a cover file exists — its cover file's presigned full-resolution URL, MIME type, and name.

#### Scenario: List with no filters
- **WHEN** a user lists artpieces with no filters
- **THEN** the system returns all artpieces belonging to the user (both references and non-references), ordered by created_at DESC, each including is_reference

#### Scenario: Filter by artist
- **WHEN** a user lists artpieces with one or more artist_ids
- **THEN** the system returns only artpieces whose artist_id is in the filter set

#### Scenario: Filter by character
- **WHEN** a user lists artpieces with one or more character_ids
- **THEN** the system returns only artpieces that have at least one matching character

#### Scenario: Filter is_reference=true
- **WHEN** a user lists artpieces with is_reference=true
- **THEN** the system returns only artpieces where is_reference is true

#### Scenario: Filter is_reference=false
- **WHEN** a user lists artpieces with is_reference=false
- **THEN** the system returns only artpieces where is_reference is false

#### Scenario: Absent is_reference param — no filter applied
- **WHEN** a user lists artpieces without the is_reference param
- **THEN** the system returns artpieces of both kinds

#### Scenario: Artpiece with a cover — cover fields populated
- **WHEN** the list includes an artpiece whose cover_file_id is non-null
- **THEN** the response for that artpiece SHALL include cover_file_url, cover_file_mime_type, and cover_file_name

#### Scenario: Artpiece with no cover — cover fields null
- **WHEN** the list includes an artpiece with cover_file_id NULL
- **THEN** the system returns that artpiece with thumbnail_url, cover_file_url, cover_file_mime_type, and cover_file_name all null

## ADDED Requirements

### Requirement: is_reference in artpiece response
Every artpiece response (create, get, list, update) SHALL include an `is_reference` boolean field.

#### Scenario: Response includes is_reference false
- **WHEN** the system returns an artpiece that has is_reference false
- **THEN** the response JSON includes "is_reference": false

#### Scenario: Response includes is_reference true
- **WHEN** the system returns an artpiece that has is_reference true
- **THEN** the response JSON includes "is_reference": true
