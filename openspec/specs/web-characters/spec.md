## Purpose

Defines the behaviour of the Characters page (`/app/characters`) in the web application: listing a user's characters, searching them client-side, and deleting individual characters.

## Requirements

### Requirement: Characters list display
The characters page at `/app/characters` SHALL fetch the authenticated user's characters from `GET /characters` and display them as a responsive grid of cards. Each card SHALL display the character's avatar image if available, or a "no image" placeholder icon on a muted background if `avatar_url` is null. The character's `name` SHALL appear as an overlay in the bottom-left of the card.

#### Scenario: Characters load and display
- **WHEN** an authenticated user navigates to `/app/characters`
- **THEN** the page SHALL fetch `GET /characters` and render one card per character in a responsive grid

#### Scenario: Character with avatar
- **WHEN** a character has a non-null `avatar_url`
- **THEN** its card SHALL display the avatar as the card's background image

#### Scenario: Character without avatar
- **WHEN** a character has a null `avatar_url`
- **THEN** its card SHALL display a muted placeholder background with a "no image" icon

### Requirement: Characters client-side search
The characters page SHALL provide a search bar that filters the displayed character list client-side by `name`. Filtering SHALL be case-insensitive.

#### Scenario: Search matches characters by name
- **WHEN** a user types into the search bar
- **THEN** only characters whose `name` contains the search text (case-insensitive) SHALL be displayed

#### Scenario: Empty search shows all characters
- **WHEN** the search bar is empty
- **THEN** all fetched characters SHALL be displayed

### Requirement: New character button
The characters page SHALL display a `+ New` button. Clicking the button SHALL open `CharacterFormModal` in create mode (empty form, no pre-populated values).

#### Scenario: New button opens modal in create mode
- **WHEN** an authenticated user clicks the `+ New` button on the characters page
- **THEN** `CharacterFormModal` SHALL open with an empty form and the title "New character"

### Requirement: Delete character
Each character card SHALL expose a `(...)` menu in the top-right corner. The menu SHALL contain a Delete option. Selecting Delete SHALL open a confirmation dialog. Confirming SHALL call `DELETE /characters/:id`, remove the character from the list, and show a success toast. Cancelling SHALL close the dialog with no change.

#### Scenario: Delete menu opens
- **WHEN** a user clicks the `(...)` button on a character card
- **THEN** a dropdown menu SHALL appear with a Delete option

#### Scenario: Confirmation dialog appears
- **WHEN** a user selects Delete from the character card menu
- **THEN** a confirmation dialog SHALL open asking the user to confirm deletion

#### Scenario: Delete confirmed
- **WHEN** a user confirms deletion in the dialog
- **THEN** the app SHALL call `DELETE /characters/:id`, close the dialog, refetch the character list, and show a success toast

#### Scenario: Delete cancelled
- **WHEN** a user cancels deletion in the dialog
- **THEN** the dialog SHALL close and the character list SHALL remain unchanged

### Requirement: Click character card to edit
Each character card SHALL be clickable. Clicking the card body (outside the `(...)` menu) SHALL open `CharacterFormModal` in edit mode, pre-populated with the character's current data.

#### Scenario: Card click opens modal in edit mode
- **WHEN** a user clicks a character card body
- **THEN** `CharacterFormModal` SHALL open with the title "Edit character" and form fields pre-populated with the character's `name`, `notes`, and `avatar_url`

### Requirement: CharacterFormModal — create character
`CharacterFormModal` in create mode SHALL allow the user to enter a name (required) and notes (optional), optionally pick an avatar image (jpg/png only), and submit. On submit:
1. `POST /characters` is called with `name` and `notes`
2. If an avatar was picked, the avatar upload sequence runs (init → PUT to presigned URL → complete)
3. On full success: modal closes, character list is refetched, success toast shown
4. If avatar upload fails after character creation: `DELETE /characters/:id` is called to clean up, an error toast is shown, and the modal stays open

The modal footer SHALL contain: Cancel (left-aligned or secondary), Save (primary).

#### Scenario: Create with name only
- **WHEN** a user enters a name, leaves notes empty, picks no avatar, and clicks Save
- **THEN** `POST /characters` is called with `{ name }`, the modal closes on success, and a success toast is shown

#### Scenario: Create with avatar
- **WHEN** a user enters a name and picks an avatar image before clicking Save
- **THEN** `POST /characters` is called first, followed by init → R2 PUT → complete; on full success the modal closes and a success toast is shown

#### Scenario: Avatar upload fails on create — character is cleaned up
- **WHEN** the avatar upload sequence fails after the character has been created
- **THEN** `DELETE /characters/:id` is called, an error toast is shown, and the modal stays open for retry

#### Scenario: Name is required
- **WHEN** a user attempts to Save with an empty name field
- **THEN** the Save button SHALL be disabled and no network call is made

#### Scenario: Cancel discards form
- **WHEN** a user clicks Cancel
- **THEN** the modal closes and no network call is made; any picked avatar preview is discarded

### Requirement: CharacterFormModal — edit character
`CharacterFormModal` in edit mode SHALL allow the user to update name, notes, and avatar. On Save:
- If a new avatar was picked: init → R2 PUT → complete runs first, then `PUT /characters/:id`
- If the avatar was cleared: `DELETE /characters/:id/avatar` runs first, then `PUT /characters/:id`
- If only metadata changed: `PUT /characters/:id` is called
On success: modal closes, character list is refetched, success toast shown.

The modal footer SHALL contain: Delete (destructive, left-aligned), Cancel, Save.

#### Scenario: Save metadata only
- **WHEN** a user edits name or notes without changing the avatar and clicks Save
- **THEN** `PUT /characters/:id` is called with the updated fields; on success the modal closes and a success toast is shown

#### Scenario: Save with new avatar
- **WHEN** a user picks a new avatar and clicks Save
- **THEN** avatar init → R2 PUT → complete runs before `PUT /characters/:id`; on full success the modal closes and a success toast is shown

#### Scenario: Save after clearing avatar
- **WHEN** a user removes the existing avatar preview and clicks Save
- **THEN** `DELETE /characters/:id/avatar` is called before `PUT /characters/:id`; on success the modal closes and a success toast is shown

#### Scenario: Delete opens confirmation dialog
- **WHEN** a user clicks the Delete button in the modal footer
- **THEN** a confirmation AlertDialog opens asking the user to confirm deletion

#### Scenario: Delete confirmed
- **WHEN** a user confirms deletion in the AlertDialog
- **THEN** `DELETE /characters/:id` is called; on success the modal closes, the list is refetched, and a success toast is shown

#### Scenario: Delete cancelled
- **WHEN** a user dismisses the confirmation AlertDialog
- **THEN** the AlertDialog closes and the character is unchanged

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

### Requirement: ResourceCard onClick
`ResourceCard` SHALL accept an optional `onClick` prop. When provided, clicking the card body (the area outside the `(...)` dropdown trigger) SHALL invoke it.

#### Scenario: Card body click invokes onClick
- **WHEN** `onClick` is provided and the user clicks the card body
- **THEN** the `onClick` handler is called

#### Scenario: Dropdown interaction does not trigger onClick
- **WHEN** a user opens or interacts with the `(...)` dropdown menu
- **THEN** the `onClick` handler is NOT called
