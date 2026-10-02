## MODIFIED Requirements

### Requirement: Notes search
The system SHALL allow the user to filter the inbox grid by `name` or `notes` content. Filtering is performed client-side on the already-fetched file list — no additional server request is made.

#### Scenario: Typing a search query
- **WHEN** a user types in the search bar
- **THEN** the grid is filtered in place to show only files where the `name` or `notes` field contains the search term (case-insensitive); no network request is made

#### Scenario: Clearing the search query
- **WHEN** the user clears the search bar
- **THEN** the full unfiltered list is shown

### Requirement: Drag-and-drop upload
The system SHALL allow the user to drop multiple files onto the grid to upload them in parallel, with immediate placeholder feedback. Each dropped file's original filename SHALL be sent as the `name` field in the initiate-upload request.

#### Scenario: Dropping files onto the grid
- **WHEN** a user drops one or more files anywhere onto the grid area
- **THEN** each file immediately gets a shimmer placeholder tile inserted into the grid, and each file's upload flow (initiate → PUT to R2 → complete) runs in parallel, with the browser `File.name` sent as `name` in the initiate request

#### Scenario: Upload completes
- **WHEN** a file's complete-upload call succeeds and the subsequent poll detects a non-null `thumbnail_url`
- **THEN** the shimmer placeholder for that file is replaced by the real tile

#### Scenario: No auto-select on upload completion
- **WHEN** a file finishes uploading
- **THEN** it is NOT automatically added to the selection
