## MODIFIED Requirements

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
