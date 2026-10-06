## Purpose

This capability covers the commissions list page — the primary UI for scanning and managing commission status, payment state, and artist contact history.

## Requirements

### Requirement: Commissions list page
The system SHALL provide a page at `/app/commissions` accessible from the sidebar. The page displays a table of commissions with a toggle to switch between two mutually exclusive views: "Waitlist / In Progress" (default) and "Finished."

#### Scenario: Default view on load
- **WHEN** a user navigates to `/app/commissions`
- **THEN** the page displays commissions with status `waitlist` or `wip`, in the "Waitlist / In Progress" toggle state

#### Scenario: Toggle to Finished
- **WHEN** a user activates the toggle
- **THEN** the page switches to show only commissions with status `done`, replacing the Last Contact column with a Time Taken column

#### Scenario: Empty state
- **WHEN** there are no commissions matching the active toggle state
- **THEN** the page displays an empty state message

---

### Requirement: Table columns
The table SHALL display the following columns in both toggle states: Title, Artist (with inline link icon), Status, Paid, Paid Date. An additional column differs by state: Last Contact in Waitlist/In Progress; Time Taken in Finished.

All columns are sortable except the artist link control (sorting raw URLs is not meaningful).

#### Scenario: Artist link shown inline
- **WHEN** a commission has an artist with a non-null artist_link
- **THEN** a small external-link icon is shown inline next to the artist name; clicking it opens the link in a new tab

#### Scenario: Artist link absent
- **WHEN** a commission has no artist or the artist's link is null
- **THEN** the artist link icon is not rendered

#### Scenario: Time Taken with both dates present
- **WHEN** a commission in the Finished view has both finish_date and paid_date set
- **THEN** Time Taken displays the difference as "N days"

#### Scenario: Time Taken with missing date
- **WHEN** a commission in the Finished view has a null finish_date or null paid_date
- **THEN** Time Taken displays an em dash (—)

---

### Requirement: Inline autosave
The Status, Paid, and Paid Date cells SHALL be inline-editable and autosave on change with no per-row Save button. A failed save SHALL show a visible inline error state on the affected cell; success is silent.

#### Scenario: Status dropdown change
- **WHEN** a user selects a new status from the inline dropdown
- **THEN** the system sends a PATCH for that commission and updates the cell on success; shows an inline error on failure

#### Scenario: Paid checkbox toggle to true
- **WHEN** a user checks the Paid checkbox
- **THEN** the system sends a PATCH with `paid: true` and `last_contacted_at: now()` in a single call; updates the cell on success; shows an inline error on failure

#### Scenario: Paid checkbox toggle to false
- **WHEN** a user unchecks the Paid checkbox
- **THEN** the system sends a PATCH with `paid: false` (last_contacted_at is not modified); updates the cell on success; shows an inline error on failure

#### Scenario: Paid Date change
- **WHEN** a user picks a date in the Paid Date datepicker
- **THEN** the system sends a PATCH with the selected date; updates the cell on success; shows an inline error on failure

---

### Requirement: Last Contact stamp
In the Waitlist / In Progress view, each row SHALL display `last_contacted_at` as a relative time string and a one-click stamp button that sets `last_contacted_at` to now.

#### Scenario: Relative time display
- **WHEN** a commission has a non-null last_contacted_at
- **THEN** the cell displays a relative string (e.g., "10 days ago", "today", "yesterday")

#### Scenario: Null last_contacted_at
- **WHEN** a commission has a null last_contacted_at
- **THEN** the cell displays an em dash (—) and the stamp button

#### Scenario: Stamp to now
- **WHEN** a user clicks the stamp button
- **THEN** the system sends a PATCH with `last_contacted_at: now()`; updates the display on success; shows an inline error on failure

---

### Requirement: New Commission button
The page SHALL display a "+ New Commission" button in the header. In this version the button is present but non-functional (no-op). The creation flow is out of scope.

#### Scenario: Button present
- **WHEN** a user views the commissions page
- **THEN** a "+ New Commission" button is visible in the page header
