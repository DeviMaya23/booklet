# Spec Delta

## MODIFIED Requirements

### Requirement: ArtistCharacterFilter shared component
The system SHALL provide a shared controlled component `ArtistCharacterFilter` at `components/ArtistCharacterFilter.tsx`. The component renders a trigger button that opens a floating panel. The panel contains: a single-select Artist combobox, a multi-select Characters token input, an optional All/Any match toggle, an optional "References only" toggle, a "Clear all" link, and an optional caller-supplied footer action rendered alongside "Clear all". All value state is owned by the caller; the component owns only the open/closed state of the panel internally.

#### Scenario: Renders with all optional props absent
- **WHEN** the component is rendered without characterMatch, onCharacterMatchChange, isReference, or onIsReferenceChange
- **THEN** the All/Any toggle and the References only toggle are not rendered; the component functions as an artist-and-character-only filter

#### Scenario: Renders References only toggle when props are provided
- **WHEN** the component is rendered with isReference and onIsReferenceChange props
- **THEN** a "References only" toggle is rendered in the panel

---

### Requirement: ArtistCharacterFilter — Clear all
The filter panel SHALL display a "Clear all" link. Clicking it calls `onArtistChange(null)`, `onCharactersChange([])`, and `onIsReferenceChange(false)` when `onIsReferenceChange` is provided.

#### Scenario: Clear all resets all filters including reference toggle
- **WHEN** the user clicks "Clear all" while artist, characters, and the "References only" toggle are all active
- **THEN** onArtistChange(null), onCharactersChange([]), and onIsReferenceChange(false) are each called

#### Scenario: Clear all without isReference prop — only artist and characters reset
- **WHEN** the component is rendered without onIsReferenceChange and the user clicks "Clear all"
- **THEN** only onArtistChange(null) and onCharactersChange([]) are called

---

## ADDED Requirements

### Requirement: ArtistCharacterFilter — References only toggle
When `isReference` and `onIsReferenceChange` props are provided, the filter panel SHALL render a "References only" toggle. Toggling it on calls `onIsReferenceChange(true)`; toggling it off calls `onIsReferenceChange(false)`.

#### Scenario: Toggle on
- **WHEN** the user activates the "References only" toggle
- **THEN** onIsReferenceChange(true) is called

#### Scenario: Toggle off
- **WHEN** the user deactivates the "References only" toggle
- **THEN** onIsReferenceChange(false) is called

### Requirement: ArtistCharacterFilter — summary label includes references
When the trigger button computes its summary label, an active `isReference` (true) SHALL contribute "references" to the label.

#### Scenario: Summary with references only active
- **WHEN** isReference is true and no artist or characters are selected
- **THEN** the trigger button label SHALL be "references" (not the plain "Filter" label)

#### Scenario: Summary with artist and references active
- **WHEN** isReference is true and one artist is selected
- **THEN** the trigger button label SHALL include both (e.g. "1 artist · references")
