## ADDED Requirements

### Requirement: Upload modal opens from + button
The system SHALL open an "Add files to dump" modal when the user clicks the floating "+" button on the Inbox page.

#### Scenario: Clicking the + button
- **WHEN** a user clicks the floating "+" button on the Inbox page
- **THEN** the "Add files to dump" modal opens with a drop zone at the top and an empty file list below

---

### Requirement: File upload via modal
The system SHALL allow the user to upload multiple files simultaneously via the modal's drop zone or file picker. Each file's original filename SHALL be sent as the `name` field in the initiate-upload request.

#### Scenario: Dropping files onto the modal drop zone
- **WHEN** a user drops one or more files onto the drop zone inside the modal
- **THEN** each file is immediately added to the file list as an uploading row (spinner thumbnail, disabled name and notes inputs, no remove button), and each file's upload flow (initiate → PUT to R2 → complete) runs in parallel

#### Scenario: Selecting files via file picker
- **WHEN** a user clicks the drop zone and selects files via the file picker
- **THEN** the same uploading row behavior as dropping applies

#### Scenario: Upload completes successfully
- **WHEN** a file's complete-upload call succeeds
- **THEN** the uploading row transitions to an uploaded row: a thumbnail slot with a spinner (pending thumbnail generation), an enabled name input pre-filled with the filename, an enabled notes input, and a remove (X) button

#### Scenario: Upload fails
- **WHEN** any step of a file's upload flow (initiate, PUT, or complete) fails
- **THEN** that file's row is removed from the list and an error toast is shown naming the file

---

### Requirement: Thumbnail pending state in modal rows
The system SHALL display a spinner in the thumbnail slot of an uploaded row until `thumbnail_gen_state` is no longer `pending`.

#### Scenario: Thumbnail generation pending
- **WHEN** an uploaded row has `thumbnail_gen_state` of `pending`
- **THEN** the thumbnail slot shows a spinner; the name and notes inputs are disabled

#### Scenario: Thumbnail generation completes
- **WHEN** the file's `thumbnail_url` becomes non-null and `thumbnail_gen_state` is no longer `pending`
- **THEN** the thumbnail slot shows the thumbnail image; the name and notes inputs become enabled

#### Scenario: Thumbnail generation fails or not applicable
- **WHEN** `thumbnail_gen_state` is `failed` or `not_applicable`
- **THEN** the thumbnail slot shows a generic file icon; the name and notes inputs are enabled

---

### Requirement: Auto-save name and notes on blur
The system SHALL persist the name and notes fields of each uploaded row independently on blur, with no global Save action.

#### Scenario: Editing and blurring the name field
- **WHEN** a user edits the name field of an uploaded row and focuses away
- **THEN** the system calls the update-file endpoint with the current name and notes values; on success, no visible feedback is shown

#### Scenario: Editing and blurring the notes field
- **WHEN** a user edits the notes field of an uploaded row and focuses away
- **THEN** the system calls the update-file endpoint with the current name and notes values; on success, no visible feedback is shown

---

### Requirement: Inline field-save error feedback
The system SHALL display an inline error on a field when its auto-save fails, and clear the error once the field saves successfully.

#### Scenario: Name field save fails
- **WHEN** the update-file call triggered by a name blur fails
- **THEN** the name input displays a red border and a short error message below it (e.g. "Couldn't save. Try again."); the notes field is unaffected

#### Scenario: Notes field save fails
- **WHEN** the update-file call triggered by a notes blur fails
- **THEN** the notes input displays a red border and a short error message below it; the name field is unaffected

#### Scenario: Retry after field-save failure
- **WHEN** a user edits a field that has an inline error and blurs again, and the save succeeds this time
- **THEN** the red border and error message are cleared

---

### Requirement: Remove file row
The system SHALL allow the user to remove an uploaded row from the modal view via an X button. Removal only dismisses the row from the modal — the file remains in the dump.

#### Scenario: Clicking X on an uploaded row
- **WHEN** a user clicks the X button on an uploaded row
- **THEN** the row is removed from the modal's file list; the file remains in the dump and appears in the inbox grid

#### Scenario: No X while uploading
- **WHEN** a file is in the uploading state
- **THEN** no X button is shown for that row

---

### Requirement: Modal close
The system SHALL close the modal when the user clicks Close, without discarding or committing any pending state (there is none).

#### Scenario: Clicking Close
- **WHEN** a user clicks the Close button in the modal footer
- **THEN** the modal closes; all successfully uploaded files remain in the dump

#### Scenario: Closing while uploads are in flight
- **WHEN** a user clicks Close while one or more files are still uploading
- **THEN** the modal closes; in-flight uploads continue to completion in the background; their results appear in the inbox grid when the grid polls
