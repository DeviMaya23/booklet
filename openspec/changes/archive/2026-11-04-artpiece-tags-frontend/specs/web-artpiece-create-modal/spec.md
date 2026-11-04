# Spec Delta

## MODIFIED Requirements

### Requirement: Create artpiece form fields
The system SHALL render a form with title, notes, artist, characters, tags, and a "This is a reference" checkbox fields. The Tags field SHALL use the tag picker and appear directly below Characters.

#### Scenario: Title is required
- **WHEN** the user attempts to save with an empty title field
- **THEN** the Save button is disabled

#### Scenario: Notes, artist, and characters are optional
- **WHEN** the user saves without filling in notes, artist, or characters
- **THEN** the artpiece is created successfully with those fields null or empty

#### Scenario: Tags are optional
- **WHEN** the user saves without selecting any tags
- **THEN** the artpiece is created successfully with no tags

#### Scenario: Artist dropdown with add action
- **WHEN** the user clicks the "+" button next to the artist dropdown
- **THEN** the existing artist-creation modal opens; on successful creation the new artist is selected in the dropdown

#### Scenario: Character chip multi-select
- **WHEN** the user types in the characters field
- **THEN** matching characters appear as suggestions; selecting one adds it as a chip; clicking the chip's remove button deselects it

#### Scenario: Tag chip multi-select
- **WHEN** the user focuses the tags field and selects a tag from the suggestions
- **THEN** the tag is added as a chip; clicking the chip's remove button deselects it

#### Scenario: Create a tag inline
- **WHEN** the user types a name that matches no existing tag and chooses the create row
- **THEN** the tag is created, added as a chip, and remains selected for the artpiece being created

#### Scenario: Reference checkbox defaults to unchecked
- **WHEN** the modal opens
- **THEN** the "This is a reference" checkbox is unchecked by default

#### Scenario: Reference checkbox checked on save
- **WHEN** the user checks "This is a reference" and saves
- **THEN** POST /artpieces is called with is_reference: true

#### Scenario: Reference checkbox unchecked on save
- **WHEN** the user saves with "This is a reference" unchecked
- **THEN** POST /artpieces is called with is_reference: false

---

### Requirement: Save artpiece
The system SHALL create the artpiece with all current form values, the current list of uploaded file IDs, the selected tag IDs, and the `is_reference` flag in a single call.

#### Scenario: Successful save
- **WHEN** the user clicks Save with a non-empty title and no uploads in-flight
- **THEN** POST /artpieces is called with title, notes, artist_id, character_ids, tag_ids, file_ids, and is_reference; on success the modal closes and a success toast is shown

#### Scenario: Save with tags
- **WHEN** the user has selected tags `outfit` and `sketch` and clicks Save
- **THEN** POST /artpieces is called with `tag_ids` containing the ids of both tags

#### Scenario: Save with no tags
- **WHEN** the user clicks Save with no tags selected
- **THEN** POST /artpieces is called with an empty `tag_ids` array

#### Scenario: Save with no files
- **WHEN** the user clicks Save with a valid title and an empty files list
- **THEN** `POST /artpieces` is called with an empty file_ids array; the artpiece is created with no attached files

#### Scenario: Save fails
- **WHEN** the `POST /artpieces` call returns an error
- **THEN** the modal remains open, an error toast is shown, and the form state is preserved, including the selected tags

## ADDED Requirements

### Requirement: Tag suggestions and refresh in the create form
The tags field in the create artpiece form SHALL request tags with `sort=recent` and exclude already-selected tags from the suggestions. After a successful save with tags selected, the tags query SHALL be invalidated so the next picker reflects the updated recency order.

#### Scenario: Tags requested in recency order
- **WHEN** the create artpiece modal opens
- **THEN** the tags request SHALL include `sort=recent`

#### Scenario: Tags refresh after saving an artpiece
- **WHEN** the user saves an artpiece with tags selected
- **THEN** the tags query SHALL be invalidated

---

### Requirement: Selected tags count as unsaved changes in the create form
The create form SHALL treat a non-empty tag selection as a dirty form, so dismissing the modal with Escape or an outside click while tags are selected asks the user to confirm discarding.

#### Scenario: Dismiss with only tags selected
- **WHEN** the user has selected a tag, entered nothing else, and presses Escape
- **THEN** the "Discard changes?" confirmation SHALL appear instead of closing the modal

#### Scenario: Tags cleared when the modal closes
- **WHEN** the modal is closed and reopened
- **THEN** the tags field SHALL be empty
