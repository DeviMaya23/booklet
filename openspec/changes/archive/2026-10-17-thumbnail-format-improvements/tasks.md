# Tasks

## 1. WebP Decoder Registration

- [x] 1.1 Add `_ "golang.org/x/image/webp"` as a blank import in `generate_file_thumbnail.go`. Run `go build ./...` from the backend directory and verify it compiles. `golang.org/x/image` is already an indirect dependency so no `go get` is needed; if the build tool promotes it to direct automatically that is the expected outcome; otherwise run `go get golang.org/x/image` explicitly.

## 2. Format Selection Logic

- [x] 2.1 Add a helper (or inline logic) in `generate_file_thumbnail.go` that accepts a mime type string and returns the thumbnail encoding format (`imaging.PNG` or `imaging.JPEG`), the R2 key extension (`.png` or `.jpg`), and the content type string (`image/png` or `image/jpeg`). Mime types `image/png`, `image/gif`, and `image/webp` map to PNG; all others map to JPEG. Update the encode, key construction, and `PutObject` calls to use the selected values. Verify `go build ./...` passes.
- [x] 2.2 Update `generate_file_thumbnail_test.go`: capture the `contentType` passed to `PutObject` in `spyFileThumbnailStorage`. Add a test `TestGenerateFileThumbnailWorker_JPEGSourceProducesJPEGThumbnail` — file has `MimeType: "image/jpeg"`, storage returns a JPEG-encoded image, assert thumbnail key ends in `.jpg` and content type is `image/jpeg`. Add a test `TestGenerateFileThumbnailWorker_PNGSourceProducesPNGThumbnail` — file has `MimeType: "image/png"`, storage returns a PNG-encoded image (use `image/png` stdlib encoder), assert thumbnail key ends in `.png` and content type is `image/png`. Update the existing `TestGenerateFileThumbnailWorker_SuccessSetsStateAndPath` test: set `MimeType: "image/jpeg"` on the file so the expected key (`.jpg`) stays valid. Verify all tests pass with `go test ./internal/worker/...`.

## 3. Decode Timeout

- [x] 3.1 Make the decode timeout duration injectable on `GenerateFileThumbnailWorker` (add an unexported `decodeTimeout time.Duration` field; set it to `30 * time.Second` in `NewGenerateFileThumbnailWorker`). Wrap the `imaging.Decode` and `imaging.Fit` calls in a goroutine; select on the result channel and a `context.WithTimeout(ctx, w.decodeTimeout)` deadline. If the deadline fires, return a `fmt.Errorf("decode timeout: %w", ctx.Err())`. Verify `go build ./...` passes.
- [x] 3.2 Add a test `TestGenerateFileThumbnailWorker_DecodeTimeoutOnFinalAttemptSetsFailedAndReturnsNil` — construct the worker with `decodeTimeout: 1 * time.Millisecond` and a storage spy whose `GetObject` returns a reader that blocks (use a pipe with no writer). Run the job at the final attempt and assert `err == nil` and `thumbnail_gen_state == "failed"`. Add `TestGenerateFileThumbnailWorker_DecodeTimeoutOnNonFinalAttemptReturnsError` — same slow reader, non-final attempt, assert error is non-nil and contains `"decode timeout"`. Verify with `go test ./internal/worker/...`.

## 4. WebP Source Test

- [x] 4.1 Add a small valid WebP file at `internal/worker/testdata/small.webp` (a minimal 1×1 or 4×4 WebP). Add a test `TestGenerateFileThumbnailWorker_WebPSourceGeneratesThumbnail` — file has `MimeType: "image/webp"`, storage spy reads and returns `testdata/small.webp` bytes, assert `thumbnail_gen_state == "done"` and thumbnail key ends in `.png`. Verify with `go test ./internal/worker/...`.

## 5. Lint

- [x] 5.1 Run `golangci-lint run ./...` from the backend directory and fix any issues introduced by this change.

## Workflow follow-up

- Archive the change after review.
