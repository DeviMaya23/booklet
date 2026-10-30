# Spec Delta

## MODIFIED Requirements

### Requirement: Filter popover
The gallery SHALL display a Filter button beside the search input that opens a floating panel containing artist and character filters, and a "References only" toggle. The Filter button and artist/character controls are provided by the shared `ArtistCharacterFilter` component.

#### Scenario: Filter button opens popover
- **WHEN** the user clicks the Filter button
- **THEN** a popover panel SHALL appear containing an artist section, a characters section, and a "References only" toggle

#### Scenario: Filter button closes popover
- **WHEN** the popover is open and the user clicks outside the popover or presses Escape
- **THEN** the popover SHALL close

#### Scenario: Filter button shows active filter summary — artist or characters
- **WHEN** at least one of artist or characters is active
- **THEN** the Filter button SHALL display a short summary (e.g. "1 artist", "2 characters, all", "1 artist · 2 characters, any") instead of the plain "Filter" label

#### Scenario: Filter button shows active filter summary — references only
- **WHEN** the "References only" toggle is on and no artist or characters are active
- **THEN** the Filter button SHALL include "references" in its summary label

#### Scenario: Filter button shows active filter summary — combined
- **WHEN** artist or characters are active AND the "References only" toggle is on
- **THEN** the Filter button summary SHALL reflect both (e.g. "1 artist · references")

---

## ADDED Requirements

### Requirement: References only filter toggle
The gallery filter popover SHALL contain a "References only" toggle (default off). When off the gallery shows only artpieces where `is_reference` is false. When on the gallery shows only artpieces where `is_reference` is true. All other filters (artist, characters, match mode, sort, search) apply in both states.

#### Scenario: Toggle off — only non-references shown
- **WHEN** the "References only" toggle is off (default)
- **THEN** the gallery SHALL display only artpieces where is_reference is false

#### Scenario: Toggle on — only references shown
- **WHEN** the user turns the "References only" toggle on
- **THEN** the gallery SHALL display only artpieces where is_reference is true

#### Scenario: Other filters still apply with toggle on
- **WHEN** the "References only" toggle is on and an artist filter is also active
- **THEN** the gallery SHALL show only artpieces that are both references and match the artist filter

#### Scenario: Toggle state persists in URL
- **WHEN** the user turns the "References only" toggle on and navigates away then back
- **THEN** the URL SHALL contain the reference param and the gallery SHALL restore the toggle to on

#### Scenario: Clear all resets toggle
- **WHEN** the user clicks "Clear all" in the filter popover while the "References only" toggle is on
- **THEN** the toggle SHALL be reset to off

### Requirement: Gallery default shows only non-references
When the gallery loads with no filter params in the URL, the "References only" toggle SHALL default to off, showing only non-reference artpieces.

#### Scenario: Default state on first load
- **WHEN** the gallery loads with no filter params in the URL
- **THEN** the "References only" toggle is off and only artpieces with is_reference false are shown

#### Scenario: Absence of ref URL param treated as off
- **WHEN** the URL does not contain the ref param
- **THEN** the gallery behaves as if the toggle is off (non-references only)
