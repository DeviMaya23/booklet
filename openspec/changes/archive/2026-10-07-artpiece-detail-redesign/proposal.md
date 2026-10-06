## Why

The current artpiece detail view is a single-column layout with a `...` overflow menu — it buries the cover image, exposes edit actions only through a dropdown, and uses a plain drag-strip for file uploads. The new design surfaces the cover prominently, makes the edit affordances direct and inline, and redesigns the file grid to support per-tile actions and pending-removal state.

## What Changes

- **View mode layout**: two-column layout — details (title, artist with external-link icon, characters as chips, notes) on the left, 360px cover tile on the right; the `...` menu is replaced by a visible "Edit" button inline with the title; a `← Artpieces` back link sits above the title
- **View mode files section**: file count shown in section header alongside a "Download all" button (no-op for now); cover tile shows a `★ Cover` chip instead of a star badge icon; file type label (PNG, JPG, PSD…) on each tile
- **Edit mode header**: title becomes an inline input; Cancel and Save replace the Edit button; an unsaved-changes caption (amber dot + summary text) appears below the title when there are pending changes
- **Edit mode fields**: Artist and Characters share a row; Notes below; cover preview stays on the right with a "Cover preview. Change it with the ☆ on a file below." caption
- **Edit mode tiles**: ☆ (set as cover) and ✕ (remove) controls always visible on each tile; the cover tile shows `★ Cover` in place of ☆; tiles pending removal stay in place but dimmed, with "Will be removed" label and an Undo button; the last grid cell is an "Add files / or drop them here" drop zone
- **Delete artpiece**: moved from the `...` menu to a quiet red "Delete artpiece" text link at the bottom of the edit view
- **BE**: add `mime_type` to `fileRef` in `GetArtpieceByID` response so the frontend can render type labels
- **FE artist link**: resolved client-side from the cached artists list (no BE change needed)

## Capabilities

### New Capabilities
<!-- none -->

### Modified Capabilities

- `web-artpiece-detail`: layout, view mode, edit mode interaction model, file grid behavior, and delete entry point all change significantly

## Impact

- `backend/internal/handler/artpiece_handler.go` — add `MimeType` field to `fileRef` struct and populate it in `GetArtpieceByID`
- `frontend/src/features/artpieces/api/useArtpiece.ts` — add `mime_type` field to `ArtpieceFile` type
- `frontend/src/features/artpieces/components/ArtpieceDetailView.tsx` — full redesign of view and edit mode layout
- `frontend/src/features/artpieces/components/ArtpieceDetailFileGrid.tsx` — redesign: type labels, cover chip, always-visible tile controls, pending-removal state with Undo, add-files-as-last-cell
- `frontend/src/features/artists/api/useArtists.ts` — already fetched; `artist_link` looked up from cache in the detail component
