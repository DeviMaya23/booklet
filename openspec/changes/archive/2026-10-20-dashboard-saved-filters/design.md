# Design

## Context

See `proposal.md — Why` for motivation.

Key existing constraints:
- `useArtpiecesFilter` reads filter state from URL params (`artist`, `characters`, `match`) — the deep-link mechanism from dashboard to gallery is already wired; this change just needs to produce the right URL.
- `filter_payload` is stored as raw JSON (`json.RawMessage`) by the backend with no schema validation. The frontend owns the shape.
- `ArtpieceSummary` exposes `thumbnail_url` (presigned) but not the raw R2 key. The backend must resolve `thumbnail_artpiece_id` server-side for PATCH, same as for CREATE.
- `ResourceCard` (characters) uses a `div onClick` for its primary action (opens edit modal). The new `ThumbnailTile` needs navigation semantics — a `<Link>` — which is why they stay separate components.
- `useArtpieces` is a globally cached React Query query; calling it inside `ChangeThumbnailModal` incurs no extra network round-trip on a warm session.

## Goals / Non-Goals

**Goals:**
- Display saved filters on the dashboard as clickable thumbnail tiles.
- Wire click-to-navigate for recent artpieces tiles.
- Introduce `ThumbnailTile` as the shared component for navigation-first tiles.
- Per-tile delete and change-thumbnail actions.
- `localStorage`-backed visibility toggle.
- Close the `PATCH /saved_filters/:id` thumbnail asymmetry by accepting `thumbnail_artpiece_id`.

**Non-Goals:**
- Artpiece count per saved filter tile (no count display).
- Reordering saved filters on the dashboard.
- Changing `ResourceCard` (characters remain on their own component).
- Any new backend field for visibility state.

## Decisions

### 1. `ThumbnailTile`: `<Link>` wraps the image area, `...` button sits outside it

**Decision**: The image area is wrapped in a React Router `<Link to={href}>`. The `...` dropdown trigger is positioned absolutely inside the container div but outside the `<Link>`, so no `stopPropagation` is needed.

```
<figure>
  <div class="relative aspect-square ...">
    <Link to={href}>          ← image area only
      <img /> | placeholder
    </Link>
    <div class="absolute top-1 right-1">   ← outside Link
      <DropdownMenu>...</DropdownMenu>      ← rendered only when menuItems present
    </div>
  </div>
  <figcaption>label + optional sublabel</figcaption>
</figure>
```

**`menuItems` prop shape:**
```typescript
interface MenuAction {
  label: string
  icon?: LucideIcon
  variant?: 'destructive'
  onClick: () => void
}
```
Absent or empty → `...` button is not rendered at all.

**Why not `ResourceCard` with a `menuItems` generalization**: `ResourceCard`'s primary action is an edit modal (`onClick` on a `div`). `ThumbnailTile`'s primary action is navigation (a `<Link>`). Merging them would require a conditional rendering of Link vs. div based on a prop, which blurs the component's contract. Two focused components are cleaner.

### 2. `ThumbnailTile` placement: `components/`

`ThumbnailTile` is shared across two dashboard sections (recent artpieces and saved filters). It is not feature-specific, so it lives at `components/ThumbnailTile.tsx` alongside `ResourceCard`.

### 3. Filter payload: defensive parse utility in `features/artpieces/lib/`

`filter_payload` from the API is `json.RawMessage` serialized as a plain object. The FE parses it defensively:

```typescript
function parseFilterPayload(raw: unknown): FilterPayload | null
```

Returns `null` on any shape mismatch. The caller maps a `null` result to an empty navigation (no URL params). This utility is colocated with the filter logic in `features/artpieces/lib/`.

The URL param mapping:
- `artist_id` (non-null string) → `artist=<id>`
- `character_ids` (non-empty array) → `characters=<id1,id2>`
- `match_mode === 'any'` → `match=any` (omit if 'all', matching `useArtpiecesFilter` convention)

### 4. Visibility state: `localStorage` key `dashboard.hiddenFilterIds`

Stored as a JSON array of string IDs. Absent key = show all. A filter not in the array is visible; adding its ID hides it.

**Why store hidden IDs rather than shown IDs**: New filters are visible by default with no action required. Storing shown IDs would require updating localStorage every time a filter is created.

**Stale IDs**: If a filter is deleted, its ID remains in the hidden set but is harmless — the filter no longer appears in the API response.

### 5. `ChangeThumbnailModal`: uses cached `useArtpieces` + FE-side filter

The modal calls `useArtpieces()` internally (returns cached data instantly on a warm session) and applies the saved filter's parsed payload with `filterAndSortArtpieces` (already exported from `useArtpiecesFilter`). Up to 10 results are shown in the picker grid. First artpiece is pre-selected.

**Why not a dedicated filtered endpoint**: The artpieces list is already cached globally. An extra API call would add latency and a loading state for data already in memory.

### 6. `PATCH /saved_filters/:id`: handler resolves `thumbnail_artpiece_id`, usecase unchanged

`patchSavedFilterRequest` gains `ThumbnailArtpieceID *string`. When set, the handler looks up the artpiece (same `ArtpieceRepository.GetByID` call as CREATE) and populates `ThumbnailR2Path` on the `UpdateSavedFilterParams`. The usecase layer does not change. When both `thumbnail_artpiece_id` and `thumbnail_r2_path` are sent (unlikely in practice), `thumbnail_artpiece_id` takes precedence.

### 7. New FE hooks location: `features/artpieces/api/`

`useSavedFilters`, `useDeleteSavedFilter`, and `useUpdateSavedFilter` join `useCreateSavedFilter` in `features/artpieces/api/`. Saved filter operations are conceptually tied to the artpiece filter feature; there is no separate `saved-filters` feature directory.

## Risks / Trade-offs

- **Stale thumbnail after artpiece deletion**: The saved filter stores an R2 path at save/patch time. If the source artpiece is later deleted and purged from R2, `thumbnail_url` will be generated but the image will 404. `ThumbnailTile` (and `ChangeThumbnailModal` picker cells) should render the placeholder on `<img>` load error. `FileThumb` already handles this pattern with an `onError` handler; follow the same approach.
- **`localStorage` unavailability**: Private browsing or blocked storage will throw on access. Wrap all reads/writes in try/catch and fall back to showing all filters.
- **`useArtpieces` not yet cached on dashboard cold load**: If the user lands directly on the dashboard without having visited the artpieces page, `useArtpieces` will be unfetched when `ChangeThumbnailModal` opens. The modal will show a loading state briefly. This is cosmetically acceptable.
