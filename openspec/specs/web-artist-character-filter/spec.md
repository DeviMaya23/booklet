## Purpose

Defines the shared `ArtistCharacterFilter` component — a controlled trigger button + floating panel that provides artist and character filtering in any context that needs it (gallery, modals). All filter state is owned by the caller; the component owns only the open/closed state of the panel internally.

## Requirements

### Requirement: ArtistCharacterFilter shared component
The system SHALL provide a shared controlled component `ArtistCharacterFilter` at `components/ArtistCharacterFilter.tsx`. The component renders a trigger button that opens a floating panel. The panel contains: a single-select Artist combobox, a multi-select Characters token input, an optional All/Any match toggle, and a "Clear all" link. All value state is owned by the caller; the component owns only the open/closed state of the panel internally.

#### Scenario: Trigger button — no active filters
- **WHEN** no artist is selected and no characters are selected
- **THEN** the trigger button SHALL display the label "Filter"

#### Scenario: Trigger button — artist active
- **WHEN** an artist is selected and no characters are selected
- **THEN** the trigger button SHALL display "1 artist"

#### Scenario: Trigger button — characters active, no toggle
- **WHEN** characters are selected and `characterMatch` is not provided
- **THEN** the trigger button SHALL display "N character(s)" with no match mode suffix

#### Scenario: Trigger button — both active with match mode
- **WHEN** an artist is selected and characters are selected and `characterMatch` is provided
- **THEN** the trigger button SHALL display "1 artist · N characters, <match>" where match is "any" or "all"

#### Scenario: Panel opens on trigger click
- **WHEN** the user clicks the trigger button
- **THEN** the filter panel SHALL open showing the Artist combobox and Characters token input

#### Scenario: Panel closes on outside click or Escape
- **WHEN** the panel is open and the user clicks outside it or presses Escape
- **THEN** the panel SHALL close without resetting filter values

---

### Requirement: ArtistCharacterFilter — Artist field
The filter panel SHALL contain a single-select Artist combobox. Selecting an artist calls `onArtistChange` with the chosen artist. Clearing calls `onArtistChange(null)`.

#### Scenario: Select artist
- **WHEN** the user selects an artist from the combobox
- **THEN** `onArtistChange` is called with that artist

#### Scenario: Clear artist
- **WHEN** the user clears the artist selection
- **THEN** `onArtistChange(null)` is called

---

### Requirement: ArtistCharacterFilter — Characters field
The filter panel SHALL contain a multi-select Characters token input. Adding a character calls `onCharactersChange` with the updated list. Removing calls `onCharactersChange` with the character removed.

#### Scenario: Add character
- **WHEN** the user selects a character in the token input
- **THEN** `onCharactersChange` is called with the character appended

#### Scenario: Remove character chip
- **WHEN** the user removes a character chip
- **THEN** `onCharactersChange` is called with that character removed

---

### Requirement: ArtistCharacterFilter — All/Any match toggle
The filter panel SHALL display an All/Any segmented toggle only when `characterMatch` and `onCharacterMatchChange` are provided by the caller AND two or more characters are currently selected. When fewer than two characters are selected, the toggle is hidden even if the props are present.

#### Scenario: Toggle hidden with fewer than two characters
- **WHEN** `characterMatch` is provided and one or zero characters are selected
- **THEN** the All/Any toggle SHALL NOT be rendered

#### Scenario: Toggle visible with two or more characters
- **WHEN** `characterMatch` is provided and two or more characters are selected
- **THEN** the All/Any toggle SHALL be rendered showing the current match mode

#### Scenario: Toggle hidden when prop omitted
- **WHEN** `characterMatch` is not provided
- **THEN** the All/Any toggle SHALL NOT be rendered regardless of how many characters are selected

#### Scenario: Toggle interaction
- **WHEN** the user clicks "Any" or "All"
- **THEN** `onCharacterMatchChange` is called with the selected mode

---

### Requirement: ArtistCharacterFilter — Clear all
The filter panel SHALL display a "Clear all" link. Clicking it calls both `onArtistChange(null)` and `onCharactersChange([])`.

#### Scenario: Clear all resets both fields
- **WHEN** the user clicks "Clear all"
- **THEN** `onArtistChange(null)` and `onCharactersChange([])` are both called
