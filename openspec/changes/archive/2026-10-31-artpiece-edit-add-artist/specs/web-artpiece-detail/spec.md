## ADDED Requirements

### Requirement: Edit mode — add artist from artist field
The detail view in edit mode SHALL show a "+" button next to the Artist combobox that opens the existing artist-creation modal, so the user can create an artist without leaving the artpiece.

#### Scenario: Add button opens artist-creation modal
- **WHEN** the detail view is in edit mode and the user clicks the "+" button next to the Artist combobox
- **THEN** the artist-creation modal SHALL open and the artpiece edit form SHALL keep its current unsaved values

#### Scenario: Created artist is selected
- **WHEN** the user creates an artist through the modal opened from the detail view's edit mode
- **THEN** the modal SHALL close, the new artist SHALL be selected in the Artist combobox, and it SHALL be available in the combobox's list of artists

#### Scenario: Created artist is a pending change
- **WHEN** a newly created artist is selected in edit mode and the user has not yet clicked Save
- **THEN** the unsaved changes caption SHALL include "artist changed" and the artpiece SHALL NOT yet be updated; clicking Save SHALL call `PATCH /artpieces/:id` with the new artist's id

#### Scenario: Cancelling artist creation leaves selection unchanged
- **WHEN** the user closes the artist-creation modal without creating an artist
- **THEN** the Artist combobox SHALL keep its previously selected value

#### Scenario: Add button disabled while saving
- **WHEN** the artpiece is being saved
- **THEN** the "+" button SHALL be disabled, as the Artist combobox is
