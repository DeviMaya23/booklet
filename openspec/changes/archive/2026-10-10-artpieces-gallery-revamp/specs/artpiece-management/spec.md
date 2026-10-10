# Spec Delta

## MODIFIED Requirements

### Requirement: List artpieces
The system SHALL return all artpieces belonging to the authenticated user, with optional filters for artist and character. Each artpiece in the list includes its cover file's thumbnail URL, its commission_id (null if not linked), and — when a cover file exists — its cover file's presigned full-resolution URL, MIME type, and name.

#### Scenario: List with no filters
- **WHEN** a user lists artpieces with no filters
- **THEN** the system returns all artpieces belonging to the user, ordered by created_at DESC, each including commission_id

#### Scenario: Filter by artist
- **WHEN** a user lists artpieces with one or more artist_ids
- **THEN** the system returns only artpieces whose artist_id is in the filter set

#### Scenario: Filter by character
- **WHEN** a user lists artpieces with one or more character_ids
- **THEN** the system returns only artpieces that have at least one matching character

#### Scenario: Artpiece with a cover — cover fields populated
- **WHEN** the list includes an artpiece whose cover_file_id is non-null
- **THEN** the response for that artpiece SHALL include `cover_file_url` (a presigned GET URL for the full-resolution cover file), `cover_file_mime_type`, and `cover_file_name` (null if the file has no name)

#### Scenario: Artpiece with no cover — cover fields null
- **WHEN** the list includes an artpiece with cover_file_id NULL
- **THEN** the system returns that artpiece with `thumbnail_url`, `cover_file_url`, `cover_file_mime_type`, and `cover_file_name` all null
