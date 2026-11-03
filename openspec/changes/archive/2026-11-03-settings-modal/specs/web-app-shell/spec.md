## MODIFIED Requirements

### Requirement: Top bar
The app shell SHALL display a persistent top bar at all times. The top bar SHALL show the application wordmark ("Booklet") on the left and a user avatar on the right. Clicking the user avatar SHALL open a profile dropdown rather than logging the user out. The dropdown contents are defined in the `web-settings-modal` spec.

#### Scenario: Top bar is visible on all authenticated routes
- **WHEN** an authenticated user is on any `/app/*` route
- **THEN** the top bar SHALL be visible with the wordmark and user avatar

#### Scenario: Clicking the avatar opens the dropdown without logging out
- **WHEN** an authenticated user clicks the user avatar in the top bar
- **THEN** the app SHALL open the profile dropdown
- **AND** the user SHALL remain authenticated and on the current route

#### Scenario: User logs out via avatar
- **WHEN** an authenticated user opens the profile dropdown from the avatar and selects "Sign out"
- **THEN** the app SHALL log the user out and redirect to the home page
