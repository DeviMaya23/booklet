## Purpose

This capability covers the CRUD lifecycle of artpieces — the core creative works in the system. An artpiece belongs to a user and may reference an artist, a set of characters, attached files, and a cover file.

## Requirements

### Requirement: Create artpiece
The system SHALL allow an authenticated user to create an artpiece with optional title, artist, notes, a set of characters, and an optional set of initial file IDs. When no file IDs are provided the artpiece is created with no files and no cover.

#### Scenario: Successful creation with all fields
- **WHEN** a user posts a valid create request with title, artist_id, notes, and character_ids
- **THEN** the system creates the artpiece owned by the user and returns it with HTTP 201

#### Scenario: Successful creation with no fields
- **WHEN** a user posts a create request with an empty body
- **THEN** the system creates an artpiece with all nullable fields null and returns it with HTTP 201

#### Scenario: Artist not owned by user
- **WHEN** a user posts a create request with an artist_id that does not belong to them
- **THEN** the system returns HTTP 422 with an error indicating the artist does not belong to the user

#### Scenario: Character not owned by user
- **WHEN** a user posts a create request with a character_id that does not belong to them
- **THEN** the system returns HTTP 422 with an error indicating the character does not belong to the user

---

### Requirement: Get artpiece by ID
The system SHALL return a single artpiece by ID, scoped to the authenticated user. The response includes commission_id (null if the artpiece is not linked to a commission).

#### Scenario: Artpiece found
- **WHEN** a user requests an artpiece by ID that belongs to them
- **THEN** the system returns the artpiece with artist, characters, cover file (including thumbnail URL), a `files` array containing each attached file's ID, presigned file URL, and presigned thumbnail URL (null if no thumbnail), and a `commission_id` field (null if not linked), with HTTP 200

#### Scenario: Artpiece not found or belongs to another user
- **WHEN** a user requests an artpiece by ID that does not exist or belongs to another user
- **THEN** the system returns HTTP 404

---

### Requirement: List artpieces
The system SHALL return all artpieces belonging to the authenticated user, with optional filters for artist and character. Each artpiece in the list includes its cover file's thumbnail URL and its commission_id (null if not linked).

#### Scenario: List with no filters
- **WHEN** a user lists artpieces with no filters
- **THEN** the system returns all artpieces belonging to the user, ordered by created_at DESC, each including commission_id

#### Scenario: Filter by artist
- **WHEN** a user lists artpieces with one or more artist_ids
- **THEN** the system returns only artpieces whose artist_id is in the filter set

#### Scenario: Filter by character
- **WHEN** a user lists artpieces with one or more character_ids
- **THEN** the system returns only artpieces that have at least one matching character

#### Scenario: Artpiece with no cover
- **WHEN** the list includes an artpiece with cover_file_id NULL
- **THEN** the system returns that artpiece with thumbnail_url as null (no broken card)

---

### Requirement: Update artpiece
The system SHALL allow an authenticated user to update an artpiece's title, artist, notes, and character set.

#### Scenario: Successful update
- **WHEN** a user sends a valid update request for an artpiece they own
- **THEN** the system updates the fields and returns the updated artpiece with HTTP 200

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

### Requirement: Delete artpiece
The system SHALL allow an authenticated user to delete an artpiece. Deleting an artpiece does NOT delete its files; files remain with artpiece_id set to NULL.

#### Scenario: Successful deletion
- **WHEN** a user deletes an artpiece they own
- **THEN** the system deletes the artpiece and returns HTTP 204, and all previously attached files remain with artpiece_id NULL

#### Scenario: Artpiece not found or belongs to another user
- **WHEN** a user deletes an artpiece that does not exist or belongs to another user
- **THEN** the system returns HTTP 404
