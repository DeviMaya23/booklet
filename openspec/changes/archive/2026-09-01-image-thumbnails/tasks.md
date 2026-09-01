## 1. Dependency

- [x] 1.1 Add `github.com/disintegration/imaging` to `go.mod` and `go.sum`

## 2. Worker

- [x] 2.1 Create `backend/internal/worker/generate_thumbnail.go` with `GenerateThumbnailArgs` (fields: `ImageID`, `UserID` as `uuid.UUID`) and `GenerateThumbnailWorker`
- [x] 2.2 Implement `Work`: fetch image record from DB, GET original from R2, decode + resize to max 600px long edge with `imaging.Fit`, encode as JPEG, PUT to `users/{user_id}/thumbnails/{image_id}.jpg`, update `thumbnail_r2_path` in DB
- [x] 2.3 Define a `thumbnailImageRepository` interface on the worker for the DB operations it needs (get by ID, update thumbnail path)

## 3. Upload Usecase

- [x] 3.1 Define `JobInserter` interface in `backend/internal/usecase/upload_usecase.go` (same shape as the one in `user_handler.go`: `Insert(ctx, args, opts)`)
- [x] 3.2 Add `jobInserter JobInserter` field to `uploadUsecase` and update `NewUploadUsecase` constructor
- [x] 3.3 In `CompleteUpload`, after `InTransaction` returns nil, call `jobInserter.Insert` with `GenerateThumbnailArgs{ImageID, UserID}`; log and continue on error

## 4. Wiring

- [x] 4.1 Register `GenerateThumbnailWorker` in `main.go` via `river.AddWorker`, passing the image repository and R2 storage it needs
- [x] 4.2 Pass `riverClient` to `NewUploadUsecase` in `main.go`

## 5. Tests

- [x] 5.1 Unit test `CompleteUpload`: assert job is enqueued with correct `ImageID` and `UserID` on success
- [x] 5.2 Unit test `CompleteUpload`: assert 201 is still returned (no error) when `jobInserter.Insert` fails
- [x] 5.3 Unit test `GenerateThumbnailWorker.Work`: assert `thumbnail_r2_path` is updated with the correct R2 key on success
- [x] 5.4 Unit test `GenerateThumbnailWorker.Work`: assert error is returned when R2 GET fails (triggering River retry)

## 6. Bruno Collection

- [x] 6.1 No new endpoint — confirm existing `complete_upload.bru` still reflects the correct flow (no changes needed to request shape)

## 7. Lint

- [x] 7.1 Run `golangci-lint run ./...` from `backend/` and fix any issues
