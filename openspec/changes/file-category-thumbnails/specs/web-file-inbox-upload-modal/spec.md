# Spec Delta

## MODIFIED Requirements

### Requirement: File upload via modal
The system SHALL allow the user to upload multiple files simultaneously via the modal's drop zone or file picker. The file extension SHALL be stripped from `File.name` before it is sent as `name` in the initiate-upload request and before the name input is pre-filled. The lowercased file extension SHALL be sent as the `extension` field in the same initiate-upload request.

#### Scenario: Dropping files onto the modal drop zone
- **WHEN** a user drops one or more files onto the drop zone inside the modal
- **THEN** each file is immediately added to the file list as an uploading row (spinner thumbnail, disabled name and notes inputs, no remove button), and each file's upload flow (initiate → PUT to R2 → complete) runs in parallel, with the lowercased extension sent as `extension` in the initiate request

#### Scenario: Selecting files via file picker
- **WHEN** a user clicks the drop zone and selects files via the file picker
- **THEN** the same uploading row behavior as dropping applies

#### Scenario: Upload completes successfully
- **WHEN** a file's complete-upload call succeeds
- **THEN** the uploading row transitions to an uploaded row: a thumbnail slot with a spinner (pending thumbnail generation), an enabled name input pre-filled with the filename with extension stripped (e.g. `pic.jpg` → `pic`), an enabled notes input, and a remove (X) button

#### Scenario: Upload fails
- **WHEN** any step of a file's upload flow (initiate, PUT, or complete) fails
- **THEN** that file's row is removed from the list and an error toast is shown naming the file
