## ADDED Requirements

### Requirement: Accounts pending deletion are blocked from authenticated endpoints
The system SHALL reject requests from users whose `is_pending_deletion` flag is `true`. This check SHALL occur after `GetOrProvision` in the auth middleware, before the request reaches any handler.

#### Scenario: Request from a pending-deletion account
- **WHEN** an authenticated request arrives for a user with `is_pending_deletion = true`
- **THEN** the auth middleware returns 401 and the request does not reach the handler

#### Scenario: Request from a normal account
- **WHEN** an authenticated request arrives for a user with `is_pending_deletion = false`
- **THEN** the auth middleware proceeds normally and sets the user ID in context
