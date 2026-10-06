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
The system SHALL return a single commission by ID, scoped to the authenticated user. The response includes all commission fields (including title and last_contacted_at), the character set, a summary of attached artpieces (id + cover thumbnail URL), and the artist's link when an artist is attached.

#### Scenario: Commission found
- **WHEN** a user requests a commission by ID that belongs to them
- **THEN** the system returns the commission with its artist_id, artist_name, artist_link, status, price, paid, paid_date, finish_date, last_contacted_at, notes, characters, and artpieces (each with id and cover thumbnail URL), with HTTP 200

#### Scenario: Commission not found or belongs to another user
- **WHEN** a user requests a commission by ID that does not exist or belongs to another user
- **THEN** the system returns HTTP 404

---

### Requirement: List commissions
The system SHALL return all commissions belonging to the authenticated user, ordered by created_at DESC. Each item in the list SHALL include artist_link and last_contacted_at.

#### Scenario: List all commissions
- **WHEN** a user lists commissions with no filters
- **THEN** the system returns all commissions belonging to the user ordered by created_at DESC, each including artist_link and last_contacted_at

---

### Requirement: Patch commission
The system SHALL allow an authenticated user to partially update a commission's inline-editable fields: `status`, `paid`, `paid_date`, and `last_contacted_at`. Only fields present in the request body are updated; absent fields are left unchanged. All fields are optional and use pointer types to distinguish explicit `null`/`false` from absent.

#### Scenario: Patch status only
- **WHEN** a user sends `PATCH /commissions/:id` with only `{ "status": "wip" }` for a commission they own
- **THEN** the system updates the status and returns the updated commission with HTTP 200, leaving all other fields unchanged

#### Scenario: Patch paid to true
- **WHEN** a user sends `{ "paid": true }` in a PATCH request
- **THEN** the system sets `paid = true` and returns the updated commission with HTTP 200

#### Scenario: Patch paid to false
- **WHEN** a user sends `{ "paid": false }` in a PATCH request
- **THEN** the system sets `paid = false` and returns the updated commission with HTTP 200 (false is a valid value, not treated as absent)

#### Scenario: Patch last_contacted_at
- **WHEN** a user sends `{ "last_contacted_at": "<timestamp>" }` in a PATCH request
- **THEN** the system sets `last_contacted_at` to the provided timestamp and returns the updated commission with HTTP 200

#### Scenario: Patch paid_date to null
- **WHEN** a user sends `{ "paid_date": null }` in a PATCH request
- **THEN** the system clears `paid_date` and returns the updated commission with HTTP 200

#### Scenario: Invalid status value
- **WHEN** a user sends a PATCH with `status` not in (`waitlist`, `wip`, `done`)
- **THEN** the system returns HTTP 422

#### Scenario: Commission not found or belongs to another user
- **WHEN** a user patches a commission that does not exist or belongs to another user
- **THEN** the system returns HTTP 404

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
