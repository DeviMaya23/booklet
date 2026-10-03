## MODIFIED Requirements

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
