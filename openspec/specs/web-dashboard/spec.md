## Purpose

<!-- TBD: Brief description of the web-dashboard capability -->

## Requirements

### Requirement: Dashboard page route
The app SHALL expose a dashboard page at `/app/dashboard`. The route at `/app` SHALL redirect to `/app/dashboard` instead of `/app/characters`.

#### Scenario: Navigating to /app redirects to dashboard
- **WHEN** an authenticated user navigates to `/app`
- **THEN** the app SHALL redirect to `/app/dashboard`

#### Scenario: Dashboard page renders at /app/dashboard
- **WHEN** an authenticated user navigates to `/app/dashboard`
- **THEN** the dashboard page SHALL render

### Requirement: Dashboard header actions
The dashboard page SHALL display an "Upload" button and a "+ New Commission" button in the top-right of the header. The Upload button SHALL open the existing Add Files modal. The "+ New Commission" button SHALL be rendered but non-functional (no-op) in this version.

#### Scenario: Upload button opens Add Files modal
- **WHEN** an authenticated user clicks the "Upload" button on the dashboard
- **THEN** the Add Files modal SHALL open

#### Scenario: New Commission button is visible but inactive
- **WHEN** an authenticated user views the dashboard
- **THEN** the "+ New Commission" button SHALL be visible but clicking it SHALL have no effect

### Requirement: Recent artpieces row
The dashboard SHALL display a "Recent artpieces" section with a thumbnail row of the five most recently created artpieces (sorted by `created_at` descending). Each thumbnail tile SHALL display the artpiece image (or a placeholder if no thumbnail), the artpiece title, and the artist name if set. A "View all" link SHALL navigate to `/app/artpieces`.

#### Scenario: Recent artpieces are displayed
- **WHEN** the dashboard loads and recent_artpieces is non-empty
- **THEN** up to five artpiece thumbnail tiles SHALL be displayed in a horizontal row

#### Scenario: View all navigates to artpieces gallery
- **WHEN** an authenticated user clicks "View all" in the Recent artpieces section
- **THEN** the app SHALL navigate to `/app/artpieces`

### Requirement: In Progress panel
The dashboard SHALL display an "In progress" panel containing all commissions with status `waitlist` or `wip`. Each row SHALL display: title, artist name (linked to artist page if artist is set, same affordance as the commission list), status chip, paid state, days elapsed since creation, and last-contacted timestamp. A "View all" link SHALL navigate to `/app/commissions`.

#### Scenario: In progress commissions are listed
- **WHEN** the dashboard loads and in_progress is non-empty
- **THEN** each commission row SHALL display title, artist, status, paid state, days elapsed, and last-contacted info

#### Scenario: View all navigates to commissions list
- **WHEN** an authenticated user clicks "View all" in the In Progress panel
- **THEN** the app SHALL navigate to `/app/commissions`

### Requirement: In Progress stale contact highlighting
Commission rows where `last_contacted_at` is more than 14 days ago SHALL render the last-contact text in a distinct warning color (orange). Commission rows where `last_contacted_at` is null SHALL render the contact field in muted color without warning styling. The 14-day threshold is a hardcoded constant.

#### Scenario: Stale contact renders in orange
- **WHEN** a commission's last_contacted_at is more than 14 days ago
- **THEN** the last-contact text SHALL render in orange

#### Scenario: Null last_contacted_at renders in muted color without orange
- **WHEN** a commission's last_contacted_at is null
- **THEN** the last-contact field SHALL render "not contacted yet" in muted color with no orange styling

### Requirement: In Progress mark-contacted action
Each In Progress commission row SHALL display a "mark contacted now" button (refresh icon). Clicking it SHALL stamp `last_contacted_at` to the current timestamp via `PATCH /commissions/:id`. The dashboard SHALL reflect the updated timestamp immediately after the mutation succeeds.

#### Scenario: Mark contacted stamps the timestamp
- **WHEN** an authenticated user clicks the mark-contacted button on a commission row
- **THEN** the system SHALL send PATCH /commissions/:id with the current timestamp as lastContactedAt, and the displayed timestamp SHALL update

### Requirement: Housekeeping panel
The dashboard SHALL display a "Housekeeping" panel with four categories, each showing a short list of items linking out to the relevant item. The categories are: "No artist" (commissions), "No artist" (artpieces), "Done, no artpieces" (commissions), and "No files" (artpieces). Each item SHALL be a link that navigates to the corresponding detail page.

#### Scenario: Housekeeping items are displayed by category
- **WHEN** the dashboard loads and one or more housekeeping arrays are non-empty
- **THEN** each non-empty category SHALL render its label and a list of linked items

### Requirement: Housekeeping all-clear empty state
When all four housekeeping arrays are empty, the Housekeeping panel SHALL display a single reassuring message ("All tidy. Nothing needs fixing.") instead of four empty list headers. The panel header SHALL still be visible.

#### Scenario: All-clear state shows single message
- **WHEN** all four housekeeping arrays are empty
- **THEN** the Housekeeping panel SHALL display "All tidy. Nothing needs fixing." with no category list headers visible

### Requirement: Dashboard data loading
The dashboard page SHALL fetch data from `GET /dashboard` on every page load. No client-side caching strategy is applied beyond React Query's default behavior. Loading and error states SHALL be handled gracefully.

#### Scenario: Dashboard shows loading state while fetching
- **WHEN** the dashboard page is mounted and the request is in flight
- **THEN** a loading indicator SHALL be displayed

#### Scenario: Dashboard renders data on successful fetch
- **WHEN** the GET /dashboard request succeeds
- **THEN** all sections SHALL populate with the returned data
