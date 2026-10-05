## MODIFIED Requirements

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
