## Purpose

This capability covers the CRUD lifecycle of commissions — records that group artpieces under a single commissioned work. A commission belongs to a user and may reference an artist, a set of characters, and an attached set of artpieces.

## Requirements

### Requirement: Create commission
The system SHALL allow an authenticated user to create a commission with optional title, artist_id, status (defaulting to `waitlist`), price, paid flag, paid_date, finish_date, notes, character_ids, and an optional array of artpiece_ids to attach in the same call.

#### Scenario: Successful creation with all fields
- **WHEN** a user posts a valid create request with artist_id, status, price, paid, paid_date, finish_date, notes, character_ids, and artpiece_ids
- **THEN** the system creates the commission owned by the user and returns it with HTTP 201

#### Scenario: Successful creation with no fields
- **WHEN** a user posts a create request with an empty body
- **THEN** the system creates a commission with status `waitlist`, paid false, and all other nullable fields null, returning HTTP 201

#### Scenario: Artist not owned by user
- **WHEN** a user creates a commission with an artist_id that does not belong to them
- **THEN** the system returns HTTP 422

#### Scenario: Character not owned by user
- **WHEN** a user creates a commission with a character_id that does not belong to them
- **THEN** the system returns HTTP 422

#### Scenario: Artpiece not owned by user
- **WHEN** a user creates a commission with an artpiece_id that does not belong to them
- **THEN** the system returns HTTP 422

#### Scenario: Artpiece already attached to a different commission
- **WHEN** a user creates a commission with an artpiece_id that already has a commission_id pointing to a different commission
- **THEN** the system returns HTTP 422 with an error indicating the artpiece is already attached to another commission

#### Scenario: Invalid status value
- **WHEN** a user creates a commission with a status value not in (`waitlist`, `wip`, `done`)
- **THEN** the system returns HTTP 422

---

### Requirement: Get commission by ID
The system SHALL return a single commission by ID, scoped to the authenticated user. The response includes all commission fields (including title), the character set, and a summary of attached artpieces (id + cover thumbnail URL).

#### Scenario: Commission found
- **WHEN** a user requests a commission by ID that belongs to them
- **THEN** the system returns the commission with its artist_id, status, price, paid, paid_date, finish_date, notes, characters, and artpieces (each with id and cover thumbnail URL), with HTTP 200

#### Scenario: Commission not found or belongs to another user
- **WHEN** a user requests a commission by ID that does not exist or belongs to another user
- **THEN** the system returns HTTP 404

---

### Requirement: List commissions
The system SHALL return all commissions belonging to the authenticated user, ordered by created_at DESC.

#### Scenario: List all commissions
- **WHEN** a user lists commissions with no filters
- **THEN** the system returns all commissions belonging to the user ordered by created_at DESC

---

### Requirement: Update commission
The system SHALL allow an authenticated user to update a commission's title, artist_id, status, price, paid, paid_date, finish_date, notes, and character set.

#### Scenario: Successful update
- **WHEN** a user sends a valid update request for a commission they own
- **THEN** the system updates the fields and returns the updated commission with HTTP 200

#### Scenario: Replace character set
- **WHEN** a user sends an update request with a new character_ids array
- **THEN** the system replaces the full character set (not merges) with the provided IDs

#### Scenario: Commission not found or belongs to another user
- **WHEN** a user updates a commission that does not exist or belongs to another user
- **THEN** the system returns HTTP 404

#### Scenario: Artist not owned by user
- **WHEN** a user updates a commission with an artist_id that does not belong to them
- **THEN** the system returns HTTP 422

#### Scenario: Character not owned by user
- **WHEN** a user updates a commission with a character_id that does not belong to them
- **THEN** the system returns HTTP 422

#### Scenario: Invalid status value
- **WHEN** a user updates a commission with a status value not in (`waitlist`, `wip`, `done`)
- **THEN** the system returns HTTP 422

---

### Requirement: Delete commission
The system SHALL allow an authenticated user to delete a commission. Deleting a commission does NOT delete its artpieces; artpieces remain with commission_id set to NULL.

#### Scenario: Successful deletion
- **WHEN** a user deletes a commission they own
- **THEN** the system deletes the commission and returns HTTP 204, and all previously attached artpieces remain with commission_id NULL

#### Scenario: Commission not found or belongs to another user
- **WHEN** a user deletes a commission that does not exist or belongs to another user
- **THEN** the system returns HTTP 404
