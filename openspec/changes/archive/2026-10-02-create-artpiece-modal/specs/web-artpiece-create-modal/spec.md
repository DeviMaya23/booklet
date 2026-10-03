## ADDED Requirements

### Requirement: Open modal from inbox entry points
The system SHALL open the "Create New Artpiece" modal when the user triggers "New Artpiece" from either the selection toolbar dropdown or the right-click context menu, passing the currently selected file IDs as the pre-populated file list.

#### Scenario: Open from selection toolbar
- **WHEN** the user selects one or more files and clicks "New Artpiece" in the "Add to artpiece ▾" dropdown
- **THEN** the "Create New Artpiece" modal opens with the selected files pre-populated in the files list

#### Scenario: Open from right-click context menu
- **WHEN** the user right-clicks one or more files and clicks "New Artpiece" in the context menu
- **THEN** the "Create New Artpiece" modal opens with the right-clicked selection pre-populated in the files list

#### Scenario: Open with no pre-selected files (empty state)
- **WHEN** the modal is opened without any pre-selected file IDs
- **THEN** the files list is empty and only the upload drop zone is shown

---

### Requirement: Create artpiece form fields
The system SHALL render a form with title, notes, artist, and characters fields.

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

---

### Requirement: Files sub-component — pre-populated files
The system SHALL populate the files list with the pre-selected file IDs when the modal opens, using name and notes from the existing file record.

#### Scenario: Pre-populated file with thumbnail ready
- **WHEN** the modal opens with a pre-selected file whose `thumbnail_url` is non-null
- **THEN** the file row shows the thumbnail, the existing name (if any), the existing notes (if any), and a remove button

#### Scenario: Pre-populated file with thumbnail pending
- **WHEN** the modal opens with a pre-selected file whose `thumbnail_url` is null
- **THEN** the file row shows a spinner in the thumbnail slot and the name/notes fields are disabled until the thumbnail becomes available

---

### Requirement: Files sub-component — upload drop zone
The system SHALL render a persistent upload affordance pinned at the top of the scrollable files region, always visible regardless of scroll position.

#### Scenario: Drop zone always visible
- **WHEN** the files list contains many rows and the user has scrolled down
- **THEN** the dashed-border drop zone remains visible at the top of the scrollable region

#### Scenario: Drop files onto drop zone
- **WHEN** the user drops one or more files onto the drop zone
- **THEN** each file immediately gets a row inserted with a spinner in the thumbnail slot and disabled name/notes fields, and the upload flow (init → PUT → complete) starts for each file in parallel

#### Scenario: Click drop zone to open file picker
- **WHEN** the user clicks the drop zone
- **THEN** the native file picker opens; selecting files triggers the same upload flow as dropping

---

### Requirement: Files sub-component — upload lifecycle
The system SHALL manage the per-file upload flow within the modal, including thumbnail readiness and error handling.

#### Scenario: Upload completes — thumbnail pending
- **WHEN** the complete-upload call succeeds but `thumbnail_url` is still null
- **THEN** the row remains with a spinner in the thumbnail slot and disabled name/notes fields until the thumbnail becomes available

#### Scenario: Upload completes — thumbnail ready
- **WHEN** the thumbnail poll detects a non-null `thumbnail_url` for the file
- **THEN** the spinner is replaced with the thumbnail image and the name/notes fields become editable

#### Scenario: Upload fails
- **WHEN** any step of the upload flow (init, PUT, or complete) throws an error
- **THEN** the file's row is removed from the list entirely and an error toast is shown with the message "Failed to upload <filename>", where filename is the browser File object's name at the time of selection

#### Scenario: Save disabled while upload in-flight
- **WHEN** at least one file in the list has an upload in progress
- **THEN** the Save button is disabled

#### Scenario: One failed upload does not affect others
- **WHEN** one file in a batch fails to upload
- **THEN** the other files in the batch continue uploading independently

---

### Requirement: Files sub-component — removing a file
The system SHALL allow the user to remove a file from the draft list without a backend call.

#### Scenario: Remove a file from the list
- **WHEN** the user clicks the "×" button on a file row
- **THEN** the row is removed from the draft list immediately; no backend call is made; the file remains in the dump with `artpiece_id` null

---

### Requirement: Files sub-component — file name and notes editing
The system SHALL allow the user to edit a file's name and notes inline, saving each change on blur.

#### Scenario: Edit file name and blur
- **WHEN** the user edits the name field of a file row and moves focus away
- **THEN** `PUT /files/:id` is called with the updated name; the field shows the saved value

#### Scenario: Edit file notes and blur
- **WHEN** the user edits the notes field of a file row and moves focus away
- **THEN** `PUT /files/:id` is called with the updated notes; the field shows the saved value

#### Scenario: Name and notes fields disabled during thumbnail pending
- **WHEN** a file row is in the thumbnail-pending state (thumbnail_url is null)
- **THEN** the name and notes fields are disabled and cannot be edited

---

### Requirement: Save artpiece
The system SHALL create the artpiece with all current form values and the current list of uploaded file IDs in a single call.

#### Scenario: Successful save
- **WHEN** the user clicks Save with a non-empty title and no uploads in-flight
- **THEN** `POST /artpieces` is called with title, notes, artist_id, character_ids, and the file_ids of all files currently in the list; on success the modal closes and a success toast is shown

#### Scenario: Save with no files
- **WHEN** the user clicks Save with a valid title and an empty files list
- **THEN** `POST /artpieces` is called with an empty file_ids array; the artpiece is created with no attached files

#### Scenario: Save fails
- **WHEN** the `POST /artpieces` call returns an error
- **THEN** the modal remains open, an error toast is shown, and the form state is preserved

---

### Requirement: Close modal while uploads in-flight
The system SHALL prevent the user from closing the modal while any file is uploading.

#### Scenario: Escape or close button while uploading
- **WHEN** the user presses Escape or clicks the dialog close button while at least one file upload is in-flight
- **THEN** the close action is ignored; the modal remains open until all in-flight uploads resolve
