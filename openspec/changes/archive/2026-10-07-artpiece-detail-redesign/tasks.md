## 1. Backend — expose mime_type on artpiece file response

- [x] 1.1 Add `MimeType string \`json:"mime_type"\`` field to `fileRef` struct in `artpiece_handler.go`
- [x] 1.2 Populate `MimeType: f.MimeType` in the `GetArtpieceByID` file loop
- [x] 1.3 Add `mime_type string` to `ArtpieceFile` interface in `useArtpiece.ts`

## 2. ArtpieceDetailFileGrid — redesign

- [x] 2.1 Remove the `...` dropdown menu from edit-mode tiles; replace with always-visible ☆ and ✕ controls
- [x] 2.2 Add type label derived from `mime_type` (uppercase subtype, e.g. PNG / JPG / PSD) to every tile
- [x] 2.3 Replace the cover badge icon with a `★ Cover` chip on the cover tile (view and edit modes)
- [x] 2.4 Change pending-removal rendering: keep removed tiles in the grid but dimmed, with "Will be removed" label and an Undo button (requires the grid to receive all files, not a pre-filtered list)
- [x] 2.5 Move the add-files drop zone from a separate strip into the last grid cell; remove the standalone drag-strip from `ArtpieceDetailView`

## 3. ArtpieceDetailView — layout and view mode

- [x] 3.1 Switch to two-column CSS grid layout (`details | 360px cover`) in view mode
- [x] 3.2 Render the `← Artpieces` back link above the title; place a visible "Edit" button inline with the title (remove the `...` overflow menu from view mode)
- [x] 3.3 Show Artist with an external-link icon when `artist_link` is set (look up from `useArtists` cache by `artist_id`)
- [x] 3.4 Render characters as chips instead of comma-separated text
- [x] 3.5 Add file count and a no-op "Download all" button to the files section header

## 4. ArtpieceDetailView — edit mode header and fields

- [x] 4.1 Keep two-column grid in edit mode; place title input and Cancel/Save in the left column header
- [x] 4.2 Compute unsaved-changes summary from local state (field diffs, `locallyRemovedIds.size`, `pendingCoverFileId`) and render amber-dot caption below the title input when changes exist
- [x] 4.3 Layout Artist (combobox) and Characters (TokenInput) side-by-side in a row; Notes textarea below
- [x] 4.4 Add cover preview tile in the right column with "Cover preview. Change it with the ☆ on a file below." caption
- [x] 4.5 Move "Delete artpiece" from view-mode `...` menu to a quiet red text link at the bottom of edit mode

## 5. Quality

- [x] 5.1 Run `npm run build` and fix any type errors
- [x] 5.2 Run `npm run lint` and fix any lint issues
- [x] 5.3 Run `golangci-lint run ./...` and fix any issues
