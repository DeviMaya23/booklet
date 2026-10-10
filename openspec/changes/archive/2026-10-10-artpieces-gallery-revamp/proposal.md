# Proposal

## Why

The artpieces gallery overlays labels on top of images and offers a single fixed density with no way to quickly browse artwork without navigating away from the grid. This redesign makes the gallery more image-first: clean tiles with captions below, a lightbox for flipping through covers without leaving the grid, and sticky view preferences so the layout persists across sessions.

## What Changes

- Replace `ResourceCard` in the gallery with a new `ArtpieceTile` component: square image tile, caption row (title + artist) rendered below the image rather than overlaid, hover lift effect, and an ⓘ button that appears on hover
- Single click on a tile opens `FileViewer` as a fullscreen lightbox navigating through artpieces by their cover image
- ⓘ hover button navigates to the artpiece detail view, replacing the existing double-click interaction
- **BREAKING**: Delete is removed from the gallery tile; it remains accessible only from the artpiece detail view
- Add a **View** popover to the gallery toolbar with:
  - Tile size toggle (Small / Medium / Large) — implemented as column count (Small = 6 cols, Medium = 4 cols, Large = 2–3 cols)
  - "Show details" switch that shows or hides the caption row below tiles
- Tile size and show-details preferences persist to `localStorage`
- **Backend**: `GET /artpieces` list response gains three new fields sourced from the already-preloaded `CoverFile`: `cover_file_url` (presigned GET URL, same strategy as `GET /artpieces/:id` file URLs), `cover_file_mime_type`, and `cover_file_name`

## Capabilities

### New Capabilities

_None._

### Modified Capabilities

- `web-artpieces-gallery`: Tile layout changes (caption below image, not overlaid); interaction model changes (single click → lightbox, ⓘ → detail); delete removed from tile; View popover added (tile size + show details); preferences persisted to localStorage
- `artpiece-management`: List response gains `cover_file_url`, `cover_file_mime_type`, and `cover_file_name` fields on each artpiece
- `web-file-viewer`: FileViewer is also opened from the gallery lightbox (one `ViewerFile` entry per artpiece cover), not only from the artpiece detail file grid

## Impact

- **Frontend**
  - `frontend/src/features/artpieces/components/ArtpiecesGrid.tsx` — rebuilt with column-count state, View popover, and tile size / show-details controls
  - `frontend/src/features/artpieces/components/ArtpieceTile.tsx` — new component
  - `frontend/src/pages/ArtpiecesPage.tsx` — FileViewer state lifted here; double-click navigation removed; no tile-level delete
  - `frontend/src/features/artpieces/api/useArtpieces.ts` — `ArtpieceSummary` gains `cover_file_url`, `cover_file_mime_type`, `cover_file_name`
  - `frontend/src/components/ResourceCard.tsx` — no change (still used elsewhere)
- **Backend**
  - `backend/internal/handler/artpiece_handler.go` — `artpieceResponse` and `ListArtpieces` gain cover file URL presigning
- **No schema migrations, no new routes, no new dependencies**
