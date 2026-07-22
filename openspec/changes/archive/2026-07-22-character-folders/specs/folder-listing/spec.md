## ADDED Requirements

### Requirement: List available folders from Bookleaf
An authenticated user SHALL be able to retrieve their available public folders from the Bookleaf service via `GET /folders`. The system SHALL proxy the Bookleaf internal API response without transformation.

The Bookleaf API called is: `GET <BOOKLEAF_HOST>/internal/users/:user_id/public-folders`

The user ID sent to Bookleaf SHALL be the authenticated user's ID from the JWT token.

The response shape from Bookleaf is:
```json
{
  "folder_list": [
    {
      "folder_id": "<uuid>",
      "token": "<string>",
      "folder_name": "<string>"
    }
  ]
}
```

This response is returned as-is to the client.

#### Scenario: Successful folder listing
- **WHEN** an authenticated user sends `GET /folders`
- **THEN** the system returns 200 with the folder list from Bookleaf, including `folder_id`, `token`, and `folder_name` for each folder

#### Scenario: Bookleaf API is unavailable
- **WHEN** the Bookleaf API returns an error or is unreachable
- **THEN** the system returns 500

---

### Requirement: Bookleaf HTTP client
The system SHALL maintain a dedicated HTTP client for the Bookleaf internal API in `internal/bookleaf/`. The client SHALL be initialised with the base URL from the `BOOKLEAF_HOST` environment variable and the internal secret from the `BOOKLEAF_INTERNAL_SECRET` environment variable, both of which are required configuration values.

All requests to the Bookleaf internal API SHALL include the header `X-Bookleaf-Internal-Secret` with the value of `BOOKLEAF_INTERNAL_SECRET`.

#### Scenario: Missing BOOKLEAF_HOST prevents startup
- **WHEN** the server starts without `BOOKLEAF_HOST` set
- **THEN** the server fails to start with a descriptive error

#### Scenario: Missing BOOKLEAF_INTERNAL_SECRET prevents startup
- **WHEN** the server starts without `BOOKLEAF_INTERNAL_SECRET` set
- **THEN** the server fails to start with a descriptive error
