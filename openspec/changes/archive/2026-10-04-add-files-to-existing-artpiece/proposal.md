## Why

The inbox's "Add to Existing Artpiece" action (context menu and toolbar dropdown) has been a no-op since the inbox UI shipped. Users who have already selected files in the inbox have no way to attach them to an artpiece they've already created without going through a different flow.

## What Changes

- New modal: `AddToExistingArtpieceModal` — lets the user pick an existing artpiece via title search (client-side filtered) with Artist and Characters filter controls, preview the picked artpiece, and confirm attachment.
- New frontend API hooks: `useArtpieces` (list all artpieces) and `useArtpiece(id)` (fetch single artpiece with files).
- Wire the two existing no-op entry points in `FilesPage` (toolbar dropdown) and `FileInboxGrid` (context menu) to open the new modal with the selected file IDs.
- `FileInboxGrid` gains an `onAddToExisting` prop (parallel to existing `onNewArtpiece`) to propagate selected file IDs to the parent.
- Save action calls `GET /artpieces/:id` on pick to retrieve current file IDs, then calls `PUT /artpieces/:id/files` with the merged set (existing + newly selected).

## Capabilities

### New Capabilities

- `web-add-files-to-artpiece`: The "Add to Existing Artpiece" modal flow — selected-file strip, artpiece title search with Artist/Characters filters, picked-artpiece preview, and save via replace-files.

### Modified Capabilities

- `web-file-inbox`: Entry points for "Add to Existing Artpiece" are wired up (were no-ops); `FileInboxGrid` gains the `onAddToExisting` prop.

## Impact

- **Frontend**: New modal component, two new API hooks, prop additions to `FileInboxGrid`, wiring in `FilesPage`.
- **Backend**: No changes — uses existing `GET /artpieces`, `GET /artpieces/:id`, and `PUT /artpieces/:id/files` endpoints.
- **No breaking changes.**
