## MODIFIED Requirements

### Requirement: List images
An authenticated user SHALL be able to retrieve images they own. Each image in the response SHALL include its associated characters (id and name). The endpoint SHALL accept optional filter query parameters; when no filters are provided, all images owned by the user are returned.

Optional query parameters:
- `q` — freetext substring search against `title` (case-insensitive); absent means no title filter
- `character_ids` — repeated UUID values; returns images tagged with any of the given character IDs; absent means no character filter
- `artist_ids` — repeated UUID values; returns images whose `artist_id` matches any of the given IDs; absent means no artist filter

Multiple filters are AND-combined. A malformed UUID in `character_ids` or `artist_ids` returns 400.

#### Scenario: Successful listing with no filters
- **WHEN** an authenticated user sends `GET /images`
- **THEN** the system returns 200 with an array of image objects (may be empty), each including a `characters` array

#### Scenario: Only owner's images returned
- **WHEN** multiple users have images and user A sends `GET /images`
- **THEN** the response contains only user A's images

#### Scenario: Filter by q returns matching images
- **WHEN** an authenticated user sends `GET /images?q=sunset`
- **THEN** the system returns only images whose title contains "sunset" case-insensitively

#### Scenario: Filter by character_ids returns tagged images
- **WHEN** an authenticated user sends `GET /images?character_ids=<uuid>`
- **THEN** the system returns only images tagged with that character, with no duplicate rows

#### Scenario: Filter by artist_ids returns matching images
- **WHEN** an authenticated user sends `GET /images?artist_ids=<uuid>`
- **THEN** the system returns only images whose artist_id matches that UUID

#### Scenario: Malformed UUID in ID filter returns 400
- **WHEN** an authenticated user sends `GET /images?character_ids=not-a-uuid`
- **THEN** the system returns 400
