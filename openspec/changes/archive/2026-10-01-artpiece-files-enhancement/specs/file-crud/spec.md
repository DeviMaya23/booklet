## ADDED Requirements

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
The system SHALL allow a user to update the notes field on a file they own.

#### Scenario: Successful update
- **WHEN** a user sends `PUT /files/:id` with a valid body
- **THEN** the system updates the notes field and returns 200 with the updated file

#### Scenario: File not found or not owned
- **WHEN** a user sends `PUT /files/:id` for a file that does not exist or belongs to another user
- **THEN** the system returns 404

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
