## MODIFIED Requirements

### Requirement: Artpieces section
The form SHALL include a collapsible "Artpieces" section, collapsed by default. The section contains: a thumbnail strip of attached/selected artpieces (each with an × button to remove), a title search input with a Filter button beside it, and a search results list. The inline filter disclosure is replaced by the shared `ArtistCharacterFilter` panel opened via the Filter button. No All/Any toggle is shown. Collapsing and expanding the section does not reset search, filter, or strip state.

#### Scenario: Filter button opens panel
- **WHEN** the artpieces section is expanded and the user clicks the Filter button
- **THEN** the filter panel SHALL open showing the Artist combobox and Characters token input

#### Scenario: Filter button closes panel
- **WHEN** the filter panel is open and the user clicks outside it or presses Escape
- **THEN** the panel SHALL close without resetting filter values

#### Scenario: Filter button shows active filter summary
- **WHEN** at least one artpiece filter is active
- **THEN** the Filter button SHALL display a summary label instead of "Filter"

#### Scenario: Characters filter — ANY match
- **WHEN** the user selects one or more characters in the filter panel
- **THEN** the artpiece search results SHALL include only artpieces featuring at least one of the selected characters (ANY logic)

#### Scenario: Section collapsed by default
- **WHEN** the modal opens (create or edit)
- **THEN** the artpieces section is collapsed and the strip, search, and filter are not visible

#### Scenario: Expand and collapse preserves state
- **WHEN** a user expands the section, applies a filter, then collapses and re-expands
- **THEN** the filter state is preserved

#### Scenario: Empty state
- **WHEN** no artpieces are attached or selected
- **THEN** the strip area shows "No artpieces attached yet"

#### Scenario: Remove artpiece from strip
- **WHEN** a user clicks the × on an artpiece thumbnail in the strip
- **THEN** the artpiece is removed from the strip immediately (not yet persisted)

#### Scenario: Search results exclude artpieces attached to a different commission
- **WHEN** a user searches for artpieces in the section
- **THEN** artpieces already attached to a different commission do not appear in the results
