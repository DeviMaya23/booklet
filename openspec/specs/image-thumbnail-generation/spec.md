# Image Thumbnail Generation

## Purpose

Defines the asynchronous thumbnail generation capability. After an image upload completes, a River worker processes a `generate_thumbnail` job that downloads the original image from R2, resizes it to at most 600px on the long edge, and stores the JPEG result back in R2, updating the image record with the thumbnail path.

---

## Requirements

### Requirement: Thumbnail job enqueued after CompleteUpload
The system SHALL enqueue a `generate_thumbnail` River job immediately after the CompleteUpload transaction commits successfully. The job SHALL carry the image ID and user ID. If the enqueue call fails, the error SHALL be logged and CompleteUpload SHALL still return success — the image record is durable regardless of enqueue outcome.

#### Scenario: Job enqueued on successful CompleteUpload
- **WHEN** `POST /images/:id/complete` succeeds and the image record is committed
- **THEN** a `generate_thumbnail` job exists in the River job queue for that image ID

#### Scenario: Enqueue failure does not fail CompleteUpload
- **WHEN** the River enqueue call returns an error after the image record commits
- **THEN** `POST /images/:id/complete` still returns 201; the error is logged; `thumbnail_r2_path` stays null
