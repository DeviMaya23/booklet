# Design

## Context

The backend has a complete, routed `saved_filters` CRUD layer (`SavedFilterHandler`, usecase, repository) with no frontend wiring. The existing `ArtistCharacterFilter` component is shared across three callers; `ArtpiecesPage` is the only one that needs the save action. Filter state in `useArtpiecesFilter` is plain React state with no URL or storage persistence. See proposal.md for motivation.

Key existing constraints:
- `SavedFilterHandler` currently has no `Presigner` or artpiece dependency.
- `savedFilterResponse` exposes `thumbnail_r2_path` (raw R2 key) — not a usable URL for the frontend.
- `ArtpieceSummary` (artpiece list response) exposes `thumbnail_url` (presigned URL) but not the raw R2 path.
- `ArtistCharacterFilter` encapsulates the popover entirely; the footer is not composable by callers today.

## Goals / Non-Goals

**Goals:**
- Wire the create saved filter flow end-to-end from the artpiece gallery.
- Add a `footerAction` slot to `ArtistCharacterFilter` without changing its two other callers.
- Persist artpiece gallery filter state in URL params.
- Update the BE response contract: all saved filter responses return `thumbnail_url` (signed) instead of the raw `thumbnail_r2_path`.

**Non-Goals:**
- Dashboard display of saved filters (Part 2).
- Deep-link navigation from the dashboard to the gallery with a prefilled filter (Part 2).
- Saved filter list, edit, and delete UI.
- Changing the `PATCH /saved_filters/:id` request shape — the patch thumbnail field stays `thumbnail_r2_path` until Part 2 wires the edit UI.

## Decisions

### 1. Thumbnail: BE resolves R2 path from `thumbnail_artpiece_id`

**Decision**: The create request accepts an optional `thumbnail_artpiece_id` (UUID string). The handler looks up the artpiece, reads `CoverFile.ThumbnailR2Path`, and stores it. The handler gains `ArtpieceRepository` (just the `GetByID` operation) and `Presigner` dependencies. All saved filter responses sign the stored path and return `thumbnail_url`.

**Why over exposing `thumbnail_r2_path` on `ArtpieceSummary`**: The R2 path is an internal storage detail. Keeping it out of the artpiece list response follows the same pattern as the dashboard handler, which also signs paths server-side and never exposes raw keys.

**Artpiece ownership check**: The handler looks up the artpiece without verifying the user owns it. In practice, the FE only surfaces the authenticated user's own artpieces in the thumbnail picker, so the risk of a user providing a foreign artpiece ID is low. An explicit ownership check can be added in Part 2 if needed.

### 2. Unique name constraint: DB + BE 422

**Decision**: Add a DB migration with a `UNIQUE (user_id, name)` constraint on `saved_filters`. The repository layer catches the unique violation (Postgres error code `23505`) and returns a typed `ErrDuplicateSavedFilterName` error. The usecase passes it through; the handler maps it to a 422 with a `name` field-level error.

**Why not FE-only check**: A FE check against a prefetched list has a race condition (another tab or device saves first) and requires an extra round-trip before the modal can render. BE constraint is authoritative and consistent with how other field validation errors work in this codebase.

### 3. `ArtistCharacterFilter`: optional `footerAction?: ReactNode` slot

**Decision**: Add a single optional prop `footerAction?: ReactNode` to `ArtistCharacterFilter`. When present, it is rendered in the existing footer row to the right of "Clear all". No logic change; the component remains purely presentational about what goes in that slot.

**Why not a compound pattern or dedicated prop**: The slot is the minimal change that satisfies the spec without affecting the other two callers (`CommissionFormModal`, `AddToExistingArtpieceModal`) or requiring a component refactor.

### 4. Filter state: URL search params via `useSearchParams`

**Decision**: Refactor `useArtpiecesFilter` to sync with `useSearchParams`. Params: `artist` (artist UUID), `characters` (comma-separated character UUIDs), `match` (`all` | `any`). On mount, the hook reads params and hydrates the artist/character objects from the already-cached `useArtists` / `useCharacters` query data.

**Hydration timing**: If the artist/character queries haven't resolved yet on mount, the filter UI will initialise to its empty state and populate once the data arrives. This is acceptable — the queries are fast and cached. An alternative of waiting for both queries before rendering would add visual delay with no real benefit.

**Why over localStorage**: URL params enable the Part 2 dashboard deep-link (navigate to `/app/artpieces?artist=<id>&...`) with no further changes to the filter state mechanism.

### 5. "Saved as X" state: ephemeral snapshot in `ArtpiecesPage`

**Decision**: After a successful save, `ArtpiecesPage` stores `{ name: string; snapshot: FilterSnapshot }` in local state. `FilterSnapshot` captures the artist ID and character IDs active at save time. The "Saved as X" indicator renders while the current filter values match the snapshot; the "Save to dashboard" button reappears as soon as values diverge.

**Why not a query or stored value**: The "Saved as X" state is purely cosmetic confirmation — no round-trip needed.

### 6. Filter payload schema

The `filter_payload` JSON blob sent to `POST /saved_filters` SHALL use the shape:

```json
{
  "artist_id": "<uuid> | null",
  "character_ids": ["<uuid>", ...],
  "match_mode": "all | any"
}
```

The backend stores this as-is without validation. The frontend owns the schema. This matches the current single-artist FE model; the shape can be revised when multi-artist support is added.

## Risks / Trade-offs

- **Stale thumbnail**: The saved filter stores the R2 path at creation time. If the source artpiece is later deleted and its file purged from R2, the presigned URL will be generated but the image will 404. The FE should render a placeholder when `<img>` fails to load. This is called out in the proposal as an accepted edge case.
- **URL param hydration race**: On first load, if `useArtists`/`useCharacters` data isn't yet in cache, the filter will appear empty for a brief moment before filling. This is cosmetically acceptable; the query results arrive quickly on a warm session.
- **`thumbnail_r2_path` in PATCH request**: The PATCH endpoint still accepts the old `thumbnail_r2_path` field name in its request (not wired to UI yet). This leaves the create and patch endpoints with inconsistent thumbnail fields until Part 2. It is a known asymmetry scoped to this proposal.
