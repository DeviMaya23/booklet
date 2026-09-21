## ADDED Requirements

### Requirement: Image create modal
The images page SHALL provide an `ImageFormModal` in create mode, opened by the `+ New` button, that allows the user to upload a new image with associated metadata.

#### Scenario: Modal opens on New button click
- **WHEN** an authenticated user clicks the `+ New` button on `/app/images`
- **THEN** the `ImageFormModal` SHALL open in create mode with all fields empty and the file picker area shown

#### Scenario: File picker area shown in create mode
- **WHEN** the modal is in create mode and no file has been selected
- **THEN** the upload area SHALL display a placeholder prompt and be clickable to open the file picker

#### Scenario: File type restricted to jpg and png
- **WHEN** the file picker opens
- **THEN** it SHALL accept only `image/jpeg` and `image/png` file types

#### Scenario: Preview shown after file selection
- **WHEN** the user selects a valid file
- **THEN** the upload area SHALL display a preview of the selected image

#### Scenario: File can be removed and replaced
- **WHEN** a file is selected and the user clicks the remove control (X icon)
- **THEN** the preview SHALL be cleared and the upload area SHALL return to the placeholder state, allowing a new file to be picked

#### Scenario: Save disabled without file
- **WHEN** no file has been selected
- **THEN** the Save button SHALL be disabled

#### Scenario: Successful create
- **WHEN** a user submits the form with a valid file selected
- **THEN** the app SHALL call `POST /images` (InitialUpload), PUT the file to the returned presigned URL, call `POST /images/:id/complete` (CompleteUpload), close the modal, invalidate the image list cache, and show a success toast

#### Scenario: Upload failure cleanup
- **WHEN** the R2 PUT or CompleteUpload call fails during create
- **THEN** the app SHALL call `DELETE /images/:id` to remove the orphaned record and show an error toast

### Requirement: Image create form fields
The create modal SHALL include the following optional metadata fields: title (text input), notes (textarea), artist (single-select Combobox with search), and characters (multi-token input). None of these fields are required.

#### Scenario: Title field accepts text
- **WHEN** the user types in the Title field
- **THEN** the value SHALL be included in the `POST /images` InitialUpload request body

#### Scenario: Notes field accepts text
- **WHEN** the user types in the Notes field
- **THEN** the value SHALL be included in the `POST /images` InitialUpload request body

#### Scenario: Artist field filters by name
- **WHEN** the user types in the Artist Combobox input
- **THEN** the dropdown SHALL show only artists whose name contains the typed text (case-insensitive, client-side)

#### Scenario: Artist selected from dropdown
- **WHEN** the user selects an artist from the Combobox
- **THEN** that artist's ID SHALL be included as `artist_id` in the `POST /images` request body

#### Scenario: Artist inline creation
- **WHEN** the user clicks the `+` button next to the Artist field
- **THEN** the `ArtistFormModal` SHALL open as a stacked dialog; on successful creation the new artist SHALL be auto-selected in the Combobox without requiring the user to search for it

#### Scenario: Characters field filters by name
- **WHEN** the user types in the Characters token input
- **THEN** the dropdown SHALL show only characters whose name contains the typed text (case-insensitive, client-side), excluding already-selected characters

#### Scenario: Characters selected as tokens
- **WHEN** the user selects characters from the dropdown
- **THEN** each selected character SHALL appear as a removable token; their IDs SHALL be included as `character_ids` in the `POST /images` request body
