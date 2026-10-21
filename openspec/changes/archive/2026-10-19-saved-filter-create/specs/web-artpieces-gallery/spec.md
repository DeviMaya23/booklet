# Spec Delta

## ADDED Requirements

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
