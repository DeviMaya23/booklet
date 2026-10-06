## Purpose

This capability covers the commission form modal — a shared dialog used for creating and editing commissions, including field management, artpiece attachment, and the delete confirmation dialog.

## Requirements

### Requirement: Commission form modal — create mode
The system SHALL provide a modal dialog for creating a new commission, opened by the "+ New Commission" button on the commissions list page. All fields are optional except that the form must be submittable without any required field (Title is optional, unlike Artpiece title). On save the system sends `POST /commissions` with all field values and the current artpiece strip IDs; on success the commissions list is refreshed and the modal closes.

#### Scenario: Open create modal
- **WHEN** a user clicks the "+ New Commission" button
- **THEN** the modal opens in create mode with all fields empty and the artpiece section collapsed

#### Scenario: Submit with no fields
- **WHEN** a user opens the create modal and clicks Save without filling any field
- **THEN** the system posts a commission with status `waitlist`, all other fields null/default, and an empty artpiece_ids array; on success the modal closes and the list refreshes

#### Scenario: Submit with all fields and artpieces
- **WHEN** a user fills all fields and attaches artpieces, then clicks Save
- **THEN** the system posts all field values plus the current artpiece_ids; on success the modal closes and the list refreshes

#### Scenario: Save failure
- **WHEN** the POST /commissions call fails
- **THEN** the modal stays open and displays an error toast; form state is preserved

#### Scenario: Cancel discards state
- **WHEN** a user closes or cancels the modal without saving
- **THEN** no changes are persisted and form state is reset on next open

---

### Requirement: Commission form modal — edit mode
The system SHALL provide the same modal in edit mode, opened by the brush icon on a commission row. The modal pre-populates all fields from the commission's current data. Artpieces are loaded from `GET /commissions/:id` (the list response does not include them). On save the system sends `PUT /commissions/:id` for field changes; if the artpiece strip changed it additionally sends `PUT /commissions/:id/artpieces` (full replace).

#### Scenario: Open edit modal — fields pre-populated
- **WHEN** a user clicks the brush icon on a commission row
- **THEN** the modal opens in edit mode with all editable fields pre-filled from the commission's current data

#### Scenario: Open edit modal — artpiece strip loading
- **WHEN** the edit modal opens
- **THEN** the artpiece strip shows a loading state while `GET /commissions/:id` resolves, then shows the commission's currently attached artpieces

#### Scenario: Save with no changes
- **WHEN** a user opens the edit modal and clicks Save without changing anything
- **THEN** `PUT /commissions/:id` is sent with the unchanged field values; `PUT /commissions/:id/artpieces` is NOT sent; on success the modal closes and the list refreshes

#### Scenario: Save with field changes only
- **WHEN** a user modifies form fields but does not change the artpiece strip, then clicks Save
- **THEN** only `PUT /commissions/:id` is sent; on success the modal closes and the list refreshes

#### Scenario: Save with artpiece set changed
- **WHEN** a user adds or removes artpieces from the strip, then clicks Save
- **THEN** `PUT /commissions/:id` is sent followed by `PUT /commissions/:id/artpieces` with the full current strip IDs; on success the modal closes and the list refreshes

#### Scenario: Save failure
- **WHEN** either update call fails
- **THEN** the modal stays open and displays an error toast; form state is preserved

---

### Requirement: Commission form fields
The form SHALL contain the following fields: Title (text, optional), Artist (combobox with + button to open artist-creation modal), Notes (textarea, optional), Status (chip dropdown, values: waitlist / wip / done, defaults to waitlist in create mode), Price (numeric, optional), Paid (checkbox), Paid Date (datepicker, optional), Finish Date (datepicker, optional, always visible regardless of status).

#### Scenario: Artist creation from modal
- **WHEN** a user clicks the + button next to the Artist combobox
- **THEN** the artist-creation modal opens; on successful creation the new artist is selected in the combobox

#### Scenario: Finish Date always visible
- **WHEN** a user views the form in any status
- **THEN** the Finish Date field is visible regardless of the Status value

