# Spec Delta

## MODIFIED Requirements

### Requirement: FileViewer component — image preview state
The system SHALL provide a `FileViewer` component that renders as a fullscreen modal overlay and displays a single file at a time with navigation controls.

#### Scenario: Open viewer on file tile click (artpiece detail)
- **WHEN** a user clicks a file tile in the artpiece detail view (view mode)
- **THEN** `FileViewer` SHALL open as a fullscreen overlay, displaying the clicked file at the initial index

#### Scenario: Open viewer from gallery lightbox
- **WHEN** a user single-clicks an artpiece tile in the gallery
- **THEN** `FileViewer` SHALL open as a fullscreen overlay, displaying that artpiece's cover as the initial entry, with the full list of filtered artpiece covers available for navigation

#### Scenario: Image preview with checkerboard background
- **WHEN** the viewer displays a file whose `mimeType` starts with `image/` and whose `thumbnailUrl` is non-null
- **THEN** the image SHALL be rendered in a centered tile with a CSS checkerboard background so transparent areas are visible against the dark scrim

#### Scenario: Top bar — image state
- **WHEN** the viewer displays a previewable file
- **THEN** the top bar SHALL show: filename (or "Untitled" if null), type badge (e.g. "PNG"), pixel dimensions (e.g. "2400 × 3200 px") when width and height are available, position counter (e.g. "1 of 4"), a Download button, and a Close (×) button

#### Scenario: Top bar — dimensions absent
- **WHEN** the viewer displays a previewable file that has no width or height
- **THEN** the top bar SHALL show filename, type badge, position counter, Download, and Close — the pixel dimensions field SHALL be omitted
