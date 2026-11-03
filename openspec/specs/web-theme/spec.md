# web-theme Specification

## Purpose

Defines the user's colour theme preference: which themes exist, how the choice is persisted and and applied across the authenticated app.

## Requirements

### Requirement: Selectable themes
The system SHALL support exactly three themes identified as `warm` (displayed as "Parchment"), `lumen`, and `sunless`. The default theme SHALL be `warm`.

#### Scenario: Default when nothing is stored
- **WHEN** the app loads and no theme preference is stored
- **THEN** the active theme SHALL be `warm`

### Requirement: Theme preference persists
Changing the theme SHALL apply it to the whole authenticated app immediately and persist it in the browser so it is restored on the next load. A stored value that is not one of the three supported themes SHALL be ignored in favour of `warm`.

#### Scenario: Changing the theme applies and persists it
- **WHEN** the user selects `sunless`
- **THEN** the document's active theme SHALL become `sunless`
- **AND** the preference SHALL be stored so a reload restores `sunless`

#### Scenario: Stored preference is restored on load
- **WHEN** the app loads with `lumen` stored as the preference
- **THEN** the active theme SHALL be `lumen`

#### Scenario: Invalid stored value falls back
- **WHEN** the app loads with an unrecognised stored theme value
- **THEN** the active theme SHALL be `warm`
