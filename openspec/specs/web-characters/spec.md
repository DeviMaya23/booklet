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

### Requirement: ResourceCard onClick
`ResourceCard` SHALL accept an optional `onClick` prop. When provided, clicking the card body (the area outside the `(...)` dropdown trigger) SHALL invoke it.

#### Scenario: Card body click invokes onClick
- **WHEN** `onClick` is provided and the user clicks the card body
- **THEN** the `onClick` handler is called

#### Scenario: Dropdown interaction does not trigger onClick
- **WHEN** a user opens or interacts with the `(...)` dropdown menu
- **THEN** the `onClick` handler is NOT called

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
