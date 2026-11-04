# web-tag-picker Specification

## Purpose

Defines the tag picker used to assign tags to an artpiece: a chip input with recency-ordered suggestions that also lets the user create a tag inline when no existing tag matches.

## Requirements

### Requirement: Selected tags shown as removable chips
The tag picker SHALL display the selected tags as chips inside the input, each with a remove control, and SHALL notify the caller with the updated list whenever a tag is added or removed. All selection state is owned by the caller.

#### Scenario: Selected tags render as chips
- **WHEN** the picker is rendered with tags `outfit` and `sketch` selected
- **THEN** a chip labelled `outfit` and a chip labelled `sketch` SHALL be shown, each with a remove control

#### Scenario: Remove a chip
- **WHEN** the user clicks the remove control on the `sketch` chip
- **THEN** the caller SHALL be notified with the list containing only `outfit`

#### Scenario: Backspace on empty input removes the last chip
- **WHEN** the input text is empty and the user presses Backspace
- **THEN** the last selected tag SHALL be removed

---

### Requirement: Tag suggestions in recency order
The picker SHALL request the user's tags with `sort=recent` and list matching suggestions in the order returned, with no client-side re-sort. Suggestions SHALL match the typed text as a case-insensitive substring and SHALL exclude tags already selected. The dropdown SHALL open when the input is focused, listing tags even before the user has typed.

#### Scenario: List requested with sort=recent
- **WHEN** the picker is rendered
- **THEN** the tags request SHALL include `sort=recent`

#### Scenario: Suggestions follow API order
- **WHEN** the API returns tags in recency order A, B, C and the user focuses the empty input
- **THEN** the suggestions SHALL be listed A, B, C

#### Scenario: Typed text filters suggestions
- **WHEN** the user types `sk` and the user's tags include `sketch`, `skirt` and `outfit`
- **THEN** the suggestions SHALL contain `sketch` and `skirt` and SHALL NOT contain `outfit`

#### Scenario: Selected tags are excluded
- **WHEN** `sketch` is already selected and the user types `sk`
- **THEN** `sketch` SHALL NOT appear in the suggestions

#### Scenario: Selecting a suggestion
- **WHEN** the user clicks a suggestion, or highlights it with the arrow keys and presses Enter
- **THEN** that tag SHALL be added as a chip and the input text SHALL be cleared

---

### Requirement: Inline tag creation
When creation is enabled and the typed text, trimmed, is not empty and does not match any existing tag name case-insensitively, the picker SHALL show a "Create" row for that text. Choosing it, or pressing Enter or comma with no suggestion highlighted, SHALL create the tag via `POST /tags` and add the returned tag as a chip. The picker SHALL not add a chip until the request has succeeded.

#### Scenario: Create row offered for an unknown name
- **WHEN** the user types `wip` and no existing tag is named `wip`
- **THEN** the dropdown SHALL show a row offering to create `wip`

#### Scenario: Create row absent for an exact existing name
- **WHEN** the user types `WIP` and a tag named `wip` exists
- **THEN** no create row SHALL be shown

#### Scenario: Creating a tag
- **WHEN** the user types `wip` and chooses the create row
- **THEN** the system SHALL call `POST /tags` with `{"name": "wip"}` and, on success, add the returned tag as a chip and clear the input

#### Scenario: Enter on an exact existing name selects it
- **WHEN** the user types `WIP`, a tag named `wip` exists, and the user presses Enter
- **THEN** the existing `wip` tag SHALL be added as a chip and `POST /tags` SHALL NOT be called

#### Scenario: Blank text creates nothing
- **WHEN** the input contains only whitespace and the user presses Enter
- **THEN** no request SHALL be made and no chip SHALL be added

#### Scenario: Newly created tag refreshes the tag list
- **WHEN** a tag is created successfully
- **THEN** the tags query SHALL be invalidated so the new tag is available to other pickers

---

### Requirement: Inline tag creation failure handling
If `POST /tags` fails, the picker SHALL NOT add a chip, SHALL keep the typed text, and SHALL show an error toast. A 409 response (the name already exists, for example created elsewhere since the list was loaded) SHALL be resolved by selecting the existing tag of that name after refetching the list, without an error toast.

#### Scenario: Server error keeps the text
- **WHEN** `POST /tags` returns a non-409 error
- **THEN** no chip SHALL be added, the typed text SHALL remain in the input, and an error toast SHALL appear

#### Scenario: Conflict selects the existing tag
- **WHEN** `POST /tags` returns 409 for the name `wip`
- **THEN** the tags list SHALL be refetched, the existing `wip` tag SHALL be added as a chip, and no error toast SHALL appear

---

### Requirement: Picker disabled state
The picker SHALL accept a disabled state in which the input cannot be focused or typed into, so no tag can be added or created.

#### Scenario: Disabled while the form is saving
- **WHEN** the picker is rendered disabled
- **THEN** the user SHALL NOT be able to type into the input or add or create a tag

---

### Requirement: Picker waits for the tag list
Until the user's tag list has loaded, the picker SHALL NOT accept typing and SHALL NOT offer to create a tag, because an existing tag cannot be told apart from one not yet loaded. If the list fails to load the picker SHALL show an error and stay unavailable for adding, while still showing already-selected tags.

#### Scenario: Loading
- **WHEN** the tag list request is still in flight
- **THEN** the input SHALL be disabled with a loading placeholder, no create row SHALL be offered, and no `POST /tags` SHALL be sent

#### Scenario: Ready once loaded
- **WHEN** the tag list request succeeds
- **THEN** the input SHALL become enabled with its normal placeholder

#### Scenario: Load failure
- **WHEN** the tag list request fails
- **THEN** the picker SHALL show "Failed to load tags.", the input SHALL be disabled, and any selected tags SHALL remain visible

### Requirement: Focus suggestion list is capped
When the dropdown is open with empty input text, the picker SHALL list at most the first 10 unselected tags in the order returned. If more unselected tags exist, it SHALL show a non-selectable hint row indicating that typing searches the rest. Once the user has typed text, matching suggestions SHALL NOT be capped.

#### Scenario: Long list is truncated on focus
- **WHEN** the user has 25 tags, none selected, and focuses the empty input
- **THEN** the first 10 tags in API order SHALL be listed, followed by a hint row indicating more tags can be found by typing

#### Scenario: Short list has no hint
- **WHEN** the user has 10 or fewer unselected tags and focuses the empty input
- **THEN** all of them SHALL be listed and no hint row SHALL be shown

#### Scenario: Selected tags do not use up the cap
- **WHEN** the user has 25 tags, 3 of them selected, and focuses the empty input
- **THEN** 10 unselected tags SHALL be listed

#### Scenario: Typing searches beyond the cap
- **WHEN** the user types text that matches a tag that was not among the first 10
- **THEN** that tag SHALL be listed

#### Scenario: Hint row is not selectable
- **WHEN** the hint row is shown and the user navigates with the arrow keys or clicks the hint row
- **THEN** the hint row SHALL NOT be highlighted or selected and no tag SHALL be added
