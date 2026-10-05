## Purpose

Defines the Artpieces gallery screen — a dedicated page for browsing all artpieces with client-side search, sort, and filter, an entry point to create a new artpiece, and the ability to delete an artpiece from the gallery.

## ADDED Requirements

### Requirement: Artpieces gallery page
The system SHALL render an Artpieces gallery page at `/app/artpieces` that displays all artpieces belonging to the authenticated user as a card grid.

#### Scenario: Page renders with artpieces
- **WHEN** an authenticated user navigates to `/app/artpieces`
- **THEN** the page SHALL display a card grid of all artpieces, each showing the cover thumbnail (or a placeholder icon if none), the title (or "Untitled" if null), and the artist name below the title (omitted if no artist)

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
The gallery SHALL display a Filter button that opens a popover panel containing artist and character filters.

#### Scenario: Filter button opens popover
- **WHEN** the user clicks the Filter button
- **THEN** a popover panel SHALL appear containing an artist section and a characters section

#### Scenario: Filter button closes popover
- **WHEN** the popover is open and the user clicks outside the popover or the Filter button again
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
The filter popover SHALL contain a character multi-select (chip-based, using `TokenInput`) with an All/Any segmented toggle.

#### Scenario: Add a character chip
- **WHEN** the user selects a character in the character input
- **THEN** a chip for that character SHALL appear and the gallery SHALL be filtered accordingly

#### Scenario: Remove a character chip
- **WHEN** the user clicks the remove button on a character chip
- **THEN** that character is removed from the filter

#### Scenario: All mode (AND)
- **WHEN** the toggle is set to "All" and one or more character chips are present
- **THEN** the gallery SHALL show only artpieces whose `characters` array contains ALL of the selected character IDs

#### Scenario: Any mode (OR)
- **WHEN** the toggle is set to "Any" and one or more character chips are present
- **THEN** the gallery SHALL show only artpieces whose `characters` array contains AT LEAST ONE of the selected character IDs

#### Scenario: Default toggle state
- **WHEN** the filter panel is first opened with no prior character selection
- **THEN** the toggle SHALL default to "All"

#### Scenario: No character chips — toggle has no effect
- **WHEN** no character chips are present
- **THEN** the All/Any toggle SHALL have no effect on the displayed results

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

### Requirement: Delete artpiece from gallery
The gallery SHALL allow the user to delete an artpiece via the card's overflow menu. Deleting an artpiece does NOT delete its files — files remain in the inbox with no artpiece association.

#### Scenario: Delete triggers confirmation dialog
- **WHEN** the user clicks "Delete" in a card's overflow menu
- **THEN** a confirmation dialog SHALL appear asking the user to confirm deletion

#### Scenario: Confirmed deletion removes the card
- **WHEN** the user confirms deletion
- **THEN** `DELETE /artpieces/:id` SHALL be called; on success the card SHALL be removed from the gallery and a toast SHALL show "Artpiece deleted"

#### Scenario: Cancelled deletion does nothing
- **WHEN** the user dismisses the confirmation dialog
- **THEN** no API call is made and the gallery is unchanged

#### Scenario: Deletion failure shows error toast
- **WHEN** the `DELETE /artpieces/:id` call fails
- **THEN** the card SHALL remain in the gallery and a toast SHALL show "Failed to delete artpiece"
