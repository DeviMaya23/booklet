## Purpose

The app shell is the authenticated entry point of the Booklet web application, served at `/app`. It acts as the root container for all authenticated views and enforces access control by redirecting unauthenticated users to the home page.

## Requirements

### Requirement: Access control
The app shell at `/app` SHALL only be accessible to authenticated users. Unauthenticated users who navigate to `/app` SHALL be redirected to `/`.

#### Scenario: Unauthenticated access attempt
- **WHEN** an unauthenticated user navigates to `/app`
- **THEN** the app SHALL redirect them to `/`

### Requirement: Authenticated shell render
Authenticated users who navigate to `/app` SHALL see the app shell page. The page is an empty placeholder at this stage; its only requirement is that it renders without error.

#### Scenario: Authenticated user accesses the shell
- **WHEN** an authenticated user navigates to `/app`
- **THEN** the app shell page SHALL render successfully
