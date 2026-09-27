## MODIFIED Requirements

### Requirement: CharacterFormModal — avatar picker
The modal SHALL display an avatar area that shows the current avatar image if one exists, or a placeholder with a camera icon and the text "Click or drag to upload" if none. Clicking the avatar area OR dropping an image file onto it SHALL open a file picker (for click) or accept the dropped file (for drop), restricted to `image/jpeg` and `image/png`.

When a file is selected or dropped, the system SHALL:
1. Open `AvatarCropDialog` with the source file
2. On crop confirmation, set `localFile` to the cropped JPEG `File` and display a local preview via `URL.createObjectURL`
3. On crop cancellation, make no changes to `localFile` or the displayed preview

The preview SHALL be revoked when the modal closes. In edit mode, a button SHALL allow the user to remove the current avatar (local state only — not committed until Save).

#### Scenario: Avatar area shows placeholder when no avatar
- **WHEN** the modal opens for a character with `avatar_url: null` (or create mode)
- **THEN** the avatar area SHALL display a placeholder icon and the text "Click or drag to upload"

#### Scenario: Avatar area shows current avatar in edit mode
- **WHEN** the modal opens for a character with a non-null `avatar_url`
- **THEN** the avatar area SHALL display the existing avatar image

#### Scenario: File picker opens on avatar area click
- **WHEN** a user clicks the avatar area
- **THEN** a file picker opens accepting only `.jpg`/`.jpeg` and `.png` files

#### Scenario: Drop accepted on avatar area
- **WHEN** a user drags and drops an image file onto the avatar area
- **THEN** the file is accepted as if it were picked via the file picker

#### Scenario: Drag-active state shown while hovering
- **WHEN** a user drags a file over the avatar area
- **THEN** the avatar area displays a visual drag-active indicator (ring highlight)

#### Scenario: Crop dialog opens after file is picked
- **WHEN** a user picks or drops a file
- **THEN** `AvatarCropDialog` opens with that file as the source image

#### Scenario: Local preview shown after crop confirmed
- **WHEN** a user confirms the crop in `AvatarCropDialog`
- **THEN** the avatar area SHALL display a local preview of the cropped image

#### Scenario: Crop cancelled leaves avatar unchanged
- **WHEN** a user cancels the crop in `AvatarCropDialog`
- **THEN** the avatar area displays whatever it showed before the file was picked (existing avatar or placeholder)

#### Scenario: Avatar removed locally in edit mode
- **WHEN** a user clicks the remove button on the avatar preview in edit mode
- **THEN** the avatar area reverts to the placeholder; no network call is made until Save
