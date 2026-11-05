# Spec Delta

## MODIFIED Requirements

### Requirement: View mode layout
The detail view in view mode SHALL display the artpiece's metadata and file thumbnails in a two-column layout: details on the left, cover preview on the right. The file thumbnail grid SHALL exclude files with `is_visible: false`.

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
- **WHEN** the detail view is open in view mode and the artpiece has attached visible files
- **THEN** the system SHALL display a thumbnail grid of all attached files with `is_visible: true`; the cover tile SHALL show a `★ Cover` chip; each tile SHALL show a type label (e.g. PNG, JPG, PSD) derived from the file's `mime_type`; each tile SHALL be clickable and open the `FileViewer` modal at that file's index

#### Scenario: Display empty file grid
- **WHEN** the detail view is open in view mode and the artpiece has no attached files, or has attached files but all are hidden
- **THEN** the files section SHALL display an empty state message

#### Scenario: Files section header
- **WHEN** the detail view is open in view mode
- **THEN** the files section header SHALL show the visible file count, a functional "Download all" button, and an "Edit details" entry point; clicking "Download all" SHALL trigger a zip download of all attached files via the `GET /artpieces/:id/download` endpoint

#### Scenario: Files section header shows a hidden-file hint
- **WHEN** the artpiece has one or more attached files with `is_visible: false`
- **THEN** the files section header SHALL show a plain-text, non-interactive hint naming the hidden count (e.g. "6 files · 2 hidden")

#### Scenario: No hidden-file hint when none are hidden
- **WHEN** the artpiece has no attached files with `is_visible: false`
- **THEN** the files section header SHALL NOT show a hidden-file hint

#### Scenario: Open edit mode from Edit button
- **WHEN** the user clicks the "Edit" button in view mode
- **THEN** the view SHALL switch to edit mode

---

### Requirement: Edit mode — per-file tile controls
Each file tile in edit mode SHALL show always-visible ☆ (set as cover) and ✕ (remove) controls. The `...` dropdown menu is removed. The edit-mode file grid SHALL exclude files with `is_visible: false`, matching view mode; the files-section header's "Edit details" entry point and hidden-file hint SHALL remain available in edit mode.

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

#### Scenario: Hidden files do not appear as edit-mode tiles
- **WHEN** the detail view is in edit mode and the artpiece has one or more attached files with `is_visible: false`
- **THEN** those files SHALL NOT render as tiles in the edit-mode grid and SHALL NOT be selectable via ☆ or ✕; they remain reachable only through the "Edit details" modal

## ADDED Requirements

### Requirement: Edit details entry point
The files-section header SHALL show an "Edit details" button in both view mode and edit mode. Clicking it opens the file-details modal (see `web-artpiece-file-details-modal`) listing every attached file, including hidden ones, independent of the detail view's current mode or any pending edit-mode staged changes.

#### Scenario: Available in view mode
- **WHEN** the detail view is in view mode
- **THEN** the files-section header SHALL show a functional "Edit details" button

#### Scenario: Available in edit mode
- **WHEN** the detail view is in edit mode
- **THEN** the files-section header SHALL show a functional "Edit details" button, unaffected by pending staged changes
