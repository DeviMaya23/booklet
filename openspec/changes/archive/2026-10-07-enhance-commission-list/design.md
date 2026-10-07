## Context

The commissions list is a single page at `/app/commissions` that currently uses local `useState` to toggle between Active and Finished views. The artpiece detail view is rendered in-page with no URL — `ArtpiecesPage` holds `selectedArtpieceId` state and swaps the grid for the detail component. This means there is no route to a specific artpiece, making cross-page linking impossible.

The backend commission list endpoint (`GET /commissions`) returns `artpieces: [{id, thumbnail_url}]` per commission. However, the original `List` repository query was missing `Preload("Artpieces.CoverFile")` — a backend change was required to add it so the list response actually includes artpiece data with presigned thumbnail URLs.

## Goals / Non-Goals

**Goals:**
- Replace toggle with tab UI and preserve tab selection in the URL
- Add amber staleness signal and tooltips to the Last Contact column
- Add undo path for accidental status→done changes
- Show linked artpiece thumbnails in the Finished tab
- Make artpiece detail reachable via `/app/artpieces/:id`

**Non-Goals:**
- Customisable staleness threshold (threshold is 14 days, hardcoded; future work)
- Commission-level URL routing (`/app/commissions/:id` not introduced)
- Artpiece gallery filtering from a commission context
- Overlapping/stacked thumbnail effect (max 2 + chip is the accepted limit)

## Decisions

### Undo: immediate PATCH + compensating PATCH
Two options considered: (a) delay the PATCH until the toast dismisses (Gmail-style), (b) PATCH immediately then offer a compensating PATCH on Undo.

**Decision: option (b) — immediate PATCH.**

Rationale: the inline autosave model already PATCHes on every cell change. Making status-to-done behave differently (delayed) would be an inconsistency the user doesn't see but the code does. Option (b) keeps the save model uniform. The Sonner toast's default duration (4 s) is enough for accidental-click recovery; intentional moves rarely get undone.

Implementation: `CommissionsTable` gains an `onCommissionDone(id, prevStatus, title)` callback prop. The table captures `commission.status` before calling `patchStatus`, then fires the callback on success. `CommissionsPage` owns a `usePatchCommission` instance for the undo PATCH and shows the toast.

### Artpiece routing: same page component, reads useParams
Options: (a) same `ArtpiecesPage` reads `useParams`, (b) new `ArtpieceDetailPage` wrapper.

**Decision: option (a).**

Rationale: `ArtpieceDetailView` already accepts `artpieceId` as a prop and `ArtpiecesPage` already has the conditional render (`if (id) return <ArtpieceDetailView>`). Swapping `useState` for `useParams` is a one-line change in the page; the view component is untouched.

Route added in `App.tsx`:
```
/app/artpieces       → ArtpiecesPage (existing)
/app/artpieces/:id   → ArtpiecesPage (reads :id param)
```

### Back navigation: navigate(-1)
Options: (a) `navigate('/app/artpieces')` (always back to gallery), (b) `navigate(-1)` (browser history back).

**Decision: option (b) — `navigate(-1)`.**

Rationale: With (a), navigating from CommissionsPage → artpiece detail → back would land on the artpieces gallery, not commissions — confusing. With (b), the browser restores the exact previous URL including the `?tab=` param, so commissions users return to the correct tab. The only edge case is a direct URL open with no history; in that case `navigate(-1)` exits the app gracefully.

The "← Artpieces" hardcoded label becomes "← Back" to remain accurate from any entry point.

### Tab state: ?tab= query param
The tab toggle moves from `useState<'active'|'done'>` to `useSearchParams`. Switching tabs calls `setSearchParams({ tab: 'finished' })`. Default (no param) is Active.

This means every navigation that pushes `/app/commissions?tab=finished` onto history preserves the tab on `navigate(-1)`.

### Backend: add Preload("Artpieces.CoverFile") to List query
The `commissionRepository.List` query was missing this preload. Without it, the API returned commissions with an empty artpieces slice. Added alongside the frontend type change.

### Commission type: add artpieces field on the frontend
The backend already returns `artpieces: [{id, thumbnail_url}]` in the list response. The frontend `Commission` type simply needs the field declared — no fetch changes, no backend changes.

## Risks / Trade-offs

- **Undo PATCH failure**: if the compensating PATCH fails, the commission stays in Finished with no visible error (the toast is already gone). Mitigation: log the failure; tolerable for this use case since the user can always re-edit.
- **navigate(-1) with no history**: on direct URL access to `/app/artpieces/:id`, clicking back exits the app. Acceptable for an internal tool; a future improvement would check `window.history.length > 1` and fall back to `/app/artpieces`.
- **Presigned URL expiry for thumbnails**: artpiece thumbnails in the Finished tab use presigned URLs from the commission list response. URLs expire after the standard TTL; stale commissions sessions may show broken thumbnails. React Query's refetch-on-focus mitigates this.
- **14-day threshold hardcoded**: the amber signal cannot be adjusted per user. Accepted as known limitation; the design leaves room for a user preference in a future ticket.
