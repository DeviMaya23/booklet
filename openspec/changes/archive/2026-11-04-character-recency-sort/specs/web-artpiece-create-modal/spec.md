# Spec Delta

## ADDED Requirements

### Requirement: Character suggestions in recency order
The characters field in the create artpiece form SHALL request the character list with `sort=recent` and display suggestions in the order returned, with no client-side re-sort. Already-selected characters SHALL be excluded from the suggestions.

#### Scenario: Suggestions follow API recency order
- **WHEN** the user focuses the characters field and the API returns characters in recency order A, B, C
- **THEN** the suggestions SHALL be listed A, B, C

#### Scenario: List requested with sort=recent
- **WHEN** the create artpiece modal opens
- **THEN** the characters request SHALL include `sort=recent`

#### Scenario: Characters refresh after saving an artpiece
- **WHEN** the user saves an artpiece with characters selected
- **THEN** the characters query SHALL be invalidated so the next picker shows the updated order
