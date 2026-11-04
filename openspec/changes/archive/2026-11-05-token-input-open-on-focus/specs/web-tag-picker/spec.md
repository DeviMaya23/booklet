## ADDED Requirements

### Requirement: Focus suggestion list is capped
When the dropdown is open with empty input text, the picker SHALL list at most the first 10 unselected tags in the order returned. If more unselected tags exist, it SHALL show a non-selectable hint row indicating that typing searches the rest. Once the user has typed text, matching suggestions SHALL NOT be capped.

#### Scenario: Long list is truncated on focus
- **WHEN** the user has 25 tags, none selected, and focuses the empty input
- **THEN** the first 10 tags in API order SHALL be listed, followed by a hint row indicating more tags can be found by typing

#### Scenario: Short list has no hint
- **WHEN** the user has 10 or fewer unselected tags and focuses the empty input
- **THEN** all of them SHALL be listed and no hint row SHALL be shown

#### Scenario: Selected tags do not use up the cap
- **WHEN** the user has 25 tags, 3 of them selected, and focuses the empty input
- **THEN** 10 unselected tags SHALL be listed

#### Scenario: Typing searches beyond the cap
- **WHEN** the user types text that matches a tag that was not among the first 10
- **THEN** that tag SHALL be listed

#### Scenario: Hint row is not selectable
- **WHEN** the hint row is shown and the user navigates with the arrow keys or clicks the hint row
- **THEN** the hint row SHALL NOT be highlighted or selected and no tag SHALL be added
