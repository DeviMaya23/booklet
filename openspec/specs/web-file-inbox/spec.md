## Purpose

Defines the File Inbox UI at `/app/files` — a grid view of all unassigned files with multi-select, search, context menus, per-tile actions, and drag-and-drop upload.

## Requirements

### Requirement: Inbox grid view
The system SHALL display a grid of all unassigned files at `/app/files` (the Inbox), accessible via a sidebar "Inbox" nav item. The grid shows one tile per file and polls for thumbnail readiness.

#### Scenario: Navigating to Inbox
- **WHEN** an authenticated user navigates to `/app/files` or clicks "Inbox" in the sidebar
- **THEN** the system displays the Inbox page with a search bar, the file grid, a floating "+" button, and a persistent "drag files here to upload" hint at the bottom

#### Scenario: File tile with thumbnail pending
- **WHEN** a file in the grid has `thumbnail_url === null` and `thumbnail_gen_state === "pending"`
- **THEN** the tile renders as a shimmer/spinner placeholder and the grid continues polling

#### Scenario: File tile with thumbnail generation settled (no thumbnail)
- **WHEN** a file in the grid has `thumbnail_url === null` and `thumbnail_gen_state !== "pending"`
- **THEN** the tile renders a static fallback icon (no shimmer, no spinner) and polling stops for that file

#### Scenario: Polling stops when all files are settled
- **WHEN** no file in the list has `thumbnail_gen_state === "pending"`
- **THEN** the grid stops polling

---

### Requirement: Notes search
The system SHALL allow the user to filter the inbox grid by `name` or `notes` content. Filtering is performed client-side on the already-fetched file list — no additional server request is made.

#### Scenario: Typing a search query
- **WHEN** a user types in the search bar
- **THEN** the grid is filtered in place to show only files where the `name` or `notes` field contains the search term (case-insensitive); no network request is made

#### Scenario: Clearing the search query
- **WHEN** the user clears the search bar
- **THEN** the full unfiltered list is shown

---

### Requirement: File selection
The system SHALL support standard multi-select behavior on file tiles.

#### Scenario: Single click
- **WHEN** a user single-clicks a tile
- **THEN** that tile becomes the only selected item (previous selection cleared)

#### Scenario: Ctrl/Cmd-click
- **WHEN** a user Ctrl/Cmd-clicks a tile
- **THEN** that tile is toggled in or out of the current selection without affecting other selected tiles

#### Scenario: Shift-click
- **WHEN** a user Shift-clicks a tile
- **THEN** all tiles between the last-clicked tile and the shift-clicked tile are added to the selection

#### Scenario: Right-click on a tile inside the current selection
- **WHEN** a user right-clicks a tile that is already selected
- **THEN** the context menu opens acting on the full current selection

#### Scenario: Click on empty grid background
- **WHEN** a user clicks on the grid area outside any tile
- **THEN** the selection is cleared

#### Scenario: Right-click on a tile outside the current selection
- **WHEN** a user right-clicks a tile that is not currently selected
- **THEN** the selection is replaced with just that tile and the context menu opens acting on that tile alone

---

### Requirement: Selection toolbar
The system SHALL show a toolbar button when one or more files are selected.

#### Scenario: At least one file selected
- **WHEN** the user has one or more files selected
- **THEN** an "Add to artpiece ▾" split-button appears in the top-right of the main area with two dropdown items: "New Artpiece" and "Add to Existing Artpiece"

#### Scenario: No files selected
- **WHEN** no files are selected
- **THEN** the "Add to artpiece ▾" button is not shown

#### Scenario: Clicking "New Artpiece" from selection toolbar
- **WHEN** the user clicks "New Artpiece" in the dropdown
- **THEN** the "Create New Artpiece" modal opens with the currently selected files pre-populated

