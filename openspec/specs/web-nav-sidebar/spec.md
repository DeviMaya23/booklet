## Purpose

<!-- TBD: Brief description of the web-nav-sidebar capability -->

## Requirements

### Requirement: Grouped navigation
The sidebar SHALL display nav items in two labeled groups. The first group, labelled "LIBRARY", SHALL contain Artpieces and Characters. The second group, labelled "MANAGE", SHALL contain Commissions and Artists. Group labels SHALL render as small muted uppercase text. There SHALL be no divider lines between groups.

#### Scenario: Library group is rendered
- **WHEN** an authenticated user views the app shell
- **THEN** the sidebar SHALL display a "LIBRARY" group heading followed by Artpieces and Characters nav items, in that order

#### Scenario: Manage group is rendered
- **WHEN** an authenticated user views the app shell
- **THEN** the sidebar SHALL display a "MANAGE" group heading followed by Commissions and Artists nav items, in that order

### Requirement: Nav row icons
Every nav row (Artpieces, Characters, Commissions, Artists) SHALL display a Lucide icon to the left of its label. Icons SHALL be 18px with strokeWidth 1.75. The icon assignments are: Artpieces → `images`, Characters → `users-round`, Commissions → `receipt-text`, Artists → `palette`.

#### Scenario: Nav row shows icon and label
- **WHEN** an authenticated user views the sidebar
- **THEN** each nav row SHALL render its assigned icon alongside its label

### Requirement: Active nav row style
The currently active nav row SHALL render with a soft rounded fill in a slightly deeper tone than the default row background, and its label SHALL use a heavier font weight (semibold). Inactive rows SHALL use normal weight with no fill.

#### Scenario: Active route highlights the correct row
- **WHEN** an authenticated user is on `/app/artpieces`
- **THEN** the Artpieces row SHALL render with the active fill and semibold label, and all other rows SHALL render without fill

### Requirement: Inbox count card
The sidebar SHALL display an Inbox card above the navigation groups. The card SHALL be a fixed 56px tall. It SHALL contain: a rounded icon tile (light background) with the `inbox` Lucide icon, a bold "Inbox" title, a contextual subtitle, and a count badge on the right edge. The count badge SHALL display the number of unassigned files. When the count is greater than zero, the subtitle SHALL read "X to sort" (where X is the count) and the badge SHALL be visible with a dark background. When the count is zero, the subtitle SHALL read "All sorted" and the badge SHALL not be rendered. The card height SHALL remain 56px regardless of count state so no layout shift occurs.

#### Scenario: Inbox card shows count when files are unassigned
- **WHEN** there are unassigned files (count > 0)
- **THEN** the Inbox card SHALL display the subtitle "X to sort", a visible dark badge with the count, and the icon tile

#### Scenario: Inbox card shows zero state
- **WHEN** there are no unassigned files (count = 0)
- **THEN** the Inbox card SHALL display the subtitle "All sorted", no badge, and maintain its 56px height

#### Scenario: Inbox card links to the Inbox page
- **WHEN** an authenticated user clicks the Inbox card
- **THEN** the app SHALL navigate to `/app/files`
