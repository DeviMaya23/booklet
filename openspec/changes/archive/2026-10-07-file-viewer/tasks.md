## 1. Backend — Extend artpiece file response

- [x] 1.1 Add `Preload("Files.ImageMetadata")` to `artpieceRepository.GetByID` in `backend/internal/repository/artpiece_repository.go`
- [x] 1.2 Extend `fileRef` struct in `backend/internal/handler/artpiece_handler.go` with `Name *string`, `Width *int`, `Height *int`
- [x] 1.3 Populate `Name`, `Width`, `Height` when building `fileRefs` in `GetArtpieceByID` (read from `f.Name` and `f.ImageMetadata`)

## 2. Backend — Single file download endpoint

- [ ] 2.1 ~~Add `DownloadURL(ctx context.Context, id uuid.UUID, userID uuid.UUID) (string, error)` method to `FileUsecase` interface in `backend/internal/handler/file_handler.go` and implement it in `backend/internal/usecase/file_usecase.go` — generates a presigned download URL with 5-minute TTL~~ **Not implemented as described** — presigning is a handler-layer concern; `GetDownloadURL` calls `fileUsecase.GetByID` + `h.presigner.GeneratePresignedDownloadURL` directly without a usecase wrapper
- [x] 2.2 Add `GetDownloadURL` handler method to `FileHandler` in `backend/internal/handler/file_handler.go` — parses `:id`, calls `DownloadURL`, returns `{ "download_url": string }`; 400 on bad UUID, 404 on not found
- [x] 2.3 Register `GET /files/:id/download` route in `backend/cmd/server/main.go`
- [x] 2.4 Add `PresignDownloadTTL = 5 * time.Minute` constant to `backend/internal/usecase/shared.go`
- [ ] 2.5 ~~Write unit tests for `FileUsecase.DownloadURL`~~ **Not applicable** — no `DownloadURL` usecase method exists; handler tests cover this path instead (see task 2.6)
- [x] 2.6 Write unit tests for `GetDownloadURL` handler — success, invalid UUID, not found
- [x] 2.7 Create `collection/files/download_file.bru` Bruno request for `GET /files/:id/download`
- [x] 2.8 Run `golangci-lint run ./...` from `backend/` and fix any issues

## 3. Frontend — Type updates

- [x] 3.1 Add `name: string | null`, `width: number | null`, `height: number | null` to `ArtpieceFile` interface in `frontend/src/features/artpieces/api/useArtpiece.ts`

## 4. Frontend — `useDownloadFile` hook

- [x] 4.1 Create `frontend/src/features/files/api/useDownloadFile.ts` — mutation that calls `GET /files/:id/download`, receives `{ download_url }`, and opens the URL via `window.open(url, '_blank')`

## 5. Frontend — `FileViewer` component

- [x] 5.1 Create `frontend/src/features/files/components/FileViewer.tsx` with the `ViewerFile` interface (`id`, `name`, `mimeType`, `thumbnailUrl`, `previewUrl`, `width?`, `height?`) and `FileViewerProps` (`files`, `initialIndex`, `onClose`)
- [x] 5.2 Implement top bar: filename, type badge (reuse/extract `mimeTypeLabel` from `ArtpieceDetailFileGrid`), pixel dimensions when available, position counter (N of M), Download button (calls `useDownloadFile`), Close button
- [x] 5.3 Implement image preview area: checkerboard CSS background on the tile, `<img>` centered with `object-contain`, dark scrim overlay behind the tile
- [x] 5.4 Implement no-preview fallback: file icon, filename, type descriptor string, "Download file" button (calls `useDownloadFile`)
- [x] 5.5 Implement prev/next arrow buttons — hidden at first/last index
- [x] 5.6 Implement thumbnail filmstrip at bottom — active thumbnail highlighted with a ring; clicking a thumbnail jumps to that index
- [x] 5.7 Implement Escape key listener to close the viewer; trap focus within the modal while open
- [x] 5.8 Render `FileViewer` as a portal/fixed overlay (use `Dialog`/`DialogContent` from shadcn if it fits, or a raw fixed-positioned div with `z-50`)

## 6. Frontend — Wire up to artpiece detail

- [x] 6.1 Add `viewerIndex: number | null` state to `ArtpieceDetailView`; open `FileViewer` when non-null, pass mapped `ViewerFile[]` array and `onClose={() => setViewerIndex(null)}`
- [x] 6.2 Pass `onOpen?: (index: number) => void` prop to `ArtpieceDetailFileGrid` and call it on tile click in view mode
- [x] 6.3 Extract `mimeTypeLabel` helper to a shared location (e.g., `frontend/src/features/files/lib/mimeTypeLabel.ts`) and import it from both `ArtpieceDetailFileGrid` and `FileViewer`

## 7. Frontend — Build and lint

- [x] 7.1 Run `npm run build` from `frontend/` and fix any type or build errors
- [x] 7.2 Run `npm run lint` from `frontend/` and fix any lint issues
