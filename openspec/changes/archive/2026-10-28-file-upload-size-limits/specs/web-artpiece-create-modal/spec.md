# Spec Delta

## ADDED Requirements

### Requirement: File size validation before upload
The system SHALL validate each selected file's size against the per-category size limit before queuing it for upload. Files that exceed the limit SHALL be silently dropped from the upload batch; for each dropped file the system SHALL display a toast identifying the file name and the applicable size limit. Files within the limit SHALL proceed to upload normally.

#### Scenario: Oversized file dropped from batch
- **WHEN** the user selects or drops one or more files into the artpiece files component and at least one exceeds its category size limit
- **THEN** the oversized file is not queued for upload, a toast is shown naming the file and the limit, and any files within the limit upload normally

#### Scenario: All files within limit
- **WHEN** the user selects or drops files and all are within their category size limits
- **THEN** all files are queued for upload with no toast shown
