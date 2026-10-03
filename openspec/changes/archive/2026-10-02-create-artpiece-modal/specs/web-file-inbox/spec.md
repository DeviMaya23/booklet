## MODIFIED Requirements

### Requirement: Selection toolbar
The system SHALL show a toolbar button when one or more files are selected.

#### Scenario: At least one file selected
- **WHEN** the user has one or more files selected
- **THEN** an "Add to artpiece ▾" split-button appears in the top-right of the main area with two dropdown items: "New Artpiece" and "Add to Existing Artpiece"

#### Scenario: No files selected
- **WHEN** no files are selected
- **THEN** the "Add to artpiece ▾" button is not shown

#### Scenario: Clicking "New Artpiece" from selection toolbar
- **WHEN** the user clicks "New Artpiece" in the dropdown
- **THEN** the "Create New Artpiece" modal opens with the currently selected files pre-populated

#### Scenario: Clicking "Add to Existing Artpiece" from selection toolbar
- **WHEN** the user clicks "Add to Existing Artpiece" in the dropdown
- **THEN** no action is taken (no-op placeholder for a future proposal)

---

### Requirement: Right-click context menu
The system SHALL show a context menu on right-click containing artpiece actions and a destructive delete action.

#### Scenario: Context menu items
- **WHEN** the context menu is open
- **THEN** it shows "New Artpiece", "Add to Existing Artpiece" (no-op), a visual divider, and "Delete"

#### Scenario: Clicking "New Artpiece" from context menu
- **WHEN** the user clicks "New Artpiece" in the context menu
- **THEN** the "Create New Artpiece" modal opens with the right-clicked selection pre-populated

#### Scenario: Clicking "Add to Existing Artpiece" from context menu
- **WHEN** the user clicks "Add to Existing Artpiece" in the context menu
- **THEN** no action is taken (no-op placeholder for a future proposal)

#### Scenario: Delete from context menu
- **WHEN** the user clicks "Delete" in the context menu
- **THEN** a confirmation dialog appears showing how many files will be deleted; on confirm, the selected files are deleted via bulk delete
