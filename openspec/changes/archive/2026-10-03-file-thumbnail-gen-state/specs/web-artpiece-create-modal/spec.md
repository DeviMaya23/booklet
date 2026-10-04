## MODIFIED Requirements

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

### Requirement: Files sub-component — upload lifecycle
The system SHALL manage the per-file upload flow within the modal, including thumbnail readiness and error handling.

#### Scenario: Upload completes — thumbnail pending
- **WHEN** the complete-upload call succeeds and `thumbnail_gen_state === "pending"`
- **THEN** the row remains with a spinner in the thumbnail slot and disabled name/notes fields

#### Scenario: Upload completes — thumbnail done
- **WHEN** `thumbnail_gen_state` becomes `"done"` for a file
- **THEN** the spinner is replaced with the thumbnail image and the name/notes fields become editable

#### Scenario: Upload completes — thumbnail failed or not applicable
- **WHEN** `thumbnail_gen_state` becomes `"failed"` or `"not_applicable"` for a file
- **THEN** the spinner is replaced with a fallback icon and the name/notes fields become editable

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

### Requirement: Files sub-component — file name and notes editing
The system SHALL allow the user to edit a file's name and notes inline, saving each change on blur.

#### Scenario: Edit file name and blur
- **WHEN** the user edits the name field of a file row and moves focus away
- **THEN** `PUT /files/:id` is called with the updated name; the field shows the saved value

#### Scenario: Edit file notes and blur
- **WHEN** the user edits the notes field of a file row and moves focus away
- **THEN** `PUT /files/:id` is called with the updated notes; the field shows the saved value

#### Scenario: Name and notes fields disabled during thumbnail pending
- **WHEN** a file row has `thumbnail_gen_state === "pending"`
- **THEN** the name and notes fields are disabled and cannot be edited
