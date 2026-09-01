## Why

Images have a `thumbnail_r2_path` column that has always been null — no code generates or stores thumbnails. The upcoming frontend needs thumbnails to render image gallery layouts at MVP.

## What Changes

- New River job `generate_thumbnail` enqueued by CompleteUpload after the image record is created
- Worker fetches the original image from R2, resizes it to max 600px on the long edge (aspect-ratio preserving), encodes as JPEG, and uploads to `users/{user_id}/thumbnails/{image_id}.jpg`
- Worker updates the image row's `thumbnail_r2_path` with the R2 key
- `thumbnail_r2_path` starts as `null` after CompleteUpload and is populated asynchronously — callers must handle the null state
- New dependency: `github.com/disintegration/imaging` for in-process image resizing
- Job is enqueued transactionally within the same pgx transaction as image creation — if the transaction rolls back, no job is enqueued

## Capabilities

### New Capabilities

- `image-thumbnail-generation`: background job that generates a 600px JPEG thumbnail from the uploaded original and stores it in R2, triggered per image on upload completion

### Modified Capabilities

- `image-upload`: CompleteUpload gains a new side effect — after committing the image record, a thumbnail generation job is enqueued; `thumbnail_r2_path` on the created image is initially null

## Impact

- **Backend**: new `worker/generate_thumbnail.go`; `upload_usecase.go` gains a `JobInserter` dependency; thumbnail job registered in `main.go`
- **Dependencies**: `github.com/disintegration/imaging` added to `go.mod`
- **API**: no endpoint signature changes; `thumbnail_r2_path` was already in the image response shape and spec — clients already need to handle null
- **DB**: no migration needed; column exists
