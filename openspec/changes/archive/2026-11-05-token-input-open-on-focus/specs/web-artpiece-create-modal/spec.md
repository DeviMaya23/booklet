## ADDED Requirements

### Requirement: Characters field lists suggestions on focus
The characters field in the create artpiece form SHALL list suggestions when focused with empty text, in recency order, excluding already-selected characters. The empty-text list SHALL be capped at 10 with a non-selectable hint row when more exist, exactly as for the tag picker. This SHALL hold whether the modal was opened from the file Inbox or the Artpieces page.

#### Scenario: Focus lists recent characters
- **WHEN** the user focuses the empty characters field and the API returns characters A, B, C in recency order
- **THEN** A, B, C SHALL be listed in that order

#### Scenario: Focus list is capped
- **WHEN** the user has 25 characters and focuses the empty characters field
- **THEN** the first 10 SHALL be listed followed by a hint row indicating more can be found by typing

#### Scenario: Opened from the Inbox
- **WHEN** the modal is opened from the file Inbox and the user focuses the empty characters field
- **THEN** characters SHALL be listed as above

#### Scenario: Typing searches all characters
- **WHEN** the user types text matching a character beyond the first 10
- **THEN** that character SHALL be listed
