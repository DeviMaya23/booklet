# Proposal

## Why

The artpieces gallery filter (artist, characters, match mode, references toggle) is stored in URL params, but those params are stripped whenever the user navigates away via the sidebar. Returning to the gallery always starts with a blank filter, even if the user had carefully composed one before leaving.

## What Changes

- `useArtpiecesFilter` gains localStorage read/write for structural filter state (artist ID, character IDs, match mode, isReference)
- On mount with no URL filter params: the hook restores the last-used structural filter from localStorage by applying it to the URL params
- On mount with URL filter params already present (e.g. arriving from a dashboard saved-filter link): the hook writes those params to localStorage, updating the "last used" state
- Each structural filter setter (artist, characters, match, isReference) also writes the full current filter state to localStorage after updating the URL
- `search` and `sort` are explicitly excluded from persistence — they remain ephemeral
- localStorage key: `artpieces.filterState`; stored as `{ artistId, characterIds, matchMode, isReference }`

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `web-artpieces-gallery`: new requirement — structural filter state persists across navigation, so returning to the gallery via the sidebar restores the last-used artist/character/reference filter

## Impact

- `frontend/src/features/artpieces/hooks/useArtpiecesFilter.ts` — the only file that changes
- No backend changes
- No new dependencies
- `ArtpiecesPage` is unaffected — persistence is fully internal to the hook
