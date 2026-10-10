## Purpose

Defines the requirements for the Artists page (`/app/artists`) in the web application, covering list display with link chips, client-side search, and CRUD operations (create, edit, delete) via a shared artist form modal.

---

## Requirements

### Requirement: Artists list display
The artists page at `/app/artists` SHALL fetch the authenticated user's artists from `GET /artists` and display them in a table with columns: Name, Links, Notes, and an edit action. Each row SHALL display the artist's `name` in bold, a links chip group, the artist's `notes` truncated with ellipsis (or "—" if absent), and an edit icon button.

#### Scenario: Artists load and display
- **WHEN** an authenticated user navigates to `/app/artists`
- **THEN** the page SHALL fetch `GET /artists` and render one table row per artist

#### Scenario: Empty list
- **WHEN** the user has no artists
- **THEN** the page SHALL render an empty table body with no rows

---

### Requirement: Artists client-side search
The artists page SHALL provide a search bar that filters the displayed artist list client-side by `name`. Filtering SHALL be case-insensitive.

#### Scenario: Search matches artists by name
- **WHEN** a user types into the search bar
- **THEN** only artists whose `name` contains the search text (case-insensitive) SHALL be displayed

#### Scenario: Empty search shows all artists
- **WHEN** the search bar is empty
- **THEN** all fetched artists SHALL be displayed

---

### Requirement: Artist link chips
Each artist row's Links cell SHALL display the artist's links as chips, showing the hostname of each URL (e.g. `bsky.app` for `https://bsky.app/@artist`). Chips SHALL be ordered with the primary link first. A maximum of 3 chips SHALL be visible; if the artist has more than 3 links, a non-interactive `+N` chip SHALL be shown after the third, where N is the count of hidden links. If the artist has no links, the cell SHALL display the text "No links" in a muted style. Each chip SHALL be an anchor that opens the full URL in a new tab.

The primary link chip SHALL be visually distinguished with a star icon. All other chips SHALL render without a star.

If a URL cannot be parsed by `new URL()`, the chip SHALL display the raw URL string as its label.

#### Scenario: Chips show hostname only
- **WHEN** an artist has links
- **THEN** each visible chip SHALL display only the hostname portion of the URL, not the full URL

#### Scenario: Primary link chip is first and distinguished
- **WHEN** an artist has a primary link alongside other links
- **THEN** the primary link chip SHALL be rendered first and SHALL show a star icon

#### Scenario: Maximum 3 chips shown with overflow badge
- **WHEN** an artist has more than 3 links
- **THEN** exactly 3 chips SHALL be visible and a non-interactive `+N` chip SHALL indicate the number of hidden links

#### Scenario: No links placeholder
- **WHEN** an artist has no links
- **THEN** the Links cell SHALL render "No links" in muted text instead of chips

#### Scenario: Chip opens URL in new tab
- **WHEN** a user clicks a link chip
- **THEN** the full URL SHALL open in a new browser tab

---

### Requirement: Artist link chip tooltip
Hovering a link chip SHALL show a tooltip containing the full URL on the first line and a descriptor on the second line. For the primary link the descriptor SHALL read "Main link · opens in a new tab". For all other links the descriptor SHALL read "Opens in a new tab".

#### Scenario: Tooltip on primary link chip
- **WHEN** a user hovers the primary link chip
- **THEN** a tooltip SHALL appear showing the full URL and "Main link · opens in a new tab"

#### Scenario: Tooltip on non-primary link chip
- **WHEN** a user hovers a non-primary link chip
- **THEN** a tooltip SHALL appear showing the full URL and "Opens in a new tab"

---

### Requirement: Artist notes truncation
The Notes column SHALL display the artist's `notes` value. When the text is too long to fit the available column width it SHALL be clipped with a trailing ellipsis. When the artist has no notes, the cell SHALL display "—".

#### Scenario: Notes truncated when too long
- **WHEN** an artist's notes text exceeds the available column width
- **THEN** the text SHALL be clipped with a trailing "…"

#### Scenario: No notes placeholder
- **WHEN** an artist has no notes
- **THEN** the Notes cell SHALL display "—"

---

### Requirement: Create artist
The artists page SHALL display a `+ New` button. Clicking it SHALL open the artist form modal in create mode. Submitting the form SHALL call `POST /artists`, close the modal, refresh the list, and show a success toast. If the backend returns 409 (name conflict), the modal SHALL remain open and show an error toast.

#### Scenario: New button opens create modal
- **WHEN** a user clicks the `+ New` button
- **THEN** the artist form modal SHALL open in create mode with all fields empty

#### Scenario: Successful create
- **WHEN** a user submits the create modal with a valid name
- **THEN** the app SHALL call `POST /artists`, close the modal, refetch the artist list, and show a success toast

#### Scenario: Name conflict on create
- **WHEN** a user submits the create modal with a name that already exists for that user
- **THEN** the backend returns 409, the modal SHALL remain open, and an error toast SHALL appear indicating the name is already taken

#### Scenario: Cancel create
- **WHEN** a user closes or cancels the create modal without submitting
- **THEN** no API call SHALL be made and the artist list SHALL remain unchanged

