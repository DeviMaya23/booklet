# Spec Delta

## MODIFIED Requirements

### Requirement: Filter summary chip
The modal SHALL display a read-only chip summarising the active filter state in human-readable form, including the verse when active (e.g. "Artist: Chev · Characters: Tanelle, Jielle (any) · Verse: Elowen").

#### Scenario: Summary reflects active filters
- **WHEN** the modal opens with artist "Chev" and characters "Tanelle" and "Jielle" with match mode "any"
- **THEN** the chip SHALL display text representing those active filters

#### Scenario: Summary includes active verse
- **WHEN** the modal opens with verse "Elowen" active (alone or alongside other filters)
- **THEN** the chip SHALL include the verse's name in the summary text
