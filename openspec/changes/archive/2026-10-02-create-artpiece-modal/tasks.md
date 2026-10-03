## 1. API Hook

- [x] 1.1 Create `frontend/src/features/artpieces/api/useCreateArtpiece.ts` — mutation calling `POST /artpieces` with `{ title, notes, artist_id, character_ids, file_ids }`
- [x] 1.2 Create `frontend/src/features/files/api/useUpdateFile.ts` — mutation calling `PUT /files/:id` with `{ name, notes }`

## 2. ArtpieceFilesInput component

- [x] 2.1 Create `frontend/src/features/files/components/ArtpieceFilesInput.tsx` with the `DraftFile` discriminated union type (`uploading` | `uploaded`)
- [x] 2.2 Implement the sticky dashed-border drop zone at the top of the scrollable region (drag-over ring, click-to-open file picker)
- [x] 2.3 Implement parallel upload flow for dropped/picked files (init → PUT → complete), inserting an `uploading` row immediately and transitioning to `uploaded` on complete; remove row and show toast on failure
- [x] 2.4 Implement thumbnail readiness: derive `thumbnailUrl` from `useFiles` data by matching file ID; show spinner when null, show thumbnail when non-null
- [x] 2.5 Implement per-row name field with save-on-blur via `useUpdateFile`
- [x] 2.6 Implement per-row notes field with save-on-blur via `useUpdateFile`
- [x] 2.7 Implement remove ("×") button — removes row from draft list, no backend call
- [x] 2.8 Disable name/notes fields when `thumbnailUrl` is null (thumbnail pending state)
- [x] 2.9 Block modal close (ignore Escape and close button) while any row has `status: 'uploading'`
- [x] 2.10 Expose `onFileIdsChange: (ids: string[]) => void` — called whenever the resolved file ID list changes

## 3. ArtpieceFormModal component

- [x] 3.1 Create `frontend/src/features/files/components/ArtpieceFormModal.tsx` — Dialog shell with title, notes, artist (Combobox + `+` button → `ArtistFormModal`), characters (`TokenInput`), `ArtpieceFilesInput`, and Save/Cancel footer
- [x] 3.2 Implement form validation: Save button disabled when title is empty or any upload is in-flight
- [x] 3.3 Accept `initialFileIds?: string[]` prop; on open, read matching entries from `useFiles` data and insert them as `uploaded` rows in `ArtpieceFilesInput`
- [x] 3.4 Implement save handler: call `useCreateArtpiece` with collected form values and current file IDs; show success toast and close on success; show error toast and keep modal open on failure

## 4. Wire up entry points

- [x] 4.1 In `FilesPage.tsx`, add `ArtpieceFormModal` state; pass selected file IDs on "New Artpiece" dropdown click
- [x] 4.2 In `FileInboxGrid.tsx`, thread `onNewArtpiece` through to open the modal with context menu's selected file IDs
- [x] 4.3 In `FileContextMenu.tsx`, confirm `onNewArtpiece` callback is already present (it is — just needs to be wired in the parent)

## 5. Bruno

- [x] 5.1 Update `collection/artpieces/create_artpiece.bru` body to include `"file_ids": []` in the example JSON

## 6. Lint

- [x] 6.1 Run `npm run build` and `npm run lint` in `frontend/`; fix any issues
