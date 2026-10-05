## Purpose

Defines the behaviour of the Images page (`/app/images`) in the web application: listing a user's images, searching them client-side, and deleting individual images.

## Requirements

### Requirement: New image button
The images page SHALL display a `+ New` button. Clicking the button SHALL open the `ImageFormModal` in create mode.

#### Scenario: New button opens create modal
- **WHEN** an authenticated user clicks the `+ New` button on `/app/images`
- **THEN** the `ImageFormModal` SHALL open in create mode

