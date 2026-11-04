# Spec Delta

## MODIFIED Requirements

### Requirement: Filter popover
The gallery SHALL display a Filter button beside the search input that opens a floating panel containing artist, character, tag, and verse filters, and a "References only" toggle. The Filter button and artist/character/tag/verse controls are provided by the shared `ArtistCharacterFilter` component.

#### Scenario: Filter button opens popover
- **WHEN** the user clicks the Filter button
- **THEN** a popover panel SHALL appear containing an artist section, a characters section, a tags section, a verse section, and a "References only" toggle

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

#### Scenario: Filter button shows active filter summary — verse
- **WHEN** a verse is active, regardless of other filters
- **THEN** the Filter button summary SHALL include it (e.g. "1 verse", "1 artist · 1 verse")

#### Scenario: Filter button shows active filter summary — tags
- **WHEN** at least one tag is active, regardless of other filters
- **THEN** the Filter button summary SHALL include it (e.g. "1 tag, all", "2 tags, any", "1 artist · 2 tags, any")

---

### Requirement: Save to dashboard action in filter popover
When at least one filter is active, the gallery filter popover SHALL display a "Save to dashboard" button in the footer alongside "Clear all". Clicking it opens the Save Filter modal with the current filter state.

#### Scenario: Button appears when filters are active
- **WHEN** at least one filter is active (artist selected, at least one character chip present, or at least one tag chip present)
- **THEN** the filter popover footer SHALL show a "Save to dashboard" button

#### Scenario: Button absent when no filters are active
- **WHEN** no artist is selected, no characters are selected and no tags are selected
- **THEN** the "Save to dashboard" button SHALL NOT appear in the footer

#### Scenario: Clicking the button opens the Save Filter modal
- **WHEN** the user clicks "Save to dashboard"
- **THEN** the Save Filter modal SHALL open pre-loaded with the current filter state

---

### Requirement: Saved-as confirmation in filter popover
After a successful save from the Save Filter modal, the filter popover footer SHALL replace the "Save to dashboard" button with a "Saved as '<name>'" confirmation for as long as the filter values remain unchanged.

#### Scenario: Confirmation shown after successful save
- **WHEN** the user saves a filter and the modal closes with a 201 response
- **THEN** the filter popover footer SHALL display "Saved as '<name>'" in place of the button

#### Scenario: Confirmation clears when filter changes
- **WHEN** the filter values change after a successful save (artist, characters, tags or either match mode modified)
- **THEN** the "Saved as '<name>'" confirmation SHALL be removed and the "Save to dashboard" button SHALL reappear

## ADDED Requirements

### Requirement: Tag filter
The filter panel SHALL contain a tag multi-select (chip-based, selecting existing tags only, with no inline creation) with an All/Any segmented toggle. The toggle SHALL be visible only when two or more tags are selected, and SHALL default to "All". Suggestions SHALL be requested with `sort=recent`, listed in the order returned, and SHALL show tag names without counts.

#### Scenario: Add a tag chip
- **WHEN** the user selects a tag in the tag input
- **THEN** a chip for that tag SHALL appear and the gallery SHALL be filtered accordingly

#### Scenario: Remove a tag chip
- **WHEN** the user clicks the remove control on a tag chip
- **THEN** that tag is removed from the filter

#### Scenario: Match All
- **WHEN** the tag toggle is "All" and two or more tag chips are present
- **THEN** the gallery SHALL show only artpieces whose `tags` array contains ALL of the selected tag IDs

#### Scenario: Match Any
- **WHEN** the tag toggle is "Any" and two or more tag chips are present
- **THEN** the gallery SHALL show only artpieces whose `tags` array contains AT LEAST ONE of the selected tag IDs

#### Scenario: Single tag
- **WHEN** exactly one tag chip is present
- **THEN** the gallery SHALL show only artpieces that have that tag, and the toggle SHALL NOT be rendered

#### Scenario: Toggle defaults to All
- **WHEN** the filter panel is first opened with no prior tag selection
- **THEN** the tag match mode SHALL be "All"

#### Scenario: No tag chips — no tag filtering
- **WHEN** no tag chips are present
- **THEN** no tag filter is applied to the gallery results

