# character-recency-sort Specification

## Purpose

Defines the `last_used_at` tracking column on characters and the opt-in recency ordering of the character list, so that character pickers can surface the most recently used characters first without changing the default listing order.

## Requirements

### Requirement: Character last_used_at tracking
The `characters` table SHALL have a `last_used_at timestamptz NULL` column, backfilled from `created_at` for existing rows. The system SHALL set it to the current timestamp when a character is created, and for each assigned character when an artpiece or commission is created or updated with a non-empty `character_ids`. It SHALL NOT change when a character's own fields are edited, and SHALL NOT appear in character responses.

#### Scenario: Backfilled from created_at
- **WHEN** the migration runs against existing characters
- **THEN** each character's `last_used_at` SHALL equal its `created_at`

#### Scenario: Set on character creation
- **WHEN** an authenticated user creates a character
- **THEN** the character's `last_used_at` SHALL be set to approximately the current timestamp

#### Scenario: Bumped on artpiece create with characters
- **WHEN** an authenticated user creates an artpiece with a non-empty `character_ids`
- **THEN** every assigned character's `last_used_at` SHALL be updated to approximately the current timestamp

#### Scenario: Bumped on artpiece update with characters
- **WHEN** an authenticated user updates an artpiece with a non-empty `character_ids`
- **THEN** every character in the request's `character_ids` SHALL have `last_used_at` updated to approximately the current timestamp, including characters already assigned

#### Scenario: Bumped on commission create with characters
- **WHEN** an authenticated user creates a commission with a non-empty `character_ids`
- **THEN** every assigned character's `last_used_at` SHALL be updated to approximately the current timestamp

#### Scenario: Bumped on commission update with characters
- **WHEN** an authenticated user updates a commission with a non-empty `character_ids`
- **THEN** every character in the request's `character_ids` SHALL have `last_used_at` updated to approximately the current timestamp

#### Scenario: Not bumped when character_ids is empty
- **WHEN** an authenticated user creates or updates an artpiece or commission with an empty or absent `character_ids`
- **THEN** no character's `last_used_at` SHALL be modified

#### Scenario: Not bumped on character profile edit
- **WHEN** an authenticated user updates a character's name, notes, avatar, folders, visibility, or verse
- **THEN** that character's `last_used_at` SHALL NOT be changed by the operation

#### Scenario: Not exposed in responses
- **WHEN** any character read endpoint returns a character
- **THEN** the response SHALL NOT contain a `last_used_at` field

---

### Requirement: Character list sorted by recency on request
`GET /characters` SHALL accept an optional `sort` query parameter. With `sort=recent` the results SHALL be ordered by `last_used_at DESC NULLS LAST`, then `name ASC`. When `sort` is absent the order SHALL remain `created_at DESC`. Any other `sort` value SHALL be rejected with 422.

#### Scenario: Most recently used character first
- **WHEN** an authenticated user sends `GET /characters?sort=recent` and character A was most recently assigned
- **THEN** character A SHALL be returned before the others

#### Scenario: Null last_used_at sorts last
- **WHEN** a character has a null `last_used_at` and the user sends `GET /characters?sort=recent`
- **THEN** it SHALL appear after all characters with a non-null `last_used_at`

#### Scenario: Tie broken alphabetically by name
- **WHEN** two characters have the same `last_used_at` and the user sends `GET /characters?sort=recent`
- **THEN** they SHALL be ordered by `name ASC`

#### Scenario: Default order unchanged
- **WHEN** an authenticated user sends `GET /characters` without `sort`
- **THEN** characters SHALL be ordered by `created_at DESC` regardless of `last_used_at`

#### Scenario: Sort combines with existing filters
- **WHEN** an authenticated user sends `GET /characters?sort=recent&verse_id=<uuid>`
- **THEN** only characters with that `verse_id` SHALL be returned, in recency order

#### Scenario: Unknown sort value rejected
- **WHEN** an authenticated user sends `GET /characters?sort=bogus`
- **THEN** the system SHALL return 422 with a structured validation error
