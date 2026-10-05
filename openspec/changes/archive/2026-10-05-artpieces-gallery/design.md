## Context

Artpieces exist in the backend and can be created from the file inbox, but there is no dedicated screen to browse them. All the data needed (`thumbnail_url`, `artist_name`, `characters[]`, `created_at`) is already returned by `GET /artpieces`. The frontend has `useArtpieces()`, `ArtpieceFormModal`, `TokenInput`, `ResourceCard`, and the `Combobox` component — all reusable without modification (except `ResourceCard` gaining a sublabel slot).

## Goals / Non-Goals

**Goals:**
- New `/app/artpieces` route and sidebar nav item
- Gallery screen with client-side search, sort, and filter
- Reuse existing `ArtpieceFormModal` for the "+ New Artpiece" entry point
- Minimal surface area: no new API calls, no new shared components beyond `ResourceCard` sublabel

**Non-Goals:**
- Server-side filtering or pagination (explicit deferral; revisit if collection grows large)
- Artpiece detail/edit view (out of scope for this pass)
- Delete from gallery (out of scope for this pass)

## Decisions

### Client-side filter and sort
The backend already supports `artist_ids` and `character_ids` query params on `GET /artpieces`, but loading all artpieces upfront and filtering in-memory is appropriate for a small personal collection. This avoids query param management, debouncing, and loading state for each filter change. Revisit when collection size becomes a concern.

### Artist filter: single-select Combobox
The filter panel uses a single-artist Combobox (same component and pattern as `ArtpieceFormModal`). Multi-select artist filter was considered but ruled out — the schema constraint (`artpieces.artist_id` is a single FK) makes selecting multiple artists a straightforward OR that is indistinguishable from "any of these", and single-select covers the main use case cleanly with zero new component work.

### Character filter: `TokenInput` with All/Any toggle
`TokenInput` already handles chip-based multi-select with keyboard and suggestion support and is used in `ArtpieceFormModal` for the same data shape (`{ id, name }`). Suggestions come from the already-fetched `useCharacters()` query (no extra network call). The All/Any toggle is a simple boolean in local state — "All" means the artpiece must include every selected character ID in its `characters[]` array; "Any" means at least one.

### Filter popover: floating div with click-outside handler
No `Popover` component exists in the UI library. Rather than introducing a new shadcn component or base-ui primitive, the filter panel is implemented as a controlled `div` with `position: absolute`, toggled by the Filter button, and dismissed on outside click via a `useEffect` with a document `mousedown` listener. This keeps the filter self-contained in a single `ArtpiecesFilterPopover` component and adds no new dependency.

### `ResourceCard` sublabel
Add an optional `sublabel?: string` prop. When present, render it as a second, smaller line below the existing label in the gradient overlay. Existing callers (`ImagesGrid`) pass no `sublabel` and are unaffected.

### Component placement
- `ArtpiecesPage` → `pages/ArtpiecesPage.tsx` (owns all filter/sort/search state)
- `ArtpiecesGrid` → `features/artpieces/components/ArtpiecesGrid.tsx` (pure rendering, takes filtered list)
- `ArtpiecesFilterPopover` → `features/artpieces/components/ArtpiecesFilterPopover.tsx` (self-contained filter panel)

## Risks / Trade-offs

- **Large collections**: Loading all artpieces upfront will degrade with scale. Mitigation: document the deferral explicitly; the switch to server-side filtering is a contained change when needed (update `useArtpieces` to accept filter params, pass them from the page).
- **Null titles in search/sort**: `title` is nullable. Mitigation: null titles are excluded from search results (not matched) and sorted last in alphabetical order. Display falls back to "Untitled".
- **Filter popover and keyboard/a11y**: A hand-rolled floating div has less a11y than a proper Popover primitive. Mitigation: acceptable for an internal tool at this stage; noted as a known gap.
