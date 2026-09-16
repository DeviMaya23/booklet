## MODIFIED Requirements

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

## ADDED Requirements

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
