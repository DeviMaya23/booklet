## Context

The artist list lives entirely in `ArtistsList.tsx` and its parent `ArtistsPage.tsx`. The data model (`Artist`, `ArtistLink`) already carries everything needed — `url`, `is_primary`, `notes` — so this is a pure display change. No backend or API contract changes are involved.

The UI kit already has `<Table>` (`/components/ui/table.tsx`) and `<Tooltip>` (`/components/ui/tooltip.tsx`). There is no `Badge` component; chips will be styled inline.

## Goals / Non-Goals

**Goals:**
- Replace the bullet list with a `<Table>` layout
- Render link chips with hostname-only label, primary-first sort, 3-chip cap + `+N` overflow
- Add a Notes column with ellipsis truncation
- Replace copy-link interaction with chip-as-link (opens in new tab)

**Non-Goals:**
- Adding a `Badge` or `Chip` shared component — chips are local to `ArtistsList`
- Any backend changes
- Changing the search or modal behaviour in `ArtistsPage.tsx`

## Decisions

### D1: Sort primary link first on the frontend, not the backend

The backend `Preload("Links")` returns links in insertion order (undefined from the frontend's perspective). Three backend consumers (`dashboard_handler`, `commission_handler`, `artpiece_handler`) iterate links looking for `IsPrimary` — none depend on order. A backend `ORDER BY` would fix order globally but is unnecessary scope.

**Decision:** Sort in `ArtistsList` before rendering: `[...links].sort((a, b) => (b.is_primary ? 1 : 0) - (a.is_primary ? 1 : 0))`. Zero backend work; self-contained in the component.

### D2: Hostname extraction via `new URL(url).hostname`

All stored URLs are normalised at entry time via `normalizeUrl()` (prepends `https://` if no scheme), so `new URL()` is safe. `hostname` strips protocol, path, and query — matching the mockup labels exactly (`bsky.app`, not `https://bsky.app/@chev`).

### D3: Notes column truncation via `max-w-0 w-full truncate` pattern

`<TableCell>` has `whitespace-nowrap` built in. Adding `max-w-0 w-full truncate` on the Notes cell makes it shrink to available width and clip with `…`. This avoids needing `table-layout: fixed` or explicit column widths.

### D4: Tooltip content differs between primary and secondary links

Primary chip tooltip: full URL on line 1, "Main link · opens in a new tab" on line 2.  
Other chip tooltip: full URL on line 1, "Opens in a new tab" on line 2.  
The existing `<TooltipContent>` accepts arbitrary children, so a two-line flex column works without modifications.

## Risks / Trade-offs

- **Hostname collision**: Two different links could share the same hostname (e.g. two `twitter.com` links). The chip label would look identical. Mitigation: none needed — the tooltip shows the full URL on hover, and the form prevents duplicate URLs at edit time.
- **Long hostnames**: An unusually long hostname (e.g. `very-long-subdomain.example.com`) could overflow the chip visually. Mitigation: cap chip label with `max-w-[10rem] truncate` if it becomes a real issue; not anticipated for normal social/portfolio URLs.
- **`new URL()` on malformed stored data**: Theoretically possible if old data pre-dates the normaliser. Mitigation: wrap in try/catch and fall back to displaying the raw URL as the label.
