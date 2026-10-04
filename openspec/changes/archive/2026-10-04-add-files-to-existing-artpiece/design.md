## Context

The inbox UI (`FilesPage`, `FileInboxGrid`, `FileContextMenu`) already has two "Add to Existing Artpiece" entry points wired up as no-ops — one in the toolbar dropdown and one in the right-click context menu. This change fills them in. No backend changes are required; all needed endpoints exist (`GET /artpieces`, `GET /artpieces/:id`, `PUT /artpieces/:id/files`).

Relevant existing code:
- `FilesPage.tsx` — owns selection state and modal open state; already has `openNewArtpieceModal` as the pattern to follow
- `FileInboxGrid.tsx` — already has `onNewArtpiece` prop; `onAddToExisting` is passed to `FileContextMenu` but propagates nothing to the parent
- `FileContextMenu.tsx` — already renders "Add to Existing Artpiece" button, calls `onAddToExisting()` with no args
- `ArtpieceFormModal.tsx` + `ArtpieceFilesInput.tsx` — pattern reference for the thumbnail strip and Artist/Characters controls

## Goals / Non-Goals

**Goals:**
- New `AddToExistingArtpieceModal` that takes selected file IDs, lets the user search and pick an artpiece, previews the pick, and saves via `PUT /artpieces/:id/files`
- Wire both existing no-op entry points to open the modal
- `FileInboxGrid` propagates selected file IDs upward via a new `onAddToExisting` prop

**Non-Goals:**
- No file upload inside the modal
- No backend changes
- No pagination or server-side filtering for artpiece search — client-side only

## Decisions

### D1: Client-side artpiece search

All artpieces are fetched once on modal open via `useArtpieces()` and filtered in memory by title substring (case-insensitive). Artist and Characters filters narrow the same in-memory list.

**Rationale:** Consistent with how the rest of the app filters (file inbox search, artist combobox in `ArtpieceFormModal`). Artpiece counts are expected to be small.

**Alternative considered:** Server-side `title` query param — deferred; no backend change needed this way.

### D2: Save via GET-then-PUT

On pick, `useArtpiece(id)` fetches the artpiece with its current file list. On Save, the current file IDs are merged with the selected inbox file IDs and sent to `PUT /artpieces/:id/files`.

**Rationale:** `PUT /artpieces/:id/files` replaces the full file list, so we need the current set. Fetching on pick (not on Save) means the preview and the merge data arrive in one step — no extra fetch at Save time.

**Trade-off:** A concurrent edit between pick and save could silently overwrite changes made in another tab. Accepted — this is an existing trade-off in the app (same pattern as `ReplaceFiles`).

### D3: Reuse existing Artist / Characters controls

The Artist dropdown reuses the same `Combobox` + `useArtists()` pattern from `ArtpieceFormModal`, without the "+" create button. The Characters chip filter reuses `TokenInput` + `useCharacters()`, without `createFromText`.

**Rationale:** Both components are self-contained and already in cache from other queries. No adaptation needed beyond omitting the create affordance.

### D4: New hooks — `useArtpieces` and `useArtpiece(id)`

- `useArtpieces()` — `GET /artpieces`, no params, standard `useQuery`
- `useArtpiece(id)` — `GET /artpieces/:id`, enabled only when `id` is non-null

These are new files under `frontend/src/features/artpieces/api/`. They return the existing `artpieceResponse` shape already defined in `useCreateArtpiece.ts` (extended with `artist_name`, `characters`, `thumbnail_url`, `files`).

### D5: FileInboxGrid `onAddToExisting` prop

`FileInboxGrid` gains `onAddToExisting?: (fileIds: string[]) => void`, mirroring `onNewArtpiece`. `FileContextMenu` already calls `onAddToExisting()` with no args — the component passes `contextMenu.ids` to it, same as `onNewArtpiece`.

## Risks / Trade-offs

- **Stale file list on save** — If the artpiece's files are modified between pick and save, the PUT silently overwrites. → Accepted trade-off per D2.
- **`ErrFileAlreadyAttached` on save** — If a selected inbox file was already attached to a *different* artpiece between modal open and save, the backend returns 422. → Surface as a toast error; user can re-select.

## Open Questions

None — all decisions resolved during exploration.
