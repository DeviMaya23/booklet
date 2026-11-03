# Spec Delta

## MODIFIED Requirements

### Requirement: Create artpiece
The system SHALL allow an authenticated user to create an artpiece with optional title, artist, notes, a set of characters, a set of tags, an optional set of initial file IDs, and an optional `is_reference` flag (boolean, default false). When no file IDs are provided the artpiece is created with no files and no cover.

#### Scenario: Successful creation with all fields
- **WHEN** a user posts a valid create request with title, artist_id, notes, character_ids, tag_ids, and is_reference true
- **THEN** the system creates the artpiece owned by the user, sets is_reference to true, assigns the characters and tags, and returns it with HTTP 201

#### Scenario: Successful creation with no fields
- **WHEN** a user posts a create request with an empty body
- **THEN** the system creates an artpiece with all nullable fields null, is_reference false, no tags, and returns it with HTTP 201

#### Scenario: Artist not owned by user
- **WHEN** a user posts a create request with an artist_id that does not belong to them
- **THEN** the system returns HTTP 422 with an error indicating the artist does not belong to the user

#### Scenario: Character not owned by user
- **WHEN** a user posts a create request with a character_id that does not belong to them
- **THEN** the system returns HTTP 422 with an error indicating the character does not belong to the user

#### Scenario: Created artpiece returns its tags
- **WHEN** a user posts a create request with tag_ids naming two tags they own
- **THEN** the response `tags` array contains both tags with their id and name

#### Scenario: Tag not owned by user
- **WHEN** a user posts a create request with a tag_id that does not exist or belongs to another user
- **THEN** the system returns HTTP 422 with an error indicating the tag does not belong to the user, and no artpiece is created

---

### Requirement: Get artpiece by ID
The system SHALL return a single artpiece by ID, scoped to the authenticated user. The response includes commission_id (null if the artpiece is not linked to a commission) and the artpiece's tags.

#### Scenario: Artpiece found
- **WHEN** a user requests an artpiece by ID that belongs to them
- **THEN** the system returns the artpiece with artist, characters, tags, cover file (including thumbnail URL), a `files` array containing each attached file's ID, presigned file URL, and presigned thumbnail URL (null if no thumbnail), and a `commission_id` field (null if not linked), with HTTP 200

#### Scenario: Artpiece with tags
- **WHEN** a user requests an artpiece that has tags `wip` and `sketch`
- **THEN** the response `tags` array contains both, each with its id and name

#### Scenario: Artpiece without tags
- **WHEN** a user requests an artpiece that has no tags
- **THEN** the response `tags` is an empty array

#### Scenario: Artpiece not found or belongs to another user
- **WHEN** a user requests an artpiece by ID that does not exist or belongs to another user
- **THEN** the system returns HTTP 404

---

### Requirement: List artpieces
The system SHALL return all artpieces belonging to the authenticated user, with optional filters for artist, character, tag, and `is_reference`. Each artpiece in the list includes its tags, its cover file's thumbnail URL, its commission_id (null if not linked), and — when a cover file exists — its cover file's presigned full-resolution URL, MIME type, and name.

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
- **THEN** the response for that artpiece SHALL include `cover_file_url` (a presigned GET URL for the full-resolution cover file), `cover_file_mime_type`, and `cover_file_name` (null if the file has no name)

#### Scenario: Artpiece with no cover — cover fields null
- **WHEN** the list includes an artpiece with cover_file_id NULL
- **THEN** the system returns that artpiece with `thumbnail_url`, `cover_file_url`, `cover_file_mime_type`, and `cover_file_name` all null

#### Scenario: List returns each artpiece's tags
- **WHEN** a user lists artpieces and one of them has tags `wip` and `sketch`
- **THEN** that artpiece's `tags` array in the list response contains both, and an artpiece without tags has an empty array

#### Scenario: Filter by a single tag
- **WHEN** a user lists artpieces with `tag_ids=<A>`
- **THEN** the system returns only artpieces that have tag A, each still showing all of its tags

#### Scenario: Multiple tags default to matching any
- **WHEN** a user lists artpieces with `tag_ids=<A>&tag_ids=<B>` and no `tag_match`
- **THEN** the system returns artpieces that have tag A or tag B, with no duplicates

#### Scenario: tag_match=any matches any
- **WHEN** a user lists artpieces with `tag_ids=<A>&tag_ids=<B>&tag_match=any`
- **THEN** the system returns artpieces that have tag A or tag B, with no duplicates

#### Scenario: tag_match=all requires every tag
- **WHEN** a user lists artpieces with `tag_ids=<A>&tag_ids=<B>&tag_match=all`
- **THEN** the system returns only artpieces that have both tag A and tag B

#### Scenario: Tag filter combines with other filters
- **WHEN** a user lists artpieces with `tag_ids=<A>&artist_ids=<X>&character_ids=<C>`
- **THEN** the system returns only artpieces that satisfy all of the filters

#### Scenario: Tag filter naming another user's tag
- **WHEN** a user lists artpieces with a `tag_ids` value belonging to another user's tag
- **THEN** no artpiece of the other user is returned and the response is HTTP 200

#### Scenario: Invalid tag filter values rejected
- **WHEN** a user lists artpieces with a non-UUID `tag_ids` value or a `tag_match` value other than `any` or `all`
- **THEN** the system returns HTTP 400

---

### Requirement: Update artpiece
The system SHALL allow an authenticated user to update an artpiece's title, artist, notes, character set, tag set, and `is_reference` flag. When `is_reference` is omitted from the request the field is left unchanged.

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

#### Scenario: Replace tag set
- **WHEN** a user sends an update request with a new tag_ids array on an artpiece that already has other tags
- **THEN** the system replaces the full tag set (not merges) with the provided IDs and the response `tags` reflects the new set

#### Scenario: Empty tag_ids clears tags
- **WHEN** a user sends an update request with an empty tag_ids array
- **THEN** the artpiece has no tags afterwards and the tags themselves still exist

#### Scenario: Artpiece not found or belongs to another user
- **WHEN** a user updates an artpiece that does not exist or belongs to another user
- **THEN** the system returns HTTP 404

#### Scenario: Artist not owned by user
- **WHEN** a user updates an artpiece with an artist_id that does not belong to them
- **THEN** the system returns HTTP 422

#### Scenario: Character not owned by user
- **WHEN** a user updates an artpiece with a character_id that does not belong to them
- **THEN** the system returns HTTP 422

#### Scenario: Tag not owned by user
- **WHEN** a user updates an artpiece with a tag_id that does not exist or belongs to another user
- **THEN** the system returns HTTP 422 and the artpiece's existing tags are unchanged
