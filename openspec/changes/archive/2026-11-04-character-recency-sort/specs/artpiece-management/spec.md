# Spec Delta

## ADDED Requirements

### Requirement: Artpiece character assignment updates character recency
Creating or updating an artpiece with a non-empty `character_ids` SHALL update `last_used_at` on each referenced character (see `character-recency-sort`). On update, character ownership SHALL be validated before any recency bump; if any character is not owned by the user, the request SHALL be rejected with 422 and no character's `last_used_at` SHALL change.

#### Scenario: Create bumps assigned characters
- **WHEN** a user creates an artpiece with `character_ids` containing characters A and B
- **THEN** A's and B's `last_used_at` SHALL be updated to approximately now

#### Scenario: Update bumps characters in the request
- **WHEN** a user updates an artpiece with `character_ids` containing characters A and C, where A was already assigned
- **THEN** A's and C's `last_used_at` SHALL be updated to approximately now

#### Scenario: Update with empty character set bumps nothing
- **WHEN** a user updates an artpiece with an empty `character_ids`
- **THEN** no character's `last_used_at` SHALL change

#### Scenario: Unowned character on update blocks the bump
- **WHEN** a user updates an artpiece with a `character_ids` entry belonging to another user
- **THEN** the system SHALL return HTTP 422 and no character's `last_used_at` SHALL change
