# Tasks

## 1. BE: Extend artpiece list response

- [x] 1.1 Add `cover_file_url *string`, `cover_file_mime_type *string`, and `cover_file_name *string` fields to `artpieceResponse` in `artpiece_handler.go`; verify the struct compiles
- [x] 1.2 In `ListArtpieces`, for each artpiece with a non-nil `CoverFile`, presign `CoverFile.FileR2Path` using `GeneratePresignedGetURL` (same TTL as `GetArtpieceByID` file URLs) and populate the three new fields; verify nil `CoverFile` produces null fields
- [x] 1.3 Update the Bruno collection file for `GET /artpieces` to document the three new response fields; verify the file reflects the new shape
- [x] 1.4 Write handler unit tests for `ListArtpieces`: (a) artpiece with cover returns non-null `cover_file_url`, `cover_file_mime_type`, `cover_file_name`; (b) artpiece with null cover returns all three as null; verify tests pass
- [x] 1.5 Write `artpieceRepository.List` integration test: seeded artpiece with a cover file returns the `CoverFile` association with `FileR2Path`, `MimeType`, and `Name` accessible; verify test passes

## 2. BE: Lint

- [x] 2.1 Run `golangci-lint run ./...` on the backend and fix any issues; verify clean output

## 3. FE: Update API type

- [x] 3.1 Add `cover_file_url: string | null`, `cover_file_mime_type: string | null`, and `cover_file_name: string | null` to the `ArtpieceSummary` interface in `useArtpieces.ts`; verify no TypeScript errors

## 4. FE: ArtpieceTile component

- [x] 4.1 Create `frontend/src/features/artpieces/components/ArtpieceTile.tsx`: square `aspect-square` image tile with `overflow-hidden rounded-lg bg-muted`; placeholder icon when no `imageUrl`; caption row (`<figcaption>`) below the tile showing title and artist (artist omitted if absent); caption row hidden when `showDetails` is false; verify renders with and without cover image and with showDetails toggled
- [x] 4.2 Add hover lift (`group/tile hover:scale-[1.02] hover:shadow-md transition-transform`) and ⓘ button (`opacity-0 group-hover/tile:opacity-100 transition-opacity`) top-right with "Open details" tooltip; ⓘ click calls `onInfoClick` without propagating to the tile's `onClick`; verify hover state and click isolation in tests
- [x] 4.3 Write unit tests for `ArtpieceTile`: (a) renders thumbnail when imageUrl present; (b) renders placeholder when imageUrl null; (c) caption row visible when showDetails true; (d) caption row hidden when showDetails false; (e) onClick fires on tile click; (f) onInfoClick fires on ⓘ click without firing onClick; verify all pass

## 5. FE: ArtpiecesGrid with view controls

- [x] 5.1 Rebuild `ArtpiecesGrid` to accept `tileSize`, `showDetails`, `onTileClick(id, index)`, and `onInfoClick(id)` props instead of `onArtpieceOpen`; render `ArtpieceTile` per artpiece; verify component renders the grid
- [x] 5.2 Derive column class from `tileSize`: Small → `grid-cols-4 sm:grid-cols-6 md:grid-cols-8`, Medium → `grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5`, Large → `grid-cols-1 sm:grid-cols-2 md:grid-cols-3`; verify each size produces the expected layout class in tests
- [x] 5.3 Add **View** button to `ArtpiecesPage` toolbar that opens a Radix Popover containing a segmented Small / Medium / Large tile size control and a "Show details" switch; wire state to `ArtpiecesGrid`; verify popover opens and controls are interactive
- [x] 5.4 Read `tileSize` and `showDetails` from `localStorage` on mount (keys `gallery.tileSize`, `gallery.showDetails`); fall back to defaults (`medium`, `true`) when absent, unreadable, or invalid; write to `localStorage` on each change; wrap all reads and writes in try/catch; verify defaults applied when storage is empty, and values restored after setting them

## 6. FE: ArtpiecesPage wiring

- [x] 6.1 Lift `viewerIndex: number | null` state to `ArtpiecesPage`; build `ViewerFile[]` from `filtered` artpieces using `cover_file_url` as `previewUrl`, `thumbnail_url` as `thumbnailUrl`, `cover_file_id` as `id`, `cover_file_name` as `name`, `cover_file_mime_type` as `mimeType`; render `<FileViewer>` when `viewerIndex !== null`; verify lightbox opens at the correct artpiece index
- [x] 6.2 Wire `onTileClick` from `ArtpiecesGrid` to set `viewerIndex` to the clicked artpiece's position in `filtered`; wire `onInfoClick` to `navigate('/app/artpieces/' + id)`; verify single click opens lightbox and ⓘ click navigates
- [x] 6.3 Remove `onDoubleClick`/`DeleteArtpieceDialog` wiring and `pendingDeleteId` state from `ArtpiecesGrid`; remove `onDeleteClick` prop from `ArtpieceTile`; verify no delete UI remains in the grid and the detail view's delete path is unaffected

## 7. FE: Build and lint

- [x] 7.1 Run `npm run build` — no type errors or build failures; fix any issues
- [x] 7.2 Run `npm run lint` — no lint errors; fix any issues

## Workflow follow-up

- Archive the change after review and smoke-testing the gallery UI.
