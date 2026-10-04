## 1. Frontend API Hooks

- [x] 1.1 Create `useArtpieces` hook (`frontend/src/features/artpieces/api/useArtpieces.ts`) — calls `GET /artpieces`, returns full artpiece list with `thumbnail_url`, `artist_name`, `characters[]`
- [x] 1.2 Create `useArtpiece` hook (`frontend/src/features/artpieces/api/useArtpiece.ts`) — calls `GET /artpieces/:id`, enabled only when `id` is non-null, returns artpiece with `files[]`
- [x] 1.3 Create `useAttachFilesToArtpiece` mutation hook (`frontend/src/features/artpieces/api/useAttachFilesToArtpiece.ts`) — calls `PUT /artpieces/:id/files` with a merged file ID array

## 2. AddToExistingArtpieceModal Component

- [x] 2.1 Create `AddToExistingArtpieceModal.tsx` in `frontend/src/features/files/components/` with props: `open`, `onOpenChange`, `fileIds: string[]`, `onSuccess?: () => void`
- [x] 2.2 Implement selected-file thumbnail strip (top section): read file data from `useFiles()` cache by ID, render thumbnail or fallback icon per tile, filename truncated with `title` attribute for hover, wrapping grid with fixed max-height and vertical scroll
- [x] 2.3 Implement artpiece search: fetch all artpieces via `useArtpieces()`, filter client-side by title substring on input change
- [x] 2.4 Implement Filters toggle (collapsed by default): Artist dropdown (Combobox + `useArtists()`, no create button) and Characters chip filter (TokenInput + `useCharacters()`, no `createFromText`)
- [x] 2.5 Apply Artist and Characters filters to the artpiece list in combination with the title search (intersection)
- [x] 2.6 Implement artpiece results list: clicking a result calls `useArtpiece(id)` to fetch with files, sets it as the picked artpiece
- [x] 2.7 Implement picked-artpiece preview section: show cover thumbnail (or placeholder), title, artist name, character names; replace on re-pick
- [x] 2.8 Implement Save: merge picked artpiece's current file IDs with `fileIds` prop, call `PUT /artpieces/:id/files`, show success toast, call `onSuccess`, close modal; show error toast on 422
- [x] 2.9 Disable Save button when no artpiece is picked or save is in flight

## 3. Wire Entry Points

- [x] 3.1 Add `onAddToExisting?: (fileIds: string[]) => void` prop to `FileInboxGrid`; pass `contextMenu.ids` to it in the context menu handler (mirroring `onNewArtpiece`)
- [x] 3.2 In `FilesPage`: add `addToArtpieceModalOpen` and `addToArtpieceModalFileIds` state; implement `openAddToArtpieceModal(fileIds)` handler; wire toolbar dropdown "Add to Existing Artpiece" item and `FileInboxGrid`'s new `onAddToExisting` prop to it
- [x] 3.3 Mount `AddToExistingArtpieceModal` in `FilesPage` with `onSuccess={() => setSelection(new Set())}`

## 4. Build & Lint

- [x] 4.1 Run `npm run build` in `frontend/` and fix any type errors
- [x] 4.2 Run `npm run lint` in `frontend/` and fix any lint issues
