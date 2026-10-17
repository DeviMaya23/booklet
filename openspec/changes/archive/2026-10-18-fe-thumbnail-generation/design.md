# Design

## Context

See `proposal.md — Why` for motivation. The current implementation is split across:

- `usecase.FileUploadUsecase.InitiateUpload` — generates presigned PUT URL for the original file
- `usecase.FileUploadUsecase.CompleteUpload` — creates the file record with `thumbnail_gen_state = "pending"`, enqueues a `GenerateFileThumbnailArgs` River job
- `worker.GenerateFileThumbnailWorker` — downloads original from R2, decodes/resizes, uploads thumbnail, updates `thumbnail_r2_path` and `thumbnail_gen_state`
- Frontend polls `GET /files` waiting for `thumbnail_gen_state` to change from `"pending"` to `"done"` before showing the thumbnail

The `thumbnail_gen_state` column exists to express the worker pipeline state. With no worker, there is no in-flight state to track.

## Goals / Non-Goals

**Goals:**
- Thumbnail generated in the browser before any PUT; available to the UI immediately, no polling delay
- Backend presigns a PUT URL for the thumbnail at initiate time and sets the path deterministically at complete time
- `thumbnail_gen_state` removed from domain, DB, and API
- Worker and its River registration removed
- UI shows thumbnail from local object URL without waiting for a server round-trip

**Non-Goals:**
- Fallback worker (no server-side retry if the client skips the thumbnail upload)
- Thumbnail regeneration for existing files
- Handling animated GIFs as animated thumbnails (single frame is fine)

## Decisions

### 1. Thumbnail key includes extension determined by mime type

The extension (`.jpg` or `.png`) depends on whether the source is alpha-capable (`image/png`, `image/gif`, `image/webp` → `.png`; others → `.jpg`). Since `mime_type` is in the initiate request, the BE can compute the correct key and content type at initiate time. The thumbnail presigned PUT URL is issued with that exact key and content type. At complete time, `CompleteUpload` derives the same key deterministically from the mime type stored in the pending record — no R2 existence check needed.

Alternative considered: always produce JPEG, fixed `.jpg` key. Rejected because it breaks alpha channels for PNG/WebP art — a specific Booklet requirement.

### 2. `thumbnail_upload_url` is null for non-image files

Non-image files get no thumbnail at any point in the flow. Returning null is explicit and consistent with how `thumbnail_url` already works in GET responses. The FE checks mime type (or the null URL) and skips thumbnail generation for non-images.

### 3. `thumbnail_gen_state` column dropped, not zeroed

There is no valid in-between state after this change. `"pending"` is unreachable (no worker to complete it); `"failed"` is unreachable (nothing to fail); `"not_applicable"` is implicit (non-image = no thumbnail path). The column carries no information the schema doesn't already express through `thumbnail_r2_path`. Dropping it is cleaner than deprecating it.

Migration: `000031_drop_thumbnail_gen_state_from_files.up.sql` drops the column. Down reverses it by re-adding `thumbnail_gen_state TEXT` with a default of `'done'` (best-effort backfill for existing rows that have a thumbnail path; rollback restores the column but cannot recover the original state values).

### 4. Local object URL used for immediate display

After `generateThumbnail(file)` returns a blob, `URL.createObjectURL(blob)` produces a local URL usable in `<img src>` before any PUT completes. This is stored in `DraftFile.localThumbnailUrl`. When the `useFiles` poll eventually returns the server presigned URL, `localThumbnailUrl` is replaced by the server URL (or the component can keep using it until unmount; both refer to the same image). Object URLs are revoked on component unmount to avoid memory leaks.

Alternative considered: use a `data:` base64 URL. Rejected — object URLs are cheaper (no encoding overhead) and garbage-collect cleanly with `revokeObjectURL`.

### 5. No `disabled` state on name/notes inputs during thumbnail generation

The previous UX disabled the name and notes fields while `thumbnail_gen_state === "pending"` to prevent saving before the thumbnail was ready. With client-side generation the thumbnail is available before `completeUpload` returns; the fields should be editable as soon as the file transitions to `status: "uploaded"`.

### 6. Original image dimensions passed in the complete call

Currently `CompleteUpload` inserts an `image_metadata` row with `Width: 0, Height: 0`, and the worker updates the values after decoding the original. With the worker removed, dimensions would stay at zero indefinitely. The artpiece detail view exposes these dimensions on file entries (`/artpieces/:id` response includes `width` and `height` per file).

The FE already decodes the image for thumbnail generation (via `createImageBitmap`). At that point, `bitmap.width` and `bitmap.height` give the original dimensions. The complete call body is extended to accept optional `width` and `height` integers; `CompleteUpload` updates `image_metadata` with these values when provided.

Alternative considered: leave dimensions as 0. Rejected because it would surface 0×0 in the artpiece API response for all new files.

## Risks / Trade-offs

- **Client skips thumbnail upload**: If thumbnail generation fails (e.g. unsupported format like HEIC) the FE catches the error, skips the thumbnail PUT, and calls complete with `thumbnail_uploaded: false`. The BE only sets `thumbnail_r2_path` when `thumbnail_uploaded` is `true` in the complete body; on `false`, `thumbnail_r2_path` remains null. This explicit handshake prevents a path pointing to a non-existent R2 object. The FE shows the category fallback icon for files with no thumbnail path.

- **OffscreenCanvas availability**: OffscreenCanvas is available in all modern browsers but not in very old ones or some WebViews. Not a concern for this app's target audience. Mitigation: none needed.

- **Large image decode in main thread vs. worker**: The current design generates the thumbnail synchronously in the upload function. For very large images this could block the UI thread momentarily. Mitigation: acceptable for single uploads; if it proves slow a Web Worker can be added later without changing the API contract.

## Migration Plan

1. Deploy migration `000031` (drop `thumbnail_gen_state`) — safe to run before or with the code deploy since neither the old nor new code reads the column after the FE deploy
2. Deploy BE changes (new initiate response, complete sets path deterministically, worker removed)
3. Deploy FE changes (thumbnail generation in upload flow, spinner removed)

**Rollback**: restore `000031.down.sql` (re-adds the column with `DEFAULT 'done'`), re-deploy previous BE and FE. The column values will be a best-effort approximation (all rows get `'done'`), not the exact previous state.

**Existing files in River queue**: any in-flight `generate_file_thumbnail` jobs at deploy time will attempt to run against the removed worker registration and will fail gracefully — River logs the unknown job kind and discards it. No data is corrupted.
