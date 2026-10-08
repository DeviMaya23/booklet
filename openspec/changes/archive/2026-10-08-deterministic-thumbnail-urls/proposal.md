## Why

Every time a list or gallery response is served, thumbnail presigned URLs are minted with the current timestamp, so each request produces a unique URL and the browser cache never hits. Thumbnails have the highest request volume and smallest payload — fixing cache behaviour here gives the most leverage.

## What Changes

- `r2Storage` gains a new `GenerateDeterministicPresignedGetURL` method that signs with the current 24h window boundary instead of `time.Now()`, producing a stable URL for the entire window (valid 48h)
- `Presigner` handler interface gets the new method added
- All handler call sites that mint thumbnail presigned GET URLs switch to the new method; file, download, and avatar presigning are unchanged
- Thumbnail objects gain `Cache-Control: private, max-age=172800` at upload time so browser disk cache persists across sessions (new thumbnails only — no migration)

## Capabilities

### New Capabilities

- `thumbnail-deterministic-url`: Deterministic presigned GET URL generation for thumbnail objects — window-based signing, stable URL per 24h window, 48h expiry

### Modified Capabilities

- `file-thumbnail-gen-state`: Thumbnail upload now sets `Cache-Control` metadata on the R2 object; no requirement on the gen-state machine itself changes, but the storage contract for thumbnail upload gains a cache header obligation

## Impact

- `backend/internal/storage/r2.go` — new method, `PutObject` signature or worker interface gains cache-control support
- `backend/internal/handler/presign.go` — new method on `Presigner` interface
- `backend/internal/handler/file_handler.go` — `presignThumbnailURL` switches method
- `backend/internal/handler/artpiece_handler.go` — `presignCoverThumbnail` switches method
- `backend/internal/handler/commission_handler.go` — cover thumbnail mint switches method
- `backend/internal/handler/file_upload_handler.go` — post-upload thumbnail URL switches method
- `backend/internal/worker/generate_file_thumbnail.go` — `PutObject` call gains cache-control
- No frontend changes, no DB migrations, no R2 key migrations
