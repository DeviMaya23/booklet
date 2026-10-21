# Spec Delta

## MODIFIED Requirements

### Requirement: ArtistCharacterFilter shared component
The system SHALL provide a shared controlled component `ArtistCharacterFilter` at `components/ArtistCharacterFilter.tsx`. The component renders a trigger button that opens a floating panel. The panel contains: a single-select Artist combobox, a multi-select Characters token input, an optional All/Any match toggle, a "Clear all" link, and an optional caller-supplied footer action rendered alongside "Clear all". All value state is owned by the caller; the component owns only the open/closed state of the panel internally.

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

#### Scenario: Footer action rendered when provided
- **WHEN** the caller passes a `footerAction` prop
- **THEN** the panel footer SHALL render it alongside the "Clear all" link

#### Scenario: Footer action absent when not provided
- **WHEN** the caller omits `footerAction`
- **THEN** the panel footer SHALL show only the "Clear all" link with no extra slot content
