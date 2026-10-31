# Spec Delta

## ADDED Requirements

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
