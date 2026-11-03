# Spec Delta

## ADDED Requirements

### Requirement: Edit mode — character suggestions in recency order
The characters field in edit mode SHALL request the character list with `sort=recent` and display suggestions in the order returned, with no client-side re-sort. Already-selected characters SHALL be excluded.

#### Scenario: Suggestions follow API recency order
- **WHEN** the user types in the characters field in edit mode and the API returns matching characters in recency order A, B
- **THEN** the suggestions SHALL be listed A, B

#### Scenario: List requested with sort=recent
- **WHEN** the detail view is open in edit mode
- **THEN** the characters request SHALL include `sort=recent`

#### Scenario: Characters refresh after saving edits
- **WHEN** the user saves artpiece edits
- **THEN** the characters query SHALL be invalidated so picker order reflects the save
