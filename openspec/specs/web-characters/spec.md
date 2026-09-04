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
The characters page SHALL display a `+ New` button. The button SHALL be present but perform no action.

#### Scenario: New button is visible but no-op
- **WHEN** an authenticated user views `/app/characters`
- **THEN** a `+ New` button SHALL be visible and clicking it SHALL have no effect

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
