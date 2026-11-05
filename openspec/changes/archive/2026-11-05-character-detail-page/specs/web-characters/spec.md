# Spec Delta

## ADDED Requirements

### Requirement: Character card opens the detail page
Each character card SHALL be clickable. Clicking the card body (outside the `(...)` menu) SHALL navigate to the character's detail page at `/app/characters/:id`. It SHALL NOT open `CharacterFormModal`.

#### Scenario: Card click navigates to the detail page
- **WHEN** a user clicks a character card body
- **THEN** the browser navigates to `/app/characters/:id` for that character and no modal opens

#### Scenario: Dropdown interaction does not navigate
- **WHEN** a user opens or interacts with the `(...)` menu on a card
- **THEN** the browser does not navigate

### Requirement: CharacterFormModal — create-mode avatar picker
The create modal SHALL display a compact avatar thumbnail beside the Name field. Clicking it or dropping a file onto it opens the crop flow. Once a cropped avatar is chosen, a "Change" overlay SHALL appear; before that, a placeholder is shown. Only `image/jpeg` and `image/png` are accepted. The preview SHALL be revoked when the modal closes.

#### Scenario: Placeholder shown before an avatar is chosen
- **WHEN** the modal opens
- **THEN** the thumbnail area displays a placeholder indicating the user can upload

#### Scenario: Chosen avatar shows a Change overlay
- **WHEN** the user has confirmed a crop
- **THEN** the thumbnail displays the cropped preview with a "Change" overlay

#### Scenario: Picker opens on thumbnail click
- **WHEN** a user clicks the avatar thumbnail
- **THEN** a file picker opens accepting only `.jpg`/`.jpeg` and `.png` files

#### Scenario: Drop accepted on thumbnail
- **WHEN** a user drops an image file onto the avatar thumbnail
- **THEN** the file is accepted as if it were picked via the file picker

#### Scenario: Drag-active indicator while hovering
- **WHEN** a user drags a file over the avatar thumbnail
- **THEN** the thumbnail displays a visual drag-active indicator

#### Scenario: Crop dialog opens after a file is picked
- **WHEN** a user picks or drops a file
- **THEN** `AvatarCropDialog` opens with that file as the source image

#### Scenario: Preview shown after crop confirmed
- **WHEN** a user confirms the crop
- **THEN** the thumbnail displays a local preview of the cropped image

#### Scenario: Crop cancelled leaves the thumbnail unchanged
- **WHEN** a user cancels the crop
- **THEN** the thumbnail displays whatever it showed before the file was picked

## REMOVED Requirements

### Requirement: Click character card to edit
**Reason**: Card click now opens the detail page instead of an edit modal.
**Migration**: See "Character card opens the detail page" in this capability and `web-character-detail`.

### Requirement: CharacterFormModal — edit character
**Reason**: Editing moves to the character detail page (`web-character-detail`). `CharacterFormModal` is create-only, so its edit mode, in-modal Delete button and edit-mode avatar removal no longer exist.
**Migration**: Use `/app/characters/:id` and its Edit button. Save semantics (avatar upload or removal before `PUT /characters/:id`) and delete confirmation are specified in `web-character-detail`.

### Requirement: CharacterFormModal — avatar picker
**Reason**: Replaced by "CharacterFormModal — create-mode avatar picker", which drops the edit-mode scenarios (existing avatar, X removal).
**Migration**: Create-mode behavior is unchanged; removal of an existing avatar is specified in `web-character-detail`.
