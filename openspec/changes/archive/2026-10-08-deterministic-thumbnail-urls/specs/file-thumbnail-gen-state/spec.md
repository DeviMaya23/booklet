## MODIFIED Requirements

### Requirement: Thumbnail generation state persisted on file record
The system SHALL store a `thumbnail_gen_state` value on every file record. The value SHALL be one of: `pending`, `done`, `failed`, `not_applicable`.

- `pending` — file is an image and thumbnail generation is in progress or queued
- `done` — thumbnail was successfully generated
- `failed` — thumbnail generation was attempted and exhausted all retries without success
- `not_applicable` — file is not an image type; thumbnail generation will not be attempted

When a thumbnail object is stored in R2, the worker SHALL set `Cache-Control: private, max-age=172800` on the object.

#### Scenario: Image file created
- **WHEN** a file upload is completed with an image mime type
- **THEN** the file record is created with `thumbnail_gen_state = "pending"`

#### Scenario: Non-image file created
- **WHEN** a file upload is completed with a non-image mime type
- **THEN** the file record is created with `thumbnail_gen_state = "not_applicable"`

#### Scenario: Thumbnail generation succeeds
- **WHEN** the generate_file_thumbnail worker successfully produces and stores a thumbnail
- **THEN** the file's `thumbnail_gen_state` is updated to `"done"`
- **AND** the thumbnail object in R2 SHALL have `Cache-Control: private, max-age=172800`

#### Scenario: Thumbnail generation fails on final attempt
- **WHEN** the generate_file_thumbnail worker encounters an error and has exhausted all retry attempts
- **THEN** the file's `thumbnail_gen_state` is updated to `"failed"` and no further retries are made

#### Scenario: Existing file with null thumbnail_gen_state
- **WHEN** a file record has a null `thumbnail_gen_state` (pre-migration row)
- **THEN** the API response treats it as `"failed"`
