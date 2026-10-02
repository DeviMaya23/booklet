## MODIFIED Requirements

### Requirement: Update file notes
The system SHALL allow a user to update the `name` and/or `notes` fields on a file they own via a single atomic request. Both fields are optional and nullable.

#### Scenario: Successful update
- **WHEN** a user sends `PUT /files/:id` with a body containing `name`, `notes`, or both (each nullable)
- **THEN** the system updates the specified fields in a single database write and returns 200 with the updated file including the new `name` and `notes` values

#### Scenario: File not found or not owned
- **WHEN** a user sends `PUT /files/:id` for a file that does not exist or belongs to another user
- **THEN** the system returns 404

## ADDED Requirements

### Requirement: File response includes name
The system SHALL include a `name` field (string or null) in every file response — GET, list, update, and complete-upload responses.

#### Scenario: File with a name
- **WHEN** a file has a non-null `name`
- **THEN** the response includes `"name": "<value>"`

#### Scenario: File without a name
- **WHEN** a file has a null `name`
- **THEN** the response includes `"name": null`
