## MODIFIED Requirements

### Requirement: Commissions list page
The system SHALL provide a page at `/app/commissions` accessible from the sidebar. The page displays two tab buttons — "Active" and "Finished" — each showing a live count of commissions in that state. The active tab is persisted in the URL as `?tab=active` or `?tab=finished`; the default (no param) is Active. A "+ New Commission" button is displayed to the right of the tabs.

#### Scenario: Default view on load
- **WHEN** a user navigates to `/app/commissions` with no `?tab` param
- **THEN** the Active tab SHALL be selected and the table SHALL display commissions with status `waitlist` or `wip`

#### Scenario: Tab counts reflect current data
- **WHEN** a user views the commissions page
- **THEN** the Active tab badge SHALL show the count of commissions with status `waitlist` or `wip`, and the Finished tab badge SHALL show the count of commissions with status `done`

#### Scenario: Switching to Finished tab
- **WHEN** a user clicks the Finished tab
- **THEN** the URL SHALL update to `?tab=finished`, the Finished tab SHALL appear selected, and the table SHALL display only commissions with status `done`

#### Scenario: Switching to Active tab
- **WHEN** a user clicks the Active tab
- **THEN** the URL SHALL update to `?tab=active` (or remove the param), the Active tab SHALL appear selected, and the table SHALL display commissions with status `waitlist` or `wip`

#### Scenario: Tab state restored from URL
- **WHEN** a user navigates to `/app/commissions?tab=finished`
- **THEN** the Finished tab SHALL be pre-selected and the Finished view SHALL be shown

#### Scenario: Empty state
- **WHEN** there are no commissions matching the selected tab
- **THEN** the page SHALL display an empty state message

---

### Requirement: Table columns
The table SHALL display columns with the following typography: header cells are smaller and muted (not bold-foreground); data cells are 14px. Only the currently sorted column's header SHALL show a chevron icon; unsorted column headers SHALL show no sort indicator.

Active tab columns: Title, Artist, Status, Paid, Paid Date, Last Contact, Actions.
Finished tab columns: Title, Artist, Status, Paid, Paid Date, Artpieces, Actions.

#### Scenario: Sort chevron on active column only
- **WHEN** a column is the active sort key
- **THEN** only that column's header SHALL show a chevron (up for asc, down for desc); all other headers SHALL show no chevron

#### Scenario: Artist link shown inline
- **WHEN** a commission has an artist with a non-null artist_link
- **THEN** a small external-link icon is shown inline next to the artist name; clicking it opens the link in a new tab

#### Scenario: Artist link absent
- **WHEN** a commission has no artist or the artist's link is null
- **THEN** the artist link icon is not rendered

#### Scenario: Artpieces column — up to two thumbnails
- **WHEN** a finished commission has one or two attached artpieces
- **THEN** the Artpieces cell SHALL display a thumbnail for each artpiece; each thumbnail SHALL be a link to `/app/artpieces/:id` for that artpiece

#### Scenario: Artpieces column — overflow chip
- **WHEN** a finished commission has more than two attached artpieces
- **THEN** the Artpieces cell SHALL display two thumbnails and a `+N` chip where N is the number of additional artpieces

#### Scenario: Artpieces column — no artpieces
- **WHEN** a finished commission has no attached artpieces
- **THEN** the Artpieces cell SHALL display an em dash (—)

---

### Requirement: Last Contact stamp
In the Active tab, each row SHALL display `last_contacted_at` as a relative time string and a one-click stamp button that sets `last_contacted_at` to now. Entries where `last_contacted_at` is more than 14 days ago (or null) SHALL render in amber to signal they need a chase. The Last Contact column header SHALL show a small info icon; hovering it SHALL display a tooltip explaining the mechanic. The stamp button SHALL show a tooltip "Mark as contacted today" on hover.

#### Scenario: Relative time display — normal
- **WHEN** a commission has a non-null `last_contacted_at` within the last 14 days
- **THEN** the cell displays a relative string (e.g., "today", "3 days ago") in default text colour with a default-coloured stamp icon

#### Scenario: Amber state — stale entry
- **WHEN** a commission has a `last_contacted_at` older than 14 days
- **THEN** the relative time string and the stamp icon SHALL both render in amber

#### Scenario: Null last_contacted_at
- **WHEN** a commission has a null `last_contacted_at`
- **THEN** the cell SHALL display an em dash in muted text (neutral colour, not amber) with a default-coloured stamp icon — null means "never logged," not "overdue"

#### Scenario: Header info tooltip
- **WHEN** a user hovers the info icon in the Last Contact column header
- **THEN** a tooltip SHALL appear with the text: "A manual reminder. Press ↻ whenever you hear from the artist. It turns amber when it's been a while."

#### Scenario: Stamp button tooltip
- **WHEN** a user hovers the stamp button on any row
- **THEN** a tooltip SHALL appear with the text: "Mark as contacted today"

#### Scenario: Stamp to now
- **WHEN** a user clicks the stamp button
- **THEN** the system sends a PATCH with `last_contacted_at: now()`; updates the display on success; shows an inline error on failure

---

## ADDED Requirements

### Requirement: Undo toast on move to Finished
When a commission's status is changed to `done` from the Active tab, the system SHALL immediately perform the PATCH, remove the row from the Active view, and show a Sonner toast with an Undo action. Clicking Undo SHALL send a compensating PATCH to restore the previous status.

#### Scenario: Status changed to done from Active tab
- **WHEN** a user selects "Done" from the inline status dropdown for a commission in the Active tab
- **THEN** the PATCH SHALL be sent immediately, the row SHALL leave the Active list, and a toast SHALL appear reading "_title_ moved to Finished" with an "Undo" action button

#### Scenario: Undo restores previous status
- **WHEN** a user clicks the "Undo" action in the toast
- **THEN** the system SHALL send a PATCH restoring the commission to its previous status (before it was set to done); the commission SHALL reappear in the Active tab

#### Scenario: Toast dismissed without undo
- **WHEN** the toast auto-dismisses or the user dismisses it manually without clicking Undo
- **THEN** no compensating PATCH is sent; the commission remains in Finished