---

### Requirement: Edit artist
Each artist row SHALL display an edit button. Clicking it SHALL open the artist form modal in edit mode, pre-populated with the artist's current `name`, `links`, and `notes`. Submitting SHALL call `PUT /artists/:id` with all fields including the full `links` array, close the modal, refresh the list, and show a success toast. If the backend returns 409, the modal SHALL remain open and show an error toast.

Empty optional fields (`notes`) SHALL be submitted as `null`. An empty `links` array SHALL be submitted as `[]`.

#### Scenario: Edit button opens edit modal pre-populated
- **WHEN** a user clicks the edit button on an artist row
- **THEN** the artist form modal SHALL open in edit mode with `name`, `links`, and `notes` pre-filled from the artist's current values

#### Scenario: Successful edit
- **WHEN** a user submits the edit modal with valid changes
- **THEN** the app SHALL call `PUT /artists/:id` with all fields, close the modal, refetch the artist list, and show a success toast

#### Scenario: Clearing all links
- **WHEN** a user removes all links and submits
- **THEN** the app SHALL send `"links": []` and the artist's stored links SHALL be cleared

#### Scenario: Name conflict on edit
- **WHEN** a user submits the edit modal with a name already used by another artist they own
- **THEN** the backend returns 409, the modal SHALL remain open, and an error toast SHALL appear indicating the name is already taken

#### Scenario: Cancel edit
- **WHEN** a user closes or cancels the edit modal without submitting
- **THEN** no API call SHALL be made and the artist list SHALL remain unchanged

---

### Requirement: Delete artist
The artist form modal in edit mode SHALL display a delete button. Clicking it SHALL open a confirmation dialog. Confirming SHALL call `DELETE /artists/:id`, close the modal, refresh the list, and show a success toast. Cancelling SHALL close the dialog and return to the edit modal.

#### Scenario: Delete button visible in edit mode only
- **WHEN** the artist form modal is opened in edit mode
- **THEN** a delete button SHALL be visible

#### Scenario: Delete button not visible in create mode
- **WHEN** the artist form modal is opened in create mode
- **THEN** no delete button SHALL be visible

#### Scenario: Confirmation dialog appears
- **WHEN** a user clicks the delete button in the edit modal
- **THEN** a confirmation dialog SHALL open asking the user to confirm deletion

#### Scenario: Delete confirmed
- **WHEN** a user confirms deletion
- **THEN** the app SHALL call `DELETE /artists/:id`, close the modal, refetch the artist list, and show a success toast

#### Scenario: Delete cancelled
- **WHEN** a user cancels the confirmation dialog
- **THEN** the dialog SHALL close and the edit modal SHALL remain open unchanged

---

### Requirement: Artist form modal fields
The artist form modal SHALL contain: `name` (required text input), a multi-link editor for `links`, and `notes` (optional textarea). The modal SHALL display Save and Cancel buttons. Submitting with an empty `name` SHALL be prevented client-side.

The multi-link editor SHALL:
- Allow adding multiple link entries via an `+ Add link` button
- Show a remove (`×`) button on each link entry
- Show a primary (star) toggle on each link entry when more than one link is present; when only one link exists, it is implicitly primary and the star toggle is hidden
- Enforce at most one primary link client-side: marking a link as primary deselects the previous primary
- Default the first link added to `is_primary: true`; when the primary link is removed, the first remaining link SHALL automatically become primary
- Auto-prepend `https://` to a link URL on blur if the user did not supply a scheme
- Validate each link URL client-side on blur using `new URL()`; display an inline error beneath the invalid input rather than a toast
- Filter out blank link entries silently on submit (entries with an empty URL are not sent)

#### Scenario: Submit blocked when name is empty
- **WHEN** a user attempts to submit the modal with an empty name field
- **THEN** the submission SHALL be blocked client-side and no API call SHALL be made

#### Scenario: Optional fields may be left empty
- **WHEN** a user submits the modal with only `name` filled and no links or notes
- **THEN** the submission SHALL proceed normally with `links: []`

#### Scenario: First link is implicitly primary
- **WHEN** a user adds the first link to the links editor
- **THEN** that link SHALL have `is_primary: true` and no star toggle SHALL be shown

#### Scenario: Star toggle appears with multiple links
- **WHEN** a user has two or more links in the editor
- **THEN** each link entry SHALL show a star toggle indicating primary status

#### Scenario: Removing the primary link promotes the first remaining
- **WHEN** a user removes the link currently marked as primary and at least one other link remains
- **THEN** the first remaining link SHALL automatically become primary

#### Scenario: Invalid URL shows inline error
- **WHEN** a user blurs a link input that contains a value that is not a valid URL (after https:// normalisation)
- **THEN** an inline error message SHALL appear beneath that input and the form SHALL not submit until corrected

#### Scenario: https:// prepended on blur
- **WHEN** a user types a URL without a scheme (e.g. `bsky.app/@artist`) and moves focus away from the input
- **THEN** the input value SHALL be updated to `https://bsky.app/@artist`
