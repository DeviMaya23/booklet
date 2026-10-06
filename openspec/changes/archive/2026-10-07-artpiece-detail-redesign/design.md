## Context

The artpiece detail view is a single-panel component rendered inline on the Artpieces page. It has two modes (view / edit) managed by local state. The file grid is a separate `ArtpieceDetailFileGrid` component. The current edit mode uses a `...` dropdown for tile actions and a drag-strip above the grid for file uploads. The delete entry point lives in the view-mode `...` menu.

The redesign is purely frontend with one small backend addition (`mime_type` on the artpiece file response). No new routes, no new hooks, no new API calls.

## Goals / Non-Goals

**Goals:**
- Two-column layout in both view and edit modes (details left, cover right)
- Inline title input in edit mode; visible Edit button in view mode
- Always-visible per-tile controls (☆ / ✕) in edit mode replacing the `...` dropdown
- Pending-removal state: tiles stay in grid, dimmed, with Undo
- Add-files drop zone as last grid cell in edit mode
- File type label on each tile (using `mime_type` from the API)
- Unsaved changes caption with amber dot below the edit header
- Artist external-link icon in view mode (resolved from cached `useArtists` data)
- Delete artpiece as a red text link at the bottom of edit mode
- `mime_type` added to the artpiece file response from the backend

**Non-Goals:**
- "Download all" (rendered as no-op button; deferred to a separate proposal)
- Any new backend endpoints
- Changes to save/cancel/upload logic (same API calls, same behavior)

## Decisions

### Two-column layout via CSS grid
Both view and edit modes use a CSS grid (`grid-cols-[1fr_360px]`) with details in the first column and the cover preview in the second. The cover column is fixed-width (360px) and the tile is square. This avoids the complexity of absolute positioning and works naturally with the existing flex-column shell.

### Pending-removal tiles stay in the grid (dimmed)
The current approach filters removed tiles out of `visibleFiles`. The new approach keeps them in the array but marks them with a `locallyRemovedIds` overlay — the tile renders dimmed with "Will be removed" and an Undo button. This requires changing `ArtpieceDetailFileGrid` to accept `locallyRemovedIds` in view of all files (not a pre-filtered list), and rendering removed tiles differently rather than hiding them.

### Add-files drop zone as last grid cell
The drag-strip above the grid is removed. The last cell of the edit-mode grid is always an "Add files / or drop them here" tile. The `dragOver` state moves to this cell. This keeps the upload affordance in context with the files.

### Artist link from cache, not from artpiece response
The artpiece detail API returns `artist_id` and `artist_name` but not `artist_link`. Rather than modifying the BE response, the detail component reads `artist_link` by looking up `artist_id` in the data already cached by `useArtists`. This avoids a BE change and is safe because `useArtists` is already fetched for the `ArtistCombobox` in edit mode.

### mime_type added to fileRef
`fileRef` in `artpiece_handler.go` currently omits `mime_type`. The field exists on `domain.File` and is always populated. Adding it is a two-line change and allows the frontend to derive the type label (PNG, JPG, PSD) without a separate lookup.

### Type label from mime_type
The tile label (e.g. "PNG", "JPG", "PSD") is derived from `mime_type` using a simple lookup: split on `/`, take the subtype, uppercase. For unknown subtypes, fall back to the part after the last `.` in the file name (if available), or omit the label.

### Unsaved changes caption
Derived from local state in `ArtpieceDetailView` — count of `locallyRemovedIds`, whether `pendingCoverFileId` differs from the current cover, and whether any field has changed. Renders below the edit header only when at least one change exists.

## Risks / Trade-offs

- **Undo after upload**: Files uploaded in edit mode immediately attach to the artpiece (existing behavior). If a user uploads a file and then hits Cancel, the file stays attached — this is the existing spec and is unchanged. The "Will be removed" / Undo pattern applies only to files that were already attached when edit mode was entered.
- **mime_type for legacy files**: Files uploaded before `mime_type` was stored would have an empty or incorrect value. In practice the field has been required and validated since the upload flow was built, so this is not expected to be an issue.
- **Cache freshness for artist_link**: If the user adds an artist in the same session and immediately opens an artpiece with that artist, the link will be available because `useArtists` is live-queried. No stale-data risk.
