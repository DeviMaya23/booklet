# Design

## Context

See proposal.md — Why.

`useArtpiecesFilter` already manages structural filter state entirely through URL search params (`useSearchParams`). The hook derives artist/characters/match/isReference from URL params and exposes setters that call `setSearchParams`. `ArtpiecesPage` consumes the hook with no knowledge of how state is stored.

The gallery already uses localStorage for two display preferences (`gallery.tileSize`, `gallery.showDetails`) following a read-at-init / write-on-change pattern. That same pattern is reused here.

## Goals / Non-Goals

**Goals:**
- Structural filter state (artist, characters, match, isReference) survives navigating away and returning
- Arriving from a dashboard saved-filter link updates the stored state so that subsequent sidebar returns land on that filter

**Non-Goals:**
- Persisting `search` or `sort`
- Any UI chrome for this feature (no "filter remembered" indicator)
- Multi-tab sync or cross-device persistence

## Decisions

### Decision: Persistence inside `useArtpiecesFilter`, not in `ArtpiecesPage`

The hook already owns all filter state and the setters that mutate it. Keeping localStorage access there avoids leaking storage concerns into the page and keeps `ArtpiecesPage` unchanged.

_Alternative considered:_ Manage localStorage in `ArtpiecesPage` and pass initial values into the hook. Rejected because the hook already knows the URL params at construction time; the page would have to know about the key names and the `StoredFilter` shape, adding unnecessary coupling.

### Decision: Two-path mount effect — restore or write

A single `useEffect` with an empty dep array handles both scenarios at mount:

- **URL has filter params** (e.g. arrived from a saved-filter link): write the current URL state to localStorage. This updates "last used" to match the filter the user just arrived with.
- **URL has no filter params** (e.g. arrived from sidebar): read localStorage; if meaningful state exists, call `setSearchParams` with it (`replace: true`).

A mount-only effect (empty deps) is the right boundary: the "check-on-arrival" decision is a one-time action, not a reaction to ongoing state changes. The `react-hooks/exhaustive-deps` eslint exception is appropriate here and is already allowed by the project's frontend conventions.

_Alternative considered:_ A sync effect that watches the URL params and continuously writes to localStorage. This creates an effect-ordering race: the sync effect fires on the first render with empty params (before the restore effect can set them), overwriting any stored state with nothing. A single mount effect eliminates the ordering problem.

### Decision: Each setter commits the whole filter, not a delta of the URL params

`setSearchParams`'s updater form does **not** chain: React Router hands it the
params of the render that created the setter, so several setters called in one
event handler each start from the same stale `prev` and each issues its own
`navigate()` — the last one wins and the others' deletions are lost. "Clear all"
(`ArtistCharacterFilter`) calls three setters in one click, so it left the URL
filtered even when only one dimension was active. (This was pre-existing, not
introduced by persistence; persistence made it stick, because the next filter
interaction re-wrote the un-cleared URL state back to localStorage.)

Every setter therefore funnels through a single `commitFilter(filter)` that
writes `filter` to localStorage and serializes **all four** filter params in one
non-updater `setSearchParams(next, { replace: true })`, built from
`pendingFilterRef` (the accumulating state) on top of `paramsRef` (the last
known params, so non-filter params survive untouched). Sequential setters in one
event then converge on the correct final URL.

_Note:_ "Clear all" does not call `setCharacterMatch`, so a `match=any` param
can outlive the chips it applied to. It is inert — `filterAndSortArtpieces`
ignores the match mode when no characters are selected.

_Alternative considered:_ a dedicated `clearFilters()` on the hook plus an
`onClearAll` prop on `ArtistCharacterFilter`. Rejected: it fixes only this one
call site, leaves every other multi-setter interaction broken, and changes a
shared component contract used by two other modals.

### Decision: Write to localStorage directly in each setter

When the user changes a filter, the setter already knows the new value and has access to the other dimensions' current values (derived from URL params in the hook's closure). Writing synchronously in the setter avoids a second effect and is consistent with how `gallery.tileSize` is written.

### Decision: Stored shape uses raw IDs, not full objects

`artist` and `characters` are stored as IDs, matching what the URL params hold. The hook already resolves IDs → full objects from React Query cache. This keeps the stored shape stable and independent of display fields (like artist name) that could change.

```ts
interface StoredFilter {
  artistId: string | null
  characterIds: string[]
  matchMode: 'all' | 'any'
  isReference: boolean
}
```

### Decision: Storage key `artpieces.filterState`

Namespaced separately from the existing display-preference keys (`gallery.tileSize`, `gallery.showDetails`) to keep concerns distinct and allow independent clearing.

### Decision: Restoration skips trivially empty stored state

If `artistId` is null, `characterIds` is empty, and `isReference` is false, there is nothing meaningful to restore — the default state is already correct. The restore path skips `setSearchParams` in that case to avoid a no-op URL update and re-render.

## Risks / Trade-offs

**Stale artist/character IDs in localStorage** → If an artist or character that was part of a stored filter is later deleted, the restored URL params will contain orphaned IDs. `useArtpiecesFilter` already handles this gracefully: it looks up IDs from the React Query cache and silently drops any that don't resolve. The stored state is not corrected, but the next filter interaction overwrites it.

**localStorage unavailable (private mode, blocked storage)** → All reads and writes are wrapped in try/catch, consistent with the existing `gallery.tileSize` pattern. The hook degrades silently to the current stateless behavior.
