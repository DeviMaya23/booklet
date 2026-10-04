## Context

`FileTile` currently renders as a pure image grid card with a `···` hover dropdown containing "View Detail" (no-op) and "Delete". There is no filename visible on the tile, no editing surface, and `onDoubleClick` is already wired to a no-op in `FileInboxGrid`. The grid reflows when any containing element changes width, which rules out fixed side panels that push content.

## Goals / Non-Goals

**Goals:**
- Show filename as a truncated read-only label beneath each tile
- Replace `···` dropdown with a direct hover trash icon (single Delete action)
- Wire double-click to open a per-file edit overlay that covers the grid without reflowing it
- Auto-save name and notes on blur; surface save failures inline (same pattern as `AddFilesModal`)

**Non-Goals:**
- File preview or full detail view beyond name and notes
- Multi-file editing
- Any backend changes

## Decisions

### Decision: Overlay is fixed-position over the grid, not a side panel that pushes content

A panel that resizes the grid causes tiles to reflow, disorienting the user's spatial reference. The overlay instead renders at a fixed position over the right portion of the grid (similar to a mobile sheet), dimming tiles behind it without moving them. The clicked tile remains visible in the grid; its image is also shown in the overlay.

**Considered**: Collapsible fixed panel (always present, zero-width when closed). Rejected — empty state is a passive cost on every session even when the feature is unused.

### Decision: Overlay state lives in `FileInboxGrid`, not `FilesPage`

The overlay is tightly coupled to tile interaction (double-click), which already lives in `FileInboxGrid`. Lifting state to `FilesPage` would require threading `onDoubleClick` and overlay open/close props through an extra layer with no benefit — `FilesPage` has no interest in which file is being edited.

### Decision: Auto-save on blur, same pattern as `AddFilesModal`

The inbox is a dump — editing is quick and incidental. A Save/Cancel flow adds unnecessary friction. Name and notes each save independently on blur via `useUpdateFile`. Save failures show a red border + inline error message; success is silent. This is identical to `AddFilesModal`'s field-save behavior.

### Decision: Trash icon replaces `···` dropdown on the tile

The dropdown had two items: "View Detail" (removed) and "Delete". A single-item dropdown is worse UX than a direct icon. A trash icon on hover is more discoverable, takes less space, and matches the directness of the inbox's purpose.

### Decision: `FileTileEditOverlay` is a new standalone component

It shares the auto-save + inline error pattern with `AddFilesModal` but differs in context: it shows a single existing file (no upload flow, no list). The component receives the file object and an `onClose` callback; it manages its own local name/notes state and error state.

## Risks / Trade-offs

- **Overlay covers grid tiles** — tiles behind the overlay are partially obscured. Acceptable because the user explicitly double-clicked to open it, and the edited file is displayed inside the overlay.
- **Trash icon discoverability** — hover-only affordance is invisible on touch devices. Acceptable given the inbox is a desktop-first surface; touch users can still use long-press/context menu if that is ever added.
- **`useFiles` thumbnail sync** — the overlay reads `thumbnail_url` from the file prop passed in. If the thumbnail resolves while the overlay is open, it won't update unless we subscribe to `useFiles`. For simplicity, the overlay shows whatever `thumbnail_url` the file had at open time — thumbnail generation is typically fast and this is an edit surface, not a preview.
