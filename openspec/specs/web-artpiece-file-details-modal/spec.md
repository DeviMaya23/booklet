# web-artpiece-file-details-modal Specification

## Purpose

Defines the "Edit details" modal on the artpiece detail view — a per-artpiece list of every attached file (including hidden ones) with auto-saving name, notes, and visibility controls on each row.

## Requirements

### Requirement: Edit details modal opens from the files-section header
The system SHALL show an "Edit details" entry point in the artpiece detail view's files-section header, in both view mode and edit mode. Clicking it opens a modal listing every file currently attached to the artpiece, including hidden ones.

#### Scenario: Opening from view mode
- **WHEN** the artpiece detail view is in view mode and the user clicks "Edit details"
- **THEN** the modal opens listing every attached file as a row, without requiring the user to enter edit mode

#### Scenario: Opening from edit mode
- **WHEN** the artpiece detail view is in edit mode and the user clicks "Edit details"
- **THEN** the modal opens listing every attached file as a row, independent of any pending staged changes (cover, remove, upload) in the edit session

#### Scenario: No attached files
- **WHEN** the artpiece has no attached files and the user clicks "Edit details"
- **THEN** the modal opens showing an empty state instead of a row list

### Requirement: Modal row shows thumbnail, name, notes, and visibility
Each row in the modal SHALL display the file's thumbnail (or fallback icon), an editable name input pre-filled with the file's current name, an editable notes input pre-filled with the file's current notes, and a visibility toggle reflecting the file's current `is_visible` value.

#### Scenario: Row reflects current values
- **WHEN** the modal opens
- **THEN** each row's name input, notes input, and visibility toggle are pre-filled from that file's current `name`, `notes`, and `is_visible` values

#### Scenario: Hidden file row is visually distinguished
- **WHEN** a row's file has `is_visible: false`
- **THEN** the row's visibility toggle indicates "Hidden" state distinctly from a visible row's toggle

### Requirement: Auto-save name and notes on blur
The system SHALL persist each row's name and notes fields independently on blur via the existing file-update endpoint, with no global Save action for the modal. A blur SHALL send the blurred field's current draft together with the other field's last successfully saved value, so an unsaved or failed draft in one field is never sent as a side effect of saving the other. A blur on a field whose value equals its last successfully saved value SHALL NOT call the endpoint. A cleared (empty) name or notes input SHALL be sent as null, not as an empty string.

#### Scenario: Editing and blurring the name field
- **WHEN** a user edits a row's name field and focuses away
- **THEN** the system calls the update-file endpoint with the name draft and the notes field's last saved value; on success, no visible feedback is shown

#### Scenario: Editing and blurring the notes field
- **WHEN** a user edits a row's notes field and focuses away
- **THEN** the system calls the update-file endpoint with the notes draft and the name field's last saved value; on success, no visible feedback is shown

#### Scenario: Failed draft in one field is not sent with the other
- **WHEN** a notes save fails and the user then edits and blurs the name field
- **THEN** the update-file call carries the new name draft and the notes field's last successfully saved value, not the failed notes draft

#### Scenario: Blurring an unchanged field
- **WHEN** a user focuses a name or notes field and blurs it without changing its value since the last successful save
- **THEN** no update-file call is made

#### Scenario: Clearing a field
- **WHEN** a user empties a row's name or notes input and focuses away
- **THEN** the system sends null for that field, and the stored value is cleared

### Requirement: Auto-save visibility on toggle
The system SHALL persist a row's visibility change immediately when the toggle is activated, with no staging and no dependency on the artpiece detail view's edit-mode Save/Cancel flow.

#### Scenario: Hiding a visible file
- **WHEN** a user toggles a visible row to hidden
- **THEN** the system calls the update-file endpoint with `is_visible: false`; on success the row's toggle reflects the hidden state

#### Scenario: Unhiding a hidden file
- **WHEN** a user toggles a hidden row to visible
- **THEN** the system calls the update-file endpoint with `is_visible: true`; on success the row's toggle reflects the visible state

#### Scenario: Hiding the current cover reassigns it
- **WHEN** a user hides a file that is the artpiece's current cover
- **THEN** the visibility toggle call succeeds, the artpiece's cover is reassigned by the backend, and the artpiece detail view reflects the new cover without a separate user action

### Requirement: Inline field-save error feedback
The system SHALL display an inline error on a field or toggle when its auto-save fails, and clear the error once that field saves successfully.

#### Scenario: Name field save fails
- **WHEN** the update-file call triggered by a name blur fails
- **THEN** the name input displays a red border and a short error message below it; the row's notes field and visibility toggle are unaffected

#### Scenario: Notes field save fails
- **WHEN** the update-file call triggered by a notes blur fails
- **THEN** the notes input displays a red border and a short error message below it; the row's name field and visibility toggle are unaffected

#### Scenario: Visibility toggle save fails
- **WHEN** the update-file call triggered by a visibility toggle fails
- **THEN** the toggle reverts to its prior state and a short error message is shown on the row; the row's name and notes fields are unaffected

#### Scenario: Retry after field-save failure
- **WHEN** a user edits a field or toggle that has an inline error and the save succeeds this time
- **THEN** the error state for that control is cleared

### Requirement: Modal close
The system SHALL close the modal when the user clicks Close, without discarding or committing any pending state (there is none — every field and toggle is saved independently as it changes).

#### Scenario: Clicking Close
- **WHEN** a user clicks the Close button in the modal footer
- **THEN** the modal closes; all previously saved name, notes, and visibility changes remain in effect
