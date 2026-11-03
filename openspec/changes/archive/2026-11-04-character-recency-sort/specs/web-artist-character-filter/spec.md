# Spec Delta

## ADDED Requirements

### Requirement: ArtistCharacterFilter — character suggestions in recency order
The Characters token input SHALL request the character list with `sort=recent` and display suggestions in the order returned, with no client-side re-sort, in every context the component is used (gallery and commission modal).

#### Scenario: Suggestions follow API recency order
- **WHEN** the filter panel is open and the API returns characters in recency order A, B, C
- **THEN** the suggestions SHALL be listed A, B, C, excluding already-selected characters

#### Scenario: List requested with sort=recent
- **WHEN** the filter panel renders the Characters field
- **THEN** the characters request SHALL include `sort=recent`
