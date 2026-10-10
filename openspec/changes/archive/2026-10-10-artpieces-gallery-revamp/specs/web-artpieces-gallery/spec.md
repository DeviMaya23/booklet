# Spec Delta

## MODIFIED Requirements

### Requirement: Artpieces gallery page
The system SHALL render an Artpieces gallery page at `/app/artpieces` that displays all artpieces belonging to the authenticated user as a card grid.

#### Scenario: Page renders with artpieces
- **WHEN** an authenticated user navigates to `/app/artpieces`
- **THEN** the page SHALL display a grid of artpiece tiles, each showing the cover thumbnail (or a placeholder icon if none), with title and artist displayed in a caption row below the tile (title shown as "Untitled" if null; artist row omitted if no artist); caption row visibility is subject to the "Show details" preference

#### Scenario: Tile hover state
- **WHEN** the user hovers over an artpiece tile
- **THEN** the tile SHALL lift slightly and an ⓘ button SHALL appear in the top-right corner of the tile with an "Open details" tooltip

#### Scenario: Empty state
- **WHEN** an authenticated user has no artpieces
- **THEN** the page SHALL display an empty state

---

### Requirement: Open artpiece detail via double-click
The gallery card SHALL support double-click to open the artpiece detail view.

#### Scenario: Double-click opens detail view
- **WHEN** the user double-clicks an artpiece card in the gallery
- **THEN** the artpiece detail view SHALL open for that artpiece

## REMOVED Requirements

### Requirement: Open artpiece detail via double-click
**Reason**: Replaced by the ⓘ hover button; single click is now reserved for the cover lightbox.
**Migration**: Users open the artpiece detail view by clicking the ⓘ button that appears on tile hover.

### Requirement: Delete artpiece from gallery
**Reason**: Tile-level delete is removed to simplify the tile interaction surface. Delete remains accessible from the artpiece detail view.
**Migration**: Open the artpiece detail view and use the delete action there.

## ADDED Requirements

### Requirement: Single-click opens cover lightbox
The system SHALL open `FileViewer` as a fullscreen lightbox when the user single-clicks an artpiece tile, navigating through artpieces by their cover image.

#### Scenario: Single click opens lightbox at correct index
- **WHEN** the user single-clicks an artpiece tile
- **THEN** `FileViewer` SHALL open in fullscreen with that artpiece's cover as the active entry, and `<` / `>` arrows SHALL navigate to adjacent artpieces in the current filtered and sorted order

#### Scenario: Artpiece with no cover in lightbox
- **WHEN** the lightbox is open and the current artpiece has no cover file
- **THEN** `FileViewer` SHALL render the no-preview fallback card for that entry

#### Scenario: Closing the lightbox returns to the gallery
- **WHEN** the user closes the lightbox (× button or Escape)
- **THEN** `FileViewer` SHALL unmount and the gallery SHALL be visible and unchanged

---

### Requirement: Open artpiece detail via ⓘ button
The system SHALL navigate to the artpiece detail view when the user clicks the ⓘ button visible on tile hover.

#### Scenario: ⓘ click navigates to detail
- **WHEN** the user clicks the ⓘ button on a tile
- **THEN** the system SHALL navigate to the artpiece detail view for that artpiece (equivalent to the previous double-click behavior)

---

### Requirement: View options popover
The gallery toolbar SHALL include a **View** button that opens a popover with tile size and caption visibility controls.

#### Scenario: Tile size — Small
- **WHEN** the user selects "Small" in the View popover
- **THEN** the grid SHALL use a higher column count (denser layout)

#### Scenario: Tile size — Medium
- **WHEN** the user selects "Medium" in the View popover
- **THEN** the grid SHALL use a medium column count (default layout)

#### Scenario: Tile size — Large
- **WHEN** the user selects "Large" in the View popover
- **THEN** the grid SHALL use a lower column count (spacious layout)

#### Scenario: Show details toggle — on
- **WHEN** "Show details" is enabled in the View popover
- **THEN** the caption row (title and artist) SHALL be visible below each tile

#### Scenario: Show details toggle — off
- **WHEN** "Show details" is disabled in the View popover
- **THEN** the caption row SHALL be hidden and tiles SHALL be image-only

---

### Requirement: View preferences persistence
The gallery SHALL persist the tile size and show-details preference to `localStorage` so they survive page reloads and session restarts.

#### Scenario: Preference survives reload
- **WHEN** the user sets a tile size or show-details value and reloads the page
- **THEN** the gallery SHALL restore the saved preference on load

#### Scenario: Default when no preference stored
- **WHEN** no preference has been saved (e.g. first visit or cleared storage)
- **THEN** the gallery SHALL default to Medium tile size and show-details enabled
