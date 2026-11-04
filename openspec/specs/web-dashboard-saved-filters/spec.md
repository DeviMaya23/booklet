# web-dashboard-saved-filters Specification

## Purpose

Defines the Saved Filters section on the dashboard — a horizontally scrolling tile row of the user's saved filters, with click-to-navigate, per-tile delete and change-thumbnail actions, and a visibility toggle that controls which filters appear.

## Requirements

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
Clicking a saved filter tile SHALL navigate to `/app/artpieces` with the filter's `artist_id`, `character_ids`, `match_mode`, `verse_id`, `tag_ids` and `tag_match_mode` mapped to the URL params `artist`, `characters`, `match`, `verse`, `tags` and `tag_match` respectively.

#### Scenario: Valid filter payload navigates with params
- **WHEN** a user clicks a saved filter tile whose `filter_payload` is a valid shape
- **THEN** the app SHALL navigate to `/app/artpieces` with the corresponding URL params applied

#### Scenario: Invalid or missing filter payload navigates without params
- **WHEN** a saved filter's `filter_payload` cannot be parsed into the expected shape
- **THEN** the app SHALL navigate to `/app/artpieces` with no filter params applied, rather than erroring

#### Scenario: Filter payload with verse_id navigates with verse param
- **WHEN** a user clicks a saved filter tile whose `filter_payload` includes a non-null `verse_id`
- **THEN** the app SHALL navigate to `/app/artpieces` with a `verse` URL param set to that id, alongside any other active params

#### Scenario: Filter payload without verse_id navigates with no verse param
- **WHEN** a user clicks a saved filter tile whose `filter_payload` has `verse_id` absent or `null`
- **THEN** the app SHALL navigate to `/app/artpieces` with no `verse` URL param

#### Scenario: Filter payload with tags navigates with tag params
- **WHEN** a user clicks a saved filter tile whose `filter_payload` has `tag_ids` `[<A>, <B>]` and `tag_match_mode` `"any"`
- **THEN** the app SHALL navigate to `/app/artpieces` with `tags=<A>,<B>` and `tag_match=any`

#### Scenario: Filter payload with tags in All mode
- **WHEN** a user clicks a saved filter tile whose `filter_payload` has non-empty `tag_ids` and `tag_match_mode` `"all"`
- **THEN** the app SHALL navigate with the `tags` param and no `tag_match` param

#### Scenario: Filter payload without tag fields
- **WHEN** a user clicks a saved filter tile whose `filter_payload` has no `tag_ids` or `tag_match_mode`
- **THEN** the app SHALL navigate with no `tags` or `tag_match` params

#### Scenario: Saved filter references a deleted tag
- **WHEN** a user clicks a saved filter tile whose `tag_ids` includes a tag that has since been deleted
- **THEN** the gallery SHALL open with the deleted tag ignored and the remaining filters applied, without error

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
Selecting "Change thumbnail" from a tile's action menu SHALL open a modal. The modal SHALL show the filter summary and a selectable grid of up to 10 artpieces that match the filter — including its tag and verse filters, if any — (from cached `GET /artpieces`, `GET /characters`, `GET /tags` and `GET /verses` data, with the filter applied). Character, tag and verse ids that no longer exist SHALL be dropped from the filter before matching. The first artpiece SHALL be pre-selected. Saving SHALL call `PATCH /saved_filters/:id` with `thumbnail_artpiece_id`. On success the tile's thumbnail SHALL update.

#### Scenario: Modal shows filtered artpiece grid
- **WHEN** the change thumbnail modal opens for a saved filter
- **THEN** it SHALL display up to 10 artpieces matching the filter's payload as a selectable grid

#### Scenario: Modal grid respects an active verse filter
- **WHEN** the change thumbnail modal opens for a saved filter whose `filter_payload` includes a non-null `verse_id`
- **THEN** the candidate grid SHALL only include artpieces that have at least one tagged character belonging to that verse

#### Scenario: Modal grid respects an active tag filter
- **WHEN** the change thumbnail modal opens for a saved filter whose `filter_payload` has `tag_ids` `[<A>, <B>]` and `tag_match_mode` `"all"`
- **THEN** the candidate grid SHALL only include artpieces that have both tag A and tag B

#### Scenario: Modal grid with tag match Any
- **WHEN** the saved filter's `tag_match_mode` is `"any"`
- **THEN** the candidate grid SHALL include artpieces that have at least one of the filter's tags

#### Scenario: Deleted character in All mode does not empty the grid
- **WHEN** the saved filter has `character_ids` `[<A>, <B>]`, `match_mode` `"all"`, character B no longer exists, and artpieces with character A exist
- **THEN** the candidate grid SHALL show artpieces that have character A, not the empty state

#### Scenario: Deleted tag in All mode does not empty the grid
- **WHEN** the saved filter has `tag_ids` `[<A>, <B>]`, `tag_match_mode` `"all"`, tag B no longer exists, and artpieces with tag A exist
- **THEN** the candidate grid SHALL show artpieces that have tag A, not the empty state

#### Scenario: Tag or character list failed to load
- **WHEN** the tag or character list request fails
- **THEN** the saved tag or character ids SHALL be applied as given, so the grid does not silently ignore that part of the filter

#### Scenario: Grid waits for tags and characters to load
- **WHEN** the tag or character list has not finished loading
- **THEN** the modal SHALL show its loading state rather than a grid or the empty state

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
