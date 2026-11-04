## Purpose

Defines the shared `ArtistCharacterFilter` component — a controlled trigger button + floating panel that provides artist and character filtering in any context that needs it (gallery, modals). All filter state is owned by the caller; the component owns only the open/closed state of the panel internally.

## Requirements

### Requirement: ArtistCharacterFilter shared component
The system SHALL provide a shared controlled component `ArtistCharacterFilter` at `components/ArtistCharacterFilter.tsx`. The component renders a trigger button that opens a floating panel. The panel contains: a single-select Artist combobox, a multi-select Characters token input, an optional All/Any match toggle, an optional "References only" toggle, a "Clear all" link, and an optional caller-supplied footer action rendered alongside "Clear all". All value state is owned by the caller; the component owns only the open/closed state of the panel internally.

#### Scenario: Renders with all optional props absent
- **WHEN** the component is rendered without characterMatch, onCharacterMatchChange, isReference, or onIsReferenceChange
- **THEN** the All/Any toggle and the References only toggle are not rendered; the component functions as an artist-and-character-only filter

#### Scenario: Renders References only toggle when props are provided
- **WHEN** the component is rendered with isReference and onIsReferenceChange props
- **THEN** a "References only" toggle is rendered in the panel

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

---

### Requirement: ArtistCharacterFilter — Clear all
The filter panel SHALL display a "Clear all" link. Clicking it calls `onArtistChange(null)`, `onCharactersChange([])`, `onIsReferenceChange(false)` when `onIsReferenceChange` is provided, `onVerseChange(null)` when `onVerseChange` is provided, and `onTagsChange([])` when `onTagsChange` is provided.

#### Scenario: Clear all resets all filters including reference toggle
- **WHEN** the user clicks "Clear all" while artist, characters, and the "References only" toggle are all active
- **THEN** onArtistChange(null), onCharactersChange([]), and onIsReferenceChange(false) are each called

#### Scenario: Clear all without isReference prop — only artist and characters reset
- **WHEN** the component is rendered without onIsReferenceChange and the user clicks "Clear all"
- **THEN** only onArtistChange(null) and onCharactersChange([]) are called

#### Scenario: Clear all resets verse when provided
- **WHEN** `onVerseChange` is provided and the user clicks "Clear all"
- **THEN** `onVerseChange` SHALL be called with `null` in addition to the existing artist/character/reference resets

#### Scenario: Clear all resets tags when provided
- **WHEN** `onTagsChange` is provided and the user clicks "Clear all"
- **THEN** `onTagsChange` SHALL be called with `[]` in addition to the existing resets

#### Scenario: Clear all link shown for a tag-only filter
- **WHEN** only tags are selected and no artist, characters, verse or references filter is active
- **THEN** the "Clear all" link SHALL be shown

---

### Requirement: ArtistCharacterFilter — References only toggle
When `isReference` and `onIsReferenceChange` props are provided, the filter panel SHALL render a "References only" toggle. Toggling it on calls `onIsReferenceChange(true)`; toggling it off calls `onIsReferenceChange(false)`.

#### Scenario: Toggle on
- **WHEN** the user activates the "References only" toggle
- **THEN** onIsReferenceChange(true) is called

#### Scenario: Toggle off
- **WHEN** the user deactivates the "References only" toggle
- **THEN** onIsReferenceChange(false) is called

---

### Requirement: ArtistCharacterFilter — summary label includes references
When the trigger button computes its summary label, an active `isReference` (true) SHALL contribute "references" to the label.

#### Scenario: Summary with references only active
- **WHEN** isReference is true and no artist or characters are selected
- **THEN** the trigger button label SHALL be "references" (not the plain "Filter" label)

#### Scenario: Summary with artist and references active
- **WHEN** isReference is true and one artist is selected
- **THEN** the trigger button label SHALL include both (e.g. "1 artist · references")

### Requirement: ArtistCharacterFilter — character suggestions in recency order
The Characters token input SHALL request the character list with `sort=recent` and display suggestions in the order returned, with no client-side re-sort, in every context the component is used (gallery and commission modal).

#### Scenario: Suggestions follow API recency order
- **WHEN** the filter panel is open and the API returns characters in recency order A, B, C
- **THEN** the suggestions SHALL be listed A, B, C, excluding already-selected characters

#### Scenario: List requested with sort=recent
- **WHEN** the filter panel renders the Characters field
- **THEN** the characters request SHALL include `sort=recent`

---

### Requirement: ArtistCharacterFilter — Tags field
The filter panel SHALL optionally contain a multi-select Tags token input, rendered only when `tags` and `onTagsChange` props are both provided by the caller, selecting existing tags only. Adding a tag calls `onTagsChange` with the updated list; removing calls it with the tag removed. Callers that omit these props SHALL see no change in rendered output or behavior.

#### Scenario: Tags field renders when props provided
- **WHEN** a caller supplies `tags` and `onTagsChange`
- **THEN** the filter panel SHALL render a Tags token input

#### Scenario: Tags field absent when props omitted
- **WHEN** a caller does not supply `tags`/`onTagsChange` (as in the commission modal and the add-to-existing-artpiece modal)
- **THEN** the panel SHALL NOT render a Tags field and the caller's existing behavior SHALL be unchanged

#### Scenario: Add tag
- **WHEN** the user selects a tag in the Tags input
- **THEN** `onTagsChange` is called with the tag appended

#### Scenario: Remove tag chip
- **WHEN** the user removes a tag chip
- **THEN** `onTagsChange` is called with that tag removed

#### Scenario: Suggestions in recency order
- **WHEN** the panel is open and the API returns tags in recency order A, B, C
- **THEN** the suggestions SHALL be listed A, B, C, excluding already-selected tags, and the tags request SHALL include `sort=recent`

---

### Requirement: ArtistCharacterFilter — Tag All/Any match toggle
The filter panel SHALL display a tag All/Any segmented toggle only when `tagMatch` and `onTagMatchChange` are provided AND two or more tags are selected. It is independent of the character toggle.

#### Scenario: Toggle hidden with fewer than two tags
- **WHEN** `tagMatch` is provided and one or zero tags are selected
- **THEN** the tag All/Any toggle SHALL NOT be rendered

#### Scenario: Toggle visible with two or more tags
- **WHEN** `tagMatch` is provided and two or more tags are selected
- **THEN** the tag All/Any toggle SHALL be rendered showing the current tag match mode

#### Scenario: Toggle interaction
- **WHEN** the user clicks "Any" or "All" in the Tags section
- **THEN** `onTagMatchChange` is called with the selected mode and the character match mode is unaffected

---

### Requirement: ArtistCharacterFilter — summary label includes tags
When the trigger button computes its summary label, selected tags SHALL contribute "N tag(s)" to the label, followed by the match mode when `tagMatch` is provided, in the same form as characters.

#### Scenario: Tags only, one tag, no match mode
- **WHEN** one tag is selected, `tagMatch` is not provided, and nothing else is active
- **THEN** the trigger button label SHALL be "1 tag"

#### Scenario: One tag with match mode
- **WHEN** one tag is selected, `tagMatch` is "all", and nothing else is active
- **THEN** the trigger button label SHALL be "1 tag, all"

#### Scenario: Tags with match mode
- **WHEN** two tags are selected, `tagMatch` is "any", and nothing else is active
- **THEN** the trigger button label SHALL be "2 tags, any"

#### Scenario: Tags combined with other filters
- **WHEN** one artist and two characters (match "all") and one tag are selected
- **THEN** the label SHALL include each part separated by " · "
