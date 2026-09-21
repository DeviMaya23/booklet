## MODIFIED Requirements

### Requirement: New image button
The images page SHALL display a `+ New` button. Clicking the button SHALL open the `ImageFormModal` in create mode.

#### Scenario: New button opens create modal
- **WHEN** an authenticated user clicks the `+ New` button on `/app/images`
- **THEN** the `ImageFormModal` SHALL open in create mode