#### Scenario: Clicking "Add to Existing Artpiece" from selection toolbar
- **WHEN** the user clicks "Add to Existing Artpiece" in the dropdown
- **THEN** the "Add to Existing Artpiece" modal opens with the currently selected files pre-populated

---

### Requirement: Right-click context menu
The system SHALL show a context menu on right-click containing artpiece actions and a destructive delete action.

#### Scenario: Context menu items
- **WHEN** the context menu is open
- **THEN** it shows "New Artpiece", "Add to Existing Artpiece", a visual divider, and "Delete"

#### Scenario: Clicking "New Artpiece" from context menu
- **WHEN** the user clicks "New Artpiece" in the context menu
- **THEN** the "Create New Artpiece" modal opens with the right-clicked selection pre-populated

#### Scenario: Clicking "Add to Existing Artpiece" from context menu
- **WHEN** the user clicks "Add to Existing Artpiece" in the context menu
- **THEN** the "Add to Existing Artpiece" modal opens with the right-clicked selection pre-populated

#### Scenario: Delete from context menu
- **WHEN** the user clicks "Delete" in the context menu
- **THEN** a confirmation dialog appears showing how many files will be deleted; on confirm, the selected files are deleted via bulk delete

---

### Requirement: Per-tile actions menu
The system SHALL show a trash icon on each tile on hover for the Delete action. The "..." dropdown menu and "View Detail" item are removed.

#### Scenario: Hovering a tile
- **WHEN** the user hovers over a file tile
- **THEN** a trash icon appears on the tile

#### Scenario: Clicking the trash icon
- **WHEN** the user clicks the trash icon on a tile
- **THEN** a confirmation dialog appears for that single file; on confirm, the file is deleted

---

### Requirement: File name label beneath tile
The system SHALL display the file's name as a read-only, truncated single-line label beneath each tile in the inbox grid.

#### Scenario: File with a name
- **WHEN** a file in the grid has a non-null `name`
- **THEN** the name is displayed in a single truncated line below the tile image

#### Scenario: File with no name
- **WHEN** a file in the grid has a null `name`
- **THEN** no label is shown beneath the tile

---

### Requirement: Double-click tile opens edit overlay
The system SHALL open the file edit overlay when the user double-clicks a tile. Double-click does not affect the tile's selection state.

#### Scenario: Double-clicking a tile
- **WHEN** a user double-clicks a file tile
- **THEN** the edit overlay opens for that file; the tile's selection state is unchanged

---

### Requirement: Drag-and-drop upload
The system SHALL allow the user to drop multiple files onto the grid to upload them in parallel, with immediate placeholder feedback. Each dropped file's original filename SHALL be sent as the `name` field in the initiate-upload request.

#### Scenario: Dropping files onto the grid
- **WHEN** a user drops one or more files anywhere onto the grid area
- **THEN** each file immediately gets a shimmer placeholder tile inserted into the grid, and each file's upload flow (initiate → PUT to R2 → complete) runs in parallel, with the browser `File.name` sent as `name` in the initiate request

#### Scenario: Upload completes — thumbnail pending
- **WHEN** a file's complete-upload call succeeds and `thumbnail_gen_state === "pending"`
- **THEN** the shimmer placeholder remains and the grid continues polling for that file

#### Scenario: Upload completes — thumbnail settled
- **WHEN** a file's complete-upload call succeeds and `thumbnail_gen_state !== "pending"`
- **THEN** the shimmer placeholder is replaced by the real tile showing either the thumbnail or the fallback icon

#### Scenario: No auto-select on upload completion
- **WHEN** a file finishes uploading
- **THEN** it is NOT automatically added to the selection

---

### Requirement: Floating + button opens upload modal
The system SHALL open the "Add files to dump" modal when the user clicks the floating "+" button. The button is always visible on the Inbox page regardless of selection state.

#### Scenario: Clicking the + button
- **WHEN** a user clicks the floating "+" button on the Inbox page
- **THEN** the "Add files to dump" modal opens
