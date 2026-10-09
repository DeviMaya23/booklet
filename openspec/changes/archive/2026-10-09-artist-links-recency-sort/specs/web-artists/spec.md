## MODIFIED Requirements

### Requirement: Copy artist link
Each artist row SHALL display a copy-link button. When the artist has a link with `is_primary = true`, clicking the button SHALL copy that link's URL to the clipboard and show a success toast. When the artist has no primary link, the button SHALL be rendered but disabled.

#### Scenario: Copy link when primary link is set
- **WHEN** a user clicks the copy-link button on a row where the artist has a primary link
- **THEN** the primary link's URL SHALL be copied to the clipboard and a success toast SHALL appear

#### Scenario: Copy link button disabled when no primary link
- **WHEN** an artist has no link with `is_primary = true`
- **THEN** the copy-link button SHALL be visible but disabled and non-interactive

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
