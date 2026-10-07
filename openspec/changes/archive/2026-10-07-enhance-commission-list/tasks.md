## 1. Artpiece URL Routing

- [x] 1.1 Add route `/app/artpieces/:id` in `App.tsx` alongside the existing `/app/artpieces` route
- [x] 1.2 Update `ArtpiecesPage` to read artpiece ID from `useParams()` instead of local `useState`; replace `setSelectedArtpieceId(id)` with `navigate('/app/artpieces/' + id)` and `setSelectedArtpieceId(null)` with `navigate(-1)` for close and `navigate('/app/artpieces')` for onDeleted
- [x] 1.3 Update `ArtpieceDetailView`: change back chevron label from "← Artpieces" to "← Back" (all occurrences); `onClose` remains a prop — no internal change needed

## 2. Commission Type Update

- [x] 2.1 Add `artpieces: { id: string; thumbnail_url: string | null }[]` field to the `Commission` interface in `useCommissions.ts`

## 3. Tab UI and URL State

- [x] 3.1 Replace `useState<View>` in `CommissionsPage` with `useSearchParams`; reading `?tab=` (default `'active'`), switching tabs via `setSearchParams`
- [x] 3.2 Replace the swap-icon toggle heading with Active / Finished tab buttons; each button shows a count badge computed from the full (unfiltered) commissions array; keep "+ New Commission" button on the right

## 4. Table Typography and Sort Chevron

- [x] 4.1 Update `CommissionsTable` header cells: apply smaller muted text style (e.g. `text-xs text-muted-foreground font-medium`)
- [x] 4.2 Update `CommissionsTable` data cells: apply `text-sm` (14px) uniformly
- [x] 4.3 Update `SortIcon`: render nothing (return `null`) when `column !== sortKey`; only show chevron on the active sort column

## 5. Last Contact Column Enhancements

- [x] 5.1 Add amber styling to the Last Contact cell: when `last_contacted_at` is null or older than 14 days, apply amber text colour to the relative-time string and stamp icon
- [x] 5.2 Add a tooltip wrapping the column header's info icon; tooltip text: "A manual reminder. Press ↻ whenever you hear from the artist. It turns amber when it's been a while."
- [x] 5.3 Add a tooltip wrapping the stamp button on each row; tooltip text: "Mark as contacted today"

## 6. Undo Toast

- [x] 6.1 Add `onCommissionDone?: (id: string, prevStatus: string, title: string | null) => void` prop to `CommissionsTable`
- [x] 6.2 In `CommissionsTable`, capture `commission.status` before calling `patchStatus`; after a successful patch to `'done'`, call `onCommissionDone` with the id, previous status, and title
- [x] 6.3 In `CommissionsPage`, add a `usePatchCommission` instance for undo; implement `handleCommissionDone` that shows a Sonner toast with the commission title and an Undo action button that PATCHes back to the previous status
- [x] 6.4 Pass `onCommissionDone={handleCommissionDone}` to `CommissionsTable`

## 7. Artpiece Thumbnails in Finished Tab

- [x] 7.1 Create `ArtpieceThumbnailStack` component in `frontend/src/features/commissions/components/`; accepts `artpieces: { id: string; thumbnail_url: string | null }[]`; renders up to 2 thumbnails as `<Link to="/app/artpieces/:id">` elements; shows `+N` chip for overflow; shows `—` when empty
- [x] 7.2 Replace the Finished tab's Time Taken column in `CommissionsTable` with an Artpieces column that renders `<ArtpieceThumbnailStack artpieces={commission.artpieces} />`

## 8. Quality

- [x] 8.1 Run `npm run build` and fix any type errors
- [x] 8.2 Run `npm run lint` and fix any lint issues
