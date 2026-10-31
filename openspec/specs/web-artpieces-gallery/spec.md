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
The gallery SHALL display a Filter button beside the search input that opens a floating panel containing artist and character filters, and a "References only" toggle. The Filter button and artist/character controls are provided by the shared `ArtistCharacterFilter` component.

#### Scenario: Filter button opens popover
- **WHEN** the user clicks the Filter button
- **THEN** a popover panel SHALL appear containing an artist section, a characters section, and a "References only" toggle

#### Scenario: Filter button closes popover
- **WHEN** the popover is open and the user clicks outside the popover or presses Escape
- **THEN** the popover SHALL close

#### Scenario: Filter button shows active filter summary — artist or characters
- **WHEN** at least one of artist or characters is active
- **THEN** the Filter button SHALL display a short summary (e.g. "1 artist", "2 characters, all", "1 artist · 2 characters, any") instead of the plain "Filter" label

#### Scenario: Filter button shows active filter summary — references only
- **WHEN** the "References only" toggle is on and no artist or characters are active
- **THEN** the Filter button SHALL include "references" in its summary label

#### Scenario: Filter button shows active filter summary — combined
- **WHEN** artist or characters are active AND the "References only" toggle is on
- **THEN** the Filter button summary SHALL reflect both (e.g. "1 artist · references")

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

---

### Requirement: Filter state persisted in URL
The gallery SHALL sync the active artist filter, character filter, and match mode to URL search params so that navigating away and back restores the same filter state.

#### Scenario: Active filter is reflected in the URL
- **WHEN** the user selects an artist or characters in the filter panel
- **THEN** the URL SHALL update to include the corresponding search params without a full page reload

#### Scenario: Filter state is restored on navigation
- **WHEN** the user navigates to the artpieces gallery and the URL contains filter params
- **THEN** the gallery SHALL initialize with those filter values active

#### Scenario: Clearing filters removes URL params
- **WHEN** the user clears all filters (via "Clear all" or by removing the last selection)
- **THEN** the filter params SHALL be removed from the URL

---

### Requirement: Save to dashboard action in filter popover
When at least one filter is active, the gallery filter popover SHALL display a "Save to dashboard" button in the footer alongside "Clear all". Clicking it opens the Save Filter modal with the current filter state.

#### Scenario: Button appears when filters are active
- **WHEN** at least one filter is active (artist selected or at least one character chip present)
- **THEN** the filter popover footer SHALL show a "Save to dashboard" button

#### Scenario: Button absent when no filters are active
- **WHEN** no artist is selected and no characters are selected
- **THEN** the "Save to dashboard" button SHALL NOT appear in the footer

#### Scenario: Clicking the button opens the Save Filter modal
- **WHEN** the user clicks "Save to dashboard"
- **THEN** the Save Filter modal SHALL open pre-loaded with the current filter state

---

### Requirement: Saved-as confirmation in filter popover
After a successful save from the Save Filter modal, the filter popover footer SHALL replace the "Save to dashboard" button with a "Saved as '<name>'" confirmation for as long as the filter values remain unchanged.

#### Scenario: Confirmation shown after successful save
- **WHEN** the user saves a filter and the modal closes with a 201 response
- **THEN** the filter popover footer SHALL display "Saved as '<name>'" in place of the button

#### Scenario: Confirmation clears when filter changes
- **WHEN** the filter values change after a successful save (artist or characters modified)
- **THEN** the "Saved as '<name>'" confirmation SHALL be removed and the "Save to dashboard" button SHALL reappear

---

### Requirement: References only filter toggle
The gallery filter popover SHALL contain a "References only" toggle (default off). When off the gallery shows only artpieces where `is_reference` is false. When on the gallery shows only artpieces where `is_reference` is true. All other filters (artist, characters, match mode, sort, search) apply in both states.

#### Scenario: Toggle off — only non-references shown
- **WHEN** the "References only" toggle is off (default)
- **THEN** the gallery SHALL display only artpieces where is_reference is false

#### Scenario: Toggle on — only references shown
- **WHEN** the user turns the "References only" toggle on
- **THEN** the gallery SHALL display only artpieces where is_reference is true

#### Scenario: Other filters still apply with toggle on
- **WHEN** the "References only" toggle is on and an artist filter is also active
- **THEN** the gallery SHALL show only artpieces that are both references and match the artist filter

#### Scenario: Toggle state persists in URL
- **WHEN** the user turns the "References only" toggle on and navigates away then back
- **THEN** the URL SHALL contain the reference param and the gallery SHALL restore the toggle to on

#### Scenario: Clear all resets toggle
- **WHEN** the user clicks "Clear all" in the filter popover while the "References only" toggle is on
- **THEN** the toggle SHALL be reset to off

---

### Requirement: Gallery default shows only non-references
When the gallery loads with no filter params in the URL, the "References only" toggle SHALL default to off, showing only non-reference artpieces.

#### Scenario: Default state on first load
- **WHEN** the gallery loads with no filter params in the URL
- **THEN** the "References only" toggle is off and only artpieces with is_reference false are shown

#### Scenario: Absence of ref URL param treated as off
- **WHEN** the URL does not contain the ref param
- **THEN** the gallery behaves as if the toggle is off (non-references only)

---

### Requirement: Structural filter state persisted across navigation
The gallery SHALL persist the structural filter state (artist ID, character IDs, match mode, isReference) to `localStorage` under the key `artpieces.filterState` so that navigating away via the sidebar and returning restores the last-used filter. Search and sort are excluded from this persistence.

#### Scenario: Filter restored after navigating away via sidebar
- **WHEN** the user has an active filter (artist, characters, or references toggle), navigates to another page via the sidebar, and returns to `/app/artpieces` with no URL filter params
- **THEN** the gallery SHALL restore the last-used structural filter from localStorage and apply it

#### Scenario: Arriving from a saved-filter link updates stored state
- **WHEN** the user arrives at the artpieces gallery with URL filter params (e.g. from a dashboard saved-filter tile)
- **THEN** the gallery SHALL write those URL params to localStorage as the new last-used filter

#### Scenario: Changing a filter writes to localStorage
- **WHEN** the user changes any structural filter (artist, characters, match mode, or references toggle)
- **THEN** the new full filter state SHALL be written to localStorage

#### Scenario: Clearing all filters persists the empty state
- **WHEN** the user clears all filters (via "Clear all" or by removing the last selection)
- **THEN** localStorage SHALL be updated to reflect the empty state, so the next visit starts with no filter

#### Scenario: No stored state — default behavior unchanged
- **WHEN** the user navigates to the gallery with no URL params and no value in localStorage
- **THEN** the gallery SHALL load with no structural filter active (default behavior)
