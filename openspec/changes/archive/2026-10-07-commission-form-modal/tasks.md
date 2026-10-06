## 1. Foundation — Type and Hook Layer

- [x] 1.1 Add `commission_id: string | null` to `ArtpieceSummary` in `frontend/src/features/artpieces/api/useArtpieces.ts`
- [x] 1.2 Create `frontend/src/features/commissions/api/useCommission.ts` — `useQuery` for `GET /commissions/:id`, returning `Commission` with artpieces (id + thumbnail_url)
- [x] 1.3 Create `frontend/src/features/commissions/api/useCreateCommission.ts` — `useMutation` for `POST /commissions`, invalidates commissions query on success
- [x] 1.4 Create `frontend/src/features/commissions/api/useUpdateCommission.ts` — `useMutation` for `PUT /commissions/:id`, invalidates commissions query on success
- [x] 1.5 Create `frontend/src/features/commissions/api/useReplaceArtpieces.ts` — `useMutation` for `PUT /commissions/:id/artpieces`, invalidates commissions query on success
- [x] 1.6 Create `frontend/src/features/commissions/api/useDeleteCommission.ts` — `useMutation` for `DELETE /commissions/:id`, invalidates commissions query on success

## 2. StatusChip Shared Component

- [x] 2.1 Create `frontend/src/features/commissions/components/StatusChip.tsx` — controlled `<StatusChip value onChange hasError? />` using `DropdownMenu`; move `STATUS_OPTIONS` and `STATUS_LABELS` consts into this file and re-export them
- [x] 2.2 Replace the inline `DropdownMenu` status cell in `CommissionsTable` with `<StatusChip value={commission.status} onChange={(s) => patchStatus(commission.id, s)} hasError={hasCellError(commission.id, 'status')} />`
- [x] 2.3 Verify the table's inline autosave behaviour is unchanged after the refactor

## 3. Delete Commission Dialog

- [x] 3.1 Create `frontend/src/features/commissions/components/DeleteCommissionDialog.tsx` — `AlertDialog` with copy "This commission will be permanently deleted. Any artpieces linked to it will be kept intact."; fires `useDeleteCommission` on confirm; shows error toast on failure

## 4. Commission Form Modal

- [x] 4.1 Create `frontend/src/features/commissions/components/CommissionFormModal.tsx` — modal scaffold with props `open`, `onOpenChange`, `mode: 'create' | 'edit'`, `commission?: Commission`
- [x] 4.2 Implement form fields: Title (text input), Artist (ArtistCombobox + + button → ArtistFormModal, same pattern as ArtpieceFormModal), Notes (textarea), Status (StatusChip), Price (numeric input), Paid (Checkbox), Paid Date (Calendar + Popover datepicker), Finish Date (Calendar + Popover datepicker)
- [x] 4.3 Implement artpieces section — collapsible wrapper, strip of attached artpieces (thumbnail + × button to remove), title search input, ArtpiecesFilterPopover (pass fixed `characterMatch='any'` and no-op `onCharacterMatchChange`)
- [x] 4.4 Wire artpiece search: filter `useArtpieces()` data client-side by title search + artist/character filter; exclude items where `commission_id !== null && commission_id !== currentCommissionId`; selecting a result adds it to the strip (deduplicate)
- [x] 4.5 Implement create mode save: POST /commissions with field values + strip artpiece IDs; show error toast on failure; close modal and reset form on success
- [x] 4.6 Implement edit mode pre-populate: on open, call `useCommission(id)` to get artpieces; show loading skeleton in strip while pending; capture initial artpiece set in a ref for change detection
- [x] 4.7 Implement edit mode save: PUT /commissions/:id with field values; if strip changed (compare current set to initial ref), also fire PUT /commissions/:id/artpieces; show error toast on any failure; close modal on success
- [x] 4.8 Reset all form state on modal close (not on unmount only — also on `open` toggling false)

## 5. Table Row Actions

- [x] 5.1 Add an actions `TableHead` column to `CommissionsTable` after the Last Contact / Time Taken column (no label or a visually hidden label for accessibility)
- [x] 5.2 Add `TableCell` per row with a brush icon button (opens edit modal for that row's commission) and a trash icon button (opens delete dialog for that row's commission)
- [x] 5.3 Wire `CommissionFormModal` open state into `CommissionsTable` (or into `CommissionsPage`) — whichever keeps the lifting minimal; pass the selected commission to the modal

## 6. Wire "+ New Commission" Button

- [x] 6.1 In `CommissionsPage` (or wherever the header button lives), replace the no-op onClick with state that opens `CommissionFormModal` in create mode

## 7. Bruno Collection

- [x] 7.1 Verify existing Bruno files cover `POST /commissions`, `PUT /commissions/:id`, `PUT /commissions/:id/artpieces`, `DELETE /commissions/:id` — create any missing files in `collection/commissions/`

## 8. Build and Lint

- [x] 8.1 Run `npm run build` in `frontend/` and fix any TypeScript or build errors
- [x] 8.2 Run `npm run lint` in `frontend/` and fix any lint errors
