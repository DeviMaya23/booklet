## ADDED Requirements

### Requirement: Inbox grid view
The system SHALL display a grid of all unassigned files at `/app/files` (the Inbox), accessible via a sidebar "Inbox" nav item. The grid shows one tile per file and polls for thumbnail readiness.

#### Scenario: Navigating to Inbox
- **WHEN** an authenticated user navigates to `/app/files` or clicks "Inbox" in the sidebar
- **THEN** the system displays the Inbox page with a search bar, the file grid, a floating "+" button, and a persistent "drag files here to upload" hint at the bottom

#### Scenario: Files with no thumbnail yet
- **WHEN** a file in the grid has a null `thumbnail_url`
- **THEN** the tile renders as a shimmer placeholder and the grid continues polling until `thumbnail_url` becomes non-null

---

### Requirement: Notes search
The system SHALL allow the user to filter the inbox grid by notes content. Filtering is performed client-side on the already-fetched file list — no additional server request is made.

#### Scenario: Typing a search query
- **WHEN** a user types in the search bar
- **THEN** the grid is filtered in place to show only files whose `notes` field contains the search term (case-insensitive); no network request is made

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

#### Scenario: Clicking "New Artpiece" or "Add to Existing Artpiece"
- **WHEN** the user clicks either dropdown item
- **THEN** no action is taken (no-op placeholder for a future proposal)

---

### Requirement: Right-click context menu
The system SHALL show a context menu on right-click containing artpiece actions (no-op) and a destructive delete action.

#### Scenario: Context menu items
- **WHEN** the context menu is open
- **THEN** it shows "New Artpiece" (no-op), "Add to Existing Artpiece" (no-op), a visual divider, and "Delete"

#### Scenario: Delete from context menu
- **WHEN** the user clicks "Delete" in the context menu
- **THEN** a confirmation dialog appears showing how many files will be deleted; on confirm, the selected files are deleted via bulk delete

---

### Requirement: Per-tile actions menu
The system SHALL show a "..." menu on each tile for single-file actions.

#### Scenario: Per-tile menu items
- **WHEN** the user opens the "..." menu on a tile
- **THEN** it shows "View Detail" (no-op) and "Delete"

#### Scenario: Delete from per-tile menu
- **WHEN** the user clicks "Delete" in the per-tile menu
- **THEN** a confirmation dialog appears for that single file; on confirm, the file is deleted

---

### Requirement: Drag-and-drop upload
The system SHALL allow the user to drop multiple files onto the grid to upload them in parallel, with immediate placeholder feedback.

#### Scenario: Dropping files onto the grid
- **WHEN** a user drops one or more files anywhere onto the grid area
- **THEN** each file immediately gets a shimmer placeholder tile inserted into the grid, and each file's upload flow (initiate → PUT to R2 → complete) runs in parallel

#### Scenario: Upload completes
- **WHEN** a file's complete-upload call succeeds and the subsequent poll detects a non-null `thumbnail_url`
- **THEN** the shimmer placeholder for that file is replaced by the real tile

#### Scenario: No auto-select on upload completion
- **WHEN** a file finishes uploading
- **THEN** it is NOT automatically added to the selection
