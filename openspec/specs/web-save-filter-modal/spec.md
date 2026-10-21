# Web Save Filter Modal Specification

## Purpose

Defines the Save Filter modal — opened from the artpiece gallery filter popover, it lets the user name a saved filter, pick a thumbnail from the current filtered artpiece results, and persist it to the dashboard.

## Requirements

### Requirement: Filter summary chip
The modal SHALL display a read-only chip summarising the active filter state in human-readable form (e.g. "Artist: Chev · Characters: Tanelle, Jielle (any)").

#### Scenario: Summary reflects active filters
- **WHEN** the modal opens with artist "Chev" and characters "Tanelle" and "Jielle" with match mode "any"
- **THEN** the chip SHALL display text representing those active filters

---

### Requirement: Name input
The modal SHALL contain a required text input for the saved filter name.

#### Scenario: Empty name blocks save
- **WHEN** the name input is empty
- **THEN** the Save button SHALL be disabled

#### Scenario: Duplicate name shows inline error
- **WHEN** the user submits and the server returns 422 for a duplicate name
- **THEN** the name input SHALL display an inline error "A saved filter named "<name>" already exists." and the modal SHALL remain open

---

### Requirement: Thumbnail picker
The modal SHALL display thumbnails of up to 10 artpieces matching the current filter as a selectable grid. Exactly one cell is selected at a time; the first artpiece is pre-selected on open.

#### Scenario: Grid shows filtered artpiece thumbnails
- **WHEN** the modal opens and filtered artpieces exist
- **THEN** the grid SHALL display thumbnails (or a placeholder icon for artpieces without one) for up to 10 filtered artpieces, with the first cell pre-selected

#### Scenario: Selecting a thumbnail
- **WHEN** the user clicks a thumbnail cell
- **THEN** that cell becomes selected and the dashboard preview updates to reflect it

---

### Requirement: No-results blocking
When no artpieces match the current filter the modal SHALL show an empty-state message instead of the thumbnail grid, and the Save button SHALL be disabled.

#### Scenario: Empty state when no results
- **WHEN** the modal opens and no artpieces match the current filter
- **THEN** the thumbnail grid SHALL be replaced by "No artpieces match this filter" and the Save button SHALL be disabled

---

### Requirement: Dashboard preview
The modal SHALL show a preview tile of how the saved filter will appear on the dashboard: the selected thumbnail alongside the entered name.

#### Scenario: Preview reflects current name and thumbnail
- **WHEN** the user changes the name or selects a different thumbnail
- **THEN** the preview tile SHALL immediately update to show the new name and thumbnail

---

### Requirement: Save action
Clicking Save POSTs to `/saved_filters` with the name, `thumbnail_artpiece_id`, and `filter_payload`. On success the modal closes and the caller receives the created saved filter. On a non-422 error a toast fires and the modal stays open.

#### Scenario: Successful save closes modal
- **WHEN** the server returns 201
- **THEN** the modal SHALL close and a success toast SHALL appear

#### Scenario: Server error shows toast
- **WHEN** the server returns a non-422 error
- **THEN** the modal SHALL remain open and an error toast SHALL appear

#### Scenario: Cancel discards without saving
- **WHEN** the user clicks Cancel
- **THEN** the modal SHALL close without POSTing to the server