#### Scenario: Unknown names cannot be created from the filter
- **WHEN** the user types a name that matches no tag in the tag filter input
- **THEN** no create row SHALL be offered and no tag SHALL be created

#### Scenario: Tag suggestions request recency order
- **WHEN** the filter panel renders the Tags field
- **THEN** the tags request SHALL include `sort=recent` and suggestions SHALL exclude already-selected tags

---

### Requirement: Tag filter combines with other filters
The tag filter SHALL be AND-combined with the artist, character, verse, references, and search filters, and SHALL be independent of the character match mode.

#### Scenario: Tags and artist
- **WHEN** the user selects an artist and one or more tags
- **THEN** the gallery SHALL show only artpieces that satisfy BOTH the artist filter AND the tag filter

#### Scenario: Tags and characters with different match modes
- **WHEN** characters are matched with "Any" and tags are matched with "All"
- **THEN** the character filter SHALL use "Any" semantics and the tag filter SHALL use "All" semantics

---

### Requirement: Tag filter state persisted in URL and localStorage
The gallery SHALL sync the selected tag ids and the tag match mode to the URL search params `tags` (comma-separated ids) and `tag_match` (present only when "Any"), and SHALL include them in the structural filter state persisted under `artpieces.filterState`. Clearing the tags SHALL remove both params.

#### Scenario: Selected tags reflected in the URL
- **WHEN** the user selects tags in the filter panel
- **THEN** the URL SHALL contain a `tags` param with their ids, and a `tag_match=any` param only when the mode is "Any"

#### Scenario: Filter state restored from URL
- **WHEN** the user opens the gallery with `tags` and `tag_match` params in the URL
- **THEN** the gallery SHALL initialise with those tags selected and that match mode

#### Scenario: Filter restored from localStorage
- **WHEN** the user navigates away and returns to `/app/artpieces` with no filter params and `localStorage` holds a tag selection
- **THEN** the gallery SHALL restore the stored tags and match mode and apply them

#### Scenario: A stored tag-only filter is restorable
- **WHEN** the only stored filter is a non-empty tag selection
- **THEN** it SHALL be restored (it is not treated as an empty filter)

#### Scenario: Stored state from before tags existed
- **WHEN** `localStorage` holds a filter state with no tag fields
- **THEN** the gallery SHALL restore it with no tags selected and tag match mode "All"

#### Scenario: Clearing removes the params
- **WHEN** the user removes the last tag chip or clicks "Clear all", including while the tag match mode is "Any"
- **THEN** the `tags` and `tag_match` params SHALL be removed from the URL, and `localStorage` SHALL store no tags and the default match mode "All"

---

### Requirement: Filter ids resolve against current data
Selected tag ids and character ids carried by the URL, `localStorage` or a saved filter SHALL be resolved against the currently fetched tag and character lists. An id that does not resolve (for example because the tag or character was deleted) SHALL be dropped from the active filter: it SHALL NOT raise an error, SHALL NOT be shown as a chip, and SHALL NOT make the filter match nothing. The tag filter SHALL NOT be applied until the tag list has loaded.

#### Scenario: Deleted tag in All mode
- **WHEN** the URL carries `tags=<A>,<B>` with `tag_match` absent (All), tag B no longer exists, and artpieces with tag A exist
- **THEN** the gallery SHALL show artpieces that have tag A, as if only tag A were selected

#### Scenario: Every selected tag deleted
- **WHEN** every tag id carried by the URL or `localStorage` no longer exists
- **THEN** the gallery SHALL behave as if no tag filter is active, without error

#### Scenario: Deleted character in All mode
- **WHEN** the URL carries `characters=<A>,<B>`, character B no longer exists, and artpieces with character A exist
- **THEN** the gallery SHALL show artpieces that have character A, as if only character A were selected

#### Scenario: Tag list still loading
- **WHEN** the gallery opens with a `tags` param and the tag list has not loaded yet
- **THEN** the tag filter SHALL NOT be applied until the list loads, and SHALL NOT exclude artpieces in the meantime

#### Scenario: Tag list failed to load
- **WHEN** the gallery opens with a `tags` param and the tag list request fails
- **THEN** the tag ids in the URL SHALL be applied as given (none can be resolved or dropped), so the gallery does not show every artpiece under a tagged filter
