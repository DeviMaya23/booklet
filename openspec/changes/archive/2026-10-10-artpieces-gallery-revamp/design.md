# Design

## Context

The gallery currently uses `ResourceCard` which renders title/artist as a gradient overlay on the image. `ArtpiecesGrid` has a hardcoded breakpoint column grid and no view controls. The page navigates away to `/app/artpieces/:id` on double-click and exposes delete via a `...` dropdown on each tile.

`FileViewer` already accepts `ViewerFile[]` with `id`, `name`, `mimeType`, `thumbnailUrl`, `previewUrl`, `width?`, `height?` — it is fully data-driven and reusable without modification.

`GET /artpieces` already preloads `CoverFile` in the repository layer. `ListArtpieces` currently calls `GenerateDeterministicPresignedGetURL` for the thumbnail path. The cover file's `FileR2Path`, `MimeType`, and `Name` are already available in the loaded `CoverFile` association — no additional SQL queries are needed.

## Goals / Non-Goals

**Goals:**
- New `ArtpieceTile` component that renders caption below the image, with hover lift and ⓘ button
- Gallery view controls (tile size, show-details) with localStorage persistence
- Single-click → FileViewer lightbox over artpiece covers; ⓘ click → detail navigation
- Backend: three new fields (`cover_file_url`, `cover_file_mime_type`, `cover_file_name`) on the list response
- Reuse `FileViewer` as-is; no changes to the component itself

**Non-Goals:**
- Changing `ResourceCard` (still used in other contexts)
- Modifying `FileViewer` internals
- Pagination or virtual scrolling for the gallery
- Server-side view preference storage

## Decisions

### New `ArtpieceTile` vs. extending `ResourceCard`

**Decision**: Create a new `ArtpieceTile` component rather than adding props to `ResourceCard`.

`ResourceCard` places its caption inside the image via an absolute overlay. Moving the caption outside the card boundary requires restructuring the DOM (the wrapping element must become a `<figure>` that includes a sibling caption row). Adding that to `ResourceCard` via a prop would make it context-aware and increase its surface area. `ArtpieceTile` starts clean and is scoped to this one use case.

**Alternative**: Add `captionBelow` / `showCaption` / `onInfoClick` props to `ResourceCard`. Rejected — `ResourceCard` is a generic component used in at least one other location (commission artpiece picker); coupling gallery-specific behavior there is brittle.

---

### Cover URL presigning strategy for the list

**Decision**: Use `GeneratePresignedGetURL` (non-deterministic, TTL-based) for `cover_file_url` in the list response, matching what `GetArtpieceByID` does for individual file URLs.

The cover file in the list is used as the lightbox preview URL (`previewUrl` in `ViewerFile`). This is consistent with how `GetArtpieceByID` presigns its file array. The existing thumbnail URL already uses the deterministic presigner for caching — combining both strategies on one endpoint is intentional: the thumbnail is cacheable (small, frequently re-fetched for cards), the full-res URL is one-time use (opened in the lightbox).

**Alternative**: Use `GenerateDeterministicPresignedGetURL` for the cover URL too. Rejected — the full-res cover file changes less predictably than thumbnails, and the TTL-based URL is the established pattern for file preview URLs in this codebase.

---

### `viewerIndex` state location: `ArtpiecesPage` vs. `ArtpiecesGrid`

**Decision**: Lift `viewerIndex` state and the `FileViewer` render into `ArtpiecesPage`.

`ArtpiecesPage` already owns the `filtered` artpieces array (from `useArtpiecesFilter`). The `ViewerFile[]` array passed to `FileViewer` must be built from that filtered list so prev/next navigation stays coherent with what the user sees in the grid. Keeping that logic in `ArtpiecesPage` avoids threading the full filtered list down through `ArtpiecesGrid`.

---

### localStorage key and structure

**Decision**: Two separate keys: `gallery.tileSize` (`"small" | "medium" | "large"`) and `gallery.showDetails` (`"true" | "false"`).

Simple string values avoid JSON parsing. Defaults applied when keys are absent or unreadable (Medium, true). Reads are wrapped in try/catch per the project's localStorage conventions.

---

### Removing delete from the tile

**Decision**: Remove the `...` dropdown and `DeleteArtpieceDialog` wiring from `ArtpiecesGrid` entirely. Delete remains in `ArtpieceDetailView`.

The `...` dropdown currently only contains "Delete". With delete moved to the detail view, the dropdown has no items — keeping it as a shell makes no sense. The ⓘ button provides a clear, discoverable path to the detail view where delete is available.

## Risks / Trade-offs

- **[Risk] Cover URL expiry during long sessions** → The presigned `cover_file_url` has a fixed TTL. If the user leaves the gallery open for longer than the TTL and then opens the lightbox, the URL may be expired. Mitigation: the TTL used for file GET URLs in this codebase (`usecase.PresignGetTTL`) is long enough to cover normal sessions; the list query refetches on focus anyway (React Query default).

- **[Risk] Large cover files in lightbox** → Full-resolution cover files can be large. The lightbox loads the full file URL, not a resized version. Mitigation: this is an intentional trade-off for image fidelity; the thumbnail still loads for the tile, and the full file only loads when the lightbox is opened.

## Migration Plan

Backend and frontend changes are independently deployable:
1. Deploy backend first: the new `cover_file_url`, `cover_file_mime_type`, `cover_file_name` fields are additive and ignored by the existing frontend.
2. Deploy frontend: the revamped gallery reads the new fields; old behaviour (double-click, tile delete) is replaced.

No rollback complexity — the backend change is additive (no field removals, no schema changes). Reverting the frontend restores the previous gallery with no data loss.
