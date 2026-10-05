## Purpose

The app shell is the authenticated entry point of the Booklet web application, served at `/app`. It acts as the root container for all authenticated views and enforces access control by redirecting unauthenticated users to the home page.

## Requirements

### Requirement: Access control
The app shell at `/app` SHALL only be accessible to authenticated users. Unauthenticated users who navigate to `/app` SHALL be redirected to `/`.

#### Scenario: Unauthenticated access attempt
- **WHEN** an unauthenticated user navigates to `/app`
- **THEN** the app SHALL redirect them to `/`

### Requirement: Authenticated shell render
Authenticated users who navigate to `/app` SHALL be redirected to `/app/characters`. All routes under `/app/*` SHALL render the app shell — a persistent top bar and collapsible sidebar — with the matched page content in the main content area.

#### Scenario: Authenticated user accesses /app
- **WHEN** an authenticated user navigates to `/app`
- **THEN** the app SHALL redirect them to `/app/characters`

#### Scenario: Authenticated user accesses a section route
- **WHEN** an authenticated user navigates to `/app/characters`, `/app/images`, or `/app/artists`
- **THEN** the app shell SHALL render with the top bar and sidebar visible and the corresponding page content in the content area

### Requirement: Top bar
The app shell SHALL display a persistent top bar at all times. The top bar SHALL show the application wordmark ("Booklet") on the left and a user avatar on the right. Clicking the user avatar SHALL log the user out.

#### Scenario: Top bar is visible on all authenticated routes
- **WHEN** an authenticated user is on any `/app/*` route
- **THEN** the top bar SHALL be visible with the wordmark and user avatar

#### Scenario: User logs out via avatar
- **WHEN** an authenticated user clicks the user avatar in the top bar
- **THEN** the app SHALL log the user out and redirect to the home page

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
