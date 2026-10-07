## MODIFIED Requirements

### Requirement: Open artpiece detail view
The system SHALL render the artpiece detail view at the route `/app/artpieces/:id`. Navigating to this route displays the detail view for the given artpiece. The artpieces gallery is shown at `/app/artpieces` (no `:id`).

#### Scenario: Open via double-click
- **WHEN** an authenticated user double-clicks an artpiece card in the gallery
- **THEN** the browser SHALL navigate to `/app/artpieces/:id` for that artpiece, displaying the detail view

#### Scenario: Close via back chevron
- **WHEN** the user clicks the back chevron in the detail view
- **THEN** the browser SHALL navigate back one step in history (`navigate(-1)`), returning the user to the previous page

#### Scenario: Back chevron label
- **WHEN** the detail view is open
- **THEN** the back chevron SHALL display the label "← Back" (not "← Artpieces")

#### Scenario: Deletion closes detail and returns to gallery
- **WHEN** the user confirms artpiece deletion from the detail view
- **THEN** the artpiece is deleted and the browser SHALL navigate to `/app/artpieces`

## REMOVED Requirements

### Requirement: Open artpiece detail view (local state)
**Reason:** Replaced by URL-based routing (`/app/artpieces/:id`). The previous behaviour — where `ArtpiecesPage` held a `selectedArtpieceId` state variable and conditionally rendered `ArtpieceDetailView` in place of the gallery — is replaced by React Router's `useParams`. No route change occurred in the old implementation; now it does.
**Migration:** `ArtpiecesPage` reads the artpiece ID from `useParams()` instead of local state. Opening an artpiece calls `navigate('/app/artpieces/' + id)`. Closing calls `navigate(-1)`. `onDeleted` navigates to `/app/artpieces`.
