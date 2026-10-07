## Why

Users can view file thumbnails on the artpiece detail page but have no way to inspect a file at full size, see its metadata, or download it individually — they must download the entire zip to get one file. Adding a fullscreen file viewer closes this gap and makes the detail page significantly more useful.

## What Changes

- New `FileViewer` component: a fullscreen modal overlay with prev/next navigation, a thumbnail filmstrip, a top bar (filename, type badge, pixel dimensions, position counter, Download, Close), and a no-preview fallback state for non-renderable file types (PSD, etc.)
- New `GET /files/:id/download` endpoint that returns a presigned download URL (`Content-Disposition: attachment`) for a single file
- Extend the artpiece detail API response (`fileRef`) to include `name`, `width`, and `height` per file
- `ImageMetadata` is preloaded from the artpiece repository so dimensions are available
- Clicking a file tile in the artpiece detail view (view mode) opens `FileViewer` at that index

## Capabilities

### New Capabilities

- `web-file-viewer`: Generic fullscreen file viewer modal — image display with checkerboard background to reveal transparency, no-preview fallback card, top bar with filename/type/dimensions/counter/download/close, prev/next arrows, thumbnail filmstrip
- `file-single-download`: New `GET /files/:id/download` endpoint returning `{ download_url }` — a short-lived presigned URL with `Content-Disposition: attachment`

### Modified Capabilities

- `web-artpiece-detail`: File tiles in view mode are now clickable and open `FileViewer`
- `artpiece-files`: Artpiece file response (`fileRef`) now includes `name`, `width`, and `height` fields

## Impact

- **Backend**: `artpiece_handler.go` — `fileRef` struct gains `Name`, `Width`, `Height`; `GetArtpieceByID` preloads `Files.ImageMetadata`; new `GET /files/:id/download` route and handler method added to `FileHandler`
- **Frontend**: New `FileViewer` component under `features/files/components/`; `ArtpieceFile` type gains `name`, `width?`, `height?`; `ArtpieceDetailFileGrid` (view mode tiles) gains `onClick` handler; new `useDownloadFile` hook
- **No breaking changes** — existing `file_url` and `thumbnail_url` fields are unchanged
