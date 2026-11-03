# Spec Delta

## ADDED Requirements

### Requirement: Commission character assignment updates character recency
Creating or updating a commission with a non-empty `character_ids` SHALL update `last_used_at` on each referenced character (see `character-recency-sort`). If any referenced character is not owned by the user, the request SHALL be rejected and no character's `last_used_at` SHALL change.

#### Scenario: Create bumps assigned characters
- **WHEN** a user creates a commission with `character_ids` containing characters A and B
- **THEN** A's and B's `last_used_at` SHALL be updated to approximately now

#### Scenario: Update bumps characters in the request
- **WHEN** a user updates a commission with a non-empty `character_ids`
- **THEN** every character in the request SHALL have `last_used_at` updated to approximately now

#### Scenario: Unowned character blocks the bump
- **WHEN** a user creates or updates a commission with a `character_ids` entry belonging to another user
- **THEN** the system SHALL reject the request and no character's `last_used_at` SHALL change
