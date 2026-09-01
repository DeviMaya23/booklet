# Image Thumbnail Generation

## Purpose

Defines the asynchronous thumbnail generation capability. After an image upload completes, a River worker processes a `generate_thumbnail` job that downloads the original image from R2, resizes it to at most 600px on the long edge, and stores the JPEG result back in R2, updating the image record with the thumbnail path.

---

## Requirements

### Requirement: generate_thumbnail job
The system SHALL provide a River worker that processes `generate_thumbnail` jobs. Each job SHALL receive the image ID and user ID of the image to thumbnail. The worker SHALL:
1. Fetch the image record from the database to obtain `image_r2_path` and `mime_type`
2. Download the original image from R2 using `image_r2_path`
3. Decode the image, resize it to a maximum of 600px on the long edge while preserving aspect ratio, and encode the result as JPEG
4. Upload the JPEG to R2 at path `users/{user_id}/thumbnails/{image_id}.jpg`
5. Update the image row's `thumbnail_r2_path` to the uploaded R2 key

The worker SHALL use `github.com/disintegration/imaging` for decode, resize, and encode operations.

#### Scenario: Successful thumbnail generation
- **WHEN** a `generate_thumbnail` job is processed for an image that exists and whose original is in R2
- **THEN** a JPEG file exists at `users/{user_id}/thumbnails/{image_id}.jpg` in R2 and the image row's `thumbnail_r2_path` is set to that key

#### Scenario: Thumbnail fits within 600px on the long edge
- **WHEN** the original image has a long edge greater than 600px
- **THEN** the generated thumbnail has a long edge of exactly 600px and the short edge is scaled proportionally

#### Scenario: Small image is not upscaled
- **WHEN** the original image has both dimensions smaller than 600px
- **THEN** the generated thumbnail dimensions are unchanged from the original

#### Scenario: Job failure is retried by River
- **WHEN** the worker returns an error (e.g. R2 is temporarily unavailable)
- **THEN** River retries the job according to its default retry schedule; `thumbnail_r2_path` remains null until a retry succeeds

#### Scenario: thumbnail_r2_path remains null on permanent failure
- **WHEN** the job exhausts all River retries
- **THEN** `thumbnail_r2_path` stays null; the image record is otherwise intact and accessible

---

### Requirement: Thumbnail job enqueued after CompleteUpload
The system SHALL enqueue a `generate_thumbnail` River job immediately after the CompleteUpload transaction commits successfully. The job SHALL carry the image ID and user ID. If the enqueue call fails, the error SHALL be logged and CompleteUpload SHALL still return success — the image record is durable regardless of enqueue outcome.

#### Scenario: Job enqueued on successful CompleteUpload
- **WHEN** `POST /images/:id/complete` succeeds and the image record is committed
- **THEN** a `generate_thumbnail` job exists in the River job queue for that image ID

#### Scenario: Enqueue failure does not fail CompleteUpload
- **WHEN** the River enqueue call returns an error after the image record commits
- **THEN** `POST /images/:id/complete` still returns 201; the error is logged; `thumbnail_r2_path` stays null
