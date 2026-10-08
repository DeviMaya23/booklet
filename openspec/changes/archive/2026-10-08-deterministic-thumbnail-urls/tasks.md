## 1. Spike — Verify R2 accepts fixed signing time

- [x] 1.1 Confirm `s3.PresignOptions.Presigner` field exists in `aws-sdk-go-v2/service/s3 v1.105.2` (check SDK source or go doc)
- [x] 1.2 Write a throwaway script: call `PresignGetObject` with a fixed `windowStart` twice (>60s apart) and assert byte-identical output, then HTTP GET the URL and confirm 200

## 2. Storage Layer — Deterministic Presigning

- [x] 2.1 Add constants to `usecase/shared.go`: `ThumbnailPresignWindowSize = 24 * time.Hour`, `ThumbnailPresignExpiry = 48 * time.Hour`
- [x] 2.2 Implement `deterministicSigner` wrapper in `storage/r2.go`: wraps `v4.NewSigner()`, overrides `signingTime` with the caller-provided fixed time in `PresignHTTP`
- [x] 2.3 Add `GenerateDeterministicPresignedGetURL(ctx context.Context, key string) (string, error)` to `r2Storage` — computes `windowStart`, injects `deterministicSigner` via `s3.PresignOptions.Presigner`, calls `PresignGetObject` with `ThumbnailPresignExpiry`

## 3. Storage Layer — Cache-Control on Thumbnail Upload

- [x] 3.1 Add `cacheControl string` parameter to `PutObject` in `r2Storage` and set `CacheControl` on `s3.PutObjectInput` when non-empty
- [x] 3.2 Update `worker.fileThumbnailStorageService` interface: add `cacheControl string` to `PutObject`
- [x] 3.3 Update `GenerateFileThumbnailWorker.Work`: pass `"private, max-age=172800"` in the `PutObject` call for the thumbnail

## 4. Handler Interface

- [x] 4.1 Add `GenerateDeterministicPresignedGetURL(ctx context.Context, key string) (string, error)` to `Presigner` interface in `handler/presign.go`
- [x] 4.2 Update handler test fakes to implement the new method (stub returning `"https://fake-deterministic-url"` and nil error)

## 5. Handler Call Sites — Switch Thumbnail Presigning

- [x] 5.1 `file_handler.go`: `presignThumbnailURL` calls `h.presigner.GenerateDeterministicPresignedGetURL` instead of `GeneratePresignedGetURL`
- [x] 5.2 `artpiece_handler.go`: `presignCoverThumbnail` and the per-file thumbnail mint in the detail response both switch to `GenerateDeterministicPresignedGetURL`
- [x] 5.3 `commission_handler.go`: cover thumbnail mint switches to `GenerateDeterministicPresignedGetURL`
- [x] 5.4 `file_upload_handler.go`: post-upload thumbnail URL mint switches to `GenerateDeterministicPresignedGetURL`

## 6. Unit Tests

- [x] 6.1 Unit test `GenerateDeterministicPresignedGetURL`: mock `s3.PresignClient`, assert that two calls with the same key and the same clock return identical URLs; assert a call after window rollover returns a different URL
- [x] 6.2 Unit test `GenerateFileThumbnailWorker`: assert the `PutObject` call for the thumbnail passes `"private, max-age=172800"` as the cache-control value

## 7. Lint & Build

- [x] 7.1 Run `golangci-lint run ./...` from `backend/` and fix any issues
