# Spec Delta

## MODIFIED Requirements

### Requirement: CharacterFormModal — avatar picker
The modal SHALL display a compact avatar thumbnail beside the Name field. The thumbnail shows the current avatar if one exists; clicking it or dropping a file onto it opens the crop flow. When an avatar exists, a "Change" label overlay SHALL appear on the thumbnail. When no avatar exists, the thumbnail area SHALL show a placeholder indicating the user can click or drag to upload. Only `image/jpeg` and `image/png` are accepted.

When a file is selected or dropped, the system SHALL:
1. Open `AvatarCropDialog` with the source file
2. On crop confirmation, set `localFile` to the cropped JPEG `File` and display a local preview via `URL.createObjectURL`
3. On crop cancellation, make no changes to `localFile` or the displayed preview

The preview SHALL be revoked when the modal closes. In edit mode, an X button SHALL allow the user to remove the current avatar (local state only — not committed until Save).

#### Scenario: Avatar area shows placeholder when no avatar
- **WHEN** the modal opens for a character with `avatar_url: null` (or create mode)
- **THEN** the thumbnail area SHALL display a placeholder indicating the user can upload

#### Scenario: Avatar area shows current avatar in edit mode
- **WHEN** the modal opens for a character with a non-null `avatar_url`
- **THEN** the thumbnail SHALL display the existing avatar image with a "Change" overlay

#### Scenario: File picker opens on avatar area click
- **WHEN** a user clicks the avatar thumbnail
- **THEN** a file picker opens accepting only `.jpg`/`.jpeg` and `.png` files

#### Scenario: Drop accepted on avatar area
- **WHEN** a user drags and drops an image file onto the avatar thumbnail
- **THEN** the file is accepted as if it were picked via the file picker

#### Scenario: Drag-active state shown while hovering
- **WHEN** a user drags a file over the avatar thumbnail
- **THEN** the thumbnail displays a visual drag-active indicator

#### Scenario: Crop dialog opens after file is picked
- **WHEN** a user picks or drops a file
- **THEN** `AvatarCropDialog` opens with that file as the source image

#### Scenario: Local preview shown after crop confirmed
- **WHEN** a user confirms the crop in `AvatarCropDialog`
- **THEN** the avatar thumbnail SHALL display a local preview of the cropped image

#### Scenario: Crop cancelled leaves avatar unchanged
- **WHEN** a user cancels the crop in `AvatarCropDialog`
- **THEN** the thumbnail displays whatever it showed before the file was picked

#### Scenario: Avatar removed locally in edit mode
- **WHEN** a user clicks the X button on the avatar thumbnail in edit mode
- **THEN** the thumbnail reverts to the placeholder; no network call is made until Save
