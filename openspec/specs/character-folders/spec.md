# Spec: Character Folders

## Purpose

Defines the rules for associating characters with Bookleaf folder IDs, including the data model, lifecycle behaviour (cascading hard-delete), and replace-all update semantics.

---

## Requirements

### Requirement: Character-folder association data model
The `character_folders` table SHALL store associations between characters and Bookleaf folder IDs.
- `character_id` — `uuid`, FK referencing `characters(id)`, not null, part of composite PK
- `folder_id` — `uuid`, not null, no FK constraint (references external Bookleaf system), part of composite PK
- Composite PK on `(character_id, folder_id)` ensures a character cannot be assigned the same folder twice

#### Scenario: Duplicate folder assignment is rejected
- **WHEN** the same `folder_id` is assigned to the same character twice
- **THEN** the database rejects the second insert with a unique constraint violation

#### Scenario: Folder rows are isolated per character
- **WHEN** two characters each have folder assignments
- **THEN** querying folders for character A does not return folders belonging to character B

---

### Requirement: Folder assignments are hard-deleted with their character
When a character is soft-deleted, all its `character_folders` rows SHALL be hard-deleted in the same database transaction.

#### Scenario: Folder rows removed on character delete
- **WHEN** an authenticated user deletes a character that has folder assignments
- **THEN** all `character_folders` rows for that character are removed from the database

---

### Requirement: Folder assignment replace-all semantics on update
When `folder_ids` is provided on a character update, the system SHALL replace all existing folder assignments for that character with the new set in a single atomic operation.

#### Scenario: Existing folders replaced with new set
- **WHEN** a character has folders [A, B] and an update is sent with `folder_ids: [C]`
- **THEN** the character's folders are [C] — A and B are removed

#### Scenario: Empty array clears all folders
- **WHEN** a character has folder assignments and an update is sent with `folder_ids: []`
- **THEN** all folder assignments for that character are removed

#### Scenario: Absent folder_ids field is a no-op
- **WHEN** a character update is sent without a `folder_ids` field
- **THEN** the character's existing folder assignments are unchanged
