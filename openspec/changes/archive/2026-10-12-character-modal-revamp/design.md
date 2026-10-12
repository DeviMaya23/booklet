# Design

## Context

See `proposal.md` for motivation. Relevant current state:

- `CharacterFormModal` is a single `Dialog` component. Folder state is managed via a diff model (`addedFolders`, `removedIds`, `renamedFolders`) against `character.folders`.
- `FolderPicker` wraps `TokenInput` for chip-based folder selection.
- `usePublicFolders` fetches `GET /folders` with `retry: false`, 15 min stale time. The `token` field is returned by the backend but currently dropped in the FE mapper.
- `useGetFolderImages` already exists for on-demand folder image fetching.
- The modal is currently a fixed-width `max-w-md` dialog.

## Goals / Non-Goals

**Goals:**
- Replace the chip picker with a list-based folder UI and an expandable Choose folders panel
- Surface invalid (auto-removed) folders as dismissible warning rows
- Add per-folder Bookleaf share links using `VITE_BOOKLEAF_HOST/share/:token`
- Add on-demand folder image previews with client-side pagination (15 / +15)
- Handle all empty and error states in the panel
- Restructure the modal header: compact avatar thumbnail inline with Name

**Non-Goals:**
- Backend changes — the API contract is unchanged
- Changing the diff-based submission semantics
- Changing the silent reconciliation logic (only adding display layer on top)

## Decisions

### 1. Panel as Dialog width expansion, not a separate modal

The Choose folders panel slides in as a right column within the same `Dialog`, widening it from `max-w-md` to a two-column layout. The panel is controlled by a local `panelOpen` boolean.

**Alternatives considered:**
- Separate Dialog/Sheet overlay: adds a stacking layer, makes it harder to see form state and folder selection simultaneously.
- Drawer (Sheet from shadcn): possible, but the side-by-side layout in the design is cleaner and keeps context visible.

### 2. `invalidFolders` as additive state alongside existing diff state

The existing `removedIds` logic (populated by the reconciliation `useEffect`) is untouched. A new `invalidFolders: PublicFolder[]` state records which folders were silently removed by reconciliation, so they can be shown as warning rows. The Remove button on a warning row calls `setInvalidFolders(prev => prev.filter(...))` only — it has no effect on `removedIds` (the folder is already there).

On `resetForm`, `invalidFolders` is cleared alongside the other diff state.

### 3. One preview open at a time via `activePreviewId` state

A `activePreviewId: string | null` state in the panel tracks which folder's preview is currently expanded. Opening a new preview sets this id; clicking Hide clears it. `useGetFolderImages` is called with `activePreviewId` (enabled only when set).

Each folder's query is cached by `['folderImages', folderID]` — reopening a previously fetched preview hits the cache, not the network.

### 4. Client-side "Show more" via a `visibleCount` counter

`useGetFolderImages` returns all images in one call. The panel holds a `visibleCount` (default 15) per active preview. "Show more" increments by 15. Resetting `activePreviewId` (closing the preview) resets `visibleCount`. No separate pagination state per folder is needed since only one preview is open at a time.

### 5. Token added to `PublicFolder` type; share URL composed in components

`PublicFolder` gains a `token: string` field. The mapper in `usePublicFolders` includes it from the response. The share URL (`${import.meta.env.VITE_BOOKLEAF_HOST}/share/${folder.token}`) is composed inline at the render site, not stored in state. Valid folder rows (from `publicFoldersQuery.data`) have tokens; invalid/warning rows (absent from the live list) have no token and show no link.

### 6. New component split

`FolderPicker` is replaced by:
- `FolderList` — the default view (list rows, warning rows, Choose folders button)
- `ChooseFoldersPanel` — the expandable right column (folder checklist, search, preview, states)
- `FolderPreview` — the inline image grid within a panel row

`CharacterFormModal` owns the `panelOpen` state and passes it down.

## Risks / Trade-offs

- **Modal width jump on panel open** → mitigate with a CSS transition on the Dialog max-width, or use an `overflow: hidden` expand. If the transition feels jarring, fall back to a Sheet.
- **`useGetFolderImages` called on every preview open** → mitigated by TanStack Query's cache (`staleTime` should be set; suggest 5 min since folder contents change infrequently).
- **Warning rows persist across reopens if form is not reset** → non-issue: `resetForm` is called on every close, which clears `invalidFolders`.

## Open Questions

- Should `useGetFolderImages` have a `staleTime` set? Currently it has none (always refetches). A 5-minute stale time would avoid redundant fetches when the user re-opens a preview in the same session.
