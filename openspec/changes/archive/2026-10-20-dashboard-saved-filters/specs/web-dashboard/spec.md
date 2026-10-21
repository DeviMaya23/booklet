# Spec Delta

## ADDED Requirements

### Requirement: Saved filters section presence
The dashboard SHALL include the Saved Filters section above the Recent artpieces section. The full behavior of that section is defined in the `web-dashboard-saved-filters` capability.

#### Scenario: Saved filters section appears above recent artpieces
- **WHEN** an authenticated user views the dashboard
- **THEN** the Saved Filters section SHALL be rendered above the Recent artpieces section

---

## MODIFIED Requirements

### Requirement: Recent artpieces row
The dashboard SHALL display a "Recent artpieces" section with a thumbnail row of the five most recently created artpieces (sorted by `created_at` descending). Each thumbnail tile SHALL display the artpiece image (or a placeholder if no thumbnail), the artpiece title, and the artist name if set. Each tile SHALL be clickable and navigate to `/app/artpieces/:id` for that artpiece. A "View all" link SHALL navigate to `/app/artpieces`.

#### Scenario: Recent artpieces are displayed
- **WHEN** the dashboard loads and recent_artpieces is non-empty
- **THEN** up to five artpiece thumbnail tiles SHALL be displayed in a horizontal row

#### Scenario: Tile click navigates to artpiece detail
- **WHEN** an authenticated user clicks a recent artpiece tile
- **THEN** the app SHALL navigate to `/app/artpieces/:id` for that artpiece

#### Scenario: View all navigates to artpieces gallery
- **WHEN** an authenticated user clicks "View all" in the Recent artpieces section
- **THEN** the app SHALL navigate to `/app/artpieces`
