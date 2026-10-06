## Why

Artpieces can currently only be created and deleted from the gallery — there is no way to view their full details (notes, artist, characters, file list) or edit them after creation. Users need a detail view to inspect and update artpiece metadata and manage attached files in place.

## What Changes

- Double-clicking a card in the Artpieces gallery opens an artpiece detail view inline (no route change).
- The detail view shows title, notes, artist, characters, and a thumbnail grid of attached files with a cover badge on the current cover.
- An Edit toggle switches the view to edit mode: all fields become editable using the same input components as the create modal; a drag-drop upload strip appears above the file grid; each thumbnail gains a `...` menu with "Set as cover" and "Remove from artpiece."
- Save persists field changes via `PUT /artpieces/:id` and, if files were removed during this edit session, via `PUT /artpieces/:id/files` (full replace). Newly uploaded files are already attached immediately on upload and do not require a separate save call. Cover changes call `PUT /artpieces/:id/cover`. Cancel reverts field edits and pending file removes/cover changes; uploads that already completed remain attached.
- The gallery card's `...` menu and the detail page's `...` menu both get a delete option that opens a shared `DeleteArtpieceDialog`. The dialog includes a checkbox "Also delete attached files" (which maps to `?delete_files=true`). The detail page variant shows the count of attached files; the gallery card variant shows a generic label (file count not available in list response).
- After deletion from the detail view, the panel closes and the gallery list refreshes.

## Capabilities

### New Capabilities
- `web-artpiece-detail`: Inline artpiece detail view with view/edit toggle, file grid, upload strip, and delete entry point.

### Modified Capabilities
- `web-artpieces-gallery`: Add double-click on card to open detail view; replace the existing inline delete AlertDialog with the shared `DeleteArtpieceDialog` (adds the "also delete files" checkbox).

## Impact

- **Frontend only** — no backend changes required; all necessary endpoints already exist.
- `frontend/src/pages/ArtpiecesPage.tsx` — add `selectedArtpieceId` state; render `ArtpieceDetailView` when set.
- `frontend/src/features/artpieces/components/ArtpiecesGrid.tsx` — wire double-click to open detail; replace inline AlertDialog with `DeleteArtpieceDialog`.
- `frontend/src/components/ResourceCard.tsx` — add `onDoubleClick` prop.
- New components:
  - `features/artpieces/components/ArtpieceDetailView.tsx`
  - `features/artpieces/components/ArtpieceDetailFileGrid.tsx`
  - `features/artpieces/components/DeleteArtpieceDialog.tsx`
- New mutation hooks:
  - `features/artpieces/api/useUpdateArtpiece.ts`
  - `features/artpieces/api/useSetCover.ts`
  - `features/artpieces/api/useDetachFileFromArtpiece.ts`
- Update `features/artpieces/api/useDeleteArtpiece.ts` to accept `{ id, deleteFiles }`.
