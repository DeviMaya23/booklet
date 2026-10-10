## Why

The current artist list is a plain bullet-point row with a copy-link icon — it predates the multi-link system and doesn't surface links or notes at a glance. The revamp brings it in line with the table-based UI used elsewhere and makes the richer data the system now stores actually visible.

## What Changes

- Replace the flat div/bullet list with a proper `<Table>` layout (Name | Links | Notes | Edit)
- Replace the copy-link button with clickable link chips per artist, showing hostname only (e.g. `bsky.app`), capped at 3 visible with a `+N` overflow badge for the rest
- The primary link chip is visually distinguished (star icon + outlined border) and always rendered first
- Each chip is an anchor that opens the full URL in a new tab; hovering shows a tooltip with the full URL and "Main link · opens in a new tab" (primary) or "Opens in a new tab" (others)
- Artists with no links show a "No links" muted placeholder instead of chips
- The Notes column shows the artist's notes, truncated with ellipsis when too long; artists with no notes show "—"
- The `+N` overflow chip is display-only (not clickable)
- **BREAKING**: The copy-link button is removed

## Capabilities

### New Capabilities

_None — this is a display revamp of an existing capability._

### Modified Capabilities

- `web-artists`: List display requirement changes substantially (table layout, link chips, notes column, no copy-link button); copy-artist-link requirement is removed entirely

## Impact

- `frontend/src/features/artists/components/ArtistsList.tsx` — full rewrite
- `frontend/src/features/artists/components/ArtistsList.test.tsx` — two existing tests obsolete; new tests for chip rendering, primary-first ordering, overflow count, link targets
- No backend changes required
- No new npm dependencies (Table and Tooltip components already exist in `/components/ui/`)
