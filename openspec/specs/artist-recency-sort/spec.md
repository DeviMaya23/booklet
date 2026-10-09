# Artist Recency Sort

## Purpose

Defines the `last_used_at` tracking column on artists and the recency-based ordering applied to the artist list endpoint and lookup combobox, so that the most recently used artists surface first.

---

## Requirements

### Requirement: Artist last_used_at tracking
The `artists` table SHALL have a `last_used_at timestamptz NULL` column. The system SHALL set this column to the current timestamp in the following situations:
- When an artist is created (set at creation time, so new artists appear at the top of the lookup immediately)
- When an artist is assigned on a commission create or update (artist_id is non-null in the request)
- When an artist is assigned on an artpiece create or update (artist_id is non-null in the request)

`last_used_at` SHALL NOT be updated when an artist's own fields (name, notes, links) are edited without an assignment event.

#### Scenario: Set on artist creation
- **WHEN** an authenticated user creates a new artist
- **THEN** the artist's `last_used_at` SHALL be set to approximately the current timestamp

#### Scenario: Updated on commission create with artist
- **WHEN** an authenticated user creates a commission with a non-null `artist_id`
- **THEN** the assigned artist's `last_used_at` SHALL be updated to approximately the current timestamp

#### Scenario: Updated on commission update with artist
- **WHEN** an authenticated user updates a commission and the request includes a non-null `artist_id`
- **THEN** the assigned artist's `last_used_at` SHALL be updated to approximately the current timestamp

#### Scenario: Updated on artpiece create with artist
- **WHEN** an authenticated user creates an artpiece with a non-null `artist_id`
- **THEN** the assigned artist's `last_used_at` SHALL be updated to approximately the current timestamp

#### Scenario: Updated on artpiece update with artist
- **WHEN** an authenticated user updates an artpiece and the request includes a non-null `artist_id`
- **THEN** the assigned artist's `last_used_at` SHALL be updated to approximately the current timestamp

#### Scenario: Not updated on artist profile edit
- **WHEN** an authenticated user updates an artist's name, notes, or links
- **THEN** `last_used_at` SHALL NOT be changed by that operation

---

### Requirement: Artist lookup sorted by recency
The artist list endpoint (`GET /artists`) SHALL order results by `last_used_at DESC NULLS LAST`, with `name ASC` as a secondary sort key to break ties.

The artist lookup combobox in the web application SHALL consume this ordering and display artists in the order returned by the API, with no client-side re-sort applied.

#### Scenario: Most recently used artist appears first
- **WHEN** an authenticated user has multiple artists and artist A was most recently assigned
- **THEN** `GET /artists` SHALL return artist A before the others

#### Scenario: Null last_used_at sorts last
- **WHEN** an artist has a null `last_used_at`
- **THEN** it SHALL appear after all artists with a non-null `last_used_at`

#### Scenario: Tie broken alphabetically by name
- **WHEN** two artists have the same `last_used_at` value
- **THEN** they SHALL be ordered by `name ASC`
