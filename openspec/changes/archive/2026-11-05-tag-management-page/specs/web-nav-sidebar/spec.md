# Spec Delta

## MODIFIED Requirements

### Requirement: Grouped navigation
The sidebar SHALL display nav items in two labeled groups. The first group, labelled "LIBRARY", SHALL contain Artpieces and Characters. The second group, labelled "MANAGE", SHALL contain Commissions, Artists and Tags. Group labels SHALL render as small muted uppercase text. There SHALL be no divider lines between groups.

#### Scenario: Library group is rendered
- **WHEN** an authenticated user views the app shell
- **THEN** the sidebar SHALL display a "LIBRARY" group heading followed by Artpieces and Characters nav items, in that order

#### Scenario: Manage group is rendered
- **WHEN** an authenticated user views the app shell
- **THEN** the sidebar SHALL display a "MANAGE" group heading followed by Commissions, Artists and Tags nav items, in that order

#### Scenario: Tags nav item navigates to /app/tags
- **WHEN** an authenticated user clicks the Tags nav item
- **THEN** the app SHALL navigate to `/app/tags` and the Tags row SHALL render with the active fill and semibold label

### Requirement: Nav row icons
Every nav row (Artpieces, Characters, Commissions, Artists, Tags) SHALL display a Lucide icon to the left of its label. Icons SHALL be 18px with strokeWidth 1.75. The icon assignments are: Artpieces → `images`, Characters → `users-round`, Commissions → `receipt-text`, Artists → `palette`, Tags → `tags`.

#### Scenario: Nav row shows icon and label
- **WHEN** an authenticated user views the sidebar
- **THEN** each nav row SHALL render its assigned icon alongside its label
