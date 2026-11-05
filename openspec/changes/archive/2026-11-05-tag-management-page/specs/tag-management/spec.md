# Spec Delta

## MODIFIED Requirements

### Requirement: List tags
The system SHALL return all tags belonging to the authenticated user from `GET /tags`, ordered by `name` ascending by default. Each tag in the response SHALL contain `id` and `name`, and SHALL NOT contain `last_used_at`. Each tag SHALL contain `artpiece_count` only when usage is requested (see "List tags with usage count"). There is no endpoint to fetch a single tag by ID.

#### Scenario: Default order is alphabetical
- **WHEN** a user with tags `sketch`, `wip` and `icon` sends `GET /tags`
- **THEN** the system returns them in the order `icon`, `sketch`, `wip`

#### Scenario: Only the user's own tags
- **WHEN** user A sends `GET /tags` and user B owns other tags
- **THEN** the response contains only user A's tags

#### Scenario: No tags
- **WHEN** a user with no tags sends `GET /tags`
- **THEN** the system returns HTTP 200 with an empty array

#### Scenario: last_used_at not exposed
- **WHEN** any tag endpoint returns a tag
- **THEN** the tag object SHALL NOT contain a `last_used_at` field

#### Scenario: No usage count by default
- **WHEN** a user sends `GET /tags` without `with_usage`
- **THEN** no tag object in the response contains an `artpiece_count` field

## ADDED Requirements

### Requirement: List tags with usage count
`GET /tags` SHALL accept an optional boolean `with_usage` query parameter. With `with_usage=true`, each tag in the response SHALL additionally contain `artpiece_count`, the number of the user's artpieces the tag is assigned to, and SHALL be 0 for an unassigned tag. It SHALL combine with `sort`. A `with_usage` value that is not a boolean SHALL be rejected with HTTP 400.

#### Scenario: Count reflects assignments
- **WHEN** tag `sketch` is assigned to 3 artpieces and a user sends `GET /tags?with_usage=true`
- **THEN** the `sketch` entry contains `"artpiece_count": 3`

#### Scenario: Unused tag reports zero
- **WHEN** a tag is assigned to no artpiece and a user sends `GET /tags?with_usage=true`
- **THEN** that tag's entry contains `"artpiece_count": 0`, not an absent field

#### Scenario: Count is not inflated by other tags or users
- **WHEN** an artpiece has several tags, or user B has artpieces with a same-named tag
- **THEN** each of user A's tags counts each of A's artpieces at most once and never counts B's artpieces

#### Scenario: Combined with recent sort
- **WHEN** a user sends `GET /tags?with_usage=true&sort=recent`
- **THEN** tags are ordered as for `sort=recent` and each contains `artpiece_count`

#### Scenario: Explicit false
- **WHEN** a user sends `GET /tags?with_usage=false`
- **THEN** the response is the same as without the parameter

#### Scenario: Invalid value rejected
- **WHEN** a user sends `GET /tags?with_usage=maybe`
- **THEN** the system returns HTTP 400

### Requirement: Merge tags
The system SHALL allow an authenticated user to merge a tag they own (the source, `:id`) into another they own (the target) via `POST /tags/:id/merge` with body `{"target_id": "<uuid>"}`. Artpieces having the source SHALL end up having the target, and the source SHALL be deleted. The operation SHALL be atomic and return HTTP 204.

#### Scenario: Target name and last_used_at untouched
- **WHEN** the user merges a tag into a target
- **THEN** the target's name and `last_used_at` are unchanged and no artpiece other than its tag assignments is modified

#### Scenario: Successful merge
- **WHEN** artpieces A and B have `sketches`, and artpiece B also has `sketch`, and the user merges `sketches` into `sketch`
- **THEN** the system returns HTTP 204, A and B both have `sketch` exactly once, and `sketches` no longer appears in `GET /tags`

#### Scenario: Target last_used_at unchanged
- **WHEN** the user merges a tag into a target
- **THEN** the target's position under `GET /tags?sort=recent` is determined by the same `last_used_at` it had before the merge

#### Scenario: Source with no artpieces
- **WHEN** the user merges an unassigned source tag into a target
- **THEN** the system returns HTTP 204 and the source is deleted

#### Scenario: Source equals target
- **WHEN** the user posts a `target_id` equal to the source `:id`
- **THEN** the system returns HTTP 422 and nothing changes

#### Scenario: Source or target not found or not owned
- **WHEN** the source or the target does not exist or belongs to another user
- **THEN** the system returns HTTP 404 and neither tag nor any assignment is changed

#### Scenario: Malformed source ID
- **WHEN** a user sends `POST /tags/not-a-uuid/merge`
- **THEN** the system returns HTTP 400

#### Scenario: Missing or malformed target_id
- **WHEN** the body is not valid JSON, or `target_id` is not a valid UUID
- **THEN** the system returns HTTP 400; if `target_id` is absent the system returns HTTP 422

#### Scenario: Failure leaves data intact
- **WHEN** any step of the merge fails
- **THEN** the source tag and its assignments, and the target's assignments, remain exactly as before
