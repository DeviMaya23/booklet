## Why

The inbox file tiles expose no metadata and the only action beyond Delete ("View Detail") is a no-op. Users need a low-friction way to see a file's name at a glance and edit its name and notes without leaving the inbox flow.

## What Changes

- File name is displayed as a truncated read-only label beneath each tile in the inbox grid
- The `···` per-tile dropdown menu is removed and replaced with a trash icon that appears on tile hover (Delete action only)
- Double-clicking a tile opens an edit overlay — a panel that slides over the right portion of the grid without causing a reflow — showing a thumbnail, editable name field, and editable notes field; fields auto-save on blur
- "View Detail" menu item is removed (was a no-op; replaced by double-click)

## Capabilities

### New Capabilities

- `web-file-inbox-tile-edit`: The double-click edit overlay — per-file name and notes editing directly from the inbox grid, with auto-save on blur and inline error feedback on save failure

### Modified Capabilities

- `web-file-inbox`: Tile visual updated (filename label below tile, trash icon replaces `···` menu); double-click interaction added

## Impact

- `frontend/src/features/files/components/FileTile.tsx` — add filename label, replace `···` dropdown with hover trash icon, wire `onDoubleClick`
- New component: `frontend/src/features/files/components/FileTileEditOverlay.tsx`
- `frontend/src/features/files/components/FileInboxGrid.tsx` — manage overlay open state, pass `onDoubleClick` to tiles
- `frontend/src/pages/FilesPage.tsx` — no changes expected
- Reuses existing hooks: `useUpdateFile`, `useFiles` (for thumbnail display in overlay)
- No backend changes required
