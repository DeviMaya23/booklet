## Purpose

Tracks the outcome of the automatic thumbnail generation pipeline for a file as a persisted enum on the file record, exposed in all file API responses. Allows clients to distinguish between "thumbnail is being generated", "generation succeeded", "generation failed", and "generation was never applicable".

## Requirements

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

### Requirement: Thumbnail output format matches source mime type
The worker SHALL encode the thumbnail as PNG (stored at a `.png` R2 key with `Content-Type: image/png`) when the source file's mime type is an alpha-capable format (`image/png`, `image/gif`, `image/webp`). For all other mime types the worker SHALL encode as JPEG (stored at a `.jpg` R2 key with `Content-Type: image/jpeg`).

#### Scenario: Opaque source produces JPEG thumbnail
- **WHEN** the source file has a non-alpha mime type (e.g. `image/jpeg`)
- **THEN** the thumbnail is encoded as JPEG, stored at a `.jpg` R2 key, and served with `Content-Type: image/jpeg`

#### Scenario: Alpha-capable source produces PNG thumbnail
- **WHEN** the source file has an alpha-capable mime type (`image/png`, `image/gif`, or `image/webp`)
- **THEN** the thumbnail is encoded as PNG, stored at a `.png` R2 key, and served with `Content-Type: image/png`

### Requirement: Thumbnail generation has a per-job decode timeout
The worker SHALL apply a hard timeout of 30 seconds to the decode and resize steps of thumbnail generation. If the timeout is exceeded, the attempt is treated as an error and follows the standard retry/exhaustion path.

#### Scenario: Decode exceeds timeout
- **WHEN** decoding or resizing the source image takes longer than 30 seconds
- **THEN** the job returns an error; River retries according to its configured policy, and sets `thumbnail_gen_state = "failed"` on final exhaustion
