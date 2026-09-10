## Purpose

Defines the requirements for the Artists page (`/app/artists`) in the web application, covering list display, client-side search, copy link, and CRUD operations (create, edit, delete) via a shared artist form modal.

---

## Requirements

### Requirement: Artists list display
The artists page at `/app/artists` SHALL fetch the authenticated user's artists from `GET /artists` and display them as a flat list of rows. Each row SHALL display a bullet point, the artist's `name`, a copy-link button, and an edit button.

#### Scenario: Artists load and display
- **WHEN** an authenticated user navigates to `/app/artists`
- **THEN** the page SHALL fetch `GET /artists` and render one row per artist

#### Scenario: Empty list
- **WHEN** the user has no artists
- **THEN** the page SHALL render an empty list with no rows

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

### Requirement: Copy artist link
Each artist row SHALL display a copy-link button. When `artist_link` is non-null and non-empty, clicking the button SHALL copy the value to the clipboard and show a success toast. When `artist_link` is null or empty, the button SHALL be rendered but disabled.

#### Scenario: Copy link when artist_link is set
- **WHEN** a user clicks the copy-link button on a row where `artist_link` is non-null and non-empty
- **THEN** the value of `artist_link` SHALL be copied to the clipboard and a success toast SHALL appear

#### Scenario: Copy link button disabled when no link
- **WHEN** an artist has a null or empty `artist_link`
- **THEN** the copy-link button SHALL be visible but disabled and non-interactive

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
Each artist row SHALL display an edit button. Clicking it SHALL open the artist form modal in edit mode, pre-populated with the artist's current `name`, `artist_link`, and `notes`. Submitting SHALL call `PUT /artists/:id` with all three fields, close the modal, refresh the list, and show a success toast. If the backend returns 409, the modal SHALL remain open and show an error toast.

Empty optional fields (`link`, `blurb`) SHALL be submitted as `null`, not as empty strings, so that clearing a field wipes its stored value.

#### Scenario: Edit button opens edit modal pre-populated
- **WHEN** a user clicks the edit button on an artist row
- **THEN** the artist form modal SHALL open in edit mode with `name`, `artist_link`, and `notes` pre-filled from the artist's current values

#### Scenario: Successful edit
- **WHEN** a user submits the edit modal with valid changes
- **THEN** the app SHALL call `PUT /artists/:id` with all three fields, close the modal, refetch the artist list, and show a success toast

#### Scenario: Clearing a previously set field
- **WHEN** a user clears the link or blurb field and submits
- **THEN** the app SHALL send `null` for that field, and the stored value SHALL be cleared

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
The artist form modal SHALL contain three fields: `name` (required text input), `link` (optional text input mapped to `artist_link`), and `blurb` (optional textarea mapped to `notes`). The modal SHALL display Save and Cancel buttons. Submitting with an empty `name` SHALL be prevented client-side.

#### Scenario: Submit blocked when name is empty
- **WHEN** a user attempts to submit the modal with an empty name field
- **THEN** the submission SHALL be blocked client-side and no API call SHALL be made

#### Scenario: Optional fields may be left empty
- **WHEN** a user submits the modal with only `name` filled and `link`/`blurb` left empty
- **THEN** the submission SHALL proceed normally
