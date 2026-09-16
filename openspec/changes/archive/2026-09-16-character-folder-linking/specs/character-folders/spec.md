## MODIFIED Requirements

### Requirement: Character-folder association data model
The `character_folders` table SHALL store associations between characters and Bookleaf folders, including the folder's name at the time of last validation.
- `character_id` — `uuid`, FK referencing `characters(id)`, not null, part of composite PK
- `folder_id` — `uuid`, not null, no FK constraint (references external Bookleaf system), part of composite PK
- `folder_name` — `text`, not null — the folder's display name as returned by Bookleaf at the time of last save; written exclusively by the backend during the Bookleaf validation flow; never accepted from FE input
- Composite PK on `(character_id, folder_id)` ensures a character cannot be assigned the same folder twice

#### Scenario: Duplicate folder assignment is rejected
- **WHEN** the same `folder_id` is assigned to the same character twice
- **THEN** the database rejects the second insert with a unique constraint violation

#### Scenario: Folder rows are isolated per character
- **WHEN** two characters each have folder assignments
- **THEN** querying folders for character A does not return folders belonging to character B

#### Scenario: folder_name is persisted alongside folder_id
- **WHEN** a character is saved with valid folder IDs
- **THEN** each `character_folders` row contains the `folder_name` returned by Bookleaf's `GetPublicFolders` response for that ID

#### Scenario: folder_name is refreshed on every successful save
- **WHEN** a character is saved with a folder ID that was previously persisted with an older name
- **THEN** the `folder_name` is overwritten with the current name from Bookleaf's response

## ADDED Requirements

### Requirement: Bookleaf validation on every folder write
When a character is created or updated with a non-empty `folder_ids` list, the system SHALL call `GetPublicFolders` once per request to validate the submitted IDs. This call is authoritative — the result determines which IDs are persisted.

- IDs present in the Bookleaf response → valid; persisted with the `folder_name` from that response (write-through, always refreshed)
- IDs absent from the Bookleaf response → confirmed gone; silently dropped and logged at INFO level; not surfaced as an error to the caller
- If `GetPublicFolders` returns any error → the entire create/update is aborted and the caller receives 500; no partial persistence

If `folder_ids` is empty or absent, `GetPublicFolders` is NOT called.

#### Scenario: Valid folder IDs are persisted with names
- **WHEN** a character is saved with folder IDs that are all present in Bookleaf's public list
- **THEN** all IDs are persisted in `character_folders` with their corresponding `folder_name` values

#### Scenario: Unknown folder IDs are silently dropped
- **WHEN** a character is saved with a `folder_ids` list that includes an ID not present in Bookleaf's public list
- **THEN** that ID is not persisted, the remaining valid IDs are saved, and no error is returned to the caller

#### Scenario: Bookleaf unavailability fails the entire operation
- **WHEN** `GetPublicFolders` returns an error during a character create or update
- **THEN** the system returns 500 and no character data is persisted or changed

#### Scenario: Empty folder_ids skips Bookleaf call
- **WHEN** a character update is sent with `folder_ids: []`
- **THEN** `GetPublicFolders` is not called and all existing folder assignments are removed