#### Scenario: Status defaults to waitlist
- **WHEN** the create modal opens
- **THEN** the Status chip shows `waitlist` as the initial value

---

### Requirement: Artpieces section
The form SHALL include a collapsible "Artpieces" section, collapsed by default. The section contains: a thumbnail strip of attached/selected artpieces (each with an × button to remove), a title search input with a Filter button beside it, and a search results list. The inline filter disclosure is replaced by the shared `ArtistCharacterFilter` panel opened via the Filter button. No All/Any toggle is shown. Collapsing and expanding the section does not reset search, filter, or strip state.

#### Scenario: Filter button opens panel
- **WHEN** the artpieces section is expanded and the user clicks the Filter button
- **THEN** the filter panel SHALL open showing the Artist combobox and Characters token input

#### Scenario: Filter button closes panel
- **WHEN** the filter panel is open and the user clicks outside it or presses Escape
- **THEN** the panel SHALL close without resetting filter values

#### Scenario: Filter button shows active filter summary
- **WHEN** at least one artpiece filter is active
- **THEN** the Filter button SHALL display a summary label instead of "Filter"

#### Scenario: Characters filter — ANY match
- **WHEN** the user selects one or more characters in the filter panel
- **THEN** the artpiece search results SHALL include only artpieces featuring at least one of the selected characters (ANY logic)

#### Scenario: Section collapsed by default
- **WHEN** the modal opens (create or edit)
- **THEN** the artpieces section is collapsed and the strip, search, and filter are not visible

#### Scenario: Expand and collapse preserves state
- **WHEN** a user expands the section, applies a filter, then collapses and re-expands
- **THEN** the filter state is preserved

#### Scenario: Empty state
- **WHEN** no artpieces are attached or selected
- **THEN** the strip area shows "No artpieces attached yet"

#### Scenario: Remove artpiece from strip
- **WHEN** a user clicks the × on an artpiece thumbnail in the strip
- **THEN** the artpiece is removed from the strip immediately (not yet persisted)

#### Scenario: Search results exclude artpieces attached to a different commission
- **WHEN** a user searches for artpieces in the section
- **THEN** artpieces already attached to a different commission do not appear in the results

#### Scenario: Search results include artpieces attached to the current commission
- **WHEN** a user searches in edit mode
- **THEN** artpieces already in this commission's strip do appear in the results (they are already in the strip and selecting them again is a no-op or deduped)

---

### Requirement: StatusChip shared component
The system SHALL provide a shared `StatusChip` controlled component used in both the `CommissionsTable` inline cell and the `CommissionFormModal`. The component accepts `value`, `onChange`, and an optional `hasError` prop. The table cell uses it to fire `patchStatus` on change; the modal form uses it to update local state.

#### Scenario: Table cell autosave unchanged
- **WHEN** a user changes the status chip in the commission table
- **THEN** the behaviour is identical to before — PATCH is fired immediately, inline error shown on failure

#### Scenario: Modal form local state
- **WHEN** a user changes the status chip in the commission form modal
- **THEN** the new value is reflected in local form state only; no network call is made until Save

---

### Requirement: Delete commission dialog
The system SHALL provide a delete confirmation dialog opened by the trash icon on a commission row. The dialog body SHALL state that linked artpieces will be kept intact. On confirm the system sends `DELETE /commissions/:id`; on success the commissions list is refreshed.

#### Scenario: Dialog opened
- **WHEN** a user clicks the trash icon on a commission row
- **THEN** a confirmation dialog opens with text indicating the commission will be permanently deleted and artpieces will be kept intact

#### Scenario: Confirm delete
- **WHEN** a user confirms deletion
- **THEN** the system sends `DELETE /commissions/:id`; on success the dialog closes and the list refreshes

#### Scenario: Cancel delete
- **WHEN** a user cancels the dialog
- **THEN** no request is sent and the commission remains

#### Scenario: Delete failure
- **WHEN** the DELETE call fails
- **THEN** the dialog stays open and an error toast is shown
