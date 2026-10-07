## Purpose

TBD

## Requirements

### Requirement: Artpiece detail is accessible via URL
The system SHALL provide a route `/app/artpieces/:id` that renders the artpiece detail view for the given artpiece ID. The Artpieces sidebar item SHALL remain active (highlighted) when on this route.

#### Scenario: Navigate to artpiece detail via URL
- **WHEN** a user navigates to `/app/artpieces/:id` with a valid artpiece ID
- **THEN** the artpiece detail view SHALL render for that artpiece

#### Scenario: Sidebar highlights Artpieces on detail route
- **WHEN** a user is on `/app/artpieces/:id`
- **THEN** the Artpieces sidebar item SHALL appear active/highlighted

#### Scenario: Double-clicking a gallery card navigates to the detail route
- **WHEN** an authenticated user double-clicks an artpiece card in the gallery
- **THEN** the browser SHALL navigate to `/app/artpieces/:id` for that artpiece

#### Scenario: Back chevron returns to previous page
- **WHEN** a user clicks the back chevron in the artpiece detail view
- **THEN** the browser SHALL navigate back one step in history (`navigate(-1)`), returning the user to wherever they came from (artpieces gallery, commissions page, etc.)
