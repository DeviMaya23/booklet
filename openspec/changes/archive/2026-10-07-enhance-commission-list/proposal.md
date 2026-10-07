## Why

The commissions list is functional but lacks visual hierarchy and key workflow features — specifically, there is no at-a-glance staleness signal for Last Contact, no undo path when accidentally marking a commission as Finished, and no way to see which artpieces belong to a finished commission without opening the edit modal. The artpiece detail view also has no URL, which makes it impossible to link to a specific artpiece from other pages.

## What Changes

- Replace the swap-icon toggle with Active / Finished tab buttons showing live counts
- Table header and cell typography: headers smaller and muted, data cells 14px, only the active sort column shows a chevron
- Last Contact column: header gains an info-icon tooltip explaining the chase-reminder mechanic; each stamp button gains a per-row tooltip; entries older than 14 days turn amber (text + stamp icon)
- **Undo toast**: when a commission is moved to `done` from the Active tab, the row leaves immediately and a Sonner toast appears ("_title_ moved to Finished · Undo") with an action button that PATCHes back to the previous status
- Finished tab: replace the Time Taken column with an Artpieces column showing up to 2 artpiece thumbnails as linked chips; overflow shown as a `+N` chip
- **BREAKING** (routing only): artpiece detail view is now reachable at `/app/artpieces/:id`; double-clicking an artpiece card navigates to that route instead of toggling local state
- Back chevron in artpiece detail changes label from "← Artpieces" to "← Back" and uses `navigate(-1)` so it returns to wherever the user came from (artpieces gallery or commissions list)
- Commissions tab state moves from local `useState` to `?tab=` query param so `navigate(-1)` from artpiece detail restores the correct tab

## Capabilities

### New Capabilities
- `web-artpiece-routing`: URL-based artpiece detail — route `/app/artpieces/:id`, App.tsx change, ArtpiecesPage reads `useParams`, navigate-based open/close

### Modified Capabilities
- `web-commissions-list`: tabs with counts, typography, last-contact amber/tooltip, undo toast, artpiece thumbnail column, tab state in query param
- `web-artpiece-detail`: route-based open (no longer local state), "← Back" label, `navigate(-1)` on close

## Impact

- `frontend/src/App.tsx` — new route `/app/artpieces/:id`
- `frontend/src/pages/ArtpiecesPage.tsx` — `useState` → `useParams`, open/close via `navigate`
- `frontend/src/pages/CommissionsPage.tsx` — tabs UI, `useState` → `useSearchParams`, undo toast wiring, `usePatchCommission` instance for undo
- `frontend/src/features/commissions/components/CommissionsTable.tsx` — typography, sort chevron logic, last-contact amber + tooltips, artpiece thumbnail column, `onCommissionDone` callback prop
- `frontend/src/features/commissions/components/ArtpieceThumbnailStack.tsx` — new component
- `frontend/src/features/artpieces/components/ArtpieceDetailView.tsx` — label + close behaviour
- `frontend/src/features/commissions/api/useCommissions.ts` — add `artpieces` field to `Commission` type (backend already returns it)
- No backend changes required
