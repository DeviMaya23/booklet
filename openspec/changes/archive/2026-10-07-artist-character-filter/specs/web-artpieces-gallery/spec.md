## MODIFIED Requirements

### Requirement: Filter popover
The gallery SHALL display a Filter button beside the search input that opens a floating panel containing artist and character filters. The Filter button and panel are provided by the shared `ArtistCharacterFilter` component.

#### Scenario: Filter button opens popover
- **WHEN** the user clicks the Filter button
- **THEN** a popover panel SHALL appear containing an artist section and a characters section

#### Scenario: Filter button closes popover
- **WHEN** the popover is open and the user clicks outside the popover or presses Escape
- **THEN** the popover SHALL close

#### Scenario: Filter button shows active filter summary
- **WHEN** at least one filter is active (artist selected or one or more character chips present)
- **THEN** the Filter button SHALL display a short summary (e.g. "1 artist", "2 characters, all", "1 artist · 2 characters, any") instead of the plain "Filter" label

---

### Requirement: Character filter
The filter panel SHALL contain a character multi-select (chip-based, using `TokenInput`) with an All/Any segmented toggle. The toggle SHALL be visible only when two or more characters are selected.

#### Scenario: Add a character chip
- **WHEN** the user selects a character in the character input
- **THEN** a chip for that character SHALL appear and the gallery SHALL be filtered accordingly

#### Scenario: Remove a character chip
- **WHEN** the user clicks the remove button on a character chip
- **THEN** that character is removed from the filter

#### Scenario: All mode (AND)
- **WHEN** the toggle is set to "All" and two or more character chips are present
- **THEN** the gallery SHALL show only artpieces whose `characters` array contains ALL of the selected character IDs

#### Scenario: Any mode (OR)
- **WHEN** the toggle is set to "Any" and two or more character chips are present
- **THEN** the gallery SHALL show only artpieces whose `characters` array contains AT LEAST ONE of the selected character IDs

#### Scenario: Toggle hidden with fewer than two characters
- **WHEN** fewer than two character chips are present
- **THEN** the All/Any toggle SHALL NOT be rendered

#### Scenario: Default toggle state
- **WHEN** the filter panel is first opened with no prior character selection
- **THEN** the toggle SHALL default to "All"

#### Scenario: No character chips — no character filtering
- **WHEN** no character chips are present
- **THEN** no character filter is applied to the gallery results
