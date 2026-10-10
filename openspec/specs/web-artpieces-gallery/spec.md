## Purpose

Defines the Artpieces gallery screen — a dedicated page for browsing all artpieces with client-side search, sort, and filter, an entry point to create a new artpiece, and a lightbox for previewing cover images.

## Requirements

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

### Requirement: Search by title
The gallery SHALL filter the displayed artpieces by title as the user types in the search input.

#### Scenario: Search matches title substring
- **WHEN** the user types a string in the search input
- **THEN** the gallery SHALL display only artpieces whose title contains the string as a case-insensitive substring

#### Scenario: Artpiece with null title during search
- **WHEN** the user has typed a search string and an artpiece has a null title
- **THEN** that artpiece SHALL NOT appear in the results

#### Scenario: Clearing search restores all results
- **WHEN** the user clears the search input
- **THEN** all artpieces (subject to any active filters) SHALL be displayed again

---

### Requirement: Sort
The gallery SHALL allow the user to sort artpieces via a sort dropdown.

#### Scenario: Sort alphabetically
- **WHEN** the user selects "Alphabetical" from the sort dropdown
- **THEN** artpieces SHALL be ordered by title ascending (case-insensitive); artpieces with null titles SHALL appear last

#### Scenario: Sort by newest
- **WHEN** the user selects "Newest first" from the sort dropdown
- **THEN** artpieces SHALL be ordered by `created_at` descending; this is the default order

---

### Requirement: Filter popover
The gallery SHALL display a Filter button beside the search input that opens a floating panel containing artist and character filters. The Filter button and panel are provided by the shared `ArtistCharacterFilter` component.

#### Scenario: Filter button opens popover
- **WHEN** the user clicks the Filter button
- **THEN** a popover panel SHALL appear containing an artist section and a characters section

#### Scenario: Filter button closes popover
- **WHEN** the popover is open and the user clicks outside the popover or presses Escape
- **THEN** the popover SHALL close

#### Scenario: Filter button shows active filter summary
- **WHEN** at least one filter is active (artist selected or one or more character chips present)
- **THEN** the Filter button SHALL display a short summary (e.g. "1 artist", "2 characters, all", "1 artist · 2 characters, any") instead of the plain "Filter" label

---

### Requirement: Artist filter
The filter popover SHALL contain a single-select artist combobox. When an artist is selected, only artpieces by that artist SHALL be shown.

#### Scenario: Select an artist
- **WHEN** the user selects an artist from the artist combobox in the filter panel
- **THEN** the gallery SHALL show only artpieces whose `artist_id` matches that artist

#### Scenario: Clear artist filter
- **WHEN** the user clears the artist selection in the combobox
- **THEN** the artist filter SHALL be removed and artpieces from all artists SHALL be eligible for display

---

### Requirement: Character filter
The filter panel SHALL contain a character multi-select (chip-based, using `TokenInput`) with an All/Any segmented toggle. The toggle SHALL be visible only when two or more characters are selected.

#### Scenario: Add a character chip
- **WHEN** the user selects a character in the character input
- **THEN** a chip for that character SHALL appear and the gallery SHALL be filtered accordingly

#### Scenario: Remove a character chip
- **WHEN** the user clicks the remove button on a character chip
- **THEN** that character is removed from the filter

#### Scenario: All mode (AND)
- **WHEN** the toggle is set to "All" and two or more character chips are present
- **THEN** the gallery SHALL show only artpieces whose `characters` array contains ALL of the selected character IDs

#### Scenario: Any mode (OR)
- **WHEN** the toggle is set to "Any" and two or more character chips are present
- **THEN** the gallery SHALL show only artpieces whose `characters` array contains AT LEAST ONE of the selected character IDs

#### Scenario: Toggle hidden with fewer than two characters
- **WHEN** fewer than two character chips are present
- **THEN** the All/Any toggle SHALL NOT be rendered

#### Scenario: Default toggle state
- **WHEN** the filter panel is first opened with no prior character selection
- **THEN** the toggle SHALL default to "All"

#### Scenario: No character chips — no character filtering
- **WHEN** no character chips are present
- **THEN** no character filter is applied to the gallery results

---

### Requirement: Combined filters
When multiple filters are active simultaneously, they SHALL be AND-combined.

#### Scenario: Artist and character filters both active
- **WHEN** the user has selected an artist AND one or more characters
- **THEN** the gallery SHALL show only artpieces that satisfy BOTH the artist filter AND the character filter

---

### Requirement: New Artpiece entry point
The gallery SHALL display a "+ New Artpiece" button that opens the existing Create New Artpiece modal with an empty file list.

#### Scenario: Open Create New Artpiece modal
- **WHEN** the user clicks "+ New Artpiece"
- **THEN** the `ArtpieceFormModal` SHALL open with no pre-populated files

#### Scenario: Artpiece list refreshes after creation
- **WHEN** the user successfully creates an artpiece via the modal
- **THEN** the gallery SHALL reflect the new artpiece without a full page reload

---

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
- **THEN** the system SHALL navigate to the artpiece detail view for that artpiece

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
