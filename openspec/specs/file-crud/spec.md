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
The system SHALL allow a user to update the `name` and/or `notes` fields on a file they own via a single atomic request. Both fields are optional and nullable.

#### Scenario: Successful update
- **WHEN** a user sends `PUT /files/:id` with a body containing `name`, `notes`, or both (each nullable)
- **THEN** the system updates the specified fields in a single database write and returns 200 with the updated file including the new `name` and `notes` values

#### Scenario: File not found or not owned
- **WHEN** a user sends `PUT /files/:id` for a file that does not exist or belongs to another user
- **THEN** the system returns 404

---

### Requirement: File response includes name
The system SHALL include a `name` field (string or null) in every file response — GET, list, update, and complete-upload responses.

#### Scenario: File with a name
- **WHEN** a file has a non-null `name`
- **THEN** the response includes `"name": "<value>"`

#### Scenario: File without a name
- **WHEN** a file has a null `name`
- **THEN** the response includes `"name": null`

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
