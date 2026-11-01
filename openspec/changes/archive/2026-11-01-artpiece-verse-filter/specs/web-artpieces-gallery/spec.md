# Spec Delta

## ADDED Requirements

### Requirement: Verse filter
The filter panel SHALL contain a single-select verse combobox. An artpiece matches the active verse filter when at least one of its tagged characters has `verse_id` equal to the selected verse's id. Characters with no verse assigned never match a verse filter.

#### Scenario: Select a verse
- **WHEN** the user selects a verse from the verse combobox in the filter panel
- **THEN** the gallery SHALL show only artpieces that have at least one tagged character whose `verse_id` matches that verse

#### Scenario: Clear verse filter
- **WHEN** the user clears the verse selection in the combobox
- **THEN** the verse filter SHALL be removed and artpieces from all verses (and artpieces with no verse-assigned characters) SHALL be eligible for display

#### Scenario: Artpiece with multiple characters from different verses
- **WHEN** an artpiece has characters from both the selected verse and other verses
- **THEN** that artpiece SHALL still match the verse filter

#### Scenario: Character with no verse never matches
- **WHEN** a verse filter is active and an artpiece's only tagged characters have no `verse_id` assigned
- **THEN** that artpiece SHALL NOT match the verse filter

### Requirement: Verse filter data resolves against the current verse list
The verse filter's selected value SHALL be resolved by looking up the stored/URL verse id against the currently fetched verse list. When the id does not resolve (e.g. the verse was deleted), the verse filter SHALL be treated as not-selected rather than erroring.

#### Scenario: Stored verse id no longer exists
- **WHEN** the URL or `localStorage` carries a `verse_id` that does not match any verse currently returned for the user
- **THEN** the gallery SHALL behave as if no verse filter is active, without error

## MODIFIED Requirements

### Requirement: Filter popover
The gallery SHALL display a Filter button beside the search input that opens a floating panel containing artist, character, and verse filters, and a "References only" toggle. The Filter button and artist/character/verse controls are provided by the shared `ArtistCharacterFilter` component.

#### Scenario: Filter button opens popover
- **WHEN** the user clicks the Filter button
- **THEN** a popover panel SHALL appear containing an artist section, a characters section, a verse section, and a "References only" toggle

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

#### Scenario: Filter button shows active filter summary — verse
- **WHEN** a verse is active, regardless of other filters
- **THEN** the Filter button summary SHALL include it (e.g. "1 verse", "1 artist · 1 verse")

### Requirement: Combined filters
When multiple filters are active simultaneously, they SHALL be AND-combined.

#### Scenario: Artist and character filters both active
- **WHEN** the user has selected an artist AND one or more characters
- **THEN** the gallery SHALL show only artpieces that satisfy BOTH the artist filter AND the character filter

#### Scenario: Verse filter combined with artist or character filters
- **WHEN** the user has selected a verse AND an artist, or a verse AND one or more characters
- **THEN** the gallery SHALL show only artpieces that satisfy ALL active filters simultaneously

### Requirement: Filter state persisted in URL
The gallery SHALL sync the active artist filter, character filter, verse filter, and match mode to URL search params so that navigating away and back restores the same filter state.

#### Scenario: Active filter is reflected in the URL
- **WHEN** the user selects an artist, characters, or a verse in the filter panel
- **THEN** the URL SHALL update to include the corresponding search params without a full page reload

#### Scenario: Filter state is restored on navigation
- **WHEN** the user navigates to the artpieces gallery and the URL contains filter params
- **THEN** the gallery SHALL initialize with those filter values active, including the verse filter

#### Scenario: Clearing filters removes URL params
- **WHEN** the user clears all filters (via "Clear all" or by removing the last selection)
- **THEN** the filter params SHALL be removed from the URL, including the verse param

### Requirement: Structural filter state persisted across navigation
The gallery SHALL persist the structural filter state (artist ID, character IDs, verse ID, match mode, isReference) to `localStorage` under the key `artpieces.filterState` so that navigating away via the sidebar and returning restores the last-used filter. Search and sort are excluded from this persistence.

#### Scenario: Filter restored after navigating away via sidebar
- **WHEN** the user has an active filter (artist, characters, verse, or references toggle), navigates to another page via the sidebar, and returns to `/app/artpieces` with no URL filter params
- **THEN** the gallery SHALL restore the last-used structural filter from localStorage, including the verse filter, and apply it

#### Scenario: Arriving from a saved-filter link updates stored state
- **WHEN** the user arrives at the artpieces gallery with URL filter params (e.g. from a dashboard saved-filter tile), including a verse param
- **THEN** the gallery SHALL write those URL params to localStorage as the new last-used filter

#### Scenario: Changing a filter writes to localStorage
- **WHEN** the user changes any structural filter (artist, characters, verse, match mode, or references toggle)
- **THEN** the new full filter state SHALL be written to localStorage

#### Scenario: Clearing all filters persists the empty state
- **WHEN** the user clears all filters (via "Clear all" or by removing the last selection)
- **THEN** localStorage SHALL be updated to reflect the empty state, so the next visit starts with no filter

#### Scenario: No stored state — default behavior unchanged
- **WHEN** the user navigates to the gallery with no URL params and no value in localStorage
- **THEN** the gallery SHALL load with no structural filter active (default behavior)
