## Context

The artpiece detail page displays file thumbnails in a grid but offers no way to view a file at full size. The backend already generates presigned GET URLs (`file_url`) per file in the artpiece response, and `ImageMetadata` (width/height) exists in the database but is not loaded by the artpiece repository or exposed in the response. The "download all" flow uses `GeneratePresignedDownloadURL` (which sets `Content-Disposition: attachment`); single-file download needs the same mechanism but without proxying the file through the server.

## Goals / Non-Goals

**Goals:**
- Fullscreen file viewer modal with prev/next navigation and thumbnail filmstrip
- Single-file download via a presigned URL with `Content-Disposition: attachment`
- Expose filename and image dimensions in the artpiece file response
- Reusable `FileViewer` component — not tied to artpieces; usable for any file array

**Non-Goals:**
- Video or audio playback
- PDF rendering in-browser
- Editing metadata from inside the viewer
- Lazy-loading or pagination of large file sets (artpieces are bounded in practice)

## Decisions

### D1: Single-file download via redirect endpoint, not inline URL

**Decision**: Add `GET /files/:id/download` which generates a presigned download URL and returns `{ download_url: string }`. The frontend opens that URL.

**Alternatives considered**:
- Embed `download_url` alongside `file_url` in the artpiece response — generates two presigned URLs per file on every artpiece fetch, even when the user never opens the viewer or downloads anything. Wasteful.
- Reuse `file_url` — the browser renders it rather than downloading it for images. Requires a client-side hack (`fetch` + blob URL) to force download.
- `GET /files/:id/download` as a 302 redirect — cleaner REST, but harder to handle from `<a>` without triggering navigation. JSON response is simpler to consume.

**Rationale**: On-demand, no wasted presigns, and the TTL only needs to be short (5 min is enough to click Download and receive the file).

### D2: FileViewer receives a typed array, not file IDs

**Decision**: `FileViewer` takes `files: ViewerFile[]` and `initialIndex: number` as props. The caller (artpiece detail) maps `ArtpieceFile[]` to `ViewerFile[]`.

```ts
interface ViewerFile {
  id: string
  name: string | null
  mimeType: string
  thumbnailUrl: string | null
  previewUrl: string        // presigned GET URL — used to display the image
  width?: number
  height?: number
}
```

**Alternatives considered**:
- Pass file IDs and let the viewer fetch details — extra requests, more complex state, defeats the purpose of already having the data.

**Rationale**: The artpiece response already carries all needed fields post-enhancement. Keeping the viewer dumb (data-in, no fetching) makes it easy to reuse for other contexts (e.g., cover image browsing) by mapping to `ViewerFile`.

### D3: Expose name + dimensions in artpiece fileRef

**Decision**: Extend `fileRef` in `artpiece_handler.go` with `Name *string`, `Width *int`, `Height *int`. `GetArtpieceByID` adds `Preload("Files.ImageMetadata")`.

**Rationale**: The viewer needs filename for the top bar and no-preview card. Dimensions are already stored and cheaply loadable — one extra preload join. Non-image files produce `null` for both.

### D4: FileViewer lives in `features/files/components/`

**Decision**: `features/files/components/FileViewer.tsx`

**Rationale**: The viewer operates on files generically and should not be scoped to artpieces. Future use cases (cover image browsing, file inbox preview) will find it here without importing across feature boundaries.

### D5: No-preview detection is mime-type based

**Decision**: A file shows the preview `<img>` when its `mimeType` starts with `image/` AND has a non-null `thumbnailUrl`. Otherwise render the no-preview card (file icon, name, Download button).

**Rationale**: PSD, PDF, and other binary types never have a `thumbnail_url` in the current system. Checking `thumbnailUrl` alone covers both "non-image type" and "image without a generated thumbnail yet" gracefully.

### D6: Checkerboard background for image tiles

**Decision**: The image display area uses a CSS checkerboard pattern on `:root` as the background, so transparent PNGs don't vanish against a dark scrim.

**Rationale**: Matches the mockup; distinguishes transparent areas from the dark overlay without adding JS complexity.

## Risks / Trade-offs

- **Stale `previewUrl` TTL**: The artpiece response is cached by React Query with a default stale time. If the page is left open > 1 hour, `file_url` expires and the viewer will show broken images. Mitigation: the viewer always uses `previewUrl` from props (which reflects the last fetch); a refetch on re-focus (React Query default) keeps it fresh. If this becomes a problem, the viewer could fall back to the thumbnail.
- **`ImageMetadata` join on every artpiece fetch**: Adding `Preload("Files.ImageMetadata")` fires a second query per artpiece GET. For artpieces with many files this is a small but nonzero cost. Acceptable given file counts are low in practice.
- **No keyboard trap in modal**: The viewer should trap focus while open (accessibility). Implementation note for the task author.

## Open Questions

- None — all design decisions above are settled.
