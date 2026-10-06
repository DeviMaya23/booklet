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
The app shell SHALL display a fixed (non-collapsible) sidebar. The sidebar SHALL contain an Inbox count card and two grouped sets of navigation links as defined in the `web-nav-sidebar` spec. The active section SHALL be visually indicated per the active row style defined in `web-nav-sidebar`.

#### Scenario: Sidebar shows inbox card and grouped nav items
- **WHEN** an authenticated user views the app shell
- **THEN** the sidebar SHALL display the Inbox card above two labeled groups (LIBRARY and MANAGE), each containing their respective nav items

#### Scenario: Active nav item is highlighted
- **WHEN** an authenticated user is on a section route (e.g. `/app/artpieces`)
- **THEN** the corresponding nav row SHALL be visually active and all others SHALL not

#### Scenario: Sidebar is always visible and not collapsible
- **WHEN** an authenticated user views any `/app/*` route
- **THEN** the sidebar SHALL be fully visible with no collapse control or rail
