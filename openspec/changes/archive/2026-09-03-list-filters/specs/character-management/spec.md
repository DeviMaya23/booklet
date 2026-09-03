## MODIFIED Requirements

### Requirement: List characters
An authenticated user SHALL be able to retrieve all non-deleted characters they own. Each character in the response SHALL include its `folder_ids`. The endpoint SHALL accept an optional `q` query parameter for freetext name filtering; when absent, all characters owned by the user are returned.

Optional query parameters:
- `q` — freetext substring search against `name` (case-insensitive); absent means no filter

#### Scenario: Successful listing with no filters
- **WHEN** an authenticated user sends `GET /characters`
- **THEN** the system returns 200 with an array of character objects (may be empty), each including `folder_ids`

#### Scenario: Only owner's characters returned
- **WHEN** multiple users have characters and user A sends `GET /characters`
- **THEN** the response contains only user A's characters

#### Scenario: Filter by q returns matching characters
- **WHEN** an authenticated user sends `GET /characters?q=aria`
- **THEN** the system returns only characters whose name contains "aria" case-insensitively
