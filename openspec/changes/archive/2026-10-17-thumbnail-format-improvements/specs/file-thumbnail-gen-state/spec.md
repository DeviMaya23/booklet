# Spec Delta

## ADDED Requirements

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
