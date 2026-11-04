# tag-management Specification

## Purpose

Defines user-owned, case-insensitive tags: their create, list, rename and delete lifecycle, and the `last_used_at` recency tracking that lets pickers surface recently used tags first. Tags are an independent entity; what they are attached to is specified by the capability that attaches them.

## Requirements

### Requirement: Create tag
The system SHALL allow an authenticated user to create a tag with a name. The name SHALL be trimmed of leading and trailing whitespace and stored in lowercase. The response SHALL contain the tag's `id` and `name`.

#### Scenario: Successful creation
- **WHEN** a user posts `{"name": "WIP"}` to `POST /tags`
- **THEN** the system creates a tag owned by the user with name `wip` and returns it with HTTP 201

#### Scenario: Name is trimmed
- **WHEN** a user posts `{"name": "  sketch  "}`
- **THEN** the stored name is `sketch`

#### Scenario: Missing or blank name rejected
- **WHEN** a user posts an empty body, an empty name, or a name containing only whitespace
- **THEN** the system returns HTTP 422 and creates no tag

#### Scenario: Duplicate name rejected case-insensitively
- **WHEN** a user who already has a tag `wip` posts `{"name": "WIP"}`
- **THEN** the system returns HTTP 409 and creates no tag

#### Scenario: Same name allowed for different users
- **WHEN** user B posts a tag named `wip` and user A already has a tag named `wip`
- **THEN** the system creates user B's tag with HTTP 201

#### Scenario: Malformed body
- **WHEN** a user posts a body that is not valid JSON
- **THEN** the system returns HTTP 400

---

### Requirement: List tags
The system SHALL return all tags belonging to the authenticated user from `GET /tags`, ordered by `name` ascending by default. Each tag in the response SHALL contain `id` and `name`, and SHALL NOT contain `last_used_at`. There is no endpoint to fetch a single tag by ID.

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

---

### Requirement: List tags sorted by recency on request
`GET /tags` SHALL accept an optional `sort` query parameter. With `sort=recent` the results SHALL be ordered by `last_used_at` descending with null values last, then `name` ascending. Any other `sort` value SHALL be rejected with HTTP 422.

#### Scenario: Most recently used tag first
- **WHEN** a user sends `GET /tags?sort=recent` and tag A was most recently assigned to an artpiece
- **THEN** tag A is returned before the other tags

#### Scenario: Null last_used_at sorts last
- **WHEN** a tag has a null `last_used_at` and the user sends `GET /tags?sort=recent`
- **THEN** it appears after every tag with a non-null `last_used_at`

#### Scenario: Tie broken alphabetically
- **WHEN** two tags have the same `last_used_at` and the user sends `GET /tags?sort=recent`
- **THEN** they are ordered by `name` ascending

#### Scenario: Default order unchanged without sort
- **WHEN** a user sends `GET /tags` without `sort`
- **THEN** tags are ordered by `name` ascending regardless of `last_used_at`

#### Scenario: Unknown sort value rejected
- **WHEN** a user sends `GET /tags?sort=bogus`
- **THEN** the system returns HTTP 422 with a structured validation error

---

### Requirement: Tag last_used_at tracking
Each tag SHALL have a `last_used_at` timestamp set to the current time when the tag is created. The system SHALL set it to the current time for each tag in `tag_ids` when an artpiece is created or updated with a non-empty `tag_ids`, including tags already attached. It SHALL NOT change when a tag is renamed.

#### Scenario: Set on tag creation
- **WHEN** a user creates a tag
- **THEN** the tag's `last_used_at` is approximately the current time

#### Scenario: Bumped on artpiece create with tags
- **WHEN** a user creates an artpiece with a non-empty `tag_ids`
- **THEN** every tag in `tag_ids` has `last_used_at` updated to approximately the current time

#### Scenario: Bumped on artpiece update with tags
- **WHEN** a user updates an artpiece with a non-empty `tag_ids`
- **THEN** every tag in `tag_ids` has `last_used_at` updated, including tags already assigned to that artpiece

#### Scenario: Not bumped when tag_ids is empty or absent
- **WHEN** a user creates or updates an artpiece with an empty or absent `tag_ids`
- **THEN** no tag's `last_used_at` is modified

#### Scenario: Not bumped on rename
- **WHEN** a user renames a tag
- **THEN** that tag's `last_used_at` is unchanged

---

### Requirement: Rename tag
The system SHALL allow an authenticated user to rename a tag they own via `PATCH /tags/:id`. The new name SHALL be trimmed and lowercased as on create. The response SHALL contain the updated tag.

#### Scenario: Successful rename
- **WHEN** a user patches their tag with `{"name": "Finished"}`
- **THEN** the system updates the name to `finished` and returns the tag with HTTP 200

#### Scenario: Rename to a name used by another of the user's tags
- **WHEN** a user renames a tag to a name, ignoring case, that another of their tags already has
- **THEN** the system returns HTTP 409 and leaves both tags unchanged

#### Scenario: Rename to its own name in a different case
- **WHEN** a user patches a tag named `wip` with `{"name": "WIP"}`
- **THEN** the system returns HTTP 200 and the name remains `wip`

#### Scenario: Blank name rejected
- **WHEN** a user patches a tag with an empty or whitespace-only name
- **THEN** the system returns HTTP 422

#### Scenario: Tag not found or belongs to another user
- **WHEN** a user patches a tag that does not exist or belongs to another user
- **THEN** the system returns HTTP 404

#### Scenario: Malformed tag ID
- **WHEN** a user patches `/tags/not-a-uuid`
- **THEN** the system returns HTTP 400

---

### Requirement: Delete tag
The system SHALL allow an authenticated user to delete a tag they own via `DELETE /tags/:id`. Deleting a tag SHALL remove it from every artpiece it was assigned to and SHALL NOT delete or otherwise modify those artpieces.

#### Scenario: Successful deletion
- **WHEN** a user deletes a tag they own
- **THEN** the system returns HTTP 204 and the tag no longer appears in `GET /tags`

#### Scenario: Deleting a tag in use detaches it
- **WHEN** a user deletes a tag assigned to two artpieces
- **THEN** both artpieces still exist and their `tags` no longer include the deleted tag

#### Scenario: Tag not found or belongs to another user
- **WHEN** a user deletes a tag that does not exist or belongs to another user
- **THEN** the system returns HTTP 404 and the other user's tag is unaffected

#### Scenario: Malformed tag ID
- **WHEN** a user sends `DELETE /tags/not-a-uuid`
- **THEN** the system returns HTTP 400
