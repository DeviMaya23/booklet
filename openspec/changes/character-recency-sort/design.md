# Design

## Context

Artists already carry `last_used_at` (migration 000029, spec `artist-recency-sort`) and verses carry one too. For artists the list endpoint always sorts by recency; the frontend does no sorting. Characters differ in one way: `GET /characters` has several consumers that need the existing `created_at DESC` order (the Characters page grid) or no particular order (id lookups in `ChangeThumbnailModal` and `useArtpiecesFilter`). Only three consumers are pickers: `ArtpieceFormModal`, `ArtpieceDetailView`, `ArtistCharacterFilter`, each using `TokenInput` with suggestions from `useCharacters()` (shared key `['characters']`). `TokenInput` filters but does not sort.

Character assignment already passes through `characterRepo.GetByIDsAndUserID` in artpiece Create and commission Create/Update. Artpiece Update is the exception: the usecase does no character validation; `artpieceRepository.Update` checks ownership itself after writing the artpiece fields.

## Goals / Non-Goals

**Goals:**
- Picker order matches artists: most recently used first.
- Everything that is not a picker keeps today's order.
- Same bump rule as artists (see `character-recency-sort`).

**Non-Goals:**
- A shared `CharacterCombobox` wrapper over the three picker call sites. The duplication stays; the cost is that a future picker must remember to pass `sort: 'recent'`.
- Exposing `last_used_at` in the API or sorting on the client.
- Changing artists or verses.

## Decisions

**1. Opt-in `sort=recent` query param instead of changing the default order.**
Add `Sort *string` (`query:"sort"`, `validate:"omitempty,oneof=recent"`) to `ListCharacterFilters`; the handler already binds it with `c.Bind`, and gains a `c.Validate` call so a bad value returns 422 (project rule: Bind + Validate). The repository switches `Order(...)` on it.
Alternatives: (a) always sort by recency, like artists — rejected, the Characters grid would reshuffle every time a character is used; (b) expose `last_used_at` and sort client-side — rejected, adds an API field and client sort logic for no benefit; (c) a separate `/characters/suggestions` endpoint — rejected, new endpoint for an ordering difference.

**2. Frontend: `useCharacters({ sort: 'recent' })`, query key `['characters', 'recent']`.**
Default call keeps `['characters']`. All existing `invalidateQueries({ queryKey: CHARACTERS_QUERY_KEY })` calls are prefix matches, so create/update/delete of characters and verse deletion already invalidate both entries. No new key constant is needed beyond building the key in the hook.
Alternative: a separate `useRecentCharacters` hook — rejected, it would duplicate the fetch function; one hook with an options arg is a smaller surface.

**3. Invalidate `['characters']` in `useCreateArtpiece`, `useUpdateArtpiece`, `useCreateCommission`, `useUpdateCommission`.**
Artists don't do this on the equivalent mutations, so their picker order can be stale until refetch. For characters we invalidate so the order is correct right after a save. Import `CHARACTERS_QUERY_KEY` from the characters feature API; cross-feature imports of query keys already exist (`useDeleteVerse`).

**4. Bump through a new `CharacterRepository.UpdateLastUsedAt(ctx, ids []uuid.UUID, userID)`.**
Unlike artists (single ID), the bump is one `UPDATE ... WHERE id IN ? AND user_id = ?` with `NOW()`. It is called from the usecase right after the existing ownership validation, before the main write — the same position as the artist bump, so a later failure of the main write still leaves characters bumped (accepted, mirrors artists).

**5. Artpiece Update validates characters in the usecase.**
When `len(params.CharacterIDs) > 0`, call `GetByIDsAndUserID` and return `ErrCharacterNotOwned` on a count mismatch, exactly like Create, then bump. The repository's own check stays as defense in depth. This is needed so a request with an unowned ID cannot bump anything, and it lets the usecase test the ownership-then-bump order.

**6. Character Create sets `LastUsedAt = now` in the usecase**, as `artist_usecase` does, so a new character appears at the top of the picker immediately.

**7. Migration `000037_add_last_used_at_to_characters`.**
Up: `ALTER TABLE characters ADD COLUMN last_used_at TIMESTAMPTZ NULL;` then `UPDATE characters SET last_used_at = created_at;`. Down: drop the column. Backfilling from `created_at` makes the recency order identical to today's order at rollout (the artist migration backfilled `NOW()` instead, which flattened ordering to name). Soft-deleted rows are backfilled too; harmless.

## Risks / Trade-offs

- [A new picker forgets `sort: 'recent'` and silently gets creation order] → accepted per the wrapper non-goal; the three call sites are listed in the proposal and covered by specs.
- [Two cached character lists (`['characters']`, `['characters','recent']`) mean two fetches when both the Characters page and a picker are used] → the list is small and per-user; invalidation by prefix keeps them consistent.
- [Bump happens before the main write, so a failed save still reorders the picker] → same behavior as artists; effect is cosmetic.
- [Re-saving an artpiece bumps every assigned character, not only newly added ones] → decided to mirror the artist rule ("non-null in the request").
- [Frontend and backend ship separately: frontend sends `sort=recent` to an old backend] → the old backend ignores unknown query params (`Bind` into a struct without the field), so it degrades to creation order. Deploy backend first.

## Migration Plan

1. Deploy backend with migration 000037 and the `sort` param.
2. Deploy frontend passing `sort: 'recent'` in the three pickers.
Rollback: revert the frontend; the down migration drops the column (the `sort` param code must be reverted with it).
