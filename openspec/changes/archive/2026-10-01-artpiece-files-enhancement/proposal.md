## Why

Files uploaded to the system have no management endpoints beyond upload initiation and completion — users cannot fetch, list, update, or delete them. Artpieces also lack a bulk mechanism for managing their full file set in one call, and their create endpoint does not accept initial file attachments.

## What Changes

- **New** `GET /files/:id` — fetch a single file by ID
- **New** `GET /files` — list the user's files, with an `unassigned=true` filter to surface files not yet attached to any artpiece (the "inbox" / dumping-ground view)
- **New** `PUT /files/:id` — update user-editable file fields (notes only for now)
- **New** `DELETE /files/:id` — delete a file record and its R2 object
- **New** `PUT /artpieces/:id/files` — full-replace the file set for an artpiece in one transaction; reconciles attaches/detaches and applies cover rules
- **Updated** `POST /artpieces` — accept an optional `file_ids` array so an artpiece can be created with files already attached

## Capabilities

### New Capabilities

- `file-crud`: Get, list (with unassigned filter), update (notes), and delete files
- `artpiece-file-set`: Artpiece create accepting an optional file array, and a bulk full-replace endpoint for managing an artpiece's file set

### Modified Capabilities

_None — no existing spec-level behavior changes._

## Impact

- **Backend handlers**: new `FileHandler` (alongside existing `FileUploadHandler`); `ArtpieceHandler` gets one new route and one updated request struct
- **Backend usecases**: new `FileUsecase` for file CRUD; `ArtpieceUsecase` extended with `ReplaceFiles` and updated `Create` params; `ArtpieceUsecase` gains a `Transactor` dependency for the bulk replace
- **Backend repositories**: `fileRepository` gains `List`, `UpdateNotes`, `Delete`, and `BulkUpdateArtpieceID` methods
- **R2 storage**: file delete calls `DeleteObject`
- **No DB schema changes** — the single nullable FK (`artpiece_id`) on `files` is sufficient
- **No impact** on the images module or the old pending-upload infrastructure
