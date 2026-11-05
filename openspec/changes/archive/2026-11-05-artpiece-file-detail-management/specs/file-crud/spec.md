# Spec Delta

## MODIFIED Requirements

### Requirement: Update file notes
The system SHALL allow a user to update the `name`, `notes`, and/or `is_visible` fields on a file they own via a single atomic request. `name` and `notes` are optional and nullable; `is_visible` is an optional boolean.

#### Scenario: Successful update
- **WHEN** a user sends `PUT /files/:id` with a body containing `name`, `notes`, `is_visible`, or any combination of them
- **THEN** the system updates the specified fields in a single database write and returns 200 with the updated file including the new values

#### Scenario: Update is_visible only
- **WHEN** a user sends `PUT /files/:id` with a body containing only `is_visible`
- **THEN** the system updates `is_visible` and leaves `name` and `notes` unchanged, returning 200 with the updated file

#### Scenario: Update is_visible together with one of name or notes
- **WHEN** a user sends `PUT /files/:id` with `is_visible` and a non-null `name` but no `notes`
- **THEN** the system updates `is_visible` and `name` and leaves `notes` unchanged, returning 200 with the updated file; the same holds with `notes` provided and `name` omitted

#### Scenario: Omitted or null name and notes without is_visible
- **WHEN** a user sends `PUT /files/:id` without `is_visible` and `name` or `notes` is omitted or null
- **THEN** the system sets that field to null (existing behavior); a null `name` or `notes` only preserves the stored value when `is_visible` is present in the same request

#### Scenario: File not found or not owned
- **WHEN** a user sends `PUT /files/:id` for a file that does not exist or belongs to another user
- **THEN** the system returns 404

---

### Requirement: File response includes name
The system SHALL include `name` (string or null) and `is_visible` (boolean) fields in every file response — GET, list, update, and complete-upload responses.

#### Scenario: File with a name
- **WHEN** a file has a non-null `name`
- **THEN** the response includes `"name": "<value>"`

#### Scenario: File without a name
- **WHEN** a file has a null `name`
- **THEN** the response includes `"name": null`

#### Scenario: File visibility in response
- **WHEN** a file has `is_visible` set to true or false
- **THEN** the response includes `"is_visible": true` or `"is_visible": false` matching the stored value
