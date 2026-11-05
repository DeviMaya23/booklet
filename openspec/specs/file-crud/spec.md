## Purpose

This capability covers CRUD operations on files owned by a user: fetching a single file, listing files (with optional unassigned filter), updating file notes, and deleting a file along with its R2 storage objects.

## Requirements

### Requirement: Get file by ID
The system SHALL allow a user to fetch a single file by ID, returning the file record with a presigned URL for the file and its thumbnail.

#### Scenario: File exists and belongs to user
- **WHEN** a user requests `GET /files/:id` for a file they own
- **THEN** the system returns 200 with the file's metadata, a presigned file URL, and a presigned thumbnail URL (null if no thumbnail)

#### Scenario: File not found or not owned
- **WHEN** a user requests `GET /files/:id` for a file that does not exist or belongs to another user
- **THEN** the system returns 404

---

### Requirement: List files
The system SHALL allow a user to list all their files, with an optional filter to return only unassigned files (those with no artpiece attached).

#### Scenario: List without filter
- **WHEN** a user requests `GET /files` with no query params
- **THEN** the system returns 200 with all files belonging to the user, each with a presigned file URL and presigned thumbnail URL

#### Scenario: List with unassigned filter
- **WHEN** a user requests `GET /files?unassigned=true`
- **THEN** the system returns 200 with only the files belonging to the user that have no artpiece attached (`artpiece_id IS NULL`)

---

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

---

### Requirement: Delete file
The system SHALL allow a user to delete a file they own. Deletion removes the file record from the database and the corresponding object from R2 storage.

#### Scenario: Successful delete
- **WHEN** a user sends `DELETE /files/:id` for a file they own
- **THEN** the system removes the file record and deletes the R2 object, returning 204

#### Scenario: File not found or not owned
- **WHEN** a user sends `DELETE /files/:id` for a file that does not exist or belongs to another user
- **THEN** the system returns 404

#### Scenario: R2 delete fails
- **WHEN** a user sends `DELETE /files/:id` and the R2 delete call fails
- **THEN** the system still removes the file record, logs the R2 error, and returns 204

### Requirement: File response includes extension
The system SHALL include an `extension` field (string or null) in every file response — GET, list, update, and complete-upload responses.

#### Scenario: File with a known extension
- **WHEN** a file has a non-null `extension` (e.g., `"clip"`)
- **THEN** every file response for that file includes `"extension": "clip"`

#### Scenario: File without a stored extension
- **WHEN** a file has a null `extension` (uploaded before this change or with no extension provided)
- **THEN** every file response for that file includes `"extension": null`
