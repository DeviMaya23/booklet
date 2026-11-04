# Spec Delta

## MODIFIED Requirements

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
