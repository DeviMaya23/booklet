# Spec Delta

## MODIFIED Requirements

### Requirement: Click navigates to filtered artpiece gallery
Clicking a saved filter tile SHALL navigate to `/app/artpieces` with the filter's `artist_id`, `character_ids`, `match_mode`, and `verse_id` mapped to the URL params `artist`, `characters`, `match`, and `verse` respectively.

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

### Requirement: Change thumbnail
Selecting "Change thumbnail" from a tile's action menu SHALL open a modal. The modal SHALL show the filter summary and a selectable grid of up to 10 artpieces that match the filter — including its verse filter, if any — (from cached `GET /artpieces` and `GET /characters` data, with the filter applied). The first artpiece SHALL be pre-selected. Saving SHALL call `PATCH /saved_filters/:id` with `thumbnail_artpiece_id`. On success the tile's thumbnail SHALL update.

#### Scenario: Modal shows filtered artpiece grid
- **WHEN** the change thumbnail modal opens for a saved filter
- **THEN** it SHALL display up to 10 artpieces matching the filter's payload as a selectable grid

#### Scenario: Modal grid respects an active verse filter
- **WHEN** the change thumbnail modal opens for a saved filter whose `filter_payload` includes a non-null `verse_id`
- **THEN** the candidate grid SHALL only include artpieces that have at least one tagged character belonging to that verse

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
