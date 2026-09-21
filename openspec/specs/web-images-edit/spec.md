## Purpose

Defines the behaviour of the `ImageFormModal` in edit mode: loading and pre-populating image data, updating metadata, downloading the full-resolution image, and deleting an image from the modal.

## Requirements

### Requirement: Image edit modal
The images page SHALL provide an `ImageFormModal` in edit mode, opened by clicking an image card, that allows the user to view and update the image's metadata.

#### Scenario: Modal opens on card click
- **WHEN** an authenticated user clicks an image card on `/app/images`
- **THEN** the `ImageFormModal` SHALL open in edit mode and fetch the image's full data from `GET /images/:id`

#### Scenario: Fields pre-populated with current values
- **WHEN** the edit modal finishes loading the image data
- **THEN** the title, notes, artist, and character fields SHALL be pre-populated with the image's current values

#### Scenario: Loading state while fetching
- **WHEN** the edit modal is open and the `GET /images/:id` request is in flight
- **THEN** the form fields SHALL be disabled until the response is received

#### Scenario: Successful update
- **WHEN** the user edits fields and clicks Save
- **THEN** the app SHALL call `PUT /images/:id` with the updated title, notes, artist_id, and character_ids, close the modal, invalidate the image list cache, and show a success toast

#### Scenario: Update failure
- **WHEN** the `PUT /images/:id` call fails
- **THEN** the modal SHALL remain open and show an error toast

#### Scenario: Image preview shown in edit mode
- **WHEN** the edit modal opens
- **THEN** the image preview area SHALL display the image's thumbnail (or a placeholder if no thumbnail exists)

#### Scenario: File upload controls not shown in edit mode
- **WHEN** the modal is in edit mode
- **THEN** there SHALL be no file picker or remove-file control visible

### Requirement: Download image in edit mode
The edit modal SHALL display a Download image button that triggers a browser download of the full-resolution image using the presigned `image_url` from `GET /images/:id`.

#### Scenario: Download button triggers file download
- **WHEN** the user clicks Download image in the edit modal
- **THEN** the browser SHALL initiate a download of the full-resolution image using the presigned URL

#### Scenario: Download button not shown in create mode
- **WHEN** the modal is in create mode
- **THEN** no Download image button SHALL be visible

### Requirement: Delete image from edit modal
The edit modal SHALL include a Delete button that opens a confirmation dialog. Confirming SHALL delete the image and close the modal.

#### Scenario: Delete button opens confirmation
- **WHEN** the user clicks Delete in the edit modal
- **THEN** a confirmation dialog SHALL open

#### Scenario: Delete confirmed
- **WHEN** the user confirms deletion
- **THEN** the app SHALL call `DELETE /images/:id`, close the modal, invalidate the image list cache, and show a success toast

#### Scenario: Delete cancelled
- **WHEN** the user cancels the confirmation dialog
- **THEN** the dialog SHALL close and the edit modal SHALL remain open unchanged
