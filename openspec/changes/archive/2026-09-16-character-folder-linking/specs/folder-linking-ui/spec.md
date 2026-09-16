## ADDED Requirements

### Requirement: Folder picker in character modal
The character create and edit modal SHALL include a Folders field that lets the user link public Bookleaf folders to the character. The picker SHALL display currently selected folders as removable chips and allow the user to search and add from the available public folder list.

- Available folders are fetched from `GET /folders` on modal open. This fetch is independent of and does not block modal open.
- The picker SHALL accept no free-text input to create folders — selection is from the available list only.
- A tooltip (question mark icon) beside the Folders label SHALL explain that the folder list comes from Bookleaf.

#### Scenario: Selected folders render as chips
- **WHEN** the character has existing folder assignments
- **THEN** each folder is displayed as a chip with its name and a remove button

#### Scenario: User can remove a folder chip
- **WHEN** the user clicks the remove button on a folder chip
- **THEN** that folder is removed from the selection and will not be submitted

#### Scenario: User can search and add a folder
- **WHEN** the user types in the folder search input
- **THEN** the available folder list is filtered by name (case-insensitive) and the user can select a folder to add it as a chip

#### Scenario: Selected folders are excluded from the add list
- **WHEN** a folder is already selected
- **THEN** it does not appear in the available folder dropdown

#### Scenario: Tooltip explains Bookleaf source
- **WHEN** the user focuses or hovers the question mark icon beside the Folders label
- **THEN** a tooltip appears explaining that the folder list comes from Bookleaf

---

### Requirement: Folder picker degraded state when Bookleaf is unreachable
If the `GET /folders` call fails on modal open, the folder picker SHALL degrade gracefully. Existing folder chips SHALL still render using the persisted names from the character's GET response. The add-new-folder input SHALL be disabled or hidden. No error is shown to the user for this failure.

#### Scenario: Existing folders render when public-folders fetch fails
- **WHEN** `GET /folders` fails on modal open and the character has existing folder assignments
- **THEN** the existing folder chips render correctly using the persisted names, with remove buttons active

#### Scenario: Add input is disabled when public-folders fetch fails
- **WHEN** `GET /folders` fails on modal open
- **THEN** the folder search input is disabled and no dropdown is shown

#### Scenario: Unrelated edits are not blocked
- **WHEN** `GET /folders` fails on modal open
- **THEN** the user can still edit character name, notes, and avatar, and submit the form

---

### Requirement: Folder diff submitted on save
The modal SHALL track folder changes as an explicit diff (adds and removes) against the original folder list from the character's GET response. The final `folder_ids` submitted on save SHALL be: original folder IDs plus added IDs minus removed IDs.

#### Scenario: Diff applied correctly on submit
- **WHEN** the original folders are [A, B], the user removes B and adds C, then saves
- **THEN** the submitted `folder_ids` is [A, C]

#### Scenario: Diff submitted even when public-folders fetch failed
- **WHEN** `GET /folders` failed on modal open, the user removes an existing folder chip, and saves
- **THEN** the submitted `folder_ids` reflects the removal (original minus removed)
