## MODIFIED Requirements

### Requirement: Inbox count card
The sidebar SHALL display an Inbox card below the Dashboard nav item and above the navigation groups. The card SHALL be a fixed 56px tall. It SHALL contain: a rounded icon tile (light background) with the `inbox` Lucide icon, a bold "Inbox" title, a contextual subtitle, and a count badge on the right edge. The count badge SHALL display the number of unassigned files. When the count is greater than zero, the subtitle SHALL read "X to sort" (where X is the count) and the badge SHALL be visible with a dark background. When the count is zero, the subtitle SHALL read "All sorted" and the badge SHALL not be rendered. The card height SHALL remain 56px regardless of count state so no layout shift occurs.

#### Scenario: Inbox card shows count when files are unassigned
- **WHEN** there are unassigned files (count > 0)
- **THEN** the Inbox card SHALL display the subtitle "X to sort", a visible dark badge with the count, and the icon tile

#### Scenario: Inbox card shows zero state
- **WHEN** there are no unassigned files (count = 0)
- **THEN** the Inbox card SHALL display the subtitle "All sorted", no badge, and maintain its 56px height

#### Scenario: Inbox card links to the Inbox page
- **WHEN** an authenticated user clicks the Inbox card
- **THEN** the app SHALL navigate to `/app/files`

## ADDED Requirements

### Requirement: Dashboard nav item
The sidebar SHALL display a Dashboard nav item above the Inbox card. It SHALL use the `layout-dashboard` Lucide icon. It SHALL link to `/app/dashboard` and follow the same active-state styling as other nav rows (soft rounded fill, semibold label when active).

#### Scenario: Dashboard nav item appears above Inbox
- **WHEN** an authenticated user views the app shell
- **THEN** a "Dashboard" nav item SHALL appear as the first item in the sidebar, above the Inbox card

#### Scenario: Dashboard nav item is active on the dashboard route
- **WHEN** an authenticated user is on `/app/dashboard`
- **THEN** the Dashboard nav item SHALL render with active fill and semibold label

#### Scenario: Dashboard nav item navigates to /app/dashboard
- **WHEN** an authenticated user clicks the Dashboard nav item
- **THEN** the app SHALL navigate to `/app/dashboard`
