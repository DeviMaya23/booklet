# Verse Management Specification

## Purpose

Defines the data model, CRUD endpoints, ownership rules, and recency tracking for verses — narrative groupings (universes, settings, canons) that a user can assign to characters.

---

## Requirements

### Requirement: Verse data model
The `verses` table SHALL use the following schema:
- `id` — `uuid`, primary key, application-generated
- `user_id` — `uuid`, FK referencing `users(id)`, not null
- `name` — `text`, not null
- `sort_order` — `integer`, not null, default `0` — explicit display order set by the user
- `last_used_at` — `timestamptz`, nullable — updated when a character is assigned to this verse
- `created_at`, `updated_at` — `timestamptz`, not null
- Unique constraint on `(user_id, name)`

#### Scenario: Schema enforces uniqueness per user
- **WHEN** two verses with the same name are created for the same user
- **THEN** the second insert is rejected with a unique-constraint violation

#### Scenario: Names are unique per user, not globally
- **WHEN** two different users each create a verse with the same name
- **THEN** both verses are persisted without error

---

### Requirement: Create verse
An authenticated user SHALL be able to create a verse with a required name. The system SHALL assign a UUID, associate the verse with the authenticated user, and set `sort_order` to one greater than the current maximum `sort_order` for that user (appending to the end). If no verses exist yet, `sort_order` is set to `0`.

#### Scenario: Verse created successfully
- **WHEN** an authenticated user submits a valid name
- **THEN** the verse is created with `sort_order = MAX(sort_order) + 1` for the user and the created verse is returned

#### Scenario: Name conflict returns 409
- **WHEN** an authenticated user submits a name that already exists for their account
- **THEN** the request is rejected with HTTP 409

#### Scenario: Name is required
- **WHEN** an authenticated user submits an empty or missing name
- **THEN** the request is rejected with HTTP 422

---

### Requirement: List verses ordered by sort_order
An authenticated user SHALL be able to retrieve all their verses. The response SHALL be ordered by `sort_order ASC`.

#### Scenario: Verses returned in sort_order
- **WHEN** a user has multiple verses with different `sort_order` values
- **THEN** the list response contains all their verses ordered from lowest to highest `sort_order`

#### Scenario: Empty list for user with no verses
- **WHEN** a user has no verses
- **THEN** the list response is an empty array

---

### Requirement: Get verse by ID
An authenticated user SHALL be able to retrieve a single verse by ID, provided they own it.

#### Scenario: Verse found and returned
- **WHEN** a user requests a verse they own by its ID
- **THEN** the verse is returned

#### Scenario: Verse not found or not owned returns 404
- **WHEN** a user requests a verse ID that does not exist or belongs to another user
- **THEN** the request is rejected with HTTP 404

---

### Requirement: Update verse name
An authenticated user SHALL be able to rename a verse they own. The name SHALL remain unique per user; a conflicting name returns 409.

#### Scenario: Verse renamed successfully
- **WHEN** a user submits a new unique name for a verse they own
- **THEN** the verse name is updated and the updated verse is returned

#### Scenario: Rename to conflicting name returns 409
- **WHEN** a user renames a verse to a name already used by another of their verses
- **THEN** the request is rejected with HTTP 409

#### Scenario: Rename not-owned verse returns 404
- **WHEN** a user attempts to rename a verse that does not exist or belongs to another user
- **THEN** the request is rejected with HTTP 404

---

### Requirement: Delete verse
An authenticated user SHALL be able to delete a verse they own. Characters assigned to the deleted verse SHALL have their `verse_id` set to `NULL` automatically (enforced by the `ON DELETE SET NULL` FK constraint).

#### Scenario: Verse deleted, characters unassigned
- **WHEN** a user deletes a verse that has characters assigned to it
- **THEN** the verse is removed and all previously assigned characters have `verse_id = NULL`

#### Scenario: Verse deleted with no characters
- **WHEN** a user deletes a verse that has no characters assigned
- **THEN** the verse is removed and the response is 204 No Content

#### Scenario: Delete not-owned verse returns 404
- **WHEN** a user attempts to delete a verse that does not exist or belongs to another user
- **THEN** the request is rejected with HTTP 404

---

### Requirement: Reorder verses
An authenticated user SHALL be able to set an explicit display order for their verses via `POST /verses/reorder`. The request body SHALL contain an ordered array of verse IDs. The system SHALL assign `sort_order = 0, 1, 2, …` to each verse in that order, in a single atomic transaction. IDs in the request that do not belong to the user SHALL be rejected, leaving no verses modified.

#### Scenario: Verses reordered successfully
- **WHEN** a user submits all their verse IDs in a new order
- **THEN** each verse's `sort_order` is updated to its index position in the submitted array, atomically, and the response is 204 No Content

#### Scenario: Partial list rejected
- **WHEN** a user submits an ID array that omits one or more of their verses
- **THEN** the request is rejected with HTTP 422 and no sort_order values are changed

#### Scenario: Foreign ID in list rejected
- **WHEN** the submitted array contains an ID that belongs to another user or does not exist
- **THEN** the request is rejected with HTTP 422 and no sort_order values are changed

---

### Requirement: last_used_at updated on verse assignment
When a character is created or updated with a non-null `verse_id`, the system SHALL update `last_used_at` on that verse to the current timestamp.

#### Scenario: last_used_at bumped on character assignment
- **WHEN** a character is saved with a `verse_id`
- **THEN** the referenced verse's `last_used_at` is updated to now

#### Scenario: last_used_at not updated when verse_id is null
- **WHEN** a character is saved with `verse_id = null`
- **THEN** no verse's `last_used_at` is modified

---

### Requirement: Verse ownership isolation
The system SHALL ensure a user cannot read or modify verses belonging to another user. All queries SHALL be scoped to the authenticated user's ID.

#### Scenario: Cross-user isolation on list
- **WHEN** user A lists verses
- **THEN** only user A's verses are returned, never user B's

---

### Requirement: Verse response shape
Every verse response SHALL include:
- `id` (string, UUID)
- `name` (string)
- `last_used_at` (string, RFC3339, or null)
- `created_at` (string, RFC3339)
- `updated_at` (string, RFC3339)

#### Scenario: Verse response fields present
- **WHEN** any verse endpoint returns a verse object
- **THEN** all five fields are present with the correct types
