# Design

## Context

See `proposal.md` - Why. Relevant current state:

- `characters.verse_id` is a single nullable FK (`ON DELETE SET NULL`) — a character belongs to at most one verse. `GET /characters` already returns `verse_id` per character.
- `GET /artpieces` returns each artpiece's `characters` as `{id, name}[]` — no `verse_id`. Adding it would require a backend repository/DTO change.
- All artpieces-gallery filtering (`frontend/src/features/artpieces/hooks/useArtpiecesFilter.ts`) is client-side: `useArtpieces()` fetches the full unfiltered list and `filterAndSortArtpieces` filters in memory. `useCharacters()` is already fetched on the same page for the existing character filter.
- Saved filters (`saved_filters.filter_payload`) are validated backend-side only generically (JSON object, ≤4096 bytes) — no field-specific schema enforcement. `is_reference` was added to this payload purely as an FE-defined convention with a documented default-on-absence; `verse_id` follows the same precedent.
- The filter payload shape and its construction/parsing is duplicated across four FE call sites, not centralized: `parseFilterPayload.ts`, `useCreateSavedFilter.ts`'s types, `SaveFilterModal.tsx`'s save handler, and `DashboardPage.tsx`'s `filterToHref()`. `ChangeThumbnailModal.tsx` independently re-invokes `filterAndSortArtpieces` and currently only imports `useArtists()`, not `useCharacters()`.

## Goals / Non-Goals

**Goals:**
- Filter the artpieces gallery to artpieces with at least one character in a selected verse, using only data already fetched on the page.
- Carry the verse filter through every existing filter-persistence surface (URL, localStorage, saved filters, dashboard tile navigation, thumbnail picker) with the same fidelity as the existing artist/character filters.
- Keep the shared `ArtistCharacterFilter` component backward-compatible for its other callers (commission form modal, add-to-artpiece modal) that don't filter by verse.

**Non-Goals:**
- No backend query-param (`verse_id` on `GET /artpieces`) or response-shape change. If gallery filtering is ever moved server-side, that's a separate change.
- No multi-select verse filter. A character has at most one verse, so a single-select verse filter mirrors the existing single-select artist filter structurally; multi-select would add UI complexity (an all/any toggle) without a corresponding data-model reason.
- No explicit "characters with no verse" filter option. Absence of a verse is not itself a filterable category in this change.

## Decisions

**1. Client-side cross-reference, not an embedded `verse_id` on artpiece characters.**
`useCharacters()` is already fetched on the gallery page for the character filter and already includes `verse_id`. The verse filter builds a `characterId -> verse_id` lookup from that data and checks artpiece membership against it, with no new network request and no backend change.
- Alternative considered: have the backend preload `Characters.Verse` and add `verse_id` to each `characterRef` in the artpiece list response. Rejected for this change because it's unneeded backend surface area for a filter that the client can already compute from data it has; revisit if a future feature needs verse-filtered *server-side* pagination/sorting.

**2. Single-select verse filter, structured like the existing artist filter.**
Reuses the existing `VerseCombobox` component as-is (no new multi-select variant needed). State shape in `useArtpiecesFilter.ts` mirrors `artistId: string | null`, not `characterIds: string[]`.

**3. `verse_id` as an optional, defaulted field on `filter_payload`, not a backend schema change.**
Matches the `is_reference` precedent exactly: FE-defined field, absent-on-read defaults to `null`, no backend validation beyond the existing generic object/size checks. Keeps the backend saved-filters contract stable.

**4. Resolve saved/stored verse ids against live `useVerses()` data; treat a miss as not-selected.**
Mirrors how `selectedArtist`/`selectedCharacters` already resolve against `useArtists()`/`useCharacters()` in `useArtpiecesFilter.ts` — a deleted artist or character already silently drops out of an old saved filter's effective state today. Verse follows the same rule rather than introducing a new "unresolvable reference" UX pattern.

## Risks / Trade-offs

- **Four-plus duplicated payload call sites** → each must be updated in lockstep (listed in proposal.md - Impact). Mitigation: no new abstraction introduced in this change (per project decision-boundary convention, introducing a shared payload builder would be a new abstraction requiring separate sign-off), but tasks.md will enumerate every call site explicitly so none is missed silently.
- **`ChangeThumbnailModal.tsx` doesn't currently import `useCharacters()`** → without adding it, the thumbnail picker would silently ignore the verse filter while the main gallery respects it. Mitigation: explicit task to add the import and pass character data into its `filterAndSortArtpieces` call.
- **Dangling `verse_id` in an old saved filter after the verse is deleted** → filter silently degrades to "no verse filter" rather than erroring, consistent with existing artist/character behavior; no user-facing indication that a referenced verse no longer exists. Accepted as consistent with current product behavior, not a regression introduced by this change.

## Migration Plan

Frontend-only, additive change. No data migration. Existing saved filters (missing `verse_id`) continue to work unchanged via the default-to-`null` parsing rule. No feature flag needed — ships as a normal FE release; rollback is a plain revert since no backend or data shape is touched.
