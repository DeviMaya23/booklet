## ADDED Requirements

### Requirement: Edit mode — characters field lists suggestions on focus
The characters field in edit mode SHALL list suggestions when focused with empty text, in recency order, excluding already-selected characters. The empty-text list SHALL be capped at 10 with a non-selectable hint row when more exist, exactly as for the tag picker.

#### Scenario: Focus lists recent characters
- **WHEN** the user focuses the empty characters field in edit mode and the API returns characters A, B in recency order
- **THEN** A, B SHALL be listed in that order

#### Scenario: Focus list is capped
- **WHEN** the user has 25 characters and focuses the empty characters field in edit mode
- **THEN** the first 10 SHALL be listed followed by a hint row indicating more can be found by typing

#### Scenario: Selected characters excluded
- **WHEN** a character is already a chip and the user focuses the empty field
- **THEN** that character SHALL NOT be listed
