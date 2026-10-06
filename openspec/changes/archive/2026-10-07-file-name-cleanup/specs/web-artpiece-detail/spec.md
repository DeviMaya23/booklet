## MODIFIED Requirements

### Requirement: Edit mode — file upload
The detail view in edit mode SHALL display an "Add files / or drop them here" drop zone as the last cell in the file grid. Files dropped or picked upload immediately and attach to the artpiece without staging. The file extension SHALL be stripped from `File.name` before it is sent as `name` in the initiate-upload request.

#### Scenario: Add-files cell is last in the grid
- **WHEN** the detail view is in edit mode
- **THEN** the last cell in the file grid SHALL be an "Add files / or drop them here" tile with a dashed border; dropping files onto it or clicking it triggers the file picker

#### Scenario: File upload attaches immediately
- **WHEN** the user drops files onto the add-files cell or picks files via the file picker
- **THEN** each file SHALL be uploaded (init → S3 PUT → complete) using the file's name with extension stripped as the `name` field (e.g. `pic.jpg` → `pic`), and then immediately attached to the artpiece via `POST /artpieces/:id/files/:file_id`; the new file thumbnail SHALL appear in the grid on attachment

#### Scenario: Uploaded files persist through Cancel
- **WHEN** the user uploads one or more files in edit mode and then clicks Cancel
- **THEN** those files SHALL remain attached to the artpiece; they SHALL appear in the file grid on the next view-mode render
