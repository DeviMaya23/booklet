## Why

File responses currently include `thumbnail_url: null` for both files that are still having their thumbnail generated and files that will never have one (non-images or images that failed thumbnail generation). The frontend cannot distinguish between these cases, causing the inbox UI to show a loading spinner indefinitely for non-image files and for image files whose thumbnail generation job exhausted all retries.

## What Changes

- A `thumbnail_gen_state` field (`"pending"` | `"done"` | `"failed"` | `"not_applicable"`) is added to the `files` table to track the outcome of the auto-generation pipeline.
- All file API responses (get by ID, list, complete upload) include `thumbnail_gen_state`.
- The `generate_file_thumbnail` worker writes `"done"` on success and `"failed"` on its final retry attempt instead of returning an error for River to discard.
- The frontend derives spinner vs. fallback icon from `thumbnail_url` + `thumbnail_gen_state`: spinner only when `thumbnail_url` is null **and** `thumbnail_gen_state === "pending"`.
- The frontend polling condition is updated to stop refetching once no files have `thumbnail_gen_state === "pending"`.

## Capabilities

### New Capabilities

- `file-thumbnail-gen-state`: Tracks the auto-generation pipeline state for a file's thumbnail as a persisted enum on the file record, exposed in all file API responses.

### Modified Capabilities

- `file-upload`: The complete-upload flow and `generate_file_thumbnail` worker gain responsibility for writing `thumbnail_gen_state`.
- `file-crud`: All file responses now include `thumbnail_gen_state`.
- `web-file-inbox`: Spinner and polling logic updated to use `thumbnail_gen_state`.
- `web-artpiece-create-modal`: Spinner and fallback icon logic in `ArtpieceFilesInput` updated to use `thumbnail_gen_state`.

## Impact

- **Database**: Migration to add `thumbnail_gen_state` column to `files` table.
- **Backend**: `domain.File`, `file_upload_usecase`, `generate_file_thumbnail` worker, `file_handler`, `file_upload_handler`, `toFileResponse` helper.
- **Frontend**: `useFiles` (polling condition), `ArtpieceFilesInput` (spinner logic), `FileTile` (fallback icon logic).
