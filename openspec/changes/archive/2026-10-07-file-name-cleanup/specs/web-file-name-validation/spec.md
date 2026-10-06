## ADDED Requirements

### Requirement: File name character validation
The system SHALL validate file name inputs against a forbidden-character set on blur and block the save call when the name is invalid.

#### Scenario: Name contains a forbidden character
- **WHEN** a user blurs a file name input that contains any of the characters `/ \ : * ? " < > |`
- **THEN** the name input displays a red border and an inline error message "Name contains invalid characters (/ \\ : * ? \" < > |)"; the save call to `PUT /files/:id` is not made

#### Scenario: Name is blank
- **WHEN** a user blurs a file name input that is empty or contains only whitespace
- **THEN** the name input displays a red border and an inline error message "Name is required"; the save call is not made

#### Scenario: Name exceeds maximum length
- **WHEN** a user blurs a file name input whose value exceeds 255 characters
- **THEN** the name input displays a red border and an inline error message "Name must be 255 characters or fewer"; the save call is not made

#### Scenario: Valid name clears error
- **WHEN** a user edits a name input that previously showed a validation error and blurs with a valid value
- **THEN** the red border and error message are cleared and the save call proceeds normally
