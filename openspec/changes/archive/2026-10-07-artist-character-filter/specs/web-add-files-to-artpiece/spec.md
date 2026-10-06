## MODIFIED Requirements

### Requirement: Filters section collapsed by default
The modal SHALL display a Filter button beside the artpiece search input. Clicking it opens the shared `ArtistCharacterFilter` panel. The filter panel is NOT shown inline; it is always accessed via the button regardless of screen size.

#### Scenario: Filter button opens panel
- **WHEN** the user clicks the Filter button beside the search input
- **THEN** the filter panel SHALL open showing the Artist combobox and Characters token input

#### Scenario: Filter button closes panel
- **WHEN** the filter panel is open and the user clicks outside it or presses Escape
- **THEN** the panel SHALL close without resetting filter values

#### Scenario: Filter button shows active filter summary
- **WHEN** at least one filter is active
- **THEN** the Filter button SHALL display a summary label instead of "Filter"

---

### Requirement: Characters filter
The filter panel SHALL contain a multi-select Characters token input. When one or more characters are selected, the artpiece list SHALL be filtered to include only artpieces whose `characters` array contains AT LEAST ONE of the selected character IDs (ANY match). No All/Any toggle is shown.

#### Scenario: Characters filter — ANY match
- **WHEN** the user selects one or more characters
- **THEN** the artpiece list SHALL include only artpieces featuring at least one of the selected characters (ANY logic)

#### Scenario: Characters filter — clear
- **WHEN** the user removes all character chips
- **THEN** the character filter is removed and all artpieces (subject to other active filters) are shown
