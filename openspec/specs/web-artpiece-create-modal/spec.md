## Purpose

Defines the "Create New Artpiece" modal — a form for creating an artpiece with a title, optional metadata fields, and a files sub-component that supports pre-populated files from the inbox, drag-and-drop upload, per-file inline editing, and upload lifecycle management.

## Requirements

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

### Requirement: Files sub-component — pre-populated files
The system SHALL populate the files list with the pre-selected file IDs when the modal opens, using name and notes from the existing file record.

#### Scenario: Pre-populated file with thumbnail ready
- **WHEN** the modal opens with a pre-selected file whose `thumbnail_url` is non-null
- **THEN** the file row shows the thumbnail, the existing name (if any), the existing notes (if any), and a remove button

#### Scenario: Pre-populated file with thumbnail pending
- **WHEN** the modal opens with a pre-selected file whose `thumbnail_gen_state === "pending"`
- **THEN** the file row shows a spinner in the thumbnail slot and the name/notes fields are disabled

#### Scenario: Pre-populated file with thumbnail settled (no thumbnail)
- **WHEN** the modal opens with a pre-selected file whose `thumbnail_url === null` and `thumbnail_gen_state !== "pending"`
- **THEN** the file row shows a fallback icon in the thumbnail slot and the name/notes fields are enabled

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
The system SHALL manage the per-file upload flow within the modal, including thumbnail readiness and error handling. The file extension SHALL be stripped from `File.name` before it is sent as `name` in the initiate-upload request and before the name input is pre-filled.

#### Scenario: Upload completes — thumbnail pending
- **WHEN** the complete-upload call succeeds and `thumbnail_gen_state === "pending"`
- **THEN** the row remains with a spinner in the thumbnail slot and disabled name/notes fields

#### Scenario: Upload completes — thumbnail done
- **WHEN** `thumbnail_gen_state` becomes `"done"` for a file
- **THEN** the spinner is replaced with the thumbnail image and the name/notes fields become editable; the name input is pre-filled with the file's name with extension stripped (e.g. `pic.jpg` → `pic`)

#### Scenario: Upload completes — thumbnail failed or not applicable
- **WHEN** `thumbnail_gen_state` becomes `"failed"` or `"not_applicable"` for a file
- **THEN** the spinner is replaced with a fallback icon and the name/notes fields become editable; the name input is pre-filled with the file's name with extension stripped

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
The system SHALL allow the user to edit a file's name and notes inline, saving each change on blur. Name inputs SHALL be validated against the shared file name validation rules before the save call is made.

#### Scenario: Edit file name and blur — valid value
- **WHEN** the user edits the name field of a file row with a valid value and moves focus away
- **THEN** `PUT /files/:id` is called with the updated name; the field shows the saved value

#### Scenario: Edit file name and blur — invalid value
- **WHEN** the user edits the name field of a file row with a forbidden character or blank value and moves focus away
- **THEN** the name input shows a red border and an inline error message; `PUT /files/:id` is not called

#### Scenario: Edit file notes and blur
- **WHEN** the user edits the notes field of a file row and moves focus away
- **THEN** `PUT /files/:id` is called with the updated notes; the field shows the saved value

#### Scenario: Name and notes fields disabled during thumbnail pending
- **WHEN** a file row has `thumbnail_gen_state === "pending"`
- **THEN** the name and notes fields are disabled and cannot be edited

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

---

### Requirement: Close modal while uploads in-flight
The system SHALL prevent the user from closing the modal while any file is uploading.

#### Scenario: Escape or close button while uploading
- **WHEN** the user presses Escape or clicks the dialog close button while at least one file upload is in-flight
- **THEN** the close action is ignored; the modal remains open until all in-flight uploads resolve

---

### Requirement: File size validation before upload
The system SHALL validate each selected file's size against the per-category size limit before queuing it for upload. Files that exceed the limit SHALL be silently dropped from the upload batch; for each dropped file the system SHALL display a toast identifying the file name and the applicable size limit. Files within the limit SHALL proceed to upload normally.

#### Scenario: Oversized file dropped from batch
- **WHEN** the user selects or drops one or more files into the artpiece files component and at least one exceeds its category size limit
- **THEN** the oversized file is not queued for upload, a toast is shown naming the file and the limit, and any files within the limit upload normally

#### Scenario: All files within limit
- **WHEN** the user selects or drops files and all are within their category size limits
- **THEN** all files are queued for upload with no toast shown

### Requirement: Character suggestions in recency order
The characters field in the create artpiece form SHALL request the character list with `sort=recent` and display suggestions in the order returned, with no client-side re-sort. Already-selected characters SHALL be excluded from the suggestions.

#### Scenario: Suggestions follow API recency order
- **WHEN** the user focuses the characters field and the API returns characters in recency order A, B, C
- **THEN** the suggestions SHALL be listed A, B, C

#### Scenario: List requested with sort=recent
- **WHEN** the create artpiece modal opens
- **THEN** the characters request SHALL include `sort=recent`

#### Scenario: Characters refresh after saving an artpiece
- **WHEN** the user saves an artpiece with characters selected
- **THEN** the characters query SHALL be invalidated so the next picker shows the updated order

---

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

### Requirement: Characters field lists suggestions on focus
The characters field in the create artpiece form SHALL list suggestions when focused with empty text, in recency order, excluding already-selected characters. The empty-text list SHALL be capped at 10 with a non-selectable hint row when more exist, exactly as for the tag picker. This SHALL hold whether the modal was opened from the file Inbox or the Artpieces page.

#### Scenario: Focus lists recent characters
- **WHEN** the user focuses the empty characters field and the API returns characters A, B, C in recency order
- **THEN** A, B, C SHALL be listed in that order

#### Scenario: Focus list is capped
- **WHEN** the user has 25 characters and focuses the empty characters field
- **THEN** the first 10 SHALL be listed followed by a hint row indicating more can be found by typing

#### Scenario: Opened from the Inbox
- **WHEN** the modal is opened from the file Inbox and the user focuses the empty characters field
- **THEN** characters SHALL be listed as above

#### Scenario: Typing searches all characters
- **WHEN** the user types text matching a character beyond the first 10
- **THEN** that character SHALL be listed
