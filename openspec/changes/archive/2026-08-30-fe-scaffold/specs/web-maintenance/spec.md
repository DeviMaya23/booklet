## ADDED Requirements

### Requirement: Maintenance flag detection
When any API response carries the header `X-Booklet-Maintenance: true`, the app SHALL enter maintenance mode immediately, without requiring a page reload.

#### Scenario: API response signals maintenance
- **WHEN** `apiFetch` receives a response with `X-Booklet-Maintenance: true`
- **THEN** `maintenanceStore` SHALL set the maintenance flag to `true` synchronously

#### Scenario: API response clears maintenance
- **WHEN** `apiFetch` receives a response without the maintenance header (or with value `false`)
- **THEN** `maintenanceStore` SHALL set the maintenance flag to `false`

### Requirement: Maintenance page display
While maintenance mode is active, all authenticated app views SHALL be replaced by the maintenance page. The maintenance page SHALL be shown in place of any routed content under `/app`.

#### Scenario: Authenticated view during maintenance
- **WHEN** maintenance mode is active and an authenticated user is on any `/app` route
- **THEN** the app SHALL display the maintenance page instead of the normal route content

#### Scenario: Authenticated view after maintenance clears
- **WHEN** maintenance mode becomes inactive
- **THEN** the app SHALL resume displaying the normal route content without a page reload

### Requirement: Maintenance bypass
If a bypass token is present in localStorage under the key `booklet-maintenance-bypass`, `apiFetch` SHALL include it as the `X-Booklet-Bypass` request header on every API call, allowing the backend to exempt privileged users from maintenance mode.

#### Scenario: Bypass token is present
- **WHEN** `apiFetch` is called and `booklet-maintenance-bypass` exists in localStorage
- **THEN** the request SHALL include the header `X-Booklet-Bypass: <token>`

#### Scenario: No bypass token
- **WHEN** `apiFetch` is called and `booklet-maintenance-bypass` is absent from localStorage
- **THEN** the request SHALL NOT include the `X-Booklet-Bypass` header
