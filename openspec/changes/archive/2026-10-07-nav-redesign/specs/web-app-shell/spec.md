## MODIFIED Requirements

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

## REMOVED Requirements

### Requirement: Sidebar can be collapsed
**Reason**: The redesigned sidebar is always expanded; no collapse control is shown.
**Migration**: Remove `collapsible="icon"` prop from `SidebarProvider` and remove `<SidebarRail />` from `AppSidebar`.
