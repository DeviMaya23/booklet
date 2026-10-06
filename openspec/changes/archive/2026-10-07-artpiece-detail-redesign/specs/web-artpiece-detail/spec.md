## MODIFIED Requirements

### Requirement: View mode layout
The detail view in view mode SHALL display the artpiece's metadata and file thumbnails in a two-column layout: details on the left, cover preview on the right.

#### Scenario: Two-column layout with cover preview
- **WHEN** the detail view is open in view mode and the artpiece has a cover file
- **THEN** the left column SHALL show a `← Artpieces` back link above the title (with a visible "Edit" button to the right of the title), followed by Artist (with external-link icon if `artist_link` is set), Characters as chips, and Notes; the right column SHALL show the cover thumbnail in a 360px square tile

#### Scenario: Two-column layout — no cover
- **WHEN** the detail view is open in view mode and the artpiece has no cover file
- **THEN** the right column cover tile SHALL render as an empty placeholder

#### Scenario: Display all metadata fields
- **WHEN** the detail view is open in view mode
- **THEN** the system SHALL show the artpiece title (or "Untitled" if null), notes paragraph (or "—" if null), artist name (or "—" if null) with an external-link icon when `artist_link` is set, and characters as individual chips (or "—" if none)

#### Scenario: Display file thumbnail grid
- **WHEN** the detail view is open in view mode and the artpiece has attached files
- **THEN** the system SHALL display a thumbnail grid of all attached files; the cover tile SHALL show a `★ Cover` chip; each tile SHALL show a type label (e.g. PNG, JPG, PSD) derived from the file's `mime_type`

#### Scenario: Display empty file grid
- **WHEN** the detail view is open in view mode and the artpiece has no attached files
- **THEN** the files section SHALL display an empty state message

#### Scenario: Files section header
- **WHEN** the detail view is open in view mode
- **THEN** the files section header SHALL show the file count and a "Download all" button; the button SHALL be present but non-functional (no-op) until download support is implemented

#### Scenario: Open edit mode from Edit button
- **WHEN** the user clicks the "Edit" button in view mode
- **THEN** the view SHALL switch to edit mode

---

### Requirement: Edit mode — header and unsaved changes
The detail view in edit mode SHALL display the title as an inline input, with Cancel and Save controls, and an unsaved-changes caption when there are pending changes.

#### Scenario: Title becomes an editable input
- **WHEN** the detail view enters edit mode
- **THEN** the title SHALL render as a text input pre-populated with the current value; Cancel and Save buttons SHALL appear in place of the Edit button

#### Scenario: Unsaved changes caption appears when changes exist
- **WHEN** the user has made at least one pending change (field edit, file removal, or cover change) in edit mode
- **THEN** a caption with an amber dot SHALL appear below the title input summarising the pending changes (e.g. "Unsaved changes: 1 file will be removed, cover changed")

#### Scenario: Unsaved changes caption absent when no changes
- **WHEN** the user has entered edit mode but made no changes
- **THEN** no unsaved changes caption SHALL be shown

#### Scenario: Save persists field changes
- **WHEN** the user clicks Save
- **THEN** the system SHALL call `PATCH /artpieces/:id` with the updated title, notes, artist, and characters; on success the view SHALL return to view mode showing the updated values and a success toast SHALL appear

#### Scenario: Save calls PUT files only if the file set changed
- **WHEN** the user clicks Save and the set of attached file IDs differs from when edit mode was entered (due to removes)
- **THEN** the system SHALL call `PUT /artpieces/:id/files` with the current working set of file IDs in addition to the field update call

#### Scenario: Save calls PUT cover only if cover changed
- **WHEN** the user clicks Save and the user selected a different cover during the edit session
- **THEN** the system SHALL call `PUT /artpieces/:id/cover` with the new cover file ID

#### Scenario: Cancel reverts field edits and pending removes
- **WHEN** the user clicks Cancel
- **THEN** all field input changes SHALL be discarded, locally-removed files SHALL reappear in the grid, any pending cover change SHALL be discarded, and the view SHALL return to view mode

