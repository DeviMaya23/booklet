# Spec Delta

## ADDED Requirements

### Requirement: ArtistCharacterFilter — Verse field
The filter panel SHALL optionally contain a single-select Verse combobox, rendered only when `selectedVerse` and `onVerseChange` props are both provided by the caller. Selecting a verse calls `onVerseChange` with the chosen verse. Clearing calls `onVerseChange(null)`. Callers that omit these props SHALL see no change in rendered output or behavior.

#### Scenario: Verse field renders when props provided
- **WHEN** a caller supplies `selectedVerse` and `onVerseChange` to `ArtistCharacterFilter`
- **THEN** the filter panel SHALL render a single-select Verse combobox

#### Scenario: Verse field absent when props omitted
- **WHEN** a caller does not supply `selectedVerse`/`onVerseChange`
- **THEN** the filter panel SHALL NOT render a Verse combobox, and existing behavior for that caller SHALL be unchanged

#### Scenario: Selecting a verse
- **WHEN** the user selects a verse from the Verse combobox
- **THEN** `onVerseChange` SHALL be called with the chosen verse

#### Scenario: Clearing the verse selection
- **WHEN** the user clears the verse selection
- **THEN** `onVerseChange` SHALL be called with `null`

## MODIFIED Requirements

### Requirement: ArtistCharacterFilter — Clear all
The filter panel SHALL display a "Clear all" link. Clicking it calls `onArtistChange(null)`, `onCharactersChange([])`, `onIsReferenceChange(false)` when `onIsReferenceChange` is provided, and `onVerseChange(null)` when `onVerseChange` is provided.

#### Scenario: Clear all resets all filters including reference toggle
- **WHEN** the user clicks "Clear all" while artist, characters, and the "References only" toggle are all active
- **THEN** onArtistChange(null), onCharactersChange([]), and onIsReferenceChange(false) are each called

#### Scenario: Clear all without isReference prop — only artist and characters reset
- **WHEN** the component is rendered without onIsReferenceChange and the user clicks "Clear all"
- **THEN** only onArtistChange(null) and onCharactersChange([]) are called

#### Scenario: Clear all resets verse when provided
- **WHEN** `onVerseChange` is provided and the user clicks "Clear all"
- **THEN** `onVerseChange` SHALL be called with `null` in addition to the existing artist/character/reference resets
