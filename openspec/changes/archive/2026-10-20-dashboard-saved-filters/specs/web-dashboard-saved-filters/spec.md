# Spec Delta

## Purpose

Defines the Saved Filters section on the dashboard — a horizontally scrolling tile row of the user's saved filters, with click-to-navigate, per-tile delete and change-thumbnail actions, and a visibility toggle that controls which filters appear.

## ADDED Requirements

### Requirement: Saved filters tile row
The dashboard SHALL display a "Saved filters" section above the Recent artpieces section. Each tile SHALL show the saved filter's thumbnail (or a placeholder) and name. The section SHALL fetch filters from `GET /saved_filters`.

#### Scenario: Saved filters render on dashboard
- **WHEN** the dashboard loads and the user has at least one saved filter
- **THEN** each visible saved filter SHALL be displayed as a thumbnail tile with its name

#### Scenario: Empty state when no saved filters exist
- **WHEN** the user has no saved filters
- **THEN** the Saved filters section SHALL display an empty state message

---

### Requirement: Click navigates to filtered artpiece gallery
Clicking a saved filter tile SHALL navigate to `/app/artpieces` with the filter's `artist_id`, `character_ids`, and `match_mode` mapped to the URL params `artist`, `characters`, and `match` respectively.

#### Scenario: Valid filter payload navigates with params
- **WHEN** a user clicks a saved filter tile whose `filter_payload` is a valid shape
- **THEN** the app SHALL navigate to `/app/artpieces` with the corresponding URL params applied

#### Scenario: Invalid or missing filter payload navigates without params
- **WHEN** a saved filter's `filter_payload` cannot be parsed into the expected shape
- **THEN** the app SHALL navigate to `/app/artpieces` with no filter params applied, rather than erroring

---

### Requirement: Per-tile action menu
Each saved filter tile SHALL expose a contextual action menu with two items: "Change thumbnail" and "Delete".

#### Scenario: Action menu is present on each tile
- **WHEN** a user hovers over a saved filter tile
- **THEN** an action button SHALL become visible that opens a menu with "Change thumbnail" and "Delete"

---

### Requirement: Delete saved filter
Selecting "Delete" from a tile's action menu SHALL open a confirmation dialog. Confirming SHALL call `DELETE /saved_filters/:id`. On success the tile SHALL be removed from the row.

#### Scenario: Delete removes tile on success
- **WHEN** a user confirms deletion of a saved filter
- **THEN** the system SHALL call `DELETE /saved_filters/:id` and the tile SHALL disappear from the dashboard

#### Scenario: Cancel leaves tile intact
- **WHEN** a user opens the delete dialog but clicks Cancel
- **THEN** no DELETE request SHALL be sent and the tile SHALL remain

---

### Requirement: Change thumbnail
Selecting "Change thumbnail" from a tile's action menu SHALL open a modal. The modal SHALL show the filter summary and a selectable grid of up to 10 artpieces that match the filter (from cached `GET /artpieces` data, with the filter applied). The first artpiece SHALL be pre-selected. Saving SHALL call `PATCH /saved_filters/:id` with `thumbnail_artpiece_id`. On success the tile's thumbnail SHALL update.

#### Scenario: Modal shows filtered artpiece grid
- **WHEN** the change thumbnail modal opens for a saved filter
- **THEN** it SHALL display up to 10 artpieces matching the filter's payload as a selectable grid

#### Scenario: First artpiece is pre-selected
- **WHEN** the change thumbnail modal opens
- **THEN** the first artpiece in the grid SHALL be selected by default

#### Scenario: No matching artpieces shows empty state
- **WHEN** the saved filter's payload matches no artpieces
- **THEN** the modal SHALL display an empty state and the Save button SHALL be disabled

#### Scenario: Save patches thumbnail
- **WHEN** the user selects an artpiece and clicks Save
- **THEN** the system SHALL call `PATCH /saved_filters/:id` with `thumbnail_artpiece_id` set to the selected artpiece's ID, and the tile SHALL update its thumbnail on success

#### Scenario: Cancel discards selection
- **WHEN** the user clicks Cancel
- **THEN** the modal SHALL close without sending a PATCH request

---

### Requirement: Visibility toggle
A "Choose shown" button in the Saved filters section header SHALL open a popover listing all saved filters with a checkbox per filter. Unchecking a filter hides it from the tile row. Visibility state SHALL be persisted in `localStorage`. All filters are shown by default; a newly created filter is visible until explicitly hidden.

#### Scenario: All filters shown by default
- **WHEN** a user has never used the visibility toggle
- **THEN** all saved filters SHALL be visible on the dashboard

#### Scenario: Hiding a filter removes its tile
- **WHEN** a user unchecks a filter in the Choose shown popover
- **THEN** that filter's tile SHALL no longer appear in the row

#### Scenario: Hidden state persists across page loads
- **WHEN** a user hides a filter and reloads the dashboard
- **THEN** the hidden filter SHALL remain hidden

#### Scenario: Empty state when all filters are hidden
- **WHEN** the user has saved filters but has hidden all of them
- **THEN** the Saved filters section SHALL display an empty state indicating no filters are shown, rather than disappearing entirely
