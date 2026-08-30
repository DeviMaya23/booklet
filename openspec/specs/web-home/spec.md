## Purpose

The home page is the public-facing entry point of the Booklet web application, served at `/`. It presents unauthenticated users with a login call-to-action and redirects already-authenticated users directly to the app shell.

## Requirements

### Requirement: Login trigger
The home page SHALL display a login button that, when clicked, initiates Kinde's OAuth login redirect. The button SHALL be the only interactive call-to-action on the page.

#### Scenario: User clicks login button
- **WHEN** an unauthenticated user clicks the login button on the home page
- **THEN** the app SHALL redirect the browser to Kinde's authorization endpoint

### Requirement: Authenticated user redirect
If a user who is already authenticated navigates to `/`, the app SHALL redirect them to `/app` without rendering the home page.

#### Scenario: Authenticated user visits home
- **WHEN** an authenticated user navigates to `/`
- **THEN** the app SHALL redirect them to `/app`

#### Scenario: Unauthenticated user visits home
- **WHEN** an unauthenticated user navigates to `/`
- **THEN** the home page SHALL render with the login button visible
