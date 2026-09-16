# Spec: Folder Linking UI

## Purpose

Defines the UI behaviour for linking Bookleaf public folders to a character from the character create and edit modal, including the folder picker component, degraded-state handling when Bookleaf is unreachable, and diff-based submission semantics.

---

## Requirements

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
If the `GET /folders` call fails on modal open, the folder picker SHALL degrade gracefully. Existing folder chips SHALL still render using the persisted names from the character's GET response. The add-new-folder input SHALL be disabled. An inline message SHALL appear below the disabled field informing the user and offering a manual retry. There is no auto-retry — `GET /folders` SHALL be configured with `retry: false`.

#### Scenario: Existing folders render when public-folders fetch fails
- **WHEN** `GET /folders` fails on modal open and the character has existing folder assignments
- **THEN** the existing folder chips render correctly using the persisted names, with remove buttons active

#### Scenario: Add input is disabled when public-folders fetch fails
- **WHEN** `GET /folders` fails on modal open
- **THEN** the folder search input is disabled and no dropdown is shown

#### Scenario: Inline error message shown on failure
- **WHEN** `GET /folders` fails on modal open
- **THEN** an inline message appears below the disabled field reading "Couldn't reach folder list" with a Retry action; the message uses a muted/neutral color (not a destructive error color)

#### Scenario: Retry triggers a new fetch
- **WHEN** the user clicks the Retry action
- **THEN** a new `GET /folders` request is issued; on success the picker re-enables, the inline message disappears, and reconciliation runs against the freshly loaded list

#### Scenario: Unrelated edits are not blocked
- **WHEN** `GET /folders` fails on modal open
- **THEN** the user can still edit character name, notes, and avatar, and submit the form

---

### Requirement: Folder name and ID reconciliation on successful fetch
When `GET /folders` succeeds (including after a manual retry), the modal SHALL silently reconcile the character's stored folders against the live Bookleaf list. No transition indicator or user notice is shown for either reconciliation action.

- A stored folder whose name differs from Bookleaf's current name SHALL have its chip name silently updated to the Bookleaf name.
- A stored folder whose ID is absent from Bookleaf's list SHALL be silently removed from the rendered selection, treated identically to a manual chip removal. This removal is not reversible from this UI. The BE will handle DB cleanup on the next PUT (per replace-all semantics).
- A folder added during the current modal session (not from stored data) already carries the current Bookleaf name; it is not subject to name reconciliation.
- A stale cache hit counts as a successful response for reconciliation purposes.

#### Scenario: Stale folder name is silently updated
- **WHEN** `GET /folders` succeeds and a stored folder's name differs from the name in the Bookleaf response
- **THEN** the chip displays the current Bookleaf name; no indicator is shown to the user

#### Scenario: Folder missing from Bookleaf is silently removed
- **WHEN** `GET /folders` succeeds and a stored folder's ID is not present in the Bookleaf response
- **THEN** that folder's chip is removed from the selection without user interaction; the removal is reflected in the diff submitted on save

#### Scenario: Reconciliation runs after manual retry succeeds
- **WHEN** `GET /folders` initially failed, the user clicked Retry, and the retry succeeded
- **THEN** reconciliation runs against the newly loaded Bookleaf list

#### Scenario: No reconciliation on fetch failure
- **WHEN** `GET /folders` fails (including after retry)
- **THEN** stored folder chips remain rendered as-is with no name or ID changes applied

---

### Requirement: Folder diff submitted on save
The modal SHALL track folder changes as an explicit diff (adds and removes) against the original folder list from the character's GET response. The final `folder_ids` submitted on save SHALL be: original folder IDs plus added IDs minus removed IDs.

#### Scenario: Diff applied correctly on submit
- **WHEN** the original folders are [A, B], the user removes B and adds C, then saves
- **THEN** the submitted `folder_ids` is [A, C]

#### Scenario: Diff submitted even when public-folders fetch failed
- **WHEN** `GET /folders` failed on modal open, the user removes an existing folder chip, and saves
- **THEN** the submitted `folder_ids` reflects the removal (original minus removed)