#### Scenario: Save failure shows error toast
- **WHEN** any Save API call fails
- **THEN** an error toast SHALL appear and the view SHALL remain in edit mode with the user's edits preserved

---

### Requirement: Edit mode — field layout
The detail view in edit mode SHALL display Artist and Characters in a shared row, Notes below, and the cover preview in the right column.

#### Scenario: Fields layout in edit mode
- **WHEN** the detail view is in edit mode
- **THEN** Artist (combobox) and Characters (TokenInput) SHALL occupy the left column in a side-by-side row; Notes (textarea) SHALL appear below them; the right column SHALL show the cover preview tile with a caption "Cover preview. Change it with the ☆ on a file below."

---

### Requirement: Edit mode — file upload
The detail view in edit mode SHALL display an "Add files / or drop them here" drop zone as the last cell in the file grid. Files dropped or picked upload immediately and attach to the artpiece without staging.

#### Scenario: Add-files cell is last in the grid
- **WHEN** the detail view is in edit mode
- **THEN** the last cell in the file grid SHALL be an "Add files / or drop them here" tile with a dashed border; dropping files onto it or clicking it triggers the file picker

#### Scenario: File upload attaches immediately
- **WHEN** the user drops files onto the add-files cell or picks files via the file picker
- **THEN** each file SHALL be uploaded (init → S3 PUT → complete) and then immediately attached to the artpiece via `POST /artpieces/:id/files/:file_id`; the new file thumbnail SHALL appear in the grid on attachment

#### Scenario: Uploaded files persist through Cancel
- **WHEN** the user uploads one or more files in edit mode and then clicks Cancel
- **THEN** those files SHALL remain attached to the artpiece; they SHALL appear in the file grid on the next view-mode render

---

### Requirement: Edit mode — per-file tile controls
Each file tile in edit mode SHALL show always-visible ☆ (set as cover) and ✕ (remove) controls. The `...` dropdown menu is removed.

#### Scenario: Set as cover via ☆
- **WHEN** the user clicks ☆ on a tile that is not the current cover
- **THEN** the selected file SHALL be marked as the pending new cover (☆ on that tile becomes `★ Cover`; the previous cover tile reverts to ☆) but the API call is deferred to Save

#### Scenario: Current cover shows ★ Cover instead of ☆
- **WHEN** a tile is the current cover (or pending cover)
- **THEN** that tile SHALL show `★ Cover` in place of the ☆ control

#### Scenario: Remove via ✕ — tile stays in grid dimmed
- **WHEN** the user clicks ✕ on a tile
- **THEN** the tile SHALL remain in the grid but render dimmed, with a "Will be removed" label and an Undo button; it SHALL NOT be immediately removed from the display

#### Scenario: Undo remove
- **WHEN** the user clicks Undo on a dimmed "Will be removed" tile
- **THEN** the tile SHALL return to its normal appearance and be excluded from the removal set on Save

---

### Requirement: Delete artpiece from detail view
The detail view SHALL provide a "Delete artpiece" text link at the bottom of the edit view that opens the shared `DeleteArtpieceDialog`. The `...` menu entry for Delete in view mode is removed.

#### Scenario: Delete artpiece link in edit mode
- **WHEN** the detail view is in edit mode
- **THEN** a quiet red "Delete artpiece" link SHALL appear at the very bottom of the view

#### Scenario: Delete opens confirmation dialog
- **WHEN** the user clicks "Delete artpiece"
- **THEN** the `DeleteArtpieceDialog` SHALL open, showing the number of currently attached files and a checkbox "Also delete X attached files"

#### Scenario: Confirmed deletion closes detail and refreshes gallery
- **WHEN** the user confirms deletion in the dialog
- **THEN** the artpiece SHALL be deleted, the detail view SHALL close, the gallery grid SHALL be shown, and a success toast SHALL appear

#### Scenario: Cancelled deletion does nothing
- **WHEN** the user dismisses the delete dialog
- **THEN** no API call is made and the detail view remains in edit mode
