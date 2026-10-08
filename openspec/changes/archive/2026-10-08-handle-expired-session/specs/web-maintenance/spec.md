## MODIFIED Requirements

### Requirement: Maintenance flag detection
When any API response carries the header `X-Booklet-Maintenance: true`, the app SHALL enter maintenance mode immediately, without requiring a page reload. `apiFetch` SHALL also check the response status and call `setSessionExpired` when it receives a 401 — both side-effects occur in the same response-handling pass.

#### Scenario: API response signals maintenance
- **WHEN** `apiFetch` receives a response with `X-Booklet-Maintenance: true`
- **THEN** `maintenanceStore` SHALL set the maintenance flag to `true` synchronously

#### Scenario: API response clears maintenance
- **WHEN** `apiFetch` receives a response without the maintenance header (or with value `false`)
- **THEN** `maintenanceStore` SHALL set the maintenance flag to `false`

#### Scenario: API response is 401
- **WHEN** `apiFetch` receives a response with status 401
- **THEN** `sessionExpiredStore` SHALL set the session-expired flag to `true`
- **AND** the maintenance header check SHALL still execute on the same response
