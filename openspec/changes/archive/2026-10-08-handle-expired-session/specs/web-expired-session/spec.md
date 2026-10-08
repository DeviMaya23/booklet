## ADDED Requirements

### Requirement: Session expiry detection
When `apiFetch` receives a 401 response, the app SHALL signal that the session has expired via `sessionExpiredStore`. Signals after the first SHALL be no-ops until the store is reset.

#### Scenario: First 401 response
- **WHEN** `apiFetch` receives a response with status 401
- **THEN** `sessionExpiredStore` SHALL set the session-expired flag to `true`

#### Scenario: Subsequent 401 responses while flag is already set
- **WHEN** `apiFetch` receives a response with status 401 and the session-expired flag is already `true`
- **THEN** `sessionExpiredStore` SHALL NOT fire listeners again (no-op)

### Requirement: Session expiry redirect
When the session-expired flag becomes `true`, the app SHALL redirect the user to `/` and display a toast with the message "Please log back in". The flag SHALL be reset immediately before navigation.

#### Scenario: Session expires while on an authenticated route
- **WHEN** the session-expired flag becomes `true` and the user is on any `/app` route
- **THEN** the app SHALL display a toast with the message "Please log back in"
- **AND** the app SHALL navigate to `/`
- **AND** `sessionExpiredStore` SHALL reset the flag to `false`

#### Scenario: Only one toast fires when multiple requests fail with 401
- **WHEN** multiple concurrent `apiFetch` calls all receive 401 responses
- **THEN** exactly one toast SHALL be shown and exactly one navigation to `/` SHALL occur
