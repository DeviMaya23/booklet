## Purpose

This capability covers the `FileViewer` fullscreen modal component for previewing and downloading individual files attached to an artpiece, with navigation between files.

## Requirements

### Requirement: FileViewer component — image preview state
The system SHALL provide a `FileViewer` component that renders as a fullscreen modal overlay and displays a single file at a time with navigation controls.

#### Scenario: Open viewer on file tile click
- **WHEN** a user clicks a file tile in the artpiece detail view (view mode)
- **THEN** `FileViewer` SHALL open as a fullscreen overlay, displaying the clicked file at the initial index

#### Scenario: Image preview with checkerboard background
- **WHEN** the viewer displays a file whose `mimeType` starts with `image/` and whose `thumbnailUrl` is non-null
- **THEN** the image SHALL be rendered in a centered tile with a CSS checkerboard background so transparent areas are visible against the dark scrim

#### Scenario: Top bar — image state
- **WHEN** the viewer displays a previewable file
- **THEN** the top bar SHALL show: filename (or "Untitled" if null), type badge (e.g. "PNG"), pixel dimensions (e.g. "2400 × 3200 px") when width and height are available, position counter (e.g. "1 of 4"), a Download button, and a Close (×) button

#### Scenario: Top bar — dimensions absent
- **WHEN** the viewer displays a previewable file that has no width or height
- **THEN** the top bar SHALL show filename, type badge, position counter, Download, and Close — the pixel dimensions field SHALL be omitted

---

### Requirement: FileViewer component — no-preview state
The system SHALL render a fallback card when a file cannot be previewed.

#### Scenario: No-preview fallback card
- **WHEN** the viewer displays a file whose `thumbnailUrl` is null (e.g. PSD, PDF)
- **THEN** the viewer SHALL render a centered card showing a file icon, the filename, a type descriptor (e.g. "Photoshop document · no preview available"), and a large "Download file" button; the top-bar Download button SHALL remain visible and functional

---

### Requirement: FileViewer navigation
The system SHALL allow the user to navigate between files in the viewer.

#### Scenario: Next file via arrow
- **WHEN** the user clicks the next (›) arrow and the current file is not the last
- **THEN** the viewer SHALL advance to the next file and update the top bar and preview area accordingly

#### Scenario: Previous file via arrow
- **WHEN** the user clicks the prev (‹) arrow and the current file is not the first
- **THEN** the viewer SHALL go back to the previous file

#### Scenario: Arrow visibility at boundaries
- **WHEN** the viewer is on the first file
- **THEN** the prev arrow SHALL be hidden; when on the last file the next arrow SHALL be hidden

#### Scenario: Navigate via thumbnail filmstrip
- **WHEN** the user clicks a thumbnail in the filmstrip at the bottom
- **THEN** the viewer SHALL jump to that file's index

#### Scenario: Active thumbnail highlighted
- **WHEN** a file is the currently displayed file
- **THEN** its filmstrip thumbnail SHALL have a visible active indicator (e.g. ring or border)

#### Scenario: Close viewer
- **WHEN** the user clicks the Close (×) button or presses Escape
- **THEN** the viewer SHALL unmount and the artpiece detail view SHALL be visible again

---

### Requirement: FileViewer download
The system SHALL allow the user to download the currently displayed file directly from the viewer.

#### Scenario: Download via top-bar button
- **WHEN** the user clicks the Download button in the top bar
- **THEN** the system SHALL call `GET /files/:id/download`, receive `{ download_url }`, and open that URL to trigger a browser file download

#### Scenario: Download via no-preview card button
- **WHEN** the user clicks "Download file" in the no-preview fallback card
- **THEN** the system SHALL perform the same download action as the top-bar Download button
