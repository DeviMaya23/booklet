## ADDED Requirements

### Requirement: Images list display
The images page at `/app/images` SHALL fetch the authenticated user's images from `GET /images` and display them as a responsive grid of cards. Each card SHALL display the image's thumbnail if `thumbnail_url` is non-null, or a "no image" placeholder icon on a muted background if null. The image's `title` SHALL appear as an overlay in the bottom-left of the card; if `title` is null, the card SHALL display `"Untitled"` as a cosmetic fallback.

#### Scenario: Images load and display
- **WHEN** an authenticated user navigates to `/app/images`
- **THEN** the page SHALL fetch `GET /images` and render one card per image in a responsive grid

#### Scenario: Image with thumbnail
- **WHEN** an image has a non-null `thumbnail_url`
- **THEN** its card SHALL display the thumbnail as the card's background image

#### Scenario: Image without thumbnail
- **WHEN** an image has a null `thumbnail_url`
- **THEN** its card SHALL display a muted placeholder background with a "no image" icon

#### Scenario: Image with null title
- **WHEN** an image has a null `title`
- **THEN** its card SHALL display the label `"Untitled"`

### Requirement: Images client-side search
The images page SHALL provide a search bar that filters the displayed image list client-side by `title`. Filtering SHALL be case-insensitive. Images with a null `title` SHALL NOT match any search query, including the literal string `"Untitled"`.

#### Scenario: Search matches images by title
- **WHEN** a user types into the search bar
- **THEN** only images whose `title` contains the search text (case-insensitive) SHALL be displayed

#### Scenario: Null-title images excluded from search results
- **WHEN** a user types any text into the search bar
- **THEN** images with a null `title` SHALL NOT appear in the results

#### Scenario: Empty search shows all images
- **WHEN** the search bar is empty
- **THEN** all fetched images SHALL be displayed

### Requirement: New image button
The images page SHALL display a `+ New` button. The button SHALL be present but perform no action.

#### Scenario: New button is visible but no-op
- **WHEN** an authenticated user views `/app/images`
- **THEN** a `+ New` button SHALL be visible and clicking it SHALL have no effect

### Requirement: Delete image
Each image card SHALL expose a `(...)` menu in the top-right corner. The menu SHALL contain a Delete option. Selecting Delete SHALL open a confirmation dialog. Confirming SHALL call `DELETE /images/:id`, remove the image from the list, and show a success toast. Cancelling SHALL close the dialog with no change.

#### Scenario: Delete menu opens
- **WHEN** a user clicks the `(...)` button on an image card
- **THEN** a dropdown menu SHALL appear with a Delete option

#### Scenario: Confirmation dialog appears
- **WHEN** a user selects Delete from the image card menu
- **THEN** a confirmation dialog SHALL open asking the user to confirm deletion

#### Scenario: Delete confirmed
- **WHEN** a user confirms deletion in the dialog
- **THEN** the app SHALL call `DELETE /images/:id`, close the dialog, refetch the image list, and show a success toast

#### Scenario: Delete cancelled
- **WHEN** a user cancels deletion in the dialog
- **THEN** the dialog SHALL close and the image list SHALL remain unchanged
