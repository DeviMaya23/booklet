## 1. Backend — Usecase

- [x] 1.1 Create `usecase/artpiece_storage.go` with `ObjectGetter` interface (`GetObject(ctx, key) (io.ReadCloser, error)`)
- [x] 1.2 Add `objectStorage ObjectGetter` field to `ArtpieceUsecase` and update `NewArtpieceUsecase` constructor
- [x] 1.3 Implement `ArtpieceUsecase.DownloadFiles(ctx, artpieceID, userID uuid.UUID, w io.Writer) error`: call `GetByID` for ownership check and file list, build zip via `archive/zip`, stream each file from `objectStorage.GetObject`, name entries `<sanitized-title>-<n>.<ext>` (extension from MimeType lookup, fallback `artpiece` if no title)
- [x] 1.4 Write unit tests for `DownloadFiles`: success with files, success with no files (empty zip), artpiece not found, R2 error mid-stream (partial zip + error logged)

## 2. Backend — Handler & Route

- [x] 2.1 Add `DownloadFiles` handler to `ArtpieceHandler`: parse `:id`, auth check, set `Content-Type: application/zip` + `Content-Disposition: attachment; filename="<title>.zip"`, call `artpieceUsecase.DownloadFiles(ctx, id, userID, c.Response())`, log error if returned
- [x] 2.2 Register route `GET /artpieces/:id/download` in `main.go` (protected group)
- [x] 2.3 Pass `r2Storage` as the `ObjectGetter` argument in `NewArtpieceUsecase` call in `main.go`

## 3. Bruno

- [x] 3.1 Create `collection/artpieces/download_all.bru` for `GET /artpieces/:id/download`

## 4. Frontend

- [x] 4.1 Add `downloadArtpieceFiles(id, getToken)` API function: fetch `GET /artpieces/:id/download` with auth header, receive blob, trigger synthetic anchor download, revoke object URL; sanitize artpiece title for `a.download` attribute
- [x] 4.2 Wire "Download all" button in `ArtpieceDetailView.tsx`: remove `disabled`, add `onClick` calling the download function with the current artpiece id and title

## 5. Linting & Build

- [x] 5.1 Run `golangci-lint run` in `backend/` and fix any issues
- [x] 5.2 Run `npm run build` and `npm run lint` in `frontend/` and fix any issues
