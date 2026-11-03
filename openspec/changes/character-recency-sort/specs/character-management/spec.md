# Spec Delta

## ADDED Requirements

### Requirement: Character last_used_at column
The `characters` table SHALL include a nullable `last_used_at` column of type `timestamptz`, used only for recency ordering (see `character-recency-sort`). It SHALL NOT be accepted in create or update request bodies and SHALL NOT be returned in character responses.

#### Scenario: last_used_at not accepted from clients
- **WHEN** an authenticated user sends `POST /characters` or `PUT /characters/:id` with a `last_used_at` field
- **THEN** the field SHALL be ignored and the stored `last_used_at` SHALL be governed only by the recency rules

---

### Requirement: List characters accepts sort parameter
`GET /characters` SHALL accept the optional `sort` query parameter described in `character-recency-sort`, alongside `q` and `verse_id`. Absent `sort`, the existing ordering applies.

#### Scenario: sort=recent returns the same set in a different order
- **WHEN** an authenticated user sends `GET /characters?sort=recent`
- **THEN** the response SHALL contain the same characters as `GET /characters`, including `folders` on each, ordered by recency

#### Scenario: Invalid sort rejected
- **WHEN** an authenticated user sends `GET /characters?sort=bogus`
- **THEN** the system SHALL return 422
