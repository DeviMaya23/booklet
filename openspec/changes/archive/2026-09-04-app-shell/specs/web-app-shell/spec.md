## MODIFIED Requirements

### Requirement: Authenticated shell render
Authenticated users who navigate to `/app` SHALL be redirected to `/app/characters`. All routes under `/app/*` SHALL render the app shell — a persistent top bar and collapsible sidebar — with the matched page content in the main content area.

#### Scenario: Authenticated user accesses /app
- **WHEN** an authenticated user navigates to `/app`
- **THEN** the app SHALL redirect them to `/app/characters`

#### Scenario: Authenticated user accesses a section route
- **WHEN** an authenticated user navigates to `/app/characters`, `/app/images`, or `/app/artists`
- **THEN** the app shell SHALL render with the top bar and sidebar visible and the corresponding page content in the content area

## ADDED Requirements

### Requirement: Top bar
The app shell SHALL display a persistent top bar at all times. The top bar SHALL show the application wordmark ("Booklet") on the left and a user avatar on the right. Clicking the user avatar SHALL log the user out.

#### Scenario: Top bar is visible on all authenticated routes
- **WHEN** an authenticated user is on any `/app/*` route
- **THEN** the top bar SHALL be visible with the wordmark and user avatar

#### Scenario: User logs out via avatar
- **WHEN** an authenticated user clicks the user avatar in the top bar
- **THEN** the app SHALL log the user out and redirect to the home page

### Requirement: Sidebar navigation
The app shell SHALL display a collapsible sidebar with navigation links to Characters, Images, and Artists. The active section SHALL be visually indicated. The sidebar SHALL be collapsible and its state SHALL persist across navigation within the session.

#### Scenario: Sidebar shows all three nav items
- **WHEN** an authenticated user views the app shell
- **THEN** the sidebar SHALL display links for Characters, Images, and Artists

#### Scenario: Active nav item is highlighted
- **WHEN** an authenticated user is on `/app/characters`
- **THEN** the Characters nav item SHALL be visually active

#### Scenario: Sidebar can be collapsed
- **WHEN** an authenticated user toggles the sidebar collapse control
- **THEN** the sidebar SHALL collapse and the main content area SHALL expand to fill the available space
