## ADDED Requirements

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

## MODIFIED Requirements

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
