## Context

`ArtpieceDetailView` has a "Download all" button that has always been disabled. Files are stored in Cloudflare R2; the backend already has `r2Storage.GetObject(ctx, key) (io.ReadCloser, error)` for direct streaming. Other usecases (`character`, `file`, `file_upload`) already hold a storage dependency, establishing the pattern. The handler layer currently only uses `Presigner` for generating presigned GET URLs.

## Goals / Non-Goals

**Goals:**
- Stream a `.zip` of all artpiece files to the client without buffering the full archive in memory
- Keep zip/storage logic in the usecase layer, consistent with bookleaf's architecture
- Wire the disabled frontend button

**Non-Goals:**
- Selective file download (subset of files)
- Progress indication or resumable downloads
- Async job-based packaging (not needed at expected file counts)

## Decisions

### Zip assembly lives in the usecase, not the handler

The handler calls `GetByID` (ownership check + file list), sets response headers from the result, then delegates to `ArtpieceUsecase.DownloadFiles(ctx, artpiece, w io.Writer)` — passing the already-fetched `*domain.Artpiece`. The usecase owns: `GetObject` calls and `archive/zip` assembly.

This avoids a double DB fetch and makes ownership verification explicit: `DownloadFiles` accepts a domain object, so any caller must have already performed an ownership-scoped fetch before calling it.

**Alternative considered**: handler iterates files and calls `GetObject` directly. Rejected — zip logic is business logic and should be testable without an HTTP context.

### `ArtpieceObjectGetter` interface on the usecase

```go
type ArtpieceObjectGetter interface {
    GetObject(ctx context.Context, key string) (io.ReadCloser, error)
}
```

Defined in `usecase/artpiece_repository.go` alongside the other repository and storage interfaces for the artpiece usecase. `r2Storage` already implements this. Passed to `NewArtpieceUsecase` as a new constructor parameter. The handler's `Presigner` interface is unchanged.

### Streaming zip — headers committed before body

`Content-Type: application/zip` and `Content-Disposition: attachment; filename="<title>.zip"` are set, then `WriteHeader(200)` is called before the usecase starts writing. Any storage error mid-stream can only be logged — it cannot be returned as an HTTP error. The client may receive a partial or malformed zip; this is the standard tradeoff for streaming archives and matches the bookleaf approach.

### Filename scheme

- **Zip file**: `<sanitized-title>.zip` — sanitize strips characters illegal in filenames (`/`, `\`, `:`, `*`, `?`, `"`, `<`, `>`, `|`) and trims surrounding whitespace. Falls back to `artpiece.zip` if title is nil.
- **Entries inside the zip**: `<sanitized-title>-<n>.<ext>` (1-indexed). Extension derived from `MimeType` via a lookup table (`image/jpeg → jpg`, `image/png → png`, `image/gif → gif`, `image/webp → webp`, `application/pdf → pdf`, `application/octet-stream` and unknowns → `bin`). No reliance on `file.Name` — it is unreliable (may be a raw filename, user label, or gibberish).

### Frontend: blob download pattern

```
fetch GET /artpieces/:id/download  (with Authorization header)
  → response.blob()
  → URL.createObjectURL(blob)
  → synthetic <a download="<title>.zip">.click()
  → URL.revokeObjectURL(...)
```

The frontend sanitizes the filename for the `a.download` attribute independently of the backend `Content-Disposition` header.

## Risks / Trade-offs

- **Partial zip on mid-stream error** → Log the error with artpiece ID and file index; close the zip writer to write a valid end-of-central-directory record for whatever completed entries exist. The client gets a partial archive rather than a corrupted one.
- **Large artpieces / slow R2** → No timeout beyond the default request timeout. Acceptable for initial cut; can add per-file context deadlines later if needed.
- **`ArtpieceUsecase` constructor change** → One new parameter (`objectStorage ObjectGetter`). Only `main.go` calls the constructor; easy to update.
