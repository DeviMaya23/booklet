# Spec Delta

## MODIFIED Requirements

### Requirement: Create artpiece form fields
The system SHALL render a form with title, notes, artist, characters, and a "This is a reference" checkbox fields.

#### Scenario: Title is required
- **WHEN** the user attempts to save with an empty title field
- **THEN** the Save button is disabled

#### Scenario: Notes, artist, and characters are optional
- **WHEN** the user saves without filling in notes, artist, or characters
- **THEN** the artpiece is created successfully with those fields null or empty

#### Scenario: Artist dropdown with add action
- **WHEN** the user clicks the "+" button next to the artist dropdown
- **THEN** the existing artist-creation modal opens; on successful creation the new artist is selected in the dropdown

#### Scenario: Character chip multi-select
- **WHEN** the user types in the characters field
- **THEN** matching characters appear as suggestions; selecting one adds it as a chip; clicking the chip's remove button deselects it

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
The system SHALL create the artpiece with all current form values, the current list of uploaded file IDs, and the `is_reference` flag in a single call.

#### Scenario: Successful save
- **WHEN** the user clicks Save with a non-empty title and no uploads in-flight
- **THEN** POST /artpieces is called with title, notes, artist_id, character_ids, file_ids, and is_reference; on success the modal closes and a success toast is shown

#### Scenario: Save with no files
- **WHEN** the user clicks Save with a valid title and an empty files list
- **THEN** POST /artpieces is called with an empty file_ids array; the artpiece is created with no attached files

#### Scenario: Save fails
- **WHEN** the POST /artpieces call returns an error
- **THEN** the modal remains open, an error toast is shown, and the form state is preserved
