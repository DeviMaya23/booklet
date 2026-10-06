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
- **THEN** the files section header SHALL show the file count and a functional "Download all" button; clicking the button SHALL trigger a zip download of all attached files via the `GET /artpieces/:id/download` endpoint

#### Scenario: Open edit mode from Edit button
- **WHEN** the user clicks the "Edit" button in view mode
- **THEN** the view SHALL switch to edit mode
