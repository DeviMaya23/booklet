## MODIFIED Requirements

### Requirement: Sidebar navigation
The app shell SHALL display a collapsible sidebar with navigation links to Inbox, Artpieces, Characters, Images, and Artists, in that order. The active section SHALL be visually indicated. The sidebar SHALL be collapsible and its state SHALL persist across navigation within the session.

#### Scenario: Sidebar shows all nav items
- **WHEN** an authenticated user views the app shell
- **THEN** the sidebar SHALL display links for Inbox, Artpieces, Characters, Images, and Artists in that order

#### Scenario: Artpieces nav item links to gallery
- **WHEN** an authenticated user clicks the Artpieces nav item
- **THEN** the app SHALL navigate to `/app/artpieces`

#### Scenario: Active nav item is highlighted
- **WHEN** an authenticated user is on `/app/characters`
- **THEN** the Characters nav item SHALL be visually active

#### Scenario: Sidebar can be collapsed
- **WHEN** an authenticated user toggles the sidebar collapse control
- **THEN** the sidebar SHALL collapse and the main content area SHALL expand to fill the available space
