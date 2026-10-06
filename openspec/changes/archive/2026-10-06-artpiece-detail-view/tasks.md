## 1. New mutation hooks

- [x] 1.1 Add `useUpdateArtpiece` — `PATCH /artpieces/:id` with `{ title, notes, artistId, characterIds }`
- [x] 1.2 Add `useSetCover` — `PUT /artpieces/:id/cover` with `{ file_id }`
- [x] 1.3 Add `useDetachFileFromArtpiece` — `DELETE /artpieces/:id/files/:file_id`
- [x] 1.4 Update `useDeleteArtpiece` to accept `{ id: string; deleteFiles: boolean }` and append `?delete_files=true` when true

## 2. Shared components and ResourceCard update

- [x] 2.1 Add `onDoubleClick` prop to `ResourceCard`
- [x] 2.2 Create `DeleteArtpieceDialog` — takes `artpieceId`, `open`, `onOpenChange`, `onSuccess`, optional `fileCount`; checkbox label is "Also delete X attached files" when count provided, "Also delete attached files" otherwise
- [x] 2.3 Update `ArtpiecesGrid` — wire double-click on each card to a new `onArtpieceOpen(id)` callback prop; replace inline `AlertDialog` with `DeleteArtpieceDialog`
- [x] 2.4 Update `ArtpiecesPage` — pass `onArtpieceOpen` to `ArtpiecesGrid`

## 3. Artpiece detail view — view mode

- [x] 3.1 Create `ArtpieceDetailFileGrid` component — thumbnail grid accepting `files: ArtpieceFile[]` and `coverFileId: string | null`; renders cover badge on matching thumbnail; display-only
- [x] 3.2 Create `ArtpieceDetailView` component (view mode only for now) — back chevron closes view; title with `...` menu (Edit, Delete); notes, artist, characters display; `ArtpieceDetailFileGrid` for files; `DeleteArtpieceDialog` wired to Delete with `fileCount` from loaded artpiece
- [x] 3.3 Update `ArtpiecesPage` — add `selectedArtpieceId` state; render `ArtpieceDetailView` when set; pass `onClose` to close it; pass `onDeleted` callback to close panel and refresh gallery list

## 4. Artpiece detail view — edit mode

- [x] 4.1 Add edit mode state to `ArtpieceDetailView`; clicking Edit in `...` menu enters edit mode; title input, notes textarea, artist combobox, characters `TokenInput` pre-populated from loaded artpiece data
- [x] 4.2 Add `locallyRemovedIds: Set<string>` and `pendingCoverFileId: string | null` state; extend `ArtpieceDetailFileGrid` to accept edit mode props — per-thumbnail `...` menu with "Set as cover" (disabled when already cover) and "Remove from artpiece"
- [x] 4.3 Add upload strip above file grid in edit mode — drag-drop + file picker; on upload complete call `POST /artpieces/:id/files/:file_id` to attach immediately; new thumbnail appears in grid
- [x] 4.4 Implement Save: call `useUpdateArtpiece` for fields; if `locallyRemovedIds` is non-empty or new files uploaded call `useAttachFilesToArtpiece` (PUT replace) with current file set minus removed IDs; if `pendingCoverFileId` set call `useSetCover`; on all success return to view mode and show success toast
- [x] 4.5 Implement Cancel: clear `locallyRemovedIds`, clear `pendingCoverFileId`, reset field inputs to loaded values; return to view mode

## 5. Final checks

- [x] 5.1 Run `npm run build` and fix any type errors
- [x] 5.2 Run `npm run lint` and fix any issues
