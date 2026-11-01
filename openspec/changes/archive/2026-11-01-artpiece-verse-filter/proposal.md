# Proposal

## Why

Artpieces have no verse of their own — a verse is only reachable through the characters tagged on an artpiece. Users with large libraries spanning many fandoms/OCs need a way to narrow the artpieces gallery to a single verse without manually selecting every character in it. The verse data model and character-verse assignment already exist; the artpieces gallery just has no filter dimension for it yet.

## What Changes

- Add a single-select verse filter to the artpieces gallery filter panel (`ArtistCharacterFilter`), alongside the existing artist and character filters.
- An artpiece matches the verse filter if at least one of its tagged characters has that `verse_id` — computed entirely client-side by cross-referencing the already-fetched `useCharacters()` data (which already returns `verse_id` per character). No backend query-param or response-shape change.
- Verse filter state follows the same lifecycle as the existing artist/character filters: synced to the URL (`?verse=<uuid>`), persisted to `localStorage` under the existing `artpieces.filterState` key, AND-combined with the other active filters, and included in the filter summary label/chip.
- Saved filters gain a `verse_id` field in `filter_payload` (nullable, defaults to `null` when absent on read — same precedent as `is_reference`). Saving, dashboard-tile navigation, and the "Change thumbnail" picker all carry the verse filter through.
- If a saved filter's `verse_id` no longer resolves against the user's current verses (e.g. the verse was deleted), the filter is treated as not-selected for that dimension — mirroring how a deleted artist/character already degrades today.

## Capabilities

### New Capabilities
(none — this extends existing filter infrastructure)

### Modified Capabilities
- `web-artist-character-filter`: the shared `ArtistCharacterFilter` component gains an optional single-select verse field (`selectedVerse`/`onVerseChange`), following the same optional-prop pattern as the existing "References only" toggle so non-gallery callers (commission form, add-to-artpiece modal) are unaffected unless they opt in.
- `web-artpieces-gallery`: the filter popover, its active-filter summary label, combined-filter semantics, URL sync, and `localStorage`-persisted structural filter state all extend to cover the verse filter.
- `web-save-filter-modal`: the filter summary chip includes the verse when active.
- `web-dashboard-saved-filters`: clicking a saved-filter tile maps `filter_payload.verse_id` to the `verse` URL param (alongside the existing `artist`/`characters`/`match` mapping); the "Change thumbnail" picker applies the verse filter when computing its candidate grid.
- `saved-filters`: `filter_payload` gains a `verse_id` field (nullable UUID string or absent), with absent-on-read treated as `null` — same shape and validation posture as the existing `is_reference` field (generic JSON object/size validation only, no field-specific backend check).

## Impact

All changes are frontend-only; no backend code, migration, or API contract changes.

- `frontend/src/features/artpieces/hooks/useArtpiecesFilter.ts` — new `verseId` filter dimension, URL param, `StoredFilter` field, and the character→verse cross-reference in `filterAndSortArtpieces`
- `frontend/src/components/ArtistCharacterFilter.tsx` — new optional verse combobox field (reusing existing `VerseCombobox`)
- `frontend/src/pages/ArtpiecesPage.tsx` — wires the new filter state into the gallery's `AristCharacterFilter` usage and `SaveFilterModal` filter state
- `frontend/src/features/artpieces/lib/parseFilterPayload.ts` — `FilterPayload` type and parsing gains `verse_id` with default-null-on-absence
- `frontend/src/features/artpieces/api/useCreateSavedFilter.ts` — `CreateSavedFilterParams`/`SavedFilter` payload types gain `verse_id`
- `frontend/src/features/artpieces/components/SaveFilterModal.tsx` — payload construction and `buildFilterSummary` include verse
- `frontend/src/features/artpieces/components/ChangeThumbnailModal.tsx` — needs `useCharacters()` wired in (not currently imported) so its independent `filterAndSortArtpieces` call can apply the verse filter; `filterChipLabel` includes verse
- `frontend/src/pages/DashboardPage.tsx` — `filterToHref()` maps `verse_id` to the `verse` URL param
