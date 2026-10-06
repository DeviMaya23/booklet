## Purpose

Defines the artpiece detail view — an inline panel rendered within the Artpieces page when an artpiece is selected, providing view and edit modes for metadata and file management.

## Requirements

### Requirement: Open artpiece detail view
The system SHALL render an inline artpiece detail view within the Artpieces page when an artpiece is selected. No route change occurs.

#### Scenario: Open via double-click
- **WHEN** an authenticated user double-clicks an artpiece card in the gallery
- **THEN** the detail view SHALL replace the gallery grid display, showing the selected artpiece's details in view mode

#### Scenario: Close via back chevron
- **WHEN** the user clicks the back chevron in the detail view
- **THEN** the detail view SHALL close and the gallery grid SHALL be shown again

---

### Requirement: View mode layout
The detail view in view mode SHALL display the artpiece's full metadata and file thumbnails.

#### Scenario: Display all metadata fields
- **WHEN** the detail view is open in view mode
- **THEN** the system SHALL show the artpiece title (or "Untitled" if null), notes (or "—" if null), artist name (or "—" if null), and characters as comma-separated names (or "—" if none)

#### Scenario: Display file thumbnail grid
- **WHEN** the detail view is open in view mode and the artpiece has attached files
- **THEN** the system SHALL display a thumbnail grid of all attached files; the current cover file SHALL show a badge icon on its thumbnail

#### Scenario: Display empty file grid
- **WHEN** the detail view is open in view mode and the artpiece has no attached files
- **THEN** the files section SHALL display an empty state

#### Scenario: Title overflow menu — Edit and Delete
- **WHEN** the user clicks the `...` menu next to the title in view mode
- **THEN** the menu SHALL show two items: "Edit" and "Delete"

#### Scenario: Open edit mode from menu
- **WHEN** the user clicks "Edit" in the `...` menu
- **THEN** the view SHALL switch to edit mode

---

### Requirement: Edit mode — field editing
The detail view in edit mode SHALL allow the user to edit the artpiece's title, notes, artist, and characters using the same input components as the create modal.

#### Scenario: Title becomes an editable input
- **WHEN** the detail view enters edit mode
- **THEN** the title SHALL render as a text input pre-populated with the current value; the `...` menu SHALL be replaced by Save and Cancel (X) controls

#### Scenario: Notes, artist, and characters become editable
- **WHEN** the detail view enters edit mode
- **THEN** notes SHALL render as a textarea, artist as a combobox dropdown, and characters as a chip-based `TokenInput`, each pre-populated with the current values

#### Scenario: Save persists field changes
- **WHEN** the user clicks Save
- **THEN** the system SHALL call `PATCH /artpieces/:id` with the updated title, notes, artist, and characters; on success the view SHALL return to view mode showing the updated values and a success toast SHALL appear

#### Scenario: Save calls PUT files only if the file set changed
- **WHEN** the user clicks Save and the set of attached file IDs differs from when edit mode was entered (due to uploads or removes)
- **THEN** the system SHALL call `PUT /artpieces/:id/files` with the current working set of file IDs in addition to the field update call

#### Scenario: Save calls PUT cover only if cover changed
- **WHEN** the user clicks Save and the user selected a different cover during the edit session
- **THEN** the system SHALL call `PUT /artpieces/:id/cover` with the new cover file ID

#### Scenario: Cancel reverts field edits and pending removes
- **WHEN** the user clicks Cancel (X)
- **THEN** all field input changes SHALL be discarded, locally-removed files SHALL reappear in the grid, any pending cover change SHALL be discarded, and the view SHALL return to view mode

#### Scenario: Save failure shows error toast
- **WHEN** any Save API call fails
- **THEN** an error toast SHALL appear and the view SHALL remain in edit mode with the user's edits preserved

---

### Requirement: Edit mode — file upload
The detail view in edit mode SHALL display a persistent drag-drop upload strip above the file thumbnail grid. Files dropped or picked upload immediately and attach to the artpiece without staging.

#### Scenario: File upload attaches immediately
- **WHEN** the user drops files onto the upload strip or picks files via the file picker
- **THEN** each file SHALL be uploaded (init → S3 PUT → complete) and then immediately attached to the artpiece via `POST /artpieces/:id/files/:file_id`; the new file thumbnail SHALL appear in the grid on attachment

#### Scenario: Uploaded files persist through Cancel
- **WHEN** the user uploads one or more files in edit mode and then clicks Cancel
- **THEN** those files SHALL remain attached to the artpiece; they SHALL appear in the file grid on the next view-mode render

---

### Requirement: Edit mode — per-file thumbnail menu
Each file thumbnail in edit mode SHALL have a `...` menu with "Set as cover" and "Remove from artpiece" actions.

#### Scenario: Set as cover
- **WHEN** the user clicks "Set as cover" on a thumbnail that is not the current cover
- **THEN** the selected file SHALL be marked as the pending new cover (visually reflected by moving the cover badge) but the API call is deferred to Save

#### Scenario: Set as cover disabled for current cover
- **WHEN** a thumbnail is the current cover
- **THEN** the "Set as cover" option in its `...` menu SHALL be disabled or visually greyed out

#### Scenario: Remove from artpiece
- **WHEN** the user clicks "Remove from artpiece" on a thumbnail
- **THEN** the thumbnail SHALL be removed from the displayed grid immediately (local-only); the API call is deferred to Save

---

### Requirement: Delete artpiece from detail view
The detail view SHALL provide a "Delete" entry point in the `...` menu that opens the shared `DeleteArtpieceDialog`.

#### Scenario: Delete opens confirmation dialog
- **WHEN** the user clicks "Delete" in the detail view's `...` menu
- **THEN** the `DeleteArtpieceDialog` SHALL open, showing the number of currently attached files and a checkbox "Also delete X attached files"

#### Scenario: Confirmed deletion closes detail and refreshes gallery
- **WHEN** the user confirms deletion in the dialog
- **THEN** the artpiece SHALL be deleted, the detail view SHALL close, the gallery grid SHALL be shown, and a success toast SHALL appear

#### Scenario: Cancelled deletion does nothing
- **WHEN** the user dismisses the delete dialog
- **THEN** no API call is made and the detail view remains open

---

### Requirement: Delete artpiece dialog
A shared `DeleteArtpieceDialog` component SHALL be used for all artpiece delete confirmations. It includes a checkbox for optionally deleting attached files, which maps to the `?delete_files=true` query parameter.

#### Scenario: Dialog with file count
- **WHEN** the dialog is opened with a known file count greater than zero
- **THEN** the checkbox label SHALL read "Also delete X attached files"

#### Scenario: Dialog without file count
- **WHEN** the dialog is opened without a file count (e.g., from the gallery card where count is unavailable)
- **THEN** the checkbox label SHALL read "Also delete attached files"

#### Scenario: Checkbox unchecked — files detached only
- **WHEN** the user confirms deletion with the checkbox unchecked
- **THEN** `DELETE /artpieces/:id` SHALL be called (no `delete_files` param); files remain as unattached inbox files

#### Scenario: Checkbox checked — files also deleted
- **WHEN** the user confirms deletion with the checkbox checked
- **THEN** `DELETE /artpieces/:id?delete_files=true` SHALL be called; files are permanently deleted along with the artpiece
