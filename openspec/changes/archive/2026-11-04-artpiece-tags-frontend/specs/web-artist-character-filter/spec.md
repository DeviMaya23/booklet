# Spec Delta

## MODIFIED Requirements

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

## ADDED Requirements

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
