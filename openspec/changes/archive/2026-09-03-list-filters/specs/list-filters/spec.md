## ADDED Requirements

### Requirement: Freetext search via `q` param
The `GET /images`, `GET /artists`, and `GET /characters` list endpoints SHALL accept an optional `q` query parameter. When present, the system SHALL return only records whose searchable field contains the value as a case-insensitive substring.

- For images: `q` is matched against `title`
- For artists: `q` is matched against `name`
- For characters: `q` is matched against `name`

When `q` is absent, all records for the user are returned (existing behavior preserved).

#### Scenario: Freetext match returns matching records
- **WHEN** an authenticated user sends a list request with `?q=foo`
- **THEN** the system returns only records whose searchable field contains "foo" case-insensitively

#### Scenario: Freetext match is case-insensitive
- **WHEN** an authenticated user sends a list request with `?q=FOO`
- **THEN** the system returns records whose searchable field contains "foo" regardless of case

#### Scenario: Freetext match is a substring match
- **WHEN** a record's searchable field is "Foobar" and the user sends `?q=oob`
- **THEN** that record is included in the response

#### Scenario: No `q` param returns all records
- **WHEN** an authenticated user sends a list request without `q`
- **THEN** all records owned by the user are returned, unchanged from current behavior

#### Scenario: `q` with no matches returns empty array
- **WHEN** an authenticated user sends a list request with `?q=zzznomatch`
- **THEN** the system returns 200 with an empty array

---

### Requirement: Multi-value ID filters on image list
The `GET /images` endpoint SHALL accept optional `character_ids` and `artist_ids` query parameters. Each parameter is multi-value (repeated keys) and accepts UUID values. When present, the system SHALL return only images associated with at least one of the given IDs.

- `character_ids`: filters images that are tagged with any of the specified character IDs
- `artist_ids`: filters images whose `artist_id` matches any of the specified artist IDs
- Multiple filters are AND-combined: if both `character_ids` and `artist_ids` are provided, an image must satisfy both
- A malformed UUID value in either parameter SHALL result in a 400 response

#### Scenario: Filter by single character_id
- **WHEN** an authenticated user sends `GET /images?character_ids=<uuid>`
- **THEN** the system returns only images tagged with that character

#### Scenario: Filter by multiple character_ids returns images matching any
- **WHEN** an authenticated user sends `GET /images?character_ids=<uuid-a>&character_ids=<uuid-b>`
- **THEN** the system returns images tagged with uuid-a OR uuid-b (no duplicates)

#### Scenario: Filter by single artist_id
- **WHEN** an authenticated user sends `GET /images?artist_ids=<uuid>`
- **THEN** the system returns only images whose artist_id matches that UUID

#### Scenario: Filter by multiple artist_ids returns images matching any
- **WHEN** an authenticated user sends `GET /images?artist_ids=<uuid-a>&artist_ids=<uuid-b>`
- **THEN** the system returns images whose artist_id is uuid-a OR uuid-b

#### Scenario: Combined character_ids and artist_ids are AND-combined
- **WHEN** an authenticated user sends `GET /images?character_ids=<uuid-c>&artist_ids=<uuid-a>`
- **THEN** the system returns only images that are both tagged with uuid-c AND have artist_id uuid-a

#### Scenario: Malformed UUID in ID filter returns 400
- **WHEN** an authenticated user sends `GET /images?character_ids=not-a-uuid`
- **THEN** the system returns 400
