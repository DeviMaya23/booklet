## Purpose

Defines the per-file edit overlay that appears over the inbox grid when a user double-clicks a file tile. Allows inline editing of the file's name and notes with auto-save on blur.

## Requirements

### Requirement: Edit overlay opens on tile double-click
The system SHALL open a per-file edit overlay when the user double-clicks a file tile. The overlay SHALL appear over the right portion of the grid without causing the grid to reflow.

#### Scenario: Double-clicking a tile
- **WHEN** a user double-clicks a file tile in the inbox grid
- **THEN** an edit overlay appears over the grid showing the file's thumbnail (or fallback icon), an editable name field pre-filled with the file's current name, and an editable notes field pre-filled with the file's current notes; the grid behind the overlay is not reflowed

#### Scenario: Closing the overlay with Escape
- **WHEN** the edit overlay is open and the user presses Escape
- **THEN** the overlay closes

#### Scenario: Closing the overlay by clicking outside
- **WHEN** the edit overlay is open and the user clicks outside it
- **THEN** the overlay closes

---

### Requirement: Auto-save name and notes in edit overlay
The system SHALL persist each field independently on blur, with no global Save action.

#### Scenario: Editing and blurring the name field
- **WHEN** a user edits the name field and focuses away
- **THEN** the system calls the update-file endpoint with the current name and notes values; on success, no visible feedback is shown

#### Scenario: Editing and blurring the notes field
- **WHEN** a user edits the notes field and focuses away
- **THEN** the system calls the update-file endpoint with the current name and notes values; on success, no visible feedback is shown

---

### Requirement: Inline field-save error feedback in edit overlay
The system SHALL display an inline error on a field when its auto-save fails, and clear the error once the field saves successfully.

#### Scenario: Name field save fails
- **WHEN** the update-file call triggered by a name blur fails
- **THEN** the name input displays a red border and a short error message below it (e.g. "Couldn't save. Try again."); the notes field is unaffected

#### Scenario: Notes field save fails
- **WHEN** the update-file call triggered by a notes blur fails
- **THEN** the notes input displays a red border and a short error message below it; the name field is unaffected

#### Scenario: Retry after field-save failure
- **WHEN** a user edits a field that has an inline error and blurs again, and the save succeeds
- **THEN** the red border and error message are cleared
