## MODIFIED Requirements

### Requirement: List artists
An authenticated user SHALL be able to retrieve artists they own, ordered by name ascending. The endpoint SHALL accept an optional `q` query parameter for freetext name filtering; when absent, all artists owned by the user are returned.

Optional query parameters:
- `q` — freetext substring search against `name` (case-insensitive); absent means no filter

#### Scenario: Successful listing with no filters
- **WHEN** an authenticated user sends `GET /artists`
- **THEN** the system returns 200 with an array of artist objects (may be empty), ordered by name ascending

#### Scenario: Only owner's artists returned
- **WHEN** multiple users have artists and user A sends `GET /artists`
- **THEN** the response contains only user A's artists

#### Scenario: Filter by q returns matching artists
- **WHEN** an authenticated user sends `GET /artists?q=jane`
- **THEN** the system returns only artists whose name contains "jane" case-insensitively, ordered by name ascending
