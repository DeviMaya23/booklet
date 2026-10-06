## Why

Users have no way to bulk-download the files attached to an artpiece. The "Download all" button exists in the detail view but is disabled. This change wires it up.

## What Changes

- New endpoint `GET /artpieces/:id/download` — streams all attached files as a `.zip` archive directly to the client, never buffering the full archive in memory
- Files inside the zip are named `<artpiece-title>-<n>.<ext>` (extension derived from MimeType); zip file is named `<artpiece-title>.zip`; both fall back to `artpiece` if the title is nil
- `ArtpieceUsecase` gains a `DownloadFiles(ctx, artpieceID, userID, w io.Writer) error` method and an `ObjectGetter` storage dependency
- Frontend "Download all" button wired: fetches the endpoint with auth, triggers a blob download via a synthetic anchor, then revokes the object URL

## Capabilities

### New Capabilities

- `artpiece-download-all`: Streaming zip download of all files attached to an artpiece via a single authenticated endpoint

### Modified Capabilities

- `web-artpiece-detail`: "Download all" button changes from permanently disabled to functional in view mode

## Impact

- **Backend**: `ArtpieceUsecase` (new method + new constructor dep), `ArtpieceHandler` (new handler method), `main.go` (new route + pass `r2Storage` as `ObjectGetter`)
- **Frontend**: `ArtpieceDetailView.tsx` (wire button), new API call in `useArtpiece.ts` or adjacent file
- **Bruno**: new `.bru` file for `GET /artpieces/:id/download`
- **No breaking changes** — existing endpoints and the artpiece response shape are unchanged
