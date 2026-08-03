# User Auth Spec

## Purpose

Covers authentication middleware behavior, including user provisioning and access control checks applied before requests reach handlers.

## Requirements

### Requirement: Accounts pending deletion are blocked from authenticated endpoints
The system SHALL reject requests from users whose `account_state` is not `active`. This check SHALL occur after `GetOrProvision` in the auth middleware, before the request reaches any handler. Both `pending_deletion` and `purged` states SHALL result in a 401 response.

#### Scenario: Request from a pending-deletion account
- **WHEN** an authenticated request arrives for a user with `account_state = 'pending_deletion'`
- **THEN** the auth middleware returns 401 and the request does not reach the handler

#### Scenario: Request from a purged (tombstoned) account
- **WHEN** an authenticated request arrives for a user with `account_state = 'purged'`
- **THEN** the auth middleware returns 401 and the request does not reach the handler

#### Scenario: Request from a normal account
- **WHEN** an authenticated request arrives for a user with `account_state = 'active'`
- **THEN** the auth middleware proceeds normally and sets the user ID in context
